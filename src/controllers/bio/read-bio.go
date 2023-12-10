package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

func ReadBio(c *gin.Context) {
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

  // get user bio
  var userBio models.UserBio
  if err := db.Preload("User").Model(&models.UserBio{}).Where("user_id = ?", id).First(&userBio).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "User bio not found", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "User bio retrieved successfully", gin.H{
    "bioId":       userBio.ID,
    "userId":      userBio.User.ID,
    "email":       userBio.User.Email,
    "firstName":   userBio.FirstName,
    "lastName":    userBio.LastName,
    "phoneNumber": userBio.PhoneNumber,
    "address":     userBio.Address,
    "city":        userBio.City,
    "province":    userBio.Province,
    "avatar":      userBio.Avatar,
  })
}
