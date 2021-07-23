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

	"github.com/apache/trafficcontrol/cache-config/t3c-request/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-log"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	log.Infoln("configuration initialized")

	// login to traffic ops.
	tccfg, err := t3cutil.TOConnect(&cfg.TCCfg)
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(2)
	}

	if cfg.GetData != "" {
		if err := t3cutil.WriteData(*tccfg); err != nil {
			log.Errorf("writing data: %s\n", err.Error())
			os.Exit(3)
		}
	}
}
