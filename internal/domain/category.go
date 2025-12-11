package domain

import "time"

// Category описывает категорию продукта
type Category struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	IsActive  bool
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}
