package varnishcfg

import (
	"path/filepath"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
)

// GetHitchConfig returns Hitch config using TO data
func GetHitchConfig(deliveryServices []atscfg.DeliveryService, sslDir string) (string, []string) {
	warnings := make([]string, 0)
	lines := []string{
		`frontend = {`,
		`	host = "*"`,
		`	port = "443"`,
		`}`,
		`backend = "[127.0.0.1]:6081"`,
		`write-proxy-v2 = on`,
		// TODO: change root user
		`user = "root"`,
	}

	dses, dsWarns := atscfg.DeliveryServicesToSSLMultiCertDSes(deliveryServices)
	warnings = append(warnings, dsWarns...)

	dses = atscfg.GetSSLMultiCertDotConfigDeliveryServices(dses)

	for dsName, ds := range dses {
		cerName, keyName := atscfg.GetSSLMultiCertDotConfigCertAndKeyName(dsName, ds)
		lines = append(lines, []string{
			`pem-file = {`,
			`	cert = "` + filepath.Join(sslDir, cerName) + `"`,
			`	private-key = "` + filepath.Join(sslDir, keyName) + `"`,
			`}`,
		}...)
	}

	txt := strings.Join(lines, "\n")
	txt += "\n"
	return txt, warnings
}
