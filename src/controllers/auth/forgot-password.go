package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
)

type ForgotPasswordInput struct {
  Email string `json:"email" binding:"required"`
}

func ForgotPassword(c *gin.Context) {
  var input ForgotPasswordInput

  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get email from request body
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // find user by email
  var user models.User
  if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
    utils.JSONResponse(c, 400, false, "Email not found", nil)
    return
  }

  // send email
  utils.SendEmail(user)

  // return response
  utils.JSONResponse(c, 200, true, "Success", nil)
}
