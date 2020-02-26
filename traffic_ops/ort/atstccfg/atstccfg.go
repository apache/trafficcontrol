// atstccfg is a tool for generating configuration files server-side on ATC cache servers.
//
// Warning: atstccfg does not have a stable command-line interface, it can and will change without warning. Scripts should avoid calling it for the time being.
//
// Usage:
//
// 	atstccfg [-u TO_URL] [-U TO_USER] [-P TO_PASSWORD] [-n] [-r N] [-e ERROR_LOCATION] [-w WARNING_LOCATION] [-i INFO_LOCATION] [-g] [-s] [-t TIMEOUT] [-a MAX_AGE] [-l] [-v] [-h]
//
// The available options are:
//
// 	-a, --cache-file-max-age-seconds                                Sets the maximum age - in seconds - a cached response can be in order to be considered "fresh" - older files will be re-generated and cached. Default: 60
// 	-e ERROR_LOCATION, --log-location-error ERROR_LOCATION          The file location to which to log errors. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
// 	-g, --print-generated-files                                     If given, the names of files generated (and not proxied to Traffic Ops) will be printed to stdout, then atstccfg will exit.
// 	-h, --help                                                      Print usage information and exit.
// 	-i INFO_LOCATION, --log-location-info INFO_LOCATION             The file location to which to log information messages. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
// 	-l, --list-plugins                                              List the loaded plugins and then exit.
// 	-n, --no-cache                                                  If given, existing cache files will not be used. Cache files will still be created, existing ones just won't be used.
// 	-P TO_PASSWORD                                                  Authenticate using this password - if not given, atstccfg will attempt to use the value of the TO_PASS environment variable
// 	-r N, --num-retries N                                           The number of times to retry getting a file if it fails. Default: 5
// 	-s, --traffic-ops-insecure                                      If given, SSL certificate errors will be ignored when communicating with Traffic Ops. NOT RECOMMENDED FOR PRODUCTION ENVIRONMENTS.
// 	-t, --traffic-ops-timeout-milliseconds                          Sets the timeout - in milliseconds - for requests made to Traffic Ops. Default: 10000
// 	-u TO_URL                                                       Request this URL, e.g. 'https://trafficops.infra.ciab.test/servers/edge/configfiles/ats'
// 	-U TO_USER                                                      Authenticate as the user TO_USER - if not given, atstccfg will attempt to use the value of the TO_USER environment variable
// 	-v, --version                                                   Print version information and exit.
// 	-w WARNING_LOCATION, --log-location-warning WARNING_LOCATION    The file location to which to log warnings. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
//
// atstccfg caches generated files in /tmp/atstccfg_cache/ for re-use.

package main

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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/cfgfile"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/plugin"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/toreq"
)

func main() {
	cfg, err := config.GetCfg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(config.ExitCodeErrGeneric)
	}

	if cfg.ListPlugins {
		fmt.Println(strings.Join(plugin.List(), "\n"))
		os.Exit(0)
	}

	if cfg.PrintGeneratedFiles {
		fmt.Println(strings.Join(GetGeneratedFilesList(), "\n"))
		os.Exit(config.ExitCodeSuccess)
	}

	plugins := plugin.Get(cfg)
	plugins.OnStartup(plugin.StartupData{Cfg: cfg})

	log.Infoln("URL: '" + cfg.TOURL.String() + "' User: '" + cfg.TOUser + "' Pass len: '" + strconv.Itoa(len(cfg.TOPass)) + "'")

	toFQDN := cfg.TOURL.Scheme + "://" + cfg.TOURL.Host
	log.Infoln("TO FQDN: '" + toFQDN + "'")
	log.Infoln("TO URL: '" + cfg.TOURL.String() + "'")

	toClient, toIP, err := toclient.LoginWithAgent(toFQDN, cfg.TOUser, cfg.TOPass, cfg.TOInsecure, config.UserAgent, false, cfg.TOTimeout)
	if err != nil {
		log.Errorln("Logging in to Traffic Ops '" + toreq.MaybeIPStr(toIP) + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	tccfg := config.TCCfg{Cfg: cfg, TOClient: &toClient}

	onReqData := plugin.OnRequestData{Cfg: tccfg}
	if handled := plugins.OnRequest(onReqData); handled {
		return
	}

	configs, err := GetAllConfigs(tccfg)
	if err != nil {
		log.Errorln("Getting config for'" + cfg.CacheHostName + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	for _, cfgFile := range configs {
		path := filepath.Join(cfg.OutputDir, cfgFile.FileNameOnDisk)
		if err := ioutil.WriteFile(path, []byte(cfgFile.Text), 0644); err != nil {
			log.Errorln("Getting config for'" + path + "': " + err.Error())
		}
	}
	os.Exit(config.ExitCodeSuccess)
}

func GetGeneratedFilesList() []string {
	names := []string{}
	for scope, fileFuncs := range ConfigFileFuncs() {
		for cfgFile, _ := range fileFuncs {
			names = append(names, scope+"/"+cfgFile)
		}
	}
	names = append(names, "profiles/url_sig_*.config")     // url_sig configs are generated, but not in the funcs because they're not a literal match
	names = append(names, "profiles/uri_signing_*.config") // uri_signing configs are generated, but not in the funcs because they're not a literal match
	names = append(names, "profiles/*")                    // unknown profiles configs are generated, a.k.a. "take-and-bake"
	return names
}

func HTTPCodeToExitCode(httpCode int) int {
	switch httpCode {
	case http.StatusNotFound:
		return config.ExitCodeNotFound
	}
	return config.ExitCodeErrGeneric
}

// GetAllConfigs returns a map[configFileName]configFileText
func GetAllConfigs(cfg config.TCCfg) ([]ATSConfigFile, error) {
	toData, err := cfgfile.GetTOData(cfg)
	if err != nil {
		return nil, errors.New("getting data from traffic ops: " + err.Error())
	}

	meta, err := cfgfile.GetMeta(toData)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	configs := []ATSConfigFile{}
	for _, fi := range meta.ConfigFiles {
		txt, _, err := GetConfigFile(toData, fi)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.APIURI + "': " + err.Error())
		}
		configs = append(configs, ATSConfigFile{ATSConfigMetaDataConfigFile: fi, Text: txt})
	}
	return configs, nil
}

type ATSConfigFile struct {
	tc.ATSConfigMetaDataConfigFile
	Text string
}
