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

t3c-preprocess - Traffic Control Cache Configuration preprocessor

# SYNOPSIS

t3c-preprocess

[\-\-help]

[\-\-version]

# DESCRIPTION

The 't3c-preprocess' app preprocesses generated config files, replacing directives with relevant data.

The stdin must be the JSON '{"data": \<data\>, "files": \<files\>}' where \<data\> is the output of 't3c-request --get-data=config' and \<files\> is the output of 't3c-generate'.

# DIRECTIVES

The following directives will be replaced. These directives may be placed anywhere in any file, by either t3c-generate(1) or Traffic Ops Parameters.

    __SERVER_TCP_PORT__ is replaced with the Server's Port from Traffic Ops; unless the server's
                        port is 80, 0, or null, in which case any occurrences preceded by a colon
                        are removed.

    __CACHE_IPV4__      is replaced with the Server's IP address from Traffic Ops.

    __HOSTNAME__        is replaced with the Server's (short) HostName from Traffic Ops.

    __FULL_HOSTNAME__   is replaced with the Server's HostName, a dot, and the Server's DomainName
                        from Traffic Ops (i.e. the Server's Fully Qualified Domain Name).

    __CACHEGROUP__      is replaced with the Server's Cachegroup name from Traffic Ops.

    __RETURN__          is replaced with a newline character, and any whitespace before or after
                        it is removed.

# OPTIONS

-h, -\-help

    Print usage information and exit

-V, -\-version

    Print version information and exit.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
