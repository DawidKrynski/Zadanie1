package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"zadanie3/controllers"
	"zadanie3/models"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"zadanie3/database"
)

func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.DB.AutoMigrate(&models.Category{}, &models.Product{}, &models.Cart{}, &models.CartItem{})
}

func seedCategory(t *testing.T) models.Category {
	cat := models.Category{Name: "Elektronika", Description: "Sprzęt elektroniczny"}
	database.DB.Create(&cat)
	assert.NotZero(t, cat.ID)
	return cat
}

func seedProduct(t *testing.T, catID uint) models.Product {
	p := models.Product{Name: "Laptop", Description: "Laptop testowy", Price: 2999.99, Quantity: 10, CategoryID: catID}
	database.DB.Create(&p)
	assert.NotZero(t, p.ID)
	return p
}

// === CATEGORY TESTS ===

func TestCreateCategory(t *testing.T) {
	setupTestDB(t)
	e := echo.New()
	body := `{"name":"Książki","description":"Kategoria książek"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.CreateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var cat models.Category
	json.Unmarshal(rec.Body.Bytes(), &cat)
	assert.Equal(t, "Książki", cat.Name)
	assert.Equal(t, "Kategoria książek", cat.Description)
	assert.NotZero(t, cat.ID)
}

func TestGetCategories(t *testing.T) {
	setupTestDB(t)
	seedCategory(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.GetCategories(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var cats []models.Category
	json.Unmarshal(rec.Body.Bytes(), &cats)
	assert.GreaterOrEqual(t, len(cats), 1)
	assert.Equal(t, "Elektronika", cats[0].Name)
}

func TestGetCategory(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.Category
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, cat.Name, result.Name)
}

func TestGetCategoryNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.GetCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCategory(t *testing.T) {
	setupTestDB(t)
	seedCategory(t)

	e := echo.New()
	body := `{"name":"Elektronika Pro","description":"Zaktualizowana"}`
	req := httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.UpdateCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.Category
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, "Elektronika Pro", result.Name)
}

func TestDeleteCategory(t *testing.T) {
	setupTestDB(t)
	seedCategory(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.DeleteCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCategoryNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/categories/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.DeleteCategory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// === PRODUCT TESTS ===

func TestCreateProduct(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)

	e := echo.New()
	body := `{"name":"Telefon","description":"Smartfon","price":1999.99,"quantity":5,"category_id":` + strings.TrimRight(strings.TrimLeft(toJSON(cat.ID), ""), "") + `}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.CreateProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var p models.Product
	json.Unmarshal(rec.Body.Bytes(), &p)
	assert.Equal(t, "Telefon", p.Name)
	assert.Equal(t, 1999.99, p.Price)
	assert.Equal(t, 5, p.Quantity)
}

func TestGetProducts(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.GetProducts(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.GreaterOrEqual(t, len(products), 1)
}

func TestGetProduct(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	prod := seedProduct(t, cat.ID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.Product
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, prod.Name, result.Name)
	assert.Equal(t, prod.Price, result.Price)
}

func TestGetProductNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateProduct(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID)

	e := echo.New()
	body := `{"name":"Laptop Pro","description":"Zaktualizowany","price":3999.99,"quantity":3}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.UpdateProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.Product
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, "Laptop Pro", result.Name)
	assert.Equal(t, 3999.99, result.Price)
}

func TestDeleteProduct(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.DeleteProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteProductNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.DeleteProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// === CART TESTS ===

func TestCreateCart(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/carts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.CreateCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var cart models.Cart
	json.Unmarshal(rec.Body.Bytes(), &cart)
	assert.NotZero(t, cart.ID)
}

func TestGetCart(t *testing.T) {
	setupTestDB(t)
	cart := models.Cart{}
	database.DB.Create(&cart)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carts/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.GetCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetCartNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carts/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.GetCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAddCartItem(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	prod := seedProduct(t, cat.ID)
	cart := models.Cart{}
	database.DB.Create(&cart)

	e := echo.New()
	body := `{"product_id":` + toJSON(prod.ID) + `,"quantity":2}`
	req := httptest.NewRequest(http.MethodPost, "/carts/1/items", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.AddCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var item models.CartItem
	json.Unmarshal(rec.Body.Bytes(), &item)
	assert.Equal(t, 2, item.Quantity)
	assert.Equal(t, prod.ID, item.ProductID)
}

func TestAddCartItemCartNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	body := `{"product_id":1,"quantity":1}`
	req := httptest.NewRequest(http.MethodPost, "/carts/999/items", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.AddCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAddCartItemProductNotFound(t *testing.T) {
	setupTestDB(t)
	cart := models.Cart{}
	database.DB.Create(&cart)

	e := echo.New()
	body := `{"product_id":999,"quantity":1}`
	req := httptest.NewRequest(http.MethodPost, "/carts/1/items", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.AddCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCartItem(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	prod := seedProduct(t, cat.ID)
	cart := models.Cart{}
	database.DB.Create(&cart)
	item := models.CartItem{CartID: cart.ID, ProductID: prod.ID, Quantity: 1}
	database.DB.Create(&item)

	e := echo.New()
	body := `{"quantity":5}`
	req := httptest.NewRequest(http.MethodPut, "/carts/1/items/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "1")

	err := controllers.UpdateCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.CartItem
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, 5, result.Quantity)
}

func TestUpdateCartItemNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	body := `{"quantity":5}`
	req := httptest.NewRequest(http.MethodPut, "/carts/1/items/999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "999")

	err := controllers.UpdateCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCartItem(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	prod := seedProduct(t, cat.ID)
	cart := models.Cart{}
	database.DB.Create(&cart)
	item := models.CartItem{CartID: cart.ID, ProductID: prod.ID, Quantity: 1}
	database.DB.Create(&item)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/carts/1/items/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "1")

	err := controllers.DeleteCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCartItemNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/carts/1/items/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "999")

	err := controllers.DeleteCartItem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCart(t *testing.T) {
	setupTestDB(t)
	cart := models.Cart{}
	database.DB.Create(&cart)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/carts/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := controllers.DeleteCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteCartNotFound(t *testing.T) {
	setupTestDB(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/carts/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := controllers.DeleteCart(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// === PRODUCT FILTER TESTS ===

func TestGetProductsFilterByCategory(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products?category_id=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames()
	c.QueryParams().Set("category_id", "1")

	err := controllers.GetProducts(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.GreaterOrEqual(t, len(products), 1)
	assert.Equal(t, cat.ID, products[0].CategoryID)
}

func TestGetProductsFilterInStock(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID)
	// produkt bez stock
	outOfStock := models.Product{Name: "Brak", Price: 10, Quantity: 0, CategoryID: cat.ID}
	database.DB.Create(&outOfStock)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products?in_stock=true", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Set("in_stock", "true")

	err := controllers.GetProducts(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	json.Unmarshal(rec.Body.Bytes(), &products)
	for _, p := range products {
		assert.Greater(t, p.Quantity, 0)
	}
}

func TestGetProductsPriceRange(t *testing.T) {
	setupTestDB(t)
	cat := seedCategory(t)
	seedProduct(t, cat.ID) // 2999.99
	cheap := models.Product{Name: "Tani", Price: 10, Quantity: 5, CategoryID: cat.ID}
	database.DB.Create(&cheap)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products?min_price=100&max_price=5000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Set("min_price", "100")
	c.QueryParams().Set("max_price", "5000")

	err := controllers.GetProducts(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, "Laptop", products[0].Name)
}

// Helper
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
