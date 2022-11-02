package main

import (
	"log"
	"net/http"
)

func main() {
	clientSet, err := NewClientSet()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	websocketHandler := websocketHandler{
		coreV1: clientSet.CoreV1(),
	}
	mux.Handle("/logs", websocketHandler)
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
