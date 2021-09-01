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
		return nil, fmt.Errorf("getting server '%s' update status: %w", cfg.CacheHostName, err)
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
		return fmt.Errorf("getting system info parameters: %w", err)
	}
	params := map[string]string{}
	for _, param := range paramArr {
		params[param.Name] = param.Value
	}
	if err := json.NewEncoder(output).Encode(params); err != nil {
		return fmt.Errorf("encoding system info parameters: %w", err)
	}
	return nil
}

// WriteStatuses writes the Traffic Ops statuses to output.
// Note this is identical to /statuses except it omits the '{response:'
// wrapper.
func WriteStatuses(cfg TCCfg, output io.Writer) error {
	statuses, _, err := cfg.TOClient.GetStatuses(nil)
	if err != nil {
		return fmt.Errorf("getting statuses: %w", err)
	}
	if err := json.NewEncoder(output).Encode(statuses); err != nil {
		return fmt.Errorf("encoding statuses: %w", err)
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
		return fmt.Errorf("encoding server update status: %w", err)
	}
	return nil
}

// WriteORTServerPackages writes the packages for serverName to output.
// Note this is identical to /ort/serverName/packages.
func WritePackages(cfg TCCfg, output io.Writer) error {
	packages, err := GetPackages(cfg)
	if err != nil {
		return fmt.Errorf("getting ORT server packages: %w", err)
	}
	if err := json.NewEncoder(output).Encode(packages); err != nil {
		return fmt.Errorf("writing packages: %w", err)
	}
	return nil
}

func GetPackages(cfg TCCfg) ([]Package, error) {
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName), nil)
	if err != nil {
		return nil, fmt.Errorf("getting server: %w", err)
	} else if len(server.Profiles) == 0 {
		return nil, errors.New("getting server: nil profile")
	} else if server.HostName == "" {
		return nil, errors.New("getting server: nil hostName")
	}
	allPackageParams, reqInf, err := cfg.TOClient.GetConfigFileParameters(atscfg.PackagesParamConfigFile, nil)
	log.Infoln(toreq.RequestInfoStr(reqInf, "GetPackages.GetConfigFileParameters("+atscfg.PackagesParamConfigFile+")"))
	if err != nil {
		return nil, fmt.Errorf("getting server '%s' package parameters: %w", server.HostName, err)
	}

	serverPackageParams, err := atscfg.GetServerParameters(server, allPackageParams)
	if err != nil {
		return nil, fmt.Errorf("calculating server '%s' package parameters: %w", server.HostName, err)
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
		return fmt.Errorf("getting chkconfig: %w", err)
	}
	if err := json.NewEncoder(output).Encode(chkconfig); err != nil {
		return fmt.Errorf("writing chkconfig: %w", err)
	}
	return nil
}

func GetChkconfig(cfg TCCfg) ([]ChkConfigEntry, error) {
	server, _, err := cfg.TOClient.GetServerByHostName(string(cfg.CacheHostName), nil)
	if err != nil {
		return nil, fmt.Errorf("getting server: %w", err)
	} else if len(server.Profiles) == 0 {
		return nil, errors.New("getting server: nil profile")
	} else if server.HostName == "" {
		return nil, errors.New("getting server: nil hostName")
	}

	allChkconfigParams, reqInf, err := cfg.TOClient.GetConfigFileParameters(atscfg.ChkconfigParamConfigFile, nil)
	log.Infoln(toreq.RequestInfoStr(reqInf, "GetChkconfig.GetConfigFileParameters("+atscfg.ChkconfigParamConfigFile+")"))
	if err != nil {
		return nil, fmt.Errorf("getting server '%s' chkconfig parameters: %w", server.HostName, err)
	}

	serverChkconfigParams, err := atscfg.GetServerParameters(server, allChkconfigParams)
	if err != nil {
		return nil, fmt.Errorf("calculating server '%s' chkconfig parameters: %w", server.HostName, err)
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
		return fmt.Errorf("setting update statuses (Traffic Ops '%s'): %w", torequtil.MaybeIPStr(reqInf.RemoteAddr), err)
	}
	return nil
}

// WriteConfig writes the Traffic Ops data necessary to generate config to output.
func WriteConfig(cfg TCCfg, output io.Writer) error {
	cfgData, err := GetConfigData(cfg.TOClient, cfg.TODisableProxy, cfg.CacheHostName, cfg.RevalOnly, cfg.OldCfg, cfg.T3CVersion)
	if err != nil {
		return fmt.Errorf("getting config data: %w", err)
	}
	if err := json.NewEncoder(output).Encode(cfgData); err != nil {
		return fmt.Errorf("encoding config data: %w", err)
	}
	return nil
}
