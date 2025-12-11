package domain

// Image описывает изображение, которое хранится в S3
type Image struct {
	ID        string // uuid
	Bucket    string
	ObjectKey string
	// Передайте значение -1 в Size, если размер потока неизвестен
	// (внимание: при передаче значения -1 будет выделен большой объем памяти).
	Size        *int64
	ContentType *string // Example: "application/text"
}

func NewImage(id string, bucket string, objectKey string, size *int64, contentType *string) *Image {
	return &Image{
		ID:          id,
		Bucket:      bucket,
		ObjectKey:   objectKey,
		Size:        size,
		ContentType: contentType,
	}
}
