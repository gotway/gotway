package middleware

import "net/http"

type Middleware interface {
	MiddlewareFunc(next http.Handler) http.Handler
}
