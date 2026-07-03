package register

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/horizon/core/services/engines/documents/internal/api"
	"github.com/horizon/core/services/engines/documents/internal/engine"
)

func RegisterRoutes(mux *http.ServeMux) {
	uploadDir := os.Getenv("DOCUMENTS_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join(os.TempDir(), "horizon-documents")
	}
	store := engine.NewDocumentStore(uploadDir)
	h := api.NewHandler(store)
	mux.HandleFunc("POST /api/v1/documents/upload", h.Upload)
	mux.HandleFunc("GET /api/v1/documents", h.List)
	mux.HandleFunc("GET /api/v1/documents/{id}", h.GetByID)
	mux.HandleFunc("GET /api/v1/documents/{id}/download", h.Download)
	mux.HandleFunc("DELETE /api/v1/documents/{id}", h.Delete)
}
