package controllers

import (
	"net/http"

	"zadanie3/database"
	"zadanie3/models"

	"github.com/labstack/echo/v4"
)

func CreateCart(c echo.Context) error {
	cart := models.Cart{}
	database.DB.Create(&cart)
	return c.JSON(http.StatusCreated, cart)
}

func GetCart(c echo.Context) error {
	id := c.Param("id")
	var cart models.Cart
	if err := database.DB.Preload("Items.Product").First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}
	return c.JSON(http.StatusOK, cart)
}

func AddCartItem(c echo.Context) error {
	cartID := c.Param("id")
	var cart models.Cart
	if err := database.DB.First(&cart, cartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}

	var item models.CartItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	item.CartID = cart.ID

	var product models.Product
	if err := database.DB.First(&product, item.ProductID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	database.DB.Create(&item)
	database.DB.Preload("Product").First(&item, item.ID)
	return c.JSON(http.StatusCreated, item)
}

func UpdateCartItem(c echo.Context) error {
	itemID := c.Param("itemId")
	var item models.CartItem
	if err := database.DB.First(&item, itemID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart item not found"})
	}

	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	item.Quantity = input.Quantity
	database.DB.Save(&item)
	database.DB.Preload("Product").First(&item, item.ID)
	return c.JSON(http.StatusOK, item)
}

func DeleteCartItem(c echo.Context) error {
	itemID := c.Param("itemId")
	var item models.CartItem
	if err := database.DB.First(&item, itemID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart item not found"})
	}
	database.DB.Delete(&item)
	return c.JSON(http.StatusNoContent, nil)
}

func DeleteCart(c echo.Context) error {
	id := c.Param("id")
	var cart models.Cart
	if err := database.DB.First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}
	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	database.DB.Delete(&cart)
	return c.JSON(http.StatusNoContent, nil)
}
