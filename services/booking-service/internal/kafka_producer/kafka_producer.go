package kafka_producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
	}
	return err
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
