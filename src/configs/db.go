package configs

import (
  "fmt"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
  dbHost := GetEnv("DB_HOST")
  dbUsername := GetEnv("DB_USERNAME")
  dbPassword := GetEnv("DB_PASSWORD")
  dbName := GetEnv("DB_NAME")
  dbPort := GetEnv("DB_PORT")

  // connect to db
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUsername, dbPassword, dbName, dbPort)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil, err
  }

  return db, nil
}
