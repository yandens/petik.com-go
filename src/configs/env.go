package configs

import (
  "io/ioutil"
)

func GetEnv(env string) string {
  //err := godotenv.Load("../.env")
  byte, err := ioutil.ReadFile("/app/.env/" + env)
  if err != nil {
    panic(err)
  }

  //return os.Getenv(env)
  return string(byte)
}
