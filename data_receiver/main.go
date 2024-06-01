package main

import (
	"fmt"
	"log"
	"net/http"
	"tolling-micorservices/types"

	"github.com/gorilla/websocket"
)

// Entry point of the program
func main() {
	// Create a new DataReceiver instance
	recv := NewDataReceiver()

	// Set up the HTTP handler for the WebSocket endpoint
	http.HandleFunc("/ws", recv.handleWS)

	// Start the HTTP server on port 30000
	http.ListenAndServe(":30000", nil)
}

// DataReceiver struct holds a channel for messages and a WebSocket connection
type DataReceiver struct {
	msgch chan types.OBUData // Channel to hold incoming OBUData
	conn  *websocket.Conn    // WebSocket connection
}

// Constructor for DataReceiver
func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128), // Initialize channel with buffer size 128
	}
}

// Handle the WebSocket connection upgrade request
func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrader to handle the WebSocket upgrade request
	u := websocket.Upgrader{
		ReadBufferSize:  1028, // Buffer size for reading
		WriteBufferSize: 1028, // Buffer size for writing
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err) // Log and terminate if upgrade fails
	}

	// Assign the WebSocket connection to the DataReceiver instance
	dr.conn = conn

	// Start the loop to receive messages from the WebSocket connection
	go dr.wsReceiveLoop()
}

// Loop to receive messages from the WebSocket connection
func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("NEW OBU client connected!")
	for {
		var data types.OBUData // Variable to hold incoming OBUData

		// Read JSON data from the WebSocket connection and unmarshal it into data
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Printf("error: %v", err) // Log any errors
			continue                     // Continue the loop on error
		}

		// Print the received data
		fmt.Printf("Received OBU Data from [%d] :: <lat %.2f,long %.2f>\n", data.OBUID, data.Lat, data.Long)

		// Send the received data to the msgch channel
		dr.msgch <- data
	}
}
