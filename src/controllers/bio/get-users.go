package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/utils"
)

func GetUsers(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get users
  var users []gin.H
  if err := db.Joins("Role").Model(&users).Find(&users).Error; err != nil {
    utils.JSONResponse(c, 400, false, "Users not found", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "Users retrieved successfully", users)
}
