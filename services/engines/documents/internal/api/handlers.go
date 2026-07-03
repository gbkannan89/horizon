package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/horizon/core/services/engines/documents/internal/engine"
)

type Handler struct {
	store *engine.DocumentStore
}

func NewHandler(store *engine.DocumentStore) *Handler {
	return &Handler{store: store}
}

func userID(r *http.Request) string {
	if u := r.URL.Query().Get("user_id"); u != "" {
		return u
	}
	return "default"
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read file error")
		return
	}

	name := header.Filename
	category := r.FormValue("category")
	description := r.FormValue("description")
	tagsStr := r.FormValue("tags")

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = detectMimeType(name)
	}

	doc, err := h.store.Save(uid, name, description, category, mimeType, data, tags)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": doc})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	docs, err := h.store.ListByUser(uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if docs == nil { docs = []*engine.Document{} }
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": docs})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	doc, err := h.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": doc})
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	doc, err := h.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	data, err := h.store.ReadFile(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read error")
		return
	}
	w.Header().Set("Content-Type", doc.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.Name))
	w.Write(data)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func detectMimeType(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".pdf": return "application/pdf"
	case ".jpg", ".jpeg": return "image/jpeg"
	case ".png": return "image/png"
	case ".doc", ".docx": return "application/msword"
	case ".xls", ".xlsx": return "application/vnd.ms-excel"
	case ".csv": return "text/csv"
	default: return "application/octet-stream"
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
