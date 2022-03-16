package k8sclient

import (
	"context"

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error)
	GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error)

	GetApis() (kuskv1.APIList, error)
	GetApi(namespace, name string) (*kuskv1.API, error)
	GetApiByEnvoyFleet(fleetNamespace, fleetName string) ([]kuskv1.API, error)
}

type kuskClient struct {
	client client.Client
}

func NewClient(c client.Client) Client {
	return &kuskClient{
		client: c,
	}
}

func (k *kuskClient) GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error) {

	list := &kuskv1.EnvoyFleetList{}

	if err := k.client.List(context.TODO(), list, &client.ListOptions{}); err != nil {
		return nil, err
	}
	return list, nil
}

func (k *kuskClient) GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error) {

	envoy := &kuskv1.EnvoyFleet{}

	if err := k.client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, envoy); err != nil {
		return nil, err
	}
	return envoy, nil
}

func (k *kuskClient) GetApis() (kuskv1.APIList, error) {

	list := kuskv1.APIList{}

	if err := k.client.List(context.TODO(), &list, &client.ListOptions{}); err != nil {
		return list, err
	}
	return list, nil
}

func (k *kuskClient) GetApi(namespace, name string) (*kuskv1.API, error) {

	api := &kuskv1.API{}

	if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, api); err != nil {
		return nil, err
	}
	return api, nil
}

// GetApiByFleet gets all APIs associated with the EnvoyFleet
func (k *kuskClient) GetApiByEnvoyFleet(fleetNamespace, fleetName string) ([]kuskv1.API, error) {
	list := kuskv1.APIList{}
	if err := k.client.List(context.TODO(), &list, &client.ListOptions{}); err != nil {
		return nil, err
	}

	toReturn := []kuskv1.API{}
	for _, api := range list.Items {
		if api.Spec.Fleet.Name == fleetName && api.Spec.Fleet.Namespace == fleetNamespace {
			toReturn = append(toReturn, api)
		}
	}
	return toReturn, nil
}
