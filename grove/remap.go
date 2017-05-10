package grove

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type HTTPRequestRemapper interface {
	// Remap returns the remapped request, the matched rule name, and whether a match was found.
	Remap(r *http.Request, scheme string) (*http.Request, string, string, bool)
	Rules() []RemapRule
}

type simpleHttpRequestRemapper struct {
	remapper Remapper
}

func (hr simpleHttpRequestRemapper) Rules() []RemapRule {
	return hr.remapper.Rules()
}

// Remap returns the given request with its URI remapped, the name of the remap rule found, the cache key, and whether a rule was found.
func (hr simpleHttpRequestRemapper) Remap(r *http.Request, scheme string) (*http.Request, string, string, bool) {
	// NewRequest(method, urlStr string, body io.Reader)
	// TODO config whether to consider query string, method, headers
	oldUri := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	fmt.Printf("DEBUG Remap oldUri: '%v'\n", oldUri)
	fmt.Printf("DEBUG request: '%+v'\n", r)
	// newUri, ruleName, options, ok :=
	rule, ok := hr.remapper.Remap(oldUri)
	if !ok {
		fmt.Printf("DEBUG Remap oldUri: '%v' NOT FOUND\n", oldUri)
		return r, "", "", false
	}

	newUri := rule.URI(oldUri)
	cacheKey := rule.CacheKey(r.Method, oldUri)
	fmt.Printf("DEBUG Remap newURI: '%v'\nDEBUG Remap cacheKey '%v'\n", newUri, cacheKey)

	newReq, err := http.NewRequest(r.Method, newUri, nil) // TODO modify given req in-place?
	if err != nil {
		fmt.Printf("Error Remap NewRequest: %v\n", err)
		return r, "", "", false
	}
	copyHeader(r.Header, &newReq.Header)
	return newReq, rule.Name, cacheKey, true
}

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

type RemapRulesJSON struct {
	Rules []RemapRule `json:"rules"`
}

type RemapRule struct {
	Name        string          `json:"name"`
	From        string          `json:"from"`
	To          string          `json:"to"`
	QueryString QueryStringRule `json:"query-string"`
	// ConcurrentRuleRequests is the number of concurrent requests permitted to a remap rule, that is, to an origin. If this is 0, the global config is used.
	ConcurrentRuleRequests int `json:"concurrent_rule_requests"`
}

type QueryStringRule struct {
	Remap bool `json:"remap"`
	Cache bool `json:"cache"`
}

func (r RemapRule) URI(fromURI string) string {
	uri := r.To + fromURI[len(r.From):]
	if !r.QueryString.Remap {
		if i := strings.Index(uri, "?"); i != -1 {
			uri = uri[:i]
		}
	}
	return uri
}

func (r RemapRule) CacheKey(method string, fromURI string) string {
	uri := r.To + fromURI[len(r.From):]
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
	remapRules := RemapRulesJSON{}
	if err := json.NewDecoder(file).Decode(&remapRules); err != nil {
		return nil, fmt.Errorf("decoding JSON: %s", err)
	}
	return remapRules.Rules, nil
}

func LoadRemapper(path string) (HTTPRequestRemapper, error) {
	rules, err := LoadRemapRules(path)
	if err != nil {
		return nil, err
	}
	return NewHTTPRequestRemapper(rules), nil
}
