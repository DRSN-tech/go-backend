package domain

import "time"

// Payload описывает дополнительную информацию вектора
type Payload map[string]any

// Embedding представляет эмбеддинг одного изображения/объекта
type Embedding struct {
	ID      string
	Vector  []float32
	Payload Payload
}

func NewEmbedding(id string, vector []float32, payload Payload) *Embedding {
	return &Embedding{
		ID:      id,
		Vector:  vector,
		Payload: payload,
	}
}

func NewPayload(productID int64, imagePath string, modelVersion string) Payload {
	return Payload{
		"product_id":    productID,
		"image_path":    imagePath,
		"created_at":    time.Now().UTC().UnixNano(),
		"model_version": modelVersion,
	}
}
