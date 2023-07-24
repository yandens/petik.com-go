package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
)

type RegisterInput struct {
  Email           string `json:"email" binding:"required"`
  Password        string `json:"password" binding:"required"`
  ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func Register(c *gin.Context) {
  var input RegisterInput

  if err := c.ShouldBindJSON(&input); err != nil {
    utils.JSONResponse(c, 400, false, "Input must be JSON", nil)
  }

}
