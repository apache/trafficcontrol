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

# CDN-in-a-Box Health Client Testing

## Building and Running

Build and run cdn-in-a-box `docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml up`, once up and running, using docker desktop, navigate to terminal tab of an edge or mid. cd into `/var/log/trafficcontrol` and run `tail -f tc-health-client.log`. Click on the `Open in external terminal` on upper right side and cd into `/usr/bin` and run `./tc-health-client`. Wait for the dispersion time to pass and then logs will start in the window where the tail command was ran. After that you may interact with it via Traffic Portal.

## Enable Debug instructions [Different from Production]

Debug is currently going to `/dev/null` to avoid filling up the logs. However, it can be redirect to show in the logs when debugging is needed. Systemd doesn't work well in Docker. Therefore the debbuging can be enabled in CDN-in-a-Box with the following steps:

If `tc-health-client` is already running `control + c` to stop it
Then run `/usr/bin/tc-health-client -vvv` in the `external terminal` this enables debug messages
Watch logs in the docker desktop tab where `tail -f tc-health-client.log` was ran

## Config files for Testing Only

For testing only the `tc-health-client.json` are the settings used to run it locally and can be changed. If changed `purge` all containers and run `docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml up` in the `infrastructure/cdn-in-a-box/` folder. Same applies if the `tc-health-client.service` and `to-creds` files are changed. The `tc-health-client.service` is set for `Debug` mode with `vvv` which is different from Production which is `vv`.

## Rebuilding the tc-health-client only 

Delete the `trafficcontrol-health-client-[version].rpm` from the `\dist` folder and from `/trafficcontrol/infrastructure/cdn-in-a-box/health` then cd into `/trafficcontrol` and run `./pkg -v -8 -b tc-health-client_build` this builds the RPM to be used with docker or `./pkg -v -8 tc-health-client_build` to build x86_64. Then copy the rpm from `/dist` into `/trafficcontrol/infrastructure/cdn-in-a-box/health` and rename it to `trafficcontrol-health-client.rpm` by removing the version. Build and run with `docker compose -f docker-compose.yml -f docker compose`.

## Example Testing Commands

Cd into `/opt/trafficserver/bin/` and run `./traffic_ctl host down --reason active mid-01.infra.ciab.test` or `./traffic_ctl host status mid-01.infra.ciab.test` update it as needed for other servers or reason codes. 

At the `/opt` level of a running containder for either edge or mid run `curl -vL http://trafficmonitor.infra.ciab.test:80` to test traffic_monitor