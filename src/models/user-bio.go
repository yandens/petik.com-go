package models

import "gorm.io/gorm"

type UserBio struct {
  gorm.Model
  User        User   `gorm:"foreignKey:UserID"`
  UserID      uint   `gorm:"not null"`
  FirstName   string `gorm:"not null"`
  LastName    string `gorm:"not null"`
  PhoneNumber string `gorm:"not null"`
  Address     string `gorm:"not null"`
  City        string `gorm:"not null"`
  Province    string `gorm:"not null"`
  Avatar      string `gorm:"not null"`
}
