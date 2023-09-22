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

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-request/config"
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
	}
	log.Infoln("configuration initialized")

	// login to traffic ops.
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

	if cfg.GetData != "" {
		if err := t3cutil.WriteData(cfg.TCCfg); err != nil {
			log.Errorf("writing data: %s\n", err.Error())
			os.Exit(3)
		}
	}
	cfg.TCCfg.TOClient.WriteFsCookie(torequtil.CookieCachePath(cfg.TOUser))
}
