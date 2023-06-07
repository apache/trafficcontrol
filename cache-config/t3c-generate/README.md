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

[\-\-help]

[\-\-version]

# DESCRIPTION

The `t3c-generate` app generates Apache Traffic Server configuration files from Traffic Ops data.

The stdin must be JSON text as output by 't3c-request --get-data=config', which contains all the data from Traffic Ops necessary to generate configuration. For the exact format, see t3c-request(1).

The output is a JSON array of objects containing the file and its metadata.

# OPTIONS

-2, -\-default-client-enable-h2

    Whether to enable HTTP/2 on Delivery Services by default, if
    they have no explicit Parameter. This is irrelevant if ATS
    records.config is not serving H2. If omitted, H2 is
    disabled.

-a, -\-ats-version

    The ATS version, e.g. 9.1.2-42.abc123.el7.x86_64. If omitted
    generation will attempt to get the ATS version from the
    Server Profile Parameters, and fall back to
    lib/go-atscfg.DefaultATSVersion.

-b, -\-dns-local-bind

    Whether to use the server's Service Addresses to set the ATS
    DNS local bind address.

-c, -\-disable-parent-config-comments

    Disable adding a comments to parent.config individual lines.

-D, -\-dir=value

    ATS config directory, used for config files without location
    parameters or with relative paths. May be blank. If blank
    and any required config file location parameter is missing
    or relative, will error.

-G, -\-go-direct=value

    [true|false|old] if omitted go_direct is set to 'false', you
    can set go_direct 'true' which is not recommended, or 'old' which
    will be based on opposite of parent_is_proxy directive. Can also be
    overridden on a per delivery service and tier basis with a parameter.
    Default is [false]

-h, -\-help

    Print usage information and exit

-i, -\-no-outgoing-ip

    Whether to not set the records.config outgoing IP to the
    server's addresses in Traffic Ops. Default is false.

-l, -\-list-plugins

    Print the list of plugins.

-r, -\-via-string-release

    Whether to use the Release value from the RPM package as a
    replacement for the ATS version specified in the build that
    is returned in the Via and Server headers from ATS.

-s, -\-silent

    Silent. Errors are not logged, and the 'verbose' flag is
    ignored. If a fatal error occurs, the return code will be
    non-zero but no text will be output to stderr

-T, -\-default-client-tls-versions=value

    Comma-delimited list of default TLS versions for Delivery
    Services with no Parameter, e.g.
    '-\-default-tls-versions=1.1,1.2,1.3'. If omitted, all
    versions are enabled.

-v, -\-verbose

    Log verbosity. Logging is output to stderr. By default,
    errors are logged. To log warnings, pass '-v'. To log info,
    pass '-vv'. To omit error logging, see '-s'.

-V, -\-version

    Print version information and exit.

-y, -\-revalidate-only

    Whether to exclude files not named 'regex_revalidate.config'

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
