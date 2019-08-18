package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIXuProvider provider for apixu.com
type APIXuProvider struct {
	APIKey string
	URL    string
}

// APIXuData data received from apixu.com
type APIXuData struct {
	Current struct {
		CelsiusTemp float64 `json:"temp_c"`
	} `json:"current"`
}

// GetTemperature implementation of WeatherProvider interface
func (p APIXuProvider) GetTemperature(city string) (float64, error) {
	res, err := http.Get(p.URL + p.APIKey + "&q=" + city)

	if err != nil || res.StatusCode != 200 {
		return 0, err
	}

	defer res.Body.Close()

	data := APIXuData{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return 0, err
	}

	fmt.Println("apixu: ", data.Current.CelsiusTemp)
	return data.Current.CelsiusTemp, err
}
