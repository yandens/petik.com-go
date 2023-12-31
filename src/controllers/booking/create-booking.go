package booking

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "math/rand"
  "time"
)

type CreateBookingInput struct {
  FlightID       uint   `json:"flightId" binding:"required"`
  FlightClass    string `json:"flightClass" binding:"required"`
  TotalPassenger int    `json:"totalPassenger" binding:"required"`
  Detail         []struct {
    PassengerName string `json:"passengerName" binding:"required"`
    PassengerAge  int    `json:"passengerAge" binding:"required"`
    NIK           string `json:"nik" binding:"required"`
    SeatNumber    string `json:"seatNumber" binding:"required"`
  } `json:"detail" binding:"required"`
}

func CreateBooking(c *gin.Context) {
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
  var input CreateBookingInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, err.Error(), nil)
    return
  }

  // get user bio for email notification
  var userBio models.UserBio
  if err := db.Preload("User").Model(&models.UserBio{}).Where("user_id = ?", userID).First(&userBio).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // check if flight exist
  var flight models.Flight
  if err := db.Model(&models.Flight{}).Where("id = ?", input.FlightID).First(&flight).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid flight", nil)
    return
  }

  // check if flight class exist
  if input.FlightClass != "Economy" && input.FlightClass != "Business" && input.FlightClass != "First" {
    helpers.JSONResponse(c, 400, false, "Invalid flight class", nil)
    return
  }

  // set total price
  var totalPrice, randomNumber int

  rand.Seed(time.Now().UnixNano())
  if input.FlightClass == "Economy" {
    randomNumber = rand.Intn(321-64) + 64
    totalPrice = randomNumber * input.TotalPassenger
  } else if input.FlightClass == "Business" {
    randomNumber = rand.Intn(642-321) + 321
    totalPrice = randomNumber * input.TotalPassenger
  } else {
    randomNumber = rand.Intn(1606-642) + 642
    totalPrice = randomNumber * input.TotalPassenger
  }

  // check if total passenger is equal to detail length
  if input.TotalPassenger != len(input.Detail) {
    helpers.JSONResponse(c, 400, false, "Invalid total passenger", nil)
    return
  }

  // get seat number from user input
  var seatNumbers []string
  for _, detail := range input.Detail {
    seatNumbers = append(seatNumbers, detail.SeatNumber)
  }

  // check if seat number is unique and valid
  if !helpers.IsUniqueSeatNumber(seatNumbers) || !helpers.IsSeatNumberValid(seatNumbers) {
    helpers.JSONResponse(c, 400, false, "Invalid seat number", nil)
    return
  }

  // check if seat number is available
  if !helpers.IsSeatNumberAvailable(seatNumbers, flight.ID, db) {
    helpers.JSONResponse(c, 400, false, "Seat number is reserved", nil)
    return
  }

  // create booking
  booking := models.Booking{
    UserID:         userID,
    FlightID:       flight.ID,
    FlightClass:    input.FlightClass,
    TotalPassenger: input.TotalPassenger,
    TotalPrice:     float64(totalPrice),
    Status:         "pending",
  }

  if err := db.Create(&booking).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // create booking detail
  for _, detail := range input.Detail {
    bookingDetail := models.BookingDetail{
      BookingID:     booking.ID,
      PassengerName: detail.PassengerName,
      PassengerAge:  detail.PassengerAge,
      NIK:           detail.NIK,
      SeatNumber:    detail.SeatNumber,
    }

    if err := db.Create(&bookingDetail).Error; err != nil {
      helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
      return
    }
  }

  // create notification
  notification := models.Notification{
    UserID:  userID,
    Title:   "Booking",
    Message: "Your booking has been created",
    IsRead:  false,
  }

  // save notification to database
  if err := db.Create(&notification).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Something went wrong", nil)
    return
  }

  // send email notification
  helpers.SendEmailBookingConfirmation(userBio, flight, booking, "Booking Confirmation")

  helpers.JSONResponse(c, 200, true, "Success create booking", gin.H{
    "bookingId":      booking.ID,
    "userId":         booking.UserID,
    "flightId":       booking.FlightID,
    "flightClass":    booking.FlightClass,
    "totalPassenger": booking.TotalPassenger,
    "totalPrice":     booking.TotalPrice,
    "status":         booking.Status,
    "detail":         input.Detail,
  })
}
