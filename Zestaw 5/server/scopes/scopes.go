package scopes

import (
	"strconv"

	"gorm.io/gorm"
)

// ByCategory filters products by category ID.
func ByCategory(categoryID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ?", categoryID)
	}
}

// PriceRange filters products within a price range.
func PriceRange(min, max float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if min > 0 && max > 0 {
			return db.Where("price BETWEEN ? AND ?", min, max)
		}
		if min > 0 {
			return db.Where("price >= ?", min)
		}
		if max > 0 {
			return db.Where("price <= ?", max)
		}
		return db
	}
}

// InStock returns only products with quantity > 0.
func InStock() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("quantity > 0")
	}
}

// Paginate returns a scope for pagination.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// WithCategory preloads the Category relation.
func WithCategory() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Category")
	}
}

// ParseFloat safely parses a string to float64, returning 0 on error.
func ParseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// ParseInt safely parses a string to int, returning 0 on error.
func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

// ParseUint safely parses a string to uint, returning 0 on error.
func ParseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
