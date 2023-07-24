package models

import "gorm.io/gorm"

type Roles struct {
  gorm.Model
  Role string `gorm:"unique;not null"`
}
