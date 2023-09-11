package flight

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/controllers/airport"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "time"
)

type UpdateFlightInput struct {
  Airline     string `json:"airline"`
  Origin      string `json:"origin"`
  Destination string `json:"destination"`
  Departure   string `json:"departure"`
  Arrival     string `json:"arrival"`
}

func UpdateFlight(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "failed to connect to database", nil)
    return
  }

  // get id from params
  id := c.Param("id")
  if id == "" {
    utils.JSONResponse(c, 400, false, "flight id is required", nil)
    return
  }

  // check if flight exists
  var flightExist models.Flight
  if err := db.Where("id = ?", id).First(&flightExist).Error; err != nil {
    utils.JSONResponse(c, 404, false, "flight not found", nil)
    return
  }

  // get input
  var input UpdateFlightInput
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "invalid input", nil)
    return
  }

  // parse departure and arrival time
  layout := "2006-01-02 15:04"
  departure, _ := time.Parse(layout, input.Departure)
  arrival, _ := time.Parse(layout, input.Arrival)

  // validate departure and arrival time
  if arrival.Before(departure) || departure.After(arrival) {
    utils.JSONResponse(c, 400, false, "invalid departure and arrival time", nil)
    return
  }

  // get airline logo based on airline name
  var airlineLogo string
  switch input.Airline {
  case "Garuda Indonesia":
    airlineLogo = "https://bit.ly/3BOGwXN"
  case "Sriwijaya Air":
    airlineLogo = "https://bit.ly/3FDlHzT"
  case "Batik Air":
    airlineLogo = "http://bit.ly/3Ytqr3w"
  case "Lion Air":
    airlineLogo = "https://bit.ly/3WbvdBj"
  case "Air Asia":
    airlineLogo = "https://bit.ly/3PF90Jm"
  case "Citilink":
    airlineLogo = "https://bit.ly/3v6DYAJ"
  case "Wings Air":
    airlineLogo = "http://bit.ly/3G4OO0p"
  default:
    airlineLogo = ""
  }

  // get origin city
  origin, err := airport.GetAirportByIATA(input.Origin)
  if err != nil {
    utils.JSONResponse(c, 400, false, "invalid input", nil)
    return
  }

  // get destination city
  destination, err := airport.GetAirportByIATA(input.Destination)
  if err != nil {
    utils.JSONResponse(c, 400, false, "invalid input", nil)
    return
  }

  // update flight
  flight := models.Flight{
    Airline:         input.Airline,
    AirlineLogo:     airlineLogo,
    Origin:          input.Origin,
    OriginCity:      origin["city"],
    Destination:     input.Destination,
    DestinationCity: destination["city"],
    Departure:       departure,
    Arrival:         arrival,
  }
  if err := db.Model(&models.Flight{}).Where("id = ?", id).Updates(flight).Error; err != nil {
    utils.JSONResponse(c, 400, false, "failed to update flight", nil)
    return
  }

  // return response
  utils.JSONResponse(c, 200, true, "flight updated successfully", nil)
}
