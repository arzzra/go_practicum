package main

import (
	"context"
	"fmt"
	"github.com/arzzra/go_practicum/internal/server"
	"github.com/arzzra/go_practicum/internal/server/api"
	"github.com/arzzra/go_practicum/internal/server/storage/local"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	host            = "localhost"
	shutdownTimeout = 4 * time.Second
)

func StartServer(ctx context.Context, host string, port string) error {
	storage := local.MakeMetricStorage()
	srv, err := server.MakeServer(storage)
	if err != nil {
		return err
	}

	h, err := api.MakeHandler(srv)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: h.Router,
	}

	go func() {
		<-ctx.Done()
		ctx2, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx2); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()
	var port string = "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	if err := StartServer(ctx, host, port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
