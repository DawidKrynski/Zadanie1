package main

import (
	"log"
	"os"
	"strings"

	"zadanie3/controllers"
	"zadanie3/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization", "Accept"},
	}))

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

	e.POST("/orders", controllers.CreateOrder)

	e.Logger.Fatal(e.Start(":8080"))
}

func allowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5173", "http://localhost:3000"}
	}

	values := strings.Split(origins, ",")
	allowed := make([]string, 0, len(values))
	for _, origin := range values {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed = append(allowed, origin)
		}
	}
	if len(allowed) == 0 {
		return []string{"http://localhost:5173", "http://localhost:3000"}
	}
	return allowed
}
