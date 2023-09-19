package payment

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
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

  // type assertion convert interface{} to uint
  var userID uint
  switch id := id.(type) {
  case float64:
    userID = uint(id)
  case uint:
    userID = id
  default:
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
  }

  // get input from user
  var input PaymentInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid Input", nil)
    return
  }

  // check if booking exists
  var booking models.Booking
  if err := db.Where("id = ? AND user_id = ?", input.BookingID, userID).First(&booking).Error; err != nil {
    helpers.JSONResponse(c, 404, false, "Booking not found", nil)
    return
  }

  // check if booking status is pending
  if booking.Status != "Pending" {
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

  // return response
  helpers.JSONResponse(c, 200, true, "Payment created successfully", payment)
}
