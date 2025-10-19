package gosltimetable

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
)

//go:embed "static/*"
var staticFiles embed.FS

type ErrorResponse struct {
	Message string
}

type Router struct {
	http.Handler
	slClient sl_api.SLClient
}

func NewRouter(slClient sl_api.SLClient) (*Router, error) {

	isDev := os.Getenv("IS_DEV") == "true"

	router := &Router{}
	router.slClient = slClient
	handler := http.NewServeMux()

	if !isDev {
		// creats a sub fs from our embedded "static/*" folder, with
		// the "static" folder as root
		staticFs, err := fs.Sub(staticFiles, "static")
		if err != nil {
			return nil, fmt.Errorf("error parsing static files %w", err)
		}

		// converts our filesustem to a http handler
		fileServer := http.FileServerFS(staticFs)

		handler.Handle("/", fileServer)

	}

	handler.Handle("/api/departures/", http.HandlerFunc(router.handleDepartures))
	handler.Handle("/api/sites", http.HandlerFunc(router.handleSites))
	router.Handler = handler

	return router, nil
}

func (router *Router) handleDepartures(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	siteId, siteIdErr := parseSiteIdFromUrl(r.URL)
	line, lineErr := parseLineFromQuery(r.URL)
	direction, directionErr := parseDirectionFromQuery(r.URL)
	transport := strings.ToUpper(r.URL.Query().Get("transport"))

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

	json.NewEncoder(w).Encode(departures)
}

func (router *Router) handleSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	searchTerm := r.URL.Query().Get("term")

	if utf8.RuneCountInString(searchTerm) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "2 or more characters needed for search"})
		return
	}

	matchingSites, _ := router.slClient.GetSites(searchTerm)

	if len(matchingSites) > 5 {
		matchingSites = matchingSites[:5]
	}
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
