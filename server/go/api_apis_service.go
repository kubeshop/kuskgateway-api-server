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
	"fmt"
	"net/http"

	kusk "github.com/GIT_USER_ID/GIT_REPO_ID/kusk"
	"github.com/kubeshop/kusk-gateway/api/v1alpha1"
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

// GetApis - Get a list of APIs
func (s *ApisApiService) GetApis(ctx context.Context, fleet string) (ImplResponse, error) {
	fmt.Println("s", s)
	apis, err := s.kuskClient.GetApis()
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, convertAPIListCRDtoAPIsModel(apis)), nil
}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (s *ApisApiService) GetPostProcessedOpenApiSpec(ctx context.Context, apiId string) (ImplResponse, error) {
	// TODO - update GetPostProcessedOpenApiSpec with the required logic for this service method.
	// Add api_apis_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	//return Response(200, map[string]interface{}{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetPostProcessedOpenApiSpec method not implemented")
}

// GetRawOpenApiSpec - Get the raw OpenAPI spec by API id
func (s *ApisApiService) GetRawOpenApiSpec(ctx context.Context, apiId string) (ImplResponse, error) {

	// TODO - update GetRawOpenApiSpec with the required logic for this service method.
	// Add api_apis_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	//return Response(200, map[string]interface{}{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetRawOpenApiSpec method not implemented")
}

func convertAPIListCRDtoAPIsModel(apis v1alpha1.APIList) []ApiItem {
	toReturn := make([]ApiItem, len(apis.Items))
	for _, api := range apis.Items {
		toReturn = append(toReturn, convertAPICRDtoAPIModel(&api))
	}
	return toReturn
}

func convertAPICRDtoAPIModel(api *v1alpha1.API) ApiItem {
	return ApiItem{
		Name: api.Name,
	}
}
