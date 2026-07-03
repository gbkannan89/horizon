package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/horizon/core/services/domains/recurring/internal/application"
	"github.com/horizon/core/services/domains/recurring/internal/application/dto/command"
	"github.com/horizon/core/services/domains/recurring/internal/application/dto/query"
	"github.com/horizon/core/services/domains/recurring/internal/domain"
	"github.com/horizon/core/services/internal/auth"
)

type Handlers struct {
	service *application.Service
}

func NewHandlers(service *application.Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateRecurring(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string            `json:"name"`
		Description  string            `json:"description,omitempty"`
		Amount       int64             `json:"amount"`
		Currency     string            `json:"currency"`
		Frequency    string            `json:"frequency"`
		Interval     int               `json:"interval"`
		StartDate    string            `json:"start_date"`
		EndDate      string            `json:"end_date,omitempty"`
		EventType    string            `json:"event_type"`
		Category     string            `json:"category"`
		Source       string            `json:"source"`
		Destination  string            `json:"destination,omitempty"`
		SkipHolidays bool              `json:"skip_holidays"`
		SkipWeekends bool              `json:"skip_weekends"`
		Tags         []string          `json:"tags,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid start_date, use YYYY-MM-DD")
		return
	}
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "Invalid end_date, use YYYY-MM-DD")
			return
		}
		endDate = &parsed
	}

	cmd := command.CreateRecurringCommand{
		UserID:       auth.UserIDFromRequest(r),
		Name:         req.Name,
		Description:  req.Description,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Frequency:    domain.Frequency(req.Frequency),
		Interval:     req.Interval,
		StartDate:    startDate,
		EndDate:      endDate,
		EventTemplate: domain.EventTemplate{
			Type: req.EventType, Category: req.Category,
			Source: req.Source, Destination: req.Destination,
		},
		SkipHolidays: req.SkipHolidays,
		SkipWeekends: req.SkipWeekends,
		Tags:         req.Tags,
		Metadata:     req.Metadata,
	}

	if err := h.service.CreateRecurring(r.Context(), cmd); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "Missing id")
		return
	}
	dto, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": dto})
}

func (h *Handlers) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	dtos, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if dtos == nil {
		dtos = []query.RecurringTransactionDTO{}
	}
	h.writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": dtos})
}

func (h *Handlers) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { h.writeError(w, http.StatusBadRequest, "Missing id"); return }
	if err := h.service.Activate(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "activated"})
}

func (h *Handlers) Pause(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { h.writeError(w, http.StatusBadRequest, "Missing id"); return }
	if err := h.service.Pause(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
}

func (h *Handlers) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { h.writeError(w, http.StatusBadRequest, "Missing id"); return }
	if err := h.service.Cancel(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handlers) Archive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" { h.writeError(w, http.StatusBadRequest, "Missing id"); return }
	if err := h.service.Archive(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "archived"})
}

func (h *Handlers) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}
