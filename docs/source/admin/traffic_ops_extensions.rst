.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

.. _admin-to-ext-script:

Managing Traffic Ops Extensions
*******************************

Each script is a separate bash script located in ``$TO_HOME/bin/checks/``. 

The extensions must be registered with Traffic Ops in order to display a column on the Server Check page. The list of currently registered extensions can be listed by running ``/opt/traffic_ops/app/bin/extensions -a``.

The below extensions are automatically registered with the Traffic Ops database (``to_extension`` table) at install time (see ``traffic_ops/app/db/seeds.sql``). However, cron must still be configured to run these checks periodically. 

The scripts are called as follows: ::

  
  $TO_HOME/bin/checks/To<name>Check.pl  -c "{\"base_url\": \",https://\"<traffic_ops_ip>\", \"check_name\": \"<check_name>\"}" -l <log level>
  where:

  <name> is the type of check script
  <traffic_ops_ip> is the IP address of the Traffic Ops Server
  <check_name> is the name of the check. For example: CDU, CHR, DSCP, MTU, etc...
  <log_level> is between 1 and 4, with 4 being the most verbose. This field is optional


Example Cron File
=================
Edit with ``crontab -e``. You may need to adjust the path to your $TO_HOME to match your system.

::

   
 PERL5LIB=/opt/traffic_ops/app/local/lib/perl5:/opt/traffic_ops/app/lib
 
 0 * * * * /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G6\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 9 * * * * /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ILO\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 18 * * * * /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 27 * * * * /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\", \"select\": [\"hostName\", \"domainName\"]}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 36 * * * * /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1g
 
 45 * * * * /opt/traffic_ops/app/bin/checks/ToMTUCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 10 * * * * /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 0,15,30,45 * * * * /opt/traffic_ops/app/bin/checks/ToCHRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"CHR\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 20 * * * * /opt/traffic_ops/app/bin/checks/ToCDUCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"CDU\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
 
 40 * * * * /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1

   
