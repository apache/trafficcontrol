package tmagent

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
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

type ParentHealthServer struct {
	parentInfo *ParentInfo
	httpServer *http.Server
}

// StartParentHealthService is a helper function that calls NewParentHealthServer, runs ListenAndServe in a goroutine, and returns the ParentHealthServer.
// If ListenAndServe returns an error, it's logged to the error log.
func StartParentHealthService(pi *ParentInfo, port int, readTimeout time.Duration) *ParentHealthServer {
	sv := NewParentHealthServer(pi, ":"+strconv.Itoa(port), readTimeout)
	go func() {
		if err := sv.ListenAndServe(); err != nil {
			log.Errorln("Parent Health Service returned: " + err.Error())
		}
	}()
	return sv
}

func NewParentHealthServer(pi *ParentInfo, addr string, readTimeout time.Duration) *ParentHealthServer {
	sv := &ParentHealthServer{
		parentInfo: pi,
		httpServer: &http.Server{
			Addr:        addr,
			ReadTimeout: readTimeout,
		},
	}
	sv.httpServer.Handler = sv
	return sv
}

func (sv *ParentHealthServer) ListenAndServe() error {
	return sv.httpServer.ListenAndServe()
}

func (sv *ParentHealthServer) ListenAndServeTLS(certFile, keyFile string) error {
	return sv.httpServer.ListenAndServeTLS(certFile, keyFile)
}

func (sv *ParentHealthServer) Shutdown(ctx context.Context) error {
	if sv.httpServer == nil {
		return nil // allow Shutdown on a default-initialized ParentHealthServer as a no-op.
	}
	return sv.httpServer.Shutdown(ctx)
}

func (sv *ParentHealthServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO implement io.Reader/+io.WriterTo?
	// TODO add 1s cache

	parentHealthL4 := sv.parentInfo.ParentHealthL4.Get()
	parentHealthL7 := sv.parentInfo.ParentHealthL7.Get()
	parentSvcHealth := sv.parentInfo.ParentServiceHealth.Get()
	combinedHealth := CombineParentHealth(parentHealthL4, parentHealthL7, parentSvcHealth)

	parentHealthBts, err := combinedHealth.Serialize(false)
	if err != nil {
		const code = http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		return
	}

	w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
	w.Write(parentHealthBts)

	// Append a newline so the response is a valid POSIX file.
	// Not strictly necessary, but may be useful to clients, and it's still valid JSON,
	// there's no reason not to.
	w.Write([]byte("\n"))
}
