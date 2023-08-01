package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
)

func WhoAmI(c *gin.Context) {
  id, _ := c.Get("id")
  email, _ := c.Get("email")
  role, _ := c.Get("role")

  if id == "" || email == "" || role == "" {
    utils.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  utils.JSONResponse(c, 200, true, "Success", gin.H{
    "id":    id,
    "email": email,
    "role":  role,
  })
}
