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

<!--

  !!!
      This file is both a Github Readme and manpage!
      Please make sure changes appear properly with man,
      and follow man conventions, such as:
      https://www.bell-labs.com/usr/dmr/www/manintro.html

      A primary goal of t3c is to follow POSIX and LSB standards
      and conventions, so it's easy to learn and use by people
      who know Linux and other *nix systems. Providing a proper
      manpage is a big part of that.
  !!!

-->
# NAME

t3c-update - Traffic Control Cache Configuration cache status updater

# SYNOPSIS

t3c-update [-ahIqv] [-d value] [-e value] [-H value] [-i value] [-l value] [-P value] [-t value] [-u value] [-U
 value]
 
[&#45;h|&#45;&#45;help]

[&#45;v|&#45;&#45;version]

# DESCRIPTION

  The t3c-update app is used to set the update and reval status in Traffic Ops.

  This is typically used after applying configuration, to set the server's "queue" or "reval" status in Traffic Ops to false.

# OPTIONS

-a, --set-reval-status

    [true | false] sets the servers revalidate status (required)

-d, --log-location-debug=value

    Where to log debugs. May be a file path, stdout, stderr

-e, --log-location-error=value

    Where to log errors. May be a file path, stdout, stderr
    [stderr]

-H, --cache-host-name=value

    Host name of the cache to generate config for. Must be the
    server host name in Traffic Ops, not a URL, and not the FQDN

-h, --help

    Print usage information and exit

-i, --log-location-info=value

    Where to log infos. May be a file path, stdout, stderr
    [stderr]

-I, --traffic-ops-insecure

    [true | false] ignore certificate errors from Traffic Ops

-l, --login-dispersion=value

    [seconds] wait a random number of seconds between 0 and
    [seconds] before login to traffic ops, default 0

-P, --traffic-ops-password=value

    Traffic Ops password. Required. May also be set with the
    environment variable TO_PASS

-q, --set-update-status

    [true | false] sets the servers update status (required)

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

-v, --version

    Print the version and exit

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
