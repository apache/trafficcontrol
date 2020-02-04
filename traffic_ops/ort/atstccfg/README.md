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
# atstccfg
atstccfg is a tool for generating configuration files server-side on ATC cache servers.

!!! Warning !!!
    <strong>atstccfg does not have a stable command-line interface, it can and will change without warning. Scripts should avoid calling it for the time being.</strong>

## Usage
```
atstccfg [-u TO_URL] [-U TO_USER] [-P TO_PASSWORD] [-n] [-r N] [-e ERROR_LOCATION] [-w WARNING_LOCATION] [-i INFO_LOCATION] [-g] [-s] [-t TIMEOUT] [-a MAX_AGE] [-l]
```
The available options are:
```
-a, --cache-file-max-age-seconds                                Sets the maximum age - in seconds - a cached response can be in order to be considered "fresh" - older files will be re-generated and cached. Default: 60
-e ERROR_LOCATION, --log-location-error ERROR_LOCATION          The file location to which to log errors. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
-g, --print-generated-files                                     If given, the names of files generated (and not proxied to Traffic Ops) will be printed to stdout, then atstccfg will exit.
-h, --help                                                      Print usage information and exit.
-i INFO_LOCATION, --log-location-info INFO_LOCATION             The file location to which to log information messages. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
-l, --list-plugins                                              List the loaded plugins and then exit.
-n, --no-cache                                                  If given, existing cache files will not be used. Cache files will still be created, existing ones just won't be used.
-P TO_PASSWORD                                                  Authenticate using this password - if not given, atstccfg will attempt to use the value of the TO_PASS environment variable
-r N, --num-retries N                                           The number of times to retry getting a file if it fails. Default: 5
-s, --traffic-ops-insecure                                      If given, SSL certificate errors will be ignored when communicating with Traffic Ops. NOT RECOMMENDED FOR PRODUCTION ENVIRONMENTS.
-t, --traffic-ops-timeout-milliseconds                          Sets the timeout - in milliseconds - for requests made to Traffic Ops. Default: 10000
-u TO_URL                                                       Request this URL, e.g. 'https://trafficops.infra.ciab.test/servers/edge/configfiles/ats'
-U TO_USER                                                      Authenticate as the user TO_USER - if not given, atstccfg will attempt to use the value of the TO_USER environment variable
-v, --version                                                   Print version information and exit.
-w WARNING_LOCATION, --log-location-warning WARNING_LOCATION    The file location to which to log warnings. Respects the special string constants of github.com/apache/trafficcontrol/lib/go-log. Default: 'stderr'
```
atstccfg caches generated files in /tmp/atstccfg_cache/ for re-use.
