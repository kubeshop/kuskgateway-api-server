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
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ApisApiController binds http requests to an api service and writes the service results to the http response
type ApisApiController struct {
	service      ApisApiServicer
	errorHandler ErrorHandler
}

// ApisApiOption for how the controller is set up.
type ApisApiOption func(*ApisApiController)

// WithApisApiErrorHandler inject ErrorHandler into controller
func WithApisApiErrorHandler(h ErrorHandler) ApisApiOption {
	return func(c *ApisApiController) {
		c.errorHandler = h
	}
}

// NewApisApiController creates a default api controller
func NewApisApiController(s ApisApiServicer, opts ...ApisApiOption) Router {
	controller := &ApisApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the ApisApiController
func (c *ApisApiController) Routes() Routes {
	return Routes{
		{
			"GetApi",
			strings.ToUpper("Get"),
			"/apis/{namespace}/{name}",
			c.GetApi,
		},
		{
			"GetApis",
			strings.ToUpper("Get"),
			"/apis",
			c.GetApis,
		},
		{
			"GetPostProcessedOpenApiSpec",
			strings.ToUpper("Get"),
			"/apis/{namespace}/{name}/postProcessedOpenApiSpec",
			c.GetPostProcessedOpenApiSpec,
		},
		{
			"GetRawOpenApiSpec",
			strings.ToUpper("Get"),
			"/apis/{namespace}/{name}/rawOpenApiSpec",
			c.GetRawOpenApiSpec,
		},
	}
}

// GetApi - Get an API instance by namespace and name
func (c *ApisApiController) GetApi(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespaceParam := params["namespace"]

	nameParam := params["name"]

	result, err := c.service.GetApi(r.Context(), namespaceParam, nameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetApis - Get a list of APIs
func (c *ApisApiController) GetApis(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fleetParam := query.Get("fleet")
	result, err := c.service.GetApis(r.Context(), fleetParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (c *ApisApiController) GetPostProcessedOpenApiSpec(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespaceParam := params["namespace"]

	nameParam := params["name"]

	result, err := c.service.GetPostProcessedOpenApiSpec(r.Context(), namespaceParam, nameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetRawOpenApiSpec - Get the raw OpenAPI spec by API id
func (c *ApisApiController) GetRawOpenApiSpec(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespaceParam := params["namespace"]

	nameParam := params["name"]

	result, err := c.service.GetRawOpenApiSpec(r.Context(), namespaceParam, nameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
