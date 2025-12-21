package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DRSN-tech/go-backend/internal/domain"
	drsnProto "github.com/DRSN-tech/go-backend/internal/proto"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/jimlawless/whereami"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	writer *kafka.Writer
	logger logger.Logger
}

func NewProducer(logger logger.Logger) (*Producer, error) {
	brokerStr := os.Getenv("KAFKA_BROKERS")
	if brokerStr == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS environment variable is required")
	}
	brokers := strings.Split(brokerStr, ",")

	topic := os.Getenv("KAFKA_TOPIC")

	if topic == "" {
		return nil, fmt.Errorf("KAFKA_TOPIC environment variable is required")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireOne,
		BatchSize:    10,
		BatchTimeout: 500 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				logger.Warnf("Kafka producer error: %s", err.Error())
			}
		},
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}, nil
}

func (p *Producer) WriteMessage(ctx context.Context, req *usecase.WriteMessageReq) error {
	value, err := p.prepareProductChangeEvent(req.ProductID, req.Embeddings)
	if err != nil {
		return e.Wrap(whereami.WhereAmI(), err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(req.ProductID, 10)),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) prepareProductChangeEvent(
	productID int64,
	embeddings []domain.Embedding,
) ([]byte, error) {
	protoEmbeddings, err := toArrProtoEmbeddings(embeddings)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	event := &drsnProto.ProductChangeEvent{
		EventId:        uuid.NewString(),
		EventTimestamp: time.Now().UnixNano(),
		Operation: &drsnProto.ProductChangeEvent_Upsert{
			Upsert: &drsnProto.UpsertEvent{
				ProductId:  productID,
				Embeddings: protoEmbeddings,
			},
		},
	}

	return proto.Marshal(event)
}

func toProtoEmbeddings(embedding domain.Embedding) (*drsnProto.Embedding, error) {
	const op = "producer.toProtoEmbedding"

	productID, ok := embedding.Payload["product_id"].(int64)
	if !ok {
		return nil, e.Wrap(op, fmt.Errorf("missing or invalid product_id in payload"))
	}

	embeddingVersion, ok := embedding.Payload["embedding_version"].(int32)
	if !ok {
		return nil, e.Wrap(op, fmt.Errorf("missing or invalid embedding_version in payload"))
	}

	imagePath, ok := embedding.Payload["image_path"].(string)
	if !ok {
		return nil, e.Wrap(op, fmt.Errorf("missing or invalid image_path in payload"))
	}

	createdAt, ok := embedding.Payload["created_at"].(int64)
	if !ok {
		return nil, e.Wrap(op, fmt.Errorf("missing or invalid created_at in payload"))
	}

	modelVersion, ok := embedding.Payload["model_version"].(string)
	if !ok {
		return nil, e.Wrap(op, fmt.Errorf("missing or invalid model_version in payload"))
	}

	return &drsnProto.Embedding{
		EmbeddingId: embedding.ID,
		Vector:      embedding.Vector,
		Metadata: &drsnProto.EmbeddingMetadata{
			ProductId:        productID,
			EmbeddingVersion: embeddingVersion,
			ImagePath:        imagePath,
			CreatedAt:        createdAt,
			ModelVersion:     modelVersion,
		},
	}, nil
}

func toArrProtoEmbeddings(embeddings []domain.Embedding) ([]*drsnProto.Embedding, error) {
	result := make([]*drsnProto.Embedding, 0, len(embeddings))
	for _, embedding := range embeddings {
		emb, err := toProtoEmbeddings(embedding)
		if err != nil {
			return nil, err
		}

		result = append(result, emb)
	}

	return result, nil
}
