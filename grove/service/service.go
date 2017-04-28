package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/golang-lru"
	"github.com/apache/incubator-trafficcontrol/grove"
)

type Config struct {
	// RFCCompliant determines whether `Cache-Control: no-cache` requests are honored. The ability to ignore `no-cache` is necessary to protect origin servers from DDOS attacks.
	RFCCompliant bool `json:"rfc_compliant"`
	// Port is the HTTP port to serve on
	Port int `json:"port"`
	// CacheSizeBytes is the size of the memory cache, in bytes.
	CacheSizeBytes int    `json:"cache_size_bytes"`
	RemapRulesFile string `json:"remap_rules_file"`
}

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	RFCCompliant:   true,
	Port:           80,
	CacheSizeBytes: bytesPerGibibyte,
	RemapRulesFile: "remap.config",
}

// Load loads the given config file. If an empty string is passed, the default config is returned.
func LoadConfig(fileName string) (Config, error) {
	cfg := DefaultConfig
	if fileName == "" {
		return cfg, nil
	}
	configBytes, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = json.Unmarshal(configBytes, &cfg)
	}
	return cfg, err
}

const bytesPerGibibyte = 1024 * 1024 * 1024

func main() {
	configFileName := flag.String("config", "", "The config file path")
	flag.Parse()

	if *configFileName == "" {
		fmt.Printf("Error starting service: The --config argument is required\n")
		os.Exit(1)
	}

	cfg, err := LoadConfig(*configFileName)
	if err != nil {
		fmt.Printf("Error starting service: loading config: %v\n", err)
		os.Exit(1)
	}

	cache, err := lru.NewLargeWithEvict(uint64(cfg.CacheSizeBytes), nil)
	if err != nil {
		fmt.Printf("Error starting service: creating cache: %v\n")
		os.Exit(1)
	}

	remapper, err := grove.LoadRemapper(cfg.RemapRulesFile)
	if err != nil {
		fmt.Printf("Error starting service: loading remap rules: %v\n", err)
		os.Exit(1)
	}

	handler := grove.NewCacheHandler(cache, remapper)

	listen := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("proxy listening on http://%s\n", listen)
	if err := http.ListenAndServe(listen, handler); err != nil {
		fmt.Printf("Error serving: %v\n", err)
	}
}

// handle makes the given request and writes it to the given writer. It's assumed the request coming from a client has had its host rewritten to some other service. DO NOT call this with an unmodified request from a client; that would cause an infinite loop of pain.
func handle(w http.ResponseWriter, r *http.Request) {
	rr := r

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
