package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"os/signal"
	"syscall"
	"golang.org/x/sync/errgroup"
	"NotificationService/internal/kafkaConsumeTools"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "BookingEventsQueue"

	c, err := kafkaConsumeTools.NewConsumer(brokers, topic)
	if err != nil {
		log.Fatalln("Failed to create consumer:", err)
	}
	defer c.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := c.PrintMessages();
		err != nil {
			fmt.Printf("error in handling messages: %s\n", err)
			return fmt.Errorf("failed to print messages: %w", err)
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done() 
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
	
		<-shutdownCtx.Done()
		c.Close()
	
		return nil
	})
}
