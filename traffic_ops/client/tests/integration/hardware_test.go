package integration

import (
	"encoding/json"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestHardware(t *testing.T) {
	resp, err := Request(*to, "GET", "/api/1.2/hwinfo.json", nil)
	if err != nil {
		t.Errorf("Could not get hardware.json reponse was: %v\n", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiHwRes traffic_ops.HardwareResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiHwRes); err != nil {
		t.Errorf("Could not decode hwinfo json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiHardware := apiHwRes.Response

	clientHardware, err := to.Hardware(0)
	if err != nil {
		t.Errorf("Could not get Hardware from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiHardware) != len(clientHardware) {
		t.Errorf("Hardware Response Length -- expected %v, got %v\n", len(apiHardware), len(clientHardware))
	}

	for _, apiHw := range apiHardware {
		match := false
		for _, clientHw := range clientHardware {
			if apiHw.HostName == clientHw.HostName && apiHw.Description == clientHw.Description {
				match = true
				compareHw(apiHw, clientHw, t)
			}
		}
		if !match {
			t.Errorf("No hardware information returned from client for %s\n", apiHw.HostName)
		}
	}
}

func TestHardwareWithLimit(t *testing.T) {
	resp, err := Request(*to, "GET", "/api/1.2/hwinfo.json?limit=10", nil)
	if err != nil {
		t.Errorf("Could not get hwinfo.json reponse was: %v\n", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiHwRes traffic_ops.HardwareResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiHwRes); err != nil {
		t.Errorf("Could not decode hwinfo json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiHardware := apiHwRes.Response

	clientHardware, err := to.Hardware(10)
	if err != nil {
		t.Errorf("Could not get Hardware from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiHardware) != len(clientHardware) {
		t.Errorf("Hardware Response Length -- expected %v, got %v\n", len(apiHardware), len(clientHardware))
	}

	for _, apiHw := range apiHardware {
		match := false
		for _, clientHw := range clientHardware {
			if apiHw.HostName == clientHw.HostName && apiHw.Description == clientHw.Description {
				match = true
				compareHw(apiHw, clientHw, t)
			}
		}
		if !match {
			t.Errorf("No hardware information returned from client for %s\n", apiHw.HostName)
		}
	}
}

func compareHw(hw1 traffic_ops.Hardware, hw2 traffic_ops.Hardware, t *testing.T) {
	if hw1.Description != hw2.Description {
		t.Errorf("Description -- Expected %v, got %v", hw1.Description, hw2.Description)
	}
	if hw1.HostName != hw2.HostName {
		t.Errorf("HostName -- Expected %v, got %v", hw1.HostName, hw2.HostName)
	}
	if hw1.ID != hw2.ID {
		t.Errorf("ID -- Expected %v, got %v", hw1.ID, hw2.ID)
	}
	if hw1.LastUpdated != hw2.LastUpdated {
		t.Errorf("LastUpdated -- Expected %v, got %v", hw1.LastUpdated, hw2.LastUpdated)
	}
	if hw1.Value != hw2.Value {
		t.Errorf("Value -- Expected %v, got %v", hw1.Value, hw2.Value)
	}
}
