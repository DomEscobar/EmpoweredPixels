package version

import (
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type versionResponse struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

func (h *Handler) GetVersion(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, versionResponse{
		AppName: "EmpoweredPixels",
		Version: "1.0.0",
	})
}
