package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type ForgotPasswordInput struct {
  Email string `json:"email" binding:"required"`
}

func ForgotPassword(c *gin.Context) {
  var input ForgotPasswordInput

  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get email from request body
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Invalid request", nil)
    return
  }

  // find user by email
  var user models.User
  if err := db.Joins("Role").Model(&models.User{}).Where("email = ?", input.Email).First(&user).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Email not found", nil)
    return
  }

  // send email
  helpers.SendEmail(user, "reset-password", "Reset Password")

  // return response
  helpers.JSONResponse(c, 200, true, "Success", nil)
}
