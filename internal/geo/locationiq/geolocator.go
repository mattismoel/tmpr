package locationiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mattismoel/tmpr/internal/model"
)

type locationIQGeolocator struct {
	apiKey string
}

type apiLocation struct {
	PlaceID     string     `json:"place_id,omitempty"`
	License     string     `json:"license,omitempty"`
	Lat         string     `json:"lat,omitempty"`
	Lon         string     `json:"lon,omitempty"`
	DisplayName string     `json:"display_name,omitempty"`
	Address     apiAddress `json:"address,omitempty"`
}

type apiAddress struct {
	HouseNumber string `json:"house_number,omitempty"`
	Road        string `json:"road,omitempty"`
	Suburb      string `json:"suburb,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Postcode    string `json:"postcode,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

func (l locationIQGeolocator) buildUrl(urlPath string, queryParams map[string]string) string {
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")

	u, _ := url.Parse(fmt.Sprintf("https://us1.locationiq.com/v1/%s", urlPath))

	query := u.Query()
	query.Set("key", l.apiKey)
	query.Set("format", "json")

	for key, val := range queryParams {
		query.Set(key, val)
	}

	u.RawQuery = query.Encode()

	return u.String()
}

func NewGeolocator(apiKey string) (locationIQGeolocator, error) {
	if apiKey == "" {
		return locationIQGeolocator{}, fmt.Errorf("no api key provided")
	}

	return locationIQGeolocator{apiKey: apiKey}, nil
}

func (l locationIQGeolocator) CoordsToLocation(coords model.Coords) (model.Location, error) {
	u := l.buildUrl(
		"reverse",
		map[string]string{
			"format":           "json",
			"lat":              strconv.FormatFloat(coords.Lat, 'f', -1, 64),
			"lon":              strconv.FormatFloat(coords.Lon, 'f', -1, 64),
			"limit":            "1",
			"addressdetails":   "1",
			"normalizeaddress": "1",
			"normalizecity":    "1",
		},
	)

	resp, err := http.Get(u)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not get location from API: %v", err)
	}

	defer resp.Body.Close()

	var apiLocation apiLocation

	err = json.NewDecoder(resp.Body).Decode(&apiLocation)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not decode response to struct: %v", err)
	}

	latFloat, err := strconv.ParseFloat(apiLocation.Lat, 64)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not parse latitude: %v", err)
	}

	lonFloat, err := strconv.ParseFloat(apiLocation.Lon, 64)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not parse longitude: %v", err)
	}

	location := model.Location{
		Name:     apiLocation.DisplayName,
		City:     apiLocation.Address.City,
		Country:  apiLocation.Address.Country,
		Postcode: apiLocation.Address.Postcode,
		Coords:   model.NewCoords(lonFloat, latFloat),
	}

	return location, nil

}

func (l locationIQGeolocator) QueryToLocation(query string) (model.Location, error) {
	u := l.buildUrl("search", map[string]string{
		"q":                query,
		"limit":            "1",
		"addressdetails":   "1",
		"normalizeaddress": "1",
		"normalizecity":    "1",
	})

	resp, err := http.Get(u)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not get location from API: %v", err)
	}

	defer resp.Body.Close()

	var apiLocations []apiLocation
	err = json.NewDecoder(resp.Body).Decode(&apiLocations)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not decode response into struct: %v", err)
	}

	if len(apiLocations) <= 0 {
		return model.Location{}, fmt.Errorf("no locations found")
	}

	apiLocation := apiLocations[0]

	lonFloat, err := strconv.ParseFloat(apiLocation.Lon, 64)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not parse longitude to float: %v", err)
	}

	latFloat, err := strconv.ParseFloat(apiLocation.Lat, 64)
	if err != nil {
		return model.Location{}, fmt.Errorf("could not parse latitude to float: %v", err)
	}

	location := model.Location{
		Name:     apiLocation.DisplayName,
		Coords:   model.NewCoords(lonFloat, latFloat),
		City:     apiLocation.Address.City,
		Country:  apiLocation.Address.Country,
		Postcode: apiLocation.Address.Postcode,
	}

	return location, nil
}
