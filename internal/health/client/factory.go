package client

import (
	"sync"

	"github.com/gotway/gotway/internal/model"
)

type Factory struct {
	clientByType map[model.ServiceType]Client
	mux          sync.Mutex
}

func (f *Factory) Get(serviceType model.ServiceType, options Options) (Client, error) {
	f.mux.Lock()
	defer f.mux.Unlock()

	if client, ok := f.clientByType[serviceType]; ok {
		return client, nil
	}

	client, err := New(serviceType, options)
	if err != nil {
		return nil, err
	}
	f.clientByType[serviceType] = client

	return client, nil
}

func NewFactory() Factory {
	return Factory{
		clientByType: make(map[model.ServiceType]Client),
	}
}
