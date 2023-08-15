package user

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

func ReadBio(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id, _ := c.Get("id")
  if id == "" {
    utils.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get user bio
  var userBio models.UserBio
  if err := db.Joins("User").Model(&models.UserBio{}).Where("user_id = ?", id).First(&userBio).Error; err != nil {
    utils.JSONResponse(c, 400, false, "User bio not found", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "User bio retrieved successfully", userBio)
}
