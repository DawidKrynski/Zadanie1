package controllers

import (
	"net/http"
	"strings"

	"zadanie8/database"
	"zadanie8/models"

	"github.com/labstack/echo/v4"
)

type orderRequest struct {
	CartID   uint `json:"cart_id"`
	Customer struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	} `json:"customer"`
	Items []struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"items"`
}

func CreateOrder(c echo.Context) error {
	var input orderRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if input.CartID == 0 || strings.TrimSpace(input.Customer.Email) == "" || len(input.Items) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "cart, customer email and items are required"})
	}

	order := models.Order{
		CartID:   input.CartID,
		Customer: strings.TrimSpace(input.Customer.Name),
		Email:    strings.TrimSpace(input.Customer.Email),
		Address:  strings.TrimSpace(input.Customer.Address),
		Items:    make([]models.OrderItem, 0, len(input.Items)),
	}

	for _, item := range input.Items {
		if item.ProductID == 0 || item.Quantity <= 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid order item"})
		}

		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product in order"})
		}

		lineTotal := product.Price * float64(item.Quantity)
		order.Total += lineTotal
		order.Items = append(order.Items, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create order"})
	}

	return c.JSON(http.StatusCreated, order)
}
