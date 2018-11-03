package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/paulbes/go-pedal/pkg/errors"
)

// Note: the transport layer ensures that the received request can be decoded

func decodeGetStationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	param := chi.URLParam(r, "identifier")
	id, err := strconv.Atoi(param)
	if err != nil {
		return nil, errors.New(err, "failed to convert id param to int", errors.Unmarshal)
	}
	return id, nil
}
