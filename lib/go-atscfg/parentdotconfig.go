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
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const ParentConfigParamQStringHandling = "psel.qstring_handling"
const ParentConfigParamMSOAlgorithm = "mso.algorithm"
const ParentConfigParamMSOParentRetry = "mso.parent_retry"
const ParentConfigParamUnavailableServerRetryResponses = "mso.unavailable_server_retry_responses"
const ParentConfigParamMaxSimpleRetries = "mso.max_simple_retries"
const ParentConfigParamMaxUnavailableServerRetries = "mso.max_unavailable_server_retries"
const ParentConfigParamAlgorithm = "algorithm"
const ParentConfigParamQString = "qstring"

const ParentConfigDSParamDefaultMSOAlgorithm = "consistent_hash"
const ParentConfigDSParamDefaultMSOParentRetry = "both"
const ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses = ""
const ParentConfigDSParamDefaultMaxSimpleRetries = "1"
const ParentConfigDSParamDefaultMaxUnavailableServerRetries = "1"

const ParentConfigCacheParamWeight = "weight"
const ParentConfigCacheParamPort = "port"
const ParentConfigCacheParamUseIP = "use_ip_address"
const ParentConfigCacheParamRank = "rank"
const ParentConfigCacheParamNotAParent = "not_a_parent"

// TODO change, this is terrible practice, using a hard-coded key. What if there were a delivery service named "all_parents" (transliterated Perl)
const DeliveryServicesAllParentsKey = "all_parents"

type ParentConfigDS struct {
	Name            tc.DeliveryServiceName
	QStringIgnore   tc.QStringIgnore
	OriginFQDN      string
	MultiSiteOrigin bool
	OriginShield    string
	Type            tc.DSType
	QStringHandling string

	RequiredCapabilities map[ServerCapability]struct{}
}

type ParentConfigDSTopLevel struct {
	ParentConfigDS
	MSOAlgorithm                       string
	MSOParentRetry                     string
	MSOUnavailableServerRetryResponses string
	MSOMaxSimpleRetries                string
	MSOMaxUnavailableServerRetries     string
}

type ParentInfo struct {
	Host            string
	Port            int
	Domain          string
	Weight          string
	UseIP           bool
	Rank            int
	IP              string
	PrimaryParent   bool
	SecondaryParent bool
	Capabilities    map[ServerCapability]struct{}
}

func (p ParentInfo) Format() string {
	host := ""
	if p.UseIP {
		host = p.IP
	} else {
		host = p.Host + "." + p.Domain
	}
	return host + ":" + strconv.Itoa(p.Port) + "|" + p.Weight + ";"
}

type OriginHost string
type OriginFQDN string

type ParentInfos map[OriginHost]ParentInfo

type ParentInfoSortByRank []ParentInfo

func (s ParentInfoSortByRank) Len() int           { return len(([]ParentInfo)(s)) }
func (s ParentInfoSortByRank) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ParentInfoSortByRank) Less(i, j int) bool { return s[i].Rank < s[j].Rank }

type ParentConfigDSTopLevelSortByName []ParentConfigDSTopLevel

func (s ParentConfigDSTopLevelSortByName) Len() int      { return len(([]ParentConfigDSTopLevel)(s)) }
func (s ParentConfigDSTopLevelSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ParentConfigDSTopLevelSortByName) Less(i, j int) bool {
	// TODO make this match the Perl sort "foreach my $ds ( sort @{ $data->{dslist} } )" ?
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
}

type ProfileCache struct {
	Weight     string
	Port       int
	UseIP      bool
	Rank       int
	NotAParent bool
}

func DefaultProfileCache() ProfileCache {
	return ProfileCache{
		Weight:     "0.999",
		Port:       0,
		UseIP:      false,
		Rank:       1,
		NotAParent: false,
	}
}

// CGServer is the server table data needed when selecting the servers assigned to a cachegroup.
type CGServer struct {
	ServerID     ServerID
	ServerHost   string
	ServerIP     string
	ServerPort   int
	CacheGroupID int
	Status       int
	Type         int
	ProfileID    ProfileID
	CDN          int
	TypeName     string
	Domain       string
	Capabilities map[ServerCapability]struct{}
}

type OriginURI struct {
	Scheme string
	Host   string
	Port   string
}

func MakeParentDotConfig(
	serverInfo *ServerInfo, // getServerInfoByHost OR getServerInfoByID
	atsMajorVer int, // GetATSMajorVersion (TODO: determine if the cache itself [ORT via Yum] should produce this data, rather than asking TO?)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	parentConfigDSes []ParentConfigDSTopLevel, // getParentConfigDSTopLevel(cdn) OR getParentConfigDS(server) (TODO determine how to handle non-top missing MSO?)
	serverParams map[string]string, // getParentConfigServerProfileParams(serverID)
	parentInfos map[OriginHost][]ParentInfo, // getParentInfo(profileID, parentCachegroupID, secondaryParentCachegroupID)
) string {

	// parentInfos := makeParentInfo(serverInfo)

	nameVersionStr := GetNameVersionStringFromToolNameAndURL(toToolName, toURL)
	hdr := HeaderCommentWithTOVersionStr(serverInfo.HostName, nameVersionStr)

	textArr := []string{}
	text := ""
	// TODO put these in separate functions. No if-statement should be this long.
	if serverInfo.IsTopLevelCache() {
		uniqueOrigins := map[string]struct{}{}

		for _, ds := range parentConfigDSes {
			parentQStr := "ignore"
			if ds.QStringHandling == "" && ds.MSOAlgorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp {
				parentQStr = "consider"
			}

			orgURIStr := ds.OriginFQDN
			orgURI, err := url.Parse(orgURIStr) // TODO verify origin is always a host:port
			if err != nil {
				log.Errorln("Malformed ds '" + string(ds.Name) + "' origin  URI: '" + orgURIStr + "', skipping! : " + err.Error())
				continue
			}
			// TODO put in function, to remove duplication
			if orgURI.Port() == "" {
				if orgURI.Scheme == "http" {
					orgURI.Host += ":80"
				} else if orgURI.Scheme == "https" {
					orgURI.Host += ":443"
				} else {
					log.Errorln("parent.config generation: delivery service '" + string(ds.Name) + "' origin  URI: '" + orgURIStr + "' is unknown scheme '" + orgURI.Scheme + "', but has no port! Using as-is! ")
				}
			}

			if _, ok := uniqueOrigins[ds.OriginFQDN]; ok {
				continue // TODO warn?
			}
			uniqueOrigins[ds.OriginFQDN] = struct{}{}

			textLine := ""

			if ds.OriginShield != "" {
				algorithm := ""
				if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; strings.TrimSpace(parentSelectAlg) != "" {
					algorithm = "round_robin=" + parentSelectAlg
				}
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " parent=" + ds.OriginShield + " " + algorithm + " go_direct=true\n"
			} else if ds.MultiSiteOrigin {
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " "
				if len(parentInfos) == 0 {
				}

				if len(parentInfos[OriginHost(orgURI.Hostname())]) == 0 {
					// TODO error? emulates Perl
					log.Warnln("ParentInfo: delivery service " + ds.Name + " has no parent servers")
				}

				parents, secondaryParents := getMSOParentStrs(ds, parentInfos[OriginHost(orgURI.Hostname())], atsMajorVer)
				textLine += parents + secondaryParents + ` round_robin=` + ds.MSOAlgorithm + ` qstring=` + parentQStr + ` go_direct=false parent_is_proxy=false`

				parentRetry := ds.MSOParentRetry
				if atsMajorVer >= 6 && parentRetry != "" {
					if unavailableServerRetryResponsesValid(ds.MSOUnavailableServerRetryResponses) {
						textLine += ` parent_retry=` + parentRetry + ` unavailable_server_retry_responses=` + ds.MSOUnavailableServerRetryResponses
					} else {
						if ds.MSOUnavailableServerRetryResponses != "" {
							log.Errorln("Malformed unavailable_server_retry_responses parameter '" + ds.MSOUnavailableServerRetryResponses + "', not using!")
						}
						textLine += ` parent_retry=` + parentRetry
					}
					textLine += ` max_simple_retries=` + ds.MSOMaxSimpleRetries + ` max_unavailable_server_retries=` + ds.MSOMaxUnavailableServerRetries
				}
				textLine += "\n" // TODO remove, and join later on "\n" instead of ""?
				textArr = append(textArr, textLine)
			}
		}
		sort.Sort(sort.StringSlice(textArr))
		text = hdr + strings.Join(textArr, "")
	} else {
		processedOriginsToDSNames := map[string]tc.DeliveryServiceName{}

		queryStringHandling := serverParams[ParentConfigParamQStringHandling] // "qsh" in Perl

		roundRobin := `round_robin=consistent_hash`
		goDirect := `go_direct=false`

		sort.Sort(ParentConfigDSTopLevelSortByName(parentConfigDSes))

		for _, ds := range parentConfigDSes {
			parents, secondaryParents := getParentStrs(ds, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)

			text := ""
			originFQDN := ds.OriginFQDN
			if originFQDN == "" {
				continue // TODO warn? (Perl doesn't)
			}

			orgURI, err := url.Parse(originFQDN) // TODO verify
			if err != nil {
				log.Errorln("Malformed ds '" + string(ds.Name) + "' origin  URI: '" + originFQDN + "': skipping!" + err.Error())
				continue
			}

			if existingDS, ok := processedOriginsToDSNames[originFQDN]; ok {
				log.Errorln("parent.config generation: duplicate origin! services '" + string(ds.Name) + "' and '" + string(existingDS) + "' share origin '" + orgURI.Host + "': skipping '" + string(ds.Name) + "'!")
				continue
			}

			// TODO put in function, to remove duplication
			if orgURI.Port() == "" {
				if orgURI.Scheme == "http" {
					orgURI.Host += ":80"
				} else if orgURI.Scheme == "https" {
					orgURI.Host += ":443"
				} else {
					log.Errorln("parent.config generation non-top-level: ds '" + string(ds.Name) + "' origin  URI: '" + originFQDN + "' is unknown scheme '" + orgURI.Scheme + "', but has no port! Using as-is! ")
				}
			}

			// TODO encode this in a DSType func, IsGoDirect() ?
			if dsType := tc.DSType(ds.Type); dsType == tc.DSTypeHTTPNoCache || dsType == tc.DSTypeHTTPLive || dsType == tc.DSTypeDNSLive {
				text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` go_direct=true` + "\n"
			} else {

				// check for profile psel.qstring_handling.  If this parameter is assigned to the server profile,
				// then edges will use the qstring handling value specified in the parameter for all profiles.

				// If there is no defined parameter in the profile, then check the delivery service profile.
				// If psel.qstring_handling exists in the DS profile, then we use that value for the specified DS only.
				// This is used only if not overridden by a server profile qstring handling parameter.

				// TODO refactor this logic, hard to understand (transliterated from Perl)
				dsQSH := queryStringHandling
				if dsQSH == "" {
					dsQSH = ds.QStringHandling
				}
				parentQStr := dsQSH
				if parentQStr == "" {
					parentQStr = "ignore"
				}
				if ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp && dsQSH == "" {
					parentQStr = "consider"
				}

				text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` ` + parents + ` ` + secondaryParents + ` ` + roundRobin + ` ` + goDirect + ` qstring=` + parentQStr + "\n"
			}
			textArr = append(textArr, text)
			processedOriginsToDSNames[originFQDN] = ds.Name
		}

		parents, secondaryParents := getParentStrs(ParentConfigDSTopLevel{}, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)
		// TODO determine if this is necessary. It's super-dangerous, and moreover ignores Server Capabilitites.
		defaultDestText := `dest_domain=. ` + parents
		if serverParams[ParentConfigParamAlgorithm] == tc.AlgorithmConsistentHash {
			defaultDestText += secondaryParents
		}
		defaultDestText += ` round_robin=consistent_hash go_direct=false`

		if qStr := serverParams[ParentConfigParamQString]; qStr != "" {
			defaultDestText += ` qstring=` + qStr
		}
		defaultDestText += "\n"

		sort.Sort(sort.StringSlice(textArr))
		text = hdr + strings.Join(textArr, "") + defaultDestText
	}
	return text
}

// getParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines.
func getParentStrs(ds ParentConfigDSTopLevel, parentInfos []ParentInfo, atsMajorVer int) (string, string) {
	parentInfo := []string{}
	secondaryParentInfo := []string{}

	sort.Sort(ParentInfoSortByRank(parentInfos))

	for _, parent := range parentInfos { // TODO fix magic key
		if !HasRequiredCapabilities(parent.Capabilities, ds.RequiredCapabilities) {
			continue
		}

		pTxt := parent.Format()
		if parent.PrimaryParent {
			parentInfo = append(parentInfo, pTxt)
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, pTxt)
		}
	}

	if len(parentInfo) == 0 {
		parentInfo = secondaryParentInfo
		secondaryParentInfo = []string{}
	}

	// TODO remove duplicate code with top level if block
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfo, seen = util.RemoveStrDuplicates(parentInfo, seen)
	secondaryParentInfo, seen = util.RemoveStrDuplicates(secondaryParentInfo, seen)

	parents := ""
	secondaryParents := "" // "secparents" in Perl
	sort.Sort(sort.StringSlice(parentInfo))
	sort.Sort(sort.StringSlice(secondaryParentInfo))

	if atsMajorVer >= 6 && len(secondaryParentInfo) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + strings.Join(secondaryParentInfo, "") + `"`
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + strings.Join(secondaryParentInfo, "") + `"`
	}
	return parents, secondaryParents
}

// getMSOParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines, for MSO.
func getMSOParentStrs(ds ParentConfigDSTopLevel, parentInfos []ParentInfo, atsMajorVer int) (string, string) {
	// TODO determine why MSO is different, and if possible, combine with getParentAndSecondaryParentStrs.

	rankedParents := ParentInfoSortByRank(parentInfos)
	sort.Sort(rankedParents)

	parentInfo := []string{}
	secondaryParentInfo := []string{}
	nullParentInfo := []string{}
	for _, parent := range ([]ParentInfo)(rankedParents) {
		if !HasRequiredCapabilities(parent.Capabilities, ds.RequiredCapabilities) {
			continue
		}

		if parent.PrimaryParent {
			parentInfo = append(parentInfo, parent.Format())
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, parent.Format())
		} else {
			nullParentInfo = append(nullParentInfo, parent.Format())
		}
	}

	if len(parentInfo) == 0 {
		// If no parents are found in the secondary parent either, then set the null parent list (parents in neither secondary or primary)
		// as the secondary parent list and clear the null parent list.
		if len(secondaryParentInfo) == 0 {
			secondaryParentInfo = nullParentInfo
			nullParentInfo = []string{}
		}
		parentInfo = secondaryParentInfo
		secondaryParentInfo = []string{} // TODO should thi be '= secondary'? Currently emulates Perl
	}

	// TODO benchmark, verify this isn't slow. if it is, it could easily be made faster
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfo, seen = util.RemoveStrDuplicates(parentInfo, seen)
	secondaryParentInfo, seen = util.RemoveStrDuplicates(secondaryParentInfo, seen)
	nullParentInfo, seen = util.RemoveStrDuplicates(nullParentInfo, seen)

	secondaryParentStr := strings.Join(secondaryParentInfo, "") + strings.Join(nullParentInfo, "")

	// If the ats version supports it and the algorithm is consistent hash, put secondary and non-primary parents into secondary parent group.
	// This will ensure that secondary and tertiary parents will be unused unless all hosts in the primary group are unavailable.

	parents := ""
	secondaryParents := ""

	if atsMajorVer >= 6 && ds.MSOAlgorithm == "consistent_hash" && len(secondaryParents) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = `" secondary_parent="` + secondaryParentStr + `"`
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + secondaryParentStr + `"`
	}
	return parents, secondaryParents
}

func MakeParentInfo(
	server *ServerInfo,
	serverDomain string, // getCDNDomainByProfileID(tx, server.ProfileID)
	profileCaches map[ProfileID]ProfileCache, // getServerParentCacheGroupProfiles(tx, server)
	originServers map[OriginHost][]CGServer, // getServerParentCacheGroupProfiles(tx, server)
) map[OriginHost][]ParentInfo {
	parentInfos := map[OriginHost][]ParentInfo{}

	// note servers also contains an "all" key
	for originHost, servers := range originServers {
		for _, row := range servers {
			profile := profileCaches[row.ProfileID]
			if profile.NotAParent {
				continue
			}
			// Perl has this check, but we only select servers ("deliveryServices" in Perl) with the right CDN in the first place
			// if profile.Domain != serverDomain {
			// 	continue
			// }

			parentInf := ParentInfo{
				Host:            row.ServerHost,
				Port:            profile.Port,
				Domain:          row.Domain,
				Weight:          profile.Weight,
				UseIP:           profile.UseIP,
				Rank:            profile.Rank,
				IP:              row.ServerIP,
				PrimaryParent:   server.ParentCacheGroupID == row.CacheGroupID,
				SecondaryParent: server.SecondaryParentCacheGroupID == row.CacheGroupID,
				Capabilities:    row.Capabilities,
			}
			if parentInf.Port < 1 {
				parentInf.Port = row.ServerPort
			}
			parentInfos[originHost] = append(parentInfos[originHost], parentInf)
		}
	}
	return parentInfos
}

// unavailableServerRetryResponsesValid returns whether a unavailable_server_retry_responses parameter is valid for an ATS parent rule.
func unavailableServerRetryResponsesValid(s string) bool {
	// optimization if param is empty
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^"(:?\d{3},)+\d{3}"\s*$`) // TODO benchmark, cache if performance matters
	return re.MatchString(s)
}

// HasRequiredCapabilities returns whether the given caps has all the required capabilities in the given reqCaps.
func HasRequiredCapabilities(caps map[ServerCapability]struct{}, reqCaps map[ServerCapability]struct{}) bool {
	for reqCap, _ := range reqCaps {
		if _, ok := caps[reqCap]; !ok {
			return false
		}
	}
	return true
}
