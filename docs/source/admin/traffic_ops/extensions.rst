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

Traffic Ops supports two types of extensions. `Check Extensions`_ are analytics scripts that collect and display information as columns in the table under :menuselection:`Monitor --> Cache Checks` in Traffic Portal. `Data Source Extensions`_ provide ways to add data to the graph views and usage APIs.

.. |checkmark| image:: images/good.png
.. |X| image:: images/bad.png

.. _to-check-ext:

Check Extensions (go version)
================
The following applies to the these checks, which have been ported from perl to go:
:abbr:`CDU (Cache Disk Usage)`
:abbr:`CHR (Cache Hit Ratio)`
:abbr:`DSCP (Differential Services CodePoint)`
:abbr:`MTU (Maximum Transmission Unit)`
:abbr:`10G (IPv4 ping check)`
:abbr:`10G6 (IPv6 ping check)`
:abbr:`ILO (Integrated Lights-Out)`
:abbr:`FQDN (Fully Qualified Domain Name)`

Check Extensions are scripts that, after registering with Traffic Ops, have a column reserved in the :menuselection:`Monitor --> Cache Checks` view and usually run periodically using :manpage:`cron(8)`. Each extension is a separate go program location in :file:`{$TO_HOME}/bin/checks/` on the Traffic Ops server. The currently registered extensions can be listed by running ``/opt/traffic_ops/app/bin/extensions -a``. Some extensions automatically registered with the Traffic Ops database (``to_extension`` table) at install time (see :atc-file:`traffic_ops/app/db/seeds.sql`). However, :manpage:`cron(8)` must still be configured to run these checks periodically. The extensions are called like so:

.. code-block:: shell
	:caption: Example Check Extension Call (go version)

  RLOG_LOG_LEVEL=<log level> $TO_HOME/bin/checks/<name> -conf <conf_path> -name <check_name> -syslog
	$TO_HOME/bin/checks/<name>  -c "{\"base_url\": \",https://\"<traffic_ops_ip>\", \"check_name\": \"<check_name>\"}" -l <log level>

:log_level: A whole number between 1 and 4 (inclusive), with 4 being the most verbose. Implementation of this field is optional
:name: The basename of the extension executable
:conf_path: The path to check-config.json e.g. '/opt/traffic_ops/conf/check-conf.json'
:check_name: The name of the check e.g. ``CDU``, ``CHR``, ``DSCP``, ``MTU``, etc...

Note that the DSCP check requires the '-iface <interface_name>' argument in order to tell the program which interface to use to capture response packets on for analysis.

Note also the following optional arguments:
'-syslog' boolean flag; result data is sent to STDOUT, but can optionally be sent to syslog using this option.
'-host <regex>' regex filter for narrowing down the set of hosts to check
'-exclude <regex>' regex filter for specifying host patters to exclude (can be combined with -host filter)
'-f' force fail boolean flag; mainly useful for testing purposes
'-reset' boolean flag; reset check status to "blank" state in Traffic Ops DB


TOAPI authentication information is set in the config file. Currently, is is mandatory to use the username "extension", which should be granted the read-only role in Traffic Ops. This user is special, however, in that it is granted special privileges for updating the extension check tables ONLY in the Traffic Ops DB. Each program receives a list of servers from Traffic Ops API, subject to the optional filters noted above, and loops through the list performing the specified check for each server. The results are then compiled into a numeric or boolean result and the script submits a ``POST`` request containing the result back to Traffic Ops using :ref:`to-api-servercheck`. A check extension can have a column of |checkmark|'s and |X|'s (CHECK_EXTENSION_BOOL) or a column that shows a number (CHECK_EXTENSION_NUM).

Check Extensions (legacy perl version)
================
Check Extensions are scripts that, after registering with Traffic Ops, have a column reserved in the :menuselection:`Monitor --> Cache Checks` view and usually run periodically using :manpage:`cron(8)`. Each extension is a separate executable located in :file:`{$TO_HOME}/bin/checks/` on the Traffic Ops server (though all of the default extensions are written in Perl, this is in *no way* a requirement; they can be any valid executable). The currently registered extensions can be listed by running ``/opt/traffic_ops/app/bin/extensions -a``. Some extensions automatically registered with the Traffic Ops database (``to_extension`` table) at install time (see :atc-file:`traffic_ops/app/db/seeds.sql`). However, :manpage:`cron(8)` must still be configured to run these checks periodically. The extensions are called like so:

.. code-block:: shell
	:caption: Example Check Extension Call (legacy perl version)

	$TO_HOME/bin/checks/<name>  -c "{\"base_url\": \",https://\"<traffic_ops_ip>\", \"check_name\": \"<check_name>\"}" -l <log level>

:name: The basename of the extension executable
:traffic_ops_ip: The IP address or :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Ops server
:check_name: The name of the check e.g. ``CDU``, ``CHR``, ``DSCP``, ``MTU``, etc...
:log_level: A whole number between 1 and 4 (inclusive), with 4 being the most verbose. Implementation of this field is optional

It is the responsibility of the check extension script to iterate over the servers it wants to check and post the results. An example script might proceed by logging into the Traffic Ops server using the HTTPS ``base_url`` provided on the command line. The script is hard-coded with an authentication token that is also provisioned in the Traffic Ops User database. This token allows the script to obtain a cookie used in later communications with the Traffic Ops API. The script then obtains a list of all :term:`cache server`\ s to be polled by accessing :ref:`to-api-servers`. This list is then iterated, running a command to gather the stats from each server. For some extensions, an HTTP ``GET`` request might be made to the :abbr:`ATS (Apache Traffic Server)` ``astats`` plugin, while for others the server might be pinged, or a command might run over :manpage:`ssh(1)`. The results are then compiled into a numeric or boolean result and the script submits a ``POST`` request containing the result back to Traffic Ops using :ref:`to-api-servercheck`. A check extension can have a column of |checkmark|'s and |X|'s (CHECK_EXTENSION_BOOL) or a column that shows a number (CHECK_EXTENSION_NUM).

Check Extensions Installed by Default
-------------------------------------
:abbr:`CDU (Cache Disk Usage)`
	This check shows how much of the available total cache disk is in use. A "warm" cache should show 100.00.

:abbr:`CHR (Cache Hit Ratio)`
	The cache hit ratio for the cache in the last 15 minutes (the interval is determined by the ``cron`` entry).

:abbr:`DSCP (Differential Services CodePoint)`
	Checks if the returning traffic from the cache has the correct :abbr:`DSCP (Differential Services CodePoint Check)` value as assigned in the :term:`Delivery Service`. (Some routers will overwrite :abbr:`DSCP (Differential Services CodePoint)`)

:abbr:`MTU (Maximum Transmission Unit)`
	Checks if the Traffic Ops host (if that is the one running the check) can send and receive 8192B packets to the ``ip_address`` of the server in the server table.

:abbr:`ORT (Operational Readiness Test)`
	The ORT column shows how many changes the :term:`ORT` script would apply if it was run. The number in this column should be 0 for caches that do not have updates pending.

10G
	Is the ``ip_address`` (the main IPv4 address) from the server table ping-able?
:abbr:`ILO (Integrated Lights-Out)`
	Is the ``ilo_ip_address`` (the lights-out-management IPv4 address) from the server table ping-able?
10G6
	Is the ``ip6_address`` (the main IPv6 address) from the server table ping-able?
:abbr:`FQDN (Fully Qualified Domain Name)`
	Is the :abbr:`FQDN (Fully Qualified Domain Name)` (the concatenation of ``host_name`` and ``.`` and ``domain_name`` from the server table) ping-able?

:abbr:`RTR (Responds to Traffic Router)`
	Checks the state of each cache as perceived by all Traffic Monitors (via Traffic Router). This extension asks each Traffic Router for the state of the cache. A check failure is indicated if one or more monitors report an error for a cache. A cache is only marked as good if all reports are positive. (This is a pessimistic approach, opposite of how TM marks a cache as up, "the optimistic approach")

.. _to-datasource-ext:

Data Source Extensions
======================
Data Source Extensions work in much the same way as `Check Extensions`_, but are implemented differently. Rather than being a totally external executable, a Data Source Extension *must* be written in Perl 5, as they are injected via manipulation of the ``$PERL5LIB`` environment variable. These extensions are not very well-documented (as you may be able to tell), and support for extending them may be phased out in future releases.

Example Cron File (go version)
=======================================
The :manpage:`cron(8)` file should be edited by running  :manpage:`crontab(1)` as the ``traffops`` user, or with :manpage:`sudo(8)`. You may need to adjust the path to your ``$TO_HOME`` to match your system.

.. code-block:: shell
	:caption: Example Cron File (go version)

  # ILO Ping
  30 */6 * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToPingCheck -name ILO -syslog >/var/log/traffic_ops/check_ext_ilo.log 2>&1

  # IPv4 Ping
  0 */6 * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToPingCheck -name 10G -syslog >/var/log/traffic_ops/check_ext_10g.log 2>&1

  # IPv6 Ping
  20 */6 * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToPingCheck -name 10G6 -syslog >/var/log/traffic_ops/check_ext_10g6.log 2>&1

  # DSCP check (with exclude filter)
  0 3 * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToDSCPCheck -name DSCP -exclude '^(devcache-|qacache-).*' -iface ens3 -syslog >/var/log/traffic_ops/check_ext_dscp.log 2>&1

  # FQDN check (with optional PTR validation and exclude filter)
  33 * * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToFQDNCheck -name FQDN -ptr -exclude '^(devcache-|qacache-).*' -syslog >/var/log/traffic_ops/check_ext_fqdn.log 2>&1

  # Cache Disk Usage (CDU) check (with host filter)
  */20 * * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToATSCheck -name CDU -host '^prodcache-.*' -syslog >/var/log/traffic_ops/check_ext_cdu.log 2>&1

  # Cache Hit Ration (CHR) check (with host filter)
  */15 * * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToATSCheck -name CHR -host '^prodcache-.*' -syslog >/var/log/traffic_ops/check_ext_chr.log 2>&1

  # MTU Ping
  50 */6 * * * root RLOG_LOG_LEVEL=INFO /opt/to_checks/bin/ToPingCheck -name MTU -syslog >/var/log/traffic_ops/check_ext_mtu.log 2>&1


Example Cron File (legacy perl version)
=======================================
The :manpage:`cron(8)` file should be edited by running  :manpage:`crontab(1)` as the ``traffops`` user, or with :manpage:`sudo(8)`. You may need to adjust the path to your ``$TO_HOME`` to match your system.

.. code-block:: shell
	:caption: Example Cron File (legacy perl version)

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
	36 * * * * root /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\", \"name\": \:term:`Delivery Service`\", \"cms_interface\": \"eth0\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

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
