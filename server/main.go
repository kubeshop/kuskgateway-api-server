/*
 * Kusk Gateway API
 *
 * This is the Kusk Gateway Management API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/handlers"
	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"github.com/kubeshop/kusk-gateway/pkg/analytics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	openapi "github.com/kubeshop/kuskgateway-api-server/go"
	"github.com/kubeshop/kuskgateway-api-server/kusk"
)

func main() {
	k8sClient, err := getClient()
	if err != nil {
		log.Fatalf(fmt.Errorf("unable to get kubernetes client: %w", err).Error())
	}

	kuskClient := kusk.NewClient(k8sClient)

	ApisApiService := openapi.NewApisApiService(kuskClient)
	ApisApiController := openapi.NewApisApiController(ApisApiService)

	FleetsApiService := openapi.NewFleetsApiService(kuskClient)
	FleetsApiController := openapi.NewFleetsApiController(FleetsApiService)

	CreateNewFleetApiService := openapi.NewCreateNewFleetApiService(kuskClient)
	CreateNewFlettApiController := openapi.NewCreateNewFleetApiController(CreateNewFleetApiService)

	CreateFleetService := openapi.NewCreateNewFleetApiService(kuskClient)
	CreateFleetController := openapi.NewCreateNewFleetApiController(CreateFleetService)

	ServicesApiService := openapi.NewServicesApiService(kuskClient)
	ServicesApiController := openapi.NewServicesApiController(ServicesApiService)

	StaticCreateRouteApiService := openapi.NewCreateNewStaticRouteApiService(kuskClient)
	StaticCreateRouteApiController := openapi.NewCreateNewStaticRouteApiController(StaticCreateRouteApiService)

	StaticRouteApiService := openapi.NewStaticRouteApiService(kuskClient)
	StaticRouteApiController := openapi.NewStaticRouteApiController(StaticRouteApiService)

	StaticRoutesApiService := openapi.NewStaticRoutesApiService(kuskClient)
	StaticRoutesApiController := openapi.NewStaticRoutesApiController(StaticRoutesApiService)

	ProbeController := openapi.NewProbeController()

	NamespacesApiService := openapi.NewNamespacesApiService(kuskClient)
	NamespaceApiController := openapi.NewNamespacesApiController(NamespacesApiService)

	router := openapi.NewRouter(
		ApisApiController,
		FleetsApiController,
		CreateNewFlettApiController,
		CreateFleetController,
		ServicesApiController,
		StaticCreateRouteApiController,
		StaticRouteApiController,
		StaticRoutesApiController,
		ProbeController,
		NamespaceApiController,
	)
	analytics.SendAnonymousInfo(context.Background(), kuskClient.K8sClient(), "kusk-api-server", "Starting kusk API server")

	log.Printf("Server started :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, methodsOk, originsOk)(router)))
}

var (
	headersOk = handlers.AllowedHeaders([]string{
		"Accept",
		"Content-Language",
		"Origin",
		"Content-Type",
		"Accept-Language",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"X-CSRF-Token",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
		"Access-Control-Allow-Origin",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials"})

	methodsOk = handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST", "PUT", "DELETE"})
	originsOk = handlers.AllowedOrigins([]string{"*"})
)

func getConfig() (*rest.Config, error) {
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
func getClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	kuskv1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)

	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return client.New(config, client.Options{Scheme: scheme})
}
