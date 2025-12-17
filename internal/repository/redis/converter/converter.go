package converter

import "github.com/DRSN-tech/go-backend/internal/usecase"

type ProductInfoConverter interface {
	ToRedisModel(entity *usecase.ProductInfo) *ProductInfoRedisModel
	ToUseCase(model *ProductInfoRedisModel) *usecase.ProductInfo
	ToArrRedisModel(entities []usecase.ProductInfo) []ProductInfoRedisModel
	ToArrUseCase(models []ProductInfoRedisModel) []usecase.ProductInfo
}
