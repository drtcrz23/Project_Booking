package main

import (
	"BookingService/internal/app"
	"BookingService/internal/handlers"
	"BookingService/internal/kafkaProduceTools"
	"BookingService/internal/repository"
	"context"
	"errors"
	"fmt"
	pb "github.com/drtcrz23/Project_Booking/services/grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
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

	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close()
	hotelClient := pb.NewHotelServiceClient(grpcConn)

	brokers := []string{"localhost:9092"}
	topic := "BookingEventsQueue"
	producer, err := kafkaProduceTools.New(brokers, topic)
	if (err != nil ) {
		log.Fatalf("Failed to create producer")
	}
	defer producer.Close()

	handler := handlers.NewHandler(db, producer, hotelClient)

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

		producer.Close();

		return nil
	})

	err = group.Wait()
	if err != nil {
		fmt.Printf("after wait: %s\n", err)
		return
	}
}
