package repository

import "github.com/Jakob-Kaae/gotth-example/internal/parking/domain/model"

type ParkingRepository interface {
	SaveAll([]model.ParkingSpot) error
	GetAll() []model.ParkingSpot
	OnChange(func([]model.ParkingSpot))
}
