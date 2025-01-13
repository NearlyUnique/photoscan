package main

import (
	"encoding/json"
	"fmt"
	"github.com/paulcager/osgridref"
	"io"
	"net/http"
	"strconv"
)

type EnrichLocation struct {
	values  map[string]any
	client  *http.Client
	baseURL string
}

func (e *EnrichLocation) Enrich(values map[string]any) error {
	if e.baseURL == "" {
		e.baseURL = "https://nominatim.openstreetmap.org"
	}
	e.values = make(map[string]any)
	var ok1, ok2 bool
	var latDegMinSec, longDegMinSec string
	latDegMinSec, ok1 = values[AttrGPSLatitude].(string)
	longDegMinSec, ok2 = values[AttrGPSLongitude].(string)
	if !ok1 || !ok2 {
		return nil
	}
	latDegMinSec += values[AttrGPSLatitudeRef].(string)
	longDegMinSec += values[AttrGPSLongitudeRef].(string)

	err := e.positionToDecimal(latDegMinSec, longDegMinSec)

	if err != nil {
		return err
	}

	return e.lookupPlace()
}

func (e *EnrichLocation) positionToDecimal(lat, long string) error {
	var err error

	var latF, lonF float64

	latF, err = osgridref.ParseDegrees(lat)
	lonF, err = osgridref.ParseDegrees(long)

	e.values[AttrLatitude] = latF
	e.values[AttrLongitude] = lonF

	return err
}

func (e *EnrichLocation) lookupPlace() error {
	url := fmt.Sprintf(e.baseURL +
		`/reverse` +
		`?lat=` + ftoa(e.values[AttrLatitude]) +
		`&lon=` + ftoa(e.values[AttrLongitude]) +
		`&format=jsonv2` +
		`&zoom=18` +
		`&addressdetails=1`,
	)
	resp, err := e.client.Get(url)
	if err != nil {
		return err
	}
	var buf []byte
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}
	e.values[AttrLocation] = data
	return nil
}

// ftoa float to ascii
func ftoa(value any) string {
	return strconv.FormatFloat(value.(float64), 'f', 5, 64)
}
