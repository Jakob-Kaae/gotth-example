package service

import (
	"github.com/Jakob-Kaae/gotth-example/internal/parking/domain/repository"
)

type PollParkingSpotsService struct {
	api        ParkingAPI
	repository repository.ParkingRepository
}

func NewPollParkingSpotsService(api ParkingAPI, repository repository.ParkingRepository) *PollParkingSpotsService {
	return &PollParkingSpotsService{
		api:        api,
		repository: repository,
	}
}

func (s *PollParkingSpotsService) Poll() error {
	spots, err := s.api.FetchParkingSpots()
	if err != nil {
		return err
	}
	return s.repository.SaveAll(spots)
}
