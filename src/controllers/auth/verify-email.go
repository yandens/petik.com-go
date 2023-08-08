package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

func VerifyEmail(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get token from url
  token := c.Query("token")
  if token == "" {
    utils.JSONResponse(c, 400, false, "Token is required", nil)
    return
  }

  // validate the user
  validToken, err := utils.ValidateToken(token)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get user claims from token
  claims, err := utils.TokenClaims(validToken)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get user
  var user models.User
  if err := db.Joins("Role").Model(&models.User{}).Where("users.id = ?", claims["id"]).First(&user).Error; err != nil {
    utils.JSONResponse(c, 400, false, "User not found", nil)
    return
  }

  // check if user is already verified
  if user.IsVerified {
    utils.JSONResponse(c, 400, false, "User already verified", nil)
    return
  }

  // update user
  err = db.Model(&models.User{}).Where("id = ?", user.ID).Update("is_verified", true).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not update user", nil)
    return
  }

  // generate token
  token, err = utils.GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not generate token", nil)
    return
  }

  // return response
  utils.JSONResponse(c, 200, true, "Success", gin.H{
    "email": user.Email,
    "token": token,
  })
}
