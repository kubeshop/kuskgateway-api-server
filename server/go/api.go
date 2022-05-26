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
	DeleteApi(http.ResponseWriter, *http.Request)
	DeployApi(http.ResponseWriter, *http.Request)
	GetApi(http.ResponseWriter, *http.Request)
	GetApiCRD(http.ResponseWriter, *http.Request)
	GetApiDefinition(http.ResponseWriter, *http.Request)
	GetApis(http.ResponseWriter, *http.Request)
}

// CreateNewFleetApiRouter defines the required methods for binding the api requests to a responses for the CreateNewFleetApi
// The CreateNewFleetApiRouter implementation should parse necessary information from the http request,
// pass the data to a CreateNewFleetApiServicer to perform the required actions, then write the service results to the http response.
type CreateNewFleetApiRouter interface {
	CreateFleet(http.ResponseWriter, *http.Request)
}

// CreateNewStaticRouteApiRouter defines the required methods for binding the api requests to a responses for the CreateNewStaticRouteApi
// The CreateNewStaticRouteApiRouter implementation should parse necessary information from the http request,
// pass the data to a CreateNewStaticRouteApiServicer to perform the required actions, then write the service results to the http response.
type CreateNewStaticRouteApiRouter interface {
	CreateStaticRoute(http.ResponseWriter, *http.Request)
}

// FleetsApiRouter defines the required methods for binding the api requests to a responses for the FleetsApi
// The FleetsApiRouter implementation should parse necessary information from the http request,
// pass the data to a FleetsApiServicer to perform the required actions, then write the service results to the http response.
type FleetsApiRouter interface {
	DeleteFleet(http.ResponseWriter, *http.Request)
	GetEnvoyFleet(http.ResponseWriter, *http.Request)
	GetEnvoyFleetCRD(http.ResponseWriter, *http.Request)
	GetEnvoyFleets(http.ResponseWriter, *http.Request)
}

// NamespacesApiRouter defines the required methods for binding the api requests to a responses for the NamespacesApi
// The NamespacesApiRouter implementation should parse necessary information from the http request,
// pass the data to a NamespacesApiServicer to perform the required actions, then write the service results to the http response.
type NamespacesApiRouter interface {
	GetNamespaces(http.ResponseWriter, *http.Request)
}

// ServicesApiRouter defines the required methods for binding the api requests to a responses for the ServicesApi
// The ServicesApiRouter implementation should parse necessary information from the http request,
// pass the data to a ServicesApiServicer to perform the required actions, then write the service results to the http response.
type ServicesApiRouter interface {
	GetService(http.ResponseWriter, *http.Request)
	GetServices(http.ResponseWriter, *http.Request)
}

// StaticRouteApiRouter defines the required methods for binding the api requests to a responses for the StaticRouteApi
// The StaticRouteApiRouter implementation should parse necessary information from the http request,
// pass the data to a StaticRouteApiServicer to perform the required actions, then write the service results to the http response.
type StaticRouteApiRouter interface {
	DeleteStaticRoute(http.ResponseWriter, *http.Request)
}

// StaticRoutesApiRouter defines the required methods for binding the api requests to a responses for the StaticRoutesApi
// The StaticRoutesApiRouter implementation should parse necessary information from the http request,
// pass the data to a StaticRoutesApiServicer to perform the required actions, then write the service results to the http response.
type StaticRoutesApiRouter interface {
	GetStaticRoute(http.ResponseWriter, *http.Request)
	GetStaticRouteCRD(http.ResponseWriter, *http.Request)
	GetStaticRoutes(http.ResponseWriter, *http.Request)
}

// ApisApiServicer defines the api actions for the ApisApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ApisApiServicer interface {
	DeleteApi(context.Context, string, string) (ImplResponse, error)
	DeployApi(context.Context, InlineObject) (ImplResponse, error)
	GetApi(context.Context, string, string) (ImplResponse, error)
	GetApiCRD(context.Context, string, string) (ImplResponse, error)
	GetApiDefinition(context.Context, string, string) (ImplResponse, error)
	GetApis(context.Context, string, string, string) (ImplResponse, error)
}

// CreateNewFleetApiServicer defines the api actions for the CreateNewFleetApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CreateNewFleetApiServicer interface {
	CreateFleet(context.Context, ServiceItem) (ImplResponse, error)
}

// CreateNewStaticRouteApiServicer defines the api actions for the CreateNewStaticRouteApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CreateNewStaticRouteApiServicer interface {
	CreateStaticRoute(context.Context, InlineObject1) (ImplResponse, error)
}

// FleetsApiServicer defines the api actions for the FleetsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type FleetsApiServicer interface {
	DeleteFleet(context.Context, string, string) (ImplResponse, error)
	GetEnvoyFleet(context.Context, string, string) (ImplResponse, error)
	GetEnvoyFleetCRD(context.Context, string, string) (ImplResponse, error)
	GetEnvoyFleets(context.Context, string) (ImplResponse, error)
}

// NamespacesApiServicer defines the api actions for the NamespacesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type NamespacesApiServicer interface {
	GetNamespaces(context.Context) (ImplResponse, error)
}

// ServicesApiServicer defines the api actions for the ServicesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ServicesApiServicer interface {
	GetService(context.Context, string, string) (ImplResponse, error)
	GetServices(context.Context, string) (ImplResponse, error)
}

// StaticRouteApiServicer defines the api actions for the StaticRouteApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StaticRouteApiServicer interface {
	DeleteStaticRoute(context.Context, string, string) (ImplResponse, error)
}

// StaticRoutesApiServicer defines the api actions for the StaticRoutesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StaticRoutesApiServicer interface {
	GetStaticRoute(context.Context, string, string) (ImplResponse, error)
	GetStaticRouteCRD(context.Context, string, string) (ImplResponse, error)
	GetStaticRoutes(context.Context, string) (ImplResponse, error)
}
