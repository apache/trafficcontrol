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
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const StorageFileName = "storage.config"
const ContentTypeStorageDotConfig = ContentTypeTextASCII
const LineCommentStorageDotConfig = LineCommentHash

// StorageDotConfigOpts contains settings to configure generation options.
type StorageDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeStorageDotConfig creates storage.config for a given ATS Profile.
// The paramData is the map of parameter names to values, for all parameters assigned to the given profile, with the config_file "storage.config".
func MakeStorageDotConfig(
	server *Server,
	serverParams []tc.Parameter,
	opt *StorageDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &StorageDotConfigOpts{}
	}
	warnings := []string{}

	if len(server.ProfileNames) == 0 {
		return Cfg{}, makeErr(warnings, "server missing Profiles")
	}

	if server.HostName == nil {
		return Cfg{}, makeErr(warnings, "server missing HostName")
	}

	paramData, paramWarns := paramsToMap(filterParams(serverParams, StorageFileName, "", "", "location"))
	warnings = append(warnings, paramWarns...)

	text := ""

	nextVolume := 1
	if drivePrefix := paramData["Drive_Prefix"]; drivePrefix != "" {
		driveLetters := strings.TrimSpace(paramData["Drive_Letters"])
		if driveLetters == "" {
			warnings = append(warnings, fmt.Sprintf("server %+v profile has Drive_Prefix parameter, but no Drive_Letters; creating anyway", *server.HostName))
		}
		text += makeStorageVolumeText(drivePrefix, driveLetters, nextVolume)
		nextVolume++
	}

	if ramDrivePrefix := paramData["RAM_Drive_Prefix"]; ramDrivePrefix != "" {
		ramDriveLetters := strings.TrimSpace(paramData["RAM_Drive_Letters"])
		if ramDriveLetters == "" {
			warnings = append(warnings, fmt.Sprintf("server %+v profile has RAM_Drive_Prefix parameter, but no RAM_Drive_Letters; creating anyway", *server.HostName))
		}
		text += makeStorageVolumeText(ramDrivePrefix, ramDriveLetters, nextVolume)
		nextVolume++
	}

	if ssdDrivePrefix := paramData["SSD_Drive_Prefix"]; ssdDrivePrefix != "" {
		ssdDriveLetters := strings.TrimSpace(paramData["SSD_Drive_Letters"])
		if ssdDriveLetters == "" {
			warnings = append(warnings, fmt.Sprintf("server %+v profile has SSD_Drive_Prefix parameter, but no SSD_Drive_Letters; creating anyway", *server.HostName))
		}
		text += makeStorageVolumeText(ssdDrivePrefix, ssdDriveLetters, nextVolume)
		nextVolume++
	}

	if text == "" {
		text = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	hdr := makeHdrComment(opt.HdrComment)
	text = hdr + text

	return Cfg{
		Text:        text,
		ContentType: ContentTypeStorageDotConfig,
		LineComment: LineCommentStorageDotConfig,
		Warnings:    warnings,
	}, nil
}

func makeStorageVolumeText(prefix string, letters string, volume int) string {
	text := ""
	for _, letter := range strings.Split(letters, ",") {
		letter = strings.TrimSpace(letter)
		if letter == "" {
			continue
		}
		text += prefix + letter + " volume=" + strconv.Itoa(volume) + "\n"
	}
	return text
}
