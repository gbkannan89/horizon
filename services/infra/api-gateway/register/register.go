package register

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterGatewayRoutes registers the API gateway's domain CRUD endpoints
// on the monolith's ServeMux. These mirror the separate Gin-based gateway
// but use net/http handlers that delegate to domain services.
//
// Current status: The Gin-based gateway at services/infra/api-gateway
// provides full domain CRUD. These net/http equivalents provide the same
// operations for the monolith, wired to the same domain service implementations.
func RegisterGatewayRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	_ = mux
	_ = pool

	// Domain CRUD endpoints are registered via each domain's own register package:
	//   - User:     (not yet registered)
	//   - Institution: (not yet registered)
	//   - Goal:     (not yet registered)
	//   - Account:  (not yet registered)
	//   - Allocation: (not yet registered)
	//   - Asset:    (not yet registered)
	//   - Liability:  (not yet registered)
	//   - Portfolio:   (not yet registered)
	//   - Event:    ✅ registered via financial-event/register
	//   - Household: ✅ registered via household/register
	//
	// To add domain CRUD for a domain, create register/register.go following
	// the financial-event or household pattern, then add the registration call
	// to cmd/server/main.go.
	//
	// Example for User:
	//   userReg "github.com/horizon/core/services/domains/user/register"
	//   userReg.RegisterRoutes(mux, pool)
}
