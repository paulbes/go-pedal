package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/paulbes/go-pedal/pkg/api"
)

// Endpoints contains all available endpoints for this API.
// Endpoints provide a useful abstraction, because they
// can be reused by other implementations, such as
// GraphQL, etc.
type Endpoints struct {
	GetStation  endpoint.Endpoint
	ListStation endpoint.Endpoint
}

// MakeEndpoints initialises the endpoints
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		GetStation:  makeGetStationEndpoint(s.Station),
		ListStation: makeListStationEndpoint(s.Station),
	}
}

// Handlers contains all available handlers for this API.
// What handler is invoked and when is setup by the router.
type Handlers struct {
	GetStation  http.Handler
	ListStation http.Handler
}

// MakeHandlers initialises the handlers with decoders, encoders, etc.
func MakeHandlers(e Endpoints, serverOptions ...kithttp.ServerOption) *Handlers {
	newServer := func(e endpoint.Endpoint, decodeRequestFn kithttp.DecodeRequestFunc) http.Handler {
		return kithttp.NewServer(
			e,
			decodeRequestFn,
			kithttp.EncodeJSONResponse,
			serverOptions...,
		)
	}

	return &Handlers{
		GetStation:  newServer(e.GetStation, decodeGetStationRequest),
		ListStation: newServer(e.ListStation, kithttp.NopRequestDecoder),
	}
}

// AttachRoutes creates a router and adds the handlers with the
// path, method, etc., that resolves to them.
func AttachRoutes(handlers *Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/stations", func(r chi.Router) {
			r.Method(http.MethodGet, "/{identifier}", handlers.GetStation)
			r.Method(http.MethodGet, "/", handlers.ListStation)
		})
	})

	return r
}

// Services contains all available services for this API
type Services struct {
	Station api.StationService
}
