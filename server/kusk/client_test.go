package k8sclient

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var testClient Client

func setup(tb testing.TB) {
	if _, fakeIt := os.LookupEnv("FAKE"); fakeIt {
		testClient = NewClient(getFakeClient())
		return
	}
	k8sclient, err := getClient()
	if err != nil {
		tb.Error(err)
		tb.Fail()
		return
	}

	testClient = NewClient(k8sclient)
}
func TestClientGetEnvoyFleets(t *testing.T) {
	setup(t)

	fleets, err := testClient.GetEnvoyFleets()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if len(fleets.Items) == 0 {
		t.Error("no data returned")
		t.Fail()
		return
	}
}

func TestClientGetEnvoyFleet(t *testing.T) {
	setup(t)

	name := "default"
	namespace := "default"
	fleet, err := testClient.GetEnvoyFleet(namespace, name)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf(`envoyfleet.gateway.kusk.io "%s" not found`, name)) {
			t.Error(err)
			t.Fail()
			return
		}
		t.Error(err)
		t.Fail()
		return
	}

	if fleet.ObjectMeta.Name != name {
		t.Error("name does not match")
		t.Fail()
		return
	}
}

func TestGetApis(t *testing.T) {
	setup(t)
	apis, err := testClient.GetApis("default")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	fmt.Println(len(apis.Items))
}

func TestGetApi(t *testing.T) {
	setup(t)
	api, err := testClient.GetApi("default", "sample")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	fmt.Println(api.Spec.Spec)
}

func TestGetSvc(t *testing.T) {
	setup(t)
	_, err := testClient.GetSvc("default", "kubernetes")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

}

func getFakeClient() client.Client {

	scheme := runtime.NewScheme()
	kuskv1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)

	initObjects := []client.Object{
		&kuskv1.API{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: "default",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubernetes",
				Namespace: "default",
			},
		},
		&kuskv1.EnvoyFleet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: "default",
			},
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(initObjects...).
		Build()
	return fakeClient
}

func getClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	kuskv1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return client.New(config, client.Options{Scheme: scheme})
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
