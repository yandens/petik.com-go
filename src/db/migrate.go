package db

import (
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
)

func MigrateAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // migrate all models
  migrator := db.Migrator()
  if err := migrator.AutoMigrate(&models.Role{}, &models.User{}, &models.UserBio{}); err != nil { // add more models here
    panic(err)
  }
}

func RollbackAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // rollback all models
  migrator := db.Migrator()
  if err := migrator.DropTable(&models.Role{}, &models.User{}, &models.UserBio{}); err != nil {
    panic(err)
  }
}
