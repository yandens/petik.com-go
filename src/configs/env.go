package configs

import (
  "io/ioutil"
)

func GetEnv(env string) string {
  //err := godotenv.Load("../.env")
  byte, err := ioutil.ReadFile("/petik-backend-api-secret/" + env)
  if err != nil {
    panic(err)
  }

  //return os.Getenv(env)
  return string(byte)
}
