package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

func GetUsers(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get users
  var users []gin.H
  if err := db.Preload("Role").Model(&models.User{}).Find(&users).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Users not found", nil)
    return
  }

  if len(users) == 0 {
    helpers.JSONResponse(c, 204, false, "Users data is empty", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Users retrieved successfully", users)
}
