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

	"github.com/apache/trafficcontrol/v6/cache-config/t3c-update/config"
	"github.com/apache/trafficcontrol/v6/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v6/cache-config/t3cutil/toreq"
	"github.com/apache/trafficcontrol/v6/lib/go-log"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

var (
	cfg config.Cfg
)

func main() {
	var err error

	cfg, err = config.InitConfig()
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
		cfg.UserAgent,
	)
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(2)
	}
	if cfg.TCCfg.TOClient.FellBack() {
		log.Warnln("Traffic Ops does not support the latest version supported by this app! Falling back to previous major Traffic Ops API version!")
	}

	err = t3cutil.SetUpdateStatus(cfg.TCCfg, tc.CacheName(cfg.TCCfg.CacheHostName), cfg.UpdatePending, cfg.RevalPending)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(3)
	}

	cur_status, err := t3cutil.GetServerUpdateStatus(cfg.TCCfg)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(4)
	}

	if cur_status.UpdatePending != cfg.UpdatePending && cfg.RevalPending != cfg.RevalPending {
		log.Errorf("ERROR: update failed, update status and/or reval status was not set.\n")
	} else {
		log.Infoln("Update successfully completed")
	}

}
