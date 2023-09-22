package dnssec_test

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
	// . "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	"flag"
	"testing"

	"github.com/apache/trafficcontrol/v8/test/router/dnssec"
	"github.com/miekg/dns"
)

var d *dnssec.DnssecClient
var nameserver string
var deliveryService string

func init() {
	flag.StringVar(&nameserver, "ns", "changeit", "ns is used to direct dns queries to a traffic router")
	flag.StringVar(&deliveryService, "ds", "changeit", "ds is used to target some dns DS and DNS queries made by traffic router")
}

// var _ = BeforeSuite(func() {
// 	d = &dnssec.DnssecClient{new(dns.Client)}
// 	d.Net = "udp"

// 	Expect(nameserver).ToNot(Equal("changeit"), "Pass in a ns flag with the hostname of the traffic router")
// 	Expect(deliveryService).ToNot(Equal("changeit"), "Pass in a ds flag with the dns label for a DNS delivery service")
// 	log.Println("Nameserver", nameserver)
// 	log.Println("DeliveryService", deliveryService)
// })

// func TestDnssec(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Dnssec Suite")
// }

func TestDNSSEC(t *testing.T) {
	d = &dnssec.DnssecClient{Client: new(dns.Client)}
	d.Net = "udp"

	if nameserver == "changeit" {
		t.Fatal("Pass in a ns flag with the hostname of th etraffic router")
	}
	if deliveryService == "changeit" {
		t.Fatal("Pass in a ds flag with the dns label for a DNS delivery service")
	}

	t.Logf("Nameserver: %s", nameserver)
	t.Logf("DeliveryService: %s", deliveryService)
}
