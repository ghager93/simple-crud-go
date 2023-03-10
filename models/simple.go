package models

import (
	"gorm.io/gorm"
)

type Simple struct {
	gorm.Model
	Name string `json:"name" gorm:"not null" validate:"required"`
	Number int `json:"number" gorm:"not null" validate:"required"`	
}
