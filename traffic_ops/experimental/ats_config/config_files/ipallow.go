package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func createIpallowDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// # DO NOT EDIT - Generated for my-edge-0 by Twelve Monkeys (https://tm.example.net/) on Thu Dec 22 23:33:13 UTC 2016
	// 	# 12M NOTE: This is running with forced volumes - the size is irrelevant
	// 	volume=1 scheme=http size=50%
	// 	volume=2 scheme=http size=50%

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	ipAllowData, err := getIPAllowData(toClient, trafficServerHost, params)
	if err != nil {
		return "", fmt.Errorf("error getting IP allow data: %v", err)
	}

	for _, d := range ipAllowData {
		s += fmt.Sprintf("src_ip=%-70s action=%-10s method=%-20s\n", d.SourceIP, d.Action, d.Method)
	}
	return s, nil
}

type IPAllowData struct {
	SourceIP string
	Action   string
	Method   string
}

func getIPAllowData(toClient towrap.ITrafficOpsSession, trafficServerHost string, params []to.Parameter) ([]IPAllowData, error) {
	allow := []IPAllowData{}

	allow = append(allow, IPAllowData{
		SourceIP: "127.0.0.1",
		Action:   "ip_allow",
		Method:   "ALL",
	})

	allow = append(allow, IPAllowData{
		SourceIP: "::1",
		Action:   "ip_allow",
		Method:   "ALL",
	})

	for _, param := range params {
		if param.Name == "purge_allow_ip" && param.ConfigFile == "ip_allow.config" {
			allow = append(allow, IPAllowData{
				SourceIP: param.Value,
				Action:   "ip_allow",
				Method:   "ALL",
			})
		}
	}

	serverTypeStr, err := getServerTypeStr(toClient, trafficServerHost)
	if err != nil {
		return nil, fmt.Errorf("error getting server '%v' type: %v", trafficServerHost, err)
	}
	isMid := strings.HasPrefix(serverTypeStr, "MID")
	if isMid {
		midAllows, err := getIPAllowDataMid(toClient, params, trafficServerHost)
		if err != nil {
			return nil, fmt.Errorf("error getting mid IP allows: %v", err)
		}
		allow = append(allow, midAllows...)
	} else {
		allow = append(allow, getIPAllowDataEdge(params)...)
	}

	return allow, nil
}

func getIPAllowDataMid(toClient towrap.ITrafficOpsSession, params []to.Parameter, serverHost string) ([]IPAllowData, error) {
	allow := []IPAllowData{}
	cacheGroups, err := toClient.CacheGroups()
	if err != nil {
		return nil, fmt.Errorf("getting cachegroups: %v", err)
	}

	server, err := getServer(toClient, serverHost)
	if err != nil {
		return nil, fmt.Errorf("getting server %s: %v", serverHost, err)
	}

	servers, err := toClient.Servers()
	if err != nil {
		return nil, fmt.Errorf("getting servers: %v", err)
	}

	allowedCachegroups := map[string]struct{}{}
	for _, cg := range cacheGroups {
		if cg.ParentName == server.Cachegroup {
			allowedCachegroups[cg.Name] = struct{}{}
		}
	}

	allowedIps := []string{}
	allowedIp6s := []string{}
	for _, allowedServer := range servers {
		if !strings.HasPrefix(allowedServer.Type, "EDGE") && allowedServer.Type != "RASCAL" {
			continue
		}
		if _, ok := allowedCachegroups[allowedServer.Cachegroup]; !ok {
			continue
		}

		if allowedServer.IPAddress != "" {
			// TODO handle subnet
			allowedIps = append(allowedIps, allowedServer.IPAddress)
		} else {
			// TODO log(allowedServer + " has an invalid IPv4 address; excluding from ip_allow data for " + server)
		}
		if allowedServer.IP6Address != "" {
			allowedIp6s = append(allowedIp6s, allowedServer.IP6Address)
		}
	}

	// TODO compact/coalesce

	for _, ip := range allowedIps {
		ip = strings.Trim(ip, " ")
		allow = append(allow, IPAllowData{
			SourceIP: ip,
			Action:   "ip_allow",
			Method:   "ALL",
		})
	}

	for _, ip := range allowedIp6s {
		ip = strings.Trim(ip, " ")
		allow = append(allow, IPAllowData{
			SourceIP: ip,
			Action:   "ip_allow",
			Method:   "ALL",
		})
	}

	// allow RFC 1918 server space - TODO JvD: parameterize

	allow = append(allow, IPAllowData{
		SourceIP: "10.0.0.0-10.255.255.255",
		Action:   "ip_allow",
		Method:   "ALL",
	})
	allow = append(allow, IPAllowData{
		SourceIP: "172.16.0.0-172.31.255.255",
		Action:   "ip_allow",
		Method:   "ALL",
	})
	allow = append(allow, IPAllowData{
		SourceIP: "192.168.0.0-192.168.255.255",
		Action:   "ip_allow",
		Method:   "ALL",
	})
	allow = append(allow, IPAllowData{
		SourceIP: "0.0.0.0-255.255.255.255",
		Action:   "ip_deny",
		Method:   "ALL",
	})
	allow = append(allow, IPAllowData{
		SourceIP: "::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		Action:   "ip_deny",
		Method:   "ALL",
	})
	return allow, nil
}

func getIPAllowDataEdge(params []to.Parameter) []IPAllowData {
	// for edges deny "PUSH|PURGE|DELETE", allow everything else to everyone.
	allow := []IPAllowData{}
	allow = append(allow, IPAllowData{
		SourceIP: "0.0.0.0-255.255.255.255",
		Action:   "ip_deny",
		Method:   "PUSH|PURGE|DELETE",
	})
	allow = append(allow, IPAllowData{
		SourceIP: "::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		Action:   "ip_deny",
		Method:   "PUSH|PURGE|DELETE",
	})
	return allow
}
