package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}
