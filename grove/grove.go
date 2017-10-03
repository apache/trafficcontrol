package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"github.com/apache/incubator-trafficcontrol/grove/cache"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/hashicorp/golang-lru"
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
	// ConnectionClose determines whether to send a `Connection: close` header. This is primarily designed for maintenance, to drain the cache of incoming requestors. This overrides rule-specific `connection-close: false` configuration, under the assumption that draining a cache is a temporary maintenance operation, and if connectionClose is true on the service and false on some rules, those rules' configuration is probably a permament setting whereas the operator probably wants to drain all connections if the global setting is true. If it's necessary to leave connection close false on some rules, set all other rules' connectionClose to true and leave the global connectionClose unset.
	ConnectionClose bool `json:"connection_close"`

	LogLocationError   string `json:"log_location_error"`
	LogLocationWarning string `json:"log_location_warning"`
	LogLocationInfo    string `json:"log_location_info"`
	LogLocationDebug   string `json:"log_location_debug"`
	LogLocationEvent   string `json:"log_location_event"`

	ReqTimeoutMS         int `json:"parent_request_timeout_ms"` // TODO rename "parent_request" to distinguish from client requests
	ReqKeepAliveMS       int `json:"parent_request_keep_alive_ms"`
	ReqMaxIdleConns      int `json:"parent_request_max_idle_connections"`
	ReqIdleConnTimeoutMS int `json:"parent_request_idle_connection_timeout_ms"`

	ServerIdleTimeoutMS  int `json:"server_idle_timeout_ms"`
	ServerWriteTimeoutMS int `json:"server_write_timeout_ms"`
	ServerReadTimeoutMS  int `json:"server_read_timeout_ms"`
}

func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.LogLocationError)
}
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.LogLocationWarning)
}
func (c Config) InfoLog() log.LogLocation {
	return log.LogLocation(c.LogLocationInfo)
}
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.LogLocationDebug)
}
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.LogLocationEvent)
}

const MSPerSec = 1000

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	RFCCompliant:           true,
	Port:                   80,
	HTTPSPort:              443,
	CacheSizeBytes:         bytesPerGibibyte,
	RemapRulesFile:         "remap.config",
	ConcurrentRuleRequests: 100000,
	ConnectionClose:        false,
	LogLocationError:       log.LogLocationStderr,
	LogLocationWarning:     log.LogLocationStdout,
	LogLocationInfo:        log.LogLocationNull,
	LogLocationDebug:       log.LogLocationNull,
	LogLocationEvent:       log.LogLocationStdout,
	ReqTimeoutMS:           30 * MSPerSec,
	ReqKeepAliveMS:         30 * MSPerSec,
	ReqMaxIdleConns:        100,
	ReqIdleConnTimeoutMS:   90 * MSPerSec,
	ServerIdleTimeoutMS:    10 * MSPerSec,
	ServerWriteTimeoutMS:   3 * MSPerSec,
	ServerReadTimeoutMS:    3 * MSPerSec,
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
	runtime.GOMAXPROCS(32) // DEBUG
	configFileName := flag.String("cfg", "", "The config file path")
	pprof := flag.Bool("pprof", false, "Whether to profile")
	flag.Parse()

	if *configFileName == "" {
		fmt.Printf("Error starting service: The -cfg argument is required\n")
		os.Exit(1)
	}

	cfg, err := LoadConfig(*configFileName)
	if err != nil {
		fmt.Printf("Error starting service: loading config: %v\n", err)
		os.Exit(1)
	}

	eventW, errW, warnW, infoW, debugW, err := log.GetLogWriters(cfg)
	if err != nil {
		fmt.Printf("Error starting service: failed to create log writers: %v\n", err)
		os.Exit(1)
	}
	log.Init(eventW, errW, warnW, infoW, debugW)

	lruCache, err := lru.NewStrLargeWithEvict(uint64(cfg.CacheSizeBytes), nil)
	if err != nil {
		log.Errorf("starting service: creating cache: %v\n")
		os.Exit(1)
	}

	remapper, err := cache.LoadRemapper(cfg.RemapRulesFile)
	if err != nil {
		log.Errorf("starting service: loading remap rules: %v\n", err)
		os.Exit(1)
	}

	certs, err := loadCerts(remapper.Rules())
	if err != nil {
		log.Errorf("starting service: loading certificates: %v\n", err)
		os.Exit(1)
	}
	defaultCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Errorf("starting service: loading default certificate: %v\n", err)
		os.Exit(1)
	}
	certs = append(certs, defaultCert)

	httpListener, httpConns, httpConnStateCallback, err := web.InterceptListen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Errorf("creating HTTP listener %v: %v\n", cfg.Port, err)
		os.Exit(1)
	}

	stats := cache.NewStats(remapper.Rules())

	buildHandler := func(scheme string, port string, conns *web.ConnMap) (http.Handler, *cache.CacheHandlerPointer) {
		statHandler := cache.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats, remapper.StatRules())
		cacheHandler := cache.NewCacheHandler(
			lruCache,
			remapper,
			uint64(cfg.ConcurrentRuleRequests),
			stats,
			scheme,
			port,
			conns,
			cfg.RFCCompliant,
			cfg.ConnectionClose,
			time.Duration(cfg.ReqTimeoutMS)*time.Millisecond,
			time.Duration(cfg.ReqKeepAliveMS)*time.Millisecond,
			cfg.ReqMaxIdleConns,
			time.Duration(cfg.ReqIdleConnTimeoutMS)*time.Millisecond,
		)
		cacheHandlerPointer := cache.NewCacheHandlerPointer(cacheHandler)

		handler := http.NewServeMux()
		handler.Handle("/_astats", statHandler)
		handler.Handle("/", cacheHandlerPointer)
		return handler, cacheHandlerPointer
	}

	httpsServer := (*http.Server)(nil)
	httpsListener := net.Listener(nil)
	httpsConns := (*web.ConnMap)(nil)
	httpsConnStateCallback := (func(net.Conn, http.ConnState))(nil)
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		if httpsListener, httpsConns, httpsConnStateCallback, err = web.InterceptListenTLS("tcp", fmt.Sprintf(":%d", cfg.HTTPSPort), certs); err != nil {
			log.Errorf("creating HTTPS listener %v: %v\n", cfg.HTTPSPort, err)
			return
		}
	}

	httpHandler, httpHandlerPointer := buildHandler("http", strconv.Itoa(cfg.Port), httpConns)
	httpsHandler, httpsHandlerPointer := buildHandler("https", strconv.Itoa(cfg.HTTPSPort), httpsConns)

	idleTimeout := time.Duration(cfg.ServerIdleTimeoutMS) * time.Millisecond
	readTimeout := time.Duration(cfg.ServerReadTimeoutMS) * time.Millisecond
	writeTimeout := time.Duration(cfg.ServerWriteTimeoutMS) * time.Millisecond

	// TODO add config to not serve HTTP (only HTTPS). If port is not set?
	httpServer := startServer(httpHandler, httpListener, httpConnStateCallback, cfg.Port, idleTimeout, readTimeout, writeTimeout, "http")

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		httpsServer = startServer(httpsHandler, httpsListener, httpsConnStateCallback, cfg.HTTPSPort, idleTimeout, readTimeout, writeTimeout, "https")
	}

	reloadConfig := func() {
		log.Infof("reloading config\n")
		err := error(nil)
		oldCfg := cfg
		cfg, err = LoadConfig(*configFileName)
		if err != nil {
			log.Errorf("reloading config: loading config file: %v\n", err)
			return
		}
		eventW, errW, warnW, infoW, debugW, err := log.GetLogWriters(cfg)
		if err != nil {
			log.Errorf("relaoding config: getting log writers '%v': %v", *configFileName, err)
		}
		log.Init(eventW, errW, warnW, infoW, debugW)

		if cfg.CacheSizeBytes != oldCfg.CacheSizeBytes {
			// TODO determine if it's ok for the cache to temporarily exceed the value. This means the cache usage could be temporarily double, as old requestors still have the old object. We could call `Purge` on the old cache, to empty it, to mitigate this.
			lruCache, err = lru.NewStrLargeWithEvict(uint64(cfg.CacheSizeBytes), nil)
			if err != nil {
				log.Errorf("reloading config: creating cache: %v\n")
				return
			}
		}

		remapper, err = cache.LoadRemapper(cfg.RemapRulesFile)
		if err != nil {
			log.Errorf("starting service: loading remap rules: %v\n", err)
			os.Exit(1)
		}

		if cfg.Port != oldCfg.Port {
			if httpListener, httpConns, httpConnStateCallback, err = web.InterceptListen("tcp", fmt.Sprintf(":%d", cfg.Port)); err != nil {
				log.Errorf("reloading config: creating HTTP listener %v: %v\n", cfg.Port, err)
				return
			}
		}

		if (cfg.CertFile != oldCfg.CertFile || cfg.KeyFile != oldCfg.KeyFile) && cfg.HTTPSPort != oldCfg.HTTPSPort {
			log.Warnf("config certificate changed, but port did not. Cannot recreate listener on same port without stopping the service. Restart the service to load the new certificate.\n")
		}

		if cfg.HTTPSPort != oldCfg.HTTPSPort {
			if httpsListener, httpsConns, httpsConnStateCallback, err = web.InterceptListenTLS("tcp", fmt.Sprintf(":%d", cfg.HTTPSPort), certs); err != nil {
				log.Errorf("creating HTTPS listener %v: %v\n", cfg.HTTPSPort, err)
			}
		}

		stats = cache.NewStats(remapper.Rules()) // TODO copy stats from old stats object?

		httpCacheHandler := cache.NewCacheHandler(
			lruCache,
			remapper,
			uint64(cfg.ConcurrentRuleRequests),
			stats,
			"http",
			strconv.Itoa(cfg.Port),
			httpConns,
			cfg.RFCCompliant,
			cfg.ConnectionClose,
			time.Duration(cfg.ReqTimeoutMS)*time.Millisecond,
			time.Duration(cfg.ReqKeepAliveMS)*time.Millisecond,
			cfg.ReqMaxIdleConns,
			time.Duration(cfg.ReqIdleConnTimeoutMS)*time.Millisecond,
		)
		httpHandlerPointer.Set(httpCacheHandler)

		httpsCacheHandler := cache.NewCacheHandler(
			lruCache,
			remapper,
			uint64(cfg.ConcurrentRuleRequests),
			stats,
			"https",
			strconv.Itoa(cfg.HTTPSPort),
			httpsConns,
			cfg.RFCCompliant,
			cfg.ConnectionClose,
			time.Duration(cfg.ReqTimeoutMS)*time.Millisecond,
			time.Duration(cfg.ReqKeepAliveMS)*time.Millisecond,
			cfg.ReqMaxIdleConns,
			time.Duration(cfg.ReqIdleConnTimeoutMS)*time.Millisecond,
		)
		httpsHandlerPointer.Set(httpsCacheHandler)

		if cfg.Port != oldCfg.Port {
			statHandler := cache.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats, remapper.StatRules())
			handler := http.NewServeMux()
			handler.Handle("/_astats", statHandler)
			handler.Handle("/", httpHandlerPointer)
			if err := httpServer.Close(); err != nil {
				log.Errorf("closing http server: %v\n", err)
			}
			httpServer = startServer(handler, httpListener, httpConnStateCallback, cfg.Port, idleTimeout, readTimeout, writeTimeout, "http")
		}

		if (httpsServer == nil || cfg.HTTPSPort != oldCfg.HTTPSPort) && cfg.CertFile != "" && cfg.KeyFile != "" {
			statHandler := cache.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats, remapper.StatRules())
			handler := http.NewServeMux()
			handler.Handle("/_astats", statHandler)
			handler.Handle("/", httpsHandlerPointer)
			if httpsServer != nil {
				if err := httpsServer.Close(); err != nil {
					log.Errorf("closing https server: %v\n", err)
				}
			}
			httpsServer = startServer(handler, httpsListener, httpsConnStateCallback, cfg.HTTPSPort, idleTimeout, readTimeout, writeTimeout, "https")
		}
	}

	if *pprof {
		profile()
	}
	signalReloader(unix.SIGHUP, reloadConfig)
}

func profile() {
	go func() {
		count := 0
		for {
			count++
			filename := fmt.Sprintf("grove%d.pprof", count)
			f, err := os.Create(filename)
			if err != nil {
				log.Errorf("creating profile: %v\n", err)
				os.Exit(1)
			}
			pprof.StartCPUProfile(f)
			time.Sleep(time.Minute)
			pprof.StopCPUProfile()
			f.Close()
		}
	}()
}

func signalReloader(sig os.Signal, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig)
	for range c {
		f()
	}
}

// startServer starts an HTTP or HTTPS server on the given port, and returns it.
func startServer(handler http.Handler, listener net.Listener, connState func(net.Conn, http.ConnState), port int, idleTimeout time.Duration, readTimeout time.Duration, writeTimeout time.Duration, protocol string) *http.Server {
	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%d", port),
		ConnState:    connState,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	go func() {
		log.Infof("listening on %s://%d\n", protocol, port)
		if err := server.Serve(listener); err != nil {
			log.Errorf("serving %s port %v: %v\n", strings.ToUpper(protocol), port, err)
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

func loadCerts(rules []cache.RemapRule) ([]tls.Certificate, error) {
	certs := []tls.Certificate{}
	for _, rule := range rules {
		if rule.CertificateFile == "" && rule.CertificateKeyFile == "" {
			continue
		}
		if rule.CertificateFile == "" {
			return nil, errors.New("rule " + rule.Name + " has a certificate but no key, using default certificate\n")
			continue
		}
		if rule.CertificateKeyFile == "" {
			return nil, errors.New("rule " + rule.Name + " has a key but no certificate, using default certificate\n")
			continue
		}

		cert, err := tls.LoadX509KeyPair(rule.CertificateFile, rule.CertificateKeyFile)
		if err != nil {
			return nil, errors.New("loading rule " + rule.Name + " certificate: " + err.Error() + "\n")
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
