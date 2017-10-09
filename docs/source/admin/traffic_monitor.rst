.. 
.. 
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
.. 
..     http://www.apache.org/licenses/LICENSE-2.0
.. 
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
.. 

***************************************
Traffic Monitor Administration (Legacy)
***************************************

.. _rl-tm-java:

* These instructions are for the legacy Java Traffic Monitor, for the new Golang version, see :ref:`here <rl-tm-golang>`.

Installing Traffic Monitor
==========================
The following are requirements to ensure an accurate set up:

* CentOS 6
* 4 vCPUs
* 8GB RAM
* Successful install of Traffic Ops
* Tomcat
* Administrative access to the Traffic Ops
* Physical address of the site
* perl-JSON
* perl-WWW-Curl

#. Add the Traffic Monitor server into Traffic Ops using 'Servers' -> 'Add Server'. Set the 'Type' field to 'RASCAL'.

#. Make sure the FQDN of the Traffic Monitor is resolvable in DNS.

#. Get the Traffic Monitor RPM.

   Sample command: ::

      wget http://traffic-control-cdn.net/downloads/1.7.0/RELEASE-1.7.0/traffic_monitor-1.7.0-3908.5b77f60f.el6.x86_64.rpm

#. Install Traffic Monitor and Perl modules: ::

    sudo yum -y install traffic_monitor-*.rpm perl-JSON perl-WWW-Curl

#. Take the config from Traffic Ops: ::

    sudo /opt/traffic_monitor/bin/traffic_monitor_config.pl https://<traffic-ops-URL> <traffic-ops-user>:<traffic-ops-password> prompt

   Sample session: ::

    traffic_mon # /opt/traffic_monitor/bin/traffic_monitor_config.pl https://traffic-ops.cdn.kabletown.net admin:kl0tevax prompt
    DEBUG: traffic_ops selected: https://traffic-ops.cdn.kabletown.net
    DEBUG: traffic_ops login: admin:kl0tevax
    DEBUG: Config write mode: prompt
    DEBUG: Found profile from traffic_ops: RASCAL_CDN
    DEBUG: Found CDN name from traffic_ops: kabletown_cdn
    DEBUG: Found location for rascal-config.txt from traffic_ops: /opt/traffic_monitor/conf
    WARN: Param not in traffic_ops: allow.config.edit                        description: Allow the running configuration to be edited through the UI                                                              Using default value of: false
    WARN: Param not in traffic_ops: default.accessControlAllowOrigin         description: The value for the header: Access-Control-Allow-Origin for published jsons... should be narrowed down to TMs              Using default value of: *
    WARN: Param not in traffic_ops: default.connection.timeout               description: Default connection time for all queries (cache, peers, TM)                                                               Using default value of: 2000
    WARN: Param not in traffic_ops: hack.forceSystemExit                     description: Call System.exit on shutdown                                                                                             Using default value of: false
    WARN: Param not in traffic_ops: hack.peerOptimistic                      description: The assumption of a caches availability when unknown by peers                                                            Using default value of: true
    WARN: Param not in traffic_ops: hack.publishDsStates                     description: If true, the delivery service states will be included in the CrStates.json                                               Using default value of: true
    WARN: Param not in traffic_ops: health.ds.interval                       description: The polling frequency for calculating the deliveryService states                                                         Using default value of: 1000
    WARN: Param not in traffic_ops: health.ds.leniency                       description: The amount of time before the deliveryService disregards the last update from a non-responsive cache                     Using default value of: 30000
    WARN: Param not in traffic_ops: health.event-count                       description: The number of historical events that will be kept                                                                        Using default value of: 200
    WARN: Param not in traffic_ops: health.polling.interval                  description: The polling frequency for getting the states from caches                                                                 Using default value of: 5000
    WARN: Param not in traffic_ops: health.startupMinCycles                  description: The number of query cycles that must be completed before this Traffic Monitor will start reporting                       Using default value of: 2
    WARN: Param not in traffic_ops: health.timepad                           description: A delay between each separate cache query                                                                                Using default value of: 10
    WARN: Param not in traffic_ops: peers.polling.interval                   description: Polling frequency for getting states from peer monitors                                                                  Using default value of: 5000
    WARN: Param not in traffic_ops: peers.polling.url                        description: The url for current, unfiltered states from peer monitors                                                                Using default value of: http://${hostname}/publish/CrStates?raw
    WARN: Param not in traffic_ops: peers.threadPool                         description: The number of threads given to the pool for querying peers                                                               Using default value of: 1
    WARN: Param not in traffic_ops: tm.auth.url                              description: The url for the authentication form                                                                                      Using default value of: https://${tmHostname}/login
    WARN: Param not in traffic_ops: tm.crConfig.json.polling.url             description: Url for the cr-config (json)                                                                                             Using default value of: https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json
    WARN: Param not in traffic_ops: tm.healthParams.polling.url              description: The url for the heath params (json)                                                                                      Using default value of: https://${tmHostname}/health/${cdnName}
    WARN: Param not in traffic_ops: tm.polling.interval                      description: The polling frequency for getting updates from TM                                                                        Using default value of: 10000
    DEBUG: allow.config.edit needed in config, but does not exist in config on disk.
    DEBUG: cdnName value on disk () does not match value needed in config (kabletown_cdn).
    DEBUG: default.accessControlAllowOrigin needed in config, but does not exist in config on disk.
    DEBUG: default.connection.timeout needed in config, but does not exist in config on disk.
    DEBUG: hack.forceSystemExit needed in config, but does not exist in config on disk.
    DEBUG: hack.peerOptimistic needed in config, but does not exist in config on disk.
    DEBUG: hack.publishDsStates needed in config, but does not exist in config on disk.
    DEBUG: health.ds.interval needed in config, but does not exist in config on disk.
    DEBUG: health.ds.leniency needed in config, but does not exist in config on disk.
    DEBUG: health.startupMinCycles needed in config, but does not exist in config on disk.
    DEBUG: health.timepad value on disk (20) does not match value needed in config (10).
    DEBUG: peers.polling.interval needed in config, but does not exist in config on disk.
    DEBUG: peers.threadPool needed in config, but does not exist in config on disk.
    DEBUG: tm.auth.password value on disk () does not match value needed in config (kl0tevax).
    DEBUG: tm.auth.username value on disk () does not match value needed in config (admin).
    DEBUG: tm.hostname value on disk () does not match value needed in config (traffic-ops.cdn.kabletown.net).
    DEBUG: Proposed traffic_monitor_config:
    {
       "traffic_monitor_config":{
          "default.accessControlAllowOrigin":"*",
          "health.startupMinCycles":"2",
          "tm.auth.password":"kl0tevax",
          "tm.auth.url":"https://${tmHostname}/login",
          "tm.healthParams.polling.url":"https://${tmHostname}/health/${cdnName}",
          "allow.config.edit":"false",
          "tm.crConfig.json.polling.url":"https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json",
          "tm.auth.username":"admin",
          "peers.polling.url":"http://${hostname}/publish/CrStates?raw",
          "health.timepad":"10",
          "hack.publishDsStates":"true",
          "default.connection.timeout":"2000",
          "health.ds.interval":"1000",
          "peers.polling.interval":"5000",
          "hack.forceSystemExit":"false",
          "health.ds.leniency":"30000",
          "cdnName":"kabletown_cdn",
          "peers.threadPool":"1",
          "tm.polling.interval":"10000",
          "health.polling.interval":"5000",
          "health.event-count":"200",
          "hack.peerOptimistic":"true",
          "tm.hostname":"traffic-ops.cdn.kabletown.net"
       }
    }
    ----------------------------------------------
    ----OK to write this config to disk? (Y/n) [n]y
    ----------------------------------------------
    ----------------------------------------------
    ----OK to write this config to disk? (Y/n) [n]Y
    ----------------------------------------------
    DEBUG: Writing /opt/traffic_monitor/conf/traffic_monitor_config.js
    traffic_mon #

#. Update the 'allow_ip' and 'allow_ip6' parameters in the profiles of all caches defined in traffic ops, both edge and mid,
   with the address of the traffic monitor being installed, so that the traffic servers will allow this Traffic Monitor
   to access the astats plugin.
   For details see :ref:`rl-param-prof` in the *Configuring Traffic Ops* section.

#. Start Tomcat: ``sudo service tomcat start`` ::

    Using CATALINA_BASE: /opt/tomcat
    Using CATALINA_HOME: /opt/tomcat
    Using CATALINA_TMPDIR: /opt/tomcat/temp
    Using JRE_HOME: /usr
    Using CLASSPATH:/opt/tomcat/bin/bootstrap.jar
    Using CATALINA_PID:/var/run/tomcat/tomcat.pid
    Starting tomcat [ OK ]

#. Configure tomcat to start automatically: ``sudo chkconfig tomcat on``

#. Verify Traffic Monitor is running by pointing your browser to port 80 on the Traffic Monitor host:

   * The 'Cache States' tab should display all Mid and Edge caches configured in Traffic Ops.
   * The 'DeliveryService States' tab should display all delivery services configured in Traffic Ops.

#. In Traffic Ops servers table, click 'Edit' for this server, then click 'Online'.


Configuring Traffic Monitor
===========================

Configuration Overview
----------------------
Traffic Monitor is configured using its JSON configuration file, ``/opt/traffic_monitor/conf/traffic_monitor_config.js``.
This file is created by ``traffic_monitor_config.pl`` script, and among other things, it contains the Traffic Ops URL and the user:password
specified during the invocation of that script.

When started, Traffic Monitor uses this basic configuration to downloads its configuration from Traffic Ops, and begins polling caches.
Once a configurable number of polling cycles completes, health protocol state is available via RESTful JSON endpoints.


Troubleshooting and log files
=============================
Traffic Monitor log files are in ``/opt/traffic_monitor/var/log/``, and tomcat log files are in ``/opt/tomcat/logs/``.
