//go:generate goverter gen github.com/DRSN-tech/go-backend/internal/repository/pgdb/converter
package converter

import (
	"time"

	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/internal/usecase"
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

// ProductEmbeddingVersionConverter преобразует сущности ProductEmbeddingVersion между domain и моделью PostgreSQL.
// goverter:converter
// goverter:extend ConvertTime
// goverter:extend ConvertPointerTime
type ProductEmbeddingVersionConverter interface {
	ToModel(entity *domain.ProductEmbeddingVersion) *ProductEmbeddingVersionModel
	ToEntity(model *ProductEmbeddingVersionModel) *domain.ProductEmbeddingVersion
}

// OutboxEventConverter преобразует сущности OutboxEvent между usecase и моделью PostgreSQL.
// goverter:converter
// goverter:extend ConvertTime
// goverter:extend ConvertPointerTime
// goverter:extend ConvertOutBoxStatus
// goverter:extend ConvertOutboxEventType
type OutboxEventConverter interface {
	ToModel(entity *usecase.OutboxEvent) *OutboxEventModel
	ToEntity(model *OutboxEventModel) *usecase.OutboxEvent
	ToArrEntity(models []*OutboxEventModel) []*usecase.OutboxEvent
}

func ConvertPointerTime(t *time.Time) *time.Time {
	return t
}

func ConvertTime(t time.Time) time.Time {
	return t
}

func ConvertOutBoxStatus(s usecase.OutboxStatus) usecase.OutboxStatus {
	return s
}

func ConvertOutboxEventType(t usecase.OutboxEventType) usecase.OutboxEventType {
	return t
}
