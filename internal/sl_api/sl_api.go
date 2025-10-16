package sl_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexdriaguine/go-sl-time-table/internal/utils"
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

	var d SLApiDepartures
	err = json.NewDecoder(res.Body).Decode(&d)

	if err != nil {
		return nil, fmt.Errorf("error decoding json for departures, %v", err)
	}

	mappedDepartures := mapDepartures(d.Departures)
	return mappedDepartures, nil
}

func (s *SLApi) GetSites(searchTerm string) ([]MappedSLSite, error) {
	res, err := s.httpClient.Get(fmt.Sprintf("%s/sites", s.baseUrl))

	if err != nil {
		return nil, fmt.Errorf("error getting sites from sl, %v", err)
	}

	defer res.Body.Close()

	var sites []SLApiSite
	err = json.NewDecoder(res.Body).Decode(&sites)

	if err != nil {
		return nil, fmt.Errorf("error decoding sites to json %v", err)
	}

	mappedSites := mapSites(sites)

	return mappedSites, nil
}

func mapSites(sites []SLApiSite) []MappedSLSite {
	mapSite := func(s SLApiSite) MappedSLSite {
		return MappedSLSite{
			Name:  s.Name,
			Id:    s.ID,
			Alias: s.Alias,
		}
	}
	return utils.Map(sites, mapSite)

}

func mapDepartures(departures []SLApiDeparture) []MappedSLDeparture {

	mapDeparture := func(d SLApiDeparture) MappedSLDeparture {
		return MappedSLDeparture{
			Destination:   d.Destination,
			Display:       d.Display,
			LineNumber:    d.Line.ID,
			TransportMode: d.Line.TransportMode,
			GroupOfLines:  d.Line.GroupOfLines,
			State:         d.State,
		}
	}
	return utils.Map(departures, mapDeparture)
}
