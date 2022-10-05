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

[\-\-version]

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

-c, -\-trafficserver-config-dir=value

    directory where ATS config files are stored.
    [/opt/trafficserver/etc/trafficserver]

-f, -\-files-adding=value

    comma-delimited list of file names being added, to not fail
    to verify if they don't already exist.

    Alternatively, this may be "input" in which case the input
    (stdin or the passed argument filename) should be JSON of
    the form:

    {"file": "config-file-text", "adding": ["files-adding"]}

    Where 'config-file-text' is the text of the config file to check
    (which otherwise would have been passed to input unadorned),
    and 'files-adding' is a JSON array of the files added
    (which otherwise would have been passed to --files-adding).
    and the full input is properly formed JSON, with the config
    file escaped for JSON.

    Callers are encouraged to use the --files-adding=input format
    to avoid errors from operating system argument length limits.

-h, -\-help

    Print usage information and exit

-p, -\-trafficserver-plugin-dir=value

    directory where ATS plugins are stored.
    [/opt/trafficserver/libexec/trafficserver]

-s, -\-silent

    Silent. Errors are not logged, and the 'verbose' flag is
    ignored. If a fatal error occurs, the return code will be
    non-zero but no text will be output to stderr

-v, -\-verbose

    Log verbosity. Logging is output to stderr. By default,
    errors are logged. To log warnings, pass '-v'. To log info,
    pass '-vv'. To omit error logging, see '-s'.

-V, -\-version

    Print version information and exit.

## EXIT CODES

Returns 0 if no missing plugin DSO or config files are found.
Otherwise the total number of missing plugin DSO and config files
are returned.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
