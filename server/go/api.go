/*
 * Kusk Gateway API
 *
 * This is the Kusk Gateway Management API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
)

// ApisApiRouter defines the required methods for binding the api requests to a responses for the ApisApi
// The ApisApiRouter implementation should parse necessary information from the http request,
// pass the data to a ApisApiServicer to perform the required actions, then write the service results to the http response.
type ApisApiRouter interface {
	GetApi(http.ResponseWriter, *http.Request)
	GetApis(http.ResponseWriter, *http.Request)
	GetPostProcessedOpenApiSpec(http.ResponseWriter, *http.Request)
	GetRawOpenApiSpec(http.ResponseWriter, *http.Request)
}

// FleetsApiRouter defines the required methods for binding the api requests to a responses for the FleetsApi
// The FleetsApiRouter implementation should parse necessary information from the http request,
// pass the data to a FleetsApiServicer to perform the required actions, then write the service results to the http response.
type FleetsApiRouter interface {
	GetEnvoyFleet(http.ResponseWriter, *http.Request)
	GetEnvoyFleets(http.ResponseWriter, *http.Request)
}

// ServicesApiRouter defines the required methods for binding the api requests to a responses for the ServicesApi
// The ServicesApiRouter implementation should parse necessary information from the http request,
// pass the data to a ServicesApiServicer to perform the required actions, then write the service results to the http response.
type ServicesApiRouter interface {
	GetService(http.ResponseWriter, *http.Request)
	GetServices(http.ResponseWriter, *http.Request)
}

// ApisApiServicer defines the api actions for the ApisApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ApisApiServicer interface {
	GetApi(context.Context, string, string) (ImplResponse, error)
	GetApis(context.Context, string) (ImplResponse, error)
	GetPostProcessedOpenApiSpec(context.Context, string, string) (ImplResponse, error)
	GetRawOpenApiSpec(context.Context, string, string) (ImplResponse, error)
}

// FleetsApiServicer defines the api actions for the FleetsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type FleetsApiServicer interface {
	GetEnvoyFleet(context.Context, string, string) (ImplResponse, error)
	GetEnvoyFleets(context.Context, string) (ImplResponse, error)
}

// ServicesApiServicer defines the api actions for the ServicesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ServicesApiServicer interface {
	GetService(context.Context, string, string) (ImplResponse, error)
	GetServices(context.Context, string) (ImplResponse, error)
}
