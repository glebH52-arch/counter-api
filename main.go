package main

import (
	"context"
	"counter-api/iternal/config"
	"counter-api/iternal/counter"
	"counter-api/iternal/handler"
	"counter-api/iternal/redis_client"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		return err
	}
	redisClient, err := redis_client.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDb)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer redisClient.Close()
	counterService := counter.NewCounterRedisService(redisClient)
	counterHandler := handler.NewCounterHandler(counterService)
	router := handler.NewRouter(*counterHandler)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %v", err)
		}
	}
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %v", err)
	}
	return nil
}
