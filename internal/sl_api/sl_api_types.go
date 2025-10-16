package sl_api

type MappedSLDeparture struct {
	Destination   string
	Display       string
	LineNumber    int
	TransportMode string
	GroupOfLines  string
	State         string
}

type MappedSLSite struct {
	Id    int
	Name  string
	Alias []string
}

// Types from the SL API Response
type SLApiSite struct {
	ID        int     `json:"id"`
	Gid       int64   `json:"gid"`
	Name      string  `json:"name"`
	Note      string  `json:"note,omitempty"`
	Lat       float64 `json:"lat,omitempty"`
	Lon       float64 `json:"lon,omitempty"`
	StopAreas []int   `json:"stop_areas"`
	Valid     struct {
		From string `json:"from"`
	} `json:"valid"`
	Abbreviation string   `json:"abbreviation,omitempty"`
	Alias        []string `json:"alias,omitempty"`
}

type SLApiDepartures struct {
	Departures     []SLApiDeparture     `json:"departures"`
	StopDeviations []SLApiStopDeviation `json:"stop_deviations"`
}

type SLApiDeparture struct {
	Destination   string                `json:"destination"`
	DirectionCode int                   `json:"direction_code"`
	Direction     string                `json:"direction"`
	State         string                `json:"state"`
	Display       string                `json:"display"`
	Scheduled     string                `json:"scheduled"`
	Expected      string                `json:"expected"`
	Journey       SLApiDepartureJourney `json:"journey"`
	StopArea      SLApiStopArea         `json:"stop_area"`
	StopPoint     SLApiStopPoint        `json:"stop_point"`
	Line          SLApiLine             `json:"line"`
	Deviations    []interface{}         `json:"deviations"`
}

type SLApiDepartureJourney struct {
	ID              int64  `json:"id"`
	State           string `json:"state"`
	PredictionState string `json:"prediction_state"`
}

type SLApiStopArea struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type SLApiStopPoint struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
}

type SLApiLine struct {
	ID                   int    `json:"id"`
	Designation          string `json:"designation"`
	TransportAuthorityID int    `json:"transport_authority_id"`
	TransportMode        string `json:"transport_mode"`
	GroupOfLines         string `json:"group_of_lines"`
}

type SLApiStopDeviation struct {
	ID              int                 `json:"id"`
	ImportanceLevel int                 `json:"importance_level"`
	Message         string              `json:"message"`
	Scope           SLApiDeviationScope `json:"scope"`
}

type SLApiDeviationScope struct {
	StopAreas  []SLApiScopeStopArea  `json:"stop_areas"`
	StopPoints []SLApiScopeStopPoint `json:"stop_points"`
	Lines      []SLApiScopeLine      `json:"lines"`
}

type SLApiScopeStopArea struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type SLApiScopeStopPoint struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
}

type SLApiScopeLine struct {
	ID                   int    `json:"id"`
	Designation          string `json:"designation"`
	TransportAuthorityID int    `json:"transport_authority_id"`
	TransportMode        string `json:"transport_mode"`
	GroupOfLines         string `json:"group_of_lines"`
}
