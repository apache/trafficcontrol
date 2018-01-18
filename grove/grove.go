package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
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
	"github.com/apache/incubator-trafficcontrol/grove/config"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/hashicorp/golang-lru"
)

const ShutdownTimeout = 60 * time.Second

func main() {
	runtime.GOMAXPROCS(32) // DEBUG
	configFileName := flag.String("cfg", "", "The config file path")
	pprof := flag.Bool("pprof", false, "Whether to profile")
	flag.Parse()

	if *configFileName == "" {
		fmt.Printf("Error starting service: The -cfg argument is required\n")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(*configFileName)
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
		log.Errorf("starting service: creating cache: %v\n", err)
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

	stats := cache.NewStats(remapper.Rules(), lruCache, uint64(cfg.CacheSizeBytes))

	buildHandler := func(scheme string, port string, conns *web.ConnMap) (http.Handler, *cache.HandlerPointer) {
		statHandler := cache.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats, remapper.StatRules())
		cacheHandler := cache.NewHandler(
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
		cacheHandlerPointer := cache.NewHandlerPointer(cacheHandler)

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
		cfg, err = config.LoadConfig(*configFileName)
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
				log.Errorf("reloading config: creating cache: %v\n", err)
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

		stats = cache.NewStats(remapper.Rules(), lruCache, uint64(cfg.CacheSizeBytes)) // TODO copy stats from old stats object?

		httpCacheHandler := cache.NewHandler(
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

		httpsCacheHandler := cache.NewHandler(
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

			ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				if err == context.DeadlineExceeded {
					log.Errorf("closing http server: connections didn't close gracefully in %v, forcefully closing.\n", ShutdownTimeout)
					httpServer.Close()
				} else {
					log.Errorf("closing http server: %v\n", err)
				}

			}

			httpServer = startServer(handler, httpListener, httpConnStateCallback, cfg.Port, idleTimeout, readTimeout, writeTimeout, "http")
		}

		if (httpsServer == nil || cfg.HTTPSPort != oldCfg.HTTPSPort) && cfg.CertFile != "" && cfg.KeyFile != "" {
			statHandler := cache.NewStatHandler(cfg.InterfaceName, remapper.Rules(), stats, remapper.StatRules())
			handler := http.NewServeMux()
			handler.Handle("/_astats", statHandler)
			handler.Handle("/", httpsHandlerPointer)

			if httpsServer != nil {
				ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
				defer cancel()
				if err := httpsServer.Shutdown(ctx); err != nil {
					if err == context.DeadlineExceeded {
						log.Errorf("closing https server: connections didn't close gracefully in %v, forcefully closing.\n", ShutdownTimeout)
						httpServer.Close()
					} else {
						log.Errorf("closing https server: %v\n", err)
					}
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

func loadCerts(rules []cache.RemapRule) ([]tls.Certificate, error) {
	certs := []tls.Certificate{}
	for _, rule := range rules {
		if rule.CertificateFile == "" && rule.CertificateKeyFile == "" {
			continue
		}
		if rule.CertificateFile == "" {
			return nil, errors.New("rule " + rule.Name + " has a certificate but no key, using default certificate\n")
		}
		if rule.CertificateKeyFile == "" {
			return nil, errors.New("rule " + rule.Name + " has a key but no certificate, using default certificate\n")
		}

		cert, err := tls.LoadX509KeyPair(rule.CertificateFile, rule.CertificateKeyFile)
		if err != nil {
			return nil, errors.New("loading rule " + rule.Name + " certificate: " + err.Error() + "\n")
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
