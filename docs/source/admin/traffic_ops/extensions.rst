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

.. |checkmark| image:: images/good.png
.. |X| image:: images/bad.png

.. _to-check-ext:

Check Extensions
================
Check Extensions are scripts that, after registering with Traffic Ops, have a column reserved in the "Monitor"->"Cache Checks" view ("Health"->"Server Checks" in the legacy UI) and usually run periodically using ``cron``. Each extension is a separate executable located in ``$TO_HOME/bin/checks/`` on the Traffic Ops server (though all of the default extensions are written in Perl, this is in *no way* a requirement; they can be any valid executable). The currently registered extensions can be listed by running ``/opt/traffic_ops/app/bin/extensions -a``. Some extensions automatically registered with the Traffic Ops database (``to_extension`` table) at install time (see ``traffic_ops/app/db/seeds.sql``). However, ``cron`` must still be configured to run these checks periodically. The extensions are called like so:

.. code-block:: shell
	:caption: Example Check Extension Call

	$TO_HOME/bin/checks/<name>  -c "{\"base_url\": \",https://\"<traffic_ops_ip>\", \"check_name\": \"<check_name>\"}" -l <log level>

:name: The basename of the extension executable
:traffic_ops_ip: The IP address or Fully Qualified Domain Name (FQDN) of the Traffic Ops server
:check_name: The name of the check e.g. ``CDU``, ``CHR``, ``DSCP``, ``MTU``, etc...
:log_level: A whole number between 1 and 4 (inclusive), with 4 being the most verbose. Implementation of this field is optional

It is the responsibility of the check extension script to iterate over the servers it wants to check and post the results. An example script might proceed by logging into the Traffic Ops server using the HTTPS ``base_url`` provided on the command line. The script is hard-coded with an authentication token that is also provisioned in the Traffic Ops User database. This token allows the script to obtain a cookie used in later communications with the Traffic Ops API. The script then obtains a list of all caches to be polled by accessing the Traffic Ops ``/api/1.3/servers`` REST endpoint. This list is then iterated, running a command to gather the stats from that cache. For some extensions, an HTTP GET request might be made to the Apache Traffic Server (ATS) ``astats`` plugin, while for others the cache might be pinged, or a command might run over SSH. The results are then compiled into a numeric or boolean result and the script POSTs the result back to the Traffic Ops using the ``/api/1.3/servercheck/`` endpoint. A check extension can have a column of |checkmark|'s and |X|'s (CHECK_EXTENSION_BOOL) or a column that shows a number (CHECK_EXTENSION_NUM).A simple example of a check extension of type CHECK_EXTENSION_NUM that will show 99.33 for all servers of type EDGE is shown below: ::

	Script here.

Currently, the following Check Extensions are available and installed by default:

Cache Disk Usage Check (CDU)
	This check shows how much of the available total cache disk is in use. A "warm" cache should show 100.00.

Cache Hit Ratio Check (CHR)
	The cache hit ratio for the cache in the last 15 minutes (the interval is determined by the ``cron`` entry).

Differential Services CodePoint Check (DSCP)
	Checks if the returning traffic from the cache has the correct DSCP value as assigned in the delivery service. (Some routers will overwrite DSCP)

Maximum Transmission Unit Check (MTU)
	Checks if the Traffic Ops host (if that is the one running the check) can send and receive 8192B packets to the ``ip_address`` of the server in the server table.

Operational Readiness Check (ORT)
	See :ref:`traffic-ops-ort` for more information on the ORT script. The ORT column shows how many changes the ``traffic_ops_ort.pl`` script would apply if it was run. The number in this column should be 0 for caches that do not have updates pending.

Ping Check - 10G, ILO, 10G6, FQDN
	The ``bin/checks/ToPingCheck.pl`` script checks basic IP connectivity, and in the default setup it checks IP connectivity to the following:

	10G
		Is the ``ip_address`` (the main IPv4 address) from the server table ping-able?
	ILO
		Is the ``ilo_ip_address`` (the lights-out-management IPv4 address) from the server table ping-able?
	10G6
		Is the ``ip6_address`` (the main IPv6 address) from the server table ping-able?
	FQDN
		Is the Fully Qualified Domain name (the concatenation of ``host_name`` and ``.`` and ``domain_name`` from the server table) ping-able?

Traffic Router Check (RTR)
	Checks the state of each cache as perceived by all Traffic Monitors (via Traffic Router). This extension asks each Traffic Router for the state of the cache. A check failure is indicated if one or more monitors report an error for a cache. A cache is only marked as good if all reports are positive. (This is a pessimistic approach, opposite of how TM marks a cache as up, "the optimistic approach")

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

