package weather

//Forecast representa el resultado del pronóstico del clima a enviar al LLM
type Forecast struct {
	Name             string
	Temperature      int
	Unit             string
	ShortForecast    string
	DetailedForecast string
}

//Alert representa una alerta meteorológica a enviar al LLM
type Alert struct {
	Event       string
	Severity    string
	Description string
}
