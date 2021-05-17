/*
Name
	t3c-request - Traffic Control cache config Traffic Ops requestor

Synopsis
	t3c-request [-hI] [-D value] [-d value] [-e value] [-H value] [-i value] \
		[-l value] [-P value] [-t value] [-u value] [-U value]

Description
  The t3c-request app is used get update status, package information, linux
  chkconfig status, system info and status from Traffic Ops, see the
  --get-data option.  If no --get-data option is specified, the servers
  system-info is fetched and returned.

Options
	-D, --get-data=value
        non-config-file Traffic Ops Data to get. Valid values are
        update-status, packages, chkconfig, system-info, statuses,
        and config.
        Default is system-info

        Note config is not versioned between t3c versions. Callers
        should only pass config to other t3c commands of the same
        version as the t3c-request used to produce it.

	-d, --log-location-debug=value
        Where to log debugs. May be a file path, stdout or stderr.
        Default is no debug logging.
	-e, --log-location-error=value
        Where to log errors. May be a file path, stdout, or stderr.
        Default is stderr.
	-i, --log-location-info=value
        Where to log infos. May be a file path, stdout or stderr.
        Default is stderr.
	-H, --cache-host-name=value
     		Host name of the cache to generate config for. Must be the
        server host name in Traffic Ops, not a URL, and not the FQDN.
        Defaults to the OS configured hostname.
	-h, --help  Print usage information and exit
 	-I, --traffic-ops-insecure
				[true | false] ignore certificate errors from Traffic Ops
	-l, --login-dispersion=value
        [seconds] wait a random number of seconds between 0 and
        [seconds] before login to traffic ops, default 0
	-P, --traffic-ops-password=value
        Traffic Ops password. Required. May also be set with the
        environment variable TO_PASS
	-t, --traffic-ops-timeout-milliseconds=value
        Timeout in milli-seconds for Traffic Ops requests, default
        is 30000 [30000]
	-u, --traffic-ops-url=value
        Traffic Ops URL. Must be the full URL, including the scheme.
        Required. May also be set with     the environment variable
        TO_URL
	-U, --traffic-ops-user=value
        Traffic Ops username. Required. May also be set with the
        environment variable TO_USER
*/

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
