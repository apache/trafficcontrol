package main

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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestPreprocessConfigFile(t *testing.T) {
	// the TCP port replacement is fundamentally different for 80 vs non-80, so test both
	{
		server := &atscfg.Server{}
		server.TCPPort = util.IntPtr(8080)
		server.Interfaces = []tc.ServerInterfaceInfo{
			tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					tc.ServerIPAddress{
						Address:        "127.0.2.1",
						ServiceAddress: true,
					},
				},
			},
		}
		server.HostName = util.StrPtr("my-edge")
		server.DomainName = util.StrPtr("example.net")
		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc8080def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}

	{
		server := &atscfg.Server{}
		server.TCPPort = util.IntPtr(80)
		server.Interfaces = []tc.ServerInterfaceInfo{
			tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					tc.ServerIPAddress{
						Address:        "127.0.2.1",
						ServiceAddress: true,
					},
				},
			},
		}
		server.HostName = util.StrPtr("my-edge")
		server.DomainName = util.StrPtr("example.net")

		cfgFile := "abc__SERVER_TCP_PORT__def__CACHE_IPV4__ghi __RETURN__  \t __HOSTNAME__ jkl __FULL_HOSTNAME__ \n__SOMETHING__ __ELSE__\nmno:__SERVER_TCP_PORT__\r\n"

		actual := PreprocessConfigFile(server, cfgFile)

		expected := "abc__SERVER_TCP_PORT__def127.0.2.1ghi\nmy-edge jkl my-edge.example.net \n__SOMETHING__ __ELSE__\nmno\r\n"

		if expected != actual {
			t.Errorf("PreprocessConfigFile expected '%v' actual '%v'", expected, actual)
		}
	}
}
