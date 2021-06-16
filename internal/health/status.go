package health

import (
	"sync"

	"github.com/gotway/gotway/internal/model"
)

type servicesByStatus map[model.ServiceStatus][]string

type statusUpdate struct {
	serviceByStatus servicesByStatus
	mux             sync.RWMutex
}

func (s *statusUpdate) Add(status model.ServiceStatus, serviceKey string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.serviceByStatus[status] == nil {
		s.serviceByStatus[status] = []string{}
	}
	s.serviceByStatus[status] = append(s.serviceByStatus[status], serviceKey)
}

func (s *statusUpdate) Get() servicesByStatus {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.serviceByStatus
}

func NewStatusUpdate() *statusUpdate {
	return &statusUpdate{
		serviceByStatus: make(servicesByStatus),
	}
}
