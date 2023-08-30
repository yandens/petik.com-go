package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/user"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func UserRoutes(router *gin.RouterGroup) {
  userRoute := router.Group("/user")
  userRoute.POST("/create-bio", middlewares.Authorized("user"), user.CreateBio)
  userRoute.PUT("/update-bio", middlewares.Authorized("user"), user.UpdateBio)
  userRoute.GET("/get-bio", middlewares.Authorized("user"), user.ReadBio)
  userRoute.POST("upload-avatar", middlewares.Authorized("user"), user.UploadAvatar)
  userRoute.GET("/get-users", middlewares.Authorized("admin"), user.GetUsers)
  userRoute.GET("/get-user/:id", middlewares.Authorized("admin"), user.GetUser)
}
