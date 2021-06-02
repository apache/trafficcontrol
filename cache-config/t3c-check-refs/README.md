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

t3c-check-refs - Traffic Control Cache Configuration generated file reference check tool

## SYNOPSIS

t3c-check-refs [-c directory] [-d location] [-e location] [-f files] [-i location] [-p directory] [file]

[\-\-help]

## DESCRIPTION

The t3c-check-refs app will read an ATS formatted plugin.config or remap.config
file line by line and verify that the plugin '.so' files are available in the
filesystem or relative to the ATS plugin installation directory by the
absolute or relative plugin filename.

In addition, any plugin parameters that end in '.config', '.cfg', or '.txt'
are considered to be plugin configuration files and there existence in the
filesystem or relative to the ATS configuration files directory is verified.

The configuration file argument is optional.  If no config file argument is
supplied, t3c-check-refs reads its config file input from stdin.

## OPTIONS

-c, --trafficserver-config-dir=value

    directory where ATS config files are stored.
    [/opt/trafficserver/etc/trafficserver]

-d, --log-location-debug=value

     Where to log debugs. May be a file path, stdout, stderr

-e, --log-location-error=value

     Where to log errors. May be a file path, stdout, stderr
     [stderr]

-f, --files-adding=value

    comma-delimited list of file names being added, to not fail
    to verify if they don't already exist.

-h, --help

    Print usage information and exit

-i, --log-location-info=value

     Where to log infos. May be a file path, stdout, stderr
     [stderr]

-p, --trafficserver-plugin-dir=value

    directory where ATS plugins are stored.
    [/opt/trafficserver/libexec/trafficserver]

## EXIT CODES

Returns 0 if no missing plugin DSO or config files are found.
Otherwise the total number of missing plugin DSO and config files
are returned.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
