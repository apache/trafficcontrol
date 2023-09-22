package tocookie

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
	"fmt"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
)

func TestParse(t *testing.T) {
	// Initialize the authentication data.
	authData := "foobar"

	// Generate a valid HTTP cookie.
	httpCookie := GetCookie(authData, 0, "fOObAR.")
	validCookie := http.Cookie{Name: httpCookie.Name, Value: httpCookie.Value}

	// Create an empty string to simulate an empty cookie value.
	emptyCookie := ""

	// Define a function that extracts the cookie value from an HTTP request.
	cookieValue := func(r *http.Request) string {
		// Attempt to extract the cookie value from the request.
		cookie, err := r.Cookie(Name)
		if err == nil && cookie != nil {
			return cookie.Value
		} else {
			return ""
		}
	}

	// Define a function that attempts to parse a cookie.
	parseCookie := func(cookieToken string) error {
		// Define the secret used to encrypt the cookie.
		secret := "fOObAR."

		// Initialize a string to hold the parsed cookie data.
		cookieData := ""

		// Attempt to parse the cookie.
		cookie, userErr, sysErr := Parse(secret, cookieToken)

		// Check if the cookie is nil.
		if cookie == nil {
			return fmt.Errorf("error: cookie data is nil")
		}

		// Extract the cookie data from the cookie object and compare it to the expected value.
		cookieData = cookie.AuthData
		if cookie != nil && userErr == nil && sysErr == nil {
			if cookieData != "foobar" {
				return fmt.Errorf("error: unable to parse cookie. expected: %v Got: %v", authData, cookieData)
			}
		}

		// Return nil to indicate success.
		return nil
	}

	// Create a new HTTP request.
	r, err := http.NewRequest("GET", "https://localhost:8888", nil)

	// Check if the request was created successfully.
	if err == nil && r != nil {

		// Test for a valid cookie.
		r.AddCookie(&validCookie)
		cookieToken := cookieValue(r)
		errValidCookie := parseCookie(cookieToken)
		assert.NoError(t, errValidCookie, "Error parsing valid cookie: %v", errValidCookie)

		// Test for an empty cookie.
		errEmptyCookie := parseCookie(emptyCookie)
		assert.Error(t, errEmptyCookie, "Expected error from parsing empty cookie, but got none")
	}
}
