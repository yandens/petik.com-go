package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/routes"
  "github.com/yandens/petik.com-go/src/utils"
)

func main() {
  configs.AutoMigrates()
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

  //authRoutes := api.Group("/auth")
  //authRoutes.POST("/register", auth.Register)
  //authRoutes.POST("/login", auth.Login)

  router.Run(configs.GetEnv("HOST") + ":" + configs.GetEnv("PORT"))
}
