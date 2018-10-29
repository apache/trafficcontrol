package ats

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
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

type ProfileConfigFunc = func(tx *sql.Tx, profile ProfileData, fileName string) (string, error)

func GetProfileConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id", "file"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	profileID := inf.IntParams["id"]
	fileName := strings.TrimSuffix(inf.Params["file"], ".json")

	scope, err := getScope(inf.Tx.Tx, fileName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetProfileConfig getting scope: "+err.Error()))
		return
	}

	if scope != tc.ATSConfigMetaDataConfigFileScopeProfiles {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("incorrect file scope for route used.  Please use the "+string(scope)+" route."), nil)
		return
	}

	profile, ok, err := getProfileData(inf.Tx.Tx, profileID)
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	fileContents := ""

	// TODO add routes for each of these, rather than dispatching ourselves
	err = nil
	switch {
	case fileName == "50-ats.rules":
		fileContents, err = atsDotRules(inf.Tx.Tx, profile, fileName)
	case fileName == "12M_facts":
		fileContents, err = facts(inf.Tx.Tx, profile, fileName)
	case fileName == "astats.config":
		fileContents, err = genericProfileConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "cache.config":
		fileContents, err = profileCacheDotConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "drop_qstring.config":
		fileContents, err = dropQStringDotConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "logs_xml.config":
		fileContents, err = logsXMLDotConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "logging.config":
		fileContents, err = loggingDotConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "plugin.config":
		fileContents, err = genericProfileConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "records.config":
		fileContents, err = genericProfileConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "storage.config":
		fileContents, err = storageDotConfig(inf.Tx.Tx, profile, fileName)
	case fileName == "sysctl.conf":
		fileContents, err = genericProfileConfig(inf.Tx.Tx, profile, fileName)
	case strings.HasPrefix(fileName, "url_sig_") && strings.HasSuffix(fileName, ".config"):
		fileContents, err = urlSigDotConfig(inf.Tx.Tx, inf.Config, profile, fileName)
	case strings.HasPrefix(fileName, "uri_signing_") && strings.HasSuffix(fileName, ".config"):
		fileContents, err = uriSigningDotConfig(inf.Tx.Tx, inf.Config, fileName)
	case fileName == "volume.config":
		fileContents, err = volumeDotConfig(inf.Tx.Tx, profile, fileName)
	default:
		// TODO move to func, "getUnknownConfig"
		params, err := GetProfileParamData(inf.Tx.Tx, profile.ID, fileName) // (map[string]string, error) {
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetProfileConfig: getting profile parameter data: "+err.Error()))
			return
		}
		fileContents, err = takeAndBakeProfile(inf.Tx.Tx, profile.Name, params)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetProfileConfig: takeAndBakeProfile: "+err.Error()))
			return
		}
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting file contents: "+err.Error()))
		return
	}

	if fileContents == "" {
		// TODO replicates old Perl; verify required.
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fileContents))
}

func takeAndBakeProfile(tx *sql.Tx, profileName string, params map[string]string) (string, error) {
	hdr, err := headerComment(tx, profileName)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	text := ""
	for paramName, paramVal := range params {
		if paramName == "header" {
			if paramVal == "none" {
				hdr = ""
			} else {
				hdr = paramVal + "\n"
			}
		} else {
			text += paramVal + "\n"
		}
	}
	text = strings.Replace(text, "__RETURN__", "\n", -1)
	return hdr + text, nil
}

func atsDotRules(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}

	// TODO add more efficient db func to only get drive params?
	paramData, err := GetProfileParamData(tx, profile.ID, "storage.config") // ats.rules is based on the storage.config params
	if err != nil {
		return "", errors.New("profile param data: " + err.Error())
	}

	drivePrefix := strings.TrimPrefix(paramData["Drive_Prefix"], `/dev/`)
	drivePostfix := strings.Split(paramData["Drive_Letters"], ",")
	for _, l := range drivePostfix {
		text += `KERNEL=="` + drivePrefix + l + `", OWNER="ats"` + "\n"
	}
	if ramPrefix, ok := paramData["RAM_Drive_Prefix"]; ok {
		ramPrefix = strings.TrimPrefix(ramPrefix, `/dev/`)
		ramPostfix := strings.Split(paramData["RAM_Drive_Letters"], ",")
		for _, l := range ramPostfix {
			text += `KERNEL=="` + ramPrefix + l + `", OWNER="ats"` + "\n"
		}
	}
	return text, nil
}

func facts(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", err
	}
	text += "profile:" + profile.Name + "\n"
	return text, nil
}

func separators() map[string]string {
	return map[string]string{
		"records.config":  " ",
		"plugin.config":   " ",
		"sysctl.conf":     " = ",
		"url_sig_.config": " = ",
		"astats.config":   "=",
	}
}

func profileCacheDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	lines := map[string]struct{}{} // use a "set" for lines, to avoid duplicates, since we're looking up by profile
	headerLine, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", err
	}

	profileDSes, err := GetProfileDS(tx, profile.ID)
	if err != nil {
		return "", errors.New("getting profile delivery services: " + err.Error())
	}

	for _, ds := range profileDSes {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}
		if ds.OriginFQDN == nil || *ds.OriginFQDN == "" {
			log.Warnf("profileCacheDotConfig ds has no origin fqdn, skipping!") // TODO add ds name to data loaded, to put it in the error here?
			continue
		}
		originFQDN, originPort := getHostPortFromURI(*ds.OriginFQDN)
		if originPort != "" {
			l := "dest_domain=" + originFQDN + " port=" + originPort + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		} else {
			l := "dest_domain=" + originFQDN + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		}
	}

	text := headerLine
	for line, _ := range lines {
		text += line
	}
	return text, nil
}

func dropQStringDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}

	dropQStringVal, hasDropQStringParam, err := GetProfileParamValue(tx, profile.ID, "drop_qstring.config", "content")
	if err != nil {
		return "", errors.New("getting profile param val: " + err.Error())
	}
	if hasDropQStringParam {
		text += dropQStringVal + "\n"
	} else {
		text += `/([^?]+) $s://$t/$1` + "\n"
	}
	return text, nil
}

const MaxLogObjects = 10

func logsXMLDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	profileParamData, err := GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	hdrComment, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	hdrComment = strings.Replace(hdrComment, `# `, ``, -1)
	hdrComment = strings.Replace(hdrComment, "\n", ``, -1)
	text := "<!-- " + hdrComment + " -->\n"

	for i := 0; i < MaxLogObjects; i++ {
		logFormatField := "LogFormat"
		logObjectField := "LogObject"
		if i > 0 {
			iStr := strconv.Itoa(i)
			logFormatField += iStr
			logObjectField += iStr
		}

		logFormatName := profileParamData[logFormatField+".Name"]
		if logFormatName != "" {
			format := profileParamData[logFormatField+".Format"]
			format = strings.Replace(format, `"`, `\"`, -1)

			text += `<LogFormat>
  <Name = "` + logFormatName + `"/>
  <Format = "` + format + `"/>
</LogFormat>
`
		}

		logObjectFileName := profileParamData[logObjectField+".Filename"]
		if logObjectFileName != "" {
			logObjectFormat := profileParamData[logObjectField+".Format"]
			logObjectRollingEnabled := profileParamData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := profileParamData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := profileParamData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := profileParamData[logObjectField+".RollingSizeMb"]
			logObjectHeader := profileParamData[logObjectField+".Header"]

			text += `<LogObject>
  <Format = "` + logObjectFormat + `"/>
  <Filename = "` + logObjectFileName + `"/>
`
			if logObjectRollingEnabled != "" {
				text += `  <RollingEnabled = ` + logObjectRollingEnabled + `/>
`
			}
			text += `<RollingIntervalSec = ` + logObjectRollingIntervalSec + `/>
  <RollingOffsetHr = ` + logObjectRollingOffsetHr + `/>
  <RollingSizeMb = ` + logObjectRollingSizeMb + `/>
`
			if logObjectHeader != "" {
				text += `  <Header = "` + logObjectHeader + `"/>
`
			}
			text += `</LogObject>
`
		}
	}
	return text, nil
}

func loggingDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	profileParamData, err := GetProfileParamData(tx, profile.ID, fileName)
	log.Errorf("DEBUG fileName: %+v\n", fileName)
	log.Errorf("DEBUG len(profileParamData): %+v\n", len(profileParamData))
	for k, v := range profileParamData {
		log.Errorf("DEBUG profileParamData[%v] = %v\n", k, v)
	}

	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	hdrComment, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	// This is an LUA file, so we need to massage the header a bit for LUA commenting.
	hdrComment = strings.Replace(hdrComment, `# `, ``, -1)
	hdrComment = strings.Replace(hdrComment, "\n", ``, -1)
	text := "-- " + hdrComment + " --\n"

	for i := 0; i < MaxLogObjects; i++ {
		logFormatField := "LogFormat"
		logObjectField := "LogObject"
		if i > 0 {
			iStr := strconv.Itoa(i)
			logFormatField += iStr
			logObjectField += iStr
		}

		logFormatName := profileParamData[logFormatField+".Name"]
		if logFormatName != "" {
			format := profileParamData[logFormatField+".Format"]
			format = strings.Replace(format, `"`, `\"`, -1)
			text += logFormatName + ` = format {
	Format = '` + format + ` '
}
`
		}

		if logObjectFileName := profileParamData[logObjectField+".Filename"]; logObjectFileName != "" {
			logObjectRollingEnabled := profileParamData[logObjectField+".RollingEnabled"]
			logObjectRollingIntervalSec := profileParamData[logObjectField+".RollingIntervalSec"]
			logObjectRollingOffsetHr := profileParamData[logObjectField+".RollingOffsetHr"]
			logObjectRollingSizeMb := profileParamData[logObjectField+".RollingSizeMb"]

			text += `
log.ascii {
  Format = ` + logFormatName + `,
  Filename = '` + logObjectFileName + `',
`
			if logObjectRollingEnabled != "" {
				text += "  RollingEnabled = " + logObjectRollingEnabled + ",\n"
			}
			text += `  RollingIntervalSec = ` + logObjectRollingIntervalSec + `,
  RollingOffsetHr = ` + logObjectRollingOffsetHr + `,
  RollingSizeMb = ` + logObjectRollingSizeMb + `
}
`
		}
	}
	return text, nil
}

func storageDotConfigVolumeText(prefix string, letters string, volume int) string {
	text := ""
	for _, letter := range strings.Split(letters, ",") {
		text += prefix + letter + " volume=" + strconv.Itoa(volume) + "\n"
	}
	return text
}

func storageDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}

	paramData, err := GetProfileParamData(tx, profile.ID, "storage.config") // ats.rules is based on the storage.config params
	if err != nil {
		return "", errors.New("profile param data: " + err.Error())
	}

	nextVolume := 1
	if drivePrefix := paramData["Drive_Prefix"]; drivePrefix != "" {
		driveLetters := strings.TrimSpace(paramData["Drive_Letters"])
		if driveLetters == "" {
			log.Warnf("Creating storage.config: profile %+v has Drive_Prefix parameter, but no Drive_Letters; creating anyway")
		}
		text += storageDotConfigVolumeText(drivePrefix, driveLetters, nextVolume)
		nextVolume++
	}

	if ramDrivePrefix := paramData["RAM_Drive_Prefix"]; ramDrivePrefix != "" {
		ramDriveLetters := strings.TrimSpace(paramData["RAM_Drive_Letters"])
		if ramDriveLetters == "" {
			log.Warnf("Creating storage.config: profile %+v has RAM_Drive_Prefix parameter, but no RAM_Drive_Letters; creating anyway")
		}
		text += storageDotConfigVolumeText(ramDrivePrefix, ramDriveLetters, nextVolume)
		nextVolume++
	}

	if ssdDrivePrefix := paramData["SSD_Drive_Prefix"]; ssdDrivePrefix != "" {
		ssdDriveLetters := strings.TrimSpace(paramData["SSD_Drive_Letters"])
		if ssdDriveLetters == "" {
			log.Warnf("Creating storage.config: profile %+v has SSD_Drive_Prefix parameter, but no SSD_Drive_Letters; creating anyway")
		}
		text += storageDotConfigVolumeText(ssdDrivePrefix, ssdDriveLetters, nextVolume)
		nextVolume++
	}
	return text, nil
}

func urlSigDotConfig(tx *sql.Tx, cfg *config.Config, profile ProfileData, fileName string) (string, error) {
	sep := " = "
	if s, ok := separators()[fileName]; ok {
		sep = s
	}

	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}

	urlSigKeys, hasURLSigKeys, err := riaksvc.GetURLSigKeysFromConfigFileKey(tx, cfg.RiakAuthOptions, fileName)
	if err != nil {
		return "", errors.New("getting url sig keys from Riak: " + err.Error())
	}
	if !hasURLSigKeys {
		return "", nil // TODO verify? Perl seems to return without returning its $text
	}

	for key, val := range urlSigKeys {
		text += key + sep + val + "\n"
	}
	return text, nil
}

func uriSigningDotConfig(tx *sql.Tx, cfg *config.Config, fileName string) (string, error) {
	riakKey := strings.TrimSuffix(strings.TrimPrefix(fileName, "uri_signing_"), ".config")
	keys, hasKeys, err := riaksvc.GetURISigningKeysRaw(tx, cfg.RiakAuthOptions, riakKey)
	if err != nil {
		return "", errors.New("getting uri signing keys from Riak: " + err.Error())
	}
	if !hasKeys {
		return "", nil // TODO verify? Perl seems to return without returning its $text
	}
	return string(keys), nil
}

func volumeDotConfigVolumeText(volume string, numVolumes int) string {
	return "volume=" + volume + " scheme=http size=" + strconv.Itoa(100/numVolumes) + "%\n"
}

func volumeDotConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	paramData, err := GetProfileParamData(tx, profile.ID, "storage.config") // volume.config is based on the storage.config params
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	numVolumes := getNumVolumes(paramData)

	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}

	text += "# TRAFFIC OPS NOTE: This is running with forced volumes - the size is irrelevant\n"
	nextVolume := 1
	if drivePrefix := paramData["Drive_Prefix"]; drivePrefix != "" {
		text += volumeDotConfigVolumeText(strconv.Itoa(nextVolume), numVolumes)
		nextVolume++
	}
	if ramDrivePrefix := paramData["RAM_Drive_Prefix"]; ramDrivePrefix != "" {
		text += volumeDotConfigVolumeText(strconv.Itoa(nextVolume), numVolumes)
		nextVolume++
	}
	if ssdDrivePrefix := paramData["SSD_Drive_Prefix"]; ssdDrivePrefix != "" {
		text += volumeDotConfigVolumeText(strconv.Itoa(nextVolume), numVolumes)
		nextVolume++
	}
	return text, nil
}

func getNumVolumes(paramData map[string]string) int {
	num := 0
	drivePrefixes := []string{"Drive_Prefix", "SSD_Drive_Prefix", "RAM_Drive_Prefix"}
	for _, pre := range drivePrefixes {
		if _, ok := paramData[pre]; ok {
			num++
		}
	}
	return num
}

func getHostPortFromURI(uriStr string) (string, string) {
	originFQDN := uriStr
	originFQDN = strings.TrimPrefix(originFQDN, "http://")
	originFQDN = strings.TrimPrefix(originFQDN, "https://")

	slashPos := strings.Index(originFQDN, "/")
	if slashPos != -1 {
		originFQDN = originFQDN[:slashPos]
	}
	portPos := strings.Index(originFQDN, ":")
	portStr := ""
	if portPos != -1 {
		portStr = originFQDN[portPos+1:]
		originFQDN = originFQDN[:portPos]
	}
	return originFQDN, portStr
}

func genericProfileConfig(tx *sql.Tx, profile ProfileData, fileName string) (string, error) {
	sep := " = "
	if s, ok := separators()[fileName]; ok {
		sep = s
	}
	profileParamData, err := GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}
	text, err := headerComment(tx, profile.Name)
	if err != nil {
		return "", err
	}
	for name, val := range profileParamData {
		name = trimParamUnderscoreNumSuffix(name)
		text += name + sep + val + "\n"
	}
	return text, nil
}

// trimParamUnderscoreNumSuffix removes any trailing "__[0-9]+" and returns the trimmed string.
// TODO unit test
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
