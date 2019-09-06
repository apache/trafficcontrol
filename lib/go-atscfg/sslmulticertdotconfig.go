package atscfg

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type SSLMultiCertDS struct {
	Type        tc.DSType
	Protocol    int
	ExampleURLs []string
}

func DeliveryServicesToSSLMultiCertDSes(dses []tc.DeliveryServiceNullable) map[tc.DeliveryServiceName]SSLMultiCertDS {
	sDSes := map[tc.DeliveryServiceName]SSLMultiCertDS{}
	for _, ds := range dses {
		if ds.Type == nil || ds.Protocol == nil || ds.XMLID == nil {
			if ds.XMLID == nil {
				log.Errorln("atscfg.DeliveryServicesToSSLMultiCertDSes got unknown DS with nil values! Skipping!")
			} else {
				log.Errorln("atscfg.DeliveryServicesToSSLMultiCertDSes got DS '" + *ds.XMLID + "' with nil values! Skipping!")
			}
			continue
		}
		sDSes[tc.DeliveryServiceName(*ds.XMLID)] = SSLMultiCertDS{Type: *ds.Type, Protocol: *ds.Protocol, ExampleURLs: ds.ExampleURLs}
	}
	return sDSes
}

func MakeSSLMultiCertDotConfig(
	cdnName tc.CDNName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	dses map[tc.DeliveryServiceName]SSLMultiCertDS,
) string {
	text := GenericHeaderComment(string(cdnName), toToolName, toURL)

	for dsName, ds := range dses {
		if ds.Type.IsSteering() {
			continue // Steering delivery service SSLs should not be on the edges.
		}
		if ds.Protocol == 0 {
			continue
		}
		if len(ds.ExampleURLs) == 0 {
			continue // TODO warn? error? Perl doesn't
		}

		hostName := ds.ExampleURLs[0] // first one is the one we want

		scheme := "https://"
		if !strings.HasPrefix(hostName, scheme) {
			scheme = "http://"
		}
		newHost := hostName
		if len(hostName) < len(scheme) {
			log.Errorln("MakeSSLMultiCertDotConfig got ds '" + string(dsName) + "' example url '" + hostName + "' with no scheme! ssl_multicert.config will likely be malformed!")
		} else {
			newHost = hostName[len(scheme):]
		}
		keyName := newHost + ".key"

		newHost = strings.Replace(newHost, ".", "_", -1)

		cerName := newHost + "_cert.cer"
		text += `ssl_cert_name=` + cerName + "\t" + ` ssl_key_name=` + keyName + "\n"
	}
	return text
}
