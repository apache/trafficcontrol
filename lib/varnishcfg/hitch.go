package varnishcfg

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
	"path/filepath"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

// GetHitchConfig returns Hitch config using TO data
func GetHitchConfig(deliveryServices []atscfg.DeliveryService, sslDir string) (string, []string) {
	warnings := make([]string, 0)
	lines := []string{
		`frontend = {`,
		`	host = "*"`,
		`	port = "443"`,
		`}`,
		`backend = "[127.0.0.1]:6081"`,
		`write-proxy-v2 = on`,
		// TODO: change root user
		`user = "root"`,
	}

	dses, dsWarns := atscfg.DeliveryServicesToSSLMultiCertDSes(deliveryServices)
	warnings = append(warnings, dsWarns...)

	dses = atscfg.GetSSLMultiCertDotConfigDeliveryServices(dses)

	for dsName, ds := range dses {
		cerName, keyName := atscfg.GetSSLMultiCertDotConfigCertAndKeyName(dsName, ds)
		lines = append(lines, []string{
			`pem-file = {`,
			`	cert = "` + filepath.Join(sslDir, cerName) + `"`,
			`	private-key = "` + filepath.Join(sslDir, keyName) + `"`,
			`}`,
		}...)
	}

	txt := strings.Join(lines, "\n")
	txt += "\n"
	return txt, warnings
}
