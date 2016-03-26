//  Started with https://github.com/nf/webfront/blob/master/main.go
// by Andrew Gerrand <adg@golang.org>

package main

import (
	// "crypto/sha1"
	"crypto/tls"
	// "encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

type TokenResponse struct {
	Token string
}

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

	handler http.Handler
}

type loginJson struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var (
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

	// override the default so we can use self-signed certs on our microservices
	// and use a self-signed cert in this server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.ListenAndServeTLS(*httpsAddr, *certFile, *keyFile, s)
}

func validateToken(tokenString string) (*jwt.Token, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("CAmeRAFiveSevenNineNine"), nil
	})

	if err == nil && token.Valid {
		log.Println("TOKEN IS GOOD -- user:", token.Claims["userid"], " role:", token.Claims["role"])
	} else {
		log.Println("TOKEN IS BAD", err)
	}
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
// request with the Rule's handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		username := ""
		password := ""
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body: ", err.Error())
			http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
			return
		}
		var lj loginJson
		log.Println(body)
		err = json.Unmarshal(body, &lj)
		if err != nil {
			log.Println("Error unmarshalling JSON: ", err.Error())
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		username = lj.User
		password = lj.Password

		// TODO JvD - check username / password against Database here!

		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["User"] = username
		token.Claims["Password"] = password
		token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString([]byte("CAmeRAFiveSevenNineNine"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(TokenResponse{Token: tokenString})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	token, err := validateToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("No valid token found!")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO JvD ^^ move into own function

	log.Println("Token:", token.Claims["userid"])

	if h := s.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not found.", http.StatusNotFound)
}

func rejectNoToken(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}
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
		// log.Println(p, "==", r.Path)
		if strings.HasPrefix(p, r.Path) {
			return r.handler
		}
	}
	log.Println("returning nil")
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
				req.URL.Scheme = "https"
				req.URL.Host = h
				req.URL.Path = "/boo1" // TODO JvD - regex to change path here
			},
		}
	}
	return nil
}
