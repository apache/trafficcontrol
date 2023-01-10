package dtp

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
	"encoding/base64"
	"net/http"
)

func init() {
	GlobalHandlerFuncs["hijack"] = Hijack
}

func Hijack(w http.ResponseWriter, r *http.Request, reqdat map[string]string) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		DebugLogf("Can't get a hijack\n")
		return
	}

	conn, buf, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		DebugLogf("hijack failed: %s\n", err)
		return
	}

	DebugLogf("Connection hijacked\n")

	if str, ok := reqdat["payload"]; ok {
		buf.WriteString(str)
	} else if encoded, ok := reqdat["payload64"]; ok {
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err == nil {
			buf.Write(data)
		}
	}
	buf.Flush()
	conn.Close()
}
