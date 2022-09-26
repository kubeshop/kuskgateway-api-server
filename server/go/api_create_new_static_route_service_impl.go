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

	kusk "github.com/GIT_USER_ID/GIT_REPO_ID/kusk"
	"github.com/kubeshop/kusk-gateway/pkg/analytics"
)

// CreateNewStaticRouteApiService is a service that implements the logic for the CreateNewStaticRouteApiServicer
// This service should implement the business logic for every endpoint for the CreateNewStaticRouteApi API.
// Include any external packages or services that will be required by this service.
type CreateNewStaticRouteApiService struct {
	kuskClient kusk.Client
}

// NewCreateNewStaticRouteApiService creates a default api service
func NewCreateNewStaticRouteApiService(kuskClient kusk.Client) CreateNewStaticRouteApiServicer {
	return &CreateNewStaticRouteApiService{kuskClient: kuskClient}
}

// CreateStaticRoute - create new static route
func (s *CreateNewStaticRouteApiService) CreateStaticRoute(ctx context.Context, staticRouteItem InlineObject1) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "kusk-api-server", "CreateStaticRoute")
	staticRoute, err := s.kuskClient.CreateStaticRoute(staticRouteItem.Namespace, staticRouteItem.Name, staticRouteItem.EnvoyFleetName, staticRouteItem.EnvoyFleetNamespace, staticRouteItem.Openapi)
	if err != nil {
		return GetResponseFromK8sError(err), err
	}

	toReturn := StaticRouteItem{
		Name:                staticRoute.Name,
		Namespace:           staticRoute.Namespace,
		EnvoyFleetName:      staticRoute.Spec.Fleet.Name,
		EnvoyFleetNamespace: staticRoute.Spec.Fleet.Namespace,
	}
	return Response(http.StatusOK, toReturn), nil
}
