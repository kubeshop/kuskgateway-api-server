package k8sclient

import (
	"fmt"
	"strings"
	"testing"
)

func TestClientGetEnvoyFleets(t *testing.T) {
	client, err := NewK8sClient()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

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
	client, err := NewK8sClient()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
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
