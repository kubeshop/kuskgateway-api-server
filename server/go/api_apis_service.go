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
	"errors"
	"net/http"

	kusk "github.com/GIT_USER_ID/GIT_REPO_ID/kusk"
)

// ApisApiService is a service that implements the logic for the ApisApiServicer
// This service should implement the business logic for every endpoint for the ApisApi API.
// Include any external packages or services that will be required by this service.
type ApisApiService struct {
	kuskClient kusk.Client
}

// NewApisApiService creates a default api service
func NewApisApiService(kuskClient kusk.Client) ApisApiServicer {
	return &ApisApiService{kuskClient: kuskClient}
}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (s *ApisApiService) GetPostProcessedOpenApiSpec(ctx context.Context, apiId string, some string) (ImplResponse, error) {
	// TODO - update GetPostProcessedOpenApiSpec with the required logic for this service method.
	// Add api_apis_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	//return Response(200, map[string]interface{}{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetPostProcessedOpenApiSpec method not implemented")
}
