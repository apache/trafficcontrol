package main

import (
	"fmt"
	"testing"
)

// Succeed is the Unicode codepoint for a check mark.
const Succeed = "\u2713"

// Failed is the Unicode codepoint for an X mark.
const Failed = "\u2717"

func TestInfluxConnect(t *testing.T) {
	startupConfig := StartupConfig{}

	runningConfig := RunningConfig{
		InfluxDBProps: []InfluxDBProps{
			InfluxDBProps{
				Fqdn: "google.com",
				Port: 80,
			},
		},
	}

	t.Log("Given the need to test a successful InfluxDB connection.")
	{
		influxClient, err := influxConnect(&startupConfig, &runningConfig)
		if err != nil {
			t.Fatal("\tShould be able to connect to InfluxDB", Failed)
		}
		t.Log("\tShould be able to connect to InfluxDB", Succeed)

		if influxClient.Addr() != fmt.Sprintf("http://%s:%d", runningConfig.InfluxDBProps[0].Fqdn, runningConfig.InfluxDBProps[0].Port) {
			t.Fatal("\tShould get back \"http://google.com:80\" for \"influxClient.Addr()\".", Failed)
		} else {
			t.Log("\tShould get back \"http://google.com:80\" for \"influxClient.Addr()\".", Succeed)
		}
	}
}

func TestRandInfluxConnect(t *testing.T) {
	startupConfig := StartupConfig{}

	runningConfig := RunningConfig{
		InfluxDBProps: []InfluxDBProps{
			InfluxDBProps{
				Fqdn: "google.com",
				Port: 80,
			},
			InfluxDBProps{
				Fqdn: "golang.org",
				Port: 80,
			},
			InfluxDBProps{
				Fqdn: "godoc.org",
				Port: 80,
			},
		},
	}

	t.Log("Given the need to test a successful connection to a random InfluxDB server.")
	{
		influxClient1, err := influxConnect(&startupConfig, &runningConfig)
		if err != nil {
			t.Fatal("\tShould be able to connect to InfluxDB", Failed)
		}
		t.Log("\tShould be able to connect to InfluxDB", Succeed)

		influxClient2, err := influxConnect(&startupConfig, &runningConfig)
		if err != nil {
			t.Fatal("\tShould be able to connect to InfluxDB", Failed)
		}
		t.Log("\tShould be able to connect to InfluxDB", Succeed)

		if influxClient1.Addr() == influxClient2.Addr() {
			t.Fatal("\tShould not connect to the same InfluxDB server twice in a row.", Failed)
		} else {
			t.Log("\tShould not connect to the same InfluxDB server twice in a row.", Succeed)
		}
	}
}
