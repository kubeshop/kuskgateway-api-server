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

// NamespacesApiService is a service that implements the logic for the NamespacesApiServicer
// This service should implement the business logic for every endpoint for the NamespacesApi API.
// Include any external packages or services that will be required by this service.
type NamespacesApiService struct {
	kuskClient kusk.Client
}

// NewNamespacesApiService creates a default api service
func NewNamespacesApiService(kuskClient kusk.Client) NamespacesApiServicer {
	return &NamespacesApiService{kuskClient: kuskClient}
}

// GetNamespaces - Get a list of namespaces
func (s *NamespacesApiService) GetNamespaces(ctx context.Context) (ImplResponse, error) {
	analytics.SendAnonymousInfo(ctx, s.kuskClient.K8sClient(), "GetNamespaces")
	namespaces, err := s.kuskClient.ListNamespaces()
	if err != nil {
		return GetResponseFromK8sError(err), err
	}
	toReturn := []NamespaceItem{}
	for _, ns := range namespaces.Items {
		toReturn = append(toReturn, NamespaceItem{ns.Name})
	}
	return Response(http.StatusOK, toReturn), nil
}
