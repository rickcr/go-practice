package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
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
func fetchWeather(city string, ch chan<- *Weather, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when this goroutine completes
	fetchstart := time.Now()

	now := uint64(time.Now().UnixNano())
	source := rand.NewPCG(now, now>>32+uint64(len(city)))
	r := rand.New(source)

	sleepDuration := time.Duration(r.IntN(7000)+500) * time.Millisecond
 
	// Simulate network delay
	time.Sleep(sleepDuration)

	// Simulate random weather data
	conditions := []string{"Sunny", "Cloudy", "Rainy", "Snowy"}
	weather := &Weather{
		City:        city,
		Temperature: 15.0 + rand.Float64()*15.0, // Random temp between 15°C and 30°C
		Condition:   conditions[rand.IntN(len(conditions))],
	}
	 
	actualDuration := time.Since(fetchstart)

	// Send the result to the channel
	fmt.Printf("Fetched weather for %s: %.1f°C, %s, (time %s)\n", weather.City, weather.Temperature, weather.Condition, actualDuration)
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

	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)                               // Increment the counter for each goroutine
		go fetchWeather(city, weatherChan, &wg) // Launch goroutine
	}
	fmt.Println("\nLaunched all the go routines to fetch weather. Now waiting for all to complete. \n(note this print immediately while goroutines run).\n")

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("\nAll goroutines completed!\n")

	// Close the channel since we're done sending
	close(weatherChan)

	// Collect results from the channel
	weatherResults := make([]*Weather, 0, len(cities))
	for weather := range weatherChan {
		fmt.Println("Received weather for", weather.City)
		weatherResults = append(weatherResults, weather)
	}

	fmt.Println("Completed collecting all weather results.")

	// Print all results
	fmt.Printf("\n=== Weather Report ===\n")
	for _, w := range weatherResults {
		fmt.Printf("%s: %.1f°C, %s\n", w.City, w.Temperature, w.Condition)
	}

	// Show how fast it was
	fmt.Printf("\nFetched weather for %d cities in %v\n", len(cities), time.Since(start))
}
