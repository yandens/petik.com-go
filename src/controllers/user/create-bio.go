package user

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

type CreateBioInput struct {
  FirstName string `json:"firstName" binding:"required"`
  LastName  string `json:"lastName" binding:"required"`
  Address   string `json:"address" binding:"required"`
  City      string `json:"city" binding:"required"`
  Province  string `json:"province" binding:"required"`
  Avatar    string `json:"avatar" binding:"required"`
}

func CreateBio(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id, _ := c.Get("id")

  // get user email
  email, _ := c.Get("email")

  if id == "" || email == "" {
    utils.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // check if user bio already exists
  isExist := db.Model(&models.UserBio{}).Where("user_id = ?", id).Take(&models.UserBio{}).RowsAffected
  if isExist == 1 {
    utils.JSONResponse(c, 400, false, "User bio already exists", nil)
    return
  }

  // get user input
  var input CreateBioInput
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // create user bio
  userBio := models.UserBio{
    UserID:    id.(uint),
    FirstName: input.FirstName,
    LastName:  input.LastName,
    Address:   input.Address,
    City:      input.City,
    Province:  input.Province,
    Avatar:    input.Avatar,
  }

  // save user bio
  if err := db.Create(&userBio).Error; err != nil {
    utils.JSONResponse(c, 500, false, "Could not create user bio", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "User bio created successfully", gin.H{
    "user_id": id,
    "email":   email,
    "userBio": userBio,
  })
}
