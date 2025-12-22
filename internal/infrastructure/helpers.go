package infrastructure

import (
	"github.com/DRSN-tech/go-backend/internal/proto"
	"github.com/DRSN-tech/go-backend/pkg/e"
)

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

func ConvertExtensionToProtoEnum(ext string) proto.ImageType {
	switch ext {
	case "jpg":
		return proto.ImageType_IMAGE_TYPE_JPEG
	case "png":
		return proto.ImageType_IMAGE_TYPE_PNG
	case "webp":
		return proto.ImageType_IMAGE_TYPE_WEBP
	default:
		return proto.ImageType_IMAGE_TYPE_UNKNOWN
	}
}
