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
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

// readDefaultUnixResolve reads the /etc/resolv.conf
// file and parses out the nameservers.
func readDefaultUnixResolve() ([]string, error) {
	fd, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return parseResolve(fd)
}

// parseResolve parses r as a resolv.conf formatted string.
// It returns all the nameservers found in the file. Any
// formatting or other issues within the file itself are ignored,
// only errors reading from r are returned.
// See following link for more information:
// http://man7.org/linux/man-pages/man5/resolv.conf.5.html
func parseResolve(r io.Reader) ([]string, error) {
	var nameservers []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if len(l) > 0 && (l[0] == '#' || l[0] == ';') {
			// Ignore comments
			continue
		}
		parts := strings.Fields(l)
		// Look for "nameserver 0.0.0.0" formatted lines
		if len(parts) < 2 || parts[0] != "nameserver" {
			continue
		}
		if net.ParseIP(parts[1]) == nil {
			// Ignore invalid IPs
			continue
		}
		nameservers = append(nameservers, parts[1])
	}

	return nameservers, s.Err()
}
