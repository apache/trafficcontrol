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

t3c - Traffic Control Cache Configuration tools

# SYNOPSIS

t3c \<command\> [\<args\>]

[\-\-help]

[\-\-version]

# DESCRIPTION

The `t3c` app generates and applies cache configuration for Apache Traffic Control.

This includes requesting Traffic Ops, generating configuration files for caching proxies such as Apache Traffic Server, verifying Traffic Ops data is valid and produces valid config files, creating a git repo for backups and history, determining whether config changes require a reload or restart of the caching proxy service, performing that restart, and more.

The latest version and documentation can be found at https://github.com/apache/trafficcontrol/cache-config.

# OPTIONS

-h, -\-help

    Prints the synopsis and usage information.

-V, -\-version

    Print the version and exit

# COMMANDS

We divide t3c into commands for each independent operation. Each command is its own application and can be called directly or via the t3c app. For example, 't3c apply' or 't3c-apply'.

t3c-apply

    Generate and apply cache configuration.

t3c-check

    Check that new config can be applied.

t3c-diff

    Diff config files, like diff or git-diff but with config-specific logic.

t3c-generate

    Generate configuration files from Traffic Ops data.

t3c-preprocess

    Preprocess generated config files.

t3c-request

    Request data from Traffic Ops.

t3c-update

    Update a server's queue and reval status in Traffic Ops.

# NOMENCLATURE

The "t3c" stands for "Traffic Control Cache Config."

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
