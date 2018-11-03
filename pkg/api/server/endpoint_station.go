package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/paulbes/go-pedal/pkg/api"
)

// Note: the endpoint layer receives the decoded request from the transport
// layer and passes it on to the service layer, i.e., it knows what parameters
// are required by the corresponding service entrypoint.

func makeGetStationEndpoint(s api.StationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int)
		return s.Get(ctx, req)
	}
}

func makeListStationEndpoint(s api.StationService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return s.List(ctx)
	}
}
