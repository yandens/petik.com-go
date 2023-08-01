package models

import "gorm.io/gorm"

type User struct {
  gorm.Model
  Role      Role   `gorm:"foreignKey:RoleID"`
  RoleID    uint   `gorm:"not null"`
  Email     string `gorm:"unique;not null"`
  Password  string `gorm:"not null"`
  IsVerifed bool   `gorm:"default:false"`
}
