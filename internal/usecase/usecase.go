package usecase

import "context"

type ProductUC interface {
	RegisterNewProduct(ctx context.Context, req *AddNewProductReq) (*OutboxEvent, error)
	GetProductsInfo(ctx context.Context, req *GetProductsReq) (*GetProductsRes, error)
}
