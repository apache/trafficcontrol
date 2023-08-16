package varnishcfg

import (
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestGetHitchConfig(t *testing.T) {
	ds1 := &atscfg.DeliveryService{}
	ds1.XMLID = util.StrPtr("ds1")
	ds1.Protocol = util.IntPtr(1)
	ds1Type := tc.DSTypeHTTP
	ds1.Type = &ds1Type
	ds1.ExampleURLs = []string{"https://ds1.example.org"}
	deliveryServices := []atscfg.DeliveryService{*ds1}
	txt, warnings := GetHitchConfig(deliveryServices, "/ssl")
	expectedTxt := strings.Join([]string{
		`frontend = {`,
		`	host = "*"`,
		`	port = "443"`,
		`}`,
		`backend = "[127.0.0.1]:6081"`,
		`write-proxy-v2 = on`,
		`user = "root"`,
		`pem-file = {`,
		`	cert = "/ssl/ds1_example_org_cert.cer"`,
		`	private-key = "/ssl/ds1.example.org.key"`,
		`}`,
	}, "\n")
	expectedTxt += "\n"
	if len(warnings) != 0 {
		t.Errorf("expected no warnings got %v", warnings)
	}
	if txt != expectedTxt {
		t.Errorf("expected: %s got: %s", expectedTxt, txt)
	}
}
