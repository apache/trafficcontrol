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

t3c-diff - Traffic Control Cache Configuration contextual diff tool

# SYNOPSIS

t3c-diff \-a \<file-a\> \-b \<file-b\> \-l \<line_comment\> \-m \<file-mode\> \-u \<file-uid\> \-g \<file-gid\>

[\-\-help]

[\-\-version]

# DESCRIPTION

The t3c-diff application compares configuration files with semantic context, omitting comments and other semantically irrelevant text.

This is useful over standard diff tools without context, for example, when the grammar of a generated comment changes, or a comment contains a date. This allows operators to avoid updating sematically identical files, undesirably updating file timestamps, effecting unnecessary reloads, and other unnecessary and undesirable results.

The input files may be file paths, or 'stdin' in which case that file is read from stdin.

Prints the diff to stdout, and returns the exit code 0 if there was no diff, 1 if there was a diff.
If one file exists but the other doesn't, it will always be a diff.

Note this means there may be no diff text printed to stdout but still exit 1 indicating a diff
if the file being created or deleted is semantically empty.

Mode is file permissions in octal format, default is 0644.
Line comment is a character that signals the line is a comment, default is #

Uid is the User id the file being checked should have, default is running process's uid.
Gid is the Group id the file being checked should have, default is running process's gid.`

# OPTIONS

-a, -\-file-a

    Path to first diff file, can also be stdin.

-b, -\-file-b
    Path to second diff file, can also be stdin.

-g, -\-file-gid
    Group id the file being checked should have.
    
-h, -\-help

    Print usage info and exit.

-l, -\-line_comment
    Symbol used to denote the line is a comment.    

-m, -\-file-mode
    Octal permissions mode for file being checked.

-u, -\-file-uid
    User id the file being checked should have.

-V, -\-version

    Print version information and exit.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
