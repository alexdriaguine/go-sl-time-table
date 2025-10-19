package gosltimetable

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
)

//go:embed "templates/*"
var templates embed.FS

type ErrorResponse struct {
	Message string
}

type Router struct {
	http.Handler
	templ    *template.Template
	slClient sl_api.SLClient
}

func NewRouter(slClient sl_api.SLClient) (*Router, error) {
	templ, err := template.ParseFS(templates, "templates/*.gohtml")

	if err != nil {
		return nil, fmt.Errorf("error parsing templates %w", err)
	}

	router := &Router{templ: templ}
	router.slClient = slClient
	handler := http.NewServeMux()

	handler.Handle("/", http.HandlerFunc(router.handleIndex))
	handler.Handle("/api/departures/", http.HandlerFunc(router.handleDepartures))
	handler.Handle("/api/sites", http.HandlerFunc(router.handleSites))
	router.Handler = handler

	return router, nil
}

func (router *Router) handleIndex(w http.ResponseWriter, r *http.Request) {

	type IndexViewModel struct {
		Title   string
		Heading string
		Message string
	}
	router.templ.ExecuteTemplate(w, "index.gohtml", IndexViewModel{Title: "Title!", Heading: "Hello", Message: "hoho"})
}

func (router *Router) handleDepartures(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	siteId, siteIdErr := parseSiteIdFromUrl(r.URL)
	line, lineErr := parseLineFromQuery(r.URL)
	direction, directionErr := parseDirectionFromQuery(r.URL)

	if siteIdErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: siteIdErr.Error()})
		return
	}

	if lineErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: lineErr.Error()})
		return
	}

	if directionErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: directionErr.Error()})
		return
	}
	transport := strings.ToUpper(r.URL.Query().Get("transport"))

	args := sl_api.GetDeparturesArgs{
		SiteId:    siteId,
		Line:      line,
		Transport: sl_api.TransportType(transport),
		Direction: direction,
	}

	departures, err := router.slClient.GetDepartures(args)

	if err != nil {
		log.Printf("error getting departures from sl, %v", err)
		if errors.Is(err, sl_api.ErrInvalidTransportType) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal Server Error"})
		return
	}

	json.NewEncoder(w).Encode(&departures)
}

func (router *Router) handleSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	searchTerm := r.URL.Query().Get("term")

	if len(searchTerm) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "3 or more characters needed for search"})
		return
	}

	matchingSites, _ := router.slClient.GetSites(searchTerm)
	json.NewEncoder(w).Encode(matchingSites)
}

func parseLineFromQuery(url *url.URL) (int, error) {
	queryLine := url.Query().Get("line")
	if queryLine == "" {
		return 0, nil
	}

	line, err := strconv.Atoi(queryLine)
	if err != nil {
		return 0, fmt.Errorf("could not parse line from value %s, %w", queryLine, err)
	}

	return line, nil
}

func parseDirectionFromQuery(url *url.URL) (int, error) {
	queryDirection := url.Query().Get("direction")
	if queryDirection == "" {
		return 0, nil
	}

	direction, err := strconv.Atoi(queryDirection)

	if err != nil {
		return 0, fmt.Errorf("could not parse direction from value %s, %w", queryDirection, err)
	}

	return direction, nil
}

func parseSiteIdFromUrl(url *url.URL) (int, error) {
	siteId, err := strconv.Atoi(strings.Replace(url.Path, "/api/departures/", "", 1))
	if err != nil {
		return 0, fmt.Errorf("could not parse url %s to siteId, %w", url, err)
	}
	return siteId, nil
}
