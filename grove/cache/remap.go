package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/chash"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// CacheType is the type (or tier) of a CDN cache.
type ParentSelectionType string

const (
	ParentSelectionTypeConsistentHash = ParentSelectionType("consistent-hash")
	ParentSelectionTypeRoundRobin     = ParentSelectionType("round-robin")
	ParentSelectionTypeInvalid        = ParentSelectionType("")
)

func (t ParentSelectionType) String() string {
	switch t {
	case ParentSelectionTypeConsistentHash:
		return "consistent-hash"
	case ParentSelectionTypeRoundRobin:
		return "round-robin"
	default:
		return "invalid"
	}
}

func ParentSelectionTypeFromString(s string) ParentSelectionType {
	s = strings.ToLower(s)
	if s == "consistent-hash" {
		return ParentSelectionTypeConsistentHash
	}
	if s == "round-robin" {
		return ParentSelectionTypeRoundRobin
	}
	return ParentSelectionTypeInvalid
}

type HTTPRequestRemapper interface {
	// Remap returns the remapped request, the matched rule name, whether the requestor's IP is allowed, whether to connection close, whether a match was found, and any error.
	// Remap(r *http.Request, scheme string, failures int) Remapping
	Rules() []RemapRule
	RemappingProducer(r *http.Request, scheme string) (*RemappingProducer, error)
}

type simpleHttpRequestRemapper struct {
	remapper Remapper
}

func (hr simpleHttpRequestRemapper) Rules() []RemapRule {
	return hr.remapper.Rules()
}

// getFQDN returns the FQDN. It tries to get the FQDN from a Remap Rule. Remap Rules should always begin with the scheme, e.g. `http://`. If the given rule does not begin with a valid scheme, behavior is undefined.
// TODO test
func getFQDN(rule string) string {
	schemeStr := "://"
	schemePos := strings.Index(rule, schemeStr)
	if schemePos == -1 {
		return rule // invalid rule, doesn't start with a scheme
	}
	schemePos += len(schemeStr)
	rule = rule[schemePos:]

	slashPos := strings.Index(rule, "/")
	if slashPos == -1 {
		return rule // rule is just the scheme+FQDN, perfectly normal
	}
	rule = rule[:slashPos] // strip off the path
	return rule
}

type Remapping struct {
	Request         *http.Request
	ProxyURL        *url.URL
	Name            string
	CacheKey        string
	ConnectionClose bool
	Timeout         time.Duration
	RetryNum        int
	RetryCodes      map[int]struct{}
}

// TODO rename
type RemappingProducer struct {
	oldURI   string
	rule     RemapRule
	cacheKey string
	failures int
}

func (r *RemappingProducer) CacheKey() string {
	return r.cacheKey
}

func (r *RemappingProducer) ConnectionClose() bool {
	return r.rule.ConnectionClose
}

func (r *RemappingProducer) Name() string {
	return r.rule.Name
}

func (r *RemappingProducer) ToFQDN() string {
	// TODO verify To is not allowed to be constructed with < 1 element
	return strings.TrimPrefix(strings.TrimPrefix(r.rule.To[0].URL, "http://"), "https://")
}

func (r *RemappingProducer) ProxyStr() string {
	if r.rule.To[0].ProxyURL != nil && r.rule.To[0].ProxyURL.Host != "" {
		return r.rule.To[0].ProxyURL.Host
	} else {
		return "NONE" // TODO const?
	}
}

var ErrRuleNotFound = errors.New("remap rule not found")
var ErrIPNotAllowed = errors.New("IP not allowed")
var ErrNoMoreRetries = errors.New("retry num exceeded")

// RequestURI returns the URI of the given request. This must be used, because Go does not populate the scheme of requests that come in from clients.
func RequestURI(r *http.Request, scheme string) string {
	return scheme + "://" + r.Host + r.RequestURI
}

func (hr simpleHttpRequestRemapper) RemappingProducer(r *http.Request, scheme string) (*RemappingProducer, error) {
	uri := RequestURI(r, scheme)
	rule, ok := hr.remapper.Remap(uri)
	if !ok {
		return nil, ErrRuleNotFound
	}

	if ip, err := GetIP(r); err != nil {
		return nil, fmt.Errorf("parsing client IP: %v", err)
	} else if !rule.Allowed(ip) {
		return nil, ErrIPNotAllowed
	} else {
		log.Debugf("Allowed %v\n", ip)
	}

	cacheKey := rule.CacheKey(r.Method, uri)

	return &RemappingProducer{
		rule:     rule,
		oldURI:   uri,
		cacheKey: cacheKey,
	}, nil
}

// GetNext returns the remapping to use to request, whether retries are allowed (i.e. if this is the last retry), or any error
func (p *RemappingProducer) GetNext(r *http.Request) (Remapping, bool, error) {
	if *p.rule.RetryNum < p.failures {
		return Remapping{}, false, ErrNoMoreRetries
	}

	newUri, proxyURL := p.rule.URI(p.oldURI, r.URL.Path, r.URL.RawQuery, p.failures)
	p.failures++
	newReq, err := http.NewRequest(r.Method, newUri, nil)
	if err != nil {
		return Remapping{}, false, fmt.Errorf("creating new request: %v\n", err)
	}
	web.CopyHeaderTo(r.Header, &newReq.Header)

	log.Debugf("GetNext oldUri: %v, Host: %v\n", p.oldURI, newReq.Header.Get("Host"))
	log.Debugf("GetNext newUri: %v, fqdn: %v\n", newUri, getFQDN(newUri))
	log.Debugf("GetNext rule name: %v\n", p.rule.Name)

	newReq.Header.Set("Host", getFQDN(newUri))

	retryAllowed := *p.rule.RetryNum < p.failures
	return Remapping{
		Request:         newReq,
		ProxyURL:        proxyURL,
		Name:            p.rule.Name,
		CacheKey:        p.cacheKey,
		ConnectionClose: p.rule.ConnectionClose,
		Timeout:         *p.rule.Timeout,
		RetryNum:        *p.rule.RetryNum,
		RetryCodes:      p.rule.RetryCodes,
	}, retryAllowed, nil
}

func RemapperToHTTP(r Remapper) HTTPRequestRemapper {
	return simpleHttpRequestRemapper{remapper: r}
}

func NewHTTPRequestRemapper(remap []RemapRule) HTTPRequestRemapper {
	return RemapperToHTTP(NewLiteralPrefixRemapper(remap))
}

// Remapper provides a function which takes strings and maps them to other strings. This is designed for URL prefix remapping, for a reverse proxy.
type Remapper interface {
	// Remap returns the given string remapped, the unique name of the rule found, and whether a remap rule was found
	Remap(uri string) (RemapRule, bool)
	// Rules returns the unique names of every remap rule.
	Rules() []RemapRule
}

// TODO change to use a prefix tree, for speed
type literalPrefixRemapper struct {
	remap []RemapRule
}

// Remap returns the remapped string, the remap rule name, the remap rule's options, and whether a remap was found
func (r literalPrefixRemapper) Remap(s string) (RemapRule, bool) {
	for _, rule := range r.remap {
		if strings.HasPrefix(s, rule.From) {
			return rule, true
		}
	}
	return RemapRule{}, false
}

func (r literalPrefixRemapper) Rules() []RemapRule {
	rules := make([]RemapRule, len(r.remap))
	for _, rule := range r.remap {
		rules = append(rules, rule)
	}
	return rules
}

func NewLiteralPrefixRemapper(remap []RemapRule) Remapper {
	return literalPrefixRemapper{remap: remap}
}

type RemapRulesBase struct {
	RetryNum *int `json:"retry_num"`
}

type RemapRulesJSON struct {
	RemapRulesBase
	Rules           []RemapRuleJSON `json:"rules"`
	RetryCodes      *[]int          `json:"retry_codes"`
	TimeoutMS       *int            `json:"timeout_ms"`
	ParentSelection *string         `json:"parent_selection"`
}

type RemapRules struct {
	RemapRulesBase
	Rules           []RemapRule
	RetryCodes      map[int]struct{}
	Timeout         *time.Duration
	ParentSelection *ParentSelectionType
}

type RemapRuleToBase struct {
	URL      string   `json:"url"`
	Weight   *float64 `json:"weight"`
	RetryNum *int     `json:"retry_num"`
}

type RemapRuleToJSON struct {
	RemapRuleToBase
	ProxyURL   *string `json:"proxy_url"`
	TimeoutMS  *int    `json:"timeout_ms"`
	RetryCodes *[]int  `json:"retry_codes"`
}

type RemapRuleTo struct {
	RemapRuleToBase
	ProxyURL   *url.URL
	Timeout    *time.Duration
	RetryCodes map[int]struct{}
}

type RemapRuleBase struct {
	Name               string          `json:"name"`
	From               string          `json:"from"`
	CertificateFile    string          `json:"certificate-file"`
	CertificateKeyFile string          `json:"certificate-key-file"`
	ConnectionClose    bool            `json:"connection-close"`
	QueryString        QueryStringRule `json:"query-string"`
	// ConcurrentRuleRequests is the number of concurrent requests permitted to a remap rule, that is, to an origin. If this is 0, the global config is used.
	ConcurrentRuleRequests int  `json:"concurrent_rule_requests"`
	RetryNum               *int `json:"retry_num"`
}

type RemapRuleJSON struct {
	RemapRuleBase
	TimeoutMS       *int              `json:"timeout_ms"`
	ParentSelection *string           `json:"parent_selection"`
	To              []RemapRuleToJSON `json:"to"`
	Allow           []string          `json:"allow"`
	Deny            []string          `json:"deny"`
	RetryCodes      *[]int            `json:"retry_codes"`
}

type RemapRule struct {
	RemapRuleBase
	Timeout         *time.Duration
	ParentSelection *ParentSelectionType
	To              []RemapRuleTo
	Allow           []*net.IPNet
	Deny            []*net.IPNet
	RetryCodes      map[int]struct{}
	ConsistentHash  chash.ATSConsistentHash
}

func GetIP(r *http.Request) (net.IP, error) {
	clientIpStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("malformed client address '%s'", r.RemoteAddr)
	}
	clientIP := net.ParseIP(clientIpStr)
	if clientIP == nil {
		return nil, fmt.Errorf("malformed client IP address '%s'", clientIpStr)
	}
	return clientIP, nil
}

func (r *RemapRule) Allowed(ip net.IP) bool {
	for _, network := range r.Deny {
		if network.Contains(ip) {
			log.Debugf("deny contains ip\n")
			return false
		}
	}
	if len(r.Allow) == 0 {
		log.Debugf("Allowed len 0\n")
		return true
	}
	for _, network := range r.Allow {
		if network.Contains(ip) {
			log.Debugf("allow contains ip\n")
			return true
		}
	}
	return false
}

type QueryStringRule struct {
	Remap bool `json:"remap"`
	Cache bool `json:"cache"`
}

// URI takes a request URI and maps it to the real URI to proxy-and-cache. The `failures` parameter indicates how many parents have tried and failed, indicating to skip to the nth hashed parent. Returns the URI to request, and the proxy URL (if any)
func (r RemapRule) URI(fromURI string, path string, query string, failures int) (string, *url.URL) {
	fromHash := path
	if r.QueryString.Remap && query != "" {
		fromHash += "?" + query
	}

	// fmt.Println("RemapRule.URI fromURI " + fromHash)
	to, proxyURI := r.uriGetTo(fromHash, failures)
	uri := to + fromURI[len(r.From):]
	if !r.QueryString.Remap {
		if i := strings.Index(uri, "?"); i != -1 {
			uri = uri[:i]
		}
	}
	return uri, proxyURI
}

// uriGetTo is a helper func for URI. It returns the To URL, based on the Parent Selection type. In the event of failure, it logs the error and returns the first parent. Also returns the URL's Proxy URI (if any).
func (r RemapRule) uriGetTo(fromURI string, failures int) (string, *url.URL) {
	switch *r.ParentSelection {
	case ParentSelectionTypeConsistentHash:
		return r.uriGetToConsistentHash(fromURI, failures)
	default:
		log.Errorf("RemapRule.URI: Rule '%v': Unknown Parent Selection type %v - using first URI in rule\n", r.Name, r.ParentSelection)
		return r.To[0].URL, r.To[0].ProxyURL
	}
}

// uriGetToConsistentHash is a helper func for URI, uriGetTo. It returns the To URL using Consistent Hashing. In the event of failure, it logs the error and returns the first parent. Also returns the Proxy URI (if any).
func (r RemapRule) uriGetToConsistentHash(fromURI string, failures int) (string, *url.URL) {
	// fmt.Printf("DEBUGL uriGetToConsistentHash RemapRule %+v\n", r)
	if r.ConsistentHash == nil {
		log.Errorf("RemapRule.URI: Rule '%v': Parent Selection Type ConsistentHash, but rule.ConsistentHash is nil! Using first parent\n", r.Name)
		return r.To[0].URL, r.To[0].ProxyURL
	}

	// fmt.Printf("DEBUGL uriGetToConsistentHash\n")
	iter, _, err := r.ConsistentHash.Lookup(fromURI)
	if err != nil {
		// if r.ConsistentHash.First() == nil {
		// 	fmt.Printf("DEBUGL uriGetToConsistentHash NodeMap empty!\n")
		// }
		// fmt.Printf("DEBUGL uriGetToConsistentHash fromURI '%v' err %v returning '%v'\n", fromURI, err, r.To[0].URL)
		log.Errorf("RemapRule.URI: Rule '%v': Error looking up Consistent Hash! Using first parent\n", r.Name)
		return r.To[0].URL, r.To[0].ProxyURL
	}

	for i := 0; i < failures; i++ {
		iter = iter.NextWrap()
	}

	return iter.Val().Name, iter.Val().ProxyURL
}

func (r RemapRule) CacheKey(method string, fromURI string) string {
	// TODO don't cache on `to`, since it's affected by Parent Selection
	// TODO add parent selection
	to := r.To[0].URL
	uri := to + fromURI[len(r.From):]
	if !r.QueryString.Cache {
		if i := strings.Index(uri, "?"); i != -1 {
			uri = uri[:i]
		}
	}
	key := method + ":" + uri
	return key
}

func LoadRemapRules(path string) ([]RemapRule, error) {
	fmt.Printf("Loading Remap Rules\n")
	defer func() {
		fmt.Printf("Loaded Remap Rules\n")
	}()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	remapRulesJSON := RemapRulesJSON{}
	if err := json.NewDecoder(file).Decode(&remapRulesJSON); err != nil {
		return nil, fmt.Errorf("decoding JSON: %s", err)
	}

	remapRules := RemapRules{RemapRulesBase: remapRulesJSON.RemapRulesBase}
	if remapRulesJSON.RetryCodes != nil {
		remapRules.RetryCodes = make(map[int]struct{}, len(*remapRulesJSON.RetryCodes))
		for _, code := range *remapRulesJSON.RetryCodes {
			if _, ok := ValidHttpCodes[code]; !ok {
				return nil, fmt.Errorf("error parsing rules: retry code invalid: %v", code)
			}
			remapRules.RetryCodes[code] = struct{}{}
		}
	}
	if remapRulesJSON.TimeoutMS != nil {
		t := time.Duration(*remapRulesJSON.TimeoutMS) * time.Millisecond
		if remapRules.Timeout = &t; *remapRules.Timeout < 0 {
			return nil, fmt.Errorf("error parsing rules: timeout must be positive: %v", remapRules.Timeout)
		}
	}
	if remapRulesJSON.ParentSelection != nil {
		ps := ParentSelectionTypeFromString(*remapRulesJSON.ParentSelection)
		if remapRules.ParentSelection = &ps; *remapRules.ParentSelection == ParentSelectionTypeInvalid {
			return nil, fmt.Errorf("error parsing rules: parent selection invalid: '%v'", remapRulesJSON.ParentSelection)
		}
	}

	rules := make([]RemapRule, len(remapRulesJSON.Rules))
	for i, jsonRule := range remapRulesJSON.Rules {
		fmt.Printf("Creating Remap Rule %v\n", jsonRule.Name)
		rule := RemapRule{RemapRuleBase: jsonRule.RemapRuleBase}

		if jsonRule.RetryCodes != nil {
			rule.RetryCodes = make(map[int]struct{}, len(*jsonRule.RetryCodes))
			for _, code := range *jsonRule.RetryCodes {
				if _, ok := ValidHttpCodes[code]; !ok {
					return nil, fmt.Errorf("error parsing rule %v retry code invalid: %v", rule.Name, code)
				}
				rule.RetryCodes[code] = struct{}{}
			}
		} else {
			rule.RetryCodes = remapRules.RetryCodes
		}

		if jsonRule.TimeoutMS != nil {
			t := time.Duration(*jsonRule.TimeoutMS) * time.Millisecond
			if rule.Timeout = &t; *rule.Timeout < 0 {
				return nil, fmt.Errorf("error parsing rule %v timeout must be positive: %v", rule.Name, rule.Timeout)
			}
		} else {
			rule.Timeout = remapRules.Timeout
		}

		if rule.RetryNum == nil {
			rule.RetryNum = remapRules.RetryNum
		}

		if rule.Allow, err = makeIPNets(jsonRule.Allow); err != nil {
			return nil, fmt.Errorf("error parsing rule %v allows: %v", rule.Name, err)
		}
		if rule.Deny, err = makeIPNets(jsonRule.Deny); err != nil {
			return nil, fmt.Errorf("error parsing rule %v denys: %v", rule.Name, err)
		}
		if rule.To, err = makeTo(jsonRule.To, rule); err != nil {
			return nil, fmt.Errorf("error parsing rule %v to: %v", rule.Name, err)
		}
		if jsonRule.ParentSelection != nil {
			ps := ParentSelectionTypeFromString(*jsonRule.ParentSelection)
			if rule.ParentSelection = &ps; *rule.ParentSelection == ParentSelectionTypeInvalid {
				return nil, fmt.Errorf("error parsing rule %v parent selection invalid: '%v'", rule.Name, jsonRule.ParentSelection)
			}
		} else {
			rule.ParentSelection = remapRules.ParentSelection
		}

		if rule.ParentSelection == nil {
			return nil, fmt.Errorf("error parsing rule %v - no parent_selection - must be set at rules or rule level", rule.Name)
		}

		if len(rule.To) == 0 {
			return nil, fmt.Errorf("error parsing rule %v - no to - must have at least one parent", rule.Name)
		}

		if *rule.ParentSelection == ParentSelectionTypeConsistentHash {
			rule.ConsistentHash = makeRuleHash(rule)
		} else {
		}
		rules[i] = rule
	}

	return rules, nil
}

const DefaultReplicas = 1024

func makeRuleHash(rule RemapRule) chash.ATSConsistentHash {
	h := chash.NewSimpleATSConsistentHash(DefaultReplicas)
	for _, to := range rule.To {
		h.Insert(&chash.ATSConsistentHashNode{Name: to.URL, ProxyURL: to.ProxyURL}, *to.Weight)
	}
	if h.First() == nil {
		fmt.Printf("DEBUGLL makeRuleHash %v NodeMap empty!\n", rule.Name)
	}

	// fmt.Println("makeRuleHash " + rule.Name + ":\n" + h.String())

	return h
}

func makeTo(tosJSON []RemapRuleToJSON, rule RemapRule) ([]RemapRuleTo, error) {
	tos := make([]RemapRuleTo, len(tosJSON))
	for i, toJSON := range tosJSON {
		if toJSON.Weight == nil {
			w := 1.0
			toJSON.Weight = &w
		}
		to := RemapRuleTo{RemapRuleToBase: toJSON.RemapRuleToBase}
		if toJSON.ProxyURL != nil {
			proxyURL, err := url.Parse(*toJSON.ProxyURL)
			if err != nil {
				return nil, fmt.Errorf("error parsing to %v proxy_url: %v", to.URL, toJSON.ProxyURL)
			}
			to.ProxyURL = proxyURL
		}
		if toJSON.TimeoutMS != nil {
			t := time.Duration(*toJSON.TimeoutMS) * time.Millisecond
			if to.Timeout = &t; *to.Timeout < 0 {
				return nil, fmt.Errorf("error parsing to %v timeout must be positive: %v", to.URL, to.Timeout)
			}
		} else {
			to.Timeout = rule.Timeout
		}
		if toJSON.RetryCodes != nil {
			to.RetryCodes = make(map[int]struct{}, len(*toJSON.RetryCodes))
			for _, code := range *toJSON.RetryCodes {
				if _, ok := ValidHttpCodes[code]; !ok {
					return nil, fmt.Errorf("error parsing to %v retry code invalid: %v", to.URL, code)
				}
				to.RetryCodes[code] = struct{}{}
			}
		} else {
			to.RetryCodes = rule.RetryCodes
		}
		if to.RetryNum == nil {
			to.RetryNum = rule.RetryNum
		}
		if to.RetryNum == nil {
			return nil, fmt.Errorf("error parsing to %v - no retry_num - must be set at rules, rule, or to level", to.URL)
		} else if to.Timeout == nil {
			return nil, fmt.Errorf("error parsing to %v - no timeout_ms - must be set at rules, rule, or to level", to.URL)
		} else if to.RetryCodes == nil {
			return nil, fmt.Errorf("error parsing to %v - no retry_codes - must be set at rules, rule, or to level", to.URL)
		}
		tos[i] = to
	}
	return tos, nil
}

func makeIPNets(netStrs []string) ([]*net.IPNet, error) {
	nets := make([]*net.IPNet, 0, len(netStrs))
	for _, netStr := range netStrs {
		netStr = strings.TrimSpace(netStr)
		if netStr == "" {
			continue
		}
		net, err := makeIPNet(netStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing CIDR %v: %v", netStr, err)
		}
		nets = append(nets, net)
	}
	return nets, nil
}

func makeIPNet(cidr string) (*net.IPNet, error) {
	_, cidrnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("error parsing CIDR '%s': %v", cidr, err)
	}
	return cidrnet, nil
}

func LoadRemapper(path string) (HTTPRequestRemapper, error) {
	rules, err := LoadRemapRules(path)
	if err != nil {
		return nil, err
	}
	return NewHTTPRequestRemapper(rules), nil
}

func RemapRulesToJSON(r RemapRules) RemapRulesJSON {
	j := RemapRulesJSON{RemapRulesBase: r.RemapRulesBase}
	if r.Timeout != nil {
		i := int(0)
		j.TimeoutMS = &i
		*j.TimeoutMS = int(*r.Timeout / time.Millisecond)
	}
	if len(r.RetryCodes) > 0 {
		rcs := []int{}
		j.RetryCodes = &rcs
		for code, _ := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, code)
		}
	}
	if r.ParentSelection != nil {
		s := ""
		j.ParentSelection = &s
		*j.ParentSelection = string(*r.ParentSelection)
	}

	for _, rule := range r.Rules {
		j.Rules = append(j.Rules, buildRemapRuleToJSON(rule))
	}
	return j
}

func buildRemapRuleToJSON(r RemapRule) RemapRuleJSON {
	j := RemapRuleJSON{RemapRuleBase: r.RemapRuleBase}
	if r.Timeout != nil {
		t := int(0)
		j.TimeoutMS = &t
		*j.TimeoutMS = int(*r.Timeout / time.Millisecond)
	}
	if r.ParentSelection != nil {
		ps := ""
		j.ParentSelection = &ps
		*j.ParentSelection = string(*r.ParentSelection)
	}
	for _, to := range r.To {
		j.To = append(j.To, RemapRuleToToJSON(to))
	}
	for _, deny := range r.Deny {
		j.Deny = append(j.Deny, deny.String())
	}
	for _, allow := range r.Allow {
		j.Allow = append(j.Allow, allow.String())
	}
	if r.RetryCodes != nil {
		rc := []int{}
		j.RetryCodes = &rc
		for retryCode, _ := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, retryCode)
		}
	}
	return j
}

func RemapRuleToToJSON(r RemapRuleTo) RemapRuleToJSON {
	j := RemapRuleToJSON{RemapRuleToBase: r.RemapRuleToBase}
	if r.ProxyURL != nil {
		s := ""
		j.ProxyURL = &s
		*j.ProxyURL = r.ProxyURL.String()
	}
	if r.Timeout != nil {
		t := int(0)
		j.TimeoutMS = &t
		*j.TimeoutMS = int(*r.Timeout / time.Millisecond)
	}
	if r.RetryCodes != nil {
		rc := []int{}
		j.RetryCodes = &rc
		for retryCode, _ := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, retryCode)
		}
	}
	return j
}
