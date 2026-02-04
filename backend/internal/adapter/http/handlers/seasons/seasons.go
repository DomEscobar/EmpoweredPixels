package seasons

import (
	"encoding/json"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	seasonsusecase "empoweredpixels/internal/usecase/seasons"
)

type Handler struct {
	service *seasonsusecase.Service
}

func NewHandler(service *seasonsusecase.Service) *Handler {
	return &Handler{service: service}
}

type seasonSummaryDto struct {
	SeasonId int `json:"seasonId"`
	Position int `json:"position"`
}

type pagingOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type pageDto[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	Items      []T `json:"items"`
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload pagingOptions
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		payload = pagingOptions{Page: 1, PageSize: 20}
	}

	summaries, err := h.service.SummaryPage(r.Context(), userID, payload.Page, payload.PageSize)
	if err != nil {
		log.Printf("seasons summary error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	items := make([]seasonSummaryDto, 0, len(summaries))
	for _, summary := range summaries {
		items = append(items, seasonSummaryDto{
			SeasonId: summary.SeasonID,
			Position: summary.Position,
		})
	}

	responses.JSON(w, http.StatusOK, pageDto[seasonSummaryDto]{
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalCount: len(items),
		Items:      items,
	})
}
