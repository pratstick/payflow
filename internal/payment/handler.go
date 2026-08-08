package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/payments", h.createPayment)
	mux.HandleFunc("GET /v1/payments/", h.getPayment)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) createPayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: ErrInvalidRequest.Error()})
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	payment, err := h.service.CreatePayment(r.Context(), req, idempotencyKey)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, payment)
}

func (h *Handler) getPayment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/payments/")
	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: ErrInvalidRequest.Error()})
		return
	}

	payment, err := h.service.GetPayment(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, ErrMissingIdempotencyKey):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, ErrDuplicateIdempotencyKey):
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, ErrPaymentNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, ErrInvalidStateTransition):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
