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
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Student API!!!!!!!!!!!!!!!!!!!"))
}

func main() {
	// Load config from file
	cfg := config.MustLoad()


	// Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
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