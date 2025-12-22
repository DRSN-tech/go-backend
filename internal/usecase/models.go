package usecase

import (
	"time"

	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/google/uuid"
)

// TODO: Task REFACTORING-49
// PRODUCT USECASE

// AddNewProductReq — запрос на добавление нового продукта.
type AddNewProductReq struct {
	Name         string
	CategoryName string
	Price        int64
	Images       []ProductImage
}

// ProductImage представляет изображение, загруженное через multipart/form-data.
type ProductImage struct {
	Data     []byte // байты изображения
	MimeType string // Content-Type из multipart (image/jpeg)
	Size     int64  // фактический размер в байтах
	Name     string // оригинальное имя файла (для логов)
}

// GetProductsReq запрос информации о продуктах по их идентификаторам.
type GetProductsReq struct {
	IDs []int64
}

// GetProductsRes — ответ с данными запрошенных продуктов.
type GetProductsRes struct {
	Products         []ProductInfo
	NotFoundProducts []int64
}

// ProductInfo — DTO с информацией о продукте для внешнего использования.
type ProductInfo struct {
	ID           int64
	Name         string
	CategoryName string
	Price        int64
}

// INFRASTUCTURE

type OutboxStatus string

const (
	Pending    OutboxStatus = "pending"
	Failed     OutboxStatus = "failed"
	Processed  OutboxStatus = "processed"
	Processing OutboxStatus = "processing"
)

type OutboxEventType string

const (
	ProductEvent OutboxEventType = "product_event"
)

type OutboxEvent struct {
	ID                  int64
	EventID             uuid.UUID
	ProductID           int64
	EventType           OutboxEventType
	Payload             []byte
	Status              OutboxStatus
	CreatedAt           time.Time
	ProcessingStartedAt *time.Time
	ProcessedAt         *time.Time
}

type WriteMessageReq struct {
	ProductID  int64
	Embeddings []domain.Embedding
}

type WriteRawMessageReq struct {
	ProductID int64
	Payload   []byte
}

// VectorizeReq — запрос на векторизацию изображений.
type VectorizeReq struct {
	Images []ProductImage
}

// VectorizeRes — результат векторизации одного изображения.
type VectorizeRes struct {
	Vector       []float32
	ModelVersion string
}

// TODO: пересмотреть структуру
// UploadImagesRes — результат загрузки изображений (ключи в MinIO).
type UploadImagesRes struct {
	ImagesKeys []string
}

// UploadImagesReq — запрос на загрузку изображений продукта.
type UploadImagesReq struct {
	Name   string
	Images []ProductImage
}

// REPOSITORIES

type UpsertProductRes struct {
	Product   *domain.Product
	NoChanges bool
}

// MAPPERS
func NewUpsertProductRes(product *domain.Product, noChanges bool) *UpsertProductRes {
	return &UpsertProductRes{
		Product:   product,
		NoChanges: noChanges,
	}
}

func NewProductInfo(id int64, name string, category string, price int64) ProductInfo {
	return ProductInfo{
		ID:           id,
		Name:         name,
		CategoryName: category,
		Price:        price,
	}
}

func NewVectorizeRes(vector []float32, modelVersion string) *VectorizeRes {
	return &VectorizeRes{
		Vector:       vector,
		ModelVersion: modelVersion,
	}
}

func NewUploadImagesReq(name string, images []ProductImage) *UploadImagesReq {
	return &UploadImagesReq{
		Name:   name,
		Images: images,
	}
}

func NewUploadImagesRes(imagesKeys []string) *UploadImagesRes {
	return &UploadImagesRes{
		ImagesKeys: imagesKeys,
	}
}

func NewVectorizeReq(images []ProductImage) *VectorizeReq {
	return &VectorizeReq{
		Images: images,
	}
}

func NewAddNewProductReq(name string, category string, price int64, images []ProductImage) *AddNewProductReq {
	return &AddNewProductReq{
		Name:         name,
		CategoryName: category,
		Price:        price,
		Images:       images,
	}
}

func NewProductImage(data []byte, mimeType string, size int64, name string) *ProductImage {
	return &ProductImage{
		Data:     data,
		MimeType: mimeType,
		Size:     size,
		Name:     name,
	}
}

func NewGetProductsRes(pr []ProductInfo, notFoundProducts []int64) *GetProductsRes {
	return &GetProductsRes{
		Products:         pr,
		NotFoundProducts: notFoundProducts,
	}
}

func NewGetProductsReq(ids []int64) *GetProductsReq {
	return &GetProductsReq{ids}
}

func NewWriteMessageReq(productID int64, embeddings []domain.Embedding) *WriteMessageReq {
	return &WriteMessageReq{
		ProductID:  productID,
		Embeddings: embeddings,
	}
}

func NewWriteRawMessageReq(productID int64, payload []byte) *WriteRawMessageReq {
	return &WriteRawMessageReq{
		ProductID: productID,
		Payload:   payload,
	}
}

func NewOutboxEvent(eventID uuid.UUID, productID int64, eventType OutboxEventType, payload []byte) *OutboxEvent {
	return &OutboxEvent{
		EventID:   eventID,
		ProductID: productID,
		EventType: eventType,
		Payload:   payload,
		Status:    Pending,
		CreatedAt: time.Now(),
	}
}
