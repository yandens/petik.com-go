package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type CreateBioInput struct {
  FirstName   string `json:"firstName" binding:"required"`
  LastName    string `json:"lastName" binding:"required"`
  PhoneNumber string `json:"phoneNumber" binding:"required,min=12,max=13"`
  Address     string `json:"address" binding:"required"`
  City        string `json:"city" binding:"required"`
  Province    string `json:"province" binding:"required"`
}

func CreateBio(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id, _ := c.Get("id")

  // get user email
  email, _ := c.Get("email")

  if id == "" || email == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // type assertion convert interface{} to uint
  var userID uint
  switch id := id.(type) {
  case float64:
    userID = uint(id)
  case uint:
    userID = id
  default:
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
  }

  // check if user bio already exists
  isExist := db.Model(&models.UserBio{}).Where("user_id = ?", userID).Take(&models.UserBio{}).RowsAffected
  if isExist == 1 {
    helpers.JSONResponse(c, 400, false, "User bio already exists", nil)
    return
  }

  // get user input
  var input CreateBioInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // create user bio
  userBio := models.UserBio{
    UserID:      userID,
    FirstName:   input.FirstName,
    LastName:    input.LastName,
    PhoneNumber: input.PhoneNumber,
    Address:     input.Address,
    City:        input.City,
    Province:    input.Province,
    Avatar:      "",
  }

  // save user bio
  if err := db.Create(&userBio).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Could not create user bio", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "User bio created successfully", gin.H{
    "bioId":       userBio.ID, // id from bio table
    "userId":      userID,
    "email":       email,
    "firstName":   userBio.FirstName,
    "lastName":    userBio.LastName,
    "phoneNumber": userBio.PhoneNumber,
    "address":     userBio.Address,
    "city":        userBio.City,
    "province":    userBio.Province,
    "avatar":      userBio.Avatar,
  })
}
