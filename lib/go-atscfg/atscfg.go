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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const HeaderCommentDateFormat = "Mon Jan 2 15:04:05 MST 2006"

type ServerInfo struct {
	CacheGroupID                  int
	CDN                           tc.CDNName
	CDNID                         int
	DomainName                    string
	HostName                      string
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

// GetATSMajorVersionFromATSVersion returns the major version of the given profile's package trafficserver parameter.
// The atsVersion is typically a Parameter on the Server's Profile, with the configFile "package" name "trafficserver".
// Returns an error if atsVersion is empty or not a number.
func GetATSMajorVersionFromATSVersion(atsVersion string) (int, error) {
	if len(atsVersion) == 0 {
		return 0, errors.New("ats version missing")
	}
	atsMajorVer, err := strconv.Atoi(atsVersion[:1])
	if err != nil {
		return 0, errors.New("ats version parameter '" + atsVersion + "' is not a number")
	}
	return atsMajorVer, nil
}

type DeliveryServiceID int
type ProfileID int
type ServerID int
