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

The Compare Tool
================
The `compare` tool is used to compare the output of a set of [Traffic Ops API](https://traffic-control-cdn.readthedocs.io/en/latest/api/) endpoints between two running instances of Traffic Ops. The idea is that two different versions of Traffic Ops with the same data will have differences in the output of their API endpoints *if and only if* either the change was intentional, or a new bug was introduced in the newer version. Typically, this isn't really true, due to rapidly changing data structures like timestamps in the API outputs, but this should offer a good starting point for identifying bugs in changes made to the Traffic Ops API.

Dependencies
------------
-   github.com/apache/trafficcontrol/lib/go-tc[1]
-   github.com/kelseyhightower/envconfig
-   golang.org/x/net/publicsuffix

Usage
-----

### `compare`

Usage: compare \[-hsV\] \[-f value\] \[--ref\_passwd value\] \[--ref\_url value\] \[--ref\_user value\] \[-r value\] \[--test\_passwd value\] \[--test\_url value\] \[--test\_user value\] \[parameters ...\]

--ref\_passwd=value        The password for logging into the reference Traffic Ops instance (overrides TO\_PASSWORD environment variable)
--ref\_url=value           The URL for the reference Traffic Ops instance (overrides TO\_URL environment variable)
--ref\_user=value          The username for logging into the reference Traffic Ops instance (overrides TO\_USER environment variable)
--test\_passwd=value       The password for logging into the testing Traffic Ops instance (overrides TEST\_PASSWORD environment variable)
--test\_url=value          The URL for the testing Traffic Ops instance (overrides TEST\_URL environment variable)
--test\_user=value         The username for logging into the testing Traffic Ops instance (overrides TEST\_USER environment variable)
-f, --file=value           File listing routes to test (will read from stdin if not given)
-h, --help                 Print usage information and exit
-r, --results\_path=value  Directory where results will be written
-s, --snapshot             Perform comparison of all CDN's snapshotted CRConfigs
-V, --version              Print version information and exit

The typical way to use `compare` is to first specify some environment variables:

TO\_URL
The URL of the reference Traffic Ops instance

TO\_USER
The username to authenticate with the reference Traffic Ops instance

TO\_PASSWORD
The password to authenticate with the reference Traffic Ops instance

TEST\_URL
The URL of the testing Traffic Ops instance

TEST\_USER
The username to authenticate with the testing Traffic Ops instance

TEST\_PASSWORD
The password to authenticate with the testing Traffic Ops instance

These can be overridden by command line switches as described above. If a username and/or password is not given for the testing instance (either via environment variables or on the command line), it/they will be assumed to be the same as the one/those specified for the reference instance.

### genConfigRoutes.py

usage: genConfigRoutes.py \[-h\] \[-k\] \[-v\] InstanceA InstanceB LoginA \[LoginB\]

-h, --help                           show this help message and exit
-k, --insecure                       Do not verify SSL certificate signatures against *either* Traffic Ops instance (default: False)
-l LOG_LEVEL, --log_level LOG_LEVEL  Sets the Python log level, one of 'DEBUG', 'INFO', 'WARN', 'ERROR', or 'CRITICAL' (default: CRITICAL)
-q, --quiet                          Suppresses all logging output - even for critical errors
-v, --version                        Print version information and exit

> **note**
>
> If you're using a CDN-in-a-Box environment for testing, it's likely that you'll need the `-k`/`--insecure` option if you're outside the Docker network

The genConfigRoutes.py script will output list of unique API routes (relative to the desired Traffic Ops URL) that point to generated configuration files for a sample set of servers common to both Traffic Ops instances. The results are printed to stdout, making the output perfect for piping directly into `compare` like so:

``` sourceCode
./genConfigRoutes.py https://trafficopsA.example.test https://trafficopsB.example.test username:password | ./compare
```

... assuming the proper environment variables have been set for `compare`.

[1] Theoretically, if you downloaded the Traffic Control repository properly (into `$GOPATH/src/github.com/apache/trafficcontrol`), this will already be satisfied.


Further Information
-------------------
For more information, please see [the official documentation](https://traffic-control-cdn.readthedocs.io/en/latest/tools/compare.html)
