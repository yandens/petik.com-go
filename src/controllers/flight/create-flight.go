package flight

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/controllers/airport"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "io/ioutil"
  "net/http"
  "time"
)

type FlightData struct {
  Hex          string  `json:"hex"`
  RegNumber    string  `json:"reg_number"`
  Flag         string  `json:"flag"`
  Lat          float64 `json:"lat"`
  Lng          float64 `json:"lng"`
  Alt          int     `json:"alt"`
  Dir          float64 `json:"dir"`
  Speed        int     `json:"speed"`
  VSpeed       float64 `json:"v_speed"`
  Squawk       string  `json:"squawk"`
  FlightNumber string  `json:"flight_number"`
  FlightICAO   string  `json:"flight_icao"`
  FlightIATA   string  `json:"flight_iata"`
  DepICAO      string  `json:"dep_icao"`
  DepIATA      string  `json:"dep_iata"`
  ArrICAO      string  `json:"arr_icao"`
  ArrIATA      string  `json:"arr_iata"`
  AirlineICAO  string  `json:"airline_icao"`
  AirlineIATA  string  `json:"airline_iata"`
  AircraftICAO string  `json:"aircraft_icao"`
  Updated      int     `json:"updated"`
  Status       string  `json:"status"`
}

type APIResponse struct {
  Lang     string `json:"lang"`
  Currency string `json:"currency"`
  Time     int    `json:"time"`
  ID       string `json:"id"`
  Server   string `json:"server"`
  Host     string `json:"host"`
  PID      int    `json:"pid"`
  Key      struct {
    ID             int    `json:"id"`
    APIKey         string `json:"api_key"`
    Type           string `json:"type"`
    Expired        string `json:"expired"`
    Registered     string `json:"registered"`
    LimitsByHour   int    `json:"limits_by_hour"`
    LimitsByMinute int    `json:"limits_by_minute"`
    LimitsByMonth  int    `json:"limits_by_month"`
    LimitsTotal    int    `json:"limits_total"`
  } `json:"key"`
  Params struct {
    Lang string `json:"lang"`
  } `json:"params"`
  Version int    `json:"version"`
  Method  string `json:"method"`
  Client  struct {
    IP  string `json:"ip"`
    Geo struct {
      CountryCode string  `json:"country_code"`
      Country     string  `json:"country"`
      Continent   string  `json:"continent"`
      City        string  `json:"city"`
      Lat         float64 `json:"lat"`
      Lng         float64 `json:"lng"`
      Timezone    string  `json:"timezone"`
    } `json:"geo"`
    Connection struct {
      Type    string `json:"type"`
      ISPCode int    `json:"isp_code"`
      ISPName string `json:"isp_name"`
    } `json:"connection"`
    Device struct{} `json:"device"`
    Agent  struct{} `json:"agent"`
    Karma  struct {
      IsBlocked bool `json:"is_blocked"`
      IsCrawler bool `json:"is_crawler"`
      IsBot     bool `json:"is_bot"`
      IsFriend  bool `json:"is_friend"`
      IsRegular bool `json:"is_regular"`
    } `json:"karma"`
  } `json:"client"`
  Response []FlightData `json:"response"`
}

type FlightInput struct {
  Airline     string `json:"airline"`
  Origin      string `json:"origin"`
  Destination string `json:"destination"`
  Departure   string `json:"departure"`
  Arrival     string `json:"arrival"`
}

func CreateFlightSeeder() {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    fmt.Println(err)
    return
  }

  // define url
  url := fmt.Sprintf("https://airlabs.co/api/v9/flights?api_key=%s", configs.GetEnv("AIRLABS_API_KEY"))

  // call third party api to get flight schedule
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer resp.Body.Close()

  // read response
  respData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
    return
  }

  // parse response to readable json
  var data APIResponse
  if err := json.Unmarshal(respData, &data); err != nil {
    fmt.Println(err)
    return
  }

  // filter flight data
  var flightData []FlightData
  for _, flight := range data.Response {
    if flight.AirlineIATA == "GA" || flight.AirlineIATA == "SJ" || flight.AirlineIATA == "ID" || flight.AirlineIATA == "JT" || flight.AirlineIATA == "AK" || flight.AirlineIATA == "QG" || flight.AirlineIATA == "IW" {
      if flight.FlightIATA != "" && flight.DepIATA != "" && flight.ArrIATA != "" {
        flightData = append(flightData, flight)
      }
    }
  }

  // loop flight data to save into database
  for _, flight := range flightData {
    var airlineLogo string
    var airlineName string

    // determine airline logo based on airline iata
    switch flight.AirlineIATA {
    case "GA":
      airlineLogo = "https://bit.ly/3BOGwXN"
      airlineName = "Garuda Indonesia"
    case "SJ":
      airlineLogo = "https://bit.ly/3FDlHzT"
      airlineName = "Sriwijaya Air"
    case "ID":
      airlineLogo = "http://bit.ly/3Ytqr3w"
      airlineName = "Batik Air"
    case "JT":
      airlineLogo = "https://bit.ly/3WbvdBj"
      airlineName = "Lion Air"
    case "AK":
      airlineLogo = "https://bit.ly/3PF90Jm"
      airlineName = "Air Asia"
    case "QG":
      airlineLogo = "https://bit.ly/3v6DYAJ"
      airlineName = "Citilink"
    case "IW":
      airlineLogo = "http://bit.ly/3G4OO0p"
      airlineName = "Wings Air"
    default:
      airlineLogo = ""
      airlineName = ""
    }

    // get origin city
    origin, err := airport.GetAirportByIATA(flight.DepIATA)
    if err != nil {
      fmt.Println(err)
      return
    }

    // get destination city
    destination, err := airport.GetAirportByIATA(flight.ArrIATA)
    if err != nil {
      fmt.Println(err)
      return
    }

    // get departure and arrival time
    departureTime, arrivalTime := utils.SetTime()

    // save to database
    flight := models.Flight{
      Airline:         airlineName,
      AirlineLogo:     airlineLogo,
      Origin:          flight.DepIATA,
      OriginCity:      origin["city"],
      Destination:     flight.ArrIATA,
      DestinationCity: destination["city"],
      Departure:       departureTime,
      Arrival:         arrivalTime,
    }
    err = db.Create(&flight).Error
    if err != nil {
      fmt.Println(err)
      return
    }
  }
}

func CreateFlightAdmin(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "failed to connect to database", nil)
  }

  // get input
  var input FlightInput
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "invalid input", nil)
    return
  }

  // parse departure and arrival time
  layout := "2006-01-02 15:04"
  departure, _ := time.Parse(layout, input.Departure)
  arrival, _ := time.Parse(layout, input.Arrival)

  // validate departure and arrival time
  if arrival.Before(departure) {
    utils.JSONResponse(c, 400, false, "invalid departure and arrival time", nil)
    return
  }

  // check if flight already exist
  isExist := db.Model(&models.Flight{}).Where("airline = ? AND origin = ? AND destination = ? AND departure = ? AND arrival = ?", input.Airline, input.Origin, input.Destination, input.Departure, input.Arrival).Take(&models.Flight{}).RowsAffected
  if isExist == 1 {
    utils.JSONResponse(c, 400, false, "flight already exist", nil)
    return
  }

  // set airline logo based on airline name
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

  // save to database
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
  err = db.Create(&flight).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "failed to save flight", flight)
    return
  }
}
