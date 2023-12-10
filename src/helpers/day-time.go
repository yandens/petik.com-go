package helpers

import (
  "math/rand"
  "time"
)

func randomDayHour(min, max int) int {
  // seed current random number generator with current time
  rand.Seed(time.Now().UnixNano())

  // random number between min and max
  return rand.Intn(max-min) + min
}

func SetTime() (time.Time, time.Time) {
  // get time now
  currentTime := time.Now()

  // get random day between 2 - 7 days
  randomDay := randomDayHour(2, 7)

  // set departure time
  departureTime := currentTime.AddDate(0, 0, randomDay).Add(time.Minute * time.Duration(randomDayHour(45, 120)))

  // set arrival time
  arrivalTime := currentTime.AddDate(0, 0, randomDay).Add(time.Minute * time.Duration(randomDayHour(45, 120)*2))

  // return departure time
  return departureTime, arrivalTime
}
