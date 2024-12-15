package kafkaConsumeTools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"NotificationService/internal/models"
)

type Consumer struct {
	client *kgo.Client
	topic  string
	topicOutput *os.File
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	groupID := uuid.New().String()
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)

	if err != nil {
		return nil, err
	}

	file, err := os.Create(topic + ".txt")
	if err != nil{
        return nil, err
    }
    defer file.Close()

	return &Consumer{client: client, topic: topic, topicOutput: file}, nil
}

func (c *Consumer) PrintMessages() (error) {
	ctx := context.Background()
	for {
		fetches := c.client.PollFetches(ctx)
		if err := fetches.Errors();  err != nil {
			return fmt.Errorf("error in fetching %v", err)
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var msg models.Message
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}
			c.topicOutput.WriteString("Send to " + msg.Email + "\n" + msg.Text + "\n")
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}