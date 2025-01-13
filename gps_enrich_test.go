package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_when_no_lat_long_skip_enrichment(t *testing.T) {
	enrich := EnrichLocation{}
	m := map[string]any{}
	err := enrich.Enrich(m)

	assert.NoError(t, err)
}
func Test_happy_path_enrich_location(t *testing.T) {
	s := httptest.NewServer(ProductionHandler())
	defer s.Close()

	enrich := EnrichLocation{
		client:  s.Client(),
		baseURL: s.URL,
	}
	/*
		Tower of london according to Google
		51°30'30.7"N 0°04'34.1"W
		51.508530, -0.076132
	*/
	expectedLat := 51.508530
	expectedLon := -0.076132
	m := map[string]any{
		"GPSLatitude":     `51° 30' 30.7"`,
		"GPSLatitudeRef":  "N",
		"GPSLongitude":    `0° 04' 34.1"`,
		"GPSLongitudeRef": "W",
	}
	err := enrich.Enrich(m)

	assert.NoError(t, err)
	require.NotEmpty(t, enrich.values[AttrLocation])
	assert.Equal(t, "Camden Town", enrich.values[AttrLocation].(map[string]any)["name"])
	assert.InDelta(t, expectedLat, enrich.values[AttrLatitude], 0.00001)
	assert.InDelta(t, expectedLon, enrich.values[AttrLongitude], 0.00001)
}

func ProductionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]any{
			"place_id":     258328585,
			"licence":      "Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright",
			"osm_type":     "node",
			"osm_id":       399607277,
			"lat":          "51.5423045",
			"lon":          "-0.1395604",
			"category":     "place",
			"type":         "town",
			"place_rank":   18,
			"importance":   0.52357864783041,
			"addresstype":  "town",
			"name":         "Camden Town",
			"display_name": "Camden Town, Greater London, England, NW1 9PJ, United Kingdom",
			"address": map[string]any{
				"town":           "Camden Town",
				"state_district": "Greater London",
				"state":          "England",
				"ISO3166-2-lvl4": "GB-ENG",
				"postcode":       "NW1 9PJ",
				"country":        "United Kingdom",
				"country_code":   "gb",
			},
		}
		buf, err := json.Marshal(payload)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write(buf)
	})
}
