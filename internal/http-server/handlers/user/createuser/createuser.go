package createuser

import (
  "net/http"
  "log/slog"
  "github.com/go-chi/render"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5/middleware"
)

type Request struct {
  Username string `json:"username"`
  ChatID int `json:"chat_id"`
  Name string `json:"name"`
  ImageURL string `json:"image_url"`
}

type Response struct {
  Status int `json:"status"`
  Message string `json:"message"`
  Error string `json:"error,omitempty"`
}

type UserCreator interface {
  CreateUser(username string, chatId int, name string, imageUrl string) error
}

func New(log *slog.Logger, uc UserCreator) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.users.create.New"

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

    err = uc.CreateUser(req.Username, req.ChatID, req.Name, req.ImageURL)

    if err != nil {
      log.Error("Failed to create user", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось создать пользователя", Error: err.Error()})
      return
    }

    log.Info("User created", slog.Any("request", req))
    render.JSON(w, r, Response{Status: http.StatusOK, Message: "Пользователь успешно создан"})
  }
}
