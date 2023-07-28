package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/auth"
)

func AuthRoutes(router *gin.RouterGroup) {
  authRoute := router.Group("/auth")
  authRoute.POST("/login", auth.Login)
  authRoute.POST("/register", auth.Register)
  authRoute.GET("/verify-email", auth.VerifyEmail)
}
