package usecase

import (
	"context"
	"strings"

	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	transaction "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TODO: добавить кафку и версию продукта
// ProductUseCase реализует бизнес-логику управления типами продуктов.
type ProductUseCase struct {
	productRepo   ProductRepository
	categoryRepo  CategoryRepository
	dbPool        transaction.Transactional
	mlService     MlServiceInfra
	imagesInfra   ImagesInfra
	embeddingRepo EmbeddingRepository
	logger        logger.Logger
	cacheRepo     CacheRepository
}

func NewProductUC(
	productRepo ProductRepository,
	categoryRepo CategoryRepository,
	dbPool transaction.Transactional,
	mlService MlServiceInfra,
	imagesInfra ImagesInfra,
	embeddingRepo EmbeddingRepository,
	logger logger.Logger,
	cacheRepo CacheRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		dbPool:        dbPool,
		mlService:     mlService,
		imagesInfra:   imagesInfra,
		embeddingRepo: embeddingRepo,
		logger:        logger,
		cacheRepo:     cacheRepo,
	}
}

// AddNewProduct обрабатывает добавление нового продукта с изображениями, категорией, векторами и сохранением в хранилища.
func (p *ProductUseCase) RegisterNewProduct(ctx context.Context, req *AddNewProductReq) error {
	const op = "ProductUseCase.RegisterNewProduct"

	// Валидация данных
	var err error
	err = p.validateProduct(req)
	if err != nil {
		return e.Wrap(op, err)
	}

	var (
		imagesRes *UploadImagesRes
		uploaded  bool
	)

	ctx, tx, err := transaction.NewTransaction(ctx, pgx.TxOptions{}, p.dbPool)
	if err != nil {
		return e.Wrap(op, err)
	}
	// Если произошла ошибка, происходит Rollback транзакции и очистка загруженных изображений
	defer func() {
		if err != nil {
			if tx.IsActive() {
				tx.Rollback(ctx)
			}

			if uploaded && imagesRes != nil {
				p.logger.Warnf(
					"Cleaning up orphaned images after transaction failure. product_name: %s, error: %v",
					req.Name,
					err,
				)

				p.imagesInfra.CleanupImages(imagesRes.ImagesKeys)
			}
		}
	}()
	ctx = context.WithValue(ctx, "tx", tx.Transaction())

	// идемпотентное создание категории
	category, err := p.createCategory(ctx, req.Name)
	if err != nil {
		return err
	}

	// идемпотентное создание продукта
	product, err := p.upsertProduct(ctx, req.Name, req.Price, category.ID)
	if err != nil {
		return err
	}

	// Отправка изображение на ML Service для получения векторов
	vectors, err := p.getVectors(ctx, req.Images)
	if err != nil {
		return err
	}

	// Сохранение изображений в MinIO
	imagesRes, err = p.uploadImages(ctx, req.Name, req.Images)
	if err != nil {
		return err
	}
	uploaded = true

	// Сохранение векторов с дополнительной информацией (S3 key, Product ID, Created At, Model Version)
	err = p.upsertEmbeddings(ctx, product.ID, imagesRes.ImagesKeys, vectors)
	if err != nil {
		return err
	}

	// Коммит изменений в бд
	err = tx.Commit(ctx)
	if err != nil {
		return e.Wrap(op, err)
	}

	// Удаление из кэша старых данных товара
	if err := p.cacheRepo.DeleteProducts(ctx, []int64{product.ID}); err != nil {
		p.logger.Warnf("Failed to delete products: %v", e.Wrap(op, err))
	}

	return nil
}

// GetProductsInfo возвращает информацию о продуктах по их идентификаторам.
func (p *ProductTypeUseCase) GetProductsInfo(ctx context.Context, req *GetProductsReq) (*GetProductsRes, error) {
	const op = "ProductTypeUseCase.GetProductsInfo"

	productsInfo, err := p.getProductsInfo(ctx, req.IDs)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	result := make([]ProductInfo, 0, len(productsInfo))
	for _, pr := range productsInfo {
		result = append(result, NewProductInfo(pr.ID, pr.Name, pr.CategoryName, pr.Price))
	}

	return &GetProductsRes{Products: result}, nil
}

// getProductsInfo делегирует запрос репозиторию продуктов.
func (p *ProductTypeUseCase) getProductsInfo(ctx context.Context, ids []int64) ([]ProductInfo, error) {
	return p.productRepo.GetProductsInfo(ctx, ids)
}

// getVectors запрашивает векторные представления изображений у ML-сервиса.
func (p *ProductTypeUseCase) getVectors(ctx context.Context, images []ProductImage) ([]VectorizeRes, error) {
	const mockImageType = 1 // TODO: убрать заглушку на реальную переменную, PROTO-46

	vectors, err := p.mlService.VectorizeRequest(ctx, NewVectorizeReq(images, mockImageType))
	if err != nil {
		return nil, err
	}

	if len(vectors) == 0 {
		return nil, e.ErrEmptyVectors
	}

	return vectors, nil
}

// upsertProductType идемпотентно создаёт или обновляет тип продукта.
func (p *ProductTypeUseCase) upsertProductType(ctx context.Context, name string, price int64, categoryID int64) (*domain.ProductType, error) {
	return p.productRepo.Upsert(ctx, domain.NewProductType(name, price, categoryID))
}

// createCategory идемпотентно создаёт категорию по имени продукта.
func (p *ProductTypeUseCase) createCategory(ctx context.Context, categoryName string) (*domain.Category, error) {
	return p.categoryRepo.Create(ctx, domain.NewCategory(categoryName))
}

// uploadImages сохраняет изображения продукта в MinIO.
func (p *ProductTypeUseCase) uploadImages(ctx context.Context, name string, images []ProductImage) (*UploadImagesRes, error) {
	return p.imagesInfra.UploadImages(ctx, NewUploadImagesReq(name, images))
}

// upsertEmbeddings сохраняет векторы изображений в Qdrant с привязкой к продукту и объектам MinIO.
func (p *ProductTypeUseCase) upsertEmbeddings(ctx context.Context, productID int64, imageKeys []string, vectors []VectorizeRes) error {
	if len(imageKeys) != len(vectors) {
		return e.ErrImageVectorMismatch
	}

	embeddings := make([]domain.Embedding, 0, len(imageKeys))
	for i, key := range imageKeys {
		if len(vectors[i].Vector) == 0 {
			return e.ErrVectorEmbeddingEmpty
		}
		payload := domain.NewPayload(productID, key, vectors[i].ModelVersion)
		embeddings = append(embeddings, *domain.NewEmbedding(uuid.NewString(), vectors[i].Vector, payload))
	}

	return p.embeddingRepo.Upsert(ctx, embeddings)
}

// validateProduct проверяет корректность входных данных запроса на добавление продукта.
func (p *ProductTypeUseCase) validateProduct(req *AddNewProductReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return e.ErrProductNameRequired
	}

	if req.Price <= 0 {
		return e.ErrPriceMustBePositive
	}

	if len(req.Images) == 0 {
		return e.ErrNoImages
	}

	return nil
}
