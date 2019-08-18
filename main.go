package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	wea "./providers"
	"github.com/gorilla/mux"
)

// ProviderList list of providers
type ProviderList []wea.WeatherProvider

// TemperatureData data received from providers
type TemperatureData struct {
	CityName       string  `json:"city_name"`
	CelsiusTemp    float64 `json:"celsius_temp"`
	KelvinTemp     float64 `json:"kelvin_temp"`
	FahrenheitTemp float64 `json:"fahrenheit_temp"`
}

func (list ProviderList) temperature(city string) float64 {

	chanTemp := make(chan float64)
	chanErr := make(chan error)

	for _, p := range list {
		go func(w wea.WeatherProvider) {
			temp, err := w.GetTemperature(city)
			if err != nil {
				chanErr <- err
				return
			}

			chanTemp <- temp
		}(p)
	}

	total := 0.0
	k := 0

	for i := 0; i < len(list); i++ {
		select {
		case temp := <-chanTemp:
			if temp > 0 {
				total += temp
				k++
			}
		case err := <-chanErr:
			panic(err)
		}
	}

	return total / float64(k)
}

func main() {

	openWeatherMap := wea.OpenWeatherMapProvider{
		APIKey: "5a42b06419f0c239f3bfb37ad803de1f",
		URL:    "http://api.openweathermap.org/data/2.5/weather?appid=",
	}

	apiXu := wea.APIXuProvider{
		APIKey: "1757ef313e404233a78165058191708",
		URL:    "http://api.apixu.com/v1/current.json?key=",
	}

	weatherBit := wea.WeatherBitProvider{
		APIKey: "acf9d480c4994a089f0f771d05544ab8",
		URL:    "http://api.weatherbit.io/v2.0/current?key=",
	}

	providerList := ProviderList{
		openWeatherMap,
		apiXu,
		weatherBit,
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/temperature/{city}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		city := vars["city"]

		tempC := providerList.temperature(city)
		tempK := tempC + 273.15
		tempF := (tempC * 1.8) + 32

		data := TemperatureData{
			CityName:       city,
			CelsiusTemp:    tempC,
			KelvinTemp:     tempK,
			FahrenheitTemp: tempF,
		}

		fmt.Printf("Temperature of %s is %f Celsius, %f Kelvin, %f Fahrenheit\n\n", city, tempC, tempK, tempF)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}).Methods("GET")

	port := 9000
	fmt.Printf("Server is listening at port: %d\n", port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), r))
}
