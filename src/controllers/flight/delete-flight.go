package flight

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

func DeleteFlightSeeder() {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // delete flights
  if err := db.Exec("DELETE FROM flights; VACUUM;").Error; err != nil {
    panic(err)
  }
}

func DeleteFlight(c *gin.Context) {
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

  // delete flight
  if err := db.Model(&models.Flight{}).Where("id = ?", id).Delete(&models.Flight{}).Error; err != nil {
    utils.JSONResponse(c, 400, false, "flight not found", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "flight deleted successfully", nil)
}
