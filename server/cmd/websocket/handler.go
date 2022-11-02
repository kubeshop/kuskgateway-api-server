package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type websocketHandler struct {
	coreV1 typedCoreV1.CoreV1Interface
}

const (
	defaultNamespaceParam = "kusk-system"
	defaultNameParam      = "kusk-gateway-envoy-fleet"
)

func (h websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	queryParams := r.URL.Query()
	namespace := queryParams.Get("namespace")
	if namespace == "" {
		namespace = defaultNamespaceParam
	}
	name := queryParams.Get("name")
	if name == "" {
		name = defaultNameParam
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer ws.Close()

	log.Println("client connected")
	stream, err := GetServiceContainerLogStream(
		r.Context(),
		namespace,
		name,
		"envoy",
		h.coreV1,
	)

	defer stream.Close()

	c := client{
		conn:      ws,
		logStream: stream,
	}

	stopCh := make(chan struct{})
	go c.readPump(r.Context(), stopCh)
	go c.writePump(r.Context(), stopCh)

	<-stopCh
	fmt.Println("request finished")
}
