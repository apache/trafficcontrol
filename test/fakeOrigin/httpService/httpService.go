package httpService

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/dtp"
	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/endpoint"
	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/transcode"
)

// EndpointRoutes contains the paths of all HTTP routes for a given config "endpoint"
// TODO come up with a better structure. The static and abr checks are hackish, the structure should only contain the necessary variables for the given route. Maybe there's a better way with interfaces?
type EndpointRoutes struct {
	MasterPath    string
	VariantPaths  []string
	FragmentPaths []string
	MetaJSONPaths []string
	IsABR         bool
}

// GetRoutes returns the map of config IDs, to the full HTTP paths, e.g. for each m3u8 and ts created for that config object.
func GetRoutes(cfg endpoint.Config) (map[string]EndpointRoutes, error) {
	// TODO return both HTTP route and file path, so the handler doesn't have to rebuild the file path to load the file from disk
	allRoutes := map[string]EndpointRoutes{}
	for _, ep := range cfg.Endpoints {
		routes := EndpointRoutes{}
		if ep.EndpointType == endpoint.Static {
			routes.MasterPath = path.Join("/", ep.ID, filepath.Base(ep.Source))
			allRoutes[ep.ID] = routes
			continue
		}

		if ep.EndpointType == endpoint.Testing {
			routes.MasterPath = path.Join("/", ep.ID)
			allRoutes[ep.ID] = routes
			continue
		}

		if ep.EndpointType == endpoint.Dir {
			fileList := []string{}
			err := filepath.Walk(ep.Source, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// Skip directories since that's not a static file
				if !f.IsDir() {
					fileList = append(fileList, path)
				}
				return nil
			})
			if err != nil {
				return nil, errors.New("reading files '" + ep.Source + "': " + err.Error())
			}

			for _, file := range fileList {
				pathsuffix := strings.TrimPrefix(file, ep.Source)
				routes.VariantPaths = append(routes.VariantPaths, path.Join("/", ep.ID, pathsuffix))
			}
			allRoutes[ep.ID] = routes
			continue
		}

		files, err := ioutil.ReadDir(ep.OutputDirectory)
		if err != nil {
			return nil, errors.New("reading files '" + ep.Source + "': " + err.Error())
		}
		routes.IsABR = false
		for _, f := range files {
			if !f.IsDir() {
				if strings.ToLower(filepath.Ext(f.Name())) == ".ts" {
					routes.FragmentPaths = append(routes.FragmentPaths, path.Join("/", ep.ID, f.Name()))
				} else if f.Name() == ep.OutputDirectory+"/"+ep.DiskID+".m3u8" {
					routes.MasterPath = path.Join("/", ep.ID, f.Name())
				} else if strings.ToLower(filepath.Ext(f.Name())) == ".m3u8" {
					routes.IsABR = true
					routes.VariantPaths = append(routes.VariantPaths, path.Join("/", ep.ID, f.Name()))
				} else if strings.HasSuffix(strings.ToLower(f.Name()), ".meta.json") {
					routes.MetaJSONPaths = append(routes.MetaJSONPaths, path.Join("/", ep.ID, ep.DiskID+".meta.json"))
				}
			}
		}

		allRoutes[ep.ID] = routes
	}
	return allRoutes, nil
}

// PrintRoutes writes the routes across multiple lines in human-readable format to the given writer.
// The header is written at the beginning of each endpoint. The routePrefix is written at the beginning of each route.
// If showFragments is false, only manifests and meta files are printed, not video fragments.
func PrintRoutes(w io.Writer, routes map[string]EndpointRoutes, endpointPrefix, routePrefix string, showFragments bool) {
	s := ""
	for endpointID, endpointRoutes := range routes {
		s += endpointPrefix + endpointID + "\n"
		if endpointRoutes.MasterPath != "" {
			s += routePrefix + endpointRoutes.MasterPath + "\n"
		}
		for _, path := range endpointRoutes.VariantPaths {
			s += routePrefix + path + "\n"
		}
		for _, path := range endpointRoutes.MetaJSONPaths {
			s += routePrefix + path + "\n"
		}
		if showFragments {
			for _, path := range endpointRoutes.FragmentPaths {
				s += routePrefix + path + "\n"
			}
		}
	}
	w.Write([]byte(s))
}

const ContentTypeJSON = `application/json`
const ContentTypeM3U8 = `application/x-mpegURL`
const ContentTypeTS = `video/MP2T`

func registerRoute(mux *http.ServeMux, e endpoint.Endpoint, httpPath string, isLive bool, contentType string, isSSL bool) error {
	startTime := time.Now()
	if e.EndpointType == endpoint.Static || e.EndpointType == endpoint.Dir {
		diskpath := e.Source
		if e.EndpointType == endpoint.Dir {
			diskpath = path.Join("/", e.Source, strings.TrimPrefix(httpPath, "/"+e.ID+"/"))
		}
		h, err := staticHandler(staticHTTPEndpoint{
			FilePath: diskpath,
			NoCache:  e.NoCache,
			Headers:  e.DefaultHeaders,
		})
		if err != nil {
			return errors.New("creating handler '" + httpPath + "': " + err.Error())
		}
		if isSSL {
			mux.Handle(httpPath, logfo(strictTransportSecurity(originHeaderManipulation(h))))
			mux.Handle(httpPath+"/", logfo(strictTransportSecurity(originHeaderManipulation(h))))
		} else {
			mux.Handle(httpPath, logfo(originHeaderManipulation(h)))
			mux.Handle(httpPath+"/", logfo(originHeaderManipulation(h)))
		}
		fmt.Println("registered for static endpoint logfo for path: ", httpPath)
		return nil
	}

	if e.EndpointType == endpoint.Testing {
		alog := log.New(os.Stderr, "", 0)
		dtpHandler := dtp.NewDTPHandler()

		dtp.GlobalConfig.Log.RequestHeaders = e.LogReqHeaders
		dtp.GlobalConfig.Log.ResponseHeaders = e.LogRespHeaders
		dtp.GlobalConfig.StallDuration = e.StallDuration * time.Second
		dtp.GlobalConfig.Debug = e.EnableDebug
		dtp.GlobalConfig.EnablePprof = e.EnablePprof

		// General DTP endpoints for testing
		mux.Handle("/"+e.ID, dtp.Logger(alog, dtpHandler))
		mux.Handle("/"+e.ID+"/", dtp.Logger(alog, dtpHandler))

		// DTP endpoints for pprof output
		if dtp.GlobalConfig.EnablePprof {
			mux.HandleFunc("/"+e.ID+"/debug/pprof/", pprof.Index)
			mux.HandleFunc("/"+e.ID+"/debug/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/"+e.ID+"/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/"+e.ID+"/debug/pprof/symbol", pprof.Symbol)
			mux.HandleFunc("/"+e.ID+"/debug/pprof/trace", pprof.Trace)
		}

		// DTP endpoint for setting various config values on the fly
		mux.HandleFunc("/"+e.ID+"/config", dtp.ConfigHandler)
		return nil
	}

	var lastTranscodeDT time.Time
	if diskmeta, err := transcode.GetMeta(e); err != nil {
		lastTranscodeDT = startTime
	} else {
		var timeOk bool
		if lastTranscodeDT, timeOk = ParseHTTPDate(diskmeta.LastTranscodeDT); !timeOk {
			lastTranscodeDT = startTime
		}
	}
	ep := httpEndpoint{
		ID:                e.ID,
		OutputDir:         e.OutputDirectory,
		Type:              e.EndpointType,
		DiskID:            e.DiskID,
		IsLive:            isLive,
		DefaultHeaders:    e.DefaultHeaders,
		NoCache:           e.NoCache,
		HTTPPath:          httpPath,
		ContentType:       contentType,
		LastTranscodeTime: lastTranscodeDT,
	}
	h, err := httpEndpointHandler(ep)
	if err != nil {
		return errors.New("registering route '" + httpPath + "': " + err.Error())
	}
	if isSSL {
		mux.Handle(httpPath, logfo(strictTransportSecurity(originHeaderManipulation(cacheOptimization(h, startTime, ep)))))
		mux.Handle(httpPath+"/", logfo(strictTransportSecurity(originHeaderManipulation(cacheOptimization(h, startTime, ep)))))
	} else {
		mux.Handle(httpPath, logfo(originHeaderManipulation(cacheOptimization(h, startTime, ep))))
		mux.Handle(httpPath+"/", logfo(originHeaderManipulation(cacheOptimization(h, startTime, ep))))
	}
	return nil
}

func registerRoutes(mux *http.ServeMux, conf endpoint.Config, routes map[string]EndpointRoutes, isSSL bool) error {
	for _, e := range conf.Endpoints {
		endpointRoutes, ok := routes[e.ID]
		if !ok {
			return errors.New("no routes found for endpoint '" + e.ID + "'")
		}
		if e.EndpointType == endpoint.Testing {
			err := registerRoute(mux, e, e.ID, false, ContentTypeJSON, isSSL)
			if err != nil {
				return errors.New("error registering endpoint '" + e.ID + "': " + err.Error())
			}
		}
		if endpointRoutes.MasterPath != "" && e.EndpointType != endpoint.Testing {
			err := registerRoute(mux, e, endpointRoutes.MasterPath, e.EndpointType == endpoint.Live && !endpointRoutes.IsABR, ContentTypeM3U8, isSSL)
			if err != nil {
				return errors.New("Error registering endpoint '" + e.ID + "': " + err.Error())
			}
		}
		for _, path := range endpointRoutes.VariantPaths {
			err := registerRoute(mux, e, path, e.EndpointType == endpoint.Live && endpointRoutes.IsABR, ContentTypeM3U8, isSSL)
			if err != nil {
				return errors.New("Error registering endpoint '" + e.ID + "': " + err.Error())
			}
		}
		for _, path := range endpointRoutes.FragmentPaths {
			err := registerRoute(mux, e, path, false, ContentTypeTS, isSSL)
			if err != nil {
				return errors.New("Error registering endpoint '" + e.ID + "': " + err.Error())
			}
		}
		for _, path := range endpointRoutes.MetaJSONPaths {
			err := registerRoute(mux, e, path, false, ContentTypeJSON, isSSL)
			if err != nil {
				return errors.New("Error registering endpoint '" + e.ID + "': " + err.Error())
			}
		}
	}
	mux.Handle("/crossdomain.xml", logfo(crossdomainHandler(conf.ServerConf.CrossdomainFile)))
	mux.Handle("/", logfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})))
	return nil
}

// StartHTTPListener kicks off the HTTPS stack
func StartHTTPListener(conf endpoint.Config, routes map[string]EndpointRoutes) error {
	mux := http.NewServeMux()
	if err := registerRoutes(mux, conf, routes, false); err != nil {
		return errors.New("registering routes: " + err.Error())
	}
	fmt.Println("Serving HTTP on " + conf.ServerConf.BindingAddress + ":" + strconv.Itoa(conf.ServerConf.HTTPListeningPort))
	srv := &http.Server{
		Addr:         conf.ServerConf.BindingAddress + ":" + strconv.Itoa(conf.ServerConf.HTTPListeningPort),
		Handler:      mux,
		ReadTimeout:  conf.ServerConf.ReadTimeout * time.Second,
		WriteTimeout: conf.ServerConf.ReadTimeout * time.Second,
	}
	return srv.ListenAndServe()
}

// StartHTTPSListener kicks off the HTTPS stack
func StartHTTPSListener(conf endpoint.Config, routes map[string]EndpointRoutes) error {
	if err := assertSSLCerts(conf.ServerConf.SSLCert, conf.ServerConf.SSLKey); err != nil {
		return fmt.Errorf("asserting SSL info Cert:'%+v' Key:'%+v': %+v", conf.ServerConf.SSLCert, conf.ServerConf.SSLKey, err)
	}
	mux := http.NewServeMux()
	if err := registerRoutes(mux, conf, routes, true); err != nil {
		return errors.New("registering routes: " + err.Error())
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         conf.ServerConf.BindingAddress + ":" + strconv.Itoa(conf.ServerConf.HTTPSListeningPort),
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		ReadTimeout:  conf.ServerConf.ReadTimeout * time.Second,
		WriteTimeout: conf.ServerConf.WriteTimeout * time.Second,
	}
	fmt.Println("Serving HTTPS on " + conf.ServerConf.BindingAddress + ":" + strconv.Itoa(conf.ServerConf.HTTPSListeningPort))
	return srv.ListenAndServeTLS(conf.ServerConf.SSLCert, conf.ServerConf.SSLKey)
}
