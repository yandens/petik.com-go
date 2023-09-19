package booking

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type CancelBookingInput struct {
  BookingID uint `json:"bookingId" binding:"required"`
}

func CancelBooking(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // get user id from middleware
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get user input
  var input CancelBookingInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // check if booking exist
  var booking models.Booking
  if err := db.Model(&models.Booking{}).Where("id = ? AND user_id = ?", input.BookingID, id).First(&booking).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid booking", nil)
    return
  }

  // check if booking status is pending to be canceled
  if booking.Status != "pending" {
    helpers.JSONResponse(c, 400, false, "Invalid booking status", nil)
    return
  }

  // update booking status to canceled
  if err := db.Model(&models.Booking{}).Where("id = ?", input.BookingID).Update("status", "canceled").Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // create notification
  notification := models.Notification{
    UserID:  booking.UserID,
    Title:   "Booking canceled",
    Message: "Your booking has been canceled",
  }

  // save notification to database
  if err := db.Create(&notification).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Booking canceled", nil)
}
