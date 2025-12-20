//go:generate goverter gen github.com/DRSN-tech/go-backend/internal/repository/pgdb/converter
package converter

import (
	"time"

	"github.com/DRSN-tech/go-backend/internal/domain"
)

// ProductConverter преобразует сущности Product между domain и моделью PostgreSQL.
// goverter:converter
// goverter:extend ConvertTime
// goverter:extend ConvertPointerTime
type ProductConverter interface {
	ToModel(entity *domain.Product) *ProductModel
	ToEntity(model *ProductModel) *domain.Product
}

// CategoryConverter преобразует сущности Category между domain и моделью PostgreSQL.
// goverter:converter
// goverter:extend ConvertTime
// goverter:extend ConvertPointerTime
type CategoryConverter interface {
	ToModel(entity *domain.Category) *CategoryModel
	ToEntity(model *CategoryModel) *domain.Category
}

func ConvertPointerTime(t *time.Time) *time.Time {
	return t
}

func ConvertTime(t time.Time) time.Time {
	return t
}
