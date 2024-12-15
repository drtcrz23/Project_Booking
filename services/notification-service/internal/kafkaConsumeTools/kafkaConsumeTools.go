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
	lastOffset  kgo.Offset
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	groupID := uuid.New().String()
	newOffset = kgo.NewOffset().AtStart()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(newOffset),
	)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(topic + ".txt")
	if err != nil{
        return nil, err
    }
    defer file.Close()

	return &Consumer{client: client, topic: topic, topicOutput: file, lastOffset: newOffset}, nil
}

func (c *Consumer) PrintMessages() (error) {
	ctx := context.Background()
	for {
		fetches := c.client.PollFetches(ctx)
		if err := fetches.Errors();  err != nil {
			return fmt.Errorf("error in fetching %v", err)
		}

		iter := fetches.RecordIter()
		var latestOffset kgo.Offset

		for !iter.Done() {
			record := iter.Next()
			var msg models.Message
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}
			c.topicOutput.WriteString("Send to " + msg.Email + "\n" + msg.Text + "\n")

			latestOffset = record.Offset + 1
		}
	}

	if latestOffset > 0 {
		err := c.client.CommitOffsets(ctx, map[string]kgo.Offset{c.topic: latestOffset,})
		if err != nil {
		 return fmt.Errorf("error committing offsets: %v", err)
		}
		c.lastOffset = latestOffset
	}

	return nil
}

func (c *Consumer) Close() {
	c.client.Close()
}