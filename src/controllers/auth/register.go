package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
  Email           string `json:"email" binding:"required"`
  Password        string `json:"password" binding:"required"`
  ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func Register(c *gin.Context) {
  var input RegisterInput
  var role models.Role

  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // validate input
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // check if email already exists
  isExist := db.Model(&models.User{}).Where("email = ?", input.Email).Take(&models.User{}).RowsAffected
  if isExist == 1 {
    utils.JSONResponse(c, 400, false, "Email already exists", nil)
    return
  }

  // check if password and confirm password are same
  if input.Password != input.ConfirmPassword {
    utils.JSONResponse(c, 400, false, "Password and confirm password must be same", nil)
  }

  // hash password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not hash password", nil)
    return
  }

  // get role id
  err = db.Model(&models.Role{}).Where("role = ?", "user").First(&role).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not find role", nil)
    return
  }

  // create user
  user := models.User{
    Email:      input.Email,
    Password:   string(hashedPassword),
    RoleID:     role.ID,
    IsVerified: false,
  }

  err = db.Create(&user).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not create user", nil)
    return
  }

  // send verify email
  utils.SendEmail(user, "verify-email")

  // return response
  utils.JSONResponse(c, 200, true, "Success", nil)
}
