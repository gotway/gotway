package requestcontext

import (
	"context"
	"errors"
	"net/http"

	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
)

type requestContextKey string

const (
	ingressKey  requestContextKey = "service"
	responseKey requestContextKey = "response"
)

func WithIngress(r *http.Request, ingress crdv1alpha1.IngressHTTP) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ingressKey, ingress))
}

func WithResponse(r *http.Request, res *http.Response) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), responseKey, res))
}

func GetIngress(r *http.Request) (crdv1alpha1.IngressHTTP, error) {
	ingress, ok := r.Context().Value(ingressKey).(crdv1alpha1.IngressHTTP)
	if !ok {
		return crdv1alpha1.IngressHTTP{}, errors.New("ingress not found in request context")
	}
	return ingress, nil
}

func GetResponse(r *http.Request) (*http.Response, error) {
	res, ok := r.Context().Value(responseKey).(*http.Response)
	if !ok {
		return nil, errors.New("response not found in request context")
	}
	return res, nil
}
