package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Route describes a single HTTP route mapping method, path, and handler.
type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []func(http.Handler) http.Handler
}

// RegisterRoutes registers a slice of Route structs into a chi.Router.
func RegisterRoutes(r chi.Router, routes []Route) {
	for _, route := range routes {
		var h http.Handler = route.Handler
		// Chain middleware from last to first
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			h = route.Middleware[i](h)
		}
		r.Method(route.Method, route.Path, h)
	}
}
