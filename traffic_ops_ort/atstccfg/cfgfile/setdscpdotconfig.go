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
	"errors"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetConfigFileCDNSetDSCP(toData *config.TOData, fileName string) (string, string, string, error) {
	if toData.Server.CDNName == nil {
		return "", "", "", errors.New("server missing CDNName")
	}

	// TODO verify prefix, suffix, and that it's a number? Perl doesn't.
	dscpNumStr := fileName
	dscpNumStr = strings.TrimPrefix(dscpNumStr, "set_dscp_")
	dscpNumStr = strings.TrimSuffix(dscpNumStr, ".config")

	return atscfg.MakeSetDSCPDotConfig(tc.CDNName(*toData.Server.CDNName), toData.TOToolName, toData.TOURL, dscpNumStr), atscfg.ContentTypeSetDSCPDotConfig, atscfg.LineCommentSetDSCPDotConfig, nil
}
