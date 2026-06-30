package handler

type Handler struct {
	// Services will be injected via DI in Phase 3.3
}

func NewHandler() *Handler {
	return &Handler{}
}
