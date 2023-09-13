% tc-health-client(1) tc-health-client 6.2.0 | ATC tc-health-client Manual
%
% 2022-03-09
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
  !!!

-->
# NAME

tc-health-client - Traffic Control Health Client service

# SYNOPSIS

tc-health-client [-f config-file]  -h  [-l logging-directory]  -v 

# DESCRIPTION

The tc-health-client command is used to manage **Apache Traffic Server** parents on a
host running **Apache Traffic Server**.  The command should be started by **systemd** 
and run as a service. On startup, the command reads its default configuration file
**/etc/trafficcontrol/tc-health-client.json**.  After reading the config
file it polls the configured **Traffic OPs** to obtain a list of **Traffic Monitors**
for the configured **CDN** and begins polling the available **Traffic Monitors** for
Traffic Server cache statuses.

On each polling cycle, defined in the configuration file, the Traffic Server parent
statuses are updated from the Traffic Server **parent.config**, **strategies.yaml** 
files, and the Traffic Server **HostStatus** subsystem.  If **Traffic Monitor** has
determined that a parent utilized by the **Traffic Server** instance is un-healthy or
otherwise unavailable, the tc-health-client will utilize the **Traffic Server** 
**traffic_ctl** tool to mark down the parent host.  If a parent host is marked down 
and **Traffic Monitor** has determined that the marked down host is now available, 
the client will then utilize the **Traffic Server** tool to mark the host back up.

Any changes to **tc-health-client.json** will require sending a SIGHUP to the running
process. This will cause the tc-health-client to read the config file and load any
changes, the **Traffic Monitors** list will be refreshed from **Traffic Ops**.
**systemctl reload tc-health-client** can also be used to send a SIGHUP.  

If errors are encountered while polling a Traffic Monitor, the error is logged
and the **Traffic Monitors** list is refreshed from **Traffic Ops**.

# REQUIREMENTS

Requires Apache TrafficServer 8.1.0 or later.

# OPTIONS

-f, -\-config-file=config-file

  Specify the config file to use.
  Defaults to /etc/trafficcontro-health-client/tc-health-client.json

-h, -\-help

  Prints command line usage and exits

-l, -\-logging-dir=logging-directory

  Specify the directory where log files are kept.  The default location
  is **/var/log/trafficcontrol/**

-v, -\-verbose

  Logging verbosity.  Errors are logged to the default log file
  **/var/log/trafficcontrol/tc-health-client.log**
  To add Warnings, use -v.  To add Warnings and Informational
  logging, use -vv.  Finally you may add Debug logging using -vvv.

# CONFIGURATION

The configuration file is a **JSON** file and is looked for by default
at **/etc/trafficcontrol/tc-health-client.json**

Sample configuration file:

```
  {
    "cdn-name": "over-the-top",
    "enable-active-markdowns": false,
    "reason-code": "active",
    "to-credential-file": "/etc/credentials",
    "to-url": "https://tp.cdn.com:443",
    "to-request-timeout-seconds": "5s",
    "tm-poll-interval-seconds": "60s",
    "tm-proxy-url": "http://sample-http-proxy.cdn.net:80",
    "to-login-dispersion-factor": 90,
    "unavailable-poll-threshold": 2,
    "markup-poll-threshold": 1,
    "trafficserver-config-dir": "/opt/trafficserver/etc/trafficserver",
    "trafficserver-bin-dir": "/opt/trafficserver/bin",
    "poll-state-json-log": "/var/log/trafficcontrol/poll-state.json",
    "enable-poll-state-log": false,
    "parent-health-poll-ms": 10000,
    "serve-parent-health": true,
    "parent-health-service-port": 31337,
    "parent-health-log-location": "/var/log/trafficcontrol/tc-health-client_parent-health.log",
    "health-methods": ["traffic-monitor", "parent-l4", "parent-l7", "parent-service"],
    "markdown-methods": ["traffic-monitor", "parent-l4", "parent-l7", "parent-service"],
    "hostname": ""
  }
```

### cdn-name

The name of the CDN that the Traffic Server host is a member of.

### Enable Debug instructions

Debug is currently going to `/dev/null` to avoid filling up the logs. However, it can be redirect to show in the logs when debugging is needed. In the desired machine, please add a 3rd `v` to `tc-health-client.service` resulting in `-vvv` which enables debugging live to trobleshoot an issue and **change it back to `-vv` once debugging is no londer needed**. 

Steps:
Cd into `/usr/lib/systemd/system/`
To edit run `sudo vi tc-health-client.service`
Save and quit with `wq!`
Reload with `sudo systemctl daemon-reload`
Run `sudo systemctl stop tc-health-client` and wait for the next prompt to show
Then run `sudo systemctl start tc-health-client`

### enable-active-markdowns

When enabled, the client will actively mark down Traffic Server parents.
When disabled, the client will only log that it would have marked down
Traffic Server parents.  Down Parents are always marked UP if Traffic Monitor
reports them available irregardless of this setting.

### reason-code

Use the reason code **active** or **local** when marking down Traffic Server
hosts in the Traffic Server **HostStatus** subsystem.

### to-credential-file

The file where **Traffic Ops** credentials are read.  The file should define the
following variables:

* TO_URL="https://trafficops.cdn.com"
* TO_USER="touser"
* TO_PASS="touser_password"

### to-url

The **Traffic Ops** URL

### to-request-timeout-seconds

The time in seconds to wait for a query response from both **Traffic Ops** and
the **Traffic Monitors**

### tm-poll-interval-seconds

The polling interval in seconds used to update **Traffic Server** parent
status.

### tm-proxy-url

If not nil, all Traffic Monitor requests will be proxied through this 
proxy endpoint.  This is useful when there are large numbers of caches
polling a Traffic Monitor and you wish to funnel queries through a caching
proxy server to limit direct direct connections to Traffic Monitor.

### to-login-dispersion-factor

This is used to calculate TrafficOps login dispersion.  It is related to the
**tm-poll-interval-seconds**.  The login dispersion is computed by multiplying
**tm-poll-interval-seconds** by the **to-login-dispersion-factor**.  For example
if to-login-dispersion-factor is 90 and the **tm-poll-interval-seconds** is 10s
the the dispersion modulo window is 900s.

### unavailable-poll-threshold

This controls when an unhealthy parent is marked down.  An unhealthy parent
will be marked down when the number of consecutive polls reaches this threshold
with the parent reported as unhealthy.  The default threshold is 2.

### markup-poll-threshold

This controls when a healthy parent is marked up.  An healthy parent
will be marked up when the number of consecutive polls reaches this threshold
with the parent reported as healthy.  The default threshold is 1.

### trafficserver-config-dir

The location on the host where **Traffic Server** configuration files are
located.

### trafficserver-bin-dir

The location on the host where **Traffic Server** **traffic_ctl** tool may
be found.

### poll-state-json-log ###

The full path to the polling state file which contains information
about the current status of parents and the health client configuration.
Polling state data is written to this file after each polling cycle when
enabled, see **enable-poll-state-log**

### enable-poll-state-log ###

Enable writing the Polling state to the **poll-state-json-log** after
eache polling cycle.  Default **false**, disabled

### markdown-min-interval-ms ###

Minimum interval between markdown processing in milliseconds. When health polls finish, they automatically signal the markdown service to mark down accordingly. This is the minimum time to wait between processing, to avoid too much processing. To always process markdowns as soon as every health poll finishes, set to 0. Default is 5 seconds.

### parent-health-l4-poll-ms ###

Interval to poll for parent health via L4 in milliseconds.

### parent-health-l7-poll-ms ###

Interval to poll for parent health via L7 in milliseconds.

### parent-health-service-poll-ms ###

Interval to poll for parent health from parents' health service in milliseconds.

### parent-health-service-port ###

The port to serve the JSON parent health data over HTTP. To disable serving, set to < 1. Default is 0, disabled.

### parent-health-log-location ###

The location to log parent health changes. May be stdout, stderr, null, or a file path.

### health-methods ###

The health methods to poll. Options are 'traffic-monitor', 'parent-l4', 'parent-l7', and 'parent-service'.

Traffic Monitor requests Traffic Monitor and uses its boolean CRStates API for health.

Parent L4 polls all parents via a HTTP request. Any valid response, including HTTP error codes, is considered a success. Only a failed HTTP request is considered unhealthy.

Parent L7 polls all parents via a TCP SYN. Any valid TCP ACK response within the timeout is considered healthy. Failure to receive an ACK before the timeout is considered unhealthy. Note this also sends a TCP Reset to aid the host in quickly releasing resources.

Parent Service polls the tc-health-client parent service, on the same port as this host's parent service, and uses a heuristic of the parent's own available parents to determine health. The heuristic is currently 50%, but that may change or be made configurable in the future. That is, if more than 50% of a parent cache's own parents are unavailable, the parent is unhealthy.

### markdown-methods ###

Markdown methods are the health methods to consider when marking down parents. See health-methods. Hence, a method may be polled via health-methods and logged and served for informational purposes, without using that method to mark down parents.

### num-health-workers ###

The number of worker microthreads (goroutines) per health poll method.

Note this only applies to Parent L4, Parent L7, and Parent Service health; the Traffic Monitor health poll is a single HTTP request and thus doesn't need workers.

### monitor-strategies-peers ###

In the **strategies.yaml** there are peers, setting this setting to `true`, which is default, it will monitor peers. If setting the value to `false` it will filter out all `&peer#` anchors within the **strategies.yaml** host section of the file.

### hostname ###

# Files

* /etc/trafficcontrol/tc-health-client.json
* /etc/logrotate.d/tc-health-client-logrotate
* /usr/bin/tc-health-client
* /usr/lib/systemd/system/tc-health-client.service
* /var/log/trafficcontrol/tc-health-client.json
* Traffic Server **parent.config**
* Traffic Server **strategies.yaml**
* Traffic Server **traffic_ctl** command
