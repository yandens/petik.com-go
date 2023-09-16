package helpers

import "github.com/gin-gonic/gin"

type Response struct {
  Status  int         `json:"status"`
  Success bool        `json:"success"`
  Message string      `json:"message"`
  Data    interface{} `json:"data"`
}

func JSONResponse(c *gin.Context, status int, success bool, message string, data interface{}) {
  response := Response{
    Status:  status,
    Success: success,
    Message: message,
    Data:    data,
  }

  c.JSON(status, response)
}
