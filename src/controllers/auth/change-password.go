package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "golang.org/x/crypto/bcrypt"
)

type ChangePasswordInput struct {
  OldPassword     string `json:"oldPassword" binding:"required"`
  NewPassword     string `json:"newPassword" binding:"required"`
  ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func ChangePassword(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get id from middleware
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get input from body
  var input ChangePasswordInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // get user
  var user models.User
  if err := db.Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Could not find user", nil)
    return
  }

  // check if old password is correct
  if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
    helpers.JSONResponse(c, 400, false, "Old password is incorrect", nil)
    return
  }

  // check if password and confirm password are same
  if input.NewPassword != input.ConfirmPassword {
    helpers.JSONResponse(c, 400, false, "Password and confirm password must be same", nil)
    return
  }

  // hash password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not hash password", nil)
    return
  }

  // update user
  if err := db.Model(&models.User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Could not update user", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Success", nil)
}
