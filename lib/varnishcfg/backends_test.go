package varnishcfg

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

func TestAddBackends(t *testing.T) {
	testCases := []struct {
		name             string
		backends         map[string]backend
		parents          []*atscfg.ParentAbstractionServiceParent
		originDomain     string
		originPort       int
		expectedBackends map[string]backend
	}{
		{
			name:         "no parents",
			backends:     make(map[string]backend),
			parents:      []*atscfg.ParentAbstractionServiceParent{},
			originDomain: "origin.example.com",
			originPort:   80,
			expectedBackends: map[string]backend{
				"origin_example_com_80": {
					host: "origin.example.com",
					port: 80,
				},
			},
		},
		{
			name:     "single parent",
			backends: make(map[string]backend),
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 444},
			},
			originDomain: "origin.example.com",
			originPort:   80,
			expectedBackends: map[string]backend{
				"parent_example_com_444": {
					host: "parent.example.com",
					port: 444,
				},
				"origin_example_com_80": {
					host: "origin.example.com",
					port: 80,
				},
			},
		},
		{
			name:     "multiple parent",
			backends: make(map[string]backend),
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 444},
				{FQDN: "parent2.example.com", Port: 80},
			},
			originDomain: "origin.example.com",
			originPort:   80,
			expectedBackends: map[string]backend{
				"parent_example_com_444": {
					host: "parent.example.com",
					port: 444,
				},
				"parent2_example_com_80": {
					host: "parent2.example.com",
					port: 80,
				},
				"origin_example_com_80": {
					host: "origin.example.com",
					port: 80,
				},
			},
		},
		{
			name: "already added parents",
			backends: map[string]backend{
				"parent_example_com_444": {
					host: "parent.example.com",
					port: 444,
				},
				"origin_example_com_80": {
					host: "origin.example.com",
					port: 80,
				},
			},
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 444},
			},
			originDomain: "origin.example.com",
			originPort:   80,
			expectedBackends: map[string]backend{
				"parent_example_com_444": {
					host: "parent.example.com",
					port: 444,
				},
				"origin_example_com_80": {
					host: "origin.example.com",
					port: 80,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			addBackends(tC.backends, tC.parents, tC.originDomain, tC.originPort)
			if !reflect.DeepEqual(tC.expectedBackends, tC.backends) {
				t.Errorf("expected %v got %v", tC.expectedBackends, tC.backends)
			}
		})
	}
}

func TestAddBackendsToDirector(t *testing.T) {
	testCases := []struct {
		name          string
		directorName  string
		retryPolicy   atscfg.ParentAbstractionServiceRetryPolicy
		parents       []*atscfg.ParentAbstractionServiceParent
		expectedLines []string
	}{
		{
			name:         "round robin",
			directorName: "dir",
			retryPolicy:  atscfg.ParentAbstractionServiceRetryPolicyRoundRobinStrict,
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 80},
				{FQDN: "parent2.example.com", Port: 80},
			},
			expectedLines: []string{
				`new dir = directors.round_robin();`,
				`dir.add_backend(parent_example_com_80);`,
				`dir.add_backend(parent2_example_com_80);`,
			},
		},
		{
			name:         "fallback",
			directorName: "dir",
			retryPolicy:  atscfg.ParentAbstractionServiceRetryPolicyFirst,
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 80},
				{FQDN: "parent2.example.com", Port: 80},
			},
			expectedLines: []string{
				`new dir = directors.fallback();`,
				`dir.add_backend(parent_example_com_80);`,
				`dir.add_backend(parent2_example_com_80);`,
			},
		},
		{
			name:         "fallback sticky",
			directorName: "dir",
			retryPolicy:  atscfg.ParentAbstractionServiceRetryPolicyLatched,
			parents: []*atscfg.ParentAbstractionServiceParent{
				{FQDN: "parent.example.com", Port: 80},
				{FQDN: "parent2.example.com", Port: 80},
			},
			expectedLines: []string{
				`new dir = directors.fallback(1);`,
				`dir.add_backend(parent_example_com_80);`,
				`dir.add_backend(parent2_example_com_80);`,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			lines := addBackendsToDirector(tC.directorName, tC.retryPolicy, tC.parents)
			if !reflect.DeepEqual(tC.expectedLines, lines) {
				t.Errorf("expected %v got %v", tC.expectedLines, lines)
			}
		})
	}
}

func TestAddDirectors(t *testing.T) {
	testCases := []struct {
		name                string
		subroutines         map[string][]string
		svc                 *atscfg.ParentAbstractionService
		expectedSubroutines map[string][]string
	}{
		{
			name:        "no parents",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:        "demo",
				RetryPolicy: atscfg.ParentAbstractionServiceRetryPolicyConsistentHash,
				Parents:     []*atscfg.ParentAbstractionServiceParent{},
				DestDomain:  "origin.example.com",
				Port:        80,
			},
			expectedSubroutines: map[string][]string{
				"vcl_init": {
					`new demo = directors.fallback();`,
					`demo.add_backend(origin_example_com_80);`,
				},
			},
		},
		{
			name:        "primary parents",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:        "demo",
				RetryPolicy: atscfg.ParentAbstractionServiceRetryPolicyConsistentHash,
				Parents: []*atscfg.ParentAbstractionServiceParent{
					{FQDN: "parent.example.com", Port: 80},
				},
				DestDomain: "origin.example.com",
				Port:       80,
			},
			expectedSubroutines: map[string][]string{
				"vcl_init": {
					`new demo_primary = directors.shard();`,
					`demo_primary.add_backend(parent_example_com_80);`,
					`new demo = directors.fallback();`,
					`demo.add_backend(demo_primary.backend());`,
					`demo.add_backend(origin_example_com_80);`,
				},
			},
		},
		{
			name:        "primary and secondary parents",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:        "demo",
				RetryPolicy: atscfg.ParentAbstractionServiceRetryPolicyLatched,
				Parents: []*atscfg.ParentAbstractionServiceParent{
					{FQDN: "parent.example.com", Port: 80},
				},
				SecondaryParents: []*atscfg.ParentAbstractionServiceParent{
					{FQDN: "parent2.example.com", Port: 80},
				},
				DestDomain: "origin.example.com",
				Port:       80,
			},
			expectedSubroutines: map[string][]string{
				"vcl_init": {
					`new demo_primary = directors.fallback(1);`,
					`demo_primary.add_backend(parent_example_com_80);`,
					`new demo_secondary = directors.fallback(1);`,
					`demo_secondary.add_backend(parent2_example_com_80);`,
					`new demo = directors.fallback();`,
					`demo.add_backend(demo_primary.backend());`,
					`demo.add_backend(demo_secondary.backend());`,
					`demo.add_backend(origin_example_com_80);`,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			addDirectors(tC.subroutines, tC.svc)
			if !reflect.DeepEqual(tC.expectedSubroutines, tC.subroutines) {
				t.Errorf("expected %v got %v", tC.expectedSubroutines, tC.subroutines)
			}
		})
	}
}

func TestAssignBackends(t *testing.T) {
	testCases := []struct {
		name                string
		subroutines         map[string][]string
		svc                 *atscfg.ParentAbstractionService
		requestFQDNs        []string
		expectedSubroutines map[string][]string
	}{
		{
			name:        "edge with one request FQDN",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:       "demo",
				DestDomain: "origin.example.com",
			},
			requestFQDNs: []string{"example.com"},
			expectedSubroutines: map[string][]string{
				"vcl_recv": {
					`if (req.http.host == "example.com") {`,
					`	set req.backend_hint = demo.backend();`,
					`}`,
				},
				"vcl_backend_fetch": {
					`if (bereq.http.host == "example.com") {`,
					`	set bereq.http.host = "origin.example.com";`,
					`}`,
				},
			},
		},
		{
			name:        "edge with multiple request FQDNs",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:       "demo",
				DestDomain: "origin.example.com",
			},
			requestFQDNs: []string{"example.com", "another.example.com"},
			expectedSubroutines: map[string][]string{
				"vcl_recv": {
					`if (req.http.host == "example.com" || req.http.host == "another.example.com") {`,
					`	set req.backend_hint = demo.backend();`,
					`}`,
				},
				"vcl_backend_fetch": {
					`if (bereq.http.host == "example.com" || bereq.http.host == "another.example.com") {`,
					`	set bereq.http.host = "origin.example.com";`,
					`}`,
				},
			},
		},
		{
			name:        "mid",
			subroutines: make(map[string][]string),
			svc: &atscfg.ParentAbstractionService{
				Name:       "demo",
				DestDomain: "origin.example.com",
			},
			requestFQDNs: []string{"origin.example.com"},
			expectedSubroutines: map[string][]string{
				"vcl_recv": {
					`if (req.http.host == "origin.example.com") {`,
					`	set req.backend_hint = demo.backend();`,
					`}`,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			assignBackends(tC.subroutines, tC.svc, tC.requestFQDNs)
			if !reflect.DeepEqual(tC.expectedSubroutines, tC.subroutines) {
				t.Errorf("expected %v got %v", tC.expectedSubroutines, tC.subroutines)
			}
		})
	}
}
