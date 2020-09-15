# ORT

ORT, traffic_ops_ort.pl, is used to apply configuration from Traffic Control, stored in Traffic Ops, to the cache.

Typical usage is to install ORT on the cache machine, and then run it periodically via a CRON job.

## Options

ORT has the following command-line options:


option                  | default | description
----------------------- | ------- | ---
dispersion              | 300     | wait a random number of seconds between 0 and <time> before starting.
retries                 | 5       | retry connection to Traffic Ops URL <number> times.
wait_for_parents        | 1       | do not update if parent_pending = 1 in the update json.
login_dispersion        | 0       | wait a random number of seconds between 0 and <time> before login.
rev_proxy_disable       | 0       | bypass the reverse proxy even if one has been configured.
skip_os_check           | 0       | bypass the check for a supported CentOS version.
override_hostname_short | ""      | override the short hostname of the OS for config generation.

# Modes

ORT can be run in a number of modes.

The syncds mode is the normal mode of operation, which should typically be run periodically via cron or a similar tool.

The badass mode is typically an emergency-fix mode, which will override and replace all files with the configuration generated from the current Traffic Ops data, regardless whether ORT (presumably incorrectly) thinks the files need updating or not. It is recommended to run this mode when something goes wrong, and the configuration on the cache is incorrect, and the data in Traffic Ops and config generation is believed to be correct. It is not recommended to run this in normal operation; use syncds mode for normal operation.

The revalidate mode will apply Revalidations from Traffic Ops (regex_revalidate.config) but no other configuration. This mode was intended to quickly apply revalidations when ORT took a long time to run. It is less relevant with ORT's current speed, but may still be useful on slow networks or very large deployments.

mode        | description
------------| ---
interactive | asks questions during config process
report      | prints config differences and exits
badass      | attempts to fix all config differences that it can
syncds      | syncs delivery services with what is configured in Traffic Ops
revalidate  | checks for updated revalidations in Traffic Ops and applies them

# Behavior

When ORT is run, it will:

1. Delete all of its temporary directories over a week old. Currently, the base temp directory is hard-coded to /tmp/ort.
1. Determine if Updates have been Queued on the server (by checking the Server's Update Pending or Revalidate Pending flag in Traffic Ops).
    1. If Updates were not queued and the script is running in syncds mode (the normal mode), exit.
1. Get the config files from Traffic Ops, via atstccfg.
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
    1. If ORT is in revalidate mode, this will only be regex_revalidate.config
    1. Perform any special processing. See [Special Processing](#special-processing).
    1. If a file exists at the path of the file, load it from disk and compare the two.
    1. If there are no changes, don't apply the new file.
    1. If there are changes, backup the existing file in the temp directory, and write the new file.
1. If configuration was changed which requires an ATS reload to apply, perform a service reload of ATS.
1. If configuration was changed which requires an ATS restart to apply, and ORT is in badass mode, perform a service restart of ATS.
1. If a sysctl.conf config file was changed, and ORT is in badass mode, run `sysctl -p`.
1. If a ntpd.conf config file was changed, and ORT is in badass mode, perform a service restart of ntpd.
1. Update Traffic Ops to unset the Update Pending or Revalidate Pending flag of this Server.

# Special Processing

Certain config files perform extra processing.

## Global replacements

All config files have certain text directives replaced. This is done by the atstccfg config generator before the file is returned to ORT.

* `__SERVER_TCP_PORT__` is replaced with the Server's Port from Traffic Ops; unless the server's port is 80, 0, or null, in which case any occurrences preceded by a colon are removed.
* `__CACHE_IPV4__` is replaced with the Server's IP address from Traffic Ops.
* `__HOSTNAME__` is replaced with the Server's (short) HostName from Traffic Ops.
* `__FULL_HOSTNAME__` is replaced with the Server's HostName, a dot, and the Server's DomainName from Traffic Ops (i.e. the Server's Fully Qualified Domain Name).
* `__RETURN__` is replaced with a newline character, and any whitespace before or after it is removed.

## remap.config

ORT processes `##OVERRIDE##` directives in the remap.config file.

The ##OVERRIDE## template string allows the Delivery Service Raw Remap Text field to override to fully override the Delivery Service’s line in the remap.config ATS configuration file, generated by Traffic Ops. The end result is the original, generated line commented out, prepended with ##OVERRIDDEN## and the ##OVERRIDE## rule is activated in its place. This behavior is used to incrementally deploy plugins used in this configuration file. Normally, this entails cloning the Delivery Service that will have the plugin, ensuring it is assigned to a subset of the cache servers that serve the Delivery Service content, then using this ##OVERRIDE## rule to create a remap.config rule that will use the plugin, overriding the normal rule. Simply grow the subset over time at the desired rate to slowly deploy the plugin. When it encompasses all cache servers that serve the original Delivery Service’s content, the “override Delivery Service” can be deleted and the original can use a non-##OVERRIDE## Raw Remap Text to add the plugin.

## 50-ats.rules

This is presumed to be a udev file for devices which are block devices to be used as disk storage by ATS.

ORT verifies all devices in the file are owned by the owner listed in the file, and logs errors otherwise.

ORT verifies all devices in the file do not have filesystems. If any device has a filesystem, ORT assumes it was a mistake to assign as an ATS storage device, and logs a fatal error.

# Logging

ORT outputs its immediate log to `stdout`, therefore it will log wherever you direct it. The recommended system location is `/var/log/ort/ort.log`.

ORT uses the helper tool `atstccfg` to generate config files. Its log is output at `/var/log/ort/atstccfg.log`.

# Trivia

ORT stands for "Operational Readiness Test." The acronym is a legacy artifact and does not reflect the current purpose, which is to apply configuration from Traffic Ops.
