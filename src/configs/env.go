package configs

import (
  "io/ioutil"
)

func GetEnv(env string) string {
  // use local .env file
  //err := godotenv.Load("../.env")

  // use kubernetes secret (mounted as volume)
  byte, err := ioutil.ReadFile("/app/.env/" + env)
  if err != nil {
    panic(err)
  }

  //return os.Getenv(env)
  return string(byte)
}
