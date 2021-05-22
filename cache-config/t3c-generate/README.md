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

t3c-generate - Traffic Control Cache Configuration generation tool

# SYNOPSIS

t3c-generate [-2bchlvVy] [-D directory] [-e location] [-i location] [-T versions] [-w location]

[&#45;h|&#45;&#45;help]

[&#45;v|&#45;&#45;version]

# DESCRIPTION

The `t3c-generate` app generates Apache Traffic Server configuration files from Traffic Ops data.

The stdin must be JSON text as output by 't3c-request --get-data=config', which contains all the data from Traffic Ops necessary to generate configuration. For the exact format, see t3c-request(1).

The output is a JSON array of objects containing the file and its metadata.

# OPTIONS

-2, --default-client-enable-h2

    Whether to enable HTTP/2 on Delivery Services by default, if
    they have no explicit Parameter. This is irrelevant if ATS
    records.config is not serving H2. If omitted, H2 is
    disabled.

-b, --dns-local-bind

    Whether to use the server's Service Addresses to set the ATS
    DNS local bind address.

-c, --disable-parent-config-comments

    Disable adding a comments to parent.config individual lines.

-D, --dir=value

    ATS config directory, used for config files without location
    parameters or with relative paths. May be blank. If blank
    and any required config file location parameter is missing
    or relative, will error.

 -e, --log-location-error=value

    Where to log errors. May be a file path, stdout, stderr, or
    null. [stderr]

-h, --help

    Print usage information and exit

-i, --log-location-info=value

    Where to log information messages. May be a file path,
    stdout, stderr, or null. [stderr]

-l, --list-plugins

    Print the list of plugins.

-T, --default-client-tls-versions=value

    Comma-delimited list of default TLS versions for Delivery
    Services with no Parameter, e.g.
    '--default-tls-versions=1.1,1.2,1.3'. If omitted, all
    versions are enabled.

-v, --version

    Print version information and exit.

-V, --via-string-release

    Whether to use the Release value from the RPM package as a
    replacement for the ATS version specified in the build that
    is returned in the Via and Server headers from ATS.

-w, --log-location-warning=value

    Where to log warnings. May be a file path, stdout, stderr,
    or null. [stderr]

-y, --revalidate-only

    Whether to exclude files not named 'regex_revalidate.config'

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
