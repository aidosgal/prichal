package createactivity

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
  ImageUrl string `json:"image_url"`
  SubcategoryID int `json:"subcategory_id"`
  SpecializationID int `json:"specialization_id"`
}

type Response struct {
  Status int `json:"status"`
  Message string `json:"message"`
  Error string `json:"error,omitempty"`
}

type ActivityCreater interface {
  CreateActivity(userID int, title string, description string, location string, categoryID int, imageUrl string, subcategoryID int, specializationID int) error 
}

func New(log *slog.Logger, ac ActivityCreater) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.activities.createactivity.New"

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

    err = ac.CreateActivity(req.UserID, req.Title, req.Description, req.Location, req.CategoryID, req.ImageUrl, req.SubcategoryID, req.SpecializationID)

    if err != nil {
      log.Error("Failed to create activity", sl.Err(err))
      render.JSON(w, r, Response{Status: http.StatusInternalServerError, Message: "Не удалось создать активность", Error: err.Error()})
      return
    }
    
    log.Info("Activity created", slog.Any("request", req))
    render.JSON(w, r, Response{Status: http.StatusOK, Message: "Активность создана"})
  }
}
