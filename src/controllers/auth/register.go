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
  var role models.Roles

  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
  }

  // validate input
  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Input must be JSON", nil)
  }

  // check if email already exists
  isExist := db.Where("email = ?", input.Email).First(&models.User{}).RowsAffected
  if isExist == 1 {
    utils.JSONResponse(c, 400, false, "Email already exists", nil)
  }

  // check if password and confirm password are same
  if input.Password != input.ConfirmPassword {
    utils.JSONResponse(c, 400, false, "Password and confirm password must be same", nil)
  }

  // hash password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not hash password", nil)
  }

  // get role id
  err = db.Model(&models.Roles{}).Where("role = ?", "user").First(&role).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not find role", nil)
  }

  // create user
  user := models.User{
    Email:    input.Email,
    Password: string(hashedPassword),
    RoleID:   role.ID,
  }

  err = db.Create(&user).Error
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not create user", nil)
  }

  // return response
  utils.JSONResponse(c, 200, true, "Success", nil)
}
