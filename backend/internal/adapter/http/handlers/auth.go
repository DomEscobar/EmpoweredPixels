package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/identity"
)

type AuthHandler struct {
	service *identity.Service
}

func NewAuthHandler(service *identity.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type refreshRequest struct {
	UserID  int64  `json:"userId"`
	Refresh string `json:"refresh"`
}

func (h *AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	var payload loginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	output, err := h.service.Token(r.Context(), identity.LoginInput{
		User:     payload.User,
		Password: payload.Password,
	})
	if err != nil {
		if err == identity.ErrInvalidCredentials {
			responses.Error(w, http.StatusBadRequest, "invalid credentials")
			return
		}
		log.Printf("token error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, output)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	var payload refreshRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	output, err := h.service.Refresh(r.Context(), payload.UserID, payload.Refresh)
	if err != nil {
		if err == identity.ErrInvalidRefresh {
			responses.Error(w, http.StatusBadRequest, "invalid refresh token")
			return
		}
		log.Printf("refresh error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, output)
}
