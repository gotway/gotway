package request

import (
	"context"
	"errors"
	"net/http"

	"github.com/gotway/gotway/internal/model"
)

const (
	serviceKey  = "service"
	responseKey = "response"
)

func WithService(r *http.Request, s model.Service) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), serviceKey, s))
}

func WithResponse(r *http.Request, res *http.Response) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), responseKey, res))
}

func GetService(r *http.Request) (model.Service, error) {
	service, ok := r.Context().Value(serviceKey).(model.Service)
	if !ok {
		return model.Service{}, errors.New("service not found in request context")
	}
	return service, nil
}

func GetResponse(r *http.Request) (*http.Response, error) {
	res, ok := r.Context().Value(responseKey).(*http.Response)
	if !ok {
		return nil, errors.New("response not found in request context")
	}
	return res, nil
}
