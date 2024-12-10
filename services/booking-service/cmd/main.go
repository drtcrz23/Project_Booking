package main

import (
	"BookingService/internal/app"
	"BookingService/internal/handlers"
	"BookingService/internal/kafka_producer"
	"BookingService/internal/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	dbName := "booking.db"
	db, err := repository.CreateDBConnection(dbName)
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}

	defer db.Close()

	if err := repository.CreateTable(db); err != nil {
		fmt.Println("Error creating booking table:", err)
		return
	}

	// Пока добавил так, дальше займусь соединением кафки с notifications
	brokers := []string{"localhost:9092"}
	producer := kafka_producer.NewKafkaProducer(brokers, "booking_event")
	defer producer.Close()

	handler := handlers.NewHandler(db, producer)

	mux := http.NewServeMux()
	mux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.AddBooking(w, r)
		case "PATCH":
			handler.SetBooking(w, r)
		case "DELETE":
			handler.DeleteBooking(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/bookings", handler.GetBookings)
	mux.HandleFunc("/bookings/user", handler.GetBookingByUser)

	err_env := app.LoadEnv()
	if err_env != nil {
		fmt.Println(err_env)
	}

	httpServer := http.Server{
		Addr:    app.GetEnvVariable("HOST"),
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("err in listen: %s\n", err)
			return fmt.Errorf("failed to serve http server: %w", err)
		}
		fmt.Println("after listener")

		return nil
	})

	group.Go(func() error {
		fmt.Println("before ctx done")
		<-ctx.Done()
		fmt.Println("after ctx done")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		fmt.Println("after server shutdown")

		return nil
	})

	err = group.Wait()
	if err != nil {
		fmt.Printf("after wait: %s\n", err)
		return
	}
}
