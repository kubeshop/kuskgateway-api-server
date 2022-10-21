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

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"github.com/kubeshop/kusk-gateway/pkg/analytics"
	"github.com/kubeshop/kusk-gateway/pkg/spec"
	"gopkg.in/yaml.v3"

	"github.com/kubeshop/kuskgateway-api-server/kusk"
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
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "GetApis")

	var apis *kuskv1.APIList
	var err error

	if fleetname == "" {
		apis, err = s.kuskClient.GetApis(namespace)
		if err != nil {
			return GetResponseFromK8sError(err), err
		}
	} else {
		apis, err = s.kuskClient.GetApiByEnvoyFleet(namespace, fleetnamespace, fleetname)
		if err != nil {
			return GetResponseFromK8sError(err), err
		}
	}
	return Response(http.StatusOK, s.convertAPIListCRDtoAPIsModel(*apis)), nil

}

// GetApi - Get an API instance by namespace and name
func (s *ApisApiService) GetApi(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "GetApi")
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	return Response(http.StatusOK, s.convertAPICRDtoAPIModel(api)), nil
}

// GetApiCRD - Get API CRD from cluster
func (s *ApisApiService) GetApiCRD(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "GetApiCRD")
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	return Response(http.StatusOK, api), nil
}

// GetPostProcessedOpenApiSpec - Get the post-processed OpenAPI spec by API id
func (s *ApisApiService) GetApiDefinition(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "GetApiDefinition")
	api, err := s.kuskClient.GetApi(namespace, name)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	return Response(http.StatusOK, api.Spec.Spec), nil
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

	if opts.Upstream != nil && opts.Upstream.Service != nil {
		apiItem.Service = ApiItemService{
			Name:      opts.Upstream.Service.Name,
			Namespace: opts.Upstream.Service.Namespace,
		}
	}

	apiItem.Fleet = ApiItemFleet{
		Name:      api.Spec.Fleet.Name,
		Namespace: api.Spec.Fleet.Namespace,
	}

	return apiItem
}

// DeployApi - Deploy new API
func (s *ApisApiService) DeployApi(ctx context.Context, payload InlineObject) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "DeployApi")
	api, err := s.kuskClient.CreateApi(payload.Namespace, payload.Name, payload.Openapi, payload.EnvoyFleetName, payload.EnvoyFleetNamespace)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	return Response(http.StatusOK, s.convertAPICRDtoAPIModel(api)), nil
}

func getApiVersion(apiSpec string) string {
	var yml map[string]interface{}
	if err := yaml.Unmarshal([]byte(apiSpec), &yml); err != nil {
		return ""
	}

	for k, v := range yml {
		if k == "info" {
			ss := v.(map[string]interface{})
			return ss["version"].(string)
		}
	}
	return ""
}

// DeleteApi - Delete an API instance by namespace and name
func (s *ApisApiService) DeleteApi(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "DeleteApi")
	if err := s.kuskClient.DeleteAPI(namespace, name); err != nil {
		return GetResponseFromK8sError(err), err

	}
	return Response(http.StatusOK, nil), nil
}

// DeleteApi - Delete an API instance by namespace and name
func (s *ApisApiService) UpdateApi(ctx context.Context, payload InlineObject) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "UpdateApi")

	api, err := s.kuskClient.UpdateApi(payload.Namespace, payload.Name, payload.Openapi, payload.EnvoyFleetName, payload.EnvoyFleetNamespace)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	return Response(http.StatusCreated, api), nil
}
