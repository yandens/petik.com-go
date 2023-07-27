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

  // return response
  utils.JSONResponse(c, 200, true, "Success", nil)
}
