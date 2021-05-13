package t3cutil

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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/toreq"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/torequtil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type TCCfg struct {
	CacheHostName string
	GetData       string
	TOClient      *toreq.TOClient
	TOInsecure    bool
	TOTimeoutMS   time.Duration
	TOPass        string
	TOUser        string
	TOURL         *url.URL
	UserAgent     string

	// TODisableProxy is whether to not use a configured Traffic Ops Proxy.
	// This is only used by WriteConfig, which is the only command that makes enough requests to matter.
	TODisableProxy bool

	// RevalOnly is whether to only fetch config data necessary to revalidate, versus all data necessary to generate config. This is only used by WriteConfig
	RevalOnly bool
}

func GetDataFuncs() map[string]func(TCCfg, io.Writer) error {
	return map[string]func(TCCfg, io.Writer) error{
		`update-status`: WriteServerUpdateStatus,
		`packages`:      WritePackages,
		`chkconfig`:     WriteChkconfig,
		`system-info`:   WriteSystemInfo,
		`statuses`:      WriteStatuses,
		`config`:        WriteConfig,
	}
}

func GetServerUpdateStatus(cfg TCCfg) (*tc.ServerUpdateStatus, error) {
	status, _, err := cfg.TOClient.GetServerUpdateStatus(tc.CacheName(cfg.CacheHostName))
	if err != nil {
		return nil, errors.New("getting server '" + cfg.CacheHostName + "' update status: " + err.Error())
	}
	return &status, nil
}

func WriteData(cfg TCCfg) error {
	log.Infoln("Getting data '" + cfg.GetData + "'")
	dataF, ok := GetDataFuncs()[cfg.GetData]
	if !ok {
		return errors.New("unknown data request '" + cfg.GetData + "'")
	}
	return dataF(cfg, os.Stdout)
}

const SystemInfoParamConfigFile = `global`

// WriteSystemInfo writes the "system info" to output.
//
// This is the same info at /api/1.x/system/info, which is actually just all Parameters with the config_file 'global'.
// Note this is different than the more common "global parameters", which usually refers to all Parameters on the Profile named 'GLOBAL'.
//
// This is identical to the /api/1.x/system/info endpoint, except it does not include a '{response: {parameters:' wrapper.
//
func WriteSystemInfo(cfg TCCfg, output io.Writer) error {
	paramArr, _, err := cfg.TOClient.GetConfigFileParameters(SystemInfoParamConfigFile)
	if err != nil {
		return errors.New("getting system info parameters: " + err.Error())
	}
	params := map[string]string{}
	for _, param := range paramArr {
		params[param.Name] = param.Value
	}
	if err := json.NewEncoder(output).Encode(params); err != nil {
		return errors.New("encoding system info parameters: " + err.Error())
	}
	return nil
}

// WriteStatuses writes the Traffic Ops statuses to output.
// Note this is identical to /api/1.x/statuses except it omits the '{response:' wrapper.
func WriteStatuses(cfg TCCfg, output io.Writer) error {
	statuses, _, err := cfg.TOClient.GetStatuses()
	if err != nil {
		return errors.New("getting statuses: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(statuses); err != nil {
		return errors.New("encoding statuses: " + err.Error())
	}
	return nil
}

// WriteUpdateStatus writes the Traffic Ops server update status to output.
// Note this is identical to /api/1.x/servers/name/update_status except it omits the '[]' wrapper.
func WriteServerUpdateStatus(cfg TCCfg, output io.Writer) error {
	status, err := GetServerUpdateStatus(cfg)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(output).Encode(status); err != nil {
		return errors.New("encoding server update status: " + err.Error())
	}
	return nil
}

// WriteORTServerPackages writes the packages for serverName to output.
// Note this is identical to /ort/serverName/packages.
func WritePackages(cfg TCCfg, output io.Writer) error {
	packages, err := GetPackages(cfg)
	if err != nil {
		return errors.New("getting ORT server packages: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(packages); err != nil {
		return errors.New("writing packages: " + err.Error())
	}
	return nil
}

func GetPackages(cfg TCCfg) ([]Package, error) {
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName))
	if err != nil {
		return nil, errors.New("getting server: " + err.Error())
	} else if server.Profile == nil {
		return nil, errors.New("getting server: nil profile")
	}
	params, _, err := cfg.TOClient.GetServerProfileParameters(*server.Profile)
	if err != nil {
		return nil, errors.New("getting server profile '" + *server.Profile + "' parameters: " + err.Error())
	}
	packages := []Package{}
	for _, param := range params {
		if param.ConfigFile != atscfg.PackagesParamConfigFile {
			continue
		}
		packages = append(packages, Package{Name: param.Name, Version: param.Value})
	}
	return packages, nil
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// WriteChkconfig writes the chkconfig for cfg.CacheHostName to output.
// Note this is identical to /ort/serverName/chkconfig.
func WriteChkconfig(cfg TCCfg, output io.Writer) error {
	chkconfig, err := GetChkconfig(cfg)
	if err != nil {
		return errors.New("getting chkconfig: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(chkconfig); err != nil {
		return errors.New("writing chkconfig: " + err.Error())
	}
	return nil
}

func GetChkconfig(cfg TCCfg) ([]ChkConfigEntry, error) {
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName))
	if err != nil {
		return nil, errors.New("getting server: " + err.Error())
	} else if server.Profile == nil {
		return nil, errors.New("getting server: nil profile")
	}
	params, _, err := cfg.TOClient.GetServerProfileParameters(*server.Profile)
	if err != nil {
		return nil, errors.New("getting server profile '" + *server.Profile + "' parameters: " + err.Error())
	}
	chkconfig := []ChkConfigEntry{}
	for _, param := range params {
		if param.ConfigFile != atscfg.ChkconfigParamConfigFile {
			continue
		}
		chkconfig = append(chkconfig, ChkConfigEntry{Name: param.Name, Val: param.Value})
	}
	return chkconfig, nil
}

type ChkConfigEntry struct {
	Name string `json:"name"`
	Val  string `json:"value"`
}

// SetUpdateStatus sets the queue and reval status of serverName in Traffic Ops.
func SetUpdateStatus(cfg TCCfg, serverName string, queue bool, revalPending bool) error {
	reqInf, err := cfg.TOClient.C.SetUpdateServerStatuses(string(serverName), &queue, &revalPending)
	if err != nil {
		return errors.New("setting update statuses (Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "'): " + err.Error())
	}
	return nil
}

// setUpdateStatusLegacy sets the queue and reval status of serverName in Traffic Ops,
// using the legacy pre-2.0 /update endpoint.
func setUpdateStatusLegacy(cfg TCCfg, serverName tc.CacheName, queue bool, revalPending bool) error {
	path := `/update/` + string(serverName) + `?updated=` + jsonBoolStr(queue) + `&reval_updated=` + jsonBoolStr(revalPending)
	// C and RawRequest should generally never be used, but the alternatve here is to manually get the cookie and do an http.Get. We need to hit a non-API endpoint, no API endpoint exists for what we need.
	resp, _, err := cfg.TOClient.C.RawRequest(http.MethodPost, path, nil)
	if err != nil {
		return errors.New("setting update statuses: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyBts, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("Traffic Ops returned %v %v", resp.StatusCode, string(bodyBts))
		}
		return fmt.Errorf("Traffic Ops returned %v (error reading body)", resp.StatusCode)
	}
	return nil
}

func jsonBoolStr(b bool) string {
	if b {
		return `true`
	}
	return `false`
}

// WriteConfig writes the Traffic Ops data necessary to generate config to output.
func WriteConfig(cfg TCCfg, output io.Writer) error {
	cfgData, err := GetConfigData(cfg.TOClient, cfg.TODisableProxy, cfg.CacheHostName, cfg.RevalOnly)
	if err != nil {
		return errors.New("getting statuses: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(cfgData); err != nil {
		return errors.New("encoding config data: " + err.Error())
	}
	return nil
}
