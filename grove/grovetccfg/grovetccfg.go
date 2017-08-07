package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/crconfig"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/grove"
)

const Version = "0.1"
const UserAgent = "grove-tc-cfg/" + Version
const TrafficOpsTimeout = time.Second * 90

func AvailableStatuses() map[string]struct{} {
	return map[string]struct{}{
		"reported": struct{}{},
		"online":   struct{}{},
	}
}

func main() {
	toURL := flag.String("tourl", "", "The Traffic Ops URL")
	toUser := flag.String("touser", "", "The Traffic Ops username")
	toPass := flag.String("topass", "", "The Traffic Ops password")
	pretty := flag.Bool("pretty", false, "Whether to pretty-print output")
	host := flag.String("host", "", "The hostname of the server whose config to generate")
	api := flag.String("api", "1.2", "API version. Determines whether to use /api/1.3/configs/ or older, less efficient 1.2 APIs")
	toInsecure := flag.Bool("insecure", false, "Whether to allow invalid certificates with Traffic Ops")
	flag.Parse()

	useCache := false
	toc, err := to.LoginWithAgent(*toURL, *toUser, *toPass, *toInsecure, UserAgent, useCache, TrafficOpsTimeout)
	if err != nil {
		fmt.Printf("Error connecting to Traffic Ops: %v\n", err)
		os.Exit(1)
	}

	rules := grove.RemapRules{}
	if *api == "1.3" {
		rules, err = createRulesNewAPI(toc, *host)
	} else {
		rules, err = createRulesOldAPI(toc, *host) // TODO remove once 1.3 / traffic_ops_golang is deployed to production.
	}
	if err != nil {
		fmt.Printf("Error creating rules: %v\n", err)
		os.Exit(1)
	}

	jsonRules := grove.RemapRulesToJSON(rules)
	bts := []byte{}
	if *pretty {
		bts, err = json.MarshalIndent(jsonRules, "", "  ")
	} else {
		bts, err = json.Marshal(jsonRules)
	}

	if err != nil {
		fmt.Printf("Error marshalling rules JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", string(bts))
	os.Exit(0)
}

func createRulesOldAPI(toc *to.Session, host string) (grove.RemapRules, error) {
	cachegroupsArr, err := toc.CacheGroups()
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Cachegroups: %v\n", err)
		os.Exit(1)
	}
	cachegroups := makeCachegroupsNameMap(cachegroupsArr)

	serversArr, err := toc.Servers()
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Servers: %v\n", err)
		os.Exit(1)
	}
	servers := makeServersHostnameMap(serversArr)

	hostServer, ok := servers[host]
	if !ok {
		fmt.Printf("Error: host '%v' not in Servers\n", host)
		os.Exit(1)
	}

	deliveryservices, err := toc.DeliveryServicesByServer(hostServer.ID)
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Deliveryservices: %v\n", err)
		os.Exit(1)
	}
	// deliveryservices := makeDeliveryservicesXMLIDMap(deliveryservicesArr)

	// deliveryserviceServers, err := toc.DeliveryServiceServer("0", "999999")
	// if err != nil {
	// 	fmt.Printf("Error getting Traffic Ops Deliveryservice Servers: %v\n", err)
	// 	os.Exit(1)
	// }

	deliveryserviceRegexArr, err := toc.DeliveryServiceRegexes()
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Deliveryservice Regexes: %v\n", err)
		os.Exit(1)
	}
	deliveryserviceRegexes := makeDeliveryserviceRegexMap(deliveryserviceRegexArr)

	cdnsArr, err := toc.CDNs()
	if err != nil {
		fmt.Printf("Error getting Traffic Ops CDNs: %v\n", err)
		os.Exit(1)
	}
	cdns := makeCDNMap(cdnsArr)

	serverParameters, err := toc.Parameters(hostServer.Profile)
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Parameters for host '%v' profile '%v': %v\n", host, hostServer.Profile, err)
		os.Exit(1)
	}


	// crconfigBts, _, err := toc.GetCRConfig()
	// if err != nil {
	// 	fmt.Printf("Error getting Traffic Ops CRConfig: %v\n", err)
	//	os.Exit(1)
	// }
	// crconfig := crconfig.CRConfig{}
	// if err := json.Unmarshal(crconfigBts, crconfig); err != nil {
	// 	fmt.Printf("Error parsing CRConfig JSON: %v\n", err)
	//	os.Exit(1)
	// }

	parents, err := getParents(host, servers, cachegroups)
	if err != nil {
		fmt.Printf("Error getting '%v' parents: %v\n", err)
		os.Exit(1)
	}

	sameCDN := func(s to.Server) bool {
		return s.CDNName == hostServer.CDNName
	}

	serverAvailable := func(s to.Server) bool {
		status := strings.ToLower(s.Status)
		statuses := AvailableStatuses()
		_, ok := statuses[status]
		return ok
	}

	// serverDses, err := getServerDeliveryservices(*host, servers, deliveryserviceServers, deliveryservices)
	// if err != nil {
	// 	fmt.Printf("Error getting '%v' parents: %v\n", err)
	// 	os.Exit(1)
	// }

	parents = filterParents(parents, sameCDN)
	parents = filterParents(parents, serverAvailable)

	// fmt.Println("Parents:")
	// for _, parent := range parents {
	// 	fmt.Println(parent.HostName)
	// }

	return createRulesOld(host, deliveryservices, parents, deliveryserviceRegexes, cdns, serverParameters)
}

func createRulesNewAPI(toc *to.Session, host string) (grove.RemapRules, error) {
	cacheCfg, err := toc.CacheConfig(host)
	if err != nil {
		fmt.Printf("Error getting Traffic Ops Cache Config: %v\n", err)
		os.Exit(1)
	}

	rules := []grove.RemapRule{}

	allowedIPs, err := makeAllowIP(cacheCfg.AllowIP)
	if err != nil {
		return grove.RemapRules{}, fmt.Errorf("creating allowed IPs: %v", err)
	}

	weight := DefaultRuleWeight
	retryNum := DefaultRetryNum
	timeout := DefaultTimeout
	parentSelection := DefaultRuleParentSelection

	for _, ds := range cacheCfg.DeliveryServices {
		protocol := ds.Protocol
		queryStringRule, err := getQueryStringRule(ds.QueryStringIgnore)
		if err != nil {
			return grove.RemapRules{}, fmt.Errorf("getting deliveryservice %v Query String Rule: %v", ds.XMLID, err)
		}

		protocolStrs := []ProtocolStr{}
		switch protocol {
		case ProtocolHTTP:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
		case ProtocolHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPAndHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPToHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "https"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		}

		dsType := strings.ToLower(ds.Type)
		if !strings.HasPrefix(dsType, "http") && !strings.HasPrefix(dsType, "dns") {
			fmt.Printf("createRules skipping deliveryservice %v - unknown type %v", ds.XMLID, ds.Type)
			continue
		}

		for _, protocolStr := range protocolStrs {
			// regexes, ok := dsRegexes[ds.XMLID]
			// if !ok {
			// 	return grove.RemapRules{}, fmt.Errorf("deliveryservice '%v' has no regexes", ds.XMLID)
			// }

			for _, dsRegex := range ds.Regexes {
				rule := grove.RemapRule{}
				pattern, patternLiteralRegex := trimLiteralRegex(dsRegex)
				rule.Name = fmt.Sprintf("%s.%s.%s.%s", ds.XMLID, protocolStr.From, protocolStr.To, pattern)
				rule.From = buildFrom(protocolStr.From, pattern, patternLiteralRegex, host, dsType, cacheCfg.Domain)
				for _, parent := range cacheCfg.Parents {
					to, proxyURLStr := buildToNew(parent, protocolStr.To, ds.OriginFQDN, dsType)
					proxyURL, err := url.Parse(proxyURLStr)
					if err != nil {
						return grove.RemapRules{}, fmt.Errorf("error parsing deliveryservice %v parent %v proxy_url: %v", ds.XMLID, parent.Host, proxyURLStr)
					}

					ruleTo := grove.RemapRuleTo{
						RemapRuleToBase: grove.RemapRuleToBase{
							URL:      to,
							Weight:   &weight,
							RetryNum: &retryNum,
						},
						ProxyURL:        proxyURL,
						RetryCodes:      DefaultRetryCodes(),
						Timeout:         &timeout,
						ParentSelection: &parentSelection,
					}
					rule.To = append(rule.To, ruleTo)
					// TODO get from TO?
					rule.RetryNum = &retryNum
					rule.Timeout = &timeout
					rule.RetryCodes = DefaultRetryCodes()
					rule.QueryString = queryStringRule
					if err != nil {
						return grove.RemapRules{}, err
					}
					rule.ConnectionClose = DefaultRuleConnectionClose
					rule.ParentSelection = &parentSelection
					rule.Allow = allowedIPs
				}
				rules = append(rules, rule)
			}
		}
	}

	remapRules := grove.RemapRules{
		Rules:           rules,
		RetryCodes:      DefaultRetryCodes(),
		Timeout:         &timeout,
		ParentSelection: &parentSelection,
	}

	return remapRules, nil
}

func makeServersHostnameMap(servers []to.Server) map[string]to.Server {
	m := map[string]to.Server{}
	for _, server := range servers {
		m[server.HostName] = server
	}
	return m
}

func makeCachegroupsNameMap(cgs []to.CacheGroup) map[string]to.CacheGroup {
	m := map[string]to.CacheGroup{}
	for _, cg := range cgs {
		m[cg.Name] = cg
	}
	return m
}

func makeDeliveryservicesXMLIDMap(dses []to.DeliveryService) map[string]to.DeliveryService {
	m := map[string]to.DeliveryService{}
	for _, ds := range dses {
		m[ds.XMLID] = ds
	}
	return m
}

func makeDeliveryservicesIDMap(dses []to.DeliveryService) map[int]to.DeliveryService {
	m := map[int]to.DeliveryService{}
	for _, ds := range dses {
		m[ds.ID] = ds
	}
	return m
}

func makeDeliveryserviceRegexMap(dsrs []to.DeliveryServiceRegexes) map[string][]to.DeliveryServiceRegex {
	m := map[string][]to.DeliveryServiceRegex{}
	for _, dsr := range dsrs {
		m[dsr.DSName] = dsr.Regexes
	}
	return m
}

func makeCDNMap(cdns []to.CDN) map[string]to.CDN {
	m := map[string]to.CDN{}
	for _, cdn := range cdns {
		m[cdn.Name] = cdn
	}
	return m
}

func getServerDeliveryservices(hostname string, servers map[string]to.Server, dssrvs []to.DeliveryServiceServer, dses []to.DeliveryService) ([]to.DeliveryService, error) {
	server, ok := servers[hostname]
	if !ok {
		return nil, fmt.Errorf("server %v not found in Traffic Ops Servers", hostname)
	}
	serverID := server.ID
	dsByID := makeDeliveryservicesIDMap(dses)
	serverDses := []to.DeliveryService{}
	for _, dssrv := range dssrvs {
		if dssrv.Server != serverID {
			continue
		}
		ds, ok := dsByID[dssrv.DeliveryService]
		if !ok {
			return nil, fmt.Errorf("delivery service ID %v not found in Traffic Ops DeliveryServices", dssrv.DeliveryService)
		}
		serverDses = append(serverDses, ds)
	}
	return serverDses, nil
}

func getParents(hostname string, servers map[string]to.Server, cachegroups map[string]to.CacheGroup) ([]to.Server, error) {
	server, ok := servers[hostname]
	if !ok {
		return nil, fmt.Errorf("hostname not found in Servers")
	}

	cachegroup, ok := cachegroups[server.Cachegroup]
	if !ok {
		return nil, fmt.Errorf("server cachegroup '%v' not found in Cachegroups", server.Cachegroup)
	}

	parents := []to.Server{}
	for _, server := range servers {
		if server.Cachegroup == cachegroup.ParentName {
			parents = append(parents, server)
		}
	}
	return parents, nil
}

func filterParents(parents []to.Server, include func(to.Server) bool) []to.Server {
	newParents := []to.Server{}
	for _, parent := range parents {
		if include(parent) {
			newParents = append(newParents, parent)
		}
	}
	return newParents
}

const ProtocolHTTP = 0
const ProtocolHTTPS = 1
const ProtocolHTTPAndHTTPS = 2
const ProtocolHTTPToHTTPS = 3

type ProtocolStr struct {
	From string
	To   string
}

// trimLiteralRegex removes the prefix and suffix in .*\.foo\.* delivery service regexes. Traffic Ops Delivery Services have regexes of this form, which aren't really regexes, and the .*\ and \.* need stripped to construct the "to" FQDN. Returns the trimmed string, and whether it was of the form `.*\.foo\.*`
func trimLiteralRegex(s string) (string, bool) {
	prefix := `.*\.`
	suffix := `\..*`
	if strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix) {
		return s[len(prefix) : len(s)-len(suffix)], true
	}
	return s, false
}

// buildFrom builds the remap "from" URI prefix. It assumes ttype is a delivery service type HTTP or DNS, behavior is undefined for any other ttype.
func buildFrom(protocol string, pattern string, patternLiteralRegex bool, host string, dsType string, cdnDomain string) string {
	if !patternLiteralRegex {
		return protocol + "://" + pattern
	}

	if isHttp := strings.HasPrefix(dsType, "http"); isHttp {
		return protocol + "://" + host + "." + pattern + "." + cdnDomain
	}

	return protocol + "://" + "edge." + pattern + "." + cdnDomain
}

func dsTypeSkipsMid(ttype string) bool {
	ttype = strings.ToLower(ttype)
	if ttype == "http_no_cache" || ttype == "http_live" || ttype == "dns_live" {
		return true
	}
	if strings.Contains(ttype, "live") && !strings.Contains(ttype, "natnl") {
		return true
	}
	return false
}

// buildTo returns the to URL, and the Proxy URL (if any)
func buildTo(parentServer to.Server, protocol string, originURI string, dsType string) (string, string) {
	// TODO add port?
	to := originURI
	proxy := ""
	if !dsTypeSkipsMid(dsType) {
		proxy = "http://" + parentServer.HostName + "." + parentServer.DomainName + ":" + strconv.Itoa(parentServer.TCPPort)
	}
	return to, proxy
}

// buildToNew returns the to URL, and the Proxy URL (if any)
func buildToNew(parent to.CacheConfigParent, protocol string, originURI string, dsType string) (string, string) {
	// TODO add port?
	to := originURI
	proxy := ""
	if !dsTypeSkipsMid(dsType) {
		proxy = "http://" + parent.Host + "." + parent.Domain + ":" + strconv.FormatUint(uint64(parent.Port), 10)
	}
	return to, proxy
}

const DeliveryServiceQueryStringCacheAndRemap = 0
const DeliveryServiceQueryStringNoCacheRemap = 1
const DeliveryServiceQueryStringNoCacheNoRemap = 2

func getQueryStringRule(dsQstringIgnore int) (grove.QueryStringRule, error) {
	switch dsQstringIgnore {
	case DeliveryServiceQueryStringCacheAndRemap:
		return grove.QueryStringRule{Remap: true, Cache: true}, nil
	case DeliveryServiceQueryStringNoCacheRemap:
		return grove.QueryStringRule{Remap: true, Cache: true}, nil
	case DeliveryServiceQueryStringNoCacheNoRemap:
		return grove.QueryStringRule{Remap: false, Cache: false}, nil
	default:
		return grove.QueryStringRule{}, fmt.Errorf("unknown delivery service qstringIgnore value '%v'", dsQstringIgnore)
	}
}

func DefaultRetryCodes() map[int]struct{} {
	return map[int]struct{}{
		404: struct{}{},
		500: struct{}{},
		501: struct{}{},
		503: struct{}{},
	}
}

const DefaultRuleWeight = 1.0
const DefaultRetryNum = 5
const DefaultTimeout = time.Millisecond * 5000
const DefaultRuleConnectionClose = false
const DefaultRuleParentSelection = grove.ParentSelectionTypeConsistentHash

func getAllowIP(params []to.Parameter) ([]*net.IPNet, error) {
	ips := []string{}
	for _, param := range params {
		if (param.Name == "allow_ip" || param.Name == "allow_ip6") && param.ConfigFile == "astats.config" {
			ips = append(ips, strings.Split(param.Value, ",")...)
		}
	}
	return makeAllowIP(ips)
}

func makeAllowIP(ips []string) ([]*net.IPNet, error) {
	cidrs := make([]*net.IPNet, len(ips))
	for i, ip := range ips {
		ip = strings.TrimSpace(ip)
		if !strings.Contains(ip, "/") {
			if strings.Contains(ip, ":") {
				ip += "/128"
			} else {
				ip += "/32"
			}
		}
		_, cidrnet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, fmt.Errorf("error parsing CIDR '%s': %v", ip, err)
		}
		cidrs[i] = cidrnet
	}
	return cidrs, nil
}

func createRulesOld(hostname string, dses []to.DeliveryService, parents []to.Server, dsRegexes map[string][]to.DeliveryServiceRegex, cdns map[string]to.CDN, hostParams []to.Parameter) (grove.RemapRules, error) {
	rules := []grove.RemapRule{}
	allowedIPs, err := getAllowIP(hostParams)
	if err != nil {
		return grove.RemapRules{}, fmt.Errorf("getting allowed IPs: %v", err)
	}

	weight := DefaultRuleWeight
	retryNum := DefaultRetryNum
	timeout := DefaultTimeout
	parentSelection := DefaultRuleParentSelection

	for _, ds := range dses {
		protocol := ds.Protocol
		queryStringRule, err := getQueryStringRule(ds.QStringIgnore)
		if err != nil {
			return grove.RemapRules{}, fmt.Errorf("getting deliveryservice %v Query String Rule: %v", ds.XMLID, err)
		}

		cdn, ok := cdns[ds.CDNName]
		if !ok {
			return grove.RemapRules{}, fmt.Errorf("deliveryservice '%v' CDN '%v' not found", ds.XMLID, ds.CDNName)
		}

		protocolStrs := []ProtocolStr{}
		switch protocol {
		case ProtocolHTTP:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
		case ProtocolHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPAndHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPToHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "https"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		}

		dsType := strings.ToLower(ds.Type)
		if !strings.HasPrefix(dsType, "http") && !strings.HasPrefix(dsType, "dns") {
			fmt.Printf("createRules skipping deliveryservice %v - unknown type %v", ds.XMLID, ds.Type)
			continue
		}

		for _, protocolStr := range protocolStrs {
			regexes, ok := dsRegexes[ds.XMLID]
			if !ok {
				return grove.RemapRules{}, fmt.Errorf("deliveryservice '%v' has no regexes", ds.XMLID)
			}

			for _, dsRegex := range regexes {
				rule := grove.RemapRule{}
				pattern, patternLiteralRegex := trimLiteralRegex(dsRegex.Pattern)
				rule.Name = fmt.Sprintf("%s.%s.%s.%s", ds.XMLID, protocolStr.From, protocolStr.To, pattern)
				rule.From = buildFrom(protocolStr.From, pattern, patternLiteralRegex, hostname, dsType, cdn.DomainName)
				for _, parent := range parents {
					to, proxyURLStr := buildTo(parent, protocolStr.To, ds.OrgServerFQDN, dsType)
					proxyURL, err := url.Parse(proxyURLStr)
					if err != nil {
						return grove.RemapRules{}, fmt.Errorf("error parsing deliveryservice %v parent %v proxy_url: %v", ds.XMLID, parent.HostName, proxyURLStr)
					}

					ruleTo := grove.RemapRuleTo{
						RemapRuleToBase: grove.RemapRuleToBase{
							URL:      to,
							Weight:   &weight,
							RetryNum: &retryNum,
						},
						ProxyURL:        proxyURL,
						RetryCodes:      DefaultRetryCodes(),
						Timeout:         &timeout,
						ParentSelection: &parentSelection,
					}
					rule.To = append(rule.To, ruleTo)
					// TODO get from TO?
					rule.RetryNum = &retryNum
					rule.Timeout = &timeout
					rule.RetryCodes = DefaultRetryCodes()
					rule.QueryString = queryStringRule
					if err != nil {
						return grove.RemapRules{}, err
					}
					rule.ConnectionClose = DefaultRuleConnectionClose
					rule.ParentSelection = &parentSelection
					rule.Allow = allowedIPs
				}
				rules = append(rules, rule)
			}
		}
	}

	remapRules := grove.RemapRules{
		Rules:           rules,
		RetryCodes:      DefaultRetryCodes(),
		Timeout:         &timeout,
		ParentSelection: &parentSelection,
	}

	return remapRules, nil
}

// if ( $remap→{type} eq "HTTP_NO_CACHE" || $remap→{type} eq "HTTP_LIVE" || $remap→{type} eq "DNS_LIVE" ) {
// $text .= "dest_domain=" . $org_uri→host . " port=" . $org_uri→port . " go_direct=true\n";

/*

http edge to



http edge to:

dns edge to:



edge got:

mid got:

remap.config contains:

parent.config contains
*/

// Name            string          `json:"name"`
// From            string          `json:"from"`
// ConnectionClose bool            `json:"connection-close"`
// QueryString     QueryStringRule `json:"query-string"`
// // ConcurrentRuleRequests is the number of concurrent requests permitted to a remap rule, that is, to an origin. If this is 0, the global config is used.
// ConcurrentRuleRequests int  `json:"concurrent_rule_requests"`
// RetryNum               *int `json:"retry_num"`
// Timeout         *time.Duration
// ParentSelection *ParentSelectionType
// To              []RemapRuleTo
// Allow           []*net.IPNet
// Deny            []*net.IPNet
// RetryCodes      map[int]struct{}
// ConsistentHash  ATSConsistentHash
