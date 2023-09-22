package ultimate_test_harness

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

	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"

	"github.com/kelseyhightower/envconfig"
)

type TrafficOpsConfig struct {
	TOURL      string `required:"true" envconfig:"TO_URL"`
	TOUser     string `required:"true" envconfig:"TO_USER"`
	TOPassword string `required:"true" envconfig:"TO_PASSWORD"`
	TOInsecure bool   `default:"true"  envconfig:"TO_INSECURE"`
	TOTimeout  int    `default:"30"    envconfig:"TO_TIMEOUT"`
}

var (
	TOConfig  TrafficOpsConfig
	TOSession *client.Session
)

func getTOConfig() {
	err := envconfig.Process("", &TOConfig)
	if err != nil {
		fmt.Printf("reading configuration from the environment: %s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	getTOConfig()
}
