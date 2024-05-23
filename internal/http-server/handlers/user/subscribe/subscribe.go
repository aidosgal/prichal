package subscribe

import (
  "net/http"
  "github.com/go-chi/render"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5/middleware"
  "log/slog"
)

type Request struct {
  UserID int `json:"user_id"`
  SubscriberId int `json:"subscriber_id"`
}

type Response struct {
  Status int `json:"status"`
  Message string `json:"message"`
  Error string `json:"error,omitempty"`
}

type Subscriber interface {
  Subscribe(userID int, subscriberID int) error
}

func New(log *slog.Logger, s Subscriber) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.users.subscribe.New"

    log = log.With(
      slog.String("op", op),
      slog.String("Request_id", middleware.GetReqID(r.Context())),
    )

    var req Request
    err := render.DecodeJSON(r.Body, &req)

    if err != nil {
      log.Error("Failed to decode request", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusBadRequest, Message: "Не удалось получить request", Error: err.Error()})
      return
    }

    log.Info("Request decoded", slog.Any("request", req))

    err = s.Subscribe(req.UserID, req.SubscriberId)

    if err != nil {
      log.Error("Failed to subscribe", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось подписаться", Error: err.Error()})
      return
    }
    
    log.Info("Subscribed", slog.Any("request", req))
    render.JSON(w, r, Response{Status: http.StatusOK, Message: "Подписка оформлена"})
  }
}
