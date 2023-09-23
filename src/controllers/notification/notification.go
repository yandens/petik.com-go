package notification

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type NotificationResponse struct {
  ID      uint   `json:"id"`
  Title   string `json:"title"`
  Message string `json:"message"`
  Date    string `json:"date"`
}

func GetNotifications(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to connect to database", nil)
  }

  // get user id from middleware
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
  }

  // type assertion convert interface{} to uint
  var userID uint
  switch id := id.(type) {
  case float64:
    userID = uint(id)
  case uint:
    userID = id
  default:
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
  }

  // get notifications
  var notifications []models.Notification
  if err := db.Model(&models.Notification{}).Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
    helpers.JSONResponse(c, 500, false, "Failed to get notifications", nil)
    return
  }

  // convert notifications to response
  var response []NotificationResponse
  for _, notification := range notifications {
    response = append(response, NotificationResponse{
      ID:      notification.ID,
      Title:   notification.Title,
      Message: notification.Message,
      Date:    notification.CreatedAt.Format("2006-01-02 15:04:05"),
    })
  }

  helpers.JSONResponse(c, 200, true, "Success", response)
}
