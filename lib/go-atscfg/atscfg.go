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
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const InvalidID = -1

const DefaultATSVersion = "5" // TODO Emulates Perl; change to 6? ATC no longer officially supports ATS 5.

const HeaderCommentDateFormat = "Mon Jan 2 15:04:05 MST 2006"

type ServerCapability string

type ServerInfo struct {
	CacheGroupID                  int
	CDN                           tc.CDNName
	CDNID                         int
	DomainName                    string
	HostName                      string
	HTTPSPort                     int
	ID                            int
	IP                            string
	ParentCacheGroupID            int
	ParentCacheGroupType          string
	ProfileID                     ProfileID
	ProfileName                   string
	Port                          int
	SecondaryParentCacheGroupID   int
	SecondaryParentCacheGroupType string
	Type                          string
}

func (s *ServerInfo) IsTopLevelCache() bool {
	return (s.ParentCacheGroupType == tc.CacheGroupOriginTypeName || s.ParentCacheGroupID == InvalidID) &&
		(s.SecondaryParentCacheGroupType == tc.CacheGroupOriginTypeName || s.SecondaryParentCacheGroupID == InvalidID)
}

func HeaderCommentWithTOVersionStr(name string, nameVersionStr string) string {
	return "# DO NOT EDIT - Generated for " + name + " by " + nameVersionStr + " on " + time.Now().Format(HeaderCommentDateFormat) + "\n"
}

func GetNameVersionStringFromToolNameAndURL(toolName string, url string) string {
	return toolName + " (" + url + ")"
}

func GenericHeaderComment(name string, toolName string, url string) string {
	return HeaderCommentWithTOVersionStr(name, GetNameVersionStringFromToolNameAndURL(toolName, url))
}

// GetATSMajorVersionFromATSVersion returns the major version of the given profile's package trafficserver parameter.
// The atsVersion is typically a Parameter on the Server's Profile, with the configFile "package" name "trafficserver".
// Returns an error if atsVersion is empty or does not start with an unsigned integer followed by a period or nothing.
func GetATSMajorVersionFromATSVersion(atsVersion string) (int, error) {
	dotPos := strings.Index(atsVersion, ".")
	if dotPos == -1 {
		dotPos = len(atsVersion) // if there's no '.' then assume the whole string is just a major version.
	}
	majorVerStr := atsVersion[:dotPos]

	majorVer, err := strconv.ParseUint(majorVerStr, 10, 64)
	if err != nil {
		return 0, errors.New("unexpected version format, expected e.g. '7.1.2.whatever'")
	}
	return int(majorVer), nil
}

type DeliveryServiceID int
type ProfileID int
type ServerID int

// GenericProfileConfig generates a generic profile config text, from the profile's parameters with the given config file name.
// This does not include a header comment, because a generic config may not use a number sign as a comment.
// If you need a header comment, it can be added manually via ats.HeaderComment, or automatically with WithProfileDataHdr.
func GenericProfileConfig(
	paramData map[string]string, // GetProfileParamData(tx, profileID, fileName)
	separator string,
) string {
	text := ""
	for name, val := range paramData {
		name = trimParamUnderscoreNumSuffix(name)
		text += name + separator + val + "\n"
	}
	return text
}

// trimParamUnderscoreNumSuffix removes any trailing "__[0-9]+" and returns the trimmed string.
func trimParamUnderscoreNumSuffix(paramName string) string {
	underscorePos := strings.LastIndex(paramName, `__`)
	if underscorePos == -1 {
		return paramName
	}
	if _, err := strconv.ParseFloat(paramName[underscorePos+2:], 64); err != nil {
		return paramName
	}
	return paramName[:underscorePos]
}

const ConfigSuffix = ".config"

func GetConfigFile(prefix string, xmlId string) string {
	return prefix + xmlId + ConfigSuffix
}
