package store

import (
	"sync"

	"github.com/Jakob-Kaae/gotth-example/internal/parking/domain/model"
)

type InMemoryParkingStore struct {
	mu       sync.RWMutex
	spots    []model.ParkingSpot
	onChange func([]model.ParkingSpot)
}

func NewInMemoryParkingStore() *InMemoryParkingStore {
	return &InMemoryParkingStore{
		spots: make([]model.ParkingSpot, 0),
	}
}

func (s *InMemoryParkingStore) SaveAll(spots []model.ParkingSpot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.spots = spots
	if s.onChange != nil {
		s.onChange(spots)
	}
	return nil
}

func (s *InMemoryParkingStore) GetAll() []model.ParkingSpot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.ParkingSpot(nil), s.spots...)
}

func (s *InMemoryParkingStore) OnChange(callback func([]model.ParkingSpot)) {
	// No-op for in-memory store, as it doesn't support change notifications.
	s.onChange = callback
}
