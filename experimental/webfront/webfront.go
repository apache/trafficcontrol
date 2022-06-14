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
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// TODO(amiry) - Handle refresh tokens

// Server implements an http.Handler that acts as a reverse proxy
type Server struct {
	mu    sync.RWMutex // guards the fields below
	last  time.Time
	rules []*Rule
}

// Rule represents a rule in a configuration file.
type Rule struct {
	Host         string            // to match against request Host header
	Path         string            // to match against a path (start)
	Forward      string            // reverse proxy map-to
	Secure       bool              // protect with jwt?
	Capabilities map[string]string // map HTTP methods to capabilitues

	handler http.Handler
}

// Config holds the configuration of the server.
type Config struct {
	RuleFile     string `json:"rule-file"`
	PollInterval int    `json:"poll-interval"`
	ListenPort   int    `json:"listen-port"`
}

var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"listen-port":   9000,
	"rule-file":     "rules.json",
	"poll-interval": 60
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

	s, err := NewServer(config.RuleFile, time.Duration(config.PollInterval)*time.Second)
	if err != nil {
		Logger.Fatal(err)
	}

	// override the default so we can use self-signed certs on our microservices
	// and use a self-signed cert in this server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if _, err := os.Stat("server.pem"); os.IsNotExist(err) {
		Logger.Fatal("server.pem file not found")
	}
	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		Logger.Fatal("server.key file not found")
	}

	Logger.Printf("Starting webfront on port %d...", config.ListenPort)
	Logger.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(int(config.ListenPort)), "server.pem", "server.key", s))
}

// NewServer constructs a Server that reads rules from file with a period
// specified by poll.
func NewServer(file string, poll time.Duration) (*Server, error) {
	s := new(Server)
	if err := s.loadRules(file); err != nil {
		Logger.Fatal("Error loading rules file: ", err)
	}
	go s.refreshRules(file, poll)
	return s, nil
}

// ServeHTTP matches the Request with a Rule and, if found, serves the
// request with the Rule's handler. If the rule's secure field is true, it will
// only allow access if the request has a valid JWT bearer token.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rule := s.getRule(r)
	if rule == nil {
		Logger.Printf("%v %v No mapping in rules file!", r.Method, r.URL.RequestURI())
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	isAuthorized := false

	if rule.Secure {
		token, err := jwt.ParseHeader(
			r.Header,
			`Authorization`,
			jwt.WithVerify(jwa.HS256, []byte(os.Args[2])),
		)
		if err != nil {
			Logger.Println("Token Error:", err.Error())
			Logger.Printf("%v %v Valid token required, but none found!", r.Method, r.URL.RequestURI())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Authorization: Check is the list of capabilities in the token's claims contains
		// the reqired capability that is listed in the rule
		var capabilities []string
		if raw, ok := token.Get(`cap`); ok {
			if caps, ok := raw.([]string); ok {
				capabilities = caps // Save this to use in the logging later
				for _, c := range caps {
					if c == rule.Capabilities[r.Method] {
						isAuthorized = true
						break
					}
				}
			}
		}

		Logger.Printf("%v %v Valid token. Subject=%v, ExpiresAt=%v, Capabilities=%v, Required=%v, Authorized=%v",
			r.Method, r.URL.RequestURI(), token.Subject(), token.Expiration(), capabilities,
			rule.Capabilities[r.Method], isAuthorized)

	} else {
		isAuthorized = true
	}

	if isAuthorized {
		if h := rule.handler; h != nil {
			h.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Not Authorized", http.StatusUnauthorized)
	return
}

func rejectNoToken(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}
}

func (s *Server) getRule(req *http.Request) *Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h := req.Host
	p := req.URL.Path

	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}

	for _, r := range s.rules {
		if strings.HasPrefix(p, r.Path) {
			// Logger.Printf("Found rule")
			return r
		}
	}

	// Logger.Printf("Rule not found")
	return nil
}

// refreshRules polls file periodically and refreshes the Server's rule
// set if the file has been modified.
func (s *Server) refreshRules(file string, poll time.Duration) {
	for {
		// Logger.Printf("loading rule file")
		if err := s.loadRules(file); err != nil {
			Logger.Println(file, ":", err)
		}
		time.Sleep(poll)
	}
}

// loadRules tests whether file has been modified since its last invocation
// and, if so, loads the rule set from file.
func (s *Server) loadRules(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	mtime := fi.ModTime()
	if !mtime.After(s.last) && s.rules != nil {
		return nil // no change
	}
	rules, err := parseRules(file)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.last = mtime
	s.rules = rules
	s.mu.Unlock()
	return nil
}

// parseRules reads rule definitions from file, constructs the Rule handlers,
// and returns the resultant Rules.
func parseRules(file string) ([]*Rule, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rules []*Rule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, err
	}
	for _, r := range rules {
		r.handler = makeHandler(r)
		if r.handler == nil {
			Logger.Printf("Bad rule: %#v", r)
		}
	}
	return rules, nil
}

// makeHandler constructs the appropriate Handler for the given Rule.
func makeHandler(r *Rule) http.Handler {
	if h := r.Forward; h != "" {
		return &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "https"
				req.URL.Host = h
				// req.URL.Path = "/boo1" // TODO JvD - regex to change path here
			},
		}
	}
	return nil
}
