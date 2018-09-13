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

*****************************
Traffic Portal Administration
*****************************
Traffic Portal is only supported on CentOS Linux distributions - version 6.7 or higher (including 7.x). It runs on NodeJS and requires version 6.0 or higher.


Installing Traffic Portal
=========================

#. Download the Traffic Portal RPM from `Apache Jenkins <https://builds.apache.org/job/trafficcontrol-master-build/>`_ or build the Traffic Portal RPM from source (e.g. by running ``./pkg -v traffic_portal_build`` as the root user or with ``sudo`` from the top level of the source repository).
#. Copy the Traffic Portal RPM to your server
#. Install NodeJS. This can be done by building it from source, installing with ``yum install nodejs`` if it happens to be in your available repositories (at version 6.0+), or using the NodeSource setup script like so:

	.. code-block:: bash

		curl --silent --location https://rpm.nodesource.com/setup_6.x | sudo bash -

#. Install the Traffic Portal RPM e.g. by running ``yum install path/to/traffic_portal.rpm`` as the root user or with ``sudo``.


Configuring Traffic Portal
==========================

- update /etc/traffic_portal/conf/config.js (if upgrade, reconcile config.js with config.js.rpmnew and then delete config.js.rpmnew)
- update /opt/traffic_portal/public/traffic_portal_properties.json (if upgrade, reconcile traffic_portal_properties.json with traffic_portal_properties.json.rpmnew and then delete traffic_portal_properties.json.rpmnew)
- [OPTIONAL] update /opt/traffic_portal/public/resources/assets/css/custom.css (to customize traffic portal skin)


Starting Traffic Portal
=======================

The Traffic Portal RPM comes with a systemd unit file, so under normal circumstances all that is necessary is to run

	.. code-block:: bash

		systemctl start traffic_portal

as the root user or with ``sudo``.

Stopping Traffic Portal
=======================

If Traffic Portal was started using ``systemctl``, simply run

	.. code-block:: bash

		systemctl stop traffic_portal

as the root user, or with ``sudo``.
