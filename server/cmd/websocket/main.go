package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()
}

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
	log.Println("starting server on :", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
