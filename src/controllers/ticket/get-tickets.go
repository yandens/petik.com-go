package ticket

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "time"
)

func GetTicket(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id from middleware
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get booking
  var booking models.Booking
  if err := db.Joins("Flight").Joins("BookingDetail").Model(&models.Booking{}).Where("user_id = ? AND status = paid", id).First(&booking).Error; err != nil {
    helpers.JSONResponse(c, 404, false, "Booking not found", nil)
    return
  }

  // loop through booking details
  var tickets []TicketData
  for _, bookingDetail := range booking.BookingDetail {
    tickets = append(tickets, TicketData{
      UserID:          booking.UserID,
      BookingID:       booking.ID,
      BookingDetailID: bookingDetail.ID,
      Name:            bookingDetail.PassengerName,
      Airline:         booking.Flight.Airline,
      AirlineLogo:     booking.Flight.AirlineLogo,
      Origin:          booking.Flight.Origin,
      Destination:     booking.Flight.Destination,
      DepartureDate:   booking.Flight.Departure.Format("02-01-2006"),
      DepartureTime:   booking.Flight.Departure.Format("15:04"),
      BoardingTime:    booking.Flight.Departure.Add(-time.Hour * 1).Format("15:04"),
      Seat:            bookingDetail.SeatNumber,
      Class:           booking.FlightClass,
      QRCode:          bookingDetail.QRCode,
    })
  }

  // return response
  helpers.JSONResponse(c, 200, true, "Success get user tickets", tickets)
}
