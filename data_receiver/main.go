package main

import (
	"fmt"
	"log"
	"net/http"
	"tolling-micorservices/types"

	"github.com/gorilla/websocket"
)

func main() {

	// websocket
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWS)
	// fmt.Println("Data receiver working file..")

	http.ListenAndServe(":30000", nil)

}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}
func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("NEW OBU client connected!")
	for {
		var data types.OBUData

		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Printf("error: %v", err)
			continue
		}
		fmt.Printf("Received OBU Data from [%d] :: <lat %.2f,long %.2f>\n", data.OBUID, data.Lat,data.Long)
		dr.msgch <- data
	}
}
