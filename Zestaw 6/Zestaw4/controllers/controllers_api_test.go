package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"zadanie3/controllers"
	"zadanie3/database"
	"zadanie3/models"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	database.DB.AutoMigrate(&models.Category{}, &models.Product{}, &models.Cart{}, &models.CartItem{})
}

func newContext(e *echo.Echo, method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// ==================== CATEGORIES API ====================

func TestAPI_CreateCategory_Success(t *testing.T) {
	setupDB(t)
	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/categories", `{"name":"Food","description":"Jedzenie"}`)

	controllers.CreateCategory(c)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var cat models.Category
	json.Unmarshal(rec.Body.Bytes(), &cat)
	assert.Equal(t, "Food", cat.Name)
}

func TestAPI_CreateCategory_BadRequest(t *testing.T) {
	setupDB(t)
	e := echo.New()
	// Niepoprawny JSON
	c, rec := newContext(e, http.MethodPost, "/categories", `{invalid}`)

	controllers.CreateCategory(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAPI_GetCategories_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Category{Name: "A"})
	database.DB.Create(&models.Category{Name: "B"})

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/categories", "")

	controllers.GetCategories(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var cats []models.Category
	json.Unmarshal(rec.Body.Bytes(), &cats)
	assert.Equal(t, 2, len(cats))
}

func TestAPI_GetCategories_Empty(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/categories", "")

	controllers.GetCategories(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var cats []models.Category
	json.Unmarshal(rec.Body.Bytes(), &cats)
	assert.Equal(t, 0, len(cats))
}

func TestAPI_GetCategory_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Category{Name: "Test"})

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/categories/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.GetCategory(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAPI_GetCategory_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/categories/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.GetCategory(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var body map[string]string
	json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Contains(t, body["error"], "not found")
}

func TestAPI_UpdateCategory_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Category{Name: "Old"})

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/categories/1", `{"name":"New"}`)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.UpdateCategory(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var cat models.Category
	json.Unmarshal(rec.Body.Bytes(), &cat)
	assert.Equal(t, "New", cat.Name)
}

func TestAPI_UpdateCategory_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/categories/999", `{"name":"X"}`)
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.UpdateCategory(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_DeleteCategory_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Category{Name: "Del"})

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/categories/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.DeleteCategory(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAPI_DeleteCategory_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/categories/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.DeleteCategory(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ==================== PRODUCTS API ====================

func TestAPI_CreateProduct_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Category{Name: "Cat"})

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/products", `{"name":"Prod","price":99.99,"quantity":10,"category_id":1}`)

	controllers.CreateProduct(c)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var p models.Product
	json.Unmarshal(rec.Body.Bytes(), &p)
	assert.Equal(t, "Prod", p.Name)
	assert.Equal(t, 99.99, p.Price)
}

func TestAPI_CreateProduct_BadRequest(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/products", `{bad json}`)

	controllers.CreateProduct(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAPI_GetProducts_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	database.DB.Create(&models.Product{Name: "P1", Price: 10, CategoryID: cat.ID})

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/products", "")

	controllers.GetProducts(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.GreaterOrEqual(t, len(products), 1)
}

func TestAPI_GetProduct_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	database.DB.Create(&models.Product{Name: "P1", Price: 10, CategoryID: cat.ID})

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/products/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.GetProduct(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var p models.Product
	json.Unmarshal(rec.Body.Bytes(), &p)
	assert.Equal(t, "P1", p.Name)
}

func TestAPI_GetProduct_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/products/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.GetProduct(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_UpdateProduct_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	database.DB.Create(&models.Product{Name: "Old", Price: 10, CategoryID: cat.ID})

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/products/1", `{"name":"Updated","price":20}`)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.UpdateProduct(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var p models.Product
	json.Unmarshal(rec.Body.Bytes(), &p)
	assert.Equal(t, "Updated", p.Name)
}

func TestAPI_UpdateProduct_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/products/999", `{"name":"X"}`)
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.UpdateProduct(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_DeleteProduct_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	database.DB.Create(&models.Product{Name: "Del", Price: 10, CategoryID: cat.ID})

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/products/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.DeleteProduct(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAPI_DeleteProduct_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/products/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.DeleteProduct(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ==================== CARTS API ====================

func TestAPI_CreateCart_Success(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/carts", "")

	controllers.CreateCart(c)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var cart models.Cart
	json.Unmarshal(rec.Body.Bytes(), &cart)
	assert.NotZero(t, cart.ID)
}

func TestAPI_GetCart_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Cart{})

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/carts/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.GetCart(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAPI_GetCart_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodGet, "/carts/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.GetCart(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_AddCartItem_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	database.DB.Create(&models.Product{Name: "P", Price: 10, Quantity: 5, CategoryID: cat.ID})
	database.DB.Create(&models.Cart{})

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/carts/1/items", `{"product_id":1,"quantity":2}`)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.AddCartItem(c)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var item models.CartItem
	json.Unmarshal(rec.Body.Bytes(), &item)
	assert.Equal(t, 2, item.Quantity)
}

func TestAPI_AddCartItem_CartNotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/carts/999/items", `{"product_id":1,"quantity":1}`)
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.AddCartItem(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_AddCartItem_ProductNotFound(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Cart{})

	e := echo.New()
	c, rec := newContext(e, http.MethodPost, "/carts/1/items", `{"product_id":999,"quantity":1}`)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.AddCartItem(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_UpdateCartItem_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	prod := models.Product{Name: "P", Price: 10, Quantity: 5, CategoryID: cat.ID}
	database.DB.Create(&prod)
	cart := models.Cart{}
	database.DB.Create(&cart)
	item := models.CartItem{CartID: cart.ID, ProductID: prod.ID, Quantity: 1}
	database.DB.Create(&item)

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/carts/1/items/1", `{"quantity":10}`)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "1")

	controllers.UpdateCartItem(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.CartItem
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, 10, result.Quantity)
}

func TestAPI_UpdateCartItem_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodPut, "/carts/1/items/999", `{"quantity":5}`)
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "999")

	controllers.UpdateCartItem(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_DeleteCartItem_Success(t *testing.T) {
	setupDB(t)
	cat := models.Category{Name: "C"}
	database.DB.Create(&cat)
	prod := models.Product{Name: "P", Price: 10, Quantity: 5, CategoryID: cat.ID}
	database.DB.Create(&prod)
	cart := models.Cart{}
	database.DB.Create(&cart)
	item := models.CartItem{CartID: cart.ID, ProductID: prod.ID, Quantity: 1}
	database.DB.Create(&item)

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/carts/1/items/1", "")
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "1")

	controllers.DeleteCartItem(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAPI_DeleteCartItem_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/carts/1/items/999", "")
	c.SetParamNames("id", "itemId")
	c.SetParamValues("1", "999")

	controllers.DeleteCartItem(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAPI_DeleteCart_Success(t *testing.T) {
	setupDB(t)
	database.DB.Create(&models.Cart{})

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/carts/1", "")
	c.SetParamNames("id")
	c.SetParamValues("1")

	controllers.DeleteCart(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAPI_DeleteCart_NotFound(t *testing.T) {
	setupDB(t)

	e := echo.New()
	c, rec := newContext(e, http.MethodDelete, "/carts/999", "")
	c.SetParamNames("id")
	c.SetParamValues("999")

	controllers.DeleteCart(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
