package configs

import (
  "github.com/joho/godotenv"
  "os"
)

func GetEnv(env string) string {
  err := godotenv.Load("../.env")
  if err != nil {
    panic("Error loading .env file")
  }

  return os.Getenv(env)
}
