package models

import "gorm.io/gorm"

type InsuranceEntity struct {
	gorm.Model
	Name string
}
