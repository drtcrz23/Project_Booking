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

	"github.com/drtcrz23/Project_Booking/services/user-service/internal/repository"

	"github.com/drtcrz23/Project_Booking/services/user-service/internal/app"
	"github.com/drtcrz23/Project_Booking/services/user-service/internal/handlers"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

func main() {
	dbName := "user.db"
	db, err := repository.CReateDBConnection(dbName)
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}
	defer db.Close()

	if err := repository.CReateTable(db); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	handler := handlers.NewHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			users, err := handler.GetAllUsers(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(users)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.AddUser(w, r)
		case "GET":
			handler.GetUserById(w, r)
		case "PATCH":
			handler.SetUser(w, r)
		case "DELETE":
			handler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err_env := app.LoadEnv()
	if err_env != nil {
		fmt.Print(err_env)
	}

	httpServer := http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Ошибка при прослушивании: %s\n", err)
			return fmt.Errorf("не удалось запустить HTTP сервер: %w", err)
		}
		fmt.Println("После запуска HTTP-сервера")
		return nil
	})
	group.Go(func() error {
		fmt.Println("Ожидание завершения контекста...")
		<-ctx.Done()
		fmt.Println("После завершения контекста")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		fmt.Println("После завершения работы HTTP-сервера")
		return nil
	})

	err = group.Wait()
	if err != nil {
		fmt.Printf("Ошибка ожидания завершения: %s\n", err)
		return
	}
}
