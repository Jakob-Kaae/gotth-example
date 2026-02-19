package httpclient

import (
	"encoding/json"
	"net/http"

	"github.com/Jakob-Kaae/gotth-example/internal/parking/domain/model"
)

type ParkingAPIClient struct {
	baseURL string
	http    *http.Client
}

type parkingSpotDto struct {
	ID             int    `json:"id"`
	ParkingSpot    string `json:"parkeringsplads"`
	NumberOfSpots  int    `json:"antalPladser"`
	AvailableSpots int    `json:"ledigePladser"`
	OccupiedSpots  int    `json:"optagedePladser"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
}

func NewParkingAPIClient(baseURL string, http *http.Client) *ParkingAPIClient {
	return &ParkingAPIClient{
		baseURL: baseURL,
		http:    http,
	}
}

func (c *ParkingAPIClient) FetchParkingSpots() ([]model.ParkingSpot, error) {
	// Call external API
	resp, err := c.http.Get(c.baseURL + "/api/ParkingOverview")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := json.NewDecoder(resp.Body)
	var apiResponse []parkingSpotDto
	if err := body.Decode(&apiResponse); err != nil {
		return nil, err
	}
	spots := make([]model.ParkingSpot, len(apiResponse))
	for i, s := range apiResponse {
		spots[i] = model.ParkingSpot{
			ID:        s.ID,
			Name:      s.ParkingSpot,
			FreeSpots: s.AvailableSpots,
		}
	}
	return spots, nil
}
