package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
   I am well aware that I can get data for both current and daily weather
   by sending one request, but I just don't really care.
*/

const urlBase = "https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&timezone=auto"

type GeocodeResponse struct {
	Results []PlaceGeocode `json:"results"`
}

type WeatherResponse struct {
	Daily   DailyWeather   `json:"daily"`
	Current CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	Time        string  `json:"time"`
	Temperature float32 `json:"temperature_2m"`
	FeelsLike   float32 `json:"apparent_temperature"`
}

type DailyWeather struct {
	Time        []string  `json:"time"`
	Max         []float32 `json:"temperature_2m_max"`
	Min         []float32 `json:"temperature_2m_min"`
	WeatherCode []int     `json:"weather_code"`
}

type PlaceGeocode struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Country   string  `json:"country"`
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

func getJsonFromUrl(url string) []byte {
	// getting the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// checking the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	// taking the response
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr.Error())
	}

	return body
}

func getPlaceGeocode(name string) PlaceGeocode {
	fixed_name := strings.Replace(name, " ", "+", -1)

	url := "https://geocoding-api.open-meteo.com/v1/search?name=" + fixed_name + "&count=10&language=en&format=json"

	body := getJsonFromUrl(url)

	// unmarshal the json
	var geocodeResponse GeocodeResponse
	err := json.Unmarshal(body, &geocodeResponse)
	if err != nil {
		log.Fatal(err.Error())
	}

	// returning 1st element
	if len(geocodeResponse.Results) > 0 {
		return geocodeResponse.Results[0]
	}

	// returning empty if no results
	return PlaceGeocode{}
}

func getDailyWeather(latitude float32, longitude float32) DailyWeather {
	url := fmt.Sprintf(
		urlBase+"&daily=weather_code,temperature_2m_max,temperature_2m_min",
		latitude,
		longitude)

	body := getJsonFromUrl(url)

	// unmarshal the json
	var weatherResponse WeatherResponse
	err := json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatal(err.Error())
	}

	// returning the weather response
	return weatherResponse.Daily
}

func getCurrentWeather(latitude float32, longitude float32) CurrentWeather {
	url := fmt.Sprintf(
		urlBase+"&current=temperature_2m,apparent_temperature",
		latitude,
		longitude)

	body := getJsonFromUrl(url)

	var weatherResponse WeatherResponse
	err := json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatal(err.Error())
	}

	return weatherResponse.Current
}

func getEmoji(wmo int) string {
	wmoWeatherMap := map[string]string{
		"00-19": "🌦",
		"20-29": "⛈",
		"30-39": "🏜",
		"40-49": "🌫️",
		"50-59": "🌦",
		"60-69": "☔",
		"70-79": "❄️",
		"80-99": "⛈",
	}

	// check each range
	for key, emoji := range wmoWeatherMap {
		var min, max int
		_, err := fmt.Sscanf(key, "%d-%d", &min, &max)
		if err != nil {
			continue
		}
		if wmo >= min && wmo <= max {
			return emoji
		}
	}

	// default case if no match is found
	return "❓"
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

func main() {
	if len(os.Args) > 1 {
		place := getPlaceGeocode(os.Args[1])

		if place == (PlaceGeocode{}) {
			fmt.Println("Place not found.")
			os.Exit(1)
		}

		daily := getDailyWeather(place.Latitude, place.Longitude)
		current := getCurrentWeather(place.Latitude, place.Longitude)

		fmt.Printf("%s, %s\n", place.Name, place.Country)
		fmt.Println(strings.Repeat("-", 10))

		fmt.Println("Current:\t", current.Temperature)
		fmt.Println("Feels like:\t", current.FeelsLike)
		fmt.Println("Icon:\t\t", getEmoji(daily.WeatherCode[0]))
		fmt.Println("Max:\t\t", daily.Max[0])
		fmt.Println("Min:\t\t", daily.Min[0])
	} else {
		fmt.Println("gopen-meteo: No parameter passed.\n\nExample usage:\n\tgopen-meteo \"Berlin\"\n\tgopen-meteo \"New York\"")
	}
}
