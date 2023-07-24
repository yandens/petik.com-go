package configs

import (
  "fmt"
  "github.com/yandens/petik.com-go/src/models"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

func ConnectToDB() {
  dbHost := GetEnv("DB_HOST")
  dbUsername := GetEnv("DB_USERNAME")
  dbPassword := GetEnv("DB_PASSWORD")
  dbName := GetEnv("DB_NAME")
  dbPort := GetEnv("DB_PORT")

  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUsername, dbPassword, dbName, dbPort)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic(err)
  }

  db.AutoMigrate(&models.User{}, &models.Roles{})
}
