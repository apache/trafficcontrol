package dnssec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/apache/incubator-trafficcontrol/test/router/dnssec"
	"github.com/miekg/dns"
	"flag"
	"log"
)

var d *dnssec.DnssecClient
var nameserver string
var deliveryService string

func init() {
	flag.StringVar(&nameserver,"ns","changeit","ns is used to direct dns queries to a traffic router")
	flag.StringVar(&deliveryService,"ds","changeit","ds is used to target some dns DS and DNS queries made by traffic router")
}

var _ = BeforeSuite(func() {
	d = &dnssec.DnssecClient{new(dns.Client)}
	d.Net = "udp"

	Expect(nameserver).ToNot(Equal("changeit"), "Pass in a ns flag with the hostname of the traffic router")
	Expect(deliveryService).ToNot(Equal("changeit"), "Pass in a ds flag with the dns label for a DNS delivery service")
	log.Println("Nameserver",nameserver)
	log.Println("DeliveryService", deliveryService)
})

func TestDnssec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dnssec Suite")
}
