package usecase

import "context"

type ProductTypeUC interface {
	AddNewProduct(ctx context.Context, req *AddNewProductReq) error
	GetProductsInfo(ctx context.Context, req *GetProductsReq) (*GetProductsRes, error)
}
