package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

func VerifyEmail(c *gin.Context) {
  // get token from url
  token := c.Query("token")

  // validate the user
  validToken, err := utils.ValidateToken(token)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid token", nil)
  }

  // get user id from token
  userID := validToken.Claims.(jwt.MapClaims)["id"]

  // update user
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
  }

  err = db.Model(&models.User{}).Where("id = ?", userID).Update("is_verified", true).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not update user", nil)
  }

  // get user
  var user models.User
  if err := db.Model(&models.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
    utils.JSONResponse(c, 400, false, "User not found", nil)
  }

  // generate token
  token, err = utils.GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not generate token", nil)
  }

  // return response
  utils.JSONResponse(c, 200, true, "Success", gin.H{
    "email": user.Email,
    "token": token,
  })
}
