package k8sclient

import (
	"fmt"
	"strings"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClientGetEnvoyFleets(t *testing.T) {
	fakeClient := fake.NewClientBuilder().Build()
	client := NewClient(fakeClient)

	fleets, err := client.GetEnvoyFleets()
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
	fakeClient := fake.NewClientBuilder().Build()
	client := NewClient(fakeClient)

	name := "default"
	namespace := "default"
	fleet, err := client.GetEnvoyFleet(namespace, name)
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
