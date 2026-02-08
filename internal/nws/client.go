package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client es la estructura que representa al cliente
type Client struct {
	httpClient *http.Client
	baseUrl    string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseUrl:    "https://api.weather.gov",
	}
}

func (c *Client) fetchJSON(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("User-Agent", "WeatherMCP-Local/1.0")
	req.Header.Set("Accept", "application/geo+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

type pointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

type forecastResponse struct {
	Properties struct {
		Periods []struct {
			Name            string `json:"name"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			DetailedForecas string `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type alertResponse struct {
	Features []struct {
		Properties struct {
			Event       string `json:"event"`
			Severity    string `json:"severity"`
			Description string `json:"description"`
		} `json:"properties"`
	} `json:"features"`
}

func (c *Client) GetForecast(lat, lon float64) (forecastResponse, error) {
	//Devuelve el pronóstico del clima para las coordenadas dadas
	var pr pointsResponse
	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", c.baseUrl, lat, lon)
	err := c.fetchJSON(pointsURL, &pr)
	if err != nil {
		return forecastResponse{}, err
	}

	var fr forecastResponse
	err = c.fetchJSON(pr.Properties.Forecast, &fr)

	if err != nil {
		return forecastResponse{}, err
	}
	return fr, nil
}

func (c *Client) GetAlerts(state string) (alertResponse, error) {
	//Devuelve las alertas meteorológicas para el estado dado como un acrónimo NY o CA
	var a alertResponse
	url := fmt.Sprintf("%s/alerts/active/area/%s", c.baseUrl, state)
	err := c.fetchJSON(url, &a)
	return a, err
}
