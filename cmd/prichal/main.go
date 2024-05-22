package main

import (
  "github.com/aidosgal/prichal/internal/config"
  "log/slog"
  "os"
  "github.com/aidosgal/prichal/internal/storage/postgre"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  mwLogger "github.com/aidosgal/prichal/internal/http-server/middleware/logger"
	"github.com/aidosgal/prichal/internal/lib/logger/handlers/slogpretty"
  "github.com/aidosgal/prichal/internal/http-server/handlers/user/create"
  "net/http"
)

const (
  envDev = "dev"
  envProd = "prod"
  envLocal = "local"
)

func main() {
  cfg := config.MustLoad()
  
  log := setupLogger(cfg.Env)

  log.Info("Starting server", slog.String("env", cfg.Env))
  log.Debug("debug messages are enabled")

  storage, err := postgre.New()
  if err != nil {
    log.Error("failed to create storage", sl.Err(err))
    os.Exit(1)
  }
  _ = storage

  router := chi.NewRouter()

  router.Use(middleware.RequestID)
  router.Use(mwLogger.New(log))
  router.Use(middleware.Recoverer)
  router.Use(middleware.URLFormat)
  
  router.Post("/users", create.New(log, storage))
  
  log.Info("Server started", slog.String("addr", cfg.Address))
  
  srv := &http.Server{
    Addr: cfg.Address,
    Handler: router,
    ReadTimeout: cfg.HTTPServer.Timeout,
    WriteTimeout: cfg.HTTPServer.Timeout,
    IdleTimeout: cfg.HTTPServer.IdleTimeout,
  }

  if err := srv.ListenAndServe(); err != nil {
    log.Error("failed to listen", sl.Err(err))
    os.Exit(1)
  }

  log.Info("Server stopped")
}

func setupLogger(env string) *slog.Logger {
  var log *slog.Logger 
  switch env {
  case envLocal:
    log = setupPrettySlog()
  case envDev:
    log = slog.New(
      slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
    )
  case envProd:
    log = slog.New(
      slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
    )
  }

  return log
}

func setupPrettySlog() *slog.Logger {
  opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
