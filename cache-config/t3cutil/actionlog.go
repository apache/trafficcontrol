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
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// ActionLogAction is an action t3c performs which affects the state of the machine.
// Things that don't affect the state of the machine are not considered actions.
//
// For example, requesting Traffic Ops is not an 'action',
// but restarting ATS or modifying a config file is.
//
// Actions also include the t3c-apply run starting and finishing,
// to help delineate the actions of different runs in the log.
type ActionLogAction string

const (
	// ActionLogActionApplyStart is the start of the t3c-apply run.
	// The status of this should always be success.
	ActionLogActionApplyStart = ActionLogAction("apply-start")

	// ActionLogActionGitInit is creating a git repo in the ATS config directory.
	ActionLogActionGitInit = ActionLogAction("git-init")

	// ActionLogActionGitCommitInitial is the initial commit at the start of the run.
	ActionLogActionGitCommitInitial = ActionLogAction("create-git-commit-initial")

	// ActionLogActionGitCommitInitial is the final commit at the end of the run.
	ActionLogActionGitCommitFinal = ActionLogAction("create-git-commit-final")

	// ActionLogActionUpdateFilesAll is writing and updating ATS config files.
	ActionLogActionUpdateFilesAll = ActionLogAction("update-files-all")

	// ActionLogActionUpdateFilesReval is writing and updating only revalidate ATS config files.
	ActionLogActionUpdateFilesReval = ActionLogAction("update-files-reval")

	// ActionLogActionATSReload is calling service reload on ATS.
	ActionLogActionATSReload = ActionLogAction("ats-reload")

	// ActionLogActionATSReload is calling service restart on ATS.
	ActionLogActionATSRestart = ActionLogAction("ats-restart")

	// ActionLogActionApplyEnd is the end of the t3c-apply run.
	ActionLogActionApplyEnd = ActionLogAction("apply-end")
)

type ActionLogStatus string

const (
	ActionLogStatusSuccess = ActionLogStatus("success")
	ActionLogStatusFailure = ActionLogStatus("failure")
)

// WriteActionLog writes the given action and status to both the info log and the given metadata object.
//
// The metaData may be nil, and should be if this is being called after the final git commit,
// to prevent modifying the file after the commit.
func WriteActionLog(action ActionLogAction, status ActionLogStatus, metaData *ApplyMetaData) {
	if metaData != nil {
		metaData.Actions = append(metaData.Actions, ApplyMetaDataAction{
			Action: string(action),
			Status: string(status),
		})
	}
	log.Infoln(`ACTION='` + string(action) + `' STATUS='` + string(status) + `'` + "\n")
}
