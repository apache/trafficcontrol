// Package routing defines the HTTP routes for Traffic Ops and provides tools to
// register those routes with appropriate middleware.
package routing

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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/plugin"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault"

	"github.com/jmoiron/sqlx"
)

// HandleBatchReq handles the request, if it's a batch request.
// If the request isn't a batch request, it does nothing.
// Returns whether the request was handled.
func HandleBatchReq(
	routes map[string][]CompiledRoute,
	versions map[api.Version]struct{},
	catchall http.Handler,
	db *sqlx.DB,
	cfg *config.Config,
	getReqID func() uint64,
	plugins plugin.Plugins,
	tv trafficvault.TrafficVault,
	w http.ResponseWriter,
	r *http.Request,
) bool {
	log.Errorln("DEBUG HandleBatchReq starting")

	if r.Method != http.MethodPost {
		log.Errorln("DEBUG HandleBatchReq not POST, returning unhandled")
		return false
	}

	pathParts := strings.Split(r.URL.Path, `/`)
	if len(pathParts) < 4 ||
		strings.ToLower(pathParts[1]) != "api" ||
		strings.ToLower(pathParts[3]) != "batch" {
		log.Errorln("DEBUG HandleBatchReq not path, returning unhandled")
		return false
	}

	apiVer, err := strconv.ParseFloat(pathParts[2], 64)
	if err != nil {
		log.Errorln("DEBUG HandleBatchReq bad api version, returning unhandled")
		// don't log or do anything here, let the normal handler deal with it
		return false
	}

	if apiVer < 5 {
		log.Errorln("DEBUG HandleBatchReq old api version, returning unhandled")
		return false // batching was introduced in API 5.0
	}

	log.Errorln("DEBUG HandleBatchReq handling")

	reqs := []BatchReq{}
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		err = errors.New("error decoding request, possibly malformed json")
		h2 := middleware.WrapAccessLog(cfg.Secrets[0], middleware.BackendErrorHandler(http.StatusBadRequest, err, nil))
		h2.ServeHTTP(w, r)
		return true
	}

	if len(reqs) == 0 {
		log.Errorln("DEBUG HandleBatchReq no reqs, returning handled error")
		// TODO this should really succeed, to make automation and ops easier
		err = errors.New("batch had no operations")
		h2 := middleware.WrapAccessLog(cfg.Secrets[0], middleware.BackendErrorHandler(http.StatusBadRequest, err, nil))
		h2.ServeHTTP(w, r)
		return true
	}

	tx, err := api.StartTx(db, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second, r.Context())
	if err != nil {
		h2 := middleware.WrapAccessLog(cfg.Secrets[0], middleware.BackendErrorHandler(http.StatusInternalServerError, nil, err))
		h2.ServeHTTP(w, r)
		return true
	}
	ctx := context.WithValue(r.Context(), api.TxContextKey, tx)
	r = r.WithContext(ctx)
	defer api.CloseTx(tx.Tx, tx.CancelF)

	responses := []BatchResp{}

	for _, req := range reqs {
		// Note sequential GETs could theoretically be parallelized, but
		// mutating methods (POST/PUT/DELETE/PATCH) MUST NOT be executed in parallel with themselves or GET requests.
		//
		// Likewise, if a batch includes a GET-POST-GET, the GETs cannot be parallelized, since the first may
		// depend on the POST not executing yet, and the second may depend on the POST having already executed.
		//
		fakeWriter := util.NewFullInterceptor()

		// copy the real request, and set the method, URL, headers, and body to the individual batch req values
		batchR := *r
		batchR.Method = req.Method

		url := *r.URL
		pathQuery := strings.SplitN(req.Path, `?`, 2)
		url.Path = pathQuery[0]
		if len(pathQuery) > 1 {
			url.RawQuery = pathQuery[1]
		}
		batchR.URL = &url

		// We overlay the batch req headers onto the global headers.
		// This should be documented.
		// If clients need to send an HTTP Header that isn't set in some batch request,
		// they should specifically set that header in the batch member to the empty string.
		batchR.Header = r.Header.Clone()
		for name, vals := range req.Headers {
			batchR.Header[name] = vals
		}

		batchR.Body = io.NopCloser(bytes.NewBuffer([]byte(req.Body)))

		Handler(routes, versions, catchall, db, cfg, getReqID, plugins, tv, fakeWriter, &batchR, true)
		responses = append(responses, BatchResp{
			Code:    fakeWriter.Code,
			Headers: fakeWriter.Header(),
			Body:    json.RawMessage(fakeWriter.Body),
		})
	}

	// Note we don't do any kind of auth, middleware, or anything.
	// Because batch ops are normal requests, they each do their own auth, middleware, headers, etc.

	respBts, err := json.Marshal(responses)
	if err != nil {
		h2 := middleware.WrapAccessLog(cfg.Secrets[0], middleware.BackendErrorHandler(http.StatusInternalServerError, nil, err))
		h2.ServeHTTP(w, r)
		return true
	}

	w.Write(respBts)
	return true
}

type BatchReq struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    json.RawMessage     `json:"body"`
}

type BatchResp struct {
	Code    int                 `json:"code"`
	Headers map[string][]string `json:"headers"`
	Body    json.RawMessage     `json:"body"`
}
