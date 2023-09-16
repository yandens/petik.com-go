package models

import "gorm.io/gorm"

type Booking struct {
  gorm.Model
  User           User            `gorm:"foreignKey:UserID"`
  UserID         uint            `gorm:"not null"`
  Flight         Flight          `gorm:"foreignKey:FlightID"`
  FlightID       uint            `gorm:"not null"`
  FlightClass    string          `gorm:"not null"`
  TotalPassenger int             `gorm:"not null"`
  TotalPrice     float64         `gorm:"not null"`
  Status         string          `gorm:"not null"`
  BookingDetail  []BookingDetail `gorm:"foreignKey:BookingID"`
}
