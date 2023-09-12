package booking

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

func GetSeatData(c *gin.Context) {
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

  // get flight id from params
  flightID := c.Param("flightID")
  if flightID == "" {
    helpers.JSONResponse(c, 400, false, "Flight ID is required", nil)
    return
  }

  // check if flight exist
  var flight models.Flight
  if err := db.Model(&models.Flight{}).Where("id = ?", flightID).First(&flight).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid flight", nil)
    return
  }

  // get seat data based on flight id
  var seatData []string
  db.Table("booking_details").Select("seat_number").Joins("JOIN bookings ON bookings.id = booking_details.booking_id").Where("bookings.flight_id = ?", flightID).Scan(&seatData)

  // check if seat data is empty
  if len(seatData) == 0 {
    helpers.JSONResponse(c, 200, true, "No reserved seats", nil)
    return
  }

  // return seat data
  helpers.JSONResponse(c, 200, true, "Success", seatData)
}
