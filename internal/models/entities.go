package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InsuranceEntity struct {
	gorm.Model
	Name string
}

type PatientEntity struct {
	gorm.Model
	Name        string
	Email       string
	Phone       string
	InsuranceID int
	Insurance   InsuranceEntity
	DateOfBirth datatypes.Date
}
