//  Started with https://github.com/nf/webfront/blob/master/main.go
// by Andrew Gerrand <adg@golang.org>

/*
Copyright 2011 Google Inc.

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

/*
webfront is an HTTP server and reverse proxy.

It reads a JSON-formatted rule file like this:

	[
		{"Host": "example.com", "Serve": "/var/www"},
		{"Host": "example.org", "Forward": "localhost:8080"}
	]

For all requests to the host example.com (or any name ending in
".example.com") it serves files from the /var/www directory.

For requests to example.org, it forwards the request to the HTTP
server listening on localhost port 8080.

Usage of webfront:
  -http=":80": HTTP listen address
  -https="": HTTPS listen address (leave empty to disable)
  -https_cert="": HTTPS certificate file
  -https_key="": HTTPS key file
  -poll=10s: file poll interval
  -rules="": rule definition file

webfront was written by Andrew Gerrand <adg@golang.org>
*/
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	httpAddr     = flag.String("http", ":80", "HTTP listen address")
	httpsAddr    = flag.String("https", "", "HTTPS listen address (leave empty to disable)")
	certFile     = flag.String("https_cert", "", "HTTPS certificate file")
	keyFile      = flag.String("https_key", "", "HTTPS key file")
	ruleFile     = flag.String("rules", "", "rule definition file")
	pollInterval = flag.Duration("poll", time.Second*10, "file poll interval")
)

func main() {
	flag.Parse()
	s, err := NewServer(*ruleFile, *pollInterval)
	if err != nil {
		log.Fatal(err)
	}
	httpFD, _ := strconv.Atoi(os.Getenv("RUNSIT_PORTFD_http"))
	httpsFD, _ := strconv.Atoi(os.Getenv("RUNSIT_PORTFD_https"))
	if httpsFD >= 3 || *httpsAddr != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}
		c := &tls.Config{Certificates: []tls.Certificate{cert}}
		l := tls.NewListener(listen(httpsFD, *httpsAddr), c)
		go func() {
			log.Fatal(http.Serve(l, s))
		}()
	}
	log.Fatal(http.Serve(listen(httpFD, *httpAddr), s))
}

func listen(fd int, addr string) net.Listener {
	var l net.Listener
	var err error
	if fd >= 3 {
		l, err = net.FileListener(os.NewFile(uintptr(fd), "http"))
	} else {
		l, err = net.Listen("tcp", addr)
	}
	if err != nil {
		log.Fatal(err)
	}
	return l
}

// Server implements an http.Handler that acts as either a reverse proxy or
// a simple file server, as determined by a rule set.
type Server struct {
	mu    sync.RWMutex // guards the fields below
	last  time.Time
	rules []*Rule
}

// Rule represents a rule in a configuration file.
type Rule struct {
	Host    string // to match against request Host header
	Forward string // non-empty if reverse proxy
	Serve   string // non-empty if file server

	handler http.Handler
}

// NewServer constructs a Server that reads rules from file with a period
// specified by poll.
func NewServer(file string, poll time.Duration) (*Server, error) {
	s := new(Server)
	if err := s.loadRules(file); err != nil {
		return nil, err
	}
	go s.refreshRules(file, poll)
	return s, nil
}

// ServeHTTP matches the Request with a Rule and, if found, serves the
// request with the Rule's handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not found.", http.StatusNotFound)
}

// handler returns the appropriate Handler for the given Request,
// or nil if none found.
func (s *Server) handler(req *http.Request) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h := req.Host
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	for _, r := range s.rules {
		if h == r.Host || strings.HasSuffix(h, "."+r.Host) {
			return r.handler
		}
	}
	return nil
}

// refreshRules polls file periodically and refreshes the Server's rule
// set if the file has been modified.
func (s *Server) refreshRules(file string, poll time.Duration) {
	for {
		if err := s.loadRules(file); err != nil {
			log.Println(err)
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
			log.Printf("bad rule: %#v", r)
		}
	}
	return rules, nil
}

// makeHandler constructs the appropriate Handler for the given Rule.
func makeHandler(r *Rule) http.Handler {
	if h := r.Forward; h != "" {
		return &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = h
			},
		}
	}
	if d := r.Serve; d != "" {
		return http.FileServer(http.Dir(d))
	}
	return nil
}
