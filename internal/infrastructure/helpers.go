package infrastructure

import "github.com/DRSN-tech/go-backend/pkg/e"

// GetExtensionFromMIME возвращает расширение файла по MIME-типу изображения.
// Поддерживает jpeg, jpg, png, webp. Возвращает ошибку e.ErrUnsupportedMediaType для неподдерживаемых типов.
func GetExtensionFromMIME(mime string) (string, error) {
	switch mime {
	case "image/jpeg", "image/jpg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/webp":
		return "webp", nil
	default:
		return "bin", e.ErrUnsupportedMediaType
	}
}
