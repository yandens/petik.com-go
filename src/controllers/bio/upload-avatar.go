package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

func UploadAvatar(c *gin.Context) {
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

  // get avatar file
  file, err := c.FormFile("avatar")
  if err != nil {
    utils.JSONResponse(c, 400, false, "Avatar is required", nil)
    return
  }

  // upload avatar to imagekit
  url, err := utils.UploadToImagekit(c, file, file.Filename)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Failed to upload avatar", nil)
    return
  }

  // update user avatar
  if err := db.Model(&models.UserBio{}).Where("id = ?", id).Update("avatar", url).Error; err != nil {
    utils.JSONResponse(c, 400, false, "Failed to update avatar", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "Avatar updated successfully", nil)
}
