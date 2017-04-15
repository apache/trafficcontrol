package grove

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type remapHandler struct {
	parent   http.Handler
	remapper HTTPRequestRemapper
}

// NewHandler returns an http.Handler objectn, which may be pipelined with other http.Handlers via `http.ListenAndServe`. If you prefer pipelining functions, use `GetHandlerFunc`.
func NewRemapHandler(parent http.Handler, remapper HTTPRequestRemapper) http.Handler {
	return &remapHandler{
		parent:   parent,
		remapper: remapper,
	}
}

// NewHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use Java. I mean, `New`.
func NewRemapHandlerFunc(parent http.HandlerFunc, remapper HTTPRequestRemapper) http.HandlerFunc {
	handler := NewRemapHandler(parent, remapper)
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func (h *remapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r, ok := h.remapper.Remap(r)
	// TODO configurable remap failure response
	if !ok {
		code := http.StatusNotFound
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		return
	}
	h.parent.ServeHTTP(w, r)
}

type HTTPRequestRemapper interface {
	Remap(*http.Request) (*http.Request, bool)
}

type simpleHttpRequestRemapper struct {
	remapper Remapper
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func (hr simpleHttpRequestRemapper) Remap(r *http.Request) (*http.Request, bool) {
	// NewRequest(method, urlStr string, body io.Reader)
	// TODO config whether to consider query string, method, headers
	oldUri := fmt.Sprintf("%s://%s%s", getScheme(r), r.Host, r.RequestURI)
	fmt.Printf("DEBUG Remap oldUri: '%v'\n", oldUri)
	fmt.Printf("DEBUG request: '%+v'\n", r)
	newUri, ok := hr.remapper.Remap(oldUri)
	if !ok {
		fmt.Printf("DEBUG Remap oldUri: '%v' NOT FOUND\n", oldUri)
		return r, false
	}
	fmt.Printf("DEBUG Remap newURI: '%v'\n", newUri)

	newReq, err := http.NewRequest(r.Method, newUri, nil) // TODO modify given req in-place?
	if err != nil {
		fmt.Printf("Error Remap NewRequest: %v\n", err)
		return r, false
	}
	copyHeader(r.Header, &newReq.Header)
	return newReq, true
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

func NewHTTPRequestRemapper(remap map[string]string) HTTPRequestRemapper {
	return RemapperToHTTP(NewLiteralPrefixRemapper(remap))
}

// Remapper provides a function which takes strings and maps them to other strings. This is designed for URL prefix remapping, for a reverse proxy.
type Remapper interface {
	// Remap returns the given string remapped, and whether a remap rule was found
	Remap(string) (string, bool)
}

// TODO change to use a prefix tree, for speed
type literalPrefixRemapper struct {
	remap map[string]string
}

func (r literalPrefixRemapper) Remap(s string) (string, bool) {
	for from, to := range r.remap {
		if strings.HasPrefix(s, from) {
			return to + s[len(from):], true
		}
	}
	return s, false
}

func NewLiteralPrefixRemapper(remap map[string]string) Remapper {
	return literalPrefixRemapper{remap: remap}
}

func LoadRemapRules(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	remap := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) < 3 {
			return nil, fmt.Errorf("malformed line '%s'", line)
		}
		rule := tokens[0]
		switch rule {
		case "map":
			from := tokens[1]
			to := tokens[2]
			remap[from] = to
		default:
			return nil, fmt.Errorf("unknown rule '%s'", line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return remap, nil
}

func LoadRemapper(path string) (HTTPRequestRemapper, error) {
	remapRules, err := LoadRemapRules(path)
	if err != nil {
		return nil, err
	}
	return NewHTTPRequestRemapper(remapRules), nil
}
