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

	"github.com/user/coc/internal/app/address"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/app/admin_management"
	"github.com/user/coc/internal/app/admin_menu"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/app/user_auth"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/config"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/middleware"
	"github.com/user/coc/internal/router"
	"github.com/user/coc/internal/validation"
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
	auditService := audit.NewService(queries)

	// User auth service (for frontend API)
	authService := user_auth.NewService(queries, auditService, cfg.JWTSecret, bearerTokenDuration)
	authHandler := user_auth.NewHandler(authService, validator)

	// User services (for frontend and admin)
	userService := user.NewService(queries, auditService)
	userAdminHandler := user.NewAdminHandler(userService, validator)
	userFrontendHandler := user.NewFrontendHandler(userService, validator)

	// Address services (separate for frontend and admin)
	addressAdminService := address.NewAdminService(queries, auditService)
	addressUserService := address.NewUserService(queries, auditService)
	addressAdminHandler := address.NewAdminHandler(addressAdminService, validator)
	addressFrontendHandler := address.NewFrontendHandler(addressUserService, validator)

	// Admin authentication service and handler (for admin login)
	adminAuthService := admin_auth.NewAuthService(queries, cfg.JWTSecret)
	adminAuthHandler := admin_auth.NewAuthHandler(adminAuthService, validator)

	// Admin CRUD service and handler (for managing admins)
	adminService := admin_management.NewService(queries, auditService)
	adminHandler := admin_management.NewHandler(adminService, validator)

	// Menu handler (for serving admin menu)
	menuHandler := admin_menu.NewHandler(queries)

	// Initialize middleware
	// User auth middleware (for frontend API)
	userAuthMiddleware := middleware.Middleware(authService, queries)

	// Admin auth middleware (for admin API)
	adminAuthMiddleware := middleware.AdminAuthMiddleware(adminAuthService)

	// Permission middleware (for granular access control)
	permissionMiddleware := middleware.NewPermissionMiddleware(queries)

	// Setup router with separate admin and frontend handlers
	r := router.New(
		userAdminHandler,
		userFrontendHandler,
		addressAdminHandler,
		addressFrontendHandler,
		authHandler,
		adminAuthHandler,
		adminHandler,
		menuHandler,
		userAuthMiddleware,
		adminAuthMiddleware,
		permissionMiddleware,
	)

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
