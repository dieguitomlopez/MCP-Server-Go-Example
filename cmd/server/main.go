package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"weather-mcp/internal/nws"
	"weather-mcp/internal/weather"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type ForecastArgs struct {
	Lat float64 `json:"latitude" jsonschema:"required,description=Latitud de la ubicación"`
	Lon float64 `json:"longitude" jsonschema:"required,description=Longitud de la ubicación"`
}

type AlertsArgs struct {
	State string `json:"state" jsonschema:"required,description=Código de estado de dos letras (ej: CA, TX)"`
}

func main() {
	nwsClient := nws.NewClient(30 * time.Second)
	weatherSvc := weather.NewService(nwsClient)
	server := mcp.NewServer(
		stdio.NewStdioServerTransport(),
		mcp.WithName("weather-go-pro"),
		mcp.WithVersion("1.0.0"),
	)

	//Tool 1: Get Forecast
	server.RegisterTool("get_forecast", "Pronóstico detallado por coordenadas",
		func(args ForecastArgs) (*mcp.ToolResponse, error) {
			forecast, err := weatherSvc.GetForecast(context.Background(), args.Lat, args.Lon)
			if err != nil {
				return nil, err
			}
			//Formateo para claude
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Pronóstico para las coordenadas %f, %f:\n\n", args.Lat, args.Lon))

			for _, f := range forecast {
				sb.WriteString(fmt.Sprintf("-***%s***: %dº%s\n %s\n\n", f.Name, f.Temperature, f.Unit, f.DetailedForecast))
			}
			return mcp.NewToolResponse(mcp.NewTextContent(sb.String())), nil
		})
	//Tool 2: Get Alerts
	server.RegisterTool("get_alerts", "Alertas meteorológicas activas por estado usando un acrónimo de dos letras del estado (ej: CA,TX)",
		func(args AlertsArgs) (*mcp.ToolResponse, error) {
			alerts, err := weatherSvc.GetAlerts(context.Background(), args.State)
			if err != nil {
				return nil, err
			}
			//Formateo para claude
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Alertas activas en %s:\n\n", args.State))
			for _, a := range alerts {
				sb.WriteString(fmt.Sprintf("### %s (%s)\n%s\n\n---\n\n", a.Event, a.Severity, a.Description))
			}
			return mcp.NewToolResponse(mcp.NewTextContent(sb.String())), nil
		})
	//Arrancamos el servidor MCP
	if err := server.Serve(); err != nil {
		log.Fatalf("MCP Server error: %v", err)
	}
	fmt.Fprintln(os.Stderr, "Servidor de clima iniciado correctamente en modo STDIO")
	// Mantener el servidor corriendo indefinidamente
	select {}
}
