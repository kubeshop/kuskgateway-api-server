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
	"strings"

	kusk "github.com/GIT_USER_ID/GIT_REPO_ID/kusk"
	"github.com/GIT_USER_ID/GIT_REPO_ID/util"
	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"github.com/kubeshop/kusk-gateway/pkg/spec"
	"gopkg.in/yaml.v3"
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
func (s *ApisApiService) GetApis(ctx context.Context, namespace string, fleetname string, fleetnamespace string) (ImplResponse, error) {
	apis := &kuskv1.APIList{}
	var err error
	if fleetname == "" {
		apis, err = s.kuskClient.GetApis(namespace)
		if err != nil {
			return Response(http.StatusInternalServerError, err), err
		}

	} else {
		apis, err = s.kuskClient.GetApiByEnvoyFleet(namespace, fleetnamespace, fleetname)
		if err != nil {
			return Response(http.StatusInternalServerError, err), err
		}
	}
	return Response(http.StatusOK, s.convertAPIListCRDtoAPIsModel(*apis)), nil

}

// GetApi - Get an API instance by namespace and name
func (s *ApisApiService) GetApi(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, s.convertAPICRDtoAPIModel(api)), nil
}

// GetApiCRD - Get API CRD from cluster
func (s *ApisApiService) GetApiCRD(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, api), nil
}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (s *ApisApiService) GetApiDefinition(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, api.Spec.Spec), nil
}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (s *ApisApiService) GetPostProcessedOpenApiSpec(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}
	rawyaml := util.ParseKuskOpenAPI(api.Spec.Spec)
	yml, _ := yaml.Marshal(rawyaml)

	return Response(http.StatusOK, string(yml)), nil
}

func (s *ApisApiService) convertAPIListCRDtoAPIsModel(apis kuskv1.APIList) []ApiItem {
	toReturn := []ApiItem{}
	for _, api := range apis.Items {
		toReturn = append(toReturn, s.convertAPICRDtoAPIModel(&api))
	}
	return toReturn
}

func (s *ApisApiService) convertAPICRDtoAPIModel(api *kuskv1.API) ApiItem {
	parser := spec.NewParser(nil)
	apiItem := ApiItem{
		Name:      api.Name,
		Namespace: api.Namespace,
	}

	apiSpec, err := parser.ParseFromReader(strings.NewReader(api.Spec.Spec))
	if err != nil {
		return apiItem
	}

	opts, err := spec.GetOptions(apiSpec)
	if err != nil {
		return apiItem
	}

	apiItem.Version = getApiVersion(api.Spec.Spec)

	apiItem.Service = ApiItemService{
		Name:      opts.Upstream.Service.Name,
		Namespace: opts.Upstream.Service.Namespace,
	}
	apiItem.Fleet = ApiItemFleet{
		Name:      api.Spec.Fleet.Name,
		Namespace: api.Spec.Fleet.Namespace,
	}

	return apiItem
}

func getApiVersion(apiSpec string) string {
	var yml map[string]interface{}
	yaml.Unmarshal([]byte(apiSpec), &yml)

	for k, v := range yml {
		if k == "info" {
			ss := v.(map[string]interface{})
			return ss["version"].(string)
		}
	}
	return ""
}
