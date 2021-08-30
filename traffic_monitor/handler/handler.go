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

type OpsConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Url           string `json:"url"`
	Insecure      bool   `json:"insecure"`
	CdnName       string `json:"cdnName"`
	HttpListener  string `json:"httpListener"`
	HttpsListener string `json:"httpsListener"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
	UsingDummyTO  bool   `json:"usingDummyTO"` // only used in the TM UI to indicate if TM started up with on-disk backup snapshots
}

type Handler interface {
	Handle(string, io.Reader, string, time.Duration, time.Time, error, uint64, bool, interface{}, chan<- uint64)
}
