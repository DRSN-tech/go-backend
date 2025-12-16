package domain

// Image описывает изображение, которое хранится в S3
type Image struct {
	ID        string // uuid
	Bucket    string
	ObjectKey string
	Bytes     []byte
	// Передайте значение -1 в Size, если размер потока неизвестен
	// (внимание: при передаче значения -1 будет выделен большой объем памяти).
	Size     *int64
	MimeType *string // Example: "application/text"
}

func NewImage(id string, bucket string, objectKey string, bytes []byte, size *int64, mimeType *string) *Image {
	return &Image{
		ID:        id,
		Bucket:    bucket,
		ObjectKey: objectKey,
		Bytes:     bytes,
		Size:      size,
		MimeType:  mimeType,
	}
}
