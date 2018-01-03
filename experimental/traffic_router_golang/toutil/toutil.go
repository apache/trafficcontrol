package toutil

import (
	"errors"
	"strconv"

	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

const TrafficMonitorType = "RASCAL"

func MonitorAvailableStatuses() map[string]struct{} {
	return map[string]struct{}{
		"ONLINE":   struct{}{},
		"REPORTED": struct{}{},
	}
}

func GetMonitorURIs(toc *to.Session, cdn string) ([]string, error) {
	servers, reqInf, err := toc.GetServersByType(map[string][]string{"type": {TrafficMonitorType}})
	if err != nil {
		return nil, errors.New("getting servers by type from '" + reqInf.RemoteAddr.String() + "':" + err.Error())
	}
	availableStatuses := MonitorAvailableStatuses()

	monitors := []string{}
	for _, server := range servers {
		if server.CDNName != cdn {
			continue
		}
		if _, ok := availableStatuses[server.Status]; !ok {
			continue
		}
		m := "http://" + server.HostName + "." + server.DomainName
		if server.TCPPort > 0 && server.TCPPort != 80 {
			m += ":" + strconv.Itoa(server.TCPPort)
		}
		monitors = append(monitors, m)
	}
	return monitors, nil
}
