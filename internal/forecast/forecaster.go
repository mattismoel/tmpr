package forecast

import (
	"context"

	"github.com/mattismoel/tmpr/internal/model"
)

type Forecaster interface {
	ForecastAtCoords(context.Context, model.Coords) (model.Forecast, error)
	ForecastAtQuery(context.Context, string) (model.Forecast, error)
}
