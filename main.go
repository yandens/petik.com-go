package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
)

func main() {
  router := gin.Default()

  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    MaxAge:       600,
  }))

  // routes
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hello world!",
    })
  })

  router.Run(configs.GetEnv("HOST") + ":" + configs.GetEnv("PORT"))
}
