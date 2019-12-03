package iso

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
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

// crypt acts Perl's built-in crypt() function, which in turn
// acts like the crypt(3) function in the C library.
// Since there isn't support in the standard library for this, the choice
// was made to shell out to `openssl`. This creates a dependency on openssl,
// which the Traffic Ops RPM already includes.
//
// Alternatives would be to use CGO to invoke the C function or
// implementing the algorithm in Go.
func crypt(password, salt string) (string, error) {
	var out strings.Builder

	// Special case to match Perl's behavior. In practice,
	// it's not expected to receive an empty password.
	if password == "" {
		password = "\n"
	}

	// https://www.openssl.org/docs/man1.1.0/man1/openssl-passwd.html
	cmd := exec.Command("openssl", "passwd", "-1", "-salt", salt, "-stdin")
	cmd.Stdin = strings.NewReader(password)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSuffix(out.String(), "\n"), nil
}

// saltChars are the possible characters rndSalt may use to generate a salt string.
const saltChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// rndSalt creates a random sequence of characters of given length.
// Suitable for use as the salt parameter with the crypt function.
func rndSalt(length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var out strings.Builder
	out.Grow(length)

	for i := 0; i < length; i++ {
		out.WriteRune(
			rune(saltChars[rng.Intn(len(saltChars))]),
		)
	}

	return out.String()
}
