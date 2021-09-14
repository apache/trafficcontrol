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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3c-apply/torequest"
	"github.com/apache/trafficcontrol/cache-config/t3c-apply/util"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-log"
	tcutil "github.com/apache/trafficcontrol/lib/go-util"
)

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

const (
	ExitCodeSuccess           = 0
	ExitCodeAlreadyRunning    = 132
	ExitCodeConfigFilesError  = 133
	ExitCodeConfigError       = 134
	ExitCodeGeneralFailure    = 135
	ExitCodePackagingError    = 136
	ExitCodeRevalidationError = 137
	ExitCodeServicesError     = 138
	ExitCodeSyncDSError       = 139
	ExitCodeUserCheckError    = 140
)

func runSysctl(cfg config.Cfg) {
	if cfg.ReportOnly {
		return
	}
	if cfg.ServiceAction == t3cutil.ApplyServiceActionFlagRestart {
		_, rc, err := util.ExecCommand("/usr/sbin/sysctl", "-p")
		if err != nil {
			log.Errorln("sysctl -p failed")
		} else if rc == 0 {
			log.Debugf("sysctl -p ran succesfully.")
		}
	}
}

const LockFilePath = "/var/run/t3c.lock"
const LockFileRetryInterval = time.Second
const LockFileRetryTimeout = time.Minute

const FailureExitMsg = `CRITICAL FAILURE, ABORTING`
const PostConfigFailureExitMsg = `CRITICAL FAILURE AFTER SETTING CONFIG, ABORTING`
const SuccessExitMsg = `SUCCESS`

func main() {
	os.Exit(LogPanic(Main))
}

// Main is the main function of t3c-apply.
// This is a separate function so defer statements behave as-expected.
// DO NOT call os.Exit within this function; return the code instead.
// Returns the application exit code.
func Main() int {
	var syncdsUpdate torequest.UpdateStatus
	var lock util.FileLock
	cfg, err := config.GetCfg(Version, GitRevision)
	if err != nil {
		fmt.Println(err)
		fmt.Println(FailureExitMsg)
		return ExitCodeConfigError
	} else if cfg == (config.Cfg{}) { // user used the --help option
		return ExitCodeSuccess
	}

	log.Infoln("Trying to acquire app lock")
	for lockStart := time.Now(); !lock.GetLock(LockFilePath); {
		if time.Since(lockStart) > LockFileRetryTimeout {
			log.Errorf("Failed to get app lock after %v seconds, another instance is running, exiting without running\n", int(LockFileRetryTimeout/time.Second))
			log.Infoln(FailureExitMsg)
			return ExitCodeAlreadyRunning
		}
		time.Sleep(LockFileRetryInterval)
	}
	log.Infoln("Acquired app lock")

	if cfg.UseGit == config.UseGitYes {
		err := util.EnsureConfigDirIsGitRepo(cfg)
		if err != nil {
			log.Errorln("Ensuring config directory '" + cfg.TsConfigDir + "' is a git repo - config may not be a git repo! " + err.Error())
		} else {
			log.Infoln("Successfully ensured ATS config directory '" + cfg.TsConfigDir + "' is a git repo")
		}
	} else {
		log.Infoln("UseGit not 'yes', not creating git repo")
	}

	if cfg.UseGit == config.UseGitYes || cfg.UseGit == config.UseGitAuto {
		// commit anything someone else changed when we weren't looking,
		// with a keyword indicating it wasn't our change
		if err := util.MakeGitCommitAll(cfg, util.GitChangeNotSelf, true); err != nil {
			log.Errorln("git committing existing changes, dir '" + cfg.TsConfigDir + "': " + err.Error())
		}
	}

	e2eSSLDir := filepath.Join(cfg.TsConfigDir, "ssl") // TODO make configurable
	hasClientCerts := false
	clientCertBasePath := filepath.Join(e2eSSLDir, util.E2ESSLPathClientBase)

	log.Infoln("DEBUG e2eSSLDir '" + e2eSSLDir + "'")

	e2eSSLCADestFile := "e2e-ssl-ca.cert"
	e2eSSLCADestDir := filepath.Join(cfg.TsConfigDir, "ssl")
	e2eCACertDestPath := filepath.Join(e2eSSLCADestDir, e2eSSLCADestFile)

	if _, err := os.Stat(cfg.E2ESSLCACertPath); os.IsNotExist(err) {
		log.Errorf("end-to-end ssl certificate authority certificate '%v' does not exist, not creating certificates", cfg.E2ESSLCACertPath)
	} else if err != nil {
		log.Errorf("end-to-end ssl certificate authority certificate '%v' error reading file: %v", cfg.E2ESSLCACertPath, err)
	} else if _, err := os.Stat(cfg.E2ESSLCAKeyPath); os.IsNotExist(err) {
		log.Errorf("end-to-end ssl certificate authority key '%v' does not exist, not creating certificates", cfg.E2ESSLCAKeyPath)
	} else if err != nil {
		log.Errorf("end-to-end ssl certificate authority key '%v' error reading file: %v", cfg.E2ESSLCACertPath, err)
	} else if err := util.E2ESSLGenerateKeysIfNotExist(e2eSSLDir, cfg.E2ESSLCAKeyPath, cfg.E2ESSLCACertPath, e2eCACertDestPath); err != nil {
		log.Errorf("generating end-to-end ssl client certificate: %v", err)
	} else {
		log.Errorf("successfully generated end-to-end ssl client certificate " + clientCertBasePath + ".cert")
		hasClientCerts = true
	}

	genInf := config.GenInf{}
	genInf.HasClientCerts = hasClientCerts
	genInf.CACertPath = e2eCACertDestPath

	trops := torequest.NewTrafficOpsReq(cfg)

	// if doing os checks, insure there is a 'systemctl' or 'service' and 'chkconfig' commands.
	if !cfg.SkipOSCheck && cfg.SvcManagement == config.Unknown {
		log.Errorln("OS checks are enabled and unable to find any know service management tools.")
	}

	// create and clean the config.TmpBase (/tmp/ort)
	if !util.MkDir(config.TmpBase, cfg) {
		log.Errorln("mkdir TmpBase '" + config.TmpBase + "' failed, cannot continue")
		log.Infoln(FailureExitMsg)
		return ExitCodeGeneralFailure
	} else if !util.CleanTmpDir(cfg) {
		log.Errorln("CleanTmpDir failed, cannot continue")
		log.Infoln(FailureExitMsg)
		return ExitCodeGeneralFailure
	}

	log.Infoln(time.Now().Format(time.RFC3339))

	if !util.CheckUser(cfg) {

		lock.Unlock()
		return ExitCodeUserCheckError
	}

	toolName := trops.GetHeaderComment()
	log.Debugf("toolname: %s\n", toolName)

	// if running in Revalidate mode, check to see if it's
	// necessary to continue
	if cfg.Files == t3cutil.ApplyFilesFlagReval {
		syncdsUpdate, err = trops.CheckRevalidateState(false, genInf)

		if err != nil {
			log.Errorln("Checking revalidate state: " + err.Error())
			return GitCommitAndExit(ExitCodeRevalidationError, FailureExitMsg, cfg)
		}
		if syncdsUpdate == torequest.UpdateTropsNotNeeded {
			log.Infoln("Checking revalidate state: returned UpdateTropsNotNeeded")
			return GitCommitAndExit(ExitCodeRevalidationError, SuccessExitMsg, cfg)
		}

	} else {
		syncdsUpdate, err = trops.CheckSyncDSState(genInf)
		if err != nil {
			log.Errorln("Checking syncds state: " + err.Error())
			return GitCommitAndExit(ExitCodeSyncDSError, FailureExitMsg, cfg)
		}
		if !cfg.IgnoreUpdateFlag && cfg.Files == t3cutil.ApplyFilesFlagAll && syncdsUpdate == torequest.UpdateTropsNotNeeded {
			// If touching remap.config fails, we want to still try to restart services
			// But log a critical-post-config-failure, which needs logged right before exit.
			postConfigFail := false
			// check for maxmind db updates even if we have no other updates
			if CheckMaxmindUpdate(cfg) {
				// We updated the db so we should touch and reload
				trops.RemapConfigReload = true
				path := cfg.TsConfigDir + "/remap.config"
				_, rc, err := util.ExecCommand("/usr/bin/touch", path)
				if err != nil {
					log.Errorf("failed to update the remap.config for reloading: %s\n", err.Error())
					postConfigFail = true
				} else if rc == 0 {
					log.Infoln("updated the remap.config for reloading.")
				}
				if err := trops.StartServices(&syncdsUpdate); err != nil {
					log.Errorln("failed to start services: " + err.Error())
					return GitCommitAndExit(ExitCodeServicesError, PostConfigFailureExitMsg, cfg)
				}
			}
			finalMsg := SuccessExitMsg
			if postConfigFail {
				finalMsg = PostConfigFailureExitMsg
			}
			return GitCommitAndExit(ExitCodeSuccess, finalMsg, cfg)
		}
	}

	if cfg.Files != t3cutil.ApplyFilesFlagAll {
		// make sure we got the data necessary to check packages
		log.Infoln("======== Didn't get all files, no package processing needed or possible ========")
	} else {
		log.Infoln("======== Start processing packages  ========")
		err = trops.ProcessPackages()
		if err != nil {
			log.Errorf("Error processing packages: %s\n", err)
			return GitCommitAndExit(ExitCodePackagingError, FailureExitMsg, cfg)
		}

		// check and make sure packages are enabled for startup
		err = trops.CheckSystemServices()
		if err != nil {
			log.Errorf("Error verifying system services: %s\n", err.Error())
			return GitCommitAndExit(ExitCodeServicesError, FailureExitMsg, cfg)
		}
	}

	log.Debugf("Preparing to fetch the config files for %s, files: %s, syncdsUpdate: %s\n", cfg.CacheHostName, cfg.Files, syncdsUpdate)
	err = trops.GetConfigFileList(genInf)
	if err != nil {
		log.Errorf("Getting config file list: %s\n", err)
		return GitCommitAndExit(ExitCodeConfigFilesError, FailureExitMsg, cfg)
	}
	syncdsUpdate, err = trops.ProcessConfigFiles()
	if err != nil {
		log.Errorf("Error while processing config files: %s\n", err.Error())
	}

	// check for maxmind db updates
	// If we've updated also reload remap to reload the plugin and pick up the new database
	if CheckMaxmindUpdate(cfg) {
		trops.RemapConfigReload = true
	}

	if trops.RemapConfigReload == true {
		cfg, ok := trops.GetConfigFile("remap.config")
		_, rc, err := util.ExecCommand("/usr/bin/touch", cfg.Path)
		if err != nil {
			log.Errorf("failed to update the remap.config for reloading: %s\n", err.Error())
		} else if rc == 0 && ok == true {
			log.Infoln("updated the remap.config for reloading.")
		}
	}

	if err := trops.StartServices(&syncdsUpdate); err != nil {
		log.Errorln("failed to start services: " + err.Error())
		return GitCommitAndExit(ExitCodeServicesError, PostConfigFailureExitMsg, cfg)
	}

	// start 'teakd' if installed.
	if trops.IsPackageInstalled("teakd") {
		svcStatus, pid, err := util.GetServiceStatus("teakd")
		if err != nil {
			log.Errorf("not starting 'teakd', error getting 'teakd' run status: %s\n", err)
		} else if svcStatus == util.SvcNotRunning {
			running, err := util.ServiceStart("teakd", "start")
			if err != nil {
				log.Errorf("'teakd' was not started: %s\n", err)
			} else if running {
				log.Infoln("service 'teakd' started.")
			} else if svcStatus == util.SvcRunning {
				log.Infof("service 'teakd' was already running, pid: %v\n", pid)
			}
		}
	}

	// reload sysctl
	if trops.SysCtlReload == true {
		runSysctl(cfg)
	}

	trops.PrintWarnings()

	if err := trops.UpdateTrafficOps(&syncdsUpdate); err != nil {
		log.Errorf("failed to update Traffic Ops: %s\n", err.Error())
	}

	return GitCommitAndExit(ExitCodeSuccess, SuccessExitMsg, cfg)
}

func LogPanic(f func() int) (exitCode int) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic: (err: %v) stacktrace:\n%s\n", err, tcutil.Stacktrace())
			log.Infoln(FailureExitMsg)
			exitCode = ExitCodeGeneralFailure
			return
		}
	}()
	return f()
}

// GitCommitAndExit attempts to git commit all changes, and logs any error.
// It then logs exitMsg at the Info level, and returns exitCode.
// This is a helper function, to reduce the duplicated commit-log-return into a single line.
func GitCommitAndExit(exitCode int, exitMsg string, cfg config.Cfg) int {
	success := exitCode == ExitCodeSuccess
	if cfg.UseGit == config.UseGitYes || cfg.UseGit == config.UseGitAuto {
		if err := util.MakeGitCommitAll(cfg, util.GitChangeIsSelf, success); err != nil {
			log.Errorln("git committing existing changes, dir '" + cfg.TsConfigDir + "': " + err.Error())
		}
	}
	log.Infoln(exitMsg)
	return exitCode
}

// CheckMaxmindUpdate will (if a url is set) check for a db on disk.
// If it exists, issue an IMS to determine if it needs to update the db.
// If no file or if an update is needed to be done it is downloaded and unpacked.
func CheckMaxmindUpdate(cfg config.Cfg) bool {
	// Check if we have a URL for a maxmind db
	// If we do, test if the file exists, do IMS based on disk time
	// and download and unpack as needed
	result := false
	if cfg.MaxMindLocation != "" {
		// Check if the maxmind db needs to be updated before reload
		result = util.UpdateMaxmind(cfg)
		if result {
			log.Infoln("maxmind database was updated from " + cfg.MaxMindLocation)
		} else {
			log.Infoln("maxmind database not updated. Either not needed or curl/gunzip failure")
		}
	} else {
		log.Infoln(("maxmindlocation is empty, not checking for DB update"))
	}

	return result
}
