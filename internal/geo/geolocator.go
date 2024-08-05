package geo

import (
	"context"

	"github.com/mattismoel/tmpr/internal/model"
)

type Geolocator interface {
	CoordsToLocation(context.Context, model.Coords) (model.Location, error)
	QueryToLocation(context.Context, string) (model.Location, error)
}
