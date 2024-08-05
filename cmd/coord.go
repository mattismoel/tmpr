/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/mattismoel/env"
	"github.com/mattismoel/tmpr/internal/forecast"
	"github.com/mattismoel/tmpr/internal/forecast/openweather"
	"github.com/mattismoel/tmpr/internal/geo/locationiq"
	"github.com/mattismoel/tmpr/internal/model"
	"github.com/spf13/cobra"
)

var (
	lon float64
	lat float64
)

// coordCmd represents the coord command
var coordCmd = &cobra.Command{
	Use:   "coord",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer := forecast.StdOutPrinter{}

		openWeatherApiKey := env.Str("OPEN_WEATHER_API_KEY", "")
		if openWeatherApiKey == "" {
			log.Fatalf("could not get API key for OpenWeather API in environment variables...")
		}

		locationIQApiKey := env.Str("LOCATION_IQ_API_KEY", "")
		if locationIQApiKey == "" {
			log.Fatalf("could not get API key for LocationIQ API in environment variables...")
		}

		openWeatherCfg, err := openweather.NewConfig(openWeatherApiKey, unit)
		if err != nil {
			log.Fatalf("could not create open weather config: %v", err)
		}

		geolocator, err := locationiq.NewGeolocator(locationIQApiKey)
		if err != nil {
			log.Fatalf("could not create geolocator instance: %v", err)
		}

		forecaster := openweather.NewForecaster(openWeatherCfg, geolocator)
		coords := model.NewCoords(lon, lat)
		fc, err := forecaster.ForecastAtCoords(coords)
		if err != nil {
			log.Fatalf("could not get forecast at coords: %v", err)
		}

		printer.Print(fc)
	},
}

func init() {
	rootCmd.AddCommand(coordCmd)

	coordCmd.Flags().Float64Var(&lon, "lon", 12.5700724, "The longitude for a given location.")
	coordCmd.Flags().Float64Var(&lat, "lat", 55.6867243, "The latitude for a given location.")

	coordCmd.MarkFlagRequired("lon")
	coordCmd.MarkFlagRequired("lat")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// coordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// coordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
