package register

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/goal/internal/application"
	"github.com/horizon/core/services/domains/goal/internal/application/dto/query"
	"github.com/horizon/core/services/domains/goal/internal/domain"
	"github.com/horizon/core/services/domains/goal/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

type noopPublisher struct{}
func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool, now func() time.Time) {
	repo := persistence.NewGoalRepository(pool)
	svc := application.NewGoalService(repo, &noopPublisher{}, now)

	mux.HandleFunc("GET /api/v1/households/{id}/goals", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing household id", http.StatusBadRequest)
			return
		}

		role := auth.HouseholdRoleFromRequest(r, id)
		if role == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		limit := 25
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		cursor := r.URL.Query().Get("cursor")

		q := query.ListByHouseholdQuery{
			HouseholdID: id,
			Limit:       limit,
			Cursor:      cursor,
		}

		res, err := svc.ListByHousehold(r.Context(), q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})
}
