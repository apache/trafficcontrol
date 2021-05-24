# t3c-apply

t3c-apply is a transliteration of traffic_ops_ort.pl script to the go language.
It is designed to replace the traffic_ops_ort.pl perl script and it is used to apply
configuration from Traffic Control, stored in Traffic Ops, to the cache.

Typical usage is to install `t3c` on the cache machine, and then run it periodically via a CRON job.

## Options

The `t3c-apply` app has the following command-line options:


long option                             | short | default | description
--------------------------------------- | ------| ------- | ------------------------------------------------------------------------------------
--cache-hostname=[hostname]                    | -H    | ""      | override the short hostname of the OS for config generation.
--dispersion=[seconds]                         | -D    | 300     | wait a random number of seconds between 0 and [seconds] before starting.
--login-dispersion=[seconds]                   | -l    | 0       | wait a random number of seconds between 0 and [seconds] before login.
--log-location-debug=[value]                   | -d    | stdout  | Where to log debugs. May be a file path, stdout, stderr, or null
--log-location-error=[value]                   | -e    | stdout  | Where to log errors. May be a file path, stdout, stderr, or null
--log-location-info=[value]                    | -i    | stdout  | Where to log info messages. May be a file path, stdout, stderr, or null
--log-location-warn=[value]                    | -w    | stdout  | Where to log warning messages. May be a file path, stdout, stderr, or null
--num-retries=[number]                         | -r    | 3       | retry connection to Traffic Ops URL [number] times.
--rev-proxy-disable=['true' or 'false']        | -p    | false   | bypass the reverse proxy even if one has been configured.
--reval-wait-time=[seconds]                    | -T    | 60      | wait a random number of seconds between 0 and [seconds] before revlidation
--run-mode=[mode]                              | -m    | report  | The mode of operation, where mode is [badass|report|revalidate|syncds].
--skip-os-check=['true' or 'false']            | -s    | false   | bypass the check for a supported CentOS version.
--traffic-ops-timeout-milliseconds=[ms]        | -t    | 30000   | The Traffic Ops request timeout in milliseconds.
--traffic-ops-password=[password]              | -P    | ""      | TrafficOps password. Required if not set with the environment variable TO_PASS
--traffic-ops-url=[url]                        | -u    | ""      | TrafficOps URL. Required if not set with the environment variable TO_URL
--traffic-ops-user=[username]                  | -U    | ""      | TrafficOps username. Required if not set with the environment variable TO_USER
--trafficserver-home=[directory]               | -R    | ""      | Used to specify an alternate install location for ATS, otherwise its set from the RPM.
--dns-local-bind=['true' or 'false']           | -b    | false   | set the ATS config to bind to the Server's Service Address in Traffic Ops for DNS.
--wait-for-parents=['true' or 'false']         | -W    | true    | do not update if parent_pending = 1 in the update json.
--git=['yes' or 'no' or 'auto']                | -g    | auto    | track changes in git. If yes, create and commit to a repo. If auto, commit if a repo exists.
--default-client-enable-h2=['true' or 'false'] | -2    | false   | Whether to enable HTTP/2 on Delivery Services by default, if they have no explicit Parameter.
--default-client-tls-versions=[versions]       | -v    | ""      | Comma-delimited list of default TLS versions for Delivery Services with no Parameter, e.g. '1.1,1.2,1.3'. If omitted, all versions are enabled.

# Modes

The `t3c-apply` app can be run in a number of modes.

The syncds mode is the normal mode of operation, which should typically be run periodically via cron or a similar tool.

The badass mode is typically an emergency-fix mode, which will override and replace all files with the configuration generated from the current Traffic Ops data, regardless whether `t3c-apply` (presumably incorrectly) thinks the files need updating or not. It is recommended to run this mode when something goes wrong, and the configuration on the cache is incorrect, and the data in Traffic Ops and config generation is believed to be correct. It is not recommended to run this in normal operation; use syncds mode for normal operation.

The revalidate mode will apply Revalidations from Traffic Ops (regex_revalidate.config) but no other configuration. This mode was intended to quickly apply revalidations when `t3c-apply` took a long time to run. It is less relevant with the current speed of `t3c-apply` but may still be useful on slow networks or very large deployments.

mode        | description
------------| ---
report      | prints config differences and exits (default)
badass      | attempts to fix all config differences that it can
syncds      | syncs delivery services with what is configured in Traffic Ops
revalidate  | checks for updated revalidations in Traffic Ops and applies them

# Behavior

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

# Special Processing

Certain config files perform extra processing.

## Global replacements

All config files have certain text directives replaced. This is done by the t3c-generate config generator before the file is returned to `t3c-apply`.

* `__SERVER_TCP_PORT__` is replaced with the Server's Port from Traffic Ops; unless the server's port is 80, 0, or null, in which case any occurrences preceded by a colon are removed.
* `__CACHE_IPV4__` is replaced with the Server's IP address from Traffic Ops.
* `__HOSTNAME__` is replaced with the Server's (short) HostName from Traffic Ops.
* `__FULL_HOSTNAME__` is replaced with the Server's HostName, a dot, and the Server's DomainName from Traffic Ops (i.e. the Server's Fully Qualified Domain Name).
* `__RETURN__` is replaced with a newline character, and any whitespace before or after it is removed.

## remap.config

The `t3c-apply` app processes `##OVERRIDE##` directives in the remap.config file.

The ##OVERRIDE## template string allows the Delivery Service Raw Remap Text field to override to fully override the Delivery Service’s line in the remap.config ATS configuration file, generated by Traffic Ops. The end result is the original, generated line commented out, prepended with ##OVERRIDDEN## and the ##OVERRIDE## rule is activated in its place. This behavior is used to incrementally deploy plugins used in this configuration file. Normally, this entails cloning the Delivery Service that will have the plugin, ensuring it is assigned to a subset of the cache servers that serve the Delivery Service content, then using this ##OVERRIDE## rule to create a remap.config rule that will use the plugin, overriding the normal rule. Simply grow the subset over time at the desired rate to slowly deploy the plugin. When it encompasses all cache servers that serve the original Delivery Service’s content, the “override Delivery Service” can be deleted and the original can use a non-##OVERRIDE## Raw Remap Text to add the plugin.

## 50-ats.rules

This is presumed to be a udev file for devices which are block devices to be used as disk storage by ATS.

The `t3c-apply` app verifies all devices in the file are owned by the owner listed in the file, and logs errors otherwise.

The `t3c-apply` app verifies all devices in the file do not have filesystems. If any device has a filesystem, `t3c-apply` assumes it was a mistake to assign as an ATS storage device, and logs a fatal error.

# Trivia

The "t3c" stands for "Traffic Control Cache Config."
