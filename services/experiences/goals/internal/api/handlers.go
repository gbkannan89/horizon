package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/horizon/core/services/internal/auth"
	"time"

	"github.com/horizon/core/services/experiences/goals/internal/aggregator"
	"github.com/horizon/core/services/experiences/goals/internal/engine"
)

type Handlers struct {
	aggregator *aggregator.Aggregator
	composer   *engine.Composer
}

func New(agg *aggregator.Aggregator, comp *engine.Composer) *Handlers {
	return &Handlers{aggregator: agg, composer: comp}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/goals/experience", h.ListGoals)
	mux.HandleFunc("GET /api/v1/goals/{id}/dashboard", h.GetDashboard)
	mux.HandleFunc("GET /api/v1/goals/{id}/progress", h.GetProgress)
	mux.HandleFunc("GET /api/v1/goals/{id}/projection", h.GetProjection)
	mux.HandleFunc("GET /api/v1/goals/{id}/recommendations", h.GetRecommendations)
	mux.HandleFunc("GET /api/v1/goals/{id}/optimization", h.GetOptimization)
	mux.HandleFunc("GET /api/v1/goals/{id}/timeline", h.GetTimeline)
	mux.HandleFunc("GET /api/v1/goals/{id}/milestones", h.GetMilestones)
}

func getDefaultUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) ListGoals(w http.ResponseWriter, r *http.Request) {
	userID := getDefaultUserID(r)
	goals, err := h.aggregator.GetGoals(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AGGREGATION_ERROR", err.Error())
		return
	}
	summaries := h.composer.BuildList(engine.Inputs{Goals: goals})
	writeJSON(w, http.StatusOK, ListResponse{
		Success: true, Data: &GoalSummaryListDTO{Goals: summaries.Goals, Total: len(summaries.Goals)},
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetDashboard(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	dash := h.composer.BuildDashboard(*goal)
	writeJSON(w, http.StatusOK, DashboardResponse{
		Success: true, Data: dash,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetProgress(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	prog := h.composer.BuildProgress(*goal)
	writeJSON(w, http.StatusOK, ProgressResponse{
		Success: true, Data: prog,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetProjection(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	proj := h.composer.BuildProjection(*goal)
	writeJSON(w, http.StatusOK, ProjectionResponse{
		Success: true, Data: proj,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	recs, err := h.aggregator.GetRecs(r.Context(), goalID)
	if err != nil {
		recs = []engine.RecItem{}
	}
	view := h.composer.BuildRecommendations(*goal, recs)
	writeJSON(w, http.StatusOK, RecsResponse{
		Success: true, Data: view,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetOptimization(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	opts, err := h.aggregator.GetOpts(r.Context(), goalID)
	if err != nil {
		opts = []engine.OptItem{}
	}
	view := h.composer.BuildOptimization(*goal, opts)
	writeJSON(w, http.StatusOK, OptResponse{
		Success: true, Data: view,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	evts, err := h.aggregator.GetEvents(r.Context(), goalID)
	if err != nil {
		evts = []engine.GoalEvent{}
	}
	view := h.composer.BuildTimeline(*goal, evts)
	writeJSON(w, http.StatusOK, TimelineResponse{
		Success: true, Data: view,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetMilestones(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("id")
	if goalID == "" { writeError(w, http.StatusBadRequest, "MISSING_ID", "goal ID required"); return }
	userID := getDefaultUserID(r)

	goal, err := h.aggregator.GetGoalData(r.Context(), userID, goalID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}
	mss, err := h.aggregator.GetMilestones(r.Context(), goalID)
	if err != nil {
		mss = []engine.Milestone{}
	}
	view := h.composer.BuildMilestones(*goal, mss)
	writeJSON(w, http.StatusOK, MilestoneResponse{
		Success: true, Data: view,
		Metadata: &Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
