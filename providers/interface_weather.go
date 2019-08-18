package weatherapi

// WeatherProvider interface to be used by providers
type WeatherProvider interface {
	GetTemperature(city string) (float64, error)
}
