package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "golang.org/x/crypto/bcrypt"
)

type ResetPasswordInput struct {
  Password        string `json:"password" binding:"required"`
  ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
  // get token from query string
  token := c.Query("token")
  if token == "" {
    utils.JSONResponse(c, 400, false, "Token is required", nil)
    return
  }

  // validate token
  validToken, err := utils.ValidateToken(token)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get user id from token
  userID := validToken.Claims.(jwt.MapClaims)["id"]

  // get password from request body
  var input ResetPasswordInput
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // check if password and confirm password is the same
  if input.Password != input.ConfirmPassword {
    utils.JSONResponse(c, 400, false, "Password and confirm password is not the same", nil)
    return
  }

  // hashed password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not hash password", nil)
    return
  }

  // update user
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  err = db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not update user", nil)
    return
  }

  // return response
  utils.JSONResponse(c, 200, true, "Password has been updated", nil)
}
