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
	"github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"github.com/kubeshop/kusk-gateway/pkg/analytics"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StaticRouteApiService is a service that implements the logic for the StaticRouteApiServicer
// This service should implement the business logic for every endpoint for the StaticRouteApi API.
// Include any external packages or services that will be required by this service.
type StaticRouteApiService struct {
	kuskClient kusk.Client
}

// NewStaticRouteApiService creates a default api service
func NewStaticRouteApiService(kuskClient kusk.Client) StaticRouteApiServicer {
	return &StaticRouteApiService{kuskClient: kuskClient}
}

// DeleteStaticRoute - Delete a StaticRoute by namespace and name
func (s *StaticRouteApiService) DeleteStaticRoute(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "DeleteStaticRoute")
	if err := s.kuskClient.DeleteStaticRoute(v1alpha1.StaticRoute{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}); err != nil {
		return GetResponseFromK8sError(err), err
	}
	return Response(http.StatusOK, nil), nil
}
