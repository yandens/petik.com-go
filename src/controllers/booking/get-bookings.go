package booking

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type BookingResponse struct {
  ID             uint    `json:"id"`
  UserID         uint    `json:"user_id"`
  FlightID       uint    `json:"flight_id"`
  FlightClass    string  `json:"flight_class"`
  TotalPassenger int     `json:"total_passenger"`
  TotalPrice     float64 `json:"total_price"`
  Status         string  `json:"status"`
  Date           string  `json:"date"`
}

func GetBookings(c *gin.Context) {
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

  // convert bookings to response
  var response []BookingResponse
  for _, booking := range bookings {
    response = append(response, BookingResponse{
      ID:             booking.ID,
      UserID:         booking.UserID,
      FlightID:       booking.FlightID,
      FlightClass:    booking.FlightClass,
      TotalPassenger: booking.TotalPassenger,
      TotalPrice:     booking.TotalPrice,
      Status:         booking.Status,
      Date:           booking.CreatedAt.Format("2006-01-02 15:04:05"),
    })
  }

  helpers.JSONResponse(c, 200, true, "Success", response)
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
