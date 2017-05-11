package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/hashicorp/golang-lru"
	"github.com/apache/incubator-trafficcontrol/grove"
)

type Config struct {
	// RFCCompliant determines whether `Cache-Control: no-cache` requests are honored. The ability to ignore `no-cache` is necessary to protect origin servers from DDOS attacks. In general, CDNs and caching proxies with the goal of origin protection should set RFCComplaint false. Cache with other goals (performance, load balancing, etc) should set RFCCompliant true.
	RFCCompliant bool `json:"rfc_compliant"`
	// Port is the HTTP port to serve on
	Port      int `json:"port"`
	HTTPSPort int `json:"https_port"`
	// CacheSizeBytes is the size of the memory cache, in bytes.
	CacheSizeBytes int    `json:"cache_size_bytes"`
	RemapRulesFile string `json:"remap_rules_file"`
	// ConcurrentRuleRequests is the number of concurrent requests permitted to a remap rule, that is, to an origin. Note this is overridden by any per-rule settings in the remap rules.
	ConcurrentRuleRequests int    `json:"concurrent_rule_requests"`
	CertFile               string `json:"cert_file"`
	KeyFile                string `json:"key_file"`
	InterfaceName          string `json:"interface_name"`
}

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	RFCCompliant:           true,
	Port:                   80,
	HTTPSPort:              443,
	CacheSizeBytes:         bytesPerGibibyte,
	RemapRulesFile:         "remap.config",
	ConcurrentRuleRequests: 100000,
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

	cache, err := lru.NewStrLargeWithEvict(uint64(cfg.CacheSizeBytes), nil)
	if err != nil {
		fmt.Printf("Error starting service: creating cache: %v\n")
		os.Exit(1)
	}

	remapper, err := grove.LoadRemapper(cfg.RemapRulesFile)
	if err != nil {
		fmt.Printf("Error starting service: loading remap rules: %v\n", err)
		os.Exit(1)
	}

	httpListener, httpConns, err := grove.InterceptListen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		fmt.Printf("Error creating HTTP listener %v: %v\n", cfg.Port, err)
		os.Exit(1)
	}

	httpsListener, httpsConns, err := grove.InterceptListenTLS("tcp", fmt.Sprintf(":%d", cfg.HTTPSPort), cfg.CertFile, cfg.KeyFile)
	if err != nil {
		fmt.Printf("Error creating HTTPS listener %v: %v\n", cfg.HTTPSPort, err)
	}

	stats := grove.NewStats(remapper.Rules())

	buildHandler := func(scheme string, conns *grove.ConnMap) (http.Handler, *grove.CacheHandlerPointer) {
		statHandler := grove.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats)
		cacheHandler := grove.NewCacheHandler(cache, remapper, uint64(cfg.ConcurrentRuleRequests), stats, scheme, conns, cfg.RFCCompliant)
		cacheHandlerPointer := grove.NewCacheHandlerPointer(cacheHandler)

		handler := http.NewServeMux()
		handler.Handle("/_astats", statHandler)
		handler.Handle("/", cacheHandlerPointer)
		return handler, cacheHandlerPointer
	}

	httpHandler, httpHandlerPointer := buildHandler("http", httpConns)
	httpsHandler, httpsHandlerPointer := buildHandler("https", httpsConns)

	// TODO add config to not serve HTTP (only HTTPS). If port is not set?
	httpServer := startHTTPServer(httpHandler, httpListener, cfg.Port)
	httpsServer := (*http.Server)(nil)
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		httpsServer = startHTTPSServer(httpsHandler, httpsListener, cfg.HTTPSPort, cfg.CertFile, cfg.KeyFile)
	}

	reloadConfig := func() {
		fmt.Printf("INFO reloading config\n")
		err := error(nil)
		oldCfg := cfg
		cfg, err = LoadConfig(*configFileName)
		if err != nil {
			fmt.Printf("Error reloading config: loading config file: %v\n", err)
			return
		}

		if cfg.CacheSizeBytes != oldCfg.CacheSizeBytes {
			// TODO determine if it's ok for the cache to temporarily exceed the value. This means the cache usage could be temporarily double, as old requestors still have the old object. We could call `Purge` on the old cache, to empty it, to mitigate this.
			cache, err = lru.NewStrLargeWithEvict(uint64(cfg.CacheSizeBytes), nil)
			if err != nil {
				fmt.Printf("Error reloading config: creating cache: %v\n")
				return
			}
		}

		remapper, err = grove.LoadRemapper(cfg.RemapRulesFile)
		if err != nil {
			fmt.Printf("Error starting service: loading remap rules: %v\n", err)
			os.Exit(1)
		}

		if cfg.Port != oldCfg.Port {
			if httpListener, httpConns, err = grove.InterceptListen("tcp", fmt.Sprintf(":%d", cfg.Port)); err != nil {
				fmt.Printf("Error reloading config: creating HTTP listener %v: %v\n", cfg.Port, err)
				return
			}
		}

		if (cfg.CertFile != oldCfg.CertFile || cfg.KeyFile != oldCfg.KeyFile) && cfg.HTTPSPort != oldCfg.HTTPSPort {
			fmt.Printf("WARNING config certificate changed, but port did not. Cannot recreate listener on same port without stopping the service. Restart the service to load the new certificate.\n")
		}

		if cfg.HTTPSPort != oldCfg.HTTPSPort {
			if httpsListener, httpsConns, err = grove.InterceptListenTLS("tcp", fmt.Sprintf(":%d", cfg.HTTPSPort), cfg.CertFile, cfg.KeyFile); err != nil {
				fmt.Printf("Error creating HTTPS listener %v: %v\n", cfg.HTTPSPort, err)
			}
		}

		stats = grove.NewStats(remapper.Rules()) // TODO copy stats from old stats object?

		httpCacheHandler := grove.NewCacheHandler(cache, remapper, uint64(cfg.ConcurrentRuleRequests), stats, "http", httpConns, cfg.RFCCompliant)
		httpHandlerPointer.Set(httpCacheHandler)

		httpsCacheHandler := grove.NewCacheHandler(cache, remapper, uint64(cfg.ConcurrentRuleRequests), stats, "https", httpsConns, cfg.RFCCompliant)
		httpsHandlerPointer.Set(httpsCacheHandler)

		if cfg.Port != oldCfg.Port {
			statHandler := grove.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats)
			handler := http.NewServeMux()
			handler.Handle("/_astats", statHandler)
			handler.Handle("/", httpHandlerPointer)
			if err := httpServer.Close(); err != nil {
				fmt.Printf("ERROR closing http server: %v\n", err)
			}
			httpServer = startHTTPServer(handler, httpListener, cfg.Port)
		}

		if (httpsServer == nil || cfg.HTTPSPort != oldCfg.HTTPSPort) && cfg.CertFile != "" && cfg.KeyFile != "" {
			statHandler := grove.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats)
			handler := http.NewServeMux()
			handler.Handle("/_astats", statHandler)
			handler.Handle("/", httpsHandlerPointer)
			if httpsServer != nil {
				if err := httpsServer.Close(); err != nil {
					fmt.Printf("ERROR closing https server: %v\n", err)
				}
			}
			httpsServer = startHTTPSServer(handler, httpsListener, cfg.HTTPSPort, cfg.CertFile, cfg.KeyFile)
		}
	}

	signalReloader(unix.SIGHUP, reloadConfig)
}

func signalReloader(sig os.Signal, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig)
	for range c {
		f()
	}
}

// startHTTPServer starts an HTTP server on the given port, and returns it
func startHTTPServer(handler http.Handler, listener net.Listener, port int) *http.Server {
	server := &http.Server{Handler: handler, Addr: fmt.Sprintf(":%d", port)}
	go func() {
		fmt.Printf("listening on http://%d\n", port)
		if err := server.Serve(listener); err != nil {
			fmt.Printf("Error serving HTTP port %v: %v\n", port, err)
		}
	}()
	return server
}

func startHTTPSServer(handler http.Handler, listener net.Listener, port int, certFile string, keyFile string) *http.Server {
	server := &http.Server{Handler: handler, Addr: fmt.Sprintf(":%d", port)}
	go func() {
		fmt.Printf("listening on https://%d\n", port)
		if err := server.Serve(listener); err != nil {
			fmt.Printf("Error serving HTTPS port %v: %v\n", port, err)
		}
	}()
	return server
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
