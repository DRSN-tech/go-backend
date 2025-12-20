package usecase

import "context"

type MlServiceInfra interface {
	VectorizeRequest(ctx context.Context, req *VectorizeReq) ([]VectorizeRes, error)
}

type ImagesInfra interface {
	UploadImages(ctx context.Context, req *UploadImagesReq) (*UploadImagesRes, error)
	CleanupImages(keys []string)
}

type Infrastructure struct {
	MlServiceInfra MlServiceInfra
	ImagesInfra    ImagesInfra
}

func NewInfrastructure(mlServiceInfra MlServiceInfra, imagesInfra ImagesInfra) *Infrastructure {
	return &Infrastructure{
		MlServiceInfra: mlServiceInfra,
		ImagesInfra:    imagesInfra,
	}
}
