package e

import "fmt"

var (
	// Внутренние ошибки с транзакциями
	ErrTransactionNotFound = fmt.Errorf("transaction not found")

	// Внутренние ошибки с векторами
	ErrEmptyVectors         = fmt.Errorf("empty vectors")
	ErrVectorEmbeddingEmpty = fmt.Errorf("vector embedding is empty")
	ErrImageVectorMismatch  = fmt.Errorf("image vector mismatch")

	// 400 Bad Request
	ErrProductNameRequired  = fmt.Errorf("product name is required")
	ErrPriceMustBePositive  = fmt.Errorf("price must be positive")
	ErrNoImages             = fmt.Errorf("no images provided")
	ErrUnsupportedMediaType = fmt.Errorf("unsupported media type")
)

// Wrap оборачивает ошибку
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
