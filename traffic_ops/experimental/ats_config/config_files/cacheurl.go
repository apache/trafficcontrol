package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

// stripProtocol takes a URL, e.g. `http://foo.example.com/bar/` and removes the protocol and any trailing slash, e.g. `foo.example.com/bar`
// TODO test
func stripProtocol(url string) string {
	// TODO move to generic file, since it's used by other configs
	// this could be made more efficient with regexp, if necessary
	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	if len(url) == 0 {
		return url
	}
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	return url
}

// createCacheurlDotConfig constructs cacheurl.config files
func createCacheurlDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// # DO NOT EDIT - Generated for odol-atsec-atl-27 by Twelve Monkeys (https://tm.comcast.net/) on Fri Dec 23 18:40:33 UTC 2016
	// http://(odol-ip-eas-origin.g.comcast.net/[^?]+)(?:\?|$)  http://$1
	// ...

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	// paramMap := createParamsMap(params)

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}
	for _, deliveryService := range deliveryServices {
		if deliveryService.QStringIgnore == 1 {
			continue
		}
		org := stripProtocol(deliveryService.OrgServerFQDN)
		s += "$1(" + org + "/[^?]+)(?:\\?|$)  $1$1\n"
	}
	s = strings.Replace(s, `__RETURN__`, `\n`, -1)
	return s, nil
}

// createCacheurlQstringDotConfig constructs cacheurl_qstring.config files
func createCacheurlQstringDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	return fmt.Sprintf(`# DO NOT EDIT - Generated for %s by Traffic Ops (%s) on %s
http://([^?]+)(?:\?|$)  http://$1
https://([^?]+)(?:\?|$)  https://$1
`,
		trafficServerHost,
		trafficOpsHost,
		time.Now().String(),
	), nil
}

// deliveryServiceByXMLID returns a traffic_ops/client.DeliveryService from its XML ID (name)
func deliveryServiceByXMLID(toClient towrap.ITrafficOpsSession, xmlId string) (to.DeliveryService, error) {
	// TODO add Traffic Ops endpoint to get a DS by name, so we don't have to fetch all delivery services.

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return to.DeliveryService{}, fmt.Errorf("error getting delivery services: %v", err)
	}

	for _, ds := range deliveryServices {
		if ds.XMLID == xmlId {
			return ds, nil
		}
	}
	return to.DeliveryService{}, fmt.Errorf("not found")
}

// createCacheurlStarDotConfig constructs cacheurl_(.*).config files
func createCacheurlStarDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	deliveryServiceXmlId := strings.TrimSuffix(strings.TrimPrefix(filename, "cacheurl_"), ".config")

	ds, err := deliveryServiceByXMLID(toClient, deliveryServiceXmlId)
	if err != nil {
		return s, nil // if the requested DS isn't found, we return the file without it, without error. This mirrors the old Traffic Ops behavior
	}
	s += ds.CacheURL + "\n"
	s = strings.Replace(s, `__RETURN__`, `\n`, -1)
	return s, nil
}
