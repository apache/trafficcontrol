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
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func (v *VCLBuilder) configureDirectors(vclFile *vclFile, parents *atscfg.ParentAbstraction) ([]string, error) {
	warnings := []string{}

	vclFile.imports = append(vclFile.imports, "directors")
	var err error
	requestFQDNs := make([]string, 0)

	for _, svc := range parents.Services {
		addBackends(vclFile.backends, append(svc.Parents, svc.SecondaryParents...), svc.DestDomain, svc.Port)
		addDirectors(vclFile.subroutines, svc)

		requestFQDNs = []string{svc.DestDomain}

		if v.toData.Server.Type == tc.CacheTypeEdge.String() {
			dsRegexes := atscfg.MakeDSRegexMap(v.toData.DeliveryServiceRegexes)
			anyCastPartners := atscfg.GetAnyCastPartners(v.toData.Server, v.toData.Servers)
			requestFQDNs, err = atscfg.GetDSRequestFQDNs(
				&svc.DS,
				dsRegexes[tc.DeliveryServiceName(svc.DS.XMLID)],
				v.toData.Server,
				anyCastPartners,
				v.toData.CDN.DomainName,
			)
			if err != nil {
				warnings = append(warnings, "error getting ds '"+svc.DS.XMLID+"' request fqdns, skipping! Error: "+err.Error())
				continue
			}
		}

		assignBackends(vclFile.subroutines, svc, requestFQDNs)
	}

	return warnings, nil
}

func assignBackends(subroutines map[string][]string, svc *atscfg.ParentAbstractionService, requestFQDNs []string) {
	lines := make([]string, 0)
	hostHeaderLines := make([]string, 0)

	conditions := make([]string, 0)
	backendConditions := make([]string, 0)
	for _, fqdn := range requestFQDNs {
		conditions = append(conditions, fmt.Sprintf(`req.http.host == "%s"`, fqdn))
		backendConditions = append(backendConditions, fmt.Sprintf(`bereq.http.host == "%s"`, fqdn))
	}

	lines = append(lines, fmt.Sprintf("if (%s) {", strings.Join(conditions, " || ")))
	lines = append(lines, fmt.Sprintf("\tset req.backend_hint = %s.backend();", svc.Name))

	// only change request host from edge servers which typically has multiple request FQDNs or
	// one request FQDN that is not the origin.
	if len(requestFQDNs) > 1 || (len(requestFQDNs) == 1 && requestFQDNs[0] != svc.DestDomain) {
		hostHeaderLines = append(hostHeaderLines, fmt.Sprintf("if (%s) {", strings.Join(backendConditions, " || ")))
		hostHeaderLines = append(hostHeaderLines, fmt.Sprintf("\tset bereq.http.host = \"%s\";", svc.DestDomain))
		hostHeaderLines = append(hostHeaderLines, "}")
	}

	lines = append(lines, "}")

	subroutines["vcl_recv"] = append(subroutines["vcl_recv"], lines...)
	if len(hostHeaderLines) == 0 {
		return
	}

	subroutines["vcl_backend_fetch"] = append(subroutines["vcl_backend_fetch"], hostHeaderLines...)
}

func addBackends(backends map[string]backend, parents []*atscfg.ParentAbstractionServiceParent, originDomain string, originPort int) {
	for _, parent := range parents {
		backendName := fmt.Sprintf("%s", getBackendName(parent.FQDN, parent.Port))
		if _, ok := backends[backendName]; ok {
			continue
		}
		backends[backendName] = backend{
			host: parent.FQDN,
			port: parent.Port,
		}
	}
	backendName := getBackendName(originDomain, originPort)
	if _, ok := backends[backendName]; ok {
		return
	}
	backends[backendName] = backend{
		host: originDomain,
		port: originPort,
	}
}

func addDirectors(subroutines map[string][]string, svc *atscfg.ParentAbstractionService) {
	lines := make([]string, 0)
	fallbackDirectorLines := make([]string, 0)
	fallbackDirectorLines = append(fallbackDirectorLines, fmt.Sprintf("new %s = directors.fallback();", svc.Name))

	if len(svc.Parents) != 0 {
		lines = append(lines, addBackendsToDirector(svc.Name+"_primary", svc.RetryPolicy, svc.Parents)...)
		fallbackDirectorLines = append(fallbackDirectorLines, fmt.Sprintf("%s.add_backend(%s_primary.backend());", svc.Name, svc.Name))
	}
	if len(svc.SecondaryParents) != 0 {
		lines = append(lines, addBackendsToDirector(svc.Name+"_secondary", svc.RetryPolicy, svc.SecondaryParents)...)
		fallbackDirectorLines = append(fallbackDirectorLines, fmt.Sprintf("%s.add_backend(%s_secondary.backend());", svc.Name, svc.Name))
	}
	fallbackDirectorLines = append(fallbackDirectorLines, fmt.Sprintf("%s.add_backend(%s);", svc.Name, getBackendName(svc.DestDomain, svc.Port)))

	lines = append(lines, fallbackDirectorLines...)

	subroutines["vcl_init"] = append(subroutines["vcl_init"], lines...)
}

func addBackendsToDirector(name string, retryPolicy atscfg.ParentAbstractionServiceRetryPolicy, parents []*atscfg.ParentAbstractionServiceParent) []string {
	lines := make([]string, 0)
	directorType, sticky := getDirectorType(retryPolicy)
	lines = append(lines, fmt.Sprintf("new %s = directors.%s(%s);", name, directorType, sticky))
	for _, parent := range parents {
		lines = append(lines, fmt.Sprintf("%s.add_backend(%s);", name, getBackendName(parent.FQDN, parent.Port)))
	}
	return lines
}

func getDirectorType(retryPolicy atscfg.ParentAbstractionServiceRetryPolicy) (director string, sticky string) {
	switch retryPolicy {
	case atscfg.ParentAbstractionServiceRetryPolicyRoundRobinIP:
		fallthrough
	case atscfg.ParentAbstractionServiceRetryPolicyRoundRobinStrict:
		director = "round_robin"
	case atscfg.ParentAbstractionServiceRetryPolicyFirst:
		director = "fallback"
	case atscfg.ParentAbstractionServiceRetryPolicyLatched:
		director = "fallback"
		sticky = "1"
	case atscfg.ParentAbstractionServiceRetryPolicyConsistentHash:
		director = "shard"
	default:
		director = "shard"
	}
	return
}

func getBackendName(host string, port int) string {
	// maybe a better way to ensure backend names are unique?

	if port <= 0 {
		return strings.ReplaceAll(host, ".", "_")
	}
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(host, ".", "_"), port)
}
