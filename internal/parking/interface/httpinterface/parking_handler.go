package httpinterface

import (
	"encoding/json"
	"net/http"

	"github.com/Jakob-Kaae/gotth-example/internal/parking/domain/repository"
)

type ParkingHandler struct {
	repo repository.ParkingRepository
}

func NewParkingHandler(repo repository.ParkingRepository) *ParkingHandler {
	return &ParkingHandler{
		repo: repo,
	}
}

func (h *ParkingHandler) GetParkingSpots(w http.ResponseWriter, r *http.Request) {
	spots := h.repo.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}
