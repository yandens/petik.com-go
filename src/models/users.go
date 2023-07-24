package models

import "gorm.io/gorm"

type User struct {
  gorm.Model
  Email    string `gorm:"unique;not null"`
  Password string `gorm:"not null"`
  Role     Roles  `gorm:"foreignKey:RoleID"`
  RoleID   uint   `gorm:"not null"`
}
