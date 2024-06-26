package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aidosgal/prichal/internal/config"
	"github.com/aidosgal/prichal/internal/http-server/handlers/telegram"
	"github.com/aidosgal/prichal/internal/http-server/handlers/user/createuser"
  "github.com/aidosgal/prichal/internal/http-server/handlers/user/getuser"
  "github.com/aidosgal/prichal/internal/http-server/handlers/activity/createactivity"
  "github.com/aidosgal/prichal/internal/http-server/handlers/request/createrequest"
  "github.com/aidosgal/prichal/internal/http-server/handlers/request/getusersbyrequest"
  "github.com/aidosgal/prichal/internal/http-server/handlers/user/subscribe"
	mwLogger "github.com/aidosgal/prichal/internal/http-server/middleware/logger"
	sl "github.com/aidosgal/prichal/internal/lib/logger/handlers/slogpretty"
	"github.com/aidosgal/prichal/internal/storage/postgre"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tb "gopkg.in/telebot.v3"
)

const (
	envDev   = "dev"
	envProd  = "prod"
	envLocal = "local"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Initialize storage
	storage, err := postgre.New()
	if err != nil {
		log.Error("failed to create storage", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize the router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Initialize Telegram bot
	bot, err := tb.NewBot(tb.Settings{
		Token:  "7079309099:AAHmlSLSxJ9OyZRd6UElcNyv1c7AIpUWIZY",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Error("failed to create bot", slog.Any("error", err))
		os.Exit(1)
	}

	// Pass the bot and storage pool to the telegram package
	telegram.New(bot, storage.Conn())

	// Define your routes
	router.Post("/users", createuser.New(log, storage))
  router.Get("/users/{chat_id}", getuser.New(log, storage))
  router.Post("/subscribe/{user_id}", subscribe.New(log, storage))
  router.Post("/activity", createactivity.New(log, storage))
  router.Post("/request", createrequest.New(log, storage))
  router.Get("/usersbyrequest", getusersbyrequest.New(log, storage))
	router.Post("/api/telegram/webhook", telegram.HandleWebhook)

	log.Info("Server started", slog.String("addr", cfg.Address))

	// Start the server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to listen", slog.Any("error", err))
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
	opts := sl.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
