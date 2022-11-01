package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	k8sConfig, err := getKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Client Connected")
	for {
		err = ws.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second * 1)
	}
}

func main() {
	http.HandleFunc("/logs", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getKubeConfig() (*rest.Config, error) {
	var err error
	var config *rest.Config
	k8sConfigExists := false
	homeDir, _ := os.UserHomeDir()
	cubeConfigPath := path.Join(homeDir, ".kube/config")

	if _, err := os.Stat(cubeConfigPath); err == nil {
		k8sConfigExists = true
	}

	if cfg, exists := os.LookupEnv("KUBECONFIG"); exists {
		config, err = clientcmd.BuildConfigFromFlags("", cfg)
	} else if k8sConfigExists {
		config, err = clientcmd.BuildConfigFromFlags("", cubeConfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	// default query per second is set to 5
	config.QPS = 40.0
	// default burst is set to 10
	config.Burst = 400.0

	return config, err
}
