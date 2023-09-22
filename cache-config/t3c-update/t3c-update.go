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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-update/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
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

	// *** Compatability requirement until ATC (v7.0+) is deployed with the timestamp features
	// Use SetUpdateStatus is preferred
	err = t3cutil.SetUpdateStatus(cfg.TCCfg, tc.CacheName(cfg.TCCfg.CacheHostName), cfg.ConfigApplyTime, cfg.RevalApplyTime)
	//err = t3cutil.SetUpdateStatusCompat(cfg.TCCfg, tc.CacheName(cfg.TCCfg.CacheHostName), cfg.ConfigApplyTime, cfg.RevalApplyTime, cfg.ConfigApplyBool, cfg.RevalApplyBool)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(3)
	}

	cur_status, err := t3cutil.GetServerUpdateStatus(cfg.TCCfg)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(4)
	}

	// When comparing equality, it must be done with microsecond precision (Round not Truncate).
	// This is because Postgres stores Microsecond precision. Round also drops the monotonic
	// clock reading.
	// t3c (Nano) -> client (Nano) -> TO (Nano) -> Postgres (Micro)
	// Postgres (Micro) -> TO (Micro) -> client (Micro) -> here / t3c (Micro)
	if cfg.ConfigApplyTime != nil && !(*cfg.ConfigApplyTime).Round(time.Microsecond).Equal((*cur_status.ConfigApplyTime).Round(time.Microsecond)) {
		log.Errorf("Failed to set config_apply_time.\nSent: %v\nRecv: %v", *cfg.ConfigApplyTime, *cur_status.ConfigApplyTime)
	}
	if cfg.RevalApplyTime != nil && !(*cfg.RevalApplyTime).Round(time.Microsecond).Equal((*cur_status.RevalidateApplyTime).Round(time.Microsecond)) {
		log.Errorf("Failed to set reval_apply_time.\nSent: %v\nRecv: %v", *cfg.RevalApplyTime, *cur_status.RevalidateApplyTime)
	}
	cfg.TCCfg.TOClient.WriteFsCookie(torequtil.CookieCachePath(cfg.TOUser))
}
