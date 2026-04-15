package main

import (
	"zadanie3/controllers"
	"zadanie3/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.Init()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Products - CRUD
	e.GET("/products", controllers.GetProducts)
	e.GET("/products/:id", controllers.GetProduct)
	e.POST("/products", controllers.CreateProduct)
	e.PUT("/products/:id", controllers.UpdateProduct)
	e.DELETE("/products/:id", controllers.DeleteProduct)

	// Categories - CRUD
	e.GET("/categories", controllers.GetCategories)
	e.GET("/categories/:id", controllers.GetCategory)
	e.POST("/categories", controllers.CreateCategory)
	e.PUT("/categories/:id", controllers.UpdateCategory)
	e.DELETE("/categories/:id", controllers.DeleteCategory)

	// Carts
	e.POST("/carts", controllers.CreateCart)
	e.GET("/carts/:id", controllers.GetCart)
	e.POST("/carts/:id/items", controllers.AddCartItem)
	e.PUT("/carts/:id/items/:itemId", controllers.UpdateCartItem)
	e.DELETE("/carts/:id/items/:itemId", controllers.DeleteCartItem)
	e.DELETE("/carts/:id", controllers.DeleteCart)

	e.Logger.Fatal(e.Start(":8080"))
}
