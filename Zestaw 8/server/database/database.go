package database

import (
	"zadanie8/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	return DB.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.User{},
		&models.OAuthAccount{},
		&models.AuthSession{},
		&models.OAuthState{},
	)
}
