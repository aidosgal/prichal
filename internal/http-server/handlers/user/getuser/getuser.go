package getuser

import (
  "github.com/aidosgal/prichal/internal/models/user"
  "net/http"
  "log/slog"
  "github.com/go-chi/render"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5/middleware"
)

type Response struct {
  User    user.User `json:"user"`
  Status  int       `json:"status"`
  Message string    `json:"message"`
  Error   string    `json:"error"`
}

type Request struct {
  ChatID int `json:"chat_id"`
}

type UserGetter interface {
  GetUserByChatID(chatID int) (user.User, error)
}

func New(log *slog.Logger, ug UserGetter) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.users.getuser.New"

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

    user, err := ug.GetUserByChatID(req.ChatID)

    if err != nil {
      log.Error("Failed to get user", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось получить пользователя", Error: err.Error()})
      return
    }

    log.Info("User fetched", slog.Any("user", user))

    render.JSON(w, r, Response{User: user, Status: http.StatusOK, Message: "Пользователь получен успешно"})
  }
}
