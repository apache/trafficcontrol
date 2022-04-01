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
 
[\-\-help]

[\-\-version]

# DESCRIPTION

  The t3c-update app is used to set the update and reval status in Traffic Ops.

  This is typically used after applying configuration, to set the server's "queue" or "reval" status in Traffic Ops to false.

# OPTIONS

-q, -\-set-config-apply-time

    [RFC3339Nano Timestamp] sets the server's config apply time.
    Either this or set-reval-apply-time must be used (Required)

-a, -\-set-reval-apply-time

    [RFC3339Nano Timestamp] sets the server's reval apply time.
    Either this or set-config-apply-time must be used (Required)

-H, -\-cache-host-name=value

    Host name of the cache to generate config for. Must be the
    server host name in Traffic Ops, not a URL, and not the FQDN

-h, -\-help

    Print usage information and exit

-I, -\-traffic-ops-insecure

    [true | false] ignore certificate errors from Traffic Ops

-l, -\-login-dispersion=value

    [seconds] wait a random number of seconds between 0 and
    [seconds] before login to traffic ops, default 0

-P, -\-traffic-ops-password=value

    Traffic Ops password. Required. May also be set with the
    environment variable TO_PASS

-s, -\-silent

    Silent. Errors are not logged, and the 'verbose' flag is
    ignored. If a fatal error occurs, the return code will be
    non-zero but no text will be output to stderr

-t, -\-traffic-ops-timeout-milliseconds=value

    Timeout in milli-seconds for Traffic Ops requests, default
    is 30000 [30000]

-u, -\-traffic-ops-url=value

    Traffic Ops URL. Must be the full URL, including the scheme.
    Required. May also be set with     the environment variable
    TO_URL

-U, -\-traffic-ops-user=value

    Traffic Ops username. Required. May also be set with the
    environment variable TO_USER

-v, -\-verbose

    Log verbosity. Logging is output to stderr. By default,
    errors are logged. To log warnings, pass '-v'. To log info,
    pass '-vv'. To omit error logging, see '-s'.

-V, -\-version

    Print the version and exit

-y, -\-set-update-status

    [true or nonexistent] Set the Update Status to false for the server.
    Compatability requirement until ATC (v7.0+) is deployed 
    with the timestamp features.

-z, -\-set-reval-status

    [true or nonexistent] Set the Reval Status to false for the server.
    Compatability requirement until ATC (v7.0+) is deployed 
    with the timestamp features.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
