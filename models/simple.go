package models

import (
	"gorm.io/gorm"
)

type Simple struct {
	gorm.Model
	Name string `json:"name"`
	Number int `json:"number"`	
}
