package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log/slog"

	"github.com/hasanm95/go-student-api/internal/config"
	"github.com/hasanm95/go-student-api/internal/http/handlers/student"
	"github.com/hasanm95/go-student-api/internal/storage/sqlite"

)

func main() {
	// Load config from file
	cfg := config.MustLoad()

	// Initialize database
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatalf("failed to initialize storage: %s", err)
	}

	// Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/students", student.New(storage))
	mux.HandleFunc("GET /api/students/{id}", student.GetStudentById(storage))
	mux.HandleFunc("GET /api/students", student.GetStudents(storage))
	slog.Info("server started on port", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: mux,
	}

	go func(){
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting http server: %s", err)
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}