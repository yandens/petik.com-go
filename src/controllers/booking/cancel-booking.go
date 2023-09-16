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

  // get user input
  var input CancelBookingInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // check if booking exist
  var booking models.Booking
  if err := db.Model(&models.Booking{}).Where("id = ?", input.BookingID).First(&booking).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid booking", nil)
    return
  }

  // check if booking belong to user
  if booking.UserID != userID {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // check if booking status is pending to be canceled
  if booking.Status == "paid" {
    helpers.JSONResponse(c, 400, false, "Invalid booking status", nil)
    return
  }

  // update booking status to canceled
  if err := db.Model(&models.Booking{}).Where("id = ?", input.BookingID).Update("status", "canceled").Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Booking canceled", nil)
}
