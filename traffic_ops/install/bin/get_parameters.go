
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// get_parameters.go gets the Traffic Ops parameters for the given profile, and prints the JSON for a new profile ready to import, with the given name and description

// Example usage:
// go run get_parameters.go  -t http://mycdn.comcast.net:3000 -u admin -p password -r GLOBAL -n NEW_GLOBAL -d "My New Global Profile" > global.profile.traffic_ops

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var trafficOpsUri string
var trafficOpsUser string
var trafficOpsPass string
var oldProfile string
var newProfile string
var newProfileDesc string

func init() {
	flag.StringVar(&trafficOpsUri, "traffic-ops-uri", "", "Traffic Ops URI, including protocol")
	flag.StringVar(&trafficOpsUri, "t", "", "Traffic Ops URI, including protocol (shorthand)")
	flag.StringVar(&trafficOpsUser, "user", "", "Traffic Ops username")
	flag.StringVar(&trafficOpsUser, "u", "", "Traffic Ops username (shorthand)")
	flag.StringVar(&trafficOpsPass, "password", "", "Traffic Ops password")
	flag.StringVar(&trafficOpsPass, "p", "", "Traffic Ops password (shorthand)")
	flag.StringVar(&oldProfile, "profile", "", "profile type to get")
	flag.StringVar(&oldProfile, "r", "", "profile to get (shorthand)")
	flag.StringVar(&newProfile, "name", "", "name of the profile to create")
	flag.StringVar(&newProfile, "n", "", "name of the profile to create (shorthand)")
	flag.StringVar(&newProfileDesc, "description", "", "description of the profile to create")
	flag.StringVar(&newProfileDesc, "d", "", "description of the profile to create (shorthand)")
	flag.Parse()
}

type JsonTrafficOpsLogin struct {
	User     string `json:"u"`
	Password string `json:"p"`
}

const LoginCookieName = "mojolicious"

func getLoginCookie(cookies []*http.Cookie) (string, error) {

	for _, cookie := range cookies {
		if cookie.Name != LoginCookieName {
			continue
		}
		return cookie.Value, nil
	}

	return "", errors.New("The " + LoginCookieName + " login cookie was not returned. Login failed?")
}

const apiPath = "/api/1.2/"

// login logs in to the Traffic Ops API and returns the login cookie
func login(uri, user, pass string) (string, error) {
	loginPath := apiPath + "user/login"
	trafficOpsUri = strings.TrimRight(trafficOpsUri, "/")

	loginBytes, err := json.Marshal(struct {
		User string `json:"u"`
		Pass string `json:"p"`
	}{User: trafficOpsUser, Pass: trafficOpsPass})
	if err != nil {
		return "", err
	}

	loginResp, err := http.Post(trafficOpsUri+loginPath, "application/json", bytes.NewReader(loginBytes))
	if err != nil {
		return "", err
	}

	return getLoginCookie(loginResp.Cookies())
}

type JsonParameter struct {
	LastUpdated string `json:"lastUpdated"`
	Value       string `json:"value"`
	Name        string `json:"name"`
	ConfigFile  string `json:"configFile"`
}

// JsonProfile represents the JSON returned by Traffic Ops /api/1.2/parameters/profile
type JsonProfile struct {
	Response []JsonParameter `json:"response"`
}

type JsonImportProfileProfile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JsonImportProfileParameter struct {
	Value      string `json:"value"`
	Name       string `json:"name"`
	ConfigFile string `json:"config_file"`
}

// JsonProfile represents the JSON expected by Traffic Ops /profile/doImport
type JsonImportProfile struct {
	Parameters []JsonImportProfileParameter `json:"parameters"`
	Profile    JsonImportProfileProfile     `json:"profile"`
}

type JsonParametersProfile struct {
	Response []JsonImportProfileParameter `json:"response"`
}

func getLatestProfile(uri, loginCookie, profileType string) (string, error) {
	if profileType == "global" {
		return "GLOBAL", nil
	}

	latestProfileName := "latest_" + profileType
	parametersProfilePath := apiPath + "parameters/profile/GLOBAL.json"

	client := &http.Client{}

	uri = strings.TrimRight(uri, "/")
	req, err := http.NewRequest("GET", uri+parametersProfilePath, nil)
	req.AddCookie(&http.Cookie{Name: LoginCookieName, Value: loginCookie})
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var paramProfile JsonParametersProfile
	err = json.NewDecoder(resp.Body).Decode(&paramProfile)
	if err != nil {
		return "", err
	}

	for _, param := range paramProfile.Response {
		if param.Name == latestProfileName {
			return param.Value, nil
		}
	}

	return "", errors.New("profile not found")
}

func getProfile(uri, loginCookie, profileName string) (JsonProfile, error) {
	profilePath := apiPath + "parameters/profile/"

	client := &http.Client{}

	uri = strings.TrimRight(uri, "/")
	req, err := http.NewRequest("GET", uri+profilePath+profileName+".json", nil)
	req.AddCookie(&http.Cookie{Name: LoginCookieName, Value: loginCookie})
	resp, err := client.Do(req)
	if err != nil {
		return JsonProfile{}, err
	}

	var profile JsonProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return JsonProfile{}, err
	}
	return profile, nil
}

func profileToImportProfile(profile JsonProfile, name string, description string) JsonImportProfile {
	var params []JsonImportProfileParameter
	for _, param := range profile.Response {
		params = append(params, JsonImportProfileParameter{Name: param.Name, Value: param.Value, ConfigFile: param.ConfigFile})
	}
	return JsonImportProfile{Parameters: params, Profile: JsonImportProfileProfile{Name: name, Description: description}}
}

func validProfileTypes() map[string]struct{} {
	return map[string]struct{}{
		"global":             struct{}{},
		"trafficserver_edge": struct{}{},
		"trafficserver_mid":  struct{}{},
		"traffic_stats":      struct{}{},
		"traffic_monitor":    struct{}{},
		"traffic_vault":      struct{}{},
		"traffic_router":     struct{}{},
	}
}

func processDelete(param JsonParameter) bool {
	configFilePrefixes := map[string]struct{}{
		"url_sig_":     struct{}{},
		"regex_remap_": struct{}{},
		"hdr_rw_":      struct{}{},
	}
	configFiles := map[string]struct{}{
		"dns.zone":              struct{}{},
		"http-log4j.properties": struct{}{},
		"dns-log4j.properties":  struct{}{},
	}
	names := map[string]struct{}{
		"allow_ip":       struct{}{},
		"allow_ip6":      struct{}{},
		"purge_allow_ip": struct{}{},
		"ramdisk_size":   struct{}{},
	}
	nameSuffixes := map[string]struct{}{
		".dnssec.inception": struct{}{},
		"_fwd_proxy":        struct{}{},
		"_graph_url":        struct{}{},
	}
	namePrefixes := map[string]struct{}{
		"visual_status_panel": struct{}{},
		"latest_":             struct{}{},
	}

	for prefix, _ := range configFilePrefixes {
		if strings.HasPrefix(param.ConfigFile, prefix) {
			return true
		}
	}
	for configFile, _ := range configFiles {
		if param.ConfigFile == configFile {
			return true
		}
	}
	for name, _ := range names {
		if param.Name == name {
			return true
		}
	}
	for prefix, _ := range namePrefixes {
		if strings.HasPrefix(param.Name, prefix) {
			return true
		}
	}

	for suffix, _ := range nameSuffixes {
		if strings.HasSuffix(param.Name, suffix) {
			return true
		}
	}

	if strings.HasPrefix(param.ConfigFile, "cacheurl_") && param.ConfigFile != "cacheurl_qstring" {
		return true
	}

	return false
}

func processModify(param JsonParameter) JsonParameter {
	if param.Name == "LogFormat.Format" {
		param.Value = strings.Replace(param.Value, `xmt=\"%<{X-MoneyTrace}cqh>\"`, "", -1)
		return param
	}

	if strings.HasPrefix(param.Name, "cron_ort_syncds_") {
		param.Value = "{{.CronOrtSyncds}}"
		param.Name = "cron_ort_syncds_cdn" // TODO(`cdn` is the name of the CDN; handle different names?)
		return param
	}

	nameValueReplacements := map[string]string{
		"domain_name":                               "{{.DomainName}}",
		"Drive_Prefix":                              "{{.DrivePrefix}}",
		"RAM_Drive_Prefix":                          "{{.RAMDrivePrefix}}",
		"RAM_Drive_Letters":                         "{{.RAMDriveLetters}}",
		"health.connection.timeout":                 "{{.HealthConnectionTimeout}}",
		"health.threshold.loadavg":                  "{{.HealthThresholdLoadavg}}",
		"health.threshold.availableBandwidthInKbps": "{{.HealthThresholdAvailableBandwidthInKbps}}",
		"health.polling.interval":                   "{{.HealthPollingInterval}}",
		"geolocation.polling.url":                   "{{.GeolocationPollingUrl}}",
		"coveragezone.polling.url":                  "{{.CoveragezonePollingUrl}}",
		"tld.soa.admin":                             "{{.TldSoaAdmin}}",
		"geolocation6.polling.url":                  "{{.Geolocation6PollingUrl}}",
		"tm.infourl":                                "{{.TmInfoUrl}}",
		"tm.instance_name":                          "{{.TmInstanceName}}",
		"tm.toolname":                               "{{.TmToolName}}",
		"tm.url":                                    "{{.TmUrl}}",
	}

	for name, valueReplacement := range nameValueReplacements {
		if param.Name != name {
			continue
		}
		param.Value = valueReplacement
		return param
	}
	return param
}

// process removes unnecessary parameters, and
// replaces CDN-specific parameters with text/template delimited template actions
func process(profile JsonProfile) JsonProfile {
	// This could keep track of the index and delete in-place, if performance mattered.
	var newProfile JsonProfile
	for _, param := range profile.Response {
		if processDelete(param) {
			continue
		}
		param = processModify(param)
		newProfile.Response = append(newProfile.Response, param)
	}
	return newProfile

	// type JsonParameter struct {
	// 	LastUpdated string `json:"lastUpdated"`
	// 	Value string `json:"value"`
	// 	Name string `json:"name"`
	// 	ConfigFile string `json:"configFile"`
	// }

}

func printExampleUsage() {
	typeStr := "{"
	for profileType, _ := range validProfileTypes() {
		typeStr += profileType + "|"
	}
	typeStr = typeStr[:len(typeStr)-1] + "}"

	fmt.Println(`Example: get_parameters -d "My Edge Cache" -n "EDGE1_CDN_520" -u admin -p mypass -r ` + typeStr + ` -t http://my-cdn-domain.com:3000 > edge_profile.traffic_ops`)
}

func main() {
	oldProfile = strings.ToLower(oldProfile)

	validTypes := validProfileTypes()
	if _, isProfileTypeValid := validTypes[oldProfile]; trafficOpsUri == "" || trafficOpsUser == "" || trafficOpsPass == "" || !isProfileTypeValid || newProfile == "" {
		flag.Usage()
		fmt.Print("\n")
		printExampleUsage()
		return
	}

	loginCookie, err := login(trafficOpsUri, trafficOpsUser, trafficOpsPass)
	if err != nil {
		fmt.Println(err)
		return
	}

	latestProfile, err := getLatestProfile(trafficOpsUri, loginCookie, oldProfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	profile, err := getProfile(trafficOpsUri, loginCookie, latestProfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	processedProfile := process(profile)
	importProfile := profileToImportProfile(processedProfile, newProfile, newProfileDesc)

	profileBytes, err := json.Marshal(importProfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", profileBytes)
}
