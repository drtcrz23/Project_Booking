package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/drtcrz23/Project_Booking/services/notification-service/internal/kafkaConsumeTools"

	"golang.org/x/sync/errgroup"
)

func main() {
	brokers := []string{"localhost:19092"}
	topic := "BookingEventsQueue"

	c, err := kafkaConsumeTools.NewConsumer(brokers, topic)
	if err != nil {
		log.Fatalln("Failed to create consumer:", err)
	}
	defer c.Close()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(cancelCtx)
	group.Go(func() error {
		for {
			err := c.PrintMessages()
			if err != nil {
				fmt.Printf("error in handling messages: %s\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	})

	group.Go(func() error {
		<-ctx.Done()
		fmt.Println("Shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		c.Close()

		<-shutdownCtx.Done()
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("Program terminated with error: %v", err)
	}

	fmt.Println("Program gracefully terminated.")
}
