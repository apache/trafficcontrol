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

.. _admin-to-ext-script:

*******************************
Managing Traffic Ops Extensions
*******************************

Traffic Ops supports two types of extensions. 'Check Extensions' are analytics scripts that collect and display information as columns in the table under 'Monitor' -> 'Cache Checks' in Traffic Portal. 'Data Source Extensions' provide ways to add data to the graph views and usage APIs.

.. _to-check-ext:

Check Extensions
================
Each extension is a separate bash script located in ``$TO_HOME/bin/checks/`` on the Traffic Ops server (though practically, they can be any executable). The extensions must be registered with Traffic Ops in order to display a column on the Cache Checks page. The currently registered extensions can be listed by running ``/opt/traffic_ops/app/bin/extensions -a``. Some extensions automatically registered with the Traffic Ops database (``to_extension`` table) at install time (see ``traffic_ops/app/db/seeds.sql``). However, cron must still be configured to run these checks periodically. The extensions are called like so:

.. code-block:: shell

	$TO_HOME/bin/checks/<name>  -c "{\"base_url\": \",https://\"<traffic_ops_ip>\", \"check_name\": \"<check_name>\"}" -l <log level>

:name: The basename of the extension executable
:traffic_ops_ip: The IP address or Fully Qualified Domain Name (FQDN) of the Traffic Ops server
:check_name: The name of the check e.g. ``CDU``, ``CHR``, ``DSCP``, ``MTU``, etc...
:log_level: A whole number between 1 and 4 (inclusive), with 4 being the most verbose. This field is optional

.. _to-datasource-ext:

Data Source Extensions
======================
Data Source Extensions work in much the same way as `Check Extensions`_, but are implemented differently. Rather than being a totally external executable, a Data Source Extension *must* be written in Perl 5, as they are injected via manipulation of the ``$PERL5LIB`` environment variable. These extensions are not very well-documented (as you may be able to tell), and support for extending them may be phased out in future releases.

Example Cron File
=================
The cron file should be edited by running  ``crontab -e`` as the ``traffops`` user, or with ``sudo``. You may need to adjust the path to your ``$TO_HOME`` to match your system.

.. code-block:: shell

	PERL5LIB=/opt/traffic_ops/app/local/lib/perl5:/opt/traffic_ops/app/lib

	# IPv4 ping examples - The 'select: ["hostName","domainName"]' works but, if you want to check DNS resolution use FQDN.
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"select\": [\"hostName\",\"domainName\"]}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"select\": \"ipAddress\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"name\": \"IPv4 Ping\", \"select\": \"ipAddress\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

	# IPv6 ping examples
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G6\", \"name\": \"IPv6 Ping\", \"select\": \"ip6Address\", \"syslog_facility\": \"local0\"}" >/dev/null 2>&1
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G6\", \"select\": \"ip6Address\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1

	# iLO ping
	18 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ILO\", \"select\": \"iloIpAddress\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	18 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ILO\", \"name\": \"ILO ping\", \"select\": \"iloIpAddress\", \"syslog_facility\": \"local0\"}" >/dev/null 2>&1

	# MTU ping
	45 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"select\": \"ipAddress\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	45 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"select\": \"ip6Address\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	45 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"name\": \"Max Trans Unit\", \"select\": \"ipAddress\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1
	45 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"name\": \"Max Trans Unit\", \"select\": \"ip6Address\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

	# FQDN
	27 * * * * root /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\""  >> /var/log/traffic_ops/extensionCheck.log 2>&1
	27 * * * * root /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\", \"name\": \"DNS Lookup\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

	# DSCP
	36 * * * * root /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\", \"cms_interface\": \"eth0\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
	36 * * * * root /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\", \"name\": \"Delivery Service\", \"cms_interface\": \"eth0\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

	# RTR
	10 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
	10 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\", \"name\": \"Content Router Check\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

	# CHR
	*/15 * * * * root /opt/traffic_ops/app/bin/checks/ToCHRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"CHR\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1

	# CDU
	20 * * * * root /opt/traffic_ops/app/bin/checks/ToCDUCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"CDU\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1

	# ORT
	40 * * * * ssh_key_edge_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\"}"  >> /var/log/traffic_ops/extensionCheck.log 2>&1
	40 * * * * ssh_key_edge_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\", \"name\": \"Operational Readiness Test\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

