package sl_api

type MappedSLDeparture struct {
	Destination   string
	Display       string
	LineNumber    int
	TransportMode string
	GroupOfLines  string
	State         string
}

// Type from the SL API Response
type SLApiDepartures struct {
	Departures []struct {
		Destination   string `json:"destination"`
		DirectionCode int    `json:"direction_code"`
		Direction     string `json:"direction"`
		State         string `json:"state"`
		Display       string `json:"display"`
		Scheduled     string `json:"scheduled"`
		Expected      string `json:"expected"`
		Journey       struct {
			ID              int64  `json:"id"`
			State           string `json:"state"`
			PredictionState string `json:"prediction_state"`
		} `json:"journey"`
		StopArea struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"stop_area"`
		StopPoint struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Designation string `json:"designation"`
		} `json:"stop_point"`
		Line struct {
			ID                   int    `json:"id"`
			Designation          string `json:"designation"`
			TransportAuthorityID int    `json:"transport_authority_id"`
			TransportMode        string `json:"transport_mode"`
			GroupOfLines         string `json:"group_of_lines"`
		} `json:"line"`
		Deviations []interface{} `json:"deviations"`
	} `json:"departures"`
	StopDeviations []struct {
		ID              int    `json:"id"`
		ImportanceLevel int    `json:"importance_level"`
		Message         string `json:"message"`
		Scope           struct {
			StopAreas []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"stop_areas"`
			StopPoints []struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Designation string `json:"designation"`
			} `json:"stop_points"`
			Lines []struct {
				ID                   int    `json:"id"`
				Designation          string `json:"designation"`
				TransportAuthorityID int    `json:"transport_authority_id"`
				TransportMode        string `json:"transport_mode"`
				GroupOfLines         string `json:"group_of_lines"`
			} `json:"lines"`
		} `json:"scope"`
	} `json:"stop_deviations"`
}
