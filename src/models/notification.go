package models

import "gorm.io/gorm"

type Notification struct {
  gorm.Model
  User    User   `gorm:"foreignKey:UserID"`
  UserID  uint   `gorm:"not null"`
  Title   string `gorm:"not null"`
  Message string `gorm:"not null"`
  IsRead  bool   `gorm:"default:false"`
}
