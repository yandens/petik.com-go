package db

import (
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "golang.org/x/crypto/bcrypt"
)

func SeedAll() {
  // connect to db
  db, err := configs.ConnectToDB()
  if err != nil {
    panic(err)
  }

  // define roles
  roles := []models.Role{
    {Role: "user"},
    {Role: "admin"},
  }

  // hash password for admin
  hashedPasswordAdmin, err := bcrypt.GenerateFromPassword([]byte(configs.GetEnv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }

  // define admin
  admin := models.User{
    Email:       configs.GetEnv("ADMIN_EMAIL"),
    Password:    string(hashedPasswordAdmin),
    RoleID:      2,
    AccountType: "basic",
    IsVerified:  true,
  }

  // hash password for user
  hashedPasswordUser, err := bcrypt.GenerateFromPassword([]byte(configs.GetEnv("USER_PASSWORD")), bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }

  // define user
  user := models.User{
    Email:       configs.GetEnv("USER_EMAIL"),
    Password:    string(hashedPasswordUser),
    RoleID:      1,
    AccountType: "basic",
    IsVerified:  true,
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

  if err := tx.Create(&admin).Error; err != nil {
    tx.Rollback()
  }

  if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
  }

  // commit transaction
  tx.Commit()
}
