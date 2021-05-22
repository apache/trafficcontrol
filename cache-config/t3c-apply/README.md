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

t3c-apply - Traffic Control Cache Configuration applicator

# SYNOPSIS

t3c-apply [-2bchIpsSvW] [-D seconds] [-d location] [-e location] [-g &lt;yes|no|auto&gt;] [-H hostname] [-i location] [-l seconds] [-M location] [-m &lt;badass|report|revalidate|syncds&gt;] [-P password] [-r retries] [-R path] [-T seconds] [-t milliseconds] [-u url] [-U username] [-V versions] [-w &lt;true|false&gt;]

[&#45;&#45;help]

# DESCRIPTION

The t3c-apply command is a transliteration of traffic_ops_ort.pl script to the go language. It is designed to replace the traffic_ops_ort.pl perl script and it is used to apply configuration from Traffic Control, stored in Traffic Ops, to the cache.

Typical usage is to install t3c on the cache machine, and then run it periodically via a CRON job.

# OPTIONS

-2, --default-client-enable-h2

    Whether to enable HTTP/2 on Delivery Services by default, if
    they have no explicit Parameter. This is irrelevant if ATS
    records.config is not serving H2. If omitted, H2 is
    disabled.

-b, --dns-local-bind

    [true | false] whether to use the server's Service Addresses
    to set the ATS DNS local bind address

-c, --disable-parent-config-comments

    Whether to disable verbose parent.config comments. Default
    false.

-D, --dispersion=value

    [seconds] wait a random number of seconds between 0 and
    [seconds] before starting, default 300 [300]

-d, --log-location-debug=value

    Where to log debugs. May be a file path, stdout, stderr, or
    null, default ''

-e, --log-location-error=value

    Where to log errors. May be a file path, stdout, stderr, or
    null, default stderr [stderr]

-g, --git=value

    Create and use a git repo in the config directory. Options
    are yes, no, and auto. If yes, create and use. If auto, use
    if it exist. Default is auto. [auto]

-H, --cache-host-name=value

    Host name of the cache to generate config for. Must be the
    server host name in Traffic Ops, not a URL, and not the FQDN

-h, --help

    Print usage information and exit

-i, --log-location-info=value

    Where to log info. May be a file path, stdout, stderr, or
    null, default stderr [stderr]

-I, --traffic-ops-insecure

    [true | false] ignore certificate errors from Traffic Ops

-l, --login-dispersion=value

    [seconds] wait a random number of seconds between 0 and
    [seconds] before login to traffic ops, default 0

-M, --maxmind-location=value

    URL of a maxmind gzipped database file, to be installed into
    the trafficserver etc directory.

-m, --run-mode=value

    [badass | report | revalidate | syncds] run mode, default is
    'report' [report]

-p, --reverse-proxy-disable

    [false | true] bypass the reverse proxy even if one has been
    configured default is false

-P, --traffic-ops-password=value

    Traffic Ops password. Required. May also be set with the
    environment variable TO_PASS

-r, --num-retries=value

    [number] retry connection to Traffic Ops URL [number] times,
    default is 3 [3]

-R, --trafficserver-home=value

    Trafficserver Package directory. May also be set with the
    environment variable TS_HOME

-s, --skip-os-check

    [false | true] skip os check, default is false

-S, --syncds-updates-ipallow

    Whether syncds mode will update ipallow. This exists because
    ATS had a bug where reloading after changing ipallow would
    block everything. Default is false.

-T, --reval-wait-time=value

    [seconds] wait a random number of seconds between 0 and
    [seconds] before revlidation, default is 60 [60]

-t, --traffic-ops-timeout-milliseconds=value

    Timeout in milli-seconds for Traffic Ops requests, default
    is 30000 [30000]

-u, --traffic-ops-url=value

    Traffic Ops URL. Must be the full URL, including the scheme.
    Required. May also be set with the environment variable
    TO_URL

-U, --traffic-ops-user=value

    Traffic Ops username. Required. May also be set with the
    environment variable TO_USER

-V, --default-client-tls-versions=value

    Comma-delimited list of default TLS versions for Delivery
    Services with no Parameter, e.g.
    --default-tls-versions='1.1,1.2,1.3'. If omitted, all
    versions are enabled.

-v, --omit-via-string-release

    Whether to set the records.config via header to the ATS
    release from the RPM. Default true.

-w, --log-location-warning=value

    Where to log warnings. May be a file path, stdout, stderr,
    or null, default stderr [stderr]

-W, --wait-for-parents

    [true | false] do not update if parent_pending = 1 in the
    update json. default is false, wait for parents

# MODES

The `t3c-apply` app can be run in a number of modes.

The syncds mode is the normal mode of operation, which should typically be run periodically via cron or a similar tool.

The badass mode is typically an emergency-fix mode, which will override and replace all files with the configuration generated from the current Traffic Ops data, regardless whether `t3c-apply` (presumably incorrectly) thinks the files need updating or not. It is recommended to run this mode when something goes wrong, and the configuration on the cache is incorrect, and the data in Traffic Ops and config generation is believed to be correct. It is not recommended to run this in normal operation; use syncds mode for normal operation.

The revalidate mode will apply Revalidations from Traffic Ops (regex_revalidate.config) but no other configuration. This mode was intended to quickly apply revalidations when `t3c-apply` took a long time to run. It is less relevant with the current speed of `t3c-apply` but may still be useful on slow networks or very large deployments.

mode        | description
----------- | ----------------------------------------------------------------
report      | prints config differences and exits (default)
badass      | attempts to fix all config differences that it can
syncds      | syncs delivery services with what is configured in Traffic Ops
revalidate  | checks for updated revalidations in Traffic Ops and applies them

# BEHAVIOR

When `t3c-apply` is run, it will:

1. Delete all of its temporary directories over a week old. Currently, the base temp directory is hard-coded to /tmp/ort.
1. Determine if Updates have been Queued on the server (by checking the Server's Update Pending or Revalidate Pending flag in Traffic Ops).
    1. If Updates were not queued and the script is running in syncds mode (the normal mode), exit.
1. Get the config files from Traffic Ops, via t3c-generate.
1. Process CentOS Yum packages.
    1. These are specified via Parameters on the Server's Profile, with the Config File 'package', where the Parameter Name is the package name, and the Parameter Value is the package version.
    1. Uninstall any packages which are installed but whose version does not match.
    1. Install all packages in the Server Profile.
1. Process chkconfig directives.
    1. These are specified via Parameters on the Server's Profile, with the Config File 'chkconfig', where the Parameter Name is the package name, and the Parameter Value is the chkconfig directive line.
    1. All chkconfig directives in the Server's Profile are applied to the CentOS chkconfig.
    1. **NOTE** the default profiles distributed by Traffic Control have an ATS chkconfig with a runlevel before networking is enabled, which is likely incorrect.
    1. **NOTE** this is not used by CentOS 7+ and ATS 7+. SystemD does not use chkconfig, and ATS 7+ uses a SystemD script not an init script.
1. Process each config file
    1. If `t3c-apply` is in revalidate mode, this will only be regex_revalidate.config
    1. Perform any special processing. See [Special Processing](#special-processing).
    1. If a file exists at the path of the file, load it from disk and compare the two.
    1. If there are no changes, don't apply the new file.
    1. If there are changes, backup the existing file in the temp directory, and write the new file.
1. If configuration was changed which requires an ATS reload to apply, perform a service reload of ATS.
1. If configuration was changed which requires an ATS restart to apply, and `t3c-apply` is in badass mode, perform a service restart of ATS.
1. If a sysctl.conf config file was changed, and `t3c-apply` is in badass mode, run `sysctl -p`.
1. If a ntpd.conf config file was changed, and `t3c-apply` is in badass mode, perform a service restart of ntpd.
1. Update Traffic Ops to unset the Update Pending or Revalidate Pending flag of this Server.

# SPECIAL PROCESSING

Certain config files perform extra processing.

Global replacements

    All config files have certain text directives replaced. This is done by the t3c-generate config generator before the file is returned to `t3c-apply`.

    __SERVER_TCP_PORT__ is replaced with the Server's Port from Traffic Ops; unless the server's port is 80, 0, or null, in which case any occurrences preceded by a colon are removed.
    __CACHE_IPV4__      is replaced with the Server's IP address from Traffic Ops.
    __HOSTNAME__        is replaced with the Server's (short) HostName from Traffic Ops.
    __FULL_HOSTNAME__   is replaced with the Server's HostName, a dot, and the Server's DomainName from Traffic Ops (i.e. the Server's Fully Qualified Domain Name).
    __RETURN__          is replaced with a newline character, and any whitespace before or after it is removed.

remap.config

    The `t3c-apply` app processes `##OVERRIDE##` directives in the remap.config file.

    The ##OVERRIDE## template string allows the Delivery Service Raw Remap Text field to override to fully override the Delivery Service’s line in the remap.config ATS configuration file, generated by Traffic Ops. The end result is the original, generated line commented out, prepended with ##OVERRIDDEN## and the ##OVERRIDE## rule is activated in its place. This behavior is used to incrementally deploy plugins used in this configuration file. Normally, this entails cloning the Delivery Service that will have the plugin, ensuring it is assigned to a subset of the cache servers that serve the Delivery Service content, then using this ##OVERRIDE## rule to create a remap.config rule that will use the plugin, overriding the normal rule. Simply grow the subset over time at the desired rate to slowly deploy the plugin. When it encompasses all cache servers that serve the original Delivery Service’s content, the “override Delivery Service” can be deleted and the original can use a non-##OVERRIDE## Raw Remap Text to add the plugin.

50-ats.rules

    This is presumed to be a udev file for devices which are block devices to be used as disk storage by ATS.

    The `t3c-apply` app verifies all devices in the file are owned by the owner listed in the file, and logs errors otherwise.

    The `t3c-apply` app verifies all devices in the file do not have filesystems. If any device has a filesystem, `t3c-apply` assumes it was a mistake to assign as an ATS storage device, and logs a fatal error.

# AUTHORS

The t3c application is maintained by Apache Traffic Control project. For help, bug reports, contributing, or anything else, see:

https://trafficcontrol.apache.org/

https://github.com/apache/trafficcontrol
