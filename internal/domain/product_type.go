package domain

import "time"

// ProductType описывает тип продукта
type ProductType struct {
	ID         int64
	Name       string
	Price      int64 // Цена хранится в копейках
	CategoryID int64
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	IsArchive  bool
}

func NewProductType(name string, price int64, categoryID int64) *ProductType {
	return &ProductType{
		Name:       name,
		Price:      price,
		CategoryID: categoryID,
	}
}
