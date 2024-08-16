package main

import (
	"encoding/json"
	"flag"
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

   creds:
        https://gist.github.com/stellasphere/9490c195ed2b53c707087c8c2db4ec0c
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
	// forced chatgpt to do it because ain't nobody got time for this xd
	wmoWeatherMap := map[int]string{
		0:  "☀️", // Clear sky
		1:  "⛱",  // Few clouds
		2:  "☁️", // Scattered clouds
		3:  "☁️", // Broken clouds
		45: "🌀",  // Tropical storm
		48: "🌀",  // Tropical storm
		51: "🌧",  // Light intensity shower rain
		53: "🌧",  // Shower rain
		55: "🌧",  // Heavy intensity shower rain
		56: "🌨",  // Light rain and snow
		57: "🌨",  // Snow
		61: "🌧",  // Light rain
		63: "🌧",  // Moderate rain
		65: "🌧",  // Heavy intensity rain
		66: "🌨",  // Light intensity drizzle
		67: "🌨",  // Drizzle
		71: "🌨",  // Light snow
		73: "🌨",  // Snow
		75: "🌨",  // Heavy snow
		77: "🌨",  // Sleet
		80: "🌧",  // Light shower rain
		81: "🌧",  // Shower rain
		82: "🌧",  // Heavy shower rain
		85: "🌨",  // Light rain and snow
		86: "🌨",  // Snow showers
		95: "⛈",  // Thunderstorm
		96: "⛈",  // Light thunderstorm
		99: "⛈",  // Heavy thunderstorm
	}

	emoji, exists := wmoWeatherMap[wmo]
	if !exists {
		return "❓"
	}

	return emoji
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

func main() {
	var short bool
	var help bool

	flag.BoolVar(&short, "s", false, "Short one line output")
	flag.BoolVar(&help, "h", false, "Displays help")

	flag.Parse()

	if help {
		fmt.Println("Example usage:\n\tgopen-meteo \"Berlin\"\n\tgopen-meteo \"New York\"")
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("gopen-meteo: No parameter passed.\n\nExample usage:\n\tgopen-meteo \"Berlin\"\n\tgopen-meteo \"New York\"")
		os.Exit(1)
	}

	placeName := args[len(args)-1]

	place := getPlaceGeocode(placeName)

	if place == (PlaceGeocode{}) {
		fmt.Println("Place not found.")
	} else {
		daily := getDailyWeather(place.Latitude, place.Longitude)
		current := getCurrentWeather(place.Latitude, place.Longitude)

		if short {
			// icon current (max/min)
			fmt.Printf("%s %v(%v)°C\n", getEmoji(daily.WeatherCode[0]), current.Temperature, current.FeelsLike)
		} else {
			fmt.Printf("%s %s, %s\n", getEmoji(daily.WeatherCode[0]), place.Name, place.Country)
			fmt.Println(strings.Repeat("-", 10))

			fmt.Printf("Current:\t %v°C\n", current.Temperature)
			fmt.Printf("Feels like:\t %v°C\n", current.FeelsLike)
			fmt.Printf("Max:\t\t %v°C\n", daily.Max[0])
			fmt.Printf("Min:\t\t %v°C\n", daily.Min[0])
		}
	}
}
