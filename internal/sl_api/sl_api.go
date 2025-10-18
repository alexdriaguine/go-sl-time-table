package sl_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexdriaguine/go-sl-time-table/internal/cache"
	"github.com/alexdriaguine/go-sl-time-table/internal/utils"
)

type SLClient interface {
	GetDepartures(int) ([]MappedSLDeparture, error)
	GetSites(string) ([]MappedSLSite, error)
}

type SLApi struct {
	httpClient      *http.Client
	baseUrl         string
	sitesCache      cache.Cacher[string, []MappedSLSite]
	departuresCache cache.Cacher[string, []MappedSLDeparture]
}

// Ensure implementing interface
var _ SLClient = (*SLApi)(nil)

func NewSLApi(httpClient *http.Client, baseUrl string) *SLApi {
	sitesCache := cache.NewCache[string, []MappedSLSite]()
	departuresCache := cache.NewCache[string, []MappedSLDeparture]()

	return &SLApi{httpClient, baseUrl, sitesCache, departuresCache}
}

const baseUrl = "https://transport.integration.sl.se/v1"
const defaultTimeout = 10 * time.Second

func NewDefaultSLApi() *SLApi {
	slApi := NewSLApi(
		&http.Client{Timeout: defaultTimeout},
		baseUrl,
	)

	log.Println("warming up sites cache")
	_, err := slApi.GetSites("")

	if err != nil {
		log.Println("error fetching sites for cache..")
	}

	return slApi
}

const departuresCacheTiime = 5 * time.Second

func (s *SLApi) GetDepartures(siteId int) ([]MappedSLDeparture, error) {

	cacheKey := fmt.Sprintf("%d-site", siteId)
	cached, found := s.departuresCache.Get(cacheKey)

	if found {
		return cached, nil
	}

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
	s.departuresCache.Set(cacheKey, mappedDepartures, departuresCacheTiime)

	return mappedDepartures, nil
}

const sitesCacheKey = "sites"
const sitesCacheTime = 10 * time.Second

func (s *SLApi) GetSites(searchTerm string) ([]MappedSLSite, error) {
	cachedSites, found := s.sitesCache.Get(sitesCacheKey)
	if found {
		log.Println("cache hit: sites")
		return filterSitesBySearchTerm(cachedSites, searchTerm), nil
	}

	log.Println("cache miss: sites")
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

	s.sitesCache.Set(sitesCacheKey, mappedSites, sitesCacheTime)
	filteredSites := filterSitesBySearchTerm(mappedSites, searchTerm)

	return filteredSites, nil
}

func mapSites(sites []SLApiSite) []MappedSLSite {
	mapSite := func(s SLApiSite) MappedSLSite {

		return MappedSLSite{
			Name: s.Name,
			Id:   s.ID,
			// copies the slice in to an empty slice, to avoid having
			// null values in the json response
			Alias: append([]string{}, s.Alias...),
		}
	}
	return utils.Map(sites, mapSite)
}

func filterSitesBySearchTerm(sites []MappedSLSite, searchTerm string) []MappedSLSite {
	filterSite := func(s MappedSLSite) bool {

		nameMatches := strings.Contains(strings.ToLower(s.Name), strings.ToLower(searchTerm))

		if nameMatches {
			return true
		}

		anyAliasMatches := utils.Filter(s.Alias, func(s string) bool {
			return strings.Contains(strings.ToLower(s), strings.ToLower(searchTerm))
		})

		return len(anyAliasMatches) > 0
	}

	return utils.Filter(sites, filterSite)
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
