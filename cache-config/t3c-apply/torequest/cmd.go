package torequest

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

// cmd.go has funcs to call t3c sub-command apps.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// generate runs t3c-generate and returns the result.
func generate(cfg config.Cfg) ([]byte, error) {
	args := []string{
		"--dir=" + config.TSConfigDir,
		"--traffic-ops-timeout-milliseconds=" + strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--traffic-ops-disable-proxy=" + strconv.FormatBool(cfg.ReverseProxyDisable),
		"--traffic-ops-user=" + cfg.TOUser,
		"--traffic-ops-password=" + cfg.TOPass,
		"--traffic-ops-url=" + cfg.TOURL,
		"--cache-host-name=" + cfg.CacheHostName,
		"--log-location-error=" + outToErr(cfg.LogLocationErr),
		"--log-location-info=" + outToErr(cfg.LogLocationInfo),
		"--log-location-warning=" + outToErr(cfg.LogLocationWarn),
	}
	if cfg.TOInsecure == true {
		args = append(args, "--traffic-ops-insecure")
	}
	if cfg.DNSLocalBind {
		args = append(args, "--dns-local-bind")
	}
	if cfg.DefaultClientEnableH2 != nil {
		args = append(args, "--default-client-enable-h2="+strconv.FormatBool(*cfg.DefaultClientEnableH2))
	}
	if cfg.DefaultClientTLSVersions != nil {
		args = append(args, "--default-client-tls-versions="+*cfg.DefaultClientTLSVersions+"")
	}
	if cfg.RunMode == config.Revalidate {
		args = append(args, "--revalidate-only")
	}
	args = append(args, "--via-string-release="+strconv.FormatBool(!cfg.OmitViaStringRelease))
	args = append(args, "--disable-parent-config-comments="+strconv.FormatBool(cfg.DisableParentConfigComments))

	stdOut, stdErr, code := t3cutil.Do(config.GenerateCmd, args...)
	if code != 0 {
		return nil, fmt.Errorf("t3c-generate returned non-zero exit code %v stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if len(bytes.TrimSpace(stdErr)) > 0 {
		log.Warnln(`t3c-generate stderr start` + "\n" + string(stdErr))
		log.Warnln(`t3c-generate stderr end`)
	}
	return stdOut, nil
}

func getStatuses(cfg config.Cfg) ([]string, error) {
	statuses := []tc.StatusNullable{}
	if err := requestJSON(cfg, "statuses", &statuses); err != nil {
		return nil, errors.New("requesting json: " + err.Error())
	}
	sl := []string{}
	for val := range statuses {
		if statuses[val].Name != nil {
			sl = append(sl, *statuses[val].Name)
		}
	}
	return sl, nil
}

func getChkconfig(cfg config.Cfg) ([]map[string]string, error) {
	result := []map[string]string{}
	if err := requestJSON(cfg, "chkconfig", &result); err != nil {
		return nil, errors.New("requesting json: " + err.Error())
	}
	return result, nil
}

func getUpdateStatus(cfg config.Cfg) (*tc.ServerUpdateStatus, error) {
	status := tc.ServerUpdateStatus{}
	if err := requestJSON(cfg, "update-status", &status); err != nil {
		return nil, errors.New("requesting json: " + err.Error())
	}
	return &status, nil
}

func getSystemInfo(cfg config.Cfg) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if err := requestJSON(cfg, "system-info", &result); err != nil {
		return nil, errors.New("requesting json: " + err.Error())
	}
	return result, nil
}

func getPackages(cfg config.Cfg) ([]Package, error) {
	pkgs := []Package{}
	if err := requestJSON(cfg, "packages", &pkgs); err != nil {
		return nil, errors.New("requesting json: " + err.Error())
	}
	return pkgs, nil
}

// sendUpdate updates the given cache's queue update and reval status in Traffic Ops.
// Note the statuses are the value to be set, not whether to set the value.
func sendUpdate(cfg config.Cfg, updateStatus bool, revalStatus bool) error {
	stdOut, stdErr, code := t3cutil.Do(`t3c-update`,
		"--traffic-ops-timeout-milliseconds="+strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--traffic-ops-user="+cfg.TOUser,
		"--traffic-ops-password="+cfg.TOPass,
		"--traffic-ops-url="+cfg.TOURL,
		"--log-location-error="+outToErr(cfg.LogLocationErr),
		"--log-location-info="+outToErr(cfg.LogLocationInfo),
		"--cache-host-name="+cfg.CacheHostName,
		"--set-update-status="+strconv.FormatBool(updateStatus),
		"--set-reval-status="+strconv.FormatBool(revalStatus),
	)
	if code != 0 {
		return fmt.Errorf("t3c-update returned non-zero exit code %v stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if len(bytes.TrimSpace(stdErr)) > 0 {
		log.Warnf("t3c-request returned code 0 but stderr '%v'", string(stdErr)) // usually warnings
	}
	log.Infoln("t3c-update succeeded")
	return nil
}

// diff calls t3c-diff to diff the given new file and the file on disk. Returns whether they're different.
// Logs the difference.
// If the file on disk doesn't exist, returns true and logs the entire file as a diff.
func diff(cfg config.Cfg, newFile string, fileLocation string) (bool, error) {
	log.Warnf("DEBUG diff calling location '" + fileLocation + "'")
	stdOut, stdErr, code := t3cutil.DoInput(newFile, `t3c-diff`, `stdin`, fileLocation)
	if code > 1 {
		return false, fmt.Errorf("t3c-update returned error code %v stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if len(bytes.TrimSpace(stdErr)) > 0 {
		log.Warnf("t3c-request returned non-error code %v but stderr '%v'", code, string(stdErr))
	}

	if code == 0 {
		log.Infof("All lines match TrOps for config file: %s\n", fileLocation)
		return false, nil // 0 is only returned if there's no diff
	}
	// code 1 means a diff, difference text will be on stdout

	lines := strings.Split(string(stdOut), "\n")
	log.Infoln("file '" + fileLocation + "' changes begin")
	for _, line := range lines {
		log.Infoln("diff: " + line)
	}
	log.Infoln("file '" + fileLocation + "' changes end")

	return true, nil
}

// verify calls t3c-verify to verify the given cfgFile.
// The cfgFile should be the full text of either a plugin.config or remap.config.
// Returns nil if t3c-verify returned no errors found, or the error found if any.
func verify(cfg config.Cfg, cfgFile string) error {
	stdOut, stdErr, code := t3cutil.DoInput(cfgFile, `t3c-verify`,
		"--log-location-error="+outToErr(cfg.LogLocationErr),
		"--log-location-info="+outToErr(cfg.LogLocationInfo),
		"--log-location-debug="+outToErr(cfg.LogLocationDebug),
	)
	if code != 0 {
		log.Errorf(`verify errors start
` + string(stdOut))
		log.Errorf(`verify errors end`)
		if strings.TrimSpace(string(stdErr)) != "" {
			log.Errorf(`verify output start
` + string(stdErr))
			log.Errorf(`verify output end`)
		}
		return fmt.Errorf("%d plugins failed to verify. See log for details.", code)
	}
	if len(bytes.TrimSpace(stdErr)) > 0 {
		log.Warnf("t3c-verify returned non-error code %v but stderr '%v'", code, string(stdErr))
	}
	if len(bytes.TrimSpace(stdOut)) > 0 {
		log.Warnf("t3c-verify returned non-error code %v but output '%v'", code, string(stdOut))
	}
	return nil
}

// requestJSON calls t3c-request with the given command, and deserializes the result as JSON into obj.
func requestJSON(cfg config.Cfg, command string, obj interface{}) error {
	stdOut, stdErr, code := t3cutil.Do(`t3c-request`,
		"--traffic-ops-timeout-milliseconds="+strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--traffic-ops-user="+cfg.TOUser,
		"--traffic-ops-password="+cfg.TOPass,
		"--traffic-ops-url="+cfg.TOURL,
		"--cache-host-name="+cfg.CacheHostName,
		"--log-location-error="+outToErr(cfg.LogLocationErr),
		"--log-location-info="+outToErr(cfg.LogLocationInfo),
		`--get-data=`+command,
	)
	if code != 0 {
		return fmt.Errorf("t3c-request returned non-zero exit code %v stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if len(bytes.TrimSpace(stdErr)) > 0 {
		log.Warnf("t3c-request returned code 0 but stderr '%v'", string(stdErr)) // usually warnings
	}
	if err := json.Unmarshal(stdOut, obj); err != nil {
		return errors.New("unmarshalling '" + string(stdOut) + "': " + err.Error())
	}
	return nil
}

// outToErr returns stderr if logLocation is stdout, otherwise returns logLocation unchanged.
// This is a helper to avoid logging to stdout for commands whose output is on stdout.
func outToErr(logLocation string) string {
	if logLocation == "stdout" {
		return "stderr"
	}
	return logLocation
}
