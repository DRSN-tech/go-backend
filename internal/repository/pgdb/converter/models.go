package converter

import "time"

// ProductTypeModel представляет запись таблицы product_types в PostgreSQL.
type ProductTypeModel struct {
	ID         int64      `db:"id"`
	Name       string     `db:"name"`
	Price      int64      `db:"price"`
	CategoryID int64      `db:"category_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	IsArchived bool       `db:"is_archived"`
}

// CategoryModel представляет запись таблицы categories в PostgreSQL.
type CategoryModel struct {
	ID         int64      `db:"id"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	IsArchived bool       `db:"is_archived"`
}
