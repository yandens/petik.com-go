package flight

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

type SearchFlightInput struct {
  Origin      string `json:"origin" binding:"required"`
  Destination string `json:"destination" binding:"required"`
  Date        string `json:"date" binding:"required"`
}

func SearchFlight(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "failed to connect to database", nil)
    return
  }

  // get input
  var input SearchFlightInput
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "invalid input", nil)
    return
  }

  // get flights
  var flights []models.Flight
  if err := db.Model(&models.Flight{}).Where("origin = ? AND destination = ? AND departure LIKE ?", input.Origin, input.Destination, input.Date+"%").Find(&flights).Error; err != nil {
    utils.JSONResponse(c, 400, false, "flight not found", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "flight retrieved successfully", flights)
}
