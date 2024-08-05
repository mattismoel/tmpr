package model

type Forecast struct {
	Location Location `json:"location"`
	Weather  Weather  `json:"weather"`
}

type Weather struct {
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feelsLike"`
	TempMin     float64 `json:"tempMin"`
	TempMax     float64 `json:"tempMax"`
	Pressure    int64   `json:"pressure"`
	Humidity    int64   `json:"humidity"`
	Description string  `json:"description"`
}
