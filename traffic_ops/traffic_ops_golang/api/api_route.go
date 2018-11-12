package api

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

import "bytes"
import "net/http"
import "regexp"
import "strconv"
//import trops "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang"
import log "github.com/apache/trafficcontrol/lib/go-log"

// todo - make the version string modular and accessible everywhere
const ServerString = "Traffic Operations/3.0.0";
const errorString = "Check the Traffic Ops log file(s) for details\n";

var APIVersions []float64;


// CompiledRoute ...
type CompiledRoute struct {
	Handler http.HandlerFunc
	Regex   *regexp.Regexp
	Params  []string
}


var AllRoutes *map[string][]CompiledRoute;
// Writes a message indicating an internal server error back to the client (in plain text)
func errorResponse(writer http.ResponseWriter) {
	err := []byte(errorString);
	writer.Header().Set("Content-Length", strconv.Itoa(len(err)));
	writer.WriteHeader(http.StatusInternalServerError);
	writer.Write(err);
}

// Handles the disallowed request methods for this endpoint
func AvailableRoutesBadMethodHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Allow", "OPTIONS");
	writer.Header().Set("Server", ServerString);
	writer.WriteHeader(http.StatusMethodNotAllowed);
}

func AvailableVersionsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Server", ServerString);

	// check that this has been set
	if APIVersions == nil {
		writer.Header().Set("Content-Length", "0");
		writer.WriteHeader(http.StatusNoContent);
		return;
	}

	var body bytes.Buffer;
	contentLength := 0;
	for _, version := range APIVersions {
		body.WriteString(strconv.FormatFloat(version, 'f', 2, 64));
		body.WriteRune('\n');
		contentLength += 3;
	}

	writer.Header().Set("Content-Length", strconv.Itoa(contentLength));
	writer.WriteHeader(http.StatusOK);
	writer.Write(body.Bytes());
}

// Writes a list of all available API routes in a response to the client in plaintext
// (or writes an error message via errorResponse if something wicked happens)
func AvailableRoutesHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Server", ServerString);

	// Check for a CORS preflight request (all headers accepted, so immediate response given)
	if request.Header.Get("Origin") != "" &&
	   request.Header.Get("Access-Control-Request-Method") != "" &&
	   request.Header.Get("Access-Control-Request-Headers") != "" {

		writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS");
		writer.Header().Set("Access-Control-Allow-Headers", "*");
		writer.WriteHeader(http.StatusNoContent);
		return;
	}

	// ... otherwise, return a list of all supported methods/routes

	// if for some reason the variable didn't get set properly, return nothing
	if AllRoutes == nil {
		writer.Header().Set("Content-Length", "0");
		writer.Header().Set("Allow", "OPTIONS");
		writer.WriteHeader(http.StatusNoContent);
		log.Warnln("API routes were requested, but weren't set!");
		return;
	}

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8");

	var body bytes.Buffer;
	contentLength := 0;
	for method, routes := range *AllRoutes {
		for _, route := range routes {
			n, err := body.WriteString(method);
			if err != nil {
				log.Errorf("Unable to append method to routes buffer: %s", err.Error());
				errorResponse(writer);
				return;
			}
			body.WriteRune(' ');
			contentLength += n + 1;

			n, err = body.WriteString(route.Regex.String());
			if err != nil {
				log.Errorf("Unable to append route to routes buffer: %s", err.Error());
				errorResponse(writer);
				return;
			}
			body.WriteRune('\n');
			contentLength += n +1;
		}
	}
	writer.Header().Set("Allow", "OPTIONS");
	writer.Header().Set("Content-Length", strconv.Itoa(contentLength));
	writer.WriteHeader(http.StatusOK);
	writer.Write(body.Bytes());
}
