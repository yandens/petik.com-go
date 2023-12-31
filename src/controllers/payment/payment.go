package payment

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/controllers/ticket"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type PaymentInput struct {
  BookingID     uint    `json:"bookingId" binding:"required"`
  PaymentMethod string  `json:"paymentMethod" binding:"required"`
  TotalPrice    float64 `json:"totalPrice" binding:"required"`
}

func CreatePayment(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to connect to database", nil)
    return
  }

  // get user id from middleware
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get input from user
  var input PaymentInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid Input", nil)
    return
  }

  // get user bio data
  var userBio models.UserBio
  if err := db.Preload("User").Model(&models.UserBio{}).Where("user_id = ?", id).First(&userBio).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to get user data", nil)
    return
  }

  // check if booking exists
  var booking models.Booking
  if err := db.Preload("Flight").Preload("BookingDetail").Model(&models.Booking{}).Where("bookings.id = ? AND user_id = ?", input.BookingID, id).First(&booking).Error; err != nil {
    helpers.JSONResponse(c, 404, false, err.Error(), nil)
    return
  }

  // check if booking status is pending
  if booking.Status != "pending" {
    helpers.JSONResponse(c, 400, false, "Booking status is not pending", nil)
    return
  }

  // check if payment method is valid
  if input.PaymentMethod != "virtual account" && input.PaymentMethod != "bank transfer" {
    helpers.JSONResponse(c, 400, false, "Invalid payment method", nil)
    return
  }

  // check if total price is valid
  if input.TotalPrice != booking.TotalPrice {
    helpers.JSONResponse(c, 400, false, "Invalid total price", nil)
    return
  }

  // create payment
  payment := models.Payment{
    BookingID:     input.BookingID,
    PaymentMethod: input.PaymentMethod,
    TotalPrice:    input.TotalPrice,
  }

  // save payment to database
  if err := db.Create(&payment).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to create payment", nil)
    return
  }

  // update booking status
  if err := db.Model(&models.Booking{}).Where("id = ?", input.BookingID).Update("status", "paid").Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to update booking status", nil)
    return
  }

  // create ticket
  if err := ticket.CreateTicket(c, db, booking); err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to create ticket", nil)
    return
  }

  // create notification
  notification := models.Notification{
    UserID:  booking.UserID,
    Title:   "Payment",
    Message: "Your payment has been created",
    IsRead:  false,
  }

  // save notification to database
  if err := db.Create(&notification).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to create notification", nil)
    return
  }

  // send notification to user
  helpers.SendEmailPaymentConfirmation(userBio, payment, "Payment Confirmation")

  // return response
  helpers.JSONResponse(c, 200, true, "Payment created successfully", gin.H{
    "paymentId":  payment.ID,
    "bookingId":  payment.BookingID,
    "method":     payment.PaymentMethod,
    "totalPrice": payment.TotalPrice,
    "date":       payment.CreatedAt.Format("2006-01-02 15:04:05"),
  })
}
