package models

import "gorm.io/gorm"

type BookingDetail struct {
  gorm.Model
  Booking       Booking `gorm:"foreignKey:BookingID"`
  BookingID     uint    `gorm:"not null"`
  PassengerName string  `gorm:"not null"`
  PassengerAge  int     `gorm:"not null"`
  NIK           string  `gorm:"not null"`
  SeatNumber    string  `gorm:"not null"`
  QRCode        string  `gorm:"not null"`
}
