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
)

// ServicesApiService is a service that implements the logic for the ServicesApiServicer
// This service should implement the business logic for every endpoint for the ServicesApi API.
// Include any external packages or services that will be required by this service.
type ServicesApiService struct {
	kuskClient kusk.Client
}

// NewServicesApiService creates a default api service
func NewServicesApiService(kuskClient kusk.Client) ServicesApiServicer {
	return &ServicesApiService{kuskClient: kuskClient}
}

// GetService - Get details for a single service
func (s *ServicesApiService) GetService(ctx context.Context, namespace string, name string) (ImplResponse, error) {
	svc, err := s.kuskClient.GetSvc(namespace, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	ports := []ServicePortItem{}
	for _, port := range svc.Spec.Ports {
		ports = append(ports, ServicePortItem{
			Port:     port.Port,
			Protocol: string(port.Protocol),
		})
	}

	return Response(http.StatusOK, ServiceItem{
		Name:      svc.Name,
		Namespace: svc.Namespace,
		Status:    "available",
		Ports:     ports,
	}), nil
}

// GetServices - Get a list of services
func (s *ServicesApiService) GetServices(ctx context.Context, namespace string) (ImplResponse, error) {
	services, err := s.kuskClient.ListServices(namespace)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}
	toReturn := []ServiceItem{}
	for _, svc := range services.Items {
		anItem := ServiceItem{
			Name:        svc.Name,
			Namespace:   svc.Namespace,
			ServiceType: string(svc.Spec.Type),
		}

		for _, p := range svc.Spec.Ports {
			anItem.Ports = append(anItem.Ports, ServicePortItem{
				Name:       p.Name,
				NodePort:   p.NodePort,
				Port:       p.Port,
				Protocol:   string(p.Protocol),
				TargetPort: p.TargetPort.StrVal,
			})
		}

		toReturn = append(toReturn, anItem)
	}
	return Response(http.StatusOK, toReturn), nil
}
