package createrequest

import (
  "net/http"
  "github.com/go-chi/render"
  "github.com/aidosgal/prichal/internal/lib/logger/sl"
  "github.com/go-chi/chi/v5/middleware"
  "log/slog"
)

type Request struct {
  UserID int `json:"user_id"`
  Title string `json:"title"`
  Description string `json:"description"`
  Location string `json:"location"`
  CategoryID int `json:"category_id"`
  SubcategoryID int `json:"subcategory_id"`
  SpecializationID int `json:"specialization_id"`
}

type Response struct {
  Status int `json:"status"`
  Message string `json:"message"`
  Error string `json:"error,omitempty"`
}

type RequestCreater interface {
  CreateRequest(userID int, title string, description string, location string, categoryID int, subcategoryID int, specializationID int) error
}

func New(log *slog.Logger, rc RequestCreater) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.requests.createrequest.New"

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

    err = rc.CreateRequest(req.UserID, req.Title, req.Description, req.Location, req.CategoryID, req.SubcategoryID, req.SpecializationID)

    if err != nil {
      log.Error("Failed to create request", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось создать запрос", Error: err.Error()})
      return
    }
    
    log.Info("Request created", slog.Any("request", req))
  }
}
