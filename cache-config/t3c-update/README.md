<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# t3c-update

## Synopsis
	t3c-update [-h] [-a value] [-d value] [-e value] [-H value] [-i value] \
		[-l value] [-P value] [-q value] [-t value] [-u value] [-U value]

## Description
  The t3c-update app is used to set the update and reval status,
  on Traffic Ops.

## Options
  -a  --set-reval-status [true | false] set the servers revalidate status (required)
  -q  --set-update-status [true | false] set the servers update queue status (required)

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
     		Host name of the cache to update the statuses of on TrafficOps.
        Server host name in Traffic Ops, not a URL, and not the FQDN.
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

