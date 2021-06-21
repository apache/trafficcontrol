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
	"github.com/apache/trafficcontrol/lib/go-log"
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

	// login to TrafficOps
	tccfg, err := t3cutil.TOConnect(&cfg.TCCfg)
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(2)
	}

	err = t3cutil.SetUpdateStatus(*tccfg, cfg.TCCfg.CacheHostName, cfg.UpdatePending, cfg.RevalPending)
	if err != nil {
		log.Errorf("%s, %s\n", err, cfg.TCCfg.CacheHostName)
		os.Exit(3)
	}

	cur_status, err := t3cutil.GetServerUpdateStatus(*tccfg)
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
