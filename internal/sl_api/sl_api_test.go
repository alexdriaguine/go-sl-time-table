package sl_api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSLApi(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(mockSLDeparturesResponse))
		}))

		slApi := sl_api.NewSLApi(server.Client(), server.URL)

		got, err := slApi.GetDepartures(sl_api.GetDeparturesArgs{SiteId: 9325})
		require.NoError(t, err)
		want := []sl_api.MappedSLDeparture{
			{
				Destination:   "Västerhaninge",
				Display:       "Nu",
				LineNumber:    43,
				TransportMode: "TRAIN",
				GroupOfLines:  "Pendeltåg",
				State:         "ATSTOP",
			},
			{
				Destination:   "Odenplan",
				Display:       "1 min",
				LineNumber:    515,
				TransportMode: "BUS",
				GroupOfLines:  "",
				State:         "EXPECTED",
			},
		}

		require.Equal(t, got, want)
	})

	t.Run("non 200 status code returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		slApi := sl_api.NewSLApi(server.Client(), server.URL)

		_, err := slApi.GetDepartures(sl_api.GetDeparturesArgs{SiteId: 9325})
		require.Error(t, err)
	})

	t.Run("can return sites", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(mockSLSitesResponse))
		}))

		slApi := sl_api.NewSLApi(server.Client(), server.URL)

		got, err := slApi.GetSites("Sundby")
		want := []sl_api.MappedSLSite{
			{
				Name: "Sundbyberg",
				Id:   9325,
				Alias: []string{
					"Sundbybergs centrum",
					"Sundbybergs station",
					"Sundbybergs torg",
				},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("incorrect transport type returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(mockSLDeparturesResponse))
		}))

		slApi := sl_api.NewSLApi(server.Client(), server.URL)

		_, err := slApi.GetDepartures(sl_api.GetDeparturesArgs{SiteId: 9325, Transport: "rocket"})

		assert.Error(t, err)
	})

}

const mockSLSitesResponse = `[
  {
    "id": 9325,
    "gid": 9091001000009325,
    "name": "Sundbyberg",
    "alias": [
      "Sundbybergs centrum",
      "Sundbybergs station",
      "Sundbybergs torg"
    ],
    "abbreviation": "SBG",
    "lat": 59.3608711069539,
    "lon": 17.9714916630653,
    "stop_areas": [
      3431,
      6031,
      12346,
      50242,
      4543
    ],
    "valid": {
      "from": "2017-10-11T00:00:00"
    }
  },
  {
    "id": 9326,
    "gid": 9091001000009326,
    "name": "Solna strand",
    "abbreviation": "SSD",
    "lat": 59.3534977796971,
    "lon": 17.9743774023631,
    "stop_areas": [
      3421,
      50053
    ],
    "valid": {
      "from": "2014-08-18T00:00:00"
    }
  },
  {
    "id": 9327,
    "gid": 9091001000009327,
    "name": "Huvudsta",
    "abbreviation": "HUV",
    "lat": 59.3496499577023,
    "lon": 17.985420470501,
    "stop_areas": [
      3411,
      12175,
      50137
    ],
    "valid": {
      "from": "2012-06-23T00:00:00"
    }
  }
]`

const mockSLDeparturesResponse = `{
  "departures": [
    {
      "destination": "Västerhaninge",
      "direction_code": 1,
      "direction": "Nynäshamn",
      "state": "ATSTOP",
      "display": "Nu",
      "scheduled": "2025-10-15T20:11:00",
      "expected": "2025-10-15T20:11:00",
      "journey": {
        "id": 2025101502865,
        "state": "NORMALPROGRESS",
        "prediction_state": "NORMAL"
      },
      "stop_area": {
        "id": 6031,
        "name": "Sundbyberg",
        "type": "RAILWSTN"
      },
      "stop_point": {
        "id": 6031,
        "name": "Sundbyberg",
        "designation": "3"
      },
      "line": {
        "id": 43,
        "designation": "43",
        "transport_authority_id": 1,
        "transport_mode": "TRAIN",
        "group_of_lines": "Pendeltåg"
      },
      "deviations": []
    },
    {
      "destination": "Odenplan",
      "direction_code": 2,
      "direction": "Odenplan",
      "state": "EXPECTED",
      "display": "1 min",
      "scheduled": "2025-10-15T20:13:00",
      "expected": "2025-10-15T20:13:00",
      "journey": {
        "id": 2025101500140,
        "state": "EXPECTED"
      },
      "stop_area": {
        "id": 12346,
        "name": "Sundbybergs station",
        "type": "BUSTERM"
      },
      "stop_point": {
        "id": 50439,
        "name": "Sundbybergs station",
        "designation": "A"
      },
      "line": {
        "id": 515,
        "designation": "515",
        "transport_authority_id": 1,
        "transport_mode": "BUS"
      },
      "deviations": []
    }
  ],
  "stop_deviations": []
}`
