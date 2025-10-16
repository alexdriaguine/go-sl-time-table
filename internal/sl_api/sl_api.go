package sl_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SLClient interface {
	GetDepartures(siteId int) ([]MappedSLDeparture, error)
}

type SLApi struct {
	httpClient *http.Client
	baseUrl    string
}

// Ensure implementing interface
var _ SLClient = (*SLApi)(nil)

func NewSLApi(httpClient *http.Client, baseUrl string) *SLApi {
	fmt.Println(baseUrl)
	return &SLApi{httpClient, baseUrl}
}

func NewDefaultSLApi() *SLApi {
	var baseUrl = "https://transport.integration.sl.se/v1"
	return NewSLApi(
		&http.Client{Timeout: 5 * time.Second},
		baseUrl,
	)
}

func (s *SLApi) GetDepartures(siteId int) ([]MappedSLDeparture, error) {

	res, err := s.httpClient.Get(fmt.Sprintf("%s/sites/%d/departures", s.baseUrl, siteId))

	if err != nil {
		return nil, fmt.Errorf("error getting departures from sl, %v", err)
	}
	defer res.Body.Close()

	var departures SLApiDepartures
	err = json.NewDecoder(res.Body).Decode(&departures)

	if err != nil {
		return nil, fmt.Errorf("error decoding json for departures, %v", err)
	}

	mappedDepartures := mapDepartures(departures)
	return mappedDepartures, nil
}

func mapDepartures(d SLApiDepartures) []MappedSLDeparture {
	mappedDepartures := []MappedSLDeparture{}

	for _, d := range d.Departures {
		mappedDepartures = append(mappedDepartures, MappedSLDeparture{
			Destination:   d.Destination,
			Display:       d.Display,
			LineNumber:    d.Line.ID,
			TransportMode: d.Line.TransportMode,
			GroupOfLines:  d.Line.GroupOfLines,
			State:         d.State,
		})
	}

	return mappedDepartures
}
