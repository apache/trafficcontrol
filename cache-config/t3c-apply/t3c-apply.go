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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/torequest"
	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/util"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	tcutil "github.com/apache/trafficcontrol/v8/lib/go-util"
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

const CacheConfigFailureExitMsg = `CACHE CONFIG FAILURE`
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
		log.Infoln(err)
		log.Errorln(FailureExitMsg)
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
		log.Errorf("Unable to acquire app lock retrying in: %v ", LockFileRetryInterval)
		time.Sleep(LockFileRetryInterval)
	}
	log.Infoln("Acquired app lock")
	defer lock.Unlock()

	// Note failing to load old metadata is not fatal!
	// oldMetaData must always be checked for nil before usage!
	oldMetaData, err := LoadMetaData(cfg)
	if err != nil {
		log.Errorln("Failed to load old metadata file, metadata-dependent functionality disabled: " + err.Error())
	}

	// Note we only write the metadata file after acquiring the app lock.
	// We don't want to write a metadata file if we didn't run because another t3c-apply
	// was already running.
	metaData := t3cutil.NewApplyMetaData()

	metaData.ServerHostName = cfg.CacheHostName

	t3cutil.WriteActionLog(t3cutil.ActionLogActionApplyStart, t3cutil.ActionLogStatusSuccess, metaData)

	if cfg.UseGit == config.UseGitYes {
		triedMakingRepo, err := util.EnsureConfigDirIsGitRepo(cfg)
		if err != nil {
			log.Errorln("Ensuring config directory '" + cfg.TsConfigDir + "' is a git repo - config may not be a git repo! " + err.Error())
			if triedMakingRepo {
				t3cutil.WriteActionLog(t3cutil.ActionLogActionGitInit, t3cutil.ActionLogStatusFailure, metaData)
			}
		} else {
			log.Infoln("Successfully ensured ATS config directory '" + cfg.TsConfigDir + "' is a git repo")
			if triedMakingRepo {
				t3cutil.WriteActionLog(t3cutil.ActionLogActionGitInit, t3cutil.ActionLogStatusSuccess, metaData)
			}
		}

	} else {
		log.Infoln("UseGit not 'yes', not creating git repo")
	}

	if cfg.UseGit == config.UseGitYes || cfg.UseGit == config.UseGitAuto {
		//need to see if there is an old lock file laying around.
		//older than 5 minutes
		const gitMaxLockAgeMinutes = 5
		const gitLock = ".git/index.lock"
		gitLockFile := filepath.Join(cfg.TsConfigDir, gitLock)
		oldLock, err := util.IsGitLockFileOld(gitLockFile, time.Now(), gitMaxLockAgeMinutes*time.Minute)
		if err != nil {
			log.Errorln("checking for git lock file: " + err.Error())
		}
		if oldLock {
			log.Errorf("removing git lock file older than %dm", gitMaxLockAgeMinutes)
			err := util.RemoveGitLock(gitLockFile)
			if err != nil {
				log.Errorf("couldn't remove git lock file: %v", err.Error())
			}
		}
		log.Infoln("Checking git for safe directory config")
		if err := util.GetGitConfigSafeDir(cfg); err != nil {
			log.Warnln("error checking git for safe directory config: " + err.Error())
		}
		// commit anything someone else changed when we weren't looking,
		// with a keyword indicating it wasn't our change
		if err := util.MakeGitCommitAll(cfg, util.GitChangeNotSelf, true); err != nil {
			log.Errorln("git committing existing changes, dir '" + cfg.TsConfigDir + "': " + err.Error())
			t3cutil.WriteActionLog(t3cutil.ActionLogActionGitCommitInitial, t3cutil.ActionLogStatusFailure, metaData)
		} else {
			t3cutil.WriteActionLog(t3cutil.ActionLogActionGitCommitInitial, t3cutil.ActionLogStatusSuccess, metaData)
		}
	}

	trops := torequest.NewTrafficOpsReq(cfg)

	// if doing os checks, insure there is a 'systemctl' or 'service' and 'chkconfig' commands.
	if !cfg.SkipOSCheck && cfg.SvcManagement == config.Unknown {
		log.Errorln("OS checks are enabled and unable to find any know service management tools.")
	}

	// create and clean the config.TmpBase (/tmp/ort)
	if !util.MkDir(config.TmpBase, cfg.ReportOnly) {
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

	// if running in Revalidate mode, check to see if it's
	// necessary to continue
	if cfg.Files == t3cutil.ApplyFilesFlagReval {
		syncdsUpdate, err = trops.CheckRevalidateState(false)

		if err != nil {
			log.Errorln("Checking revalidate state: " + err.Error())
			return GitCommitAndExit(ExitCodeRevalidationError, FailureExitMsg, cfg, metaData, oldMetaData)
		}
		if syncdsUpdate == torequest.UpdateTropsNotNeeded {
			log.Infoln("Checking revalidate state: returned UpdateTropsNotNeeded")
			metaData.Succeeded = true
			return GitCommitAndExit(ExitCodeRevalidationError, SuccessExitMsg, cfg, metaData, oldMetaData)
		}

	} else {
		syncdsUpdate, err = trops.CheckSyncDSState(metaData, cfg)
		if err != nil {
			log.Errorln("Checking syncds state: " + err.Error())
			return GitCommitAndExit(ExitCodeSyncDSError, FailureExitMsg, cfg, metaData, oldMetaData)
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
				if err := trops.StartServices(&syncdsUpdate, metaData, cfg); err != nil {
					log.Errorln("failed to start services: " + err.Error())
					metaData.PartialSuccess = true
					return GitCommitAndExit(ExitCodeServicesError, PostConfigFailureExitMsg, cfg, metaData, oldMetaData)
				}
			}
			finalMsg := SuccessExitMsg
			if postConfigFail {
				finalMsg = PostConfigFailureExitMsg
			}
			metaData.Succeeded = true
			return GitCommitAndExit(ExitCodeSuccess, finalMsg, cfg, metaData, oldMetaData)
		}
	}

	if cfg.Files != t3cutil.ApplyFilesFlagAll {
		// make sure we got the data necessary to check packages
		log.Infoln("======== Didn't get all files, no package processing needed or possible ========")
		metaData.InstalledPackages = oldMetaData.InstalledPackages
	} else if cfg.RpmDBOk {
		log.Infoln("======== Start processing packages  ========")
		err = trops.ProcessPackages()
		if err != nil {
			log.Errorf("Error processing packages: %s\n", err)
			return GitCommitAndExit(ExitCodePackagingError, FailureExitMsg, cfg, metaData, oldMetaData)
		}
		metaData.InstalledPackages = t3cutil.PackagesToMetaData(trops.Pkgs)

		// check and make sure packages are enabled for startup
		err = trops.CheckSystemServices()
		if err != nil {
			log.Errorf("Error verifying system services: %s\n", err.Error())
			return GitCommitAndExit(ExitCodeServicesError, FailureExitMsg, cfg, metaData, oldMetaData)
		}
	} else {
		log.Warnln("======== RPM DB checks failed, package processing not possible, using installed packages from  metadata if available========")
		trops.ProcessPackagesWithMetaData(oldMetaData.InstalledPackages)
	}

	log.Debugf("Preparing to fetch the config files for %s, files: %s, syncdsUpdate: %s\n", cfg.CacheHostName, cfg.Files, syncdsUpdate)
	err = trops.GetConfigFileList()
	if err != nil {
		log.Errorf("Getting config file list: %s\n", err)
		return GitCommitAndExit(ExitCodeConfigFilesError, FailureExitMsg, cfg, metaData, oldMetaData)
	}
	syncdsUpdate, err = trops.ProcessConfigFiles(metaData)
	if err != nil {
		log.Errorf("Error while processing config files: %s\n", err.Error())
		t3cutil.WriteActionLog(t3cutil.ActionLogActionUpdateFilesAll, t3cutil.ActionLogStatusFailure, metaData)

	} else {
		t3cutil.WriteActionLog(t3cutil.ActionLogActionUpdateFilesAll, t3cutil.ActionLogStatusSuccess, metaData)
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

	if err := trops.StartServices(&syncdsUpdate, metaData, cfg); err != nil {
		log.Errorln("failed to start services: " + err.Error())
		metaData.PartialSuccess = true
		return GitCommitAndExit(ExitCodeServicesError, PostConfigFailureExitMsg, cfg, metaData, oldMetaData)
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

	if trops.HitchReload {
		svcStatus, _, err := util.GetServiceStatus("hitch")
		cmd := "start"
		running := false
		if err != nil {
			log.Errorf("not starting 'hitch', error getting 'hitch' run status: %s\n", err)
		} else if svcStatus != util.SvcNotRunning {
			cmd = "reload"
		}
		running, err = util.ServiceStart("hitch", cmd)
		if err != nil {
			log.Errorf("'hitch' was not %sed: %s\n", cmd, err)
		} else if running {
			log.Infof("service 'hitch' %sed", cmd)
		} else {
			log.Infoln("service 'hitch' already running")
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

	metaData.Succeeded = true
	if syncdsUpdate == torequest.UpdateTropsFailed {
		return GitCommitAndExit(ExitCodeSuccess, CacheConfigFailureExitMsg, cfg, metaData, oldMetaData)
	} else {
		return GitCommitAndExit(ExitCodeSuccess, SuccessExitMsg, cfg, metaData, oldMetaData)
	}
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
func GitCommitAndExit(exitCode int, exitMsg string, cfg config.Cfg, metaData *t3cutil.ApplyMetaData, oldMetaData *t3cutil.ApplyMetaData) int {

	// metadata isn't actually part of git, but we always want to write it before committing to git, so this is the right place

	// files previously dropped never become "unmanaged",
	// and if we delete them they're removed from oldMetaData as well as the new,
	// so add the old files to the new metadata.
	// This is especially important for reval runs, which don't add all files.
	metaData.OwnedFilePaths = t3cutil.CombineOwnedFilePaths(metaData, oldMetaData)
	if len(metaData.InstalledPackages) == 0 && oldMetaData != nil {
		metaData.InstalledPackages = oldMetaData.InstalledPackages
	}
	WriteMetaData(cfg, metaData)
	success := exitCode == ExitCodeSuccess
	if cfg.UseGit == config.UseGitYes || cfg.UseGit == config.UseGitAuto {
		if err := util.MakeGitCommitAll(cfg, util.GitChangeIsSelf, success); err != nil {
			log.Errorln("git committing existing changes, dir '" + cfg.TsConfigDir + "': " + err.Error())
			// nil metadata to prevent modifying the file after the final git commit
			t3cutil.WriteActionLog(t3cutil.ActionLogActionGitCommitFinal, t3cutil.ActionLogStatusFailure, nil)
		} else {
			// nil metadata to prevent modifying the file after the final git commit
			t3cutil.WriteActionLog(t3cutil.ActionLogActionGitCommitFinal, t3cutil.ActionLogStatusSuccess, nil)
		}
	}

	if metaData.Succeeded {
		// nil metadata to prevent modifying the file after the final git commit
		t3cutil.WriteActionLog(t3cutil.ActionLogActionApplyEnd, t3cutil.ActionLogStatusSuccess, nil)
	} else {
		// nil metadata to prevent modifying the file after the final git commit
		t3cutil.WriteActionLog(t3cutil.ActionLogActionApplyEnd, t3cutil.ActionLogStatusFailure, nil)
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
	retresult := false
	if cfg.MaxMindLocation != "" {
		// Check if the maxmind db needs to be updated before reload
		MaxMindList := strings.Split(cfg.MaxMindLocation, ",")
		for _, v := range MaxMindList {
			result := util.UpdateMaxmind(v, cfg.TsConfigDir, cfg.ReportOnly)
			if result {
				log.Infoln("maxmind database was updated from " + v)
			} else {
				log.Infoln("maxmind database not updated. Either not needed or curl/gunzip failure: " + v)
			}
			if result {
				// If we've seen any database updates then return true to update ATS
				retresult = true
			}
		}
	} else {
		log.Infoln(("maxmindlocation is empty, not checking for DB update"))
	}

	return retresult
}

const MetaDataFileName = `t3c-apply-metadata.json`
const MetaDataFileMode = 0600

// WriteMetaData writes the metaData file.
//
// The metadata file is written in the ATS config directory, so it's versioned
// with git.
//
// On error, an error is written to the log, but no error is returned.
func WriteMetaData(cfg config.Cfg, metaData *t3cutil.ApplyMetaData) {
	bts, err := metaData.Format()
	if err != nil {
		log.Errorln("formatting metadata file: " + err.Error())
		return
	}

	metaDataFilePath := GetMetaDataFilePath(cfg)

	if err := os.WriteFile(metaDataFilePath, bts, MetaDataFileMode); err != nil {
		log.Errorln("writing metadata file '" + metaDataFilePath + "': " + err.Error())
		return
	}
}

func LoadMetaData(cfg config.Cfg) (*t3cutil.ApplyMetaData, error) {
	metaDataFilePath := GetMetaDataFilePath(cfg)

	bts, err := os.ReadFile(metaDataFilePath)
	if err != nil {
		return nil, errors.New("reading metadata file '" + metaDataFilePath + "': " + err.Error())
	}

	metaData := &t3cutil.ApplyMetaData{}

	if err := json.Unmarshal(bts, &metaData); err != nil {
		return nil, errors.New("unmarshalling metadata file: " + err.Error())
	}

	return metaData, nil
}

func GetMetaDataFilePath(cfg config.Cfg) string {
	return filepath.Join(cfg.TsConfigDir, MetaDataFileName)
}
