package k8sclient

import (
	"context"

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error)
	GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error)
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
