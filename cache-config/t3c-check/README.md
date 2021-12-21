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

t3c-check - Traffic Control Cache Configuration generated file check tool

# SYNOPSIS

t3c-check \<command\> [\<args\>]

[\-\-help]

[\-\-version]

# DESCRIPTION

The t3c-check application has commands for checking things about new config files, such as
whether they can be safely applied or if a service reload or restart will be required.

For the arguments of a command, see 't3c-check \<command\> \-\-help'.

# COMMANDS

We divide t3c-check into commands for each independent operation. Each command is its own application and can be called directly or via the t3c app. For example, 't3c check refs' or 't3c-check refs' or 't3c-check-refs'.

t3c-check-reload

    Check if a reload or restart is needed

t3c-check-refs

    Check if a config file's referenced plugins and files are valid

# OPTIONS
-h, -\-help

    Print usage information and exit

-V, -\-version

    Print version information and exit.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
