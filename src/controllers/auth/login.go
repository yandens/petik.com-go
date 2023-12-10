package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
  Email    string `json:"email" binding:"required"`
  Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
  var input LoginInput
  var user models.User

  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // validate input
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // find user by email
  err = db.Preload("Role").Model(&models.User{}).Where("email = ?", input.Email).First(&user).Error
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Incorrect Email or Password", nil)
    return
  }

  // check account type
  if user.AccountType != "basic" {
    helpers.JSONResponse(c, 400, false, "Incorrect Email or Password", nil)
  }

  // compare password
  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
  if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
    helpers.JSONResponse(c, 400, false, "Incorrect Email or Password", nil)
    return
  }

  // generate token
  token, err := helpers.GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not generate token", nil)
    return
  }

  // return response
  helpers.JSONResponse(c, 200, true, "Success", gin.H{
    "email": user.Email,
    "token": token,
  })
}
