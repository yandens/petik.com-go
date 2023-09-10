package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/routes"
  "github.com/yandens/petik.com-go/src/utils"
)

func main() {
  //db.MigrateAll()
  //db.RollbackAll()
  //db.SeedAll()
  router := gin.Default()

  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    MaxAge:       600,
  }))

  // routes
  // default route
  router.GET("/", func(c *gin.Context) {
    utils.JSONResponse(c, 200, true, "Server Running Well", nil)
  })

  // api group route
  api := router.Group("/api")

  // auth routes
  routes.AuthRoutes(api)

  // user routes
  routes.UserRoutes(api)

  // admin routes
  routes.AdminRoutes(api)

  // airport routes
  routes.AirportRoutes(api)

  // flight routes
  routes.FlightRoutes(api)

  router.Run(configs.GetEnv("HOST") + ":" + configs.GetEnv("PORT"))
}
