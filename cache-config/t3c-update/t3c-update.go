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

	"github.com/apache/trafficcontrol/cache-config/t3c-update/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/t3cutil/toreq"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

func main() {
	cfg, err := config.InitConfig(Version, GitRevision)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	} else {
		log.Infoln("configuration initialized")
	}

	cfg.TCCfg.TOClient, err = toreq.New(
		cfg.TOURL,
		cfg.TOUser,
		cfg.TOPass,
		cfg.TOInsecure,
		cfg.TOTimeoutMS,
		cfg.UserAgent(),
	)
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(2)
	}
	if cfg.TCCfg.TOClient.FellBack() {
		log.Warnln("Traffic Ops does not support the latest version supported by this app! Falling back to previous major Traffic Ops API version!")
	}

	err = t3cutil.SetUpdateStatus(cfg.TCCfg, tc.CacheName(cfg.TCCfg.CacheHostName), cfg.ConfigApplyTime, cfg.RevalApplyTime)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(3)
	}

	cur_status, err := t3cutil.GetServerUpdateStatus(cfg.TCCfg)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(4)
	}

	if cfg.ConfigApplyTime != nil && !(*cfg.ConfigApplyTime).Equal(*cur_status.ConfigApplyTime) {
		log.Errorf("Failed to set config_apply_time.\nSent: %v\nRecv: %v", *cfg.ConfigApplyTime, *cur_status.ConfigApplyTime)
	}

	if cfg.RevalApplyTime != nil && !(*cfg.RevalApplyTime).Equal(*cur_status.RevalidateApplyTime) {
		log.Errorf("Failed to set reval_apply_time.\nSent: %v\nRecv: %v", *cfg.RevalApplyTime, *cur_status.RevalidateApplyTime)
	}

	if (*cur_status.ConfigUpdateTime).After(*cur_status.ConfigApplyTime) {
		// Another update appears to have been queued. Should we run again?
		log.Warnf("Config Update and Apply times do not match")
	}

	if (*cur_status.RevalidateUpdateTime).After(*cur_status.RevalidateApplyTime) {
		// Another reval appears to have been queued. Should we run again?
		log.Warnf("Reval Update and Apply times do not match")
	}

	log.Infoln("Update completed")
}
