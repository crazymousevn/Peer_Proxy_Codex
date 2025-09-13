package main

import (
	"log"
	"net/http"

	"peerproxy/internal/signaling"
)

func main() {
	srv := signaling.NewServer()
	log.Println("Starting signaling server on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
