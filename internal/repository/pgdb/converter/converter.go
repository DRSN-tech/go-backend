//go:generate goverter gen go-backend/internal/repository/pgdb/converter
package converter

import (
	"github.com/DRSN-tech/go-backend/internal/domain"
)

// ProductTypeConverter преобразует сущности ProductType между domain и моделью PostgreSQL.
type ProductTypeConverter interface {
	ToModel(entity *domain.ProductType) *ProductTypeModel
	ToEntity(model *ProductTypeModel) *domain.ProductType
}

// CategoryConverter преобразует сущности Category между domain и моделью PostgreSQL.
type CategoryConverter interface {
	ToModel(entity *domain.Category) *CategoryModel
	ToEntity(model *CategoryModel) *domain.Category
}
