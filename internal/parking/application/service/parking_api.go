package service

import "github.com/Jakob-Kaae/gotth-example/internal/parking/domain/model"

type ParkingAPI interface {
	FetchParkingSpots() ([]model.ParkingSpot, error)
}
