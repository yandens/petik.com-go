package models

import "gorm.io/gorm"

type Role struct {
  gorm.Model
  Role string `gorm:"unique;not null"`
  User []User `gorm:"foreignKey:RoleID"`
}
