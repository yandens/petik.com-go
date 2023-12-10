package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/helpers"
)

func WhoAmI(c *gin.Context) {
  id, _ := c.Get("id")
  email, _ := c.Get("email")
  role, _ := c.Get("role")

  if id == "" || email == "" || role == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "Success", gin.H{
    "id":    id,
    "email": email,
    "role":  role,
  })
}
