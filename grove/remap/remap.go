package remap

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
	"github.com/apache/incubator-trafficcontrol/grove/icache"
	"github.com/apache/incubator-trafficcontrol/grove/plugin"
	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

type HTTPRequestRemapper interface {
	// Remap returns the remapped request, the matched rule name, whether the requestor's IP is allowed, whether to connection close, whether a match was found, and any error.
	// Remap(r *http.Request, scheme string, failures int) Remapping
	Rules() []remapdata.RemapRule
	RemappingProducer(r *http.Request, scheme string) (*RemappingProducer, error)
	StatRules() remapdata.RemapRulesStats
	PluginCfg() map[string]interface{} // global plugins, outside the individual remap rules
}

type simpleHTTPRequestRemapper struct {
	remapper Remapper
	stats    *remapdata.RemapRulesStats
}

func (hr simpleHTTPRequestRemapper) Rules() []remapdata.RemapRule         { return hr.remapper.Rules() }
func (hr simpleHTTPRequestRemapper) StatRules() remapdata.RemapRulesStats { return *hr.stats }
func (hr simpleHTTPRequestRemapper) PluginCfg() map[string]interface{}    { return hr.remapper.PluginCfg() }

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
	Cache           icache.Cache
}

// RemappingProducer takes an HTTP Request and returns a Remapping to be used for that request.
// TODO rename? interface?
type RemappingProducer struct {
	oldURI   string
	rule     remapdata.RemapRule
	cacheKey string
	failures int
}

func (p *RemappingProducer) CacheKey() string                  { return p.cacheKey }
func (p *RemappingProducer) ConnectionClose() bool             { return p.rule.ConnectionClose }
func (p *RemappingProducer) Name() string                      { return p.rule.Name }
func (p *RemappingProducer) DSCP() int                         { return p.rule.DSCP }
func (p *RemappingProducer) PluginCfg() map[string]interface{} { return p.rule.Plugins }
func (p *RemappingProducer) Cache() icache.Cache               { return p.rule.Cache }
func (p *RemappingProducer) ToFQDN() string {
	// TODO verify To is not allowed to be constructed with < 1 element
	return strings.TrimPrefix(strings.TrimPrefix(p.rule.To[0].URL, "http://"), "https://")
}
func (p *RemappingProducer) ProxyStr() string {
	if p.rule.To[0].ProxyURL != nil && p.rule.To[0].ProxyURL.Host != "" {
		return p.rule.To[0].ProxyURL.Host
	}
	return "NONE" // TODO const?
}

var ErrRuleNotFound = errors.New("remap rule not found")
var ErrIPNotAllowed = errors.New("IP not allowed")
var ErrNoMoreRetries = errors.New("retry num exceeded")

// RequestURI returns the URI of the given request. This must be used, because Go does not populate the scheme of requests that come in from clients.
func RequestURI(r *http.Request, scheme string) string {
	return scheme + "://" + r.Host + r.RequestURI
}
func (hr simpleHTTPRequestRemapper) RemappingProducer(r *http.Request, scheme string) (*RemappingProducer, error) {
	uri := RequestURI(r, scheme)
	rule, ok := hr.remapper.Remap(uri)
	if !ok {
		return nil, ErrRuleNotFound
	}

	if ip, err := web.GetIP(r); err != nil {
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

	newURI, proxyURL := p.rule.URI(p.oldURI, r.URL.Path, r.URL.RawQuery, p.failures)
	p.failures++
	newReq, err := http.NewRequest(r.Method, newURI, nil)
	if err != nil {
		return Remapping{}, false, fmt.Errorf("creating new request: %v\n", err)
	}
	web.CopyHeaderTo(r.Header, &newReq.Header)

	log.Debugf("GetNext oldUri: %v, Host: %v\n", p.oldURI, newReq.Header.Get("Host"))
	log.Debugf("GetNext newURI: %v, fqdn: %v\n", newURI, getFQDN(newURI))
	log.Debugf("GetNext rule name: %v\n", p.rule.Name)

	newReq.Header.Set("Host", getFQDN(newURI))

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
		Cache:           p.rule.Cache,
	}, retryAllowed, nil
}

func RemapperToHTTP(r Remapper, statRules *remapdata.RemapRulesStats) HTTPRequestRemapper {
	return simpleHTTPRequestRemapper{remapper: r, stats: statRules}
}

func NewHTTPRequestRemapper(remap []remapdata.RemapRule, plugins map[string]interface{}, statRules *remapdata.RemapRulesStats) HTTPRequestRemapper {
	return RemapperToHTTP(NewLiteralPrefixRemapper(remap, plugins), statRules)
}

// Remapper provides a function which takes strings and maps them to other strings. This is designed for URL prefix remapping, for a reverse proxy.
type Remapper interface {
	// Remap returns the given string remapped, the unique name of the rule found, and whether a remap rule was found
	Remap(uri string) (remapdata.RemapRule, bool)
	// Rules returns the unique names of every remap rule.
	Rules() []remapdata.RemapRule
	PluginCfg() map[string]interface{} // global plugins, outside the individual remap rules
}

// TODO change to use a prefix tree, for speed
type literalPrefixRemapper struct {
	remap   []remapdata.RemapRule
	plugins map[string]interface{}
}

func (r literalPrefixRemapper) PluginCfg() map[string]interface{} { return r.plugins }

// Remap returns the remapped string, the remap rule name, the remap rule's options, and whether a remap was found
func (r literalPrefixRemapper) Remap(s string) (remapdata.RemapRule, bool) {
	for _, rule := range r.remap {
		if strings.HasPrefix(s, rule.From) {
			return rule, true
		}
	}
	return remapdata.RemapRule{}, false
}

func (r literalPrefixRemapper) Rules() []remapdata.RemapRule {
	rules := make([]remapdata.RemapRule, len(r.remap))
	for _, rule := range r.remap {
		rules = append(rules, rule)
	}
	return rules
}

func NewLiteralPrefixRemapper(remap []remapdata.RemapRule, plugins map[string]interface{}) Remapper {
	return literalPrefixRemapper{remap: remap, plugins: plugins}
}

type RemapRulesStatsJSON struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

type RemapRulesBase struct {
	RetryNum *int `json:"retry_num"`
}

type RemapRulesJSON struct {
	RemapRulesBase
	Rules           []RemapRuleJSON            `json:"rules"`
	RetryCodes      *[]int                     `json:"retry_codes"`
	TimeoutMS       *int                       `json:"timeout_ms"`
	ParentSelection *string                    `json:"parent_selection"`
	Stats           RemapRulesStatsJSON        `json:"stats"`
	Plugins         map[string]json.RawMessage `json:"plugins"`
}

type RemapRules struct {
	RemapRulesBase
	Rules           []remapdata.RemapRule
	RetryCodes      map[int]struct{}
	Timeout         *time.Duration
	ParentSelection *remapdata.ParentSelectionType
	Stats           remapdata.RemapRulesStats
	Plugins         map[string]interface{}
	Cache           icache.Cache
}

type RemapRuleToJSON struct {
	remapdata.RemapRuleToBase
	ProxyURL   *string `json:"proxy_url"`
	TimeoutMS  *int    `json:"timeout_ms"`
	RetryCodes *[]int  `json:"retry_codes"`
}

type HdrModder interface {
	Mod(h *http.Header)
}

type RemapRuleJSON struct {
	remapdata.RemapRuleBase
	TimeoutMS       *int                       `json:"timeout_ms"`
	ParentSelection *string                    `json:"parent_selection"`
	To              []RemapRuleToJSON          `json:"to"`
	Allow           []string                   `json:"allow"`
	Deny            []string                   `json:"deny"`
	RetryCodes      *[]int                     `json:"retry_codes"`
	CacheName       *string                    `json:"cache_name"`
	Plugins         map[string]json.RawMessage `json:"plugins"`
}

// LoadRemapRules returns the loaded rules, the global plugins, the Stats remap rules, and any error
func LoadRemapRules(path string, pluginConfigLoaders map[string]plugin.LoadFunc, caches map[string]icache.Cache) ([]remapdata.RemapRule, map[string]interface{}, *remapdata.RemapRulesStats, error) {
	fmt.Printf("Loading Remap Rules\n")
	defer func() {
		fmt.Printf("Loaded Remap Rules\n")
	}()
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	remapRulesJSON := RemapRulesJSON{}
	if err := json.NewDecoder(file).Decode(&remapRulesJSON); err != nil {
		return nil, nil, nil, fmt.Errorf("decoding JSON: %s", err)
	}

	remapRules := RemapRules{RemapRulesBase: remapRulesJSON.RemapRulesBase}

	if remapRulesJSON.RetryCodes != nil {
		remapRules.RetryCodes = make(map[int]struct{}, len(*remapRulesJSON.RetryCodes))
		for _, code := range *remapRulesJSON.RetryCodes {
			if _, ok := ValidHTTPCodes[code]; !ok {
				return nil, nil, nil, fmt.Errorf("error parsing rules: retry code invalid: %v", code)
			}
			remapRules.RetryCodes[code] = struct{}{}
		}
	}
	if remapRulesJSON.TimeoutMS != nil {
		t := time.Duration(*remapRulesJSON.TimeoutMS) * time.Millisecond
		if remapRules.Timeout = &t; *remapRules.Timeout < 0 {
			return nil, nil, nil, fmt.Errorf("error parsing rules: timeout must be positive: %v", remapRules.Timeout)
		}
	}
	if remapRulesJSON.ParentSelection != nil {
		ps := remapdata.ParentSelectionTypeFromString(*remapRulesJSON.ParentSelection)
		if remapRules.ParentSelection = &ps; *remapRules.ParentSelection == remapdata.ParentSelectionTypeInvalid {
			return nil, nil, nil, fmt.Errorf("error parsing rules: parent selection invalid: '%v'", remapRulesJSON.ParentSelection)
		}
	}
	if remapRulesJSON.Stats.Allow != nil {
		if remapRules.Stats.Allow, err = makeIPNets(remapRulesJSON.Stats.Allow); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing rules allows: %v", err)
		}
	}
	if remapRulesJSON.Stats.Deny != nil {
		if remapRules.Stats.Deny, err = makeIPNets(remapRulesJSON.Stats.Deny); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing rules denys: %v", err)
		}
	}

	remapRules.Plugins = make(map[string]interface{}, len(remapRulesJSON.Plugins))
	for name, b := range remapRulesJSON.Plugins {
		if loadF := pluginConfigLoaders[name]; loadF != nil {
			remapRules.Plugins[name] = loadF(b)
		}
	}

	rules := make([]remapdata.RemapRule, len(remapRulesJSON.Rules))
	for i, jsonRule := range remapRulesJSON.Rules {
		fmt.Printf("Creating Remap Rule %v\n", jsonRule.Name)
		rule := remapdata.RemapRule{RemapRuleBase: jsonRule.RemapRuleBase}

		rule.Plugins = make(map[string]interface{}, len(jsonRule.Plugins))
		for name, b := range jsonRule.Plugins {
			if loadF := pluginConfigLoaders[name]; loadF != nil {
				rule.Plugins[name] = loadF(b)
			}
		}
		for name, loader := range remapRules.Plugins {
			if _, ok := rule.Plugins[name]; !ok {
				rule.Plugins[name] = loader
			}
		}

		if jsonRule.RetryCodes != nil {
			rule.RetryCodes = make(map[int]struct{}, len(*jsonRule.RetryCodes))
			for _, code := range *jsonRule.RetryCodes {
				if _, ok := ValidHTTPCodes[code]; !ok {
					return nil, nil, nil, fmt.Errorf("error parsing rule %v retry code invalid: %v", rule.Name, code)
				}
				rule.RetryCodes[code] = struct{}{}
			}
		} else {
			rule.RetryCodes = remapRules.RetryCodes
		}

		if jsonRule.TimeoutMS != nil {
			t := time.Duration(*jsonRule.TimeoutMS) * time.Millisecond
			if rule.Timeout = &t; *rule.Timeout < 0 {
				return nil, nil, nil, fmt.Errorf("error parsing rule %v timeout must be positive: %v", rule.Name, rule.Timeout)
			}
		} else {
			rule.Timeout = remapRules.Timeout
		}

		if rule.RetryNum == nil {
			rule.RetryNum = remapRules.RetryNum
		}

		cacheName := "" // default string is the default cache
		if jsonRule.CacheName != nil {
			cacheName = *jsonRule.CacheName
		}
		ok := false
		if rule.Cache, ok = caches[cacheName]; !ok {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v: cache name %v not found", rule.Name, cacheName)
		}

		if rule.Allow, err = makeIPNets(jsonRule.Allow); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v allows: %v", rule.Name, err)
		}
		if rule.Deny, err = makeIPNets(jsonRule.Deny); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v denys: %v", rule.Name, err)
		}
		if rule.To, err = makeTo(jsonRule.To, rule); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v to: %v", rule.Name, err)
		}
		if jsonRule.ParentSelection != nil {
			ps := remapdata.ParentSelectionTypeFromString(*jsonRule.ParentSelection)
			if rule.ParentSelection = &ps; *rule.ParentSelection == remapdata.ParentSelectionTypeInvalid {
				return nil, nil, nil, fmt.Errorf("error parsing rule %v parent selection invalid: '%v'", rule.Name, jsonRule.ParentSelection)
			}
		} else {
			rule.ParentSelection = remapRules.ParentSelection
		}

		if rule.ParentSelection == nil {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v - no parent_selection - must be set at rules or rule level", rule.Name)
		}

		if len(rule.To) == 0 {
			return nil, nil, nil, fmt.Errorf("error parsing rule %v - no to - must have at least one parent", rule.Name)
		}

		if *rule.ParentSelection == remapdata.ParentSelectionTypeConsistentHash {
			rule.ConsistentHash = makeRuleHash(rule)
		} else {
		}
		rules[i] = rule
	}

	return rules, remapRules.Plugins, &remapRules.Stats, nil
}

const DefaultReplicas = 1024

func makeRuleHash(rule remapdata.RemapRule) chash.ATSConsistentHash {
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

func makeTo(tosJSON []RemapRuleToJSON, rule remapdata.RemapRule) ([]remapdata.RemapRuleTo, error) {
	tos := make([]remapdata.RemapRuleTo, len(tosJSON))
	for i, toJSON := range tosJSON {
		if toJSON.Weight == nil {
			w := 1.0
			toJSON.Weight = &w
		}
		to := remapdata.RemapRuleTo{RemapRuleToBase: toJSON.RemapRuleToBase}
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
				if _, ok := ValidHTTPCodes[code]; !ok {
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

func LoadRemapper(path string, pluginConfigLoaders map[string]plugin.LoadFunc, caches map[string]icache.Cache) (HTTPRequestRemapper, error) {
	rules, plugins, statRules, err := LoadRemapRules(path, pluginConfigLoaders, caches)
	if err != nil {
		return nil, err
	}
	return NewHTTPRequestRemapper(rules, plugins, statRules), nil
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
		for code := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, code)
		}
	}
	if r.ParentSelection != nil {
		s := ""
		j.ParentSelection = &s
		*j.ParentSelection = string(*r.ParentSelection)
	}
	for _, deny := range r.Stats.Deny {
		j.Stats.Deny = append(j.Stats.Deny, deny.String())
	}
	for _, allow := range r.Stats.Allow {
		j.Stats.Allow = append(j.Stats.Allow, allow.String())
	}

	for _, rule := range r.Rules {
		j.Rules = append(j.Rules, buildRemapRuleToJSON(rule))
	}
	j.Plugins = make(map[string]json.RawMessage)
	for name, plugin := range r.Plugins {
		clientHeadersJSONBytes, _ := json.Marshal(plugin)
		j.Plugins[name] = clientHeadersJSONBytes
	}
	return j
}

func buildRemapRuleToJSON(r remapdata.RemapRule) RemapRuleJSON {
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
		for retryCode := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, retryCode)
		}
	}
	j.Plugins = make(map[string]json.RawMessage)
	for name, plugin := range r.Plugins {
		clientHeadersJSONBytes, _ := json.Marshal(plugin)
		j.Plugins[name] = clientHeadersJSONBytes
	}
	return j
}

func RemapRuleToToJSON(r remapdata.RemapRuleTo) RemapRuleToJSON {
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
		for retryCode := range r.RetryCodes {
			*j.RetryCodes = append(*j.RetryCodes, retryCode)
		}
	}
	return j
}
