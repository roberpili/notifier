package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/database"
	routes "github.com/rober0xf/notifier/internal/delivery/http"
	email "github.com/rober0xf/notifier/internal/infraestructure/email"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/internal/infraestructure/scheduler"
	"github.com/rober0xf/notifier/internal/middleware"
	"github.com/rober0xf/notifier/internal/session"
)

type APIServer struct {
	addr        string
	router      http.Handler
	cancelClean context.CancelFunc
	scheduler   gocron.Scheduler
	logger      *slog.Logger
}

func NewAPIServer(addr string, logger *slog.Logger) (*APIServer, error) {
	db, err := database.InitPostgres()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	config := Config{
		BaseURL:        os.Getenv("BASE_URL"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL is not set")
	}
	if config.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is not set")
	}

	// init repos
	userRepo := postgres.NewUserRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	tokenGen := auth.NewJWTGenerator(database.JwtKey, 24)

	emailConfig := email.NewConfigFromEnv()
	emailSender, err := email.NewSender(*emailConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create email sender: %v", err)
	}

	alerter := scheduler.NewAlerter(emailSender, logger)
	sched, err := alerter.InitCron()
	if err != nil {
		return nil, fmt.Errorf("cron setup: %w", err)
	}

	// handlers
	userHandler := buildUserRoutes(userRepo, tokenGen, emailSender, config, logger)
	paymentHandler := buildPaymentRoutes(paymentRepo, userRepo, logger)
	router := routes.SetupRoutes(userHandler, paymentHandler, middleware.AuthMiddleware(tokenGen, session.SessionCookieName))

	// clean tokens
	ctx, cancel := context.WithCancel(context.Background())
	database.StartTokenCleanJob(ctx, userRepo)

	return &APIServer{
		addr:        addr,
		router:      router,
		cancelClean: cancel,
		scheduler:   sched,
		logger:      logger,
	}, nil
}

func (s *APIServer) Run() error {
	server := &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		s.logger.Info("shutting down server")

		s.cancelClean()
		if err := s.scheduler.Shutdown(); err != nil {
			s.logger.Error("scheduler shutdown failed", "error", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error("server shutdown failed", "error", err)
		}
	}()

	s.logger.Info("server running", "addr", s.addr)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v", s.addr, err)
	}

	return nil
}
