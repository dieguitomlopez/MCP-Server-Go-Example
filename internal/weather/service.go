package weather

import (
	"context"
	"weather-mcp/internal/nws"
)

type Service struct {
	nwsClient *nws.Client
}

func NewService(client *nws.Client) *Service {
	return &Service{nwsClient: client}
}

func (s *Service) GetForecast(ctx context.Context, lat, lon float64) ([]Forecast, error) {
	//1. LLamar a NWS para obtener los puntos
	data, err := s.nwsClient.GetForecast(lat, lon)

	if err != nil {
		return nil, err
	}

	//2. Mapear a weather forecast
	var results []Forecast

	for _, p := range data.Properties.Periods {
		results = append(results, Forecast{
			Name:             p.Name,
			Temperature:      p.Temperature,
			Unit:             p.TemperatureUnit,
			DetailedForecast: p.DetailedForecas,
		})
	}
	return results, nil
}

func (s *Service) GetAlerts(ctx context.Context, state string) ([]Alert, error) {
	data, err := s.nwsClient.GetAlerts(state)

	if err != nil {
		return nil, err
	}
	var results []Alert
	for _, p := range data.Features {
		results = append(results, Alert{
			Event:       p.Properties.Event,
			Severity:    p.Properties.Severity,
			Description: p.Properties.Description,
		})
	}
	return results, nil
}
