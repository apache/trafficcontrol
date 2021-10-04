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
	"io/ioutil"
	golog "log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type ServerAndConfigs struct {
	ConfigData  json.RawMessage
	ConfigFiles json.RawMessage
}

// generate runs t3c-generate and returns the result.
func generate(cfg config.Cfg) ([]t3cutil.ATSConfigFile, error) {
	configData, err := requestConfig(cfg)
	if err != nil {
		return nil, errors.New("requesting: " + err.Error())
	}
	args := []string{
		"--dir=" + cfg.TsConfigDir,
	}

	if cfg.LogLocationErr == log.LogLocationNull {
		args = append(args, "-s")
	}
	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
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
	if cfg.Files == t3cutil.ApplyFilesFlagReval {
		args = append(args, "--revalidate-only")
	}
	args = append(args, "--via-string-release="+strconv.FormatBool(!cfg.OmitViaStringRelease))
	args = append(args, "--no-outgoing-ip="+strconv.FormatBool(cfg.NoOutgoingIP))
	args = append(args, "--disable-parent-config-comments="+strconv.FormatBool(cfg.DisableParentConfigComments))

	generatedFiles, stdErr, code := t3cutil.DoInput(configData, config.GenerateCmd, args...)
	if code != 0 {
		logSubAppErr(`t3c-generate stdout`, generatedFiles)
		logSubAppErr(`t3c-generate stderr`, stdErr)
		return nil, fmt.Errorf("t3c-generate returned non-zero exit code %v, see log for output", code)
	}
	logSubApp(`t3c-generate`, stdErr)

	preprocessedBytes, err := preprocess(cfg, configData, generatedFiles)
	if err != nil {
		return nil, errors.New("preprocessing config files: " + err.Error())
	}

	allFiles := []t3cutil.ATSConfigFile{}
	if err := json.Unmarshal(preprocessedBytes, &allFiles); err != nil {
		return nil, errors.New("unmarshalling generated files: " + err.Error())
	}

	return allFiles, nil
}

// preprocess takes the to Data from 't3c-request --get-data=config' and the generated files from 't3c-generate', passes them to `t3c-preprocess`, and returns the result.
func preprocess(cfg config.Cfg, configData []byte, generatedFiles []byte) ([]byte, error) {
	args := []string{}

	if cfg.LogLocationErr == log.LogLocationNull {
		args = append(args, "-s")
	}
	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
	}

	cmd := exec.Command(`t3c-preprocess`, args...)
	outbuf := bytes.Buffer{}
	errbuf := bytes.Buffer{}
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.New("getting command pipe: " + err.Error())
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.New("starting command: " + err.Error())
	}

	if _, err := stdinPipe.Write([]byte(`{"data":`)); err != nil {
		return nil, errors.New("writing opening json to input: " + err.Error())
	} else if _, err := stdinPipe.Write(configData); err != nil {
		return nil, errors.New("writing config data to input: " + err.Error())
	} else if _, err := stdinPipe.Write([]byte(`,"files":`)); err != nil {
		return nil, errors.New("writing files key to input: " + err.Error())
	} else if _, err := stdinPipe.Write(generatedFiles); err != nil {
		return nil, errors.New("writing generated files to input: " + err.Error())
	} else if _, err := stdinPipe.Write([]byte(`}`)); err != nil {
		return nil, errors.New("writing closing json input: " + err.Error())
	} else if err := stdinPipe.Close(); err != nil {
		return nil, errors.New("closing stdin writer: " + err.Error())
	}

	code := 0 // if cmd.Wait returns no error, that means the command returned 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, errors.New("error running command: " + err.Error())
		} else {
			code = exitErr.ExitCode()
		}
	}

	stdOut := outbuf.Bytes()
	stdErr := errbuf.Bytes()
	if code != 0 {
		logSubAppErr(`t3c-preprocess stdout`, stdOut)
		logSubAppErr(`t3c-preprocess stderr`, stdErr)
		return nil, fmt.Errorf("t3c-preprocess returned non-zero exit code %v, see log for output", code)
	}
	logSubApp(`t3c-preprocess`, stdErr)
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
	args := []string{
		"--traffic-ops-timeout-milliseconds=" + strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--traffic-ops-user=" + cfg.TOUser,
		"--traffic-ops-password=" + cfg.TOPass,
		"--traffic-ops-url=" + cfg.TOURL,
		"--traffic-ops-insecure=" + strconv.FormatBool(cfg.TOInsecure),
		"--cache-host-name=" + cfg.CacheHostName,
		"--set-update-status=" + strconv.FormatBool(updateStatus),
		"--set-reval-status=" + strconv.FormatBool(revalStatus),
	}

	if cfg.LogLocationErr == log.LogLocationNull {
		args = append(args, "-s")
	}
	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
	}

	if _, used := os.LookupEnv("TO_USER"); !used {
		args = append(args, "--traffic-ops-user="+cfg.TOUser)
	}
	if _, used := os.LookupEnv("TO_PASS"); !used {
		args = append(args, "--traffic-ops-password="+cfg.TOPass)
	}
	if _, used := os.LookupEnv("TO_URL"); !used {
		args = append(args, "--traffic-ops-url="+cfg.TOURL)
	}
	stdOut, stdErr, code := t3cutil.Do(`t3c-update`, args...)
	if code != 0 {
		logSubAppErr(`t3c-update stdout`, stdOut)
		logSubAppErr(`t3c-update stderr`, stdErr)
		return fmt.Errorf("t3c-update returned non-zero exit code %v, see log for output", code)
	}
	logSubApp(`t3c-update`, stdErr)
	log.Infoln("t3c-update succeeded")
	return nil
}

// diff calls t3c-diff to diff the given new file and the file on disk. Returns whether they're different.
// Logs the difference.
// If the file on disk doesn't exist, returns true and logs the entire file as a diff.
func diff(cfg config.Cfg, newFile []byte, fileLocation string, reportOnly bool, perm os.FileMode) (bool, error) {
	args := []string{
		"--file-a=stdin",
		"--file-b=" + fileLocation,
		"--file-mode=" + fmt.Sprintf("%#o", perm),
	}
	
	diffMsg := ""
	stdOut, stdErr, code := t3cutil.DoInput(newFile, `t3c-diff`, args...)

	if code > 1 {
		return false, fmt.Errorf("t3c-diff returned error code %v stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	logSubApp(`t3c-diff`, stdErr)

	if code == 0 {
		diffMsg += fmt.Sprintf("All lines and file permissions match TrOps for config file: %s\n", fileLocation)
		return false, nil // 0 is only returned if there's no diff
	}
	// code 1 means a diff, difference text will be on stdout

	stdOut = bytes.TrimSpace(stdOut) // the shell output includes a trailing newline that isn't part of the diff; remove it
	lines := strings.Split(string(stdOut), "\n")
	diffMsg += "file '" + fileLocation + "' changes begin\n"
	for _, line := range lines {
		diffMsg += "diff: " + line + "\n"
	}
	diffMsg += "file '" + fileLocation + "' changes end" // no trailing newline, becuase we're using log*ln, the last line will get a newline appropriately

	if reportOnly {
		// Create our own info logger, to log the diff.
		// We can't use the logger initialized in the config package because
		// we don't want to log all the other Info logs.
		// But we want the standard log.Info prefix, timestamp, etc.
		reportLocation := os.Stdout
		goLogger := golog.New(reportLocation, log.InfoPrefix, log.InfoFlags)
		for _, line := range strings.Split(diffMsg, "\n") {
			log.Logln(goLogger, line)
		}
	} else {
		for _, line := range strings.Split(diffMsg, "\n") {
			log.Infoln(line)
		}
	}

	return true, nil
}

// checkRefs calls t3c-check-refs to verify the given cfgFile.
// The cfgFile should be the full text of either a plugin.config or remap.config.
// Returns nil if t3c-check-refs returned no errors found, or the error found if any.
func checkRefs(cfg config.Cfg, cfgFile []byte, filesAdding []string) error {
	args := []string{`check`, `refs`,
		"--files-adding=" + strings.Join(filesAdding, ","),
	}
	if cfg.LogLocationErr == log.LogLocationNull {
		args = append(args, "-s")
	}
	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
	}

	stdOut, stdErr, code := t3cutil.DoInput(cfgFile, `t3c`, args...)

	if code != 0 {
		logSubAppErr(`t3c-check-refs stdout`, stdOut)
		logSubAppErr(`t3c-check-refs stderr`, stdErr)
		return fmt.Errorf("%d plugins failed to verify. See log for details.", code)
	}
	logSubApp(`t3c-check-refs stdout`, stdOut)
	logSubApp(`t3c-check-refs stderr`, stdErr)
	return nil
}

// checkReload is a helper for the sub-command t3c-check-reload.
func checkReload(pluginPackagesInstalled []string, changedConfigFiles []string) (t3cutil.ServiceNeeds, error) {
	log.Infof("t3c-check-reload calling with pluginPackagesInstalled '%v' changedConfigFiles '%v'\n", pluginPackagesInstalled, changedConfigFiles)

	changedFiles := []byte(strings.Join(changedConfigFiles, ","))
	installedPlugins := []byte(strings.Join(pluginPackagesInstalled, ","))

	cmd := exec.Command(`t3c-check-reload`)
	outBuf := bytes.Buffer{}
	errBuf := bytes.Buffer{}
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("getting command pipe: " + err.Error())
	}

	if err := cmd.Start(); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("starting command: " + err.Error())
	}

	if _, err := stdinPipe.Write([]byte(`{"changed_files":"`)); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("writing opening json to input: " + err.Error())
	} else if _, err := stdinPipe.Write(changedFiles); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("writing changed files to input: " + err.Error())
	} else if _, err := stdinPipe.Write([]byte(`","installed_plugins":"`)); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("writing installed_plugins key to input: " + err.Error())
	} else if _, err := stdinPipe.Write(installedPlugins); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("writing plugins to input: " + err.Error())
	} else if _, err := stdinPipe.Write([]byte(`"}`)); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("writing closing json input: " + err.Error())
	} else if err := stdinPipe.Close(); err != nil {
		return t3cutil.ServiceNeedsInvalid, errors.New("closing stdin writer: " + err.Error())
	}

	code := 0 // if cmd.Wait returns no error, that means the command returned 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return t3cutil.ServiceNeedsInvalid, errors.New("error running command: " + err.Error())
		} else {
			code = exitErr.ExitCode()
		}
	}

	stdOut := outBuf.Bytes()
	stdErr := errBuf.Bytes()

	if code != 0 {
		logSubAppErr(`t3c-check-reload stdout`, stdOut)
		logSubAppErr(`t3c-check-reload stderr`, stdErr)
		return t3cutil.ServiceNeedsInvalid, fmt.Errorf("t3c-check-reload returned error code %d - see log for details.", code)
	}

	logSubApp(`t3c-check-reload`, stdErr)

	needs := t3cutil.StrToServiceNeeds(strings.TrimSpace(string(stdOut)))
	if needs == t3cutil.ServiceNeedsInvalid {
		return t3cutil.ServiceNeedsInvalid, errors.New("t3c-check-reload returned unknown string '" + string(stdOut) + "'")
	}
	return needs, nil
}

// requestJSON calls t3c-request with the given command, and deserializes the result as JSON into obj.
func requestJSON(cfg config.Cfg, command string, obj interface{}) error {
	stdOut, err := request(cfg, command)
	if err != nil {
		return errors.New("requesting: " + err.Error())
	}
	if err := json.Unmarshal(stdOut, obj); err != nil {
		return errors.New("unmarshalling '" + string(stdOut) + "': " + err.Error())
	}
	return nil
}

// request calls t3c-request with the given command, and returns the stdout bytes.
func request(cfg config.Cfg, command string) ([]byte, error) {
	args := []string{
		"--traffic-ops-insecure=" + strconv.FormatBool(cfg.TOInsecure),
		"--traffic-ops-timeout-milliseconds=" + strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--cache-host-name=" + cfg.CacheHostName,
		`--get-data=` + command,
	}

	if cfg.LogLocationErr == log.LogLocationNull {
		args = append(args, "-s")
	}
	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
	}

	if _, used := os.LookupEnv("TO_USER"); !used {
		args = append(args, "--traffic-ops-user="+cfg.TOUser)
	}
	if _, used := os.LookupEnv("TO_PASS"); !used {
		args = append(args, "--traffic-ops-password="+cfg.TOPass)
	}
	if _, used := os.LookupEnv("TO_URL"); !used {
		args = append(args, "--traffic-ops-url="+cfg.TOURL)
	}
	stdOut, stdErr, code := t3cutil.Do(`t3c-request`, args...)
	if code != 0 {
		logSubAppErr(`t3c-request stdout`, stdOut)
		logSubAppErr(`t3c-request stderr`, stdErr)
		return nil, fmt.Errorf("t3c-request returned non-zero exit code %v, see log for output", code)
	}

	logSubApp(`t3c-request`, stdErr)

	return stdOut, nil
}

// requestConfig calls t3c-request and returns the stdout bytes.
// It also caches the config in /var/lib/trafficcontrol-cache-config and uses the cache to issue IMS requests.
func requestConfig(cfg config.Cfg) ([]byte, error) {
	// TODO support /opt

	cacheBts := ([]byte)(nil)
	if !cfg.NoCache {
		err := error(nil)
		if cacheBts, err = ioutil.ReadFile(t3cutil.ApplyCachePath); err != nil {
			// don't log an error if the cache didn't exist
			if !os.IsNotExist(err) {
				log.Errorln("getting cached config data failed, not using cache! Error: " + err.Error())
			}
			cacheBts = []byte{}
		}
	}

	log.Infof("config cache bytes: %v\n", len(cacheBts))

	args := []string{
		"request",
		"--traffic-ops-insecure=" + strconv.FormatBool(cfg.TOInsecure),
		"--traffic-ops-timeout-milliseconds=" + strconv.FormatInt(int64(cfg.TOTimeoutMS), 10),
		"--cache-host-name=" + cfg.CacheHostName,
		`--get-data=config`,
	}
	if len(cacheBts) > 0 {
		args = append(args, `--old-config=stdin`)
	}

	if cfg.LogLocationWarn != log.LogLocationNull {
		args = append(args, "-v")
	}
	if cfg.LogLocationInfo != log.LogLocationNull {
		args = append(args, "-v")
	}

	if _, used := os.LookupEnv("TO_USER"); !used {
		args = append(args, "--traffic-ops-user="+cfg.TOUser)
	}
	if _, used := os.LookupEnv("TO_PASS"); !used {
		args = append(args, "--traffic-ops-password="+cfg.TOPass)
	}
	if _, used := os.LookupEnv("TO_URL"); !used {
		args = append(args, "--traffic-ops-url="+cfg.TOURL)
	}

	stdOut := ([]byte)(nil)
	stdErr := ([]byte)(nil)
	code := 0
	if len(cacheBts) > 0 {
		stdOut, stdErr, code = t3cutil.DoInput(cacheBts, `t3c`, args...)
	} else {
		stdOut, stdErr, code = t3cutil.Do(`t3c`, args...)
	}
	if code != 0 {
		logSubAppErr(`t3c-request stdout`, stdOut)
		logSubAppErr(`t3c-request stderr`, stdErr)
		return nil, fmt.Errorf("t3c-request returned non-zero exit code %v, see log for details", code)
	}
	logSubApp(`t3c-request`, stdErr)

	if err := ioutil.WriteFile(t3cutil.ApplyCachePath, stdOut, 0600); err != nil {
		log.Errorln("writing config data to cache failed: " + err.Error())
	}

	return stdOut, nil
}

func logSubApp(appName string, stdErr []byte)    { logSubAppWarnOrErr(appName, stdErr, false) }
func logSubAppErr(appName string, stdErr []byte) { logSubAppWarnOrErr(appName, stdErr, true) }
func logSubAppWarnOrErr(appName string, stdErr []byte, isErr bool) {
	stdErr = bytes.TrimSpace(stdErr)
	if len(stdErr) == 0 {
		return
	}

	logF := log.Warnln
	logger := log.Warning
	if isErr {
		logF = log.Errorln
		logger = log.Error
	}

	if logger == nil {
		return
	}

	logF(appName + ` start`)
	// write directly to the underlying logger, to avoid duplicating the prefix
	for _, msg := range bytes.Split(stdErr, []byte{'\n'}) {
		logger.Writer().Write([]byte(string(msg) + "\n"))
	}
	logF(appName + ` end`)
}

// outToErr returns stderr if logLocation is stdout, otherwise returns logLocation unchanged.
// This is a helper to avoid logging to stdout for commands whose output is on stdout.
func outToErr(logLocation string) string {
	if logLocation == "stdout" {
		return "stderr"
	}
	return logLocation
}
