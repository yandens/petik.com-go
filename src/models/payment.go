package models

import "gorm.io/gorm"

type Payment struct {
  gorm.Model
  Booking       Booking `gorm:"foreignKey:BookingID"`
  BookingID     uint    `gorm:"not null"`
  PaymentMethod string  `gorm:"not null"`
  TotalPrice    float64 `gorm:"not null"`
}
