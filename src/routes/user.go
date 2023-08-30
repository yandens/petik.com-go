package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/bio"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func BioRoutes(router *gin.RouterGroup) {
  bioRoute := router.Group("/bio")
  bioRoute.POST("/create-bio", middlewares.Authorized("user"), bio.CreateBio)
  bioRoute.PUT("/update-bio", middlewares.Authorized("user"), bio.UpdateBio)
  bioRoute.GET("/get-bio", middlewares.Authorized("user"), bio.ReadBio)
  bioRoute.POST("upload-avatar", middlewares.Authorized("user"), bio.UploadAvatar)
  bioRoute.GET("/get-users", middlewares.Authorized("admin"), bio.GetUsers)
  bioRoute.GET("/get-user/:id", middlewares.Authorized("admin"), bio.GetUser)
}
