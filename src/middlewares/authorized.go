package middlewares

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
  "strings"
)

func Authorized(roles ...string) gin.HandlerFunc {
  return func(c *gin.Context) {
    // get token from header
    token := c.GetHeader("Authorization")

    // check if token is missing or not
    if token == "" {
      utils.JSONResponse(c, 401, false, "Authorization header missing", nil)
      c.Abort()
      return
    }

    // extract the token from the Authorization header (Bearer token)
    tokenParts := strings.Split(token, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
      utils.JSONResponse(c, 401, false, "Invalid token format", nil)
      c.Abort()
      return
    }

    // validate the token
    validToken, err := utils.ValidateToken(tokenParts[1])
    if err != nil {
      utils.JSONResponse(c, 401, false, "Invalid token", nil)
      c.Abort()
      return
    }

    // get user claims from token
    user, err := utils.TokenClaims(validToken)
    if err != nil {
      utils.JSONResponse(c, 401, false, "Invalid token", nil)
      c.Abort()
      return
    }

    // check roles length
    if len(roles) < 0 {
      utils.JSONResponse(c, 401, false, "Unauthorized", nil)
      c.Abort()
      return
    }

    for _, role := range roles {
      if user["role"] == role {
        c.Set("id", user["id"])
        c.Set("email", user["email"])
        c.Set("role", user["role"])
        c.Next()
        return
      }

      utils.JSONResponse(c, 401, false, "Unauthorized", nil)
      c.Abort()
      return
    }
  }
}
