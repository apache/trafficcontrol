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
# ORT Rewrite in UNIX Philosophy

## Problem Description
ORT is:
- Difficult to maintain. Writing Perl is difficult, and reading it is even more difficult.
- Dangerous to modify. Perl is not compiled, and even validity checks (`perl –c`) fail to verify dynamic runtime errors. This makes it very easy to introduce a bug in seldom-executed areas.
- Untested. Perl ORT has no unit or integration tests.
- Opaque. Nobody really knows everything it does, or when, or why.

## Proposed Change

ORT will be rewritten into a series of standalone executables, in the "UNIX Philosophy"

> 1. Make each program do one thing well. To do a new job, build afresh rather than complicate old programs by adding new "features".
> 2. Expect the output of every program to become the input to another, as yet unknown, program. Don't clutter output with extraneous information. Avoid stringently columnar or binary input formats. Don't insist on interactive input.

- Each executable should do exactly 1 thing, and if a new "thing" becomes necessary, a new executable will be created.
- The input and output of executables should be text which is easily parseable, so the executables can easily be pipelined (passing the output of one to the input of another), as well as easily read by humans and manipulated by standard Linux/POSIX tools.

This makes ORT:
- Easier to maintain. Each binary does one thing, is much smaller, and is more obvious. Presumably they’re also written in a language easier to read and write, such as Go.
- Safer to modify. If each component is smaller, it’s more obvious what it does. We also presume the apps will be written with good development practices (such as modularization), with a language which verifies more at compile-time, and with tests.
- Clear and easy for operators to understand what each app does. We assume clean interfaces, and good documentation (ideally in the app itself, via help flags, printing usage when no arguments are received, and/or man pages).

#### Implementation

The implementation should adhere to the "UNIX Philosophy," POSIX, Linux Standard Base (LSB), and GNU as much as possible.

ORT will continue to consist of a single OS package (e.g. RPM), which installs all executables.

ORT will require the following executables:
- **Aggregator**. This is the “primary application” which will emulate the existing ORT script, and be called by CRON or operators to deploy all configs, as ORT does today. Note this is similar to how git works, and several other common Linux CLI utilities. Will include a "Report Mode" which executes the pipeline of commands necessary to emulate the existing ORT Report Mode.
  This app will have no logic itself, except to call the other executables.
    - INPUT: configuration and specification to fetch and emplace config files.
    - BEHAVIOR: fetches and places config files
    - OUTPUT: success or failure message

- **Traffic Ops Requestor**. This will fetch data needed from Traffic Ops, such as the Update Pending flag, packages, etc. This should never modify TO data, and should be guaranteed read-only. Any status modifications should go in the Traffic Ops Updater.
    - INPUT: Traffic Ops URL and credentials, and data to fetch
    - BEHAVIOR: Requests data from Traffic Ops
    - OUTPUT: Traffic Ops data requested
        - Format is probably multipart/mixed, but format may be different if a better format is determined. Ideal "UNIX Philosophy" format is line-delimited text, but the complexity may preclude that. The more complex and difficult to parse, the further from the "UNIX Philosophy." E.g. multipart/mixed is preferable to JSON.
- **Config File Generator**. This will take TO data and produce config files.
    - INPUT: Traffic Ops data, and config file(s) to generate or exclude (typically all, possibly “reval only” or other behaviors of ORT)
    - BEHAVIOR: No side effects. Computationally: builds requested files.
    - OUTPUT: Config files
        - Format is probably multipart/mixed, but as above, may be different, multipart is preferable to JSON, etc.
- **Config File Preprocessor**. Preprocesses generated config files, making post-generation modifications such as underscore directive replacements, and remap OVERRIDE replacements. TODO: determine if this should be rolled into the Config File Generator.
    - INPUT: Config files, TO data/metadata.
    - BEHAVIOR: No side effects; computationally: modifies files per rules and data.
    - OUTPUT: modified config files.
- **Server Config Readiness Verifier**. Verifies the operating system is ready and safe to apply the given config. Currently, this is just checking udev rules and verifying ATS block devices don’t have filessytems. But it may be more in the future. TODO: determine if necessary; should this even be ORT’s job? Will we ever need anything besides udev/50-ats.rules? 
    - INPUT: config files
    - BEHAVIOR: No side effects; reads configs and inspects server state.
    - OUTPUT: whether server is safe and ready to apply config files.
- **ATS Plugin Readiness Verifier**. Verifies ATS has the necessary plugins for the config files.
    - INPUT: config files
    - BEHAVIOR: No side effects; reads configs and inspects server state.
    - OUTPUT: whether ATS has all necessary plugins, or which files require which missing plugins.
- **Diff Tool**. This will take two config files (presumably an existing and new file) and return their differences. TODO: determine if this is necessary; if configs are deterministic, can POSIX diff be used?
    - INPUT: Two config files
    - BEHAVIOR: No side effects. Computationally: diffs given files.
    - OUTPUT: file diff. Ideally in a standard format.
- **Backup Tool**. This will take a file and copy it to a backup location. TODO: determine if necessary; are standard POSIX cp/mv/etc enough?
    - INPUT: config file to backup, and backup location. May be text of new file or path of existing file.
    - BEHAVIOR: Save or copy given file to given backup location. Should be atomic.
    - OUTPUT: success or failure message.
- **Restart Determiner**. Takes the config to be applied (only the changed files, after diffing), encapsulates the logic of what changes require a reload or restart, and returns whether a restart, reload, or neither is required.
    - INPUT: config files to be applied.
    - BEHAVIOR: no side effects; computationally: inspects files and determines action.
    - OUTPUT: whether ATS needs reloaded or restarted.
- **Service Reloader**. Takes the name of the service (possibly only ATS) to reload or restart, and reloads or restarts as necessary. TODO determine if necessary; is this any logic beyond calling service restart and/or traffic_ctl?
    - INPUT: service, whether to reload or restart
    - BEHAVIOR: reloads or restarts the service
    - OUTPUT: success or failure
- **Traffic Ops Updater**. This will set the server’s update status in Traffic Ops. This should be only the Update and Reval Pending flags; ORT should never modify server configuration data, only ever server configuration status data.
    - INPUT: Traffic Ops URL and credentials, and status to set
    - BEHAVIOR: Makes a POST request to TO setting the status
    - OUTPUT: success or failure message

#### Features Omitted

The following features of the current ORT are specifically not being implemented in the redesign:

- **Chkconfig**. Chkconfig is not used by CentOS 7+; specifically, SystemD does not use it. It is misleading that ORT sets it today.
- **Ntpd**. ORT currently has custom logic to restart ntpd if an ntpd.conf is changed. This should be managed by whatever system is managing the server (Ansible/Puppet/Manual/etc). Network time should not be the responsibility of Traffic Control or its config applicator. 
- **Interactive mode**. This mode is rarely possibly never used today. Further, by dividing ORT into UNIX-style apps for each function, an operator can easily see what results would be before running them.
- **Revalidate Mode**. ORT is now fast enough to make a separate Revalidate unnecessary. It should always check and apply all files.
- **Package Installation**. ORT will cease to perform this. OS (RPM) package installation will no longer be done by Traffic Control, but rather by whoever or whatever is managing the machine and operating system (Ansible, Puppet, human system administrators, etc).
    - Whatever is managing the other hundreds of packages on the operating system should also manage ATS and its plugins. ORT's job is to manage Traffic Control configuration data, not the operating system.

#### Additional Utilities.

Shell scripts which are “one-liners” pipelining common operations should be provided with the OS Package. TODO: add a list of scripts under Implementation heading.

Additionally, a .pl script which emulates the existing ORT behavior will be provided in the old location, to preserve backwards-compatibility. This script should be very small, and essentially translate old calls and flags to the new Aggregator.

#### Notes

I started to list requirements for all apps, such as unit tests, integration tests, modular design, argument/manpage for usage info, etc. But then I realized I was just listing good design principles. So I decided to omit that.

### Traffic Portal Impact
None.

### Traffic Ops Impact
None.

#### REST API Impact
None.

#### Client Impact
None.

#### Data Model / Database Impact
None.

### ORT Impact
Completely rewrites ORT. Backward-compatibility for safe upgrades will be preserved.

The Interface will be similar, but compatibility is not a goal. Further, LSB-compliant options and parameters is a goal, and will require incompatibility.

A `traffic_ops_ort.pl` script will be provided, whose interface _does_ preserve backward compatibility, and calls the new ORT "aggregator". This will preserve existing CRON or other tools users are using with ORT, and make an upgrade not break a production system.

### Traffic Monitor Impact
None.

### Traffic Router Impact
None.

### Traffic Stats Impact
None.

### Traffic Vault Impact
None.

### Documentation Impact
Rewrite will provide MAN pages for ORT. ORT arguments are not currently documented, and that will not be changed, to avoid duplicate documentation. TC "read-the-docs" may include a small comment pointing users to the MAN page.

### Testing Impact

ORT does not currently include an Integration Test framework. Creating one is orthogonal to this project, but this project is a good opportunity to do so at the same time.

ORT Integration Tests need exactly the same thing as the Traffic Ops API tests: a running TO instance populated with data. It is suggested that an ORT Integration Test framework copy or abstract the TO API Tests.

### Performance Impact
Not significant. ORT in master as of this writing takes less than a minute to run, on a large CDN with at least 1000 Delivery Services and 1000 Servers. The Rewrite should still take less than a minute to run.

It should be a goal for ORT to take less than 15 seconds to run, given a fast Traffic Ops. Because a faster ORT means TC changes propogate faster, which is ideal. But ORT is not in the Request Path, so performance is not critical to TC operation. Further, the previous release of ORT took 5-8 minutes on a large CDN.

### Security Impact
None.

### Upgrade Impact
None. The rewrite will provide a compatibility script, so existing tools and runs continue to work as before.

### Operations Impact
Operators should learn the new real ORT interface, to better operate Traffic Control. However, it's not strictly required, they can continue to use the compatibility script for the immediate future.

### Developer Impact
Developers on ORT will have to learn the new codebase and write in Go.

The new code should not be logically larger or more complex than the existing ORT. The Go language is vertically verbose, but this project does not significantly increase the actual logic, or add new logical features.

## Alternatives
There are innumerable ways ORT could be rewritten. Some possibilities:
- Transliterated in-place to another language.
    - Disadvantages: Monolithic, more difficult to develop, more difficult to understand each component and "thing" done, more difficult to understand all things done.
    - Advantages: faster, depending on the developer. Execution may be faster (avoids Inter-Process Communication).
- Refactored into the UNIX Philosophy with the current Perl code.
    - Disadvantages: still Perl, still unsafe, not compiled, archaic, difficult to write and read.
    - Advantages: faster, depending on the developer.
- Written as a Service.
    - Disadvantages: another single-point-of-failure, another running app Operators have to be aware of and maintain and monitor, uses constant resources like CPU and memory on every Cache, if it communicates with Traffic Ops it potentially creates a undirected graph of the TC network communication which is much more difficult to comprehend and ensure correctness.
    - More power, can receive requests from Traffic Ops, can run irrespective of CRON or operator deployment.
- Removed, changed to a server orchestration tool like Ansible or Puppet.
    - Disadvantages: radical change, makes TC dependent on an orchestration tool and platform, tool will require nearly as much code to do the same thing in its own instruction language, requires operators to learn a particular and very large and complex third-party tool, requires even small CDN Operators to deploy a large orchestration tool they don't otherwise need to learn or deploy, requires CDN operators with a different orchestration tool to deploy and learn multiple large redundant systems.
    - Advantages: Simpler, an orchestration config language is (hopefully) simpler than a Turing-Complete app. Easier for Operators to learn, if they're already using the orchestrator everywhere else.

## Dependencies
None anticipated. Small Go library dependencies may be needed, which will be considered as they arise. Any dependencies should be compile-time.

## References
https://en.wikipedia.org/wiki/Unix_philosophy

https://refspecs.linuxfoundation.org/lsb.shtml

https://www.gnu.org/prep/standards/html_node/Command_002dLine-Interfaces.html

https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html

https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html#tag_12_02
