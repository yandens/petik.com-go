package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
)

func WhoAmI(c *gin.Context) {
  id := c.GetString("id")
  email := c.GetString("email")
  role := c.GetString("role")

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
