package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
  "golang.org/x/crypto/bcrypt"
)

type ResetPasswordInput struct {
  Password        string `json:"password" binding:"required"`
  ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get token from query string
  token := c.Query("token")
  if token == "" {
    helpers.JSONResponse(c, 400, false, "Token is required", nil)
    return
  }

  // validate token
  validToken, err := helpers.ValidateToken(token)
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get user claims from token
  claims, err := helpers.TokenClaims(validToken)
  if err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get password from request body
  var input ResetPasswordInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // check if password and confirm password is the same
  if input.Password != input.ConfirmPassword {
    helpers.JSONResponse(c, 400, false, "Password and confirm password is not the same", nil)
    return
  }

  // hashed password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not hash password", nil)
    return
  }

  // update user
  err = db.Model(&models.User{}).Where("id = ?", claims["id"]).Update("password", string(hashedPassword)).Error
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not update user", nil)
    return
  }

  // return response
  helpers.JSONResponse(c, 200, true, "Password has been updated", nil)
}
