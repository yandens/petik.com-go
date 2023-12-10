package bio

import (
  "encoding/base64"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "io"
)

func UploadAvatar(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get avatar file
  file, err := c.FormFile("avatar")
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Avatar is required", nil)
    return
  }

  // open file
  f, err := file.Open()
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Failed to open avatar", nil)
  }
  defer f.Close()

  // read file
  fileBytes, err := io.ReadAll(f)
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Failed to read avatar", nil)
    return
  }

  // encode file into base64
  fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)

  // upload avatar to imagekit
  url, err := helpers.UploadToImagekit(c, fileBase64, file.Filename, "/avatars")
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Failed to upload avatar", nil)
    return
  }

  // update user avatar
  if err := db.Model(&models.UserBio{}).Where("id = ?", id).Update("avatar", url).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Failed to update avatar", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Avatar updated successfully", nil)
}
