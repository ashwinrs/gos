package models

import "gorm.io/gorm"

type PetEntity struct {
	gorm.Model
	Id   int64
	Name string
	Tag  *string
}
