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
	"log"
	"net/http"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
)

func main() {
	log.Printf("Server started")

	ApisApiService := openapi.NewApisApiService()
	ApisApiController := openapi.NewApisApiController(ApisApiService)

	FleetsApiService := openapi.NewFleetsApiService()
	FleetsApiController := openapi.NewFleetsApiController(FleetsApiService)

	ServicesApiService := openapi.NewServicesApiService()
	ServicesApiController := openapi.NewServicesApiController(ServicesApiService)

	router := openapi.NewRouter(ApisApiController, FleetsApiController, ServicesApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
