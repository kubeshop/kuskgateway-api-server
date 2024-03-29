package kusk

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kubeshop/kusk-gateway/api/v1alpha1"
	kuskv1 "github.com/kubeshop/kusk-gateway/api/v1alpha1"
)

var testClient Client

func setup(t *testing.T) {
	if fake, isFakeDefined := os.LookupEnv("FAKE"); isFakeDefined && (fake == "true" || fake == "TRUE" || fake == "1") {
		testClient = NewClient(getFakeClient())
		return
	}

	k8sclient, err := getClient()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	testClient = NewClient(k8sclient)
}

func TestCreateEnvoyFleet(t *testing.T) {
	require := require.New(t)
	setup(t)
	fleet := v1alpha1.EnvoyFleet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: kuskv1.EnvoyFleetSpec{
			Service: &kuskv1.ServiceConfig{
				Type: corev1.ServiceTypeLoadBalancer,
			},
		},
	}

	_, err := testClient.CreateFleet(fleet)
	require.NoError(err)
}

func TestDeleteFleet(t *testing.T) {
	require := require.New(t)
	setup(t)

	fleet := v1alpha1.EnvoyFleet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
	}

	require.NoError(testClient.DeleteFleet(fleet))
}
func TestClientGetEnvoyFleets(t *testing.T) {
	require := require.New(t)
	setup(t)

	fleets, err := testClient.GetEnvoyFleets()
	require.NoError(err)

	require.NotEqual(len(fleets.Items), 0)
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
func TestGetNotFoundApi(t *testing.T) {
	setup(t)
	_, err := testClient.GetApi("default", "not-found")
	if err == ErrNotFound {
		return
	}
	t.Errorf("%s does not equal to ErrNotFound", err)
	t.Fail()
}

func TestDeleteAPI(t *testing.T) {
	require := require.New(t)
	setup(t)

	err := testClient.DeleteAPI("default", "sample")
	require.NoError(err)
}

func TestUpdateAPI(t *testing.T) {
	require := require.New(t)

	tc := NewClient(getFakeClient())

	_, err := tc.UpdateApi("default", "non-existent", "", "test", "default")
	require.Error(err)

	_, err = tc.UpdateApi("default", "sample", "", "test", "default")
	require.NoError(err)
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

func TestListServices(t *testing.T) {
	setup(t)
	_, err := testClient.ListServices("default")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestCreateStaticRoute(t *testing.T) {
	require := require.New(t)

	setup(t)

	namespace := "default"
	name := "static-route-example-1-top-level-upstream"
	fleetNamespace := "default"
	fleetName := "default"
	specs := `
spec:
  upstream:
    service:
      name: static-route-example-1-top-level-upstream
      namespace: default
      port: 80
`
	staticRoute, err := testClient.CreateStaticRoute(namespace, name, fleetNamespace, fleetName, specs)

	require.NoError(err)
	require.NotNil(staticRoute)

	require.Equal(fleetNamespace+"."+fleetName, staticRoute.Spec.Fleet.String())
}

func TestGetStaticRoutes(t *testing.T) {
	require := require.New(t)

	setup(t)

	namespace := "default"
	staticRoutes, err := testClient.GetStaticRoutes(namespace)

	require.NoError(err)
	require.NotNil(staticRoutes)
	require.Len(staticRoutes.Items, 1)
}

func TestDeleteStaticRoute(t *testing.T) {
	require := require.New(t)

	setup(t)
	name := "static-route-1"
	namespace := "default"
	err := testClient.DeleteStaticRoute(v1alpha1.StaticRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})

	require.NoError(err)
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
		&kuskv1.StaticRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "static-route-1",
				Namespace: "default",
			},
			Spec: kuskv1.StaticRouteSpec{
				Fleet: &kuskv1.EnvoyFleetID{
					Name:      "/static-route-1",
					Namespace: "default",
				},
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
	if err := kuskv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
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
