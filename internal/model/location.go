package model

import "fmt"

type Coords struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func NewCoords(lon, lat float64) Coords {
	return Coords{Lon: lon, Lat: lat}
}

type Location struct {
	Name     string
	City     string
	Country  string
	Postcode string
	Coords   Coords
}

func (l Location) String() string {
	return fmt.Sprintf("%s, %s, %s", l.Postcode, l.City, l.Country)
}
