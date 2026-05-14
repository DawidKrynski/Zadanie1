package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CartID   uint        `json:"cart_id"`
	Customer string      `json:"customer"`
	Email    string      `json:"email"`
	Address  string      `json:"address"`
	Total    float64     `json:"total"`
	Items    []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
