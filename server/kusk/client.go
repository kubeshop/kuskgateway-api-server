package kusk

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrNotFound = errors.New("error not found")

type Client interface {
	GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error)
	GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error)
	CreateFleet(kuskv1.EnvoyFleet) (*kuskv1.EnvoyFleet, error)
	DeleteFleet(kuskv1.EnvoyFleet) error

	GetApis(namespace string) (*kuskv1.APIList, error)
	GetApi(namespace, name string) (*kuskv1.API, error)
	GetApiByEnvoyFleet(namespace, fleetNamespace, fleetName string) (*kuskv1.APIList, error)
	CreateApi(namespace, name, openapispec, fleetName, fleetnamespace string) (*kuskv1.API, error)
	UpdateApi(namespace, name, openapispec, fleetName, fleetnamespace string) (*kuskv1.API, error)
	DeleteAPI(namespace, name string) error

	GetStaticRoute(namespace, name string) (*kuskv1.StaticRoute, error)
	GetStaticRoutes(namespace string) (*kuskv1.StaticRouteList, error)
	CreateStaticRoute(namespace, name, fleetName, fleetNamespace, specs string) (*kuskv1.StaticRoute, error)
	UpdateStaticRoute(namespace, name, fleetName, fleetNamespace, specs string) (*kuskv1.StaticRoute, error)
	DeleteStaticRoute(kuskv1.StaticRoute) error

	GetSvc(namespace, name string) (*corev1.Service, error)
	ListServices(namespace string) (*corev1.ServiceList, error)
	ListNamespaces() (*corev1.NamespaceList, error)

	K8sClient() client.Client

	GetLogs(name, namespace string, logs chan []byte) error
}

type kuskClient struct {
	client              client.Client
	config              *rest.Config
	kuskManagedSelector labels.Selector
}

func NewClient(c client.Client, config *rest.Config) Client {
	// for use when querying apis and static routes to filter out those that are managed by kusk
	r, _ := labels.NewRequirement("kusk-managed", selection.NotIn, []string{"true"})

	return &kuskClient{
		client:              c,
		config:              config,
		kuskManagedSelector: labels.NewSelector().Add(*r),
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
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

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
	if err := k.client.List(
		context.TODO(),
		list,
		&client.ListOptions{
			Namespace:     namespace,
			LabelSelector: k.kuskManagedSelector,
		},
	); err != nil {
		return nil, err
	}

	return list, nil
}

func (k *kuskClient) GetApi(namespace, name string) (*kuskv1.API, error) {
	api := &kuskv1.API{}

	if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, api); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return api, nil
}

// GetApiByFleet gets all APIs associated with the EnvoyFleet
func (k *kuskClient) GetApiByEnvoyFleet(namespace, fleetNamespace, fleetName string) (*kuskv1.APIList, error) {
	list := kuskv1.APIList{}
	if err := k.client.List(
		context.TODO(),
		&list,
		&client.ListOptions{
			Namespace:     namespace,
			LabelSelector: k.kuskManagedSelector,
		},
	); err != nil {
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

func (k *kuskClient) UpdateApi(namespace, name, openapispec string, fleetName string, fleetnamespace string) (*kuskv1.API, error) {
	api := &kuskv1.API{}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of API before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, api); err != nil {
			return err
		}

		api.Spec = kuskv1.APISpec{
			Spec: openapispec,
			Fleet: &kuskv1.EnvoyFleetID{
				Name:      fleetName,
				Namespace: fleetnamespace,
			},
		}

		if err := k.client.Update(context.TODO(), api, &client.UpdateOptions{}); err != nil {
			return err
		}

		return nil
	})

	return api, retryErr
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
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

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
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return staticRoute, nil
}

func (k *kuskClient) GetStaticRoutes(namespace string) (*kuskv1.StaticRouteList, error) {
	list := &kuskv1.StaticRouteList{}
	if err := k.client.List(
		context.TODO(),
		list,
		&client.ListOptions{
			Namespace:     namespace,
			LabelSelector: k.kuskManagedSelector,
		},
	); err != nil {
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

func (k *kuskClient) UpdateStaticRoute(namespace, name, fleetName, fleetNamespace, specs string) (*kuskv1.StaticRoute, error) {
	staticRoute := &kuskv1.StaticRoute{}

	// marshal the paths and hosts separately for the static route
	// to use later
	tmp := &kuskv1.StaticRoute{}
	err := yaml.Unmarshal([]byte(specs), tmp)
	if err != nil {
		fmt.Println(err)
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := k.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, staticRoute); err != nil {
			return err
		}

		staticRoute.Spec = kuskv1.StaticRouteSpec{
			Fleet: &kuskv1.EnvoyFleetID{
				Name:      fleetName,
				Namespace: fleetNamespace,
			},
			Paths: tmp.Spec.Paths,
			Hosts: tmp.Spec.Hosts,
		}

		return k.client.Update(context.TODO(), staticRoute, &client.UpdateOptions{})
	})

	return staticRoute, retryErr
}

func (k *kuskClient) DeleteStaticRoute(sroute kuskv1.StaticRoute) error {
	return k.client.Delete(context.TODO(), &sroute, &client.DeleteOptions{})
}

func (k *kuskClient) ListNamespaces() (*corev1.NamespaceList, error) {
	list := &corev1.NamespaceList{}
	if err := k.client.List(context.TODO(), list, &client.ListOptions{}); err != nil {
		return nil, err
	}
	return list, nil
}

func (k *kuskClient) GetLogs(name, namespace string, logs chan []byte) error {
	clientset, err := kubernetes.NewForConfig(k.K8sConfig())
	if err != nil {
		return err
	}

	pods, _ := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/instance=%s", name),
	})

	var pod corev1.Pod
	if len(pods.Items) == 1 {
		pod = pods.Items[0]
	}
	count := int64(1)
	fmt.Println("TADA")
	// go func() {
	defer close(logs)
	for _, container := range pod.Spec.Containers {
		podLogOptions := corev1.PodLogOptions{
			Follow:    true,
			TailLines: &count,
			Container: container.Name,
		}
		podLogRequest := clientset.CoreV1().
			Pods(pod.Namespace).
			GetLogs(pod.Name, &podLogOptions)

		stream, err := podLogRequest.Stream(context.Background())
		if err != nil {
			fmt.Println("stream error", "error", err)
			continue
		}
		reader := bufio.NewReader(stream)
		for {
			b, err := ReadLongLine(reader)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			fmt.Println("blah blah TailPodLogs stream scan", "out", string(b), "pod", pod.Name)
			logs <- b
		}

		if err != nil {
			fmt.Println("scanner error", "error", err)
		}
	}
	// }()
	return nil
}

func (k *kuskClient) K8sClient() client.Client {
	return k.client
}

func (k *kuskClient) K8sConfig() *rest.Config {
	return k.config
}

// ReadLongLine reads long line
func ReadLongLine(r *bufio.Reader) (line []byte, err error) {
	var buffer []byte
	var isPrefix bool

	for {
		buffer, isPrefix, err = r.ReadLine()
		line = append(line, buffer...)
		if err != nil {
			break
		}

		if !isPrefix {
			break
		}
	}

	return line, err
}
