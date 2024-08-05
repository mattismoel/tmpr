package forecast

import "github.com/mattismoel/tmpr/internal/model"

type Forecaster interface {
	ForecastAtCoords(model.Coords) (model.Forecast, error)
	ForecastAtQuery(string) (model.Forecast, error)
}
