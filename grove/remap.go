package grove

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
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
	newUri := p.rule.URI(p.oldURI, p.failures)
	p.failures++
	newReq, err := http.NewRequest(r.Method, newUri, nil)
	if err != nil {
		return Remapping{}, false, fmt.Errorf("creating new request: %v\n", err)
	}
	copyHeader(r.Header, &newReq.Header)

	log.Errorf("DEBUGQ oldUri: %v, Host: %v\n", p.oldURI, newReq.Header.Get("Host"))
	log.Errorf("DEBUGQ newUri: %v, fqdn: %v\n", newUri, getFQDN(newUri))
	log.Errorf("DEBUGQ rule name: %v\n", p.rule.Name)

	// if newReq.Header.Get("Host") == p.oldURI {
	// 	log.Errorf("DEBUGQ setting Host header\n")
	newReq.Header.Set("Host", getFQDN(newUri))
	// } else {
	// 	fmt.Printf("DEBUGL leaving Host header: %v\n", newReq.Header.Get("Host"))
	// }

	retryAllowed := p.failures >= *p.rule.RetryNum
	return Remapping{
		Request:         newReq,
		Name:            p.rule.Name,
		CacheKey:        p.cacheKey,
		ConnectionClose: p.rule.ConnectionClose,
		Timeout:         *p.rule.Timeout,
		RetryNum:        *p.rule.RetryNum,
		RetryCodes:      p.rule.RetryCodes,
	}, retryAllowed, nil
}

// Remap returns the given request with its URI remapped, the name of the remap rule found, the cache key, whether the requestor's IP is allowed, whether the rule calls for sending a connection close header, whether a rule was found, and any error. The `failures` parameter indicates the number of failed parents which have occurred, so those parents won't be remapped.
// func (hr simpleHttpRequestRemapper) Remap(r *http.Request, scheme string, failures int) Remapping {
// 	// NewRequest(method, urlStr string, body io.Reader)
// 	// TODO config whether to consider query string, method, headers

// 	oldUri := scheme + "://" + r.Host + r.RequestURI
// 	log.Debugf("Remap oldUri: '%v'\n", oldUri)
// 	log.Debugf("request: '%+v'\n", r)
// 	rule, ok := hr.remapper.Remap(oldUri)
// 	if !ok {
// 		log.Debugf("Remap oldUri: '%v' NOT FOUND\n", oldUri)
// 		return Remapping{RuleNotFound: true}
// 	}

// 	ip, err := GetIP(r)
// 	if err != nil {
// 		return Remapping{Err: fmt.Errorf("parsing client IP: %v", err)}
// 	}

// 	if !rule.Allowed(ip) {
// 		return Remapping{IPNotAllowed: true}
// 	}

// 	log.Debugf("Allowed %v\n", ip)

// 	newUri := rule.URI(oldUri, failures)
// 	cacheKey := rule.CacheKey(r.Method, oldUri)
// 	log.Debugf("Remap newURI: '%v'\nDEBUG Remap cacheKey '%v'\n", newUri, cacheKey)

// 	newReq, err := http.NewRequest(r.Method, newUri, nil) // TODO modify given req in-place?
// 	if err != nil {
// 		log.Errorf("Remap NewRequest: %v\n", err)
// 		return Remapping{RuleNotFound: true} // TODO return err?
// 	}
// 	copyHeader(r.Header, &newReq.Header)

// 	log.Errorf("DEBUGQ oldUri: %v, Host: %v\n", oldUri, newReq.Header.Get("Host"))
// 	log.Errorf("DEBUGQ newUri: %v, fqdn: %v\n", newUri, getFQDN(newUri))

// 	if newReq.Header.Get("Host") == oldUri {
// 		log.Errorf("DEBUGQ setting Host header\n")
// 		newReq.Header.Set("Host", getFQDN(newUri))
// 	}
// 	// newReq.Header.Add("If-Modified-Since", cacheObj.respTime.Format(time.RFC1123))
// 	return Remapping{
// 		Request:         newReq,
// 		Name:            rule.Name,
// 		CacheKey:        cacheKey,
// 		ConnectionClose: rule.ConnectionClose,
// 		Timeout:         *rule.Timeout,
// 		RetryNum:        *rule.RetryNum,
// 		RetryCodes:      rule.RetryCodes,
// 	}
// }

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
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
	Rules    []RemapRuleJSON `json:"rules"`
	RetryNum *int            `json:"retry_num"`
}

type RemapRulesJSON struct {
	RemapRulesBase
	RetryCodes      *[]int  `json:"retry_codes"`
	TimeoutMS       *int    `json:"timeout_ms"`
	ParentSelection *string `json:"parent_selection"`
}

type RemapRules struct {
	RemapRulesBase
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
	ParentSelection *string `json:"parent_selection"`
	TimeoutMS       *int    `json:"timeout_ms"`
	RetryCodes      *[]int  `json:"retry_codes"`
}

type RemapRuleTo struct {
	RemapRuleToBase
	ParentSelection *ParentSelectionType
	Timeout         *time.Duration
	RetryCodes      map[int]struct{}
}

type RemapRuleBase struct {
	Name            string          `json:"name"`
	From            string          `json:"from"`
	ConnectionClose bool            `json:"connection-close"`
	QueryString     QueryStringRule `json:"query-string"`
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
	ConsistentHash  ATSConsistentHash
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

// URI takes a request URI and maps it to the real URI to proxy-and-cache. The `failures` parameter indicates how many parents have tried and failed, indicating to skip to the nth hashed parent
func (r RemapRule) URI(fromURI string, failures int) string {
	to := r.uriGetTo(fromURI, failures)
	uri := to + fromURI[len(r.From):]
	if !r.QueryString.Remap {
		if i := strings.Index(uri, "?"); i != -1 {
			uri = uri[:i]
		}
	}
	return uri
}

// uriGetTo is a helper func for URI. It returns the To URL, based on the Parent Selection type. In the event of failure, it logs the error and returns the first parent.
func (r RemapRule) uriGetTo(fromURI string, failures int) string {
	switch *r.ParentSelection {
	case ParentSelectionTypeConsistentHash:
		return r.uriGetToConsistentHash(fromURI, failures)
	default:
		log.Errorf("RemapRule.URI: Rule '%v': Unknown Parent Selection type %v - using first URI in rule\n", r.Name, r.ParentSelection)
		return r.To[0].URL
	}
}

// uriGetToConsistentHash is a helper func for URI, uriGetTo. It returns the To URL using Consistent Hashing. In the event of failure, it logs the error and returns the first parent.
func (r RemapRule) uriGetToConsistentHash(fromURI string, failures int) string {
	fmt.Printf("DEBUGL uriGetToConsistentHash RemapRule %+v\n", r)
	if r.ConsistentHash == nil {
		log.Errorf("RemapRule.URI: Rule '%v': Parent Selection Type ConsistentHash, but rule.ConsistentHash is nil! Using first parent\n", r.Name)
		return r.To[0].URL
	}

	fmt.Printf("DEBUGL uriGetToConsistentHash\n")
	iter, _, err := r.ConsistentHash.Lookup(fromURI)
	if err != nil {
		if r.ConsistentHash.First() == nil {
			fmt.Printf("DEBUGL uriGetToConsistentHash NodeMap empty!\n")
		}
		fmt.Printf("DEBUGL uriGetToConsistentHash fromURI '%v' err %v returning '%v'\n", fromURI, err, r.To[0].URL)
		log.Errorf("RemapRule.URI: Rule '%v': Error looking up Consistent Hash! Using first parent\n", r.Name)
		return r.To[0].URL
	}

	for i := 0; i < failures; i++ {
		iter = iter.NextWrap()
	}
	fmt.Printf("DEBUGL uriGetToConsistentHash returning iter.Val().Name %v\n", iter.Val().Name)
	return iter.Val().Name
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

	rules := make([]RemapRule, len(remapRules.Rules))
	for i, jsonRule := range remapRules.Rules {
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
			fmt.Printf("DEBUGLL making rule hash %v\n", rule.Name)
			rule.ConsistentHash = makeRuleHash(rule)
		} else {
			fmt.Printf("DEBUGLL NOT making rule hash %v\n", rule.Name)
		}
		rules[i] = rule
	}

	return rules, nil
}

func makeRuleHash(rule RemapRule) ATSConsistentHash {
	replicas := 100 // TODO put in config?
	h := NewSimpleATSConsistentHash(replicas)
	fmt.Printf("DEBUGLL makeRuleHash %v len(rule.To) %v\n", rule.Name, len(rule.To))
	for _, to := range rule.To {
		fmt.Printf("DEBUGLL makeRuleHash %v inserting %v\n", rule.Name, to.URL)
		h.Insert(&ATSConsistentHashNode{Name: to.URL}, *to.Weight)
	}

	if h.First() == nil {
		fmt.Printf("DEBUGLL makeRuleHash %v NodeMap empty!\n", rule.Name)
	}
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
		if toJSON.TimeoutMS != nil {
			t := time.Duration(*toJSON.TimeoutMS) * time.Millisecond
			if to.Timeout = &t; *to.Timeout < 0 {
				return nil, fmt.Errorf("error parsing to %v timeout must be positive: %v", to.URL, to.Timeout)
			}
		} else {
			to.Timeout = rule.Timeout
		}
		if toJSON.ParentSelection != nil {
			ps := ParentSelectionTypeFromString(*toJSON.ParentSelection)
			if to.ParentSelection = &ps; *to.ParentSelection == ParentSelectionTypeInvalid {
				return nil, fmt.Errorf("error parsing to %v parent selection invalid: '%v'", to.URL, *toJSON.ParentSelection)
			}
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
