//go:generate goverter gen github.com/DRSN-tech/go-backend/internal/repository/redis/converter

package converter

import (
	"time"

	"github.com/DRSN-tech/go-backend/internal/usecase"
)

// goverter:converter
// goverter:extend ConvertTime
// goverter:extend ConvertPointerTime
type ProductInfoConverter interface {
	ToRedisModel(entity *usecase.ProductInfo) *ProductInfoRedisModel
	ToUseCase(model *ProductInfoRedisModel) *usecase.ProductInfo
	ToArrRedisModel(entities []usecase.ProductInfo) []ProductInfoRedisModel
	ToArrUseCase(models []ProductInfoRedisModel) []usecase.ProductInfo
}

func ConvertPointerTime(t *time.Time) *time.Time {
	return t
}

func ConvertTime(t time.Time) time.Time {
	return t
}
