package flight

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "math"
  "strconv"
)

func GetFlights(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "failed to connect to database", nil)
    return
  }

  // get query params
  page, _ := strconv.Atoi(c.Query("page"))
  limit, _ := strconv.Atoi(c.Query("limit"))

  // set offset
  offset := page*limit - limit

  // count total flights
  var flightCount int64
  if err := db.Model(&models.Flight{}).Count(&flightCount).Error; err != nil {
    utils.JSONResponse(c, 500, false, "failed to count flights", nil)
    return
  }

  // count total pages
  totalPages := math.Ceil(float64(flightCount) / float64(limit))

  // get flights
  var flights []models.Flight
  if err := db.Model(&models.Flight{}).Limit(limit).Offset(offset).Find(&flights).Error; err != nil {
    utils.JSONResponse(c, 500, false, "failed to get flights", nil)
    return
  }

  // check if flights is empty
  if len(flights) == 0 {
    utils.JSONResponse(c, 404, false, "flights not found", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "flights retrieved successfully", gin.H{
    "flights":      flights,
    "totalPages":   totalPages,
    "currentPage":  page,
    "totalFlights": flightCount,
    "limit":        limit,
  })
}
