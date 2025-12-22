package domain

import "time"

// Product описывает продукт
type Product struct {
	ID         int64
	Name       string
	Price      int64 // Цена хранится в копейках
	CategoryID int64
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	IsArchived bool
}

func NewProduct(name string, price int64, categoryID int64) *Product {
	return &Product{
		Name:       name,
		Price:      price,
		CategoryID: categoryID,
	}
}
