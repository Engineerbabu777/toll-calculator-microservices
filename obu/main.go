package main

import (
	"fmt"                      // Package for formatted I/O
	"log"                      // Package for logging
	"math"                     // Package for mathematical constants and functions
	"math/rand"                // Package for generating pseudo-random numbers
	"time"                     // Package for time-related functions
	"tolling-micorservices/types" // Custom package for types

	"github.com/gorilla/websocket" // Package for WebSocket connections
)

// Interval between sending data
var sendInterval = time.Second

// WebSocket endpoint URL
const wsEndpoint = "ws://127.0.0.1:30000/ws"

// Generate a random coordinate
func genCoord() float64 {
	n := float64(rand.Intn(100) + 1) // Generate a random integer between 1 and 100
	f := rand.Float64()              // Generate a random float between 0 and 1
	return n + f                     // Combine the integer and float to create a coordinate
}

// Generate random latitude and longitude
func genLatLong() (float64, float64) {
	return genCoord(), genCoord() // Generate and return random latitude and longitude
}

// Main function
func main() {
	// Generate a list of random OBUIDs
	obuIDS := generateOBUIDS(20)

	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err) // Log and terminate if connection fails
	}

	// Infinite loop to continuously send data
	for {
		// Iterate over each OBUID
		for i := 0; i < len(obuIDS); i++ {
			// Generate random latitude and longitude
			lat, long := genLatLong()

			// Create OBUData with generated coordinates and OBUID
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}

			// Send OBUData as JSON over WebSocket connection
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err) // Log and terminate if sending fails
			}

			// Print the sent data
			fmt.Printf("%+v\n", data)
		}

		// Wait for the specified interval before sending the next batch
		time.Sleep(sendInterval)
	}
}

// Generate a list of random OBUIDs
func generateOBUIDS(n int) []int {
	ids := make([]int, n) // Create a slice to hold the OBUIDs

	// Populate the slice with random integers
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt) // Generate a random integer
	}

	return ids // Return the slice of OBUIDs
}
