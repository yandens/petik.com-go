package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/airport"
)

func AirportRoutes(router *gin.RouterGroup) {
  airportRoute := router.Group("/airport")

  // airport routes
  airportRoute.GET("/search/:search", airport.GetAirportBySearch)
}
