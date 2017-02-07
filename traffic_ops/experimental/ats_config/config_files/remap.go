package config_files

import (
	"encoding/json"
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

func edgeHdrRwFile(ds to.DeliveryService) string {
	return fmt.Sprintf("hdr_rw_%s.config", ds.XMLID)
}

func midHdrRwFile(ds to.DeliveryService) string {
	return fmt.Sprintf("hdr_rw_mid_%s.config", ds.XMLID)
}

func cacheurlFile(ds to.DeliveryService) string {
	return fmt.Sprintf("cacheurl_%s.config", ds.XMLID)
}

func createRemapDotConfigMid(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	midRemap := map[string]string{} // map[originFQDN]remapStr
	for _, ds := range deliveryServices {
		orgFqdn := ds.OrgServerFQDN

		isLiveLocal := strings.HasSuffix(ds.Type, "_LIVE") // if '_Live' is the suffix, it's not National ('_LIVE_NATNL')

		if isLiveLocal {
			continue // Live local delivery services skip mids
		}
		if _, ok := midRemap[orgFqdn]; ok {
			continue
		}

		if ds.MidHeaderRewrite != "" {
			midRemap[orgFqdn] += fmt.Sprintf(" @plugin=header_rewrite.so @pparam=%s", midHdrRwFile(ds))
		}
		if ds.QStringIgnore == 1 {
			midRemap[orgFqdn] += " @plugin=cacheurl.so @pparam=cacheurl_qstring.config"
		}
		// TODO warn for invalid QStringIgnore
		if ds.CacheURL != "" {
			midRemap[orgFqdn] += fmt.Sprintf(" @plugin=cacheurl.so @pparam=", cacheurlFile(ds))
		}
		if ds.RangeRequestHandling == 2 {
			midRemap[orgFqdn] += " @plugin=cache_range_requests.so"
		}
		// TODO warn for invalid RangeRequestHandling
	}

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"
	for originFqdn, remapStr := range midRemap {
		s += fmt.Sprintf("map %s %s %s\n", originFqdn, originFqdn, remapStr)
	}
	return s, nil
}

// type DeliveryServiceMatch struct {
// 	Type      string `json:"type"`
// 	SetNumber string `json:"setNumber"`
// 	Pattern   string `json:"pattern"`
// }

func getDomain(paramMap map[string]map[string]string) (string, error) {
	if crcParams, ok := paramMap["CRConfig.json"]; ok {
		if domainName, ok := crcParams["domain_name"]; ok {
			return domainName, nil
		}
		return "", fmt.Errorf("domain_name parameter does not exist")
	}
	return "", fmt.Errorf("CRConfig.json parameter file does not exist")
}

// remapLines returns a map[mapTo]mapFrom. Note this requires the regexes, which as of this writing (2016-12) are not populated in the API endpoint, and hence the DeliveryService must be retreived from the CRConfig.
func getRemapLines(server to.Server, ds to.DeliveryService, matchset Matchset, paramMap map[string]map[string]string) (map[string]string, error) {
	domainName, err := getDomain(paramMap)
	if err != nil {
		return nil, fmt.Errorf("error getting domain: %v", err)
	}

	remapLines := map[string]string{}
	// TODO fix to.client to populate MatchList?
matchsetFor:
	for _, dsMatch := range matchset.MatchList {
		// TODO determine if HOST should be HOST_REGEXP and from somewhere else?
		if dsMatch.MatchType != "HOST" || ds.Type == "ANY_MAP" {
			continue matchsetFor
		}
		hostRe := dsMatch.Regex
		mapTo := ds.OrgServerFQDN
		if strings.HasSuffix(hostRe, `.*`) {
			re := hostRe
			re = strings.Replace(re, `\`, ``, -1)
			re = strings.Replace(re, `.*`, ``, -1)

			hname := ""
			if strings.HasPrefix(ds.Type, "DNS") {
				hname = "edge"
			} else {
				hname = "ccr"
			}

			portStr := ""
			serverPort := server.TCPPort
			if err != nil {
				return nil, fmt.Errorf("server %v port %v is not a number", server.HostName, server.TCPPort)
			}
			if hname == "ccr" && serverPort > 0 && serverPort != 80 {
				portStr = fmt.Sprintf(":%d", serverPort)
			}

			if ds.Protocol == 0 || ds.Protocol == 2 {
				mapFrom := fmt.Sprintf(`http://%s%s%s%s/`, hname, re, domainName, portStr)
				remapLines[mapFrom] = mapTo
			}
			if ds.Protocol == 1 || ds.Protocol == 3 {
				mapFrom := fmt.Sprintf(`https://%s%s%s/`, hname, re, domainName)
				remapLines[mapFrom] = mapTo
			}
			// TODO log invalid protocol
		} else {
			if ds.Protocol == 0 || ds.Protocol == 2 {
				mapFrom := fmt.Sprintf(`http://%s/`, hostRe)
				remapLines[mapFrom] = mapTo
			}
			if ds.Protocol == 1 || ds.Protocol == 3 {
				mapFrom := fmt.Sprintf(`https://%s/`, hostRe)
				remapLines[mapFrom] = mapTo
			}
			// TODO log invalid protocol
		}
	}
	return remapLines, nil
}

func getServer(toClient towrap.ITrafficOpsSession, serverToFind string) (to.Server, error) {
	// TODO add TO endpoint to get a single server's data, for efficiency.
	servers, err := toClient.Servers()
	if err != nil {
		return to.Server{}, fmt.Errorf("error getting servers from Traffic Ops: %v", err)
	}
	for _, server := range servers {
		if server.HostName == serverToFind {
			return server, nil
		}
	}
	return to.Server{}, fmt.Errorf("not found")
}

type Matchset struct {
	Protocol  string `json:"protocol"`
	MatchList []struct {
		Regex     string `json:"regex"`
		MatchType string `json:"match-type"`
	} `json:"matchlist"`
}
type CRConfigDeliveryService struct {
	Matchsets []Matchset `json:"matchsets"`
}

// CRConfig is the CrConfig data needed by TOData. Note this is not all data in the CRConfig.
// TODO change strings to type?
type CRConfig struct {
	ContentServers map[string]struct {
		DeliveryServices map[string][]string `json:"deliveryServices"`
		CacheGroup       string              `json:"cacheGroup"`
		Type             string              `json:"type"`
	} `json:"contentServers"`
	DeliveryServices map[string]CRConfigDeliveryService `json:"deliveryServices"`
}

// getServerDeliveryServices returns the Delivery Services for the given server with their regex data. This works around TO not having an API to get DS regexes.
// TODO add TO API for DS regexes
func getServerDeliveryServices(toClient towrap.ITrafficOpsSession, server to.Server) (map[string]CRConfigDeliveryService, error) {
	crConfigBytes, err := toClient.CRConfigRaw(server.CDNName)
	if err != nil {
		return nil, fmt.Errorf("error getting CRConfig: %v", err)
	}

	crConfig := CRConfig{}
	if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
		return nil, fmt.Errorf("Error unmarshalling CRConfig: %v", err)
	}

	serverDeliveryServiceNames := crConfig.ContentServers[server.HostName]
	serverDeliveryServices := map[string]CRConfigDeliveryService{}
	for name, _ := range serverDeliveryServiceNames.DeliveryServices {
		serverDeliveryServices[name] = crConfig.DeliveryServices[name]
	}
	return serverDeliveryServices, nil
}

// TODO change to use param map, for efficiency
func hasDscpRemap(paramMap map[string]map[string]string) bool {
	if packageParams, ok := paramMap["package"]; ok {
		if _, ok := packageParams["dscp_remap"]; ok {
			return true
		}
	}
	return false
}

func buildRemapLine(toClient towrap.ITrafficOpsSession, trafficServerHost string, ds to.DeliveryService, matchset Matchset, paramMap map[string]map[string]string, mapFrom string, mapTo string) string {
	if ds.Type == "ANY_MAP" {
		return fmt.Sprintf("%s\n", ds.RemapText)
	}

	s := ""
	if hasDscpRemap(paramMap) {
		s += fmt.Sprintf("map %s     %s @plugin=dscp_remap.so @pparam=%s", mapFrom, mapTo, ds.DSCP)
	} else {
		s += fmt.Sprintf(`map	%s     %s @plugin=header_rewrite.so @pparam=dscp/set_dscp_%s.config`, mapFrom, mapTo, ds.DSCP)
	}

	if ds.EdgeHeaderRewrite != "" {
		s += fmt.Sprintf(" @plugin=header_rewrite.so @pparam=%s", edgeHdrRwFile(ds))
	}

	if ds.Signed {
		s += fmt.Sprintf(` @plugin=url_sig.so @pparam=url_sig_%s.config`, ds.XMLID)
	}

	if ds.QStringIgnore == 2 {
		dqsFile := "drop_qstring.config"
		s += fmt.Sprintf(` @plugin=regex_remap.so @pparam=%s`, dqsFile)
	} else if ds.QStringIgnore == 1 {
		globalExists := false
		if locationParams, ok := paramMap["location"]; ok {
			if _, cacheUrlExists := locationParams["cacheurl.config"]; cacheUrlExists {
				globalExists = true
			}
		}
		if globalExists {
			// log("qstring_ignore == 1, but global cacheurl.config param exists, so skipping remap rename config_file=cacheurl.config parameter if you want to change")
		} else {
			s += ` @plugin=cacheurl.so @pparam=cacheurl_qstring.config`
		}
	}

	if ds.CacheURL != "" {
		s += fmt.Sprintf(` @plugin=cacheurl.so @pparam=%s`, cacheurlFile(ds))
	}

	// Note: should use full path here?
	if ds.RegexRemap != "" {
		s += fmt.Sprintf(` @plugin=regex_remap.so @pparam=regex_remap_%s.config`, ds.XMLID)
	}

	if ds.RangeRequestHandling == 1 {
		s += ` @plugin=background_fetch.so @pparam=bg_fetch.config`
	} else if ds.RangeRequestHandling == 2 {
		s += ` @plugin=cache_range_requests.so `
	}

	if ds.RemapText != "" {
		s += fmt.Sprintf(` %s`, ds.RemapText)
	}

	s += "\n"
	return s
}

func getDS(name string, dses []to.DeliveryService) (to.DeliveryService, error) {
	// TODO put in map
	for _, ds := range dses {
		if ds.XMLID == name {
			return ds, nil
		}
	}
	return to.DeliveryService{}, fmt.Errorf("not found")
}

func createRemapDotConfigEdge(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	deliveryServices, err := toClient.DeliveryServices()
	if err != nil {
		return "", fmt.Errorf("error getting delivery services: %v", err)
	}

	paramMap := createParamsMap(params) // TODO pass param map
	//map[string]map[string]string

	server, err := getServer(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("error getting server: %v", err)
	}

	regexDSes, err := getServerDeliveryServices(toClient, server)
	if err != nil {
		return "", fmt.Errorf("error getting server delivery services: %v", err)
	}
	for dsName, crcDs := range regexDSes {
		ds, err := getDS(dsName, deliveryServices)
		if err != nil {
			// TODO log and continue, instead of failing for one missing ds
			return "", fmt.Errorf("error getting delivery service %v: %v", dsName, err)
		}

		for _, matchset := range crcDs.Matchsets {
			remapLines, err := getRemapLines(server, ds, matchset, paramMap)
			if err != nil {
				return "", fmt.Errorf("error getting remap lines: %v", err)
			}
			for mapFrom, mapTo := range remapLines {
				s += buildRemapLine(toClient, trafficServerHost, ds, matchset, paramMap, mapFrom, mapTo)
			}
		}
	}

	return s, nil
}

func createRemapDotConfig(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	serverTypeStr, err := getServerTypeStr(toClient, trafficServerHost)
	if err != nil {
		return "", fmt.Errorf("error getting server '%v' type: %v", trafficServerHost, err)
	}

	if isMid := strings.HasPrefix(serverTypeStr, "MID"); isMid {
		return createRemapDotConfigMid(toClient, filename, trafficOpsHost, trafficServerHost, params)
	} else {
		return createRemapDotConfigEdge(toClient, filename, trafficOpsHost, trafficServerHost, params)
	}
}
