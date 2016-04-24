//  Started with https://github.com/nf/webfront/blob/master/main.go
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
	"strings"
	"sync"
	"time"
)

// Server implements an http.Handler that acts as a reverse proxy
type Server struct {
	mu    sync.RWMutex // guards the fields below
	last  time.Time
	rules []*Rule
}

// Rule represents a rule in a configuration file.
type Rule struct {
	Host    string // to match against request Host header
	Path    string // to match against a path (start)
	Forward string // reverse proxy map-to
	Secure  bool   // protect with jwt?

	handler http.Handler
}

type loginJson struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// Config holds the configuration of the server.
type Config struct {
	HTTPSAddr    string `json:"httpsAddr"`
	RuleFile     string `json:"ruleFile"`
	PollInterval int    `json:"pollInterval"`
}

var Logger *log.Logger

func printUsage() {
	exampleConfig := `{
	"httpsAddr": ":9000",
	"ruleFile": "rules.json",
	"pollInterval": 60
}`
	Logger.Println("Usage: " + path.Base(os.Args[0]) + " configfile")
	Logger.Println("")
	Logger.Println("Example config file:")
	Logger.Println(exampleConfig)
}

func main() {
	if len(os.Args) < 2 {
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
	Logger.Println("Starting webfront...")
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
	http.ListenAndServeTLS(config.HTTPSAddr, "server.pem", "server.key", s)
}

func validateToken(tokenString string) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("CAmeRAFiveSevenNineNine"), nil
	})
	return token, err
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
// request with the Rule's handler. If the rule's secure field is true, it will
// only allow access if the request has a valid JWT bearer token.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isSecure := s.isSecure(r)
	tokenValid := false
	token, err := validateToken(r.Header.Get("Authorization"))
	if err == nil {
		tokenValid = true
	} else {
		Logger.Println("Token Error:", err.Error())
	}

	if isSecure {
		if !tokenValid {
			Logger.Println(r.URL.Path + ": valid token required, but none found!")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		tokenUser := token.Claims["User"]

		re := regexp.MustCompile("[^/]+")
		params := re.FindAllString(r.URL.Path, -1)
		// Logger.Println(">>>", r.URL.Path, " >>> ", len(params))
		if len(params) < 2 {
			Logger.Println("Invalid path: ", r.URL.Path)
			// TODO add root user exemption here.
			http.Error(w, "Invalid request - user not found.", http.StatusBadRequest)
			return
		}
		pathUser := params[1]
		if pathUser != tokenUser {
			if tokenUser != "root" {
				Logger.Println(r.Method+" "+r.URL.RequestURI()+": valid token found, identified user:", tokenUser, " != ", pathUser, " - deny")
				http.Error(w, "Not Authorized", http.StatusUnauthorized)
				return
			} else {
				Logger.Println(r.Method+" "+r.URL.RequestURI()+": valid token found, identified user:", tokenUser, " - allow")
			}
		}
		Logger.Println(r.Method+" "+r.URL.RequestURI()+": valid token found, identified user:", token.Claims["User"], " matches path - allow")
	} else {
		Logger.Println(r.Method + " " + r.URL.RequestURI() + ": no token required - allow")
	}

	if h := s.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}
	Logger.Println(r.Method + " " + r.URL.Path + ": no mapping in rules file!")
	http.Error(w, "Not found.", http.StatusNotFound)
}

func rejectNoToken(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}
}

// isSecure returns the true if this path should be protected by a jwt
func (s *Server) isSecure(req *http.Request) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h := req.Host
	p := req.URL.Path
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	for _, r := range s.rules {
		// Logger.Println(p, "==", r.Path)
		if strings.HasPrefix(p, r.Path) {
			return r.Secure
		}
	}
	Logger.Println("returning hard false")
	return true
}

// handler returns the appropriate Handler for the given Request,
// or nil if none found.
func (s *Server) handler(req *http.Request) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h := req.Host
	p := req.URL.Path
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	for _, r := range s.rules {
		// Logger.Println(p, "==", r.Path)
		if strings.HasPrefix(p, r.Path) {
			return r.handler
		}
	}
	Logger.Println("returning nil")
	return nil
}

// refreshRules polls file periodically and refreshes the Server's rule
// set if the file has been modified.
func (s *Server) refreshRules(file string, poll time.Duration) {
	for {
		// Logger.Println("loading file")
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
			Logger.Printf("bad rule: %#v", r)
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
