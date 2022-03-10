package k8sclient

import (
	"context"
	"os"
	"path"

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClient struct {
	Client client.Client
}

func NewK8sClient() (*K8sClient, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}

	k8sclient := &K8sClient{
		Client: c,
	}

	return k8sclient, nil
}

func (k *K8sClient) GetEnvoyFleets() (*kuskv1.EnvoyFleetList, error) {

	list := &kuskv1.EnvoyFleetList{}

	if err := k.Client.List(context.TODO(), list, &client.ListOptions{}); err != nil {
		return nil, err
	}
	return list, nil
}

func (k *K8sClient) GetEnvoyFleet(namespace, name string) (*kuskv1.EnvoyFleet, error) {

	envoy := &kuskv1.EnvoyFleet{}

	if err := k.Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, envoy); err != nil {
		return nil, err
	}
	return envoy, nil
}

func getConfig() (*rest.Config, error) {
	var err error
	var config *rest.Config
	k8sConfigExists := false
	homeDir, _ := os.UserHomeDir()
	cubeConfigPath := path.Join(homeDir, ".kube/config")

	if _, err := os.Stat(cubeConfigPath); err == nil {
		k8sConfigExists = true
	}

	if cfg, exists := os.LookupEnv("KUBECONFIG"); exists {
		config, err = clientcmd.BuildConfigFromFlags("", cfg)
	} else if k8sConfigExists {
		config, err = clientcmd.BuildConfigFromFlags("", cubeConfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	// default query per second is set to 5
	config.QPS = 40.0
	// default burst is set to 10
	config.Burst = 400.0

	return config, err
}
func getClient() (client.Client, error) {
	scheme := runtime.NewScheme()

	kuskv1.AddToScheme(scheme)
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return client.New(config, client.Options{Scheme: scheme})
}
