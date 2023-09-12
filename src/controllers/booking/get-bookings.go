package booking

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

func GetUserBookings(c *gin.Context) {
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

  // get user bookings
  var bookings []models.Booking
  if err := db.Model(&models.Booking{}).Where("user_id = ?", id).Find(&bookings).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Success", bookings)
}

func GetTotalBooking(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // get id
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get total booking by user id
  var totalBooking int64
  if err := db.Model(&models.Booking{}).Where("user_id = ?", id).Count(&totalBooking).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Success", gin.H{
    "totalBooking": totalBooking,
    "userId":       id,
  })
}
