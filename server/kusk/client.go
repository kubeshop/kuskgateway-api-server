package k8sclient

import (
	"context"
	"fmt"

	"github.com/kubeshop/kusk-gateway/api/v1alpha1"
	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error)
	GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error)
	CreateFleet(kuskv1.EnvoyFleet) (*kuskv1.EnvoyFleet, error)
	DeleteFleet(kuskv1.EnvoyFleet) error

	GetApis(namespace string) (*kuskv1.APIList, error)
	GetApi(namespace, name string) (*kuskv1.API, error)
	GetApiByEnvoyFleet(namespace, fleetNamespace, fleetName string) (*kuskv1.APIList, error)
	CreateApi(namespace, name, openapispec, fleetName, fleetnamespace string) (*kuskv1.API, error)
	DeleteAPI(namespace, name string) error

	GetStaticRoute(namespace, name string) (*kuskv1.StaticRoute, error)
	GetStaticRoutes(namespace string) (*kuskv1.StaticRouteList, error)
	CreateStaticRoute(namespace, name, fleetName, fleetNamespace, specs string) (*kuskv1.StaticRoute, error)
	DeleteStaticRoute(kuskv1.StaticRoute) error

	GetSvc(namespace, name string) (*corev1.Service, error)
	ListServices(namespace string) (*corev1.ServiceList, error)
	ListNamespaces() (*corev1.NamespaceList, error)

	K8sClient() client.Client
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

func (k *kuskClient) CreateFleet(fleet kuskv1.EnvoyFleet) (*kuskv1.EnvoyFleet, error) {
	if err := k.client.Create(context.TODO(), &fleet, &client.CreateOptions{}); err != nil {
		return nil, err
	}

	return &fleet, nil
}

func (k *kuskClient) DeleteFleet(fleet kuskv1.EnvoyFleet) error {
	return k.client.Delete(context.TODO(), &fleet, &client.DeleteOptions{})
}

func (k *kuskClient) GetApis(namespace string) (*kuskv1.APIList, error) {
	list := &kuskv1.APIList{}
	if err := k.client.List(context.TODO(), list, &client.ListOptions{Namespace: namespace}); err != nil {
		return nil, err
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
func (k *kuskClient) GetApiByEnvoyFleet(namespace, fleetNamespace, fleetName string) (*kuskv1.APIList, error) {
	list := kuskv1.APIList{}
	if err := k.client.List(context.TODO(), &list, &client.ListOptions{Namespace: namespace}); err != nil {
		return nil, err
	}

	toReturn := []kuskv1.API{}
	for _, api := range list.Items {
		if api.Spec.Fleet.Name == fleetName && api.Spec.Fleet.Namespace == fleetNamespace {
			toReturn = append(toReturn, api)
		}
	}

	return &kuskv1.APIList{Items: toReturn}, nil
}

func (k *kuskClient) CreateApi(namespace, name, openapispec string, fleetName string, fleetnamespace string) (*kuskv1.API, error) {
	api := &kuskv1.API{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: kuskv1.APISpec{
			Spec: openapispec,
			Fleet: &kuskv1.EnvoyFleetID{
				Name:      fleetName,
				Namespace: fleetnamespace,
			},
		},
	}
	if err := k.client.Create(context.TODO(), api, &client.CreateOptions{}); err != nil {
		return nil, err
	}
	return api, nil
}

func (k *kuskClient) DeleteAPI(namespace, name string) error {
	api := &kuskv1.API{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return k.client.Delete(context.TODO(), api, &client.DeleteOptions{})
}

func (k *kuskClient) GetSvc(namespace, name string) (*corev1.Service, error) {
	svc := &corev1.Service{}
	if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, svc); err != nil {
		return nil, err
	}

	return svc, nil
}

func (k *kuskClient) ListServices(namespace string) (*corev1.ServiceList, error) {
	list := &corev1.ServiceList{}

	if err := k.client.List(context.TODO(), list, &client.ListOptions{Namespace: namespace}); err != nil {
		return nil, err
	}

	return list, nil

}

func (k *kuskClient) GetStaticRoute(namespace, name string) (*kuskv1.StaticRoute, error) {
	staticRoute := &kuskv1.StaticRoute{}
	if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, staticRoute); err != nil {
		return nil, err
	}

	return staticRoute, nil
}

func (k *kuskClient) GetStaticRoutes(namespace string) (*kuskv1.StaticRouteList, error) {
	list := &kuskv1.StaticRouteList{}
	if err := k.client.List(context.TODO(), list, &client.ListOptions{}); err != nil {
		return nil, err
	}

	return list, nil
}

func (k *kuskClient) CreateStaticRoute(namespace, name, fleetName, fleetNamespace, specs string) (*kuskv1.StaticRoute, error) {
	staticRoute := &kuskv1.StaticRoute{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: kuskv1.StaticRouteSpec{
			Fleet: &kuskv1.EnvoyFleetID{
				Name:      fleetName,
				Namespace: fleetNamespace,
			},
		},
	}

	tmp := &kuskv1.StaticRoute{}

	err := yaml.Unmarshal([]byte(specs), tmp)
	if err != nil {
		fmt.Println(err)
	}

	staticRoute.Spec.Paths = tmp.Spec.Paths
	staticRoute.Spec.Hosts = tmp.Spec.Hosts

	if err := k.client.Create(context.TODO(), staticRoute, &client.CreateOptions{}); err != nil {
		return nil, err
	}

	return staticRoute, nil
}

func (k *kuskClient) DeleteStaticRoute(sroute v1alpha1.StaticRoute) error {
	return k.client.Delete(context.TODO(), &sroute, &client.DeleteOptions{})
}

func (k *kuskClient) ListNamespaces() (*corev1.NamespaceList, error) {
	list := &corev1.NamespaceList{}
	if err := k.client.List(context.TODO(), list, &client.ListOptions{}); err != nil {
		return nil, err
	}
	return list, nil
}

func (k *kuskClient) K8sClient() client.Client {
	return k.client
}
