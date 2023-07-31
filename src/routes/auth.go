package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/auth"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func AuthRoutes(router *gin.RouterGroup) {
  authRoute := router.Group("/auth")
  authRoute.POST("/login", auth.Login)
  authRoute.POST("/register", auth.Register)
  authRoute.GET("/verify-email", auth.VerifyEmail)
  authRoute.POST("/forgot-password", auth.ForgotPassword)
  authRoute.POST("/reset-password", auth.ResetPassword)
  authRoute.GET("/whoami", middlewares.Authorized("user", "admin"), auth.WhoAmI)
}
