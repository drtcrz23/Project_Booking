package kafkaProduceTools

import (
	"context"
	"encoding/json"

	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func New(brokers []string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{client: client, topic: topic}, nil
}
func (p *Producer) SendMessage(email, text string) {
	ctx := context.Background()
	msg := model.Message{Email: email, Text: text}
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, nil)
}
func (p *Producer) Close() {
	p.client.Close()
}
