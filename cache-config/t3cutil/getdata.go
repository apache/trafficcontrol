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
	"io"
	"net/url"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
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

	// OldCfg is the previously fetched ConfigData, for 'config' requests. May be nil.
	OldCfg *ConfigData

	// T3CVersion is the version of the t3c app ecosystem
	// This value will be the same for any t3c app.
	T3CVersion string
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

func GetServerUpdateStatus(cfg TCCfg) (*atscfg.ServerUpdateStatus, error) {
	status, _, err := cfg.TOClient.GetServerUpdateStatus(tc.CacheName(cfg.CacheHostName), nil)
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
// This is the same info as the Traffic Ops API's /system/info endpoint, which
// is actually just all Parameters with the config_file 'global'.
// Note this is different than the more common "global parameters", which
// usually refers to all Parameters on the Profile named 'GLOBAL'.
//
// This is identical to the /system/info endpoint, except it does not include a
//
//	'{response: {parameters:' wrapper.
func WriteSystemInfo(cfg TCCfg, output io.Writer) error {
	paramArr, _, err := cfg.TOClient.GetConfigFileParameters(SystemInfoParamConfigFile, nil)
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
// Note this is identical to /statuses except it omits the '{response:'
// wrapper.
func WriteStatuses(cfg TCCfg, output io.Writer) error {
	statuses, _, err := cfg.TOClient.GetStatuses(nil)
	if err != nil {
		return errors.New("getting statuses: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(statuses); err != nil {
		return errors.New("encoding statuses: " + err.Error())
	}
	return nil
}

// WriteServerUpdateStatus writes the Traffic Ops server update status to
// output.
// Note this is identical to the Traffic Ops API's
// /servers/{{host name}}/update_status endpoint except it omits the '[]'
// wrapper.
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
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName), nil)
	if err != nil {
		return nil, errors.New("getting server: " + err.Error())
	} else if len(server.Profiles) == 0 {
		return nil, errors.New("getting server: nil profile")
	} else if server.HostName == "" {
		return nil, errors.New("getting server: nil hostName")
	}
	allPackageParams, reqInf, err := cfg.TOClient.GetConfigFileParameters(atscfg.PackagesParamConfigFile, nil)
	log.Infoln(toreq.RequestInfoStr(reqInf, "GetPackages.GetConfigFileParameters("+atscfg.PackagesParamConfigFile+")"))
	if err != nil {
		return nil, errors.New("getting server '" + server.HostName + "' package parameters: " + err.Error())
	}

	serverPackageParams, err := atscfg.GetServerParameters(server, allPackageParams)
	if err != nil {
		return nil, errors.New("calculating server '" + server.HostName + "' package parameters: " + err.Error())
	}

	packages := []Package{}
	for _, param := range serverPackageParams {
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
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName), nil)
	if err != nil {
		return nil, errors.New("getting server: " + err.Error())
	} else if len(server.Profiles) == 0 {
		return nil, errors.New("getting server: nil profile")
	} else if server.HostName == "" {
		return nil, errors.New("getting server: nil hostName")
	}

	allChkconfigParams, reqInf, err := cfg.TOClient.GetConfigFileParameters(atscfg.ChkconfigParamConfigFile, nil)
	log.Infoln(toreq.RequestInfoStr(reqInf, "GetChkconfig.GetConfigFileParameters("+atscfg.ChkconfigParamConfigFile+")"))
	if err != nil {
		return nil, errors.New("getting server '" + server.HostName + "' chkconfig parameters: " + err.Error())
	}

	serverChkconfigParams, err := atscfg.GetServerParameters(server, allChkconfigParams)
	if err != nil {
		return nil, errors.New("calculating server '" + server.HostName + "' chkconfig parameters: " + err.Error())
	}

	chkconfig := []ChkConfigEntry{}
	for _, param := range serverChkconfigParams {
		chkconfig = append(chkconfig, ChkConfigEntry{Name: param.Name, Val: param.Value})
	}
	return chkconfig, nil
}

type ChkConfigEntry struct {
	Name string `json:"name"`
	Val  string `json:"value"`
}

// SetUpdateStatus sets the queue and reval status of serverName in Traffic Ops.
func SetUpdateStatus(cfg TCCfg, serverName tc.CacheName, configApply, revalApply *time.Time) error {
	// TODO need to move to toreq, add fallback
	reqInf, err := cfg.TOClient.SetServerUpdateStatus(serverName, configApply, revalApply)
	if err != nil {
		return errors.New("setting update statuses (Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "'): " + err.Error())
	}
	return nil
}

// SetUpdateStatusCompat sets the queue and reval status of serverName in Traffic Ops.
// *** Compatability requirement until ATC (v7.0+) is deployed with the timestamp features
/*func SetUpdateStatusCompat(cfg TCCfg, serverName tc.CacheName, configApply, revalApply *time.Time, configApplyBool, revalApplyBool *bool) error {
	// TODO need to move to toreq, add fallback
	reqInf, err := cfg.TOClient.SetServerUpdateStatusBoolCompat(serverName, configApply, revalApply, configApplyBool, revalApplyBool)
	if err != nil {
		return errors.New("setting update statuses (Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "'): " + err.Error())
	}
	return nil
}*/

// WriteConfig writes the Traffic Ops data necessary to generate config to output.
func WriteConfig(cfg TCCfg, output io.Writer) error {
	cfgData, err := GetConfigData(cfg.TOClient, cfg.TODisableProxy, cfg.CacheHostName, cfg.RevalOnly, cfg.OldCfg, cfg.T3CVersion)
	if err != nil {
		return errors.New("getting config data: " + err.Error())
	}
	if err := json.NewEncoder(output).Encode(cfgData); err != nil {
		return errors.New("encoding config data: " + err.Error())
	}
	return nil
}
