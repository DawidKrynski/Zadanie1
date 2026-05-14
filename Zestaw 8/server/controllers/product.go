package controllers

import (
	"net/http"

	"zadanie8/database"
	"zadanie8/models"
	"zadanie8/scopes"

	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {
	var products []models.Product

	query := database.DB.Scopes(scopes.WithCategory())

	if catID := scopes.ParseUint(c.QueryParam("category_id")); catID > 0 {
		query = query.Scopes(scopes.ByCategory(catID))
	}

	minPrice := scopes.ParseFloat(c.QueryParam("min_price"))
	maxPrice := scopes.ParseFloat(c.QueryParam("max_price"))
	if minPrice > 0 || maxPrice > 0 {
		query = query.Scopes(scopes.PriceRange(minPrice, maxPrice))
	}

	if c.QueryParam("in_stock") == "true" {
		query = query.Scopes(scopes.InStock())
	}

	page := scopes.ParseInt(c.QueryParam("page"))
	pageSize := scopes.ParseInt(c.QueryParam("page_size"))
	if page > 0 {
		query = query.Scopes(scopes.Paginate(page, pageSize))
	}

	if err := query.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not load products"})
	}
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.Scopes(scopes.WithCategory()).First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := database.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create product"})
	}
	return c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err := database.DB.Save(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update product"})
	}
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	if err := database.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not delete product"})
	}
	return c.JSON(http.StatusNoContent, nil)
}
