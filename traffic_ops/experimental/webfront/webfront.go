/*

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

// Started with https://github.com/nf/webfront/blob/master/main.go
// by Andrew Gerrand <adg@golang.org>

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO(amiry) - Refresh tokens
// TODO(amiry) - Rewrite rules
// TODO(amiry) - Access log
// TODO(amiry) - Cantralized, managed route configuration
// TODO(amiry) - Auth server: Legacy token expiration should be longer than JWT expiration

// TODO(amiry) - Test regex match performance
// TODO(amiry) - Test/Document: Deprecate API with empty "auth" object in json
// TODO(amiry) - Add "/" route for admin user? This will cause non existant routes to return Forbidden instead of Not Found
//    { "match": "/.*", "auth": { "GET": ["all-read"], "POST": ["all-write"], "PUT": ["all-write"], "PATCH": ["all-write"], "DELETE": ["all-write"] }}


// Config holds the configuration of the server.
type Config struct {
	ListenPort   			int    	`json:"listen-port"`
	RuleFile     			string 	`json:"rule-file"`
	PollInterval 			int    	`json:"poll-interval"`
	CrtFile					string  `json:"crt-file"`
	KeyFile					string  `json:"key-file"`
	InsecureSkipVerify 		bool	`json:"insecure-skip-verify"`
}

// Server implements an http.Handler that acts as a reverse proxy
type Server struct {
	mu    sync.RWMutex // guards the fields below
	last  time.Time
	Rules []*FwdRule
}

// FwdRule represents a FwdRule in a configuration file
type FwdRule struct {
	Host           string     							// to match against request Host header
	Path           string     							// to match against a path (start)
	Forward        string     							// reverse proxy map-to
	Scheme         string     							// reverse proxy URL scheme (HTTP)
	Auth           bool       							// protect with jwt?
	RoutesFile     string     `json:"routes-file"`		// path to routes file

	routes         []*Route
	handler        http.Handler
}

type Route struct {
	Match   string 										// the route's path regex
	Auth    map[string]([]string)						// map a HTTP method to the capabilities that are required
														// to perform the method on this route. Methods that are not 
														// in this list are forbidden.
	// A compiled regex for "Match"
	matchRegexp *regexp.Regexp
}

type Claims struct {
    Capabilities []string      `json:"cap"`
    LegacyCookie string        `json:"legacy-cookie"`	// LEGACY: The legacy cookie to be set upon request
    jwt.StandardClaims
}

var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"listen-port":   8080,
	"rule-file":     "rules.json",
	"poll-interval": 5,
	"crt-file":      "server.crt",
	"key-file":      "server.key",
	"insecure-skip-verify": false
}`
	fmt.Println("Usage: " + path.Base(os.Args[0]) + " config-file secret")
	fmt.Println("")
	fmt.Println("Example config-file:")
	fmt.Println(exampleConfig)
}

func main() {

	if len(os.Args) < 3 {
		printUsage()
		return
	}

	Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	file, err := os.Open(os.Args[1])
	if err != nil {
		Logger.Println("Error opening config file:", err)
		return
	}

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		Logger.Println("Error reading config file:", err)
		return
	}

	if _, err := os.Stat(config.CrtFile); os.IsNotExist(err) {
		Logger.Fatalf("%s file not found", config.CrtFile)
	}
	if _, err := os.Stat(config.KeyFile); os.IsNotExist(err) {
		Logger.Fatalf("%s file not found", config.KeyFile)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = makeTLSConfig(&config)
	s, err := NewServer(config.RuleFile, time.Duration(config.PollInterval)*time.Second)
	if err != nil {
		Logger.Fatal(err)
	}

	Logger.Printf("Starting webfront on port %d...", config.ListenPort)
	Logger.Fatal(http.ListenAndServeTLS(":" + strconv.Itoa(int(config.ListenPort)), config.CrtFile, config.KeyFile, s))
}

// NewServer constructs a Server that reads Rules from file with a period 
// specified by poll
func NewServer(file string, poll time.Duration) (*Server, error) {
	s := new(Server)
	if err := s.loadRules(file); err != nil {
		Logger.Fatal(fmt.Errorf("Load rules failed: %s", err))
	}

	// TODO(amiry) - Reload config using NOHUP signal instead of poll for changes
	go s.refreshRules(file, poll)

	return s, nil
}

func makeTLSConfig(config *Config) *tls.Config {

	s := false 
	if config.InsecureSkipVerify == true {
		Logger.Printf("NOTICE: Skip certificate verification")
		s = true
	}
	return &tls.Config{InsecureSkipVerify: s}
}

// loadRules tests whether file has been modified since its last invocation
// and, if so, loads the rule set from file.
func (s *Server) loadRules(file string) error {

	fi, err := os.Stat(file)
	if err != nil {
		return err
	}

	mtime := fi.ModTime()
	if !mtime.After(s.last) && s.Rules != nil {
		return nil // no change
	}

	Rules, err := parseRules(file)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.last = mtime
	s.Rules = Rules
	s.mu.Unlock()
	return nil
}

// refreshRules polls file periodically and refreshes the Server's rule set
// if the file has been modified.
func (s *Server) refreshRules(file string, poll time.Duration) {
	for {
		if err := s.loadRules(file); err != nil {
			Logger.Printf("Refresh rules failed: %s", err)
		}
		time.Sleep(poll)
	}
}

// parseRules reads rule definitions from file, constructs the rule handlers,
// and returns the resultant rules.
func parseRules(file string) ([]*FwdRule, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	Logger.Printf("Loading rules file: %s", file)

	var rules []*FwdRule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, err
	}

	for _, r := range rules {

		if r.Auth {
			r.routes, err = parseRoutes(r.RoutesFile)
			if err != nil {
				Logger.Printf("Skip rule %s ERROR: %s", r.Path, err)
				continue
			}			
		}

		r.handler, err = makeHandler(r)
		if err != nil {
			Logger.Printf("Skip rule %s ERROR: %s", r.Path, err)
			continue
		}

		// Logger.Printf("Loaded rule: %s", r.Path)
	}

	return rules, nil
}

// parseRoutes reads route definitions from file, constructs the route auth handler,
// and returns the resultant routes.
func parseRoutes(file string) ([]*Route, error) {

	// If the rule defines a routes file, we load the routes and enforce access.
	// Routes than are not present in this file are forbidden.

	// Note that there is currently no mechanism to trigger an update on a change in the route files. 
	// To trigger an update, one needs to touch rules.json 

	cf, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer cf.Close()

	Logger.Printf("Loading routes file: %s", file)

	var routes []*Route
	if err := json.NewDecoder(cf).Decode(&routes); err != nil {
		return nil, err
	}

	for _, r := range routes {

		/*
		// If the match ends with a slash, it is treated as a prefix. 
		// If not, it is an exact match
		if !strings.EndsWith(r.Match, "/") {
			r.Match = r.Match + '$'
		}
		*/

		r.matchRegexp, err = regexp.Compile(r.Match + "$")
		if err != nil {
			Logger.Printf("Skip route %s ERROR: %s", r.Match, err)
			continue
		}

		// Logger.Printf("Loaded route: %s", r.Match)
	}

	return routes, nil
}

// makeHandler constructs the appropriate Handler for the given FwdRule.
func makeHandler(r *FwdRule) (http.Handler, error) {

	host := r.Forward
	pathPrefix := "/"

	if i := strings.Index(r.Forward, "/"); i >= 0 {
		host = r.Forward[:i]
		pathPrefix = r.Forward[i:]
	}

	if host == "" {
		return nil, fmt.Errorf("Not a forward rule")
	}

	return &httputil.ReverseProxy {
		Director: func(req *http.Request) {
			req.URL.Scheme = r.Scheme
			req.URL.Host = host
			req.URL.Path = pathPrefix + strings.TrimPrefix(req.URL.Path, r.Path)
			Logger.Printf("Proxy: HOST: %s PATH: %s", req.URL.Host, req.URL.Path)
		},
	}, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// ServeHTTP matches the Request with a forward rule and, if found, serves the
// request with the rule's handler. If the rule's secure field is true, it will
// only allow access if the request has a valid JWT bearer token.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	rule := s.matchRule(req)
	if rule == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if rule.Auth {
		authorized := rule.authorize(w, req)
		if !authorized {
			return
		}			
	}

	if h := rule.handler; h != nil {
		h.ServeHTTP(w, req)
		return
	}

	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return
}

func (rule *FwdRule) authorize(w http.ResponseWriter, req *http.Request) bool {

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// LEGACY: If request contains a Mojo cookie instead of a JWT, we bypass token authorization 
	// and let legacy TO handle all authorization. 
	var cookie, err = req.Cookie("mojolicious")
	if cookie != nil {
		Logger.Printf("LEGACY: Found mojolicious cookie. Bypass authorization")
		return true
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////

	token, err := validateToken(req.Header.Get("Authorization"))

	if err != nil {
		Logger.Printf("%v %v Token error: %s", req.Method, req.URL.RequestURI(), err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		Logger.Printf("%v %v Token valid but cannot parse claims", req.Method, req.URL.RequestURI())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}

	route := rule.matchRoute(req)
	if route == nil {
		Logger.Printf("%v %v Route not found", req.Method, req.URL.RequestURI())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false
	}

    method := route.Auth[req.Method]
    if method == nil {
		Logger.Printf("%v %v Route found but method forbidden", req.Method, req.URL.RequestURI())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
    }

    // method is actually a list of capabilities required to perform this method.
    // Re performance - the lists are VERY short
    satisfied := len(method)
	for _, need := range method {
		for _, has := range claims.Capabilities {
        	if has == need {
        		satisfied--
        	}
        }
    }

    if (satisfied > 0) {
		Logger.Printf("%v %v Route found but required capabilities not satisfied. HAS %v, NEED %v", 
			req.Method, req.URL.RequestURI(), claims.Capabilities, method)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
    }

	Logger.Printf("%v %v Authorized. Subject=%v, ExpiresAt=%v, Rule=%s, Route=%s, Has=%v, Need=%v", 
		req.Method, req.URL.RequestURI(), claims.Subject, claims.ExpiresAt, rule.Path, route.Match, method, claims.Capabilities)

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// LEGACY: Pass legacy authentication token upon every secured request...
	legacyCookie := claims.LegacyCookie;
	req.Header.Add("Cookie", legacyCookie)			
	/////////////////////////////////////////////////////////////////////////////////////////////////

	return true
}

func validateToken(tokenString string) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Args[2]), nil
	})
	return token, err
}

func (s *Server) matchRule(req *http.Request) *FwdRule {

	s.mu.RLock()
	defer s.mu.RUnlock()

	// h := req.Host
	p := req.URL.Path

	/*
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	*/

	for _, r := range s.Rules {

		// Rules are matched in order! Longer rules should take precedence in rule file.
		// Logger.Printf("CHECK RULE: PATH %s BEGINS WITH %s ?", p, r.Path)		
		if strings.HasPrefix(p, r.Path) {
			// Logger.Printf("FOUND RULE: %s", r.Path)
			return r
		}
	}

	// Logger.Printf("Rule not found for path: %s", p)
	return nil
}

func (r *FwdRule) matchRoute(req *http.Request) *Route {

	// TODO(amiry) - Naive implementation

	Logger.Printf("MATCH ROUTE: PATH %s", req.URL.Path)

	// h := req.Host
	p := req.URL.Path

	/*
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	*/

	for _, r := range r.routes {

		// Routes are matched in order! Longer routes should take precedence in rule file.
		// Logger.Printf("CHECK ROUTE: PATH %s MATCHES %s ?", p, r.Match)
		if r.matchRegexp.MatchString(p) {
			// Logger.Printf("FOUND ROUTE: %s", r.matchRegexp)
			return r
		}
	}

	// Logger.Printf("Route not found for path: %s", p)
	return nil
}
