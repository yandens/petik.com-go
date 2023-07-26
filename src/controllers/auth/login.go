package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
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
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
  }

  // validate input
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Input must be JSON", nil)
  }

  // find user by email
  err = db.Where("email = ?", input.Email).First(&user).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Incorrect Email or Password", nil)
  }

  // compare password
  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
  if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
    utils.JSONResponse(c, 400, false, "Incorrect Email or Password", nil)
  }

  // generate token
  token, err := utils.GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not generate token", nil)
  }

  // return response
  utils.JSONResponse(c, 200, true, "Success", gin.H{
    "email": user.Email,
    "token": token,
  })
}
