package getusersbyrequest

import (
  "github.com/aidosgal/prichal/internal/models/user"
  "net/http"
  "log/slog"
  "github.com/go-chi/render"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5/middleware"
)

type Response struct {
  Users   []user.User `json:"users"`
  Status  int         `json:"status"`
  Message string      `json:"message"`
  Error   string      `json:"error"`
}

type Request struct {
  Title string `json:"title"`
  Location string `json:"location"`
}

type UsersGetterByRequest interface {
  GetUsersByRequest(title, location string) ([]user.User, error)
}

func New(log *slog.Logger, ug UsersGetterByRequest) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.users.getusersbyrequest.New"

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

    users, err := ug.GetUsersByRequest(req.Title, req.Location)

    if err != nil {
      log.Error("Failed to get users", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось получить пользователей", Error: err.Error()})
      return
    }

    log.Info("Users fetched", slog.Any("users", users))

    render.JSON(w, r, Response{Users: users, Status: http.StatusOK, Message: "Пользователи получены успешно"})
  }
}
