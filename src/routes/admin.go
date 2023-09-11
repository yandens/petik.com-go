package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/bio"
  "github.com/yandens/petik.com-go/src/controllers/flight"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func AdminRoutes(router *gin.RouterGroup) {
  adminRoute := router.Group("/admin")
  adminRoute.Use(middlewares.Authorized("admin"))

  // bio routes
  adminRoute.GET("/get-users", bio.GetUsers)
  adminRoute.GET("/get-user/:id", bio.GetUser)

  // flight routes
  adminRoute.POST("/create-flight", flight.CreateFlight)
  adminRoute.GET("/get-flights", flight.GetFlights)
  adminRoute.GET("/get-flight/:id", flight.GetFlight)
  adminRoute.PUT("/update-flight/:id", flight.UpdateFlight)
  adminRoute.DELETE("/delete-flight/:id", flight.DeleteFlight)
}
