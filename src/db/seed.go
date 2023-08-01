package db

import (
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
)

func seedAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // define roles
  roles := []models.Roles{
    {Role: "user"},
    {Role: "admin"},
  }

  // start transaction
  tx := db.Begin()
  defer func() {
    if r := recover(); r != nil {
      tx.Rollback()
    }
  }()

  for _, role := range roles {
    if err := tx.Create(&role).Error; err != nil {
      tx.Rollback()
    }
  }

  // commit transaction
  tx.Commit()
}
