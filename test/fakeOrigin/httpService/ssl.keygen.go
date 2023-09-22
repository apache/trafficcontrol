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
	"errors"
	"os"

	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/transcode"
)

func assertSSLCerts(crtPath, keyPath string) error {
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		if err = transcode.RunSynchronousCmd("openssl", []string{"genrsa", "-out", keyPath, "2048"}); err != nil {
			return errors.New("generating new ssl key '" + os.Args[0] + "': " + err.Error())
		}
	}
	if _, err := os.Stat(crtPath); os.IsNotExist(err) {
		if err = transcode.RunSynchronousCmd("openssl", []string{"req", "-new", "-x509", "-sha256", "-key", keyPath, "-out", crtPath, "-days", "3650", "-subj", "/CN=localhost"}); err != nil {
			return errors.New("generating new ssl cert '" + os.Args[0] + "': " + err.Error())
		}
	}

	return nil
}
