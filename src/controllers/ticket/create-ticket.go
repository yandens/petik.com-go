package ticket

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "gorm.io/gorm"
  "time"
)

type TicketData struct {
  UserID          uint
  BookingID       uint
  BookingDetailID uint
  Name            string
  Airline         string
  AirlineLogo     string
  Origin          string
  Destination     string
  DepartureDate   string
  DepartureTime   string
  BoardingTime    string
  Seat            string
  Class           string
  QRCode          string
}

func CreateTicket(c *gin.Context, db *gorm.DB, booking models.Booking) error {
  // loop through booking details to create ticket data
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
    })
  }

  // loop through tickets to generate qr code
  for _, ticket := range tickets {
    // generate qr code
    qrCode, err := helpers.GenerateQRCode(ticket)
    if err != nil {
      return err
    }

    // upload qr code to imagekit
    url, err := helpers.UploadToImagekit(c, qrCode, "qr-code", "/qr-codes")
    if err != nil {
      return err
    }

    // update booking detail
    if err := db.Model(&models.BookingDetail{}).Where("id = ?", ticket.BookingDetailID).Update("qr_code", url).Error; err != nil {
      return err
    }
  }

  return nil
}
