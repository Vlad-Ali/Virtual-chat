package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ollama/ollama/api"
)

const (
	ConfigVar = "config.yml"
)

type App struct {
	logger   *slog.Logger
	config   *Config
	handlers *Handlers
	services *Services
	repos    *Repositories
	server   *http.Server
	appCtx   context.Context
	cancel   context.CancelFunc
}

func NewApp() (*App, error) {
	setLogger()
	logger := slog.Default().With("component", "app")
	cfg, err := LoadConfig(ConfigVar)
	if err != nil {
		logger.Error("error loading config", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ollamaBaseURL, err := url.Parse(cfg.OllamaConfig.URL)
	if err != nil {
		logger.Error("error parsing ollama url", "error", err)
	}

	ollamaClient := api.NewClient(ollamaBaseURL, http.DefaultClient)

	repos := NewRepositories(cfg, ollamaClient)
	services := NewServices(repos.Provider)
	handlers, handler := NewHandlers(cfg, services.Service, ctx)

	return &App{logger: logger, config: cfg, repos: repos, services: services, handlers: handlers, server: &http.Server{Addr: cfg.HTTPConfig.Address, Handler: handler, WriteTimeout: cfg.HTTPConfig.WriteTimeout, ReadTimeout: cfg.HTTPConfig.ReadTimeout},
		appCtx: ctx, cancel: cancel}, nil
}

func (app *App) Run() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	app.logger.Info("Starting application...")

	go func() {
		app.logger.Info(fmt.Sprintf("Listening on %s", app.config.HTTPConfig.Address))
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-signalChan
	app.logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	serverErr := app.server.Shutdown(ctx)
	if serverErr != nil {
		app.logger.Error("Server shutdown failed", "error", serverErr)
		return serverErr
	}

	app.cancel()
	time.Sleep(1 * time.Second)

	app.logger.Info("Application stopped gracefully")

	return nil
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
