package peer

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCrStates(t *testing.T) {
	t.Log("Running Peer Tests")

	text, err := ioutil.ReadFile("crstates.json")
	if err != nil {
		t.Log(err)
	}
	crStates, err := CrstatesUnMarshall(text)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(len(crStates.Caches), "caches found")
	for cacheName, crState := range crStates.Caches {
		t.Logf("%v -> %v", cacheName, crState.IsAvailable)
	}

	fmt.Println(len(crStates.Deliveryservice), "deliveryservices found")
	for dsName, deliveryService := range crStates.Deliveryservice {
		t.Logf("%v -> %v (len:%v)", dsName, deliveryService.IsAvailable, len(deliveryService.DisabledLocations))
	}

}
