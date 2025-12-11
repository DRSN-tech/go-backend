package domain

// QdrantPoint описывает запись в Qdrant
type QdrantPoint struct {
	ID       uint64
	Vectors  []float32
	Payloads map[string]any
}

func NewQdrantPoint(id uint64, vectors []float32, payloads map[string]any) *QdrantPoint {
	return &QdrantPoint{
		ID:       id,
		Vectors:  vectors,
		Payloads: payloads,
	}
}
