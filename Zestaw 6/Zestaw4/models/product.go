package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
