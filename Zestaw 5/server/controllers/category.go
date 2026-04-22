package controllers

import (
	"net/http"

	"zadanie3/database"
	"zadanie3/models"

	"github.com/labstack/echo/v4"
)

func GetCategories(c echo.Context) error {
	var categories []models.Category
	database.DB.Find(&categories)
	return c.JSON(http.StatusOK, categories)
}

func GetCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.Preload("Products").First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Category not found"})
	}
	return c.JSON(http.StatusOK, category)
}

func CreateCategory(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	database.DB.Create(&category)
	return c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Category not found"})
	}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	database.DB.Save(&category)
	return c.JSON(http.StatusOK, category)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Category not found"})
	}
	database.DB.Delete(&category)
	return c.JSON(http.StatusNoContent, nil)
}
