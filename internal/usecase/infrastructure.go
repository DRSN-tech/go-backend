package usecase

import "context"

type MlServiceInfra interface {
	VectorizeRequest(ctx context.Context, req *VectorizeReq) ([]VectorizeRes, error)
}

type ImagesInfra interface {
	UploadImages(ctx context.Context, req *UploadImagesReq) (*UploadImagesRes, error)
	CleanupImages(keys []string)
}

type MessageProducer interface {
	WriteMessage(ctx context.Context, req *WriteMessageReq) error
	GetPayloadBytes(req *WriteMessageReq) ([]byte, error)
	WriteRawMessage(ctx context.Context, req *WriteRawMessageReq) error
}
