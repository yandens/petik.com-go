package helpers

import (
  "fmt"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
)

func LoadNotifications(userID int) []models.Notification {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    fmt.Println(err)
    return nil
  }

  // get notifications
  var notifications []models.Notification
  if err := db.Model(&models.Notification{}).Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
    fmt.Println(err)
    return nil
  }

  return notifications
}

func UpdateNotifications(userID int) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    fmt.Println(err)
    return
  }

  // update notifications
  if err := db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error; err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println("Success")
}
