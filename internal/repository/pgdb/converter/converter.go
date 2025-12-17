//go:generate goverter gen go-backend/internal/repository/pgdb/converter
package converter

import (
	"github.com/DRSN-tech/go-backend/internal/domain"
)

// ProductConverter преобразует сущности Product между domain и моделью PostgreSQL.
type ProductConverter interface {
	ToModel(entity *domain.Product) *ProductModel
	ToEntity(model *ProductModel) *domain.Product
}

// CategoryConverter преобразует сущности Category между domain и моделью PostgreSQL.
type CategoryConverter interface {
	ToModel(entity *domain.Category) *CategoryModel
	ToEntity(model *CategoryModel) *domain.Category
}
