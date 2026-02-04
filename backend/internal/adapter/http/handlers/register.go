package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/identity"
)

type RegisterHandler struct {
	service *identity.Service
}

func NewRegisterHandler(service *identity.Service) *RegisterHandler {
	return &RegisterHandler{service: service}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload registerRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.service.Register(r.Context(), identity.RegisterInput{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		if err == identity.ErrUserExists {
			responses.Error(w, http.StatusBadRequest, "user already exists")
			return
		}
		log.Printf("register error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RegisterHandler) Verify(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	var payload struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || payload.Value == "" {
		var raw string
		if err := json.Unmarshal(body, &raw); err != nil || raw == "" {
			responses.Error(w, http.StatusBadRequest, "invalid payload")
			return
		}
		payload.Value = raw
	}

	if err := h.service.Verify(r.Context(), payload.Value); err != nil {
		if err == identity.ErrInvalidVerification {
			responses.Error(w, http.StatusBadRequest, "invalid verification")
			return
		}
		log.Printf("verify error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
