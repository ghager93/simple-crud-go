package models

import (
	"gorm.io/gorm"
)

type Simple struct {
	gorm.Model
	Name string
	Number int
}
