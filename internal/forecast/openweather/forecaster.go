package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/mattismoel/tmpr/internal/geo"
	"github.com/mattismoel/tmpr/internal/model"
	"golang.org/x/sync/errgroup"
)

type apiConfig struct {
	Unit   string
	APIKey string
}

type apiForecast struct {
	Coord   model.Coords   `json:"coord"`
	Weather []apiWeather   `json:"weather"`
	Base    string         `json:"base"`
	Main    apiMainWeather `json:"main"`
}

type apiWeather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
}

type apiMainWeather struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int64   `json:"pressure"`
	Humidity  int64   `json:"humidity"`
}

func (cfg *apiConfig) validate() error {
	units := []string{"standard", "metric", "imperial"}

	if cfg.Unit == "" {
		cfg.Unit = "standard"
	}

	if !slices.Contains(units, cfg.Unit) {
		return fmt.Errorf("invalid unit %q. Needs to be standard, metric or imperial.", cfg.Unit)
	}

	return nil
}

func NewConfig(apiKey, unit string) (apiConfig, error) {
	cfg := apiConfig{APIKey: apiKey, Unit: unit}
	err := cfg.validate()
	if err != nil {
		return apiConfig{}, fmt.Errorf("could not validate config: %v", err)
	}

	return cfg, nil
}

type openWeatherForecaster struct {
	cfg     apiConfig
	geolctr geo.Geolocator
}

func (f openWeatherForecaster) buildURL(urlPath string, queryParams map[string]string) string {
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")

	u, _ := url.Parse(fmt.Sprintf("https://api.openweathermap.org/data/2.5/%s", urlPath))

	query := u.Query()
	query.Set("appid", f.cfg.APIKey)
	query.Set("units", f.cfg.Unit)

	for key, val := range queryParams {
		query.Set(key, val)
	}

	u.RawQuery = query.Encode()

	return u.String()
}

func NewForecaster(cfg apiConfig, geolocator geo.Geolocator) *openWeatherForecaster {
	return &openWeatherForecaster{cfg: cfg, geolctr: geolocator}
}

func (f openWeatherForecaster) ForecastAtCoords(ctx context.Context, coords model.Coords) (model.Forecast, error) {
	grp, ctx := errgroup.WithContext(ctx)

	u := f.buildURL("weather",
		map[string]string{
			"lat": strconv.FormatFloat(coords.Lat, 'f', -1, 64),
			"lon": strconv.FormatFloat(coords.Lon, 'f', -1, 64),
		})

	var apiForecast apiForecast
	grp.Go(func() error {
		resp, err := http.Get(u)
		if err != nil {
			return fmt.Errorf("could not get forecast: %v", err)
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&apiForecast)
		if err != nil {
			return fmt.Errorf("could not decode json to struct: %v", err)
		}

		return nil
	})

	var location model.Location
	grp.Go(func() error {
		var err error
		location, err = f.geolctr.CoordsToLocation(ctx, coords)
		if err != nil {
			return fmt.Errorf("could not get location from coords: %v", err)
		}

		return nil
	})

	err := grp.Wait()
	if err != nil {
		return model.Forecast{}, err
	}

	forecast := model.Forecast{
		Location: location,
		Weather: model.Weather{
			Temperature: apiForecast.Main.Temp,
			FeelsLike:   apiForecast.Main.FeelsLike,
			TempMin:     apiForecast.Main.TempMin,
			TempMax:     apiForecast.Main.TempMax,
			Pressure:    apiForecast.Main.Pressure,
			Humidity:    apiForecast.Main.Humidity,
			Description: apiForecast.Weather[0].Description,
		},
	}

	return forecast, nil
}

func (f openWeatherForecaster) ForecastAtQuery(ctx context.Context, query string) (model.Forecast, error) {
	grp, ctx := errgroup.WithContext(ctx)

	var location model.Location
	grp.Go(func() error {
		var err error
		location, err = f.geolctr.QueryToLocation(ctx, query)
		if err != nil {
			return fmt.Errorf("could not get location: %v", err)
		}

		return nil
	})

	var forecast model.Forecast
	grp.Go(func() error {
		var err error
		forecast, err = f.ForecastAtCoords(ctx, location.Coords)
		if err != nil {
			return fmt.Errorf("could not get forecast: %v", err)
		}

		return nil
	})

	err := grp.Wait()
	if err != nil {
		return model.Forecast{}, err
	}

	forecast.Location = location

	return forecast, nil
}
