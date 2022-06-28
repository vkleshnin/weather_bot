package main

import (
	"encoding/json"
	"github.com/essentialkaos/translit/v2"
	"io/ioutil"
	"net/http"
	"strings"
)

func weatherAPI(city string) (*WeatherData, error) {
	weatherToken := "167cfa7007134746a00124044222806"
	weatherApi := "https://api.weatherapi.com/v1/current.json"
	weatherUrl := weatherApi + "?key=" + weatherToken
	weatherData, err := getWeather(weatherUrl, city)
	if err != nil {
		return nil, err
	}
	return weatherData, nil
}

func getWeather(weatherUrl string, city string) (*WeatherData, error) {
	resp, err := http.Get(weatherUrl + "&q=" + city + "&aqi=no")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, err
	}
	return &weatherData, nil
}

func trueCityName(city string) string {
	buf := translit.EncodeToICAO(city)
	trueName := strings.Replace(buf, " ", "-", -1)
	return trueName
}
