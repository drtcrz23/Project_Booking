package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"HotelService/internal/app"
	"HotelService/internal/handlers"
	"HotelService/internal/repository"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

func main() {
	dbName := "hotel.db"
	db, err := repository.CreateDBConnection(dbName)
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}
	defer db.Close()

	if err := repository.CreateTable(db); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	if err := repository.InsertInTable(db); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	handler := handlers.NewHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/hotels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			hotels, err := handler.GetAllHotels(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(hotels)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/hotel", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.AddHotel(w, r)
		case "PATCH":
			handler.SetHotel(w, r)
		case "GET":
			handler.GetHotelsByHotelier(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.AddRoom(w, r)
		case "GET":
			handler.GetRoomsByHotel(w, r)
		case "PATCH":
			handler.SetRoom(w, r)
		case "DELETE":
			handler.DeleteRoom(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err_env := app.LoadEnv()
	if err_env != nil {
		fmt.Print(err_env)
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
