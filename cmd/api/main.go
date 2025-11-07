package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/user/coc/internal/auth"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/order"
	"github.com/user/coc/internal/user"
	"github.com/user/coc/pkg/config"
	"github.com/user/coc/pkg/router"
	"github.com/user/coc/pkg/validation"
)

func main() {

	// Setup logger (level can be controlled with LOG_LEVEL env var)
	level := slog.LevelInfo
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))); v != "" {
		switch v {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			// unknown value -> keep default (info)
		}
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	// use config.go to load .env and validate config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("starting application", "port", cfg.Port)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection pool
	pool, err := db.NewPool(ctx, cfg.DatabaseURL, cfg.DBMaxConnection)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize queries
	queries := db.New(pool)

	// Initialize validator with i18n support (en, zh)
	validator := validation.New()

	// Parse bearer token duration
	bearerTokenDuration, err := time.ParseDuration(cfg.BearerTokenDuration)
	if err != nil {
		slog.Error("invalid BEARER_TOKEN_DURATION format", "error", err)
		os.Exit(1)
	}

	// Initialize services
	userService := user.NewService(queries)
	orderService := order.NewService(queries)
	authService := auth.NewService(queries, cfg.JWTSecret, bearerTokenDuration)

	// Initialize handlers (pass shared validator instance)
	userHandler := user.NewHandler(userService, validator)
	orderHandler := order.NewHandler(orderService, validator)
	authHandler := auth.NewHandler(authService, validator)

	// Initialize auth middleware
	authMiddleware := auth.Middleware(authService, queries)

	// Setup router
	r := router.New(userHandler, orderHandler, authHandler, authMiddleware)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
