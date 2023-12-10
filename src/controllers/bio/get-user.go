package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

func GetUser(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id := c.Param("id")
  if id == "" {
    helpers.JSONResponse(c, 400, false, "User id is required", nil)
    return
  }

  // get user
  var user models.User
  if err := db.Preload("Role").Preload("UserBio").Model(&models.User{}).Where("users.id = ?", id).First(&user).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "User not found", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "User retrieved successfully", user)
}
