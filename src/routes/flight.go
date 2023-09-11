package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/flight"
)

func FlightRoutes(router *gin.RouterGroup) {
  flightRoute := router.Group("/flight")

  // flight routes
  flightRoute.GET("/search-flight", flight.SearchFlight)
}
