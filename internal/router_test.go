package gosltimetable_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gosltimetable "github.com/alexdriaguine/go-sl-time-table/internal"
	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
	"github.com/stretchr/testify/assert"
)

var (
	siteIdExists = 1337
)

type slApiClientStub struct {
	departures []sl_api.MappedSLDeparture
	sites      []sl_api.MappedSLSite
	err        error
}

func (s *slApiClientStub) GetDepartures(args sl_api.GetDeparturesArgs) ([]sl_api.MappedSLDeparture, error) {
	if s.err != nil {
		return nil, s.err
	}
	if args.SiteId == siteIdExists {
		return s.departures, nil
	}
	return []sl_api.MappedSLDeparture{}, nil
}

func (s *slApiClientStub) GetSites(searchTerm string) ([]sl_api.MappedSLSite, error) {
	return s.sites, nil
}

func TestRouter(t *testing.T) {

	t.Run("departures route with existing site", func(t *testing.T) {
		slApiMock, departuresJson := buildSLClientStub(false)
		router, _ := gosltimetable.NewRouter(slApiMock)

		request := newGetRequest(fmt.Sprintf("/api/departures/%d", siteIdExists))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := departuresJson
		got := response.Body.String()

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, got, string(want))
	})

	t.Run("departures with unkown siteId returns empty array", func(t *testing.T) {
		slApiMock, _ := buildSLClientStub(false)
		router, _ := gosltimetable.NewRouter(slApiMock)

		request := newGetRequest(fmt.Sprintf("/api/departures/%d", 404))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := "[]"
		got := response.Body.String()

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, want, got)
	})

	t.Run("returns 500 on error", func(t *testing.T) {
		slApiMock, _ := buildSLClientStub(true)
		router, _ := gosltimetable.NewRouter(slApiMock)

		request := newGetRequest(fmt.Sprintf("/api/departures/%d", siteIdExists))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("sites endpoint search", func(t *testing.T) {
		slApiMock, _ := buildSLClientStub(false)
		router, _ := gosltimetable.NewRouter(slApiMock)

		request := newGetRequest("/api/sites?term=sundby")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("content-type"))
	})

}

func buildSLClientStub(shouldError bool) (*slApiClientStub, string) {
	mockDepartures := []sl_api.MappedSLDeparture{
		{
			Destination:   "Mock Destination",
			Display:       "Nu",
			LineNumber:    123,
			TransportMode: "BUS",
			GroupOfLines:  "",
			State:         "EXPECTED",
		},
		{
			Destination:   "Mock Destination",
			Display:       "Nu",
			LineNumber:    123,
			TransportMode: "BUS",
			GroupOfLines:  "",
			State:         "EXPECTED",
		},
	}

	mockSites := []sl_api.MappedSLSite{
		{Id: 1, Name: "Sundbyberg", Alias: []string{"Sundbybergs centrum"}},
		{Id: 2, Name: "Solna", Alias: []string{"Blåkulla"}},
	}
	stub := &slApiClientStub{mockDepartures, mockSites, nil}
	if shouldError {
		stub.err = fmt.Errorf("error")
	}
	mockDeparturesJson, _ := json.Marshal(mockDepartures)

	return stub, string(mockDeparturesJson)

}

func newGetRequest(path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return req
}
