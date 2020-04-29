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
atstccfg [-u TO_URL] [-U TO_USER] [-P TO_PASSWORD] [-n] [-r N] [-e ERROR_LOCATION] [-w WARNING_LOCATION] [-i INFO_LOCATION] [-g] [-s] [-t TIMEOUT] [-a REVAL_STATUS] [-l]
```
The available options are:
```
-a, --set-reval-status string
    Sets the reval_pending property of the server in Traffic Ops. Must be 'true'
    or 'false'. Requires --set-queue-status also be set. This disables normal
    output.
-e, --log-location-error string
    A location for error-level logging. Passing "stderr" causes it to log to
    STDERR, "stdout" causes logging to STDOUT, "null" disables error-level
    logging, and anything else is treated as a path to a file which will contain
    the logs. (Default: stderr)
-d, --get-data string
    Specifies non-configuration-file data to retrieve from Traffic Ops. This
    disables normal output. Valid values are update-status, packages, chkconfig,
    system-info, and statuses. Output is in JSON-encoded format. For specifics,
    refer to the official documentation.
-h, --help
    Print usage information and exit.
-i, --log-location-info string
    A location for informative-level logging. Passing "stderr" causes it to log
    to STDERR, "stdout" causes logging to STDOUT, "null" disables
    informative-level logging, and anything else is treated as a path to a file
    which will contain the logs. (Default: stderr)
-l, --list-plugins
    Print the list of plug-ins and exit.
-n, --cache-host-name string
    Required. Specifies the (short) hostname of the cache server for which
    output will be generated. Must be the server host name in Traffic Ops, not a
    URL, or Fully Qualified Domain Name. Behavior when more than one server
    exists with the passed hostname is undefined.
-p, --traffic-ops-disable-proxy
    Bypass the Traffic Ops caching proxy and make requests directly to Traffic
    Ops. Has no effect if no such proxy exists.
-P, --traffic-ops-password string
    The password to use when authenticating with Traffic Ops password. If not
    given, the value of the TO_PASSWORD environment variable is used. If that
    environment variable is not set, this option-argument is required.
-q, --set-queue-status string
    Sets the upd_pending property of the server in Traffic Ops. Must be 'true'
    or 'false'. Requires --set-reval-status also be set. This disables normal
    output.
-r, --num-retries int
    The number of times to retry getting a file if it fails. (Default 5)
-s, --traffic-ops-insecure
    Ignore HTTPS certificate errors from Traffic Ops. It is HIGHLY RECOMMENDED
    to never use this in a production environment, but only for debugging.
-t, --traffic-ops-timeout-milliseconds int
    Timeout in milliseconds for Traffic Ops requests. (Default 30000)
-u, --traffic-ops-url string
    The full URL, including scheme and optionally port number, of the Traffic
    Ops server. If not given, the value of the TO_URL environment variable will
    be used. If that environment variable is not properly set, this
    option-argument is required.
-U, --traffic-ops-user string
    The username to use when authenticating with Traffic Ops. If not given, the
    value of the TO_USER environment variable will be used. If that environment
    variable is not set, this option-argument is required.
-v, --version
    Print version information and exit.
-w, --log-location-warning string
    A location for warning-level logging. Passing "stderr" causes it to log to
    STDERR, "stdout" causes logging to STDOUT, "null" disables warning-level
    logging, and anything else is treated as a path to a file which will contain
    the logs. (Default: stderr)
-y, --revalidate-only
    When given, atstccfg will only emit files relevant for updating content
    invalidation jobs. For Apache Traffic Server implementations, this limits
    the output to be only files named 'regex_revalidate.config'.
```
atstccfg caches generated files in /tmp/atstccfg_cache/ for re-use.

# Development

## Updating for new Traffic Control Versions

After a new Traffic Control release, the Traffic Ops client from the new release
branch should be vendored at `toreq/vendor`, and all usages of
`config.TOClientNew` should be changed to `config.TOClient`.

There's a "script" to do this at
[`./update-to-client/update-to-client.go`](./update-to-client). Run the "script"
with no arguments for usage information.
