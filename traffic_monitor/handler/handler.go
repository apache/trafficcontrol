package handler

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
	"io"
	"time"
)

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

// OpsConfig holds configuration for a Traffic Monitor relating to its
// connections with Traffic Ops **and** settings for its API/web UI server.
type OpsConfig struct {
	// The name of the CDN to which this Traffic Monitor belongs.
	CdnName string `json:"cdnName"`
	// The path to an SSL certificate to use with KeyFile to provide HTTP
	// encryption for the TM API and web UI.
	CertFile string `json:"certFile"`
	// The address on which to listen for HTTP requests.
	HttpListener string `json:"httpListener"`
	// The address on which to listen for HTTPS requests. If not set, TM serves
	// its API and UI over HTTP. If this is set, the HTTP server is only used to
	// redirect traffic to HTTPS.
	HttpsListener string `json:"httpsListener"`
	// Controls whether to validate the HTTPS certificate prevented by the
	// Traffic Ops server.
	Insecure bool `json:"insecure"`
	// The path to an SSL key to use with CertFile to provide HTTP encryption
	// for the TM API and web UI.
	KeyFile string `json:"keyFile"`
	// The password of the user identified by Username.
	Password string `json:"password"`
	// The URL at which Traffic Ops may be reached.
	Url string `json:"url"`
	// The username of the user as whom to authenticate with Traffic Ops.
	Username string `json:"username"`
	// Only used in the TM UI to indicate if TM started up with on-disk backup
	// Snapshots.
	UsingDummyTO bool `json:"usingDummyTO"`
}

type Handler interface {
	Handle(string, io.Reader, string, time.Duration, time.Time, error, uint64, bool, interface{}, chan<- uint64)
}
