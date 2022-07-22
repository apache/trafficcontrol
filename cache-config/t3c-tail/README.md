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

t3c-tail - Traffic Control Cache Configuration tail tool

# SYNOPSIS

t3c-tail \-f \<path to file\> \-m \<regex to match\> \-e \<regex match to exit\> \-t \<timeout in ms\>

[\-\-help]

[\-\-version]

# DESCRIPTION

The t3c-tail application will tail a file, usually a log file.
Provide a file name to watch, a regex to filter or .* is the default,
a regex match to exit tail (if omitted will exit on timeout),
timeout in milliseconds for how long you want it to run, default is 15000 milliseconds.

# OPTIONS

-e, -\-end-match

    Regex pattern that will cause tail to exit before timeout.

-f, -\-file
    Path to file to watch.

-h, -\-help

    Print usage info and exit.

-m, -\-match
    Regex pattern you want to match while running tail default is .*.

-t -\-timeout-ms
    Timeout in milliseconds that will cause tail to exit default is 15000 MS.

-V, -\-version

    Print version information and exit.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
