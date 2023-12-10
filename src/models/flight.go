package models

import (
  "gorm.io/gorm"
  "time"
)

type Flight struct {
  gorm.Model
  Airline         string    `gorm:"not null"`
  AirlineLogo     string    `gorm:"not null"`
  Origin          string    `gorm:"not null"`
  OriginCity      string    `gorm:"not null"`
  Destination     string    `gorm:"not null"`
  DestinationCity string    `gorm:"not null"`
  Departure       time.Time `gorm:"not null"`
  Arrival         time.Time `gorm:"not null"`
  Booking         []Booking `gorm:"foreignKey:FlightID"`
}
