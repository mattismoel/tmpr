package geo

import "github.com/mattismoel/tmpr/internal/model"

type Geolocator interface {
	CoordsToLocation(model.Coords) (model.Location, error)
	QueryToLocation(string) (model.Location, error)
}
