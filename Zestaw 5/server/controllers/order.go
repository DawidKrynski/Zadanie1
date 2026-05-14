package controllers

import (
	"net/http"

	"zadanie3/database"
	"zadanie3/models"

	"github.com/labstack/echo/v4"
)

type orderRequest struct {
	CartID   uint `json:"cart_id"`
	Customer struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	} `json:"customer"`
	Items []orderItemRequest `json:"items"`
}

type orderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func CreateOrder(c echo.Context) error {
	var input orderRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if input.CartID == 0 || input.Customer.Name == "" || input.Customer.Email == "" || len(input.Items) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid order payload"})
	}

	var cart models.Cart
	if err := database.DB.First(&cart, input.CartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}

	order := models.Order{
		CartID:          input.CartID,
		CustomerName:    input.Customer.Name,
		CustomerEmail:   input.Customer.Email,
		CustomerAddress: input.Customer.Address,
		Items:           make([]models.OrderItem, 0, len(input.Items)),
	}
	for _, item := range input.Items {
		if item.ProductID == 0 || item.Quantity <= 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid order item"})
		}
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
		}
		order.Total += product.Price * float64(item.Quantity)
		order.Items = append(order.Items, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create order"})
	}
	return c.JSON(http.StatusCreated, order)
}
