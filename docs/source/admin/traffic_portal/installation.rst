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
Traffic Portal is only supported on CentOS Linux distributions version 7.x and 8.x. It runs on `NodeJS <https://nodejs.org/>`_ and requires version 12 or higher.


Installing Traffic Portal
=========================

#. Download the Traffic Portal RPM from `Apache Jenkins <https://builds.apache.org/job/trafficcontrol-master-build/>`_ or build the Traffic Portal RPM from source using the instructions in :ref:`dev-building`.
#. Copy the Traffic Portal RPM to your server
#. Install NodeJS. This can be done by building it from source, installing with :manpage:`yum(8)` if it happens to be in your available repositories (at version 12+), or using the NodeSource setup script.

	.. code-block:: bash
		:caption: Installing NodeJS using the NodeSource Setup Script

		curl --silent --location https://rpm.nodesource.com/setup_12.x | sudo bash -

#. Install the Traffic Portal RPM with :manpage:`yum(8)` or :manpage:`rpm(8)` e.g. by running ``yum install path/to/traffic_portal.rpm`` as the root user or with :manpage:`sudo(8)`.


Configuring Traffic Portal
==========================
- update :file:`/etc/traffic_portal/conf/config.js` (if Traffic Portal is being upgraded, reconcile :file:`config.js` with :file:`config.js.rpmnew` and then delete :file:`config.js.rpmnew`)
- update :file:`/opt/traffic_portal/public/traffic_portal_properties.json` (if Traffic Portal is being upgraded, reconcile :file:`traffic_portal_properties.json` with :file:`traffic_portal_properties.json.rpmnew` and then delete :file:`traffic_portal_properties.json.rpmnew`)
- Optional: update :file:`/opt/traffic_portal/public/resources/assets/css/custom.css` to customize Traffic Portal styling.


Configuring OAuth Through Traffic Portal
========================================
See :ref:`oauth_login`.


Starting Traffic Portal
=======================
The Traffic Portal RPM comes with a :manpage:`systemd(1)` unit file, so under normal circumstances Traffic Portal may be started with :manpage:`systemctl(1)`.

.. code-block:: bash
	:caption: Starting Traffic Portal

	systemctl start traffic_portal

Stopping Traffic Portal
=======================
The Traffic Portal RPM comes with a :manpage:`systemd(1)` unit file, so under normal circumstances Traffic Portal may be stopped with :manpage:`systemctl(1)`.

.. code-block:: bash
	:caption: Stopping Traffic Portal

	systemctl stop traffic_portal
