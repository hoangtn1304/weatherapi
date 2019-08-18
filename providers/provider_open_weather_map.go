package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OpenWeatherMapProvider provider for openweathermap.org
type OpenWeatherMapProvider struct {
	APIKey string
	URL    string
}

// OpenWeatherMapData data received from openweathermap.org
type OpenWeatherMapData struct {
	Current struct {
		KelvinTemp float64 `json:"temp"`
	} `json:"main"`
}

// GetTemperature implementation of WeatherProvider interface
func (p OpenWeatherMapProvider) GetTemperature(city string) (float64, error) {
	res, err := http.Get(p.URL + p.APIKey + "&q=" + city)

	if err != nil || res.StatusCode != 200 {
		return 0, err
	}

	defer res.Body.Close()

	data := OpenWeatherMapData{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return 0, err
	}

	tempC := data.Current.KelvinTemp - 273.15
	fmt.Println("openweathermap: ", tempC)

	return tempC, err
}
