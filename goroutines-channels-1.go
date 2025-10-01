package main

import (
	"fmt"
	"math/rand"
	"time"
)

//from Grok and I tweaked

// Weather represents the weather data for a city
type Weather struct {
	City        string
	Temperature float64
	Condition   string
}

// fetchWeather simulates an API call to get weather data for a city
func fetchWeather(city string, ch chan<- *Weather) {
	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(4500)+4000) * time.Millisecond)

	// Simulate random weather data
	conditions := []string{"Sunny", "Cloudy", "Rainy", "Snowy"}
	weather := &Weather{
		City:        city,
		Temperature: 15.0 + rand.Float64()*15.0, // Random temp between 15°C and 30°C
		Condition:   conditions[rand.Intn(len(conditions))],
	}

	// Send the result to the channel
	fmt.Printf("Fetched weather for %s: %.1f°C, %s\n", weather.City, weather.Temperature, weather.Condition)
	ch <- weather
}

func main() {

	fmt.Println("=== Weather Fetcher with Go Concurrency ===")

	// List of cities to fetch weather for
	cities := []string{"New York", "London", "Tokyo", "Sydney", "Paris"}

	// Create a channel to receive Weather structs (using pointers for efficiency)
	weatherChan := make(chan *Weather, len(cities))

	// Start a goroutine for each city
	fmt.Println("Starting weather fetches...\n")

	start := time.Now()
	for _, city := range cities {
		go fetchWeather(city, weatherChan) // Launch goroutine
	}
	fmt.Println("\nLaunched all the go routines to fetch weather \n(note this print immediately while goroutines run).\n")

	//better approach is use wait group to wait for all goroutines to finish
	// Collect results from the channel
	weatherResults := make([]*Weather, 0, len(cities))
	for i := 0; i < len(cities); i++ {
		// Receive weather data from channel
		weather := <-weatherChan
		fmt.Println("Received weather for", weather.City)
		weatherResults = append(weatherResults, weather)
	}

	fmt.Println("Completed for loop over cities waitForWeatherResults.")

	// Close the channel (optional here, but good practice)
	close(weatherChan)

	// Print all results
	fmt.Printf("\n=== Weather Report ===\n")
	for _, w := range weatherResults {
		fmt.Printf("%s: %.1f°C, %s\n", w.City, w.Temperature, w.Condition)
	}

	// Show how fast it was
	fmt.Printf("\nFetched weather for %d cities in %v\n", len(cities), time.Since(start))
}
