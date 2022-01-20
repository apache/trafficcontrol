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

Also on each polling cycle the configuration file, **tc-health-client.json** is 
checked and a new config is reloaded if the file has changed since the last 
polling cycle.  The **Traffic Monitors** list is refreshed from **Traffic Ops**.

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

Sample configuarion file:

```
  {
    "cdn-name": "over-the-top",
    "enable-active-markdowns": false,
    "reason-code": "active",
    "to-credential-file": "/etc/credentials",
    "to-url": "https://tp.cdn.com:443", 
    "to-request-timeout-seconds": "5s",
    "tm-poll-interval-seconds": "60s",
    "tm-proxy-url", "http://sample-http-proxy.cdn.net:80",
    "tm-update-cycles": 5,
    "unavailable-poll-threshold": 2,
    "trafficserver-config-dir": "/opt/trafficserver/etc/trafficserver",
    "trafficserver-bin-dir": "/opt/trafficserver/bin",
  }
```

### cdn-name 

  The name of the CDN that the Traffic Server host is a member of.

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

### tm-update-cycles

  Each time a polling cycle completes a count is incremented. When the count
  reaches **tm-update-cycles**, TrafficOps is polled for a new list of available
  TrafficMonitors for the CDN and the poll count is reset to 0.

### unavailable-poll-threshold

  This controls when an unhealthy parent is marked down.  An unhealthy parent
  will be marked down when the number of consecutive polls reaches this threshold
  with the parent reported as unhealthy.  The default threshold is 2.

### trafficserver-config-dir

  The location on the host where **Traffic Server** configuration files are 
  located.

### trafficserver-bin-dir

  The location on the host where **Traffic Server** **traffic_ctl** tool may
  be found.

# Files

* /etc/trafficcontrol/tc-health-client.json
* /etc/logrotate.d/tc-health-client-logrotate
* /usr/bin/tc-health-client
* /usr/lib/systemd/system/tc-health-client.service
* /var/log/trafficcontrol/tc-health-client.json
* Traffic Server **parent.config**
* Traffic Server **strategies.yaml**
* Traffic Server **traffic_ctl** command
  
   
