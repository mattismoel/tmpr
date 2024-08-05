/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"time"

	"github.com/mattismoel/env"
	"github.com/mattismoel/tmpr/internal/forecast"
	"github.com/mattismoel/tmpr/internal/forecast/openweather"
	"github.com/mattismoel/tmpr/internal/geo/locationiq"
	"github.com/spf13/cobra"
)

var (
	query string
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query for the forecast at the given location.",
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

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		fc, err := forecaster.ForecastAtQuery(ctx, query)
		if err != nil {
			log.Fatalf("could not get forecast: %v", err)
		}

		printer.Print(fc)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&query, "query", "q", "Copenhagen", "The query for the forecast. E.g. a city or specific location.")
	queryCmd.MarkFlagRequired("query")
}
