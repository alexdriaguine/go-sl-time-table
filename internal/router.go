package gosltimetable

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
)

//go:embed "templates/*"
var templates embed.FS

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
	siteId, _ := parseSiteIdFromUrl(r.URL.Path)
	line, _ := strconv.Atoi(r.URL.Query().Get("line"))
	transport := r.URL.Query().Get("transport")
	direction, _ := strconv.Atoi(r.URL.Query().Get("direction"))

	args := sl_api.GetDeparturesArgs{
		SiteId:    siteId,
		Line:      line,
		Transport: sl_api.TransportType(transport),
		Direction: direction,
	}

	departures, err := router.slClient.GetDepartures(args)

	if err != nil {
		log.Printf("error getting departures from sl, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Message string }{Message: "Internal Server ERror"})
		return
	}

	json.NewEncoder(w).Encode(&departures)
}

func (router *Router) handleSites(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	searchTerm := r.URL.Query().Get("term")

	if len(searchTerm) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Message string }{Message: "3 or more characters needed for search"})
		return
	}

	matchingSites, _ := router.slClient.GetSites(searchTerm)
	json.NewEncoder(w).Encode(matchingSites)

}

func parseSiteIdFromUrl(url string) (int, error) {
	siteId, err := strconv.Atoi(strings.Replace(url, "/api/departures/", "", 1))
	if err != nil {
		return 0, fmt.Errorf("could not parse url %s to siteId, %w", url, err)
	}
	return siteId, nil
}
