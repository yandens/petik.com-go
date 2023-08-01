package db

import (
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
)

func migrateAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // migrate all models
  migrator := db.Migrator()
  if err := migrator.AutoMigrate(&models.User{}, &models.Roles{}, &models.UserBio{}); err != nil { // add more models here
    panic(err)
  }
}

func rollbackAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // rollback all models
  migrator := db.Migrator()
  if err := migrator.DropTable(&models.User{}, &models.Roles{}); err != nil {
    panic(err)
  }
}
