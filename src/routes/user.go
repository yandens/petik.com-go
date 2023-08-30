package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/bio"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func UserRoutes(router *gin.RouterGroup) {
  userRoute := router.Group("/user")
  userRoute.Use(middlewares.Authorized("user"))

  // bio routes
  userRoute.POST("/create-bio", bio.CreateBio)
  userRoute.PUT("/update-bio", bio.UpdateBio)
  userRoute.GET("/get-bio", bio.ReadBio)
  userRoute.POST("/upload-avatar", bio.UploadAvatar)
}
