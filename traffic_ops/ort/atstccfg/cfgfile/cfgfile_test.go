package cfgfile

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

func TestWriteConfigs(t *testing.T) {
	buf := &bytes.Buffer{}
	configs := []config.ATSConfigFile{
		{
			ATSConfigMetaDataConfigFile: tc.ATSConfigMetaDataConfigFile{
				FileNameOnDisk: "config0.txt",
				Location:       "/my/config0/location",
			},
			Text:        "config0",
			ContentType: "text/plain",
		},
		{
			ATSConfigMetaDataConfigFile: tc.ATSConfigMetaDataConfigFile{
				FileNameOnDisk: "config1.txt",
				Location:       "/my/config1/location",
			},
			Text:        "config2,foo",
			ContentType: "text/csv",
		},
	}

	if err := WriteConfigs(configs, buf); err != nil {
		t.Fatalf("WriteConfigs error expected nil, actual: %v", err)
	}

	actual := buf.String()

	expected0 := "Content-Type: text/plain\r\nLine-Comment: \r\nPath: /my/config0/location/config0.txt\r\n\r\nconfig0\r\n"

	if !strings.Contains(actual, expected0) {
		t.Errorf("WriteConfigs expected '%v' actual '%v'", expected0, actual)
	}

	expected1 := "Content-Type: text/csv\r\nLine-Comment: \r\nPath: /my/config1/location/config1.txt\r\n\r\nconfig2,foo\r\n"
	if !strings.Contains(actual, expected1) {
		t.Errorf("WriteConfigs expected config1 '%v' actual '%v'", expected1, actual)
	}

	expectedPrefix := "MIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary="
	if !strings.HasPrefix(actual, expectedPrefix) {
		t.Errorf("WriteConfigs expected prefix '%v' actual '%v'", expectedPrefix, actual)
	}
}

func TestPreprocessConfigFile(t *testing.T) {
	// the TCP port replacement is fundamentally different for 80 vs non-80, so test both
	{
		server := tc.Server{
			TCPPort:    8080,
			IPAddress:  "127.0.2.1",
			HostName:   "my-edge",
			DomainName: "example.net",
		}
		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc8080def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}

	{
		server := tc.Server{
			TCPPort:    80,
			IPAddress:  "127.0.2.1",
			HostName:   "my-edge",
			DomainName: "example.net",
		}
		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno:__SERVER_TCP_PORT__\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc__SERVER_TCP_PORT__def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}
}
