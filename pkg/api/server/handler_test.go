package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/paulbes/go-pedal/pkg/api/mock"
	"github.com/paulbes/go-pedal/pkg/errors"
	"github.com/sebdah/goldie"
)

func TestRoutes(t *testing.T) {
	testCases := []struct {
		Name         string
		Method       string
		Path         string
		Err          error
		ExpectCode   int
		ExpectGolden string
	}{
		{
			Name:         "Get station ok",
			Method:       http.MethodGet,
			Path:         "/v1/stations/1",
			ExpectCode:   http.StatusOK,
			ExpectGolden: "get.200",
		},
		{
			Name:         "Get station method not allowed",
			Method:       http.MethodPost,
			Path:         "/v1/stations/1",
			ExpectCode:   http.StatusMethodNotAllowed,
			ExpectGolden: "get.405",
		},
		{
			Name:         "Get station bad request",
			Method:       http.MethodGet,
			Path:         "/v1/stations/gimme",
			ExpectCode:   http.StatusBadRequest,
			ExpectGolden: "get.400",
		},
		{
			Name:         "Get station not found",
			Method:       http.MethodGet,
			Path:         "/v1/stations/1000",
			Err:          errors.New(fmt.Errorf("no such id: 1000"), "could not find station", errors.NotFound),
			ExpectCode:   http.StatusNotFound,
			ExpectGolden: "get.404",
		},
		{
			Name:         "List station ok",
			Method:       http.MethodGet,
			Path:         "/v1/stations/",
			ExpectCode:   http.StatusOK,
			ExpectGolden: "list.200",
		},
		{
			Name:         "List station method not allowed",
			Method:       http.MethodConnect,
			Path:         "/v1/stations/",
			ExpectCode:   http.StatusMethodNotAllowed,
			ExpectGolden: "list.405",
		},
		{
			Name:         "List station internal error",
			Method:       http.MethodGet,
			Path:         "/v1/stations/",
			Err:          errors.New(fmt.Errorf("uh oh"), "io error", errors.IO),
			ExpectCode:   http.StatusInternalServerError,
			ExpectGolden: "list.500",
		},
	}

	for _, tc := range testCases {
		// Here we could have created a mocked service instead,
		// but now we get to test more with fewer tests :p
		station := mock.NewStation()
		store := mock.NewStationStore(station, tc.Err)
		service := NewStationService(store)
		endpoints := MakeEndpoints(Services{
			Station: service,
		})
		handlers := MakeHandlers(endpoints)
		router := AttachRoutes(handlers)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(tc.Method, tc.Path, nil)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, tc.ExpectCode, tc.Name)
		goldie.Assert(t, tc.ExpectGolden, recorder.Body.Bytes())
	}
}
