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

**************
Traffic Portal
**************
Traffic Portal is only supported on CentOS Linux distributions version 7.x and 8.x. It runs on `NodeJS <https://nodejs.org/>`_ and requires version 16 or higher.

Installing
==========
#. Build the Traffic Portal RPM using the instructions in :ref:`dev-building`.
#. Copy the Traffic Portal RPM to your server
#. Install NodeJS. This can be done by building it from source, installing with :manpage:`yum(8)` if it happens to be in your available repositories (at version 16+), or using the NodeSource setup script.

	.. code-block:: bash
		:caption: Installing NodeJS using the NodeSource Setup Script

		curl --silent --location https://rpm.nodesource.com/setup_16.x | sudo bash -

#. Install the Traffic Portal RPM with :manpage:`yum(8)` or :manpage:`rpm(8)` e.g. by running ``yum install path/to/traffic_portal.rpm`` as the root user or with :manpage:`sudo(8)`.


Configuring
===========
Traffic Portal is primarily configured through three different files that affect different parts of its behavior. Those files are detailed in this section.

``/etc/traffic_portal/conf/config.js``
--------------------------------------
This file controls the behavior of the Traffic Portal server. It is a JavaScript source file which **MUST** set ``module.exports`` to the value of a configuration object. Configuration objects have the following allowed properties.

.. tip:: If Traffic Portal is being upgraded, reconcile ``config.js`` with ``config.js.rpmnew`` and then delete ``config.js.rpmnew``.

:api: The properties of this object control Traffic Portal's interactions with the :ref:`to-api`.

	:base_url: The **full** URL that points to the root of the :ref:`to-api`, e.g. ``https://trafficops.infra.ciab.test:443/api/`` **not** ``https://trafficops.infra.ciab.test:443`` or ``https://trafficops.infra.ciab.test:443/api`` (trailing ``/`` required).

:files: The properties of this object describe system paths important to Traffic Portal (other than SSL-related files, which are kept in the ``ssl`` object).

	:static: The directory where the built Traffic Portal front-end files are kept. In most cases, changing this from the default is unnecessary/discouraged.

:log: The properties of this object describe the logging behavior of Traffic Portal. If this property is missing, or is ``null`` (or other "falsey" values), then logging is done simply on STDOUT.

	:stream: Defines a file location for Traffic Portal's access logs. If this property is missing, or is ``null`` (or other "falsey" values), then logs will go to STDOUT instead.

:port:                A port number that sets the port on which Traffic Portal will listen for insecure (HTTP) connections. If ``useSSL`` is ``false``, requests made to this port will be redirected to use HTTPS on the ``sslPort`` instead of just serving content insecurely.
:reject_unauthorized: A boolean that defines whether or not Traffic Portal will reject SSL certificates that fail to validate as trusted.

	.. caution:: Setting this to ``false`` exposes Traffic Portal to security vulnerabilities such as man-in-the-middle attacks, and should *never* be done in a production setting.

:ssl: This object has properties that set the locations of SSL keys and certificates. Has no effect if ``useSSL`` is ``false``.

	:key:  The file location of the SSL certificate private key.
	:cert: The file location of the x509 SSL certificate that Traffic Portal will use.
	:ca:   The file locations of the full certificate chain for the certificate authority that signed the SSL key (in order).

:sslPort: A port number that sets the port on which Traffic Portal will listen for secure (HTTPS) connections. Has no effect if ``useSSL`` is ``false``.
:useSSL: A boolean that defines whether or not the Traffic Portal instance will offer secure (HTTPS) connections.

	.. caution:: Setting this to ``false`` can expose sensitive data such as authentication credentials. Do not *ever* do that in a production setting.

:timeout: This property defines the maximum time for which Traffic Portal will process requests before sending a timeout response. It can be either a number of milliseconds, or a duration string accepted by `the ms library <https://www.npmjs.com/package/ms#readme>`_.

	.. warning:: Slow requests will continue to use CPU and memory, even though a response has already been sent to the client.

``/opt/traffic_portal/public/traffic_portal_properties.json``
-------------------------------------------------------------
- update :file:`/opt/traffic_portal/public/traffic_portal_properties.json` (if Traffic Portal is being upgraded, reconcile :file:`traffic_portal_properties.json` with :file:`traffic_portal_properties.json.rpmnew` and then delete :file:`traffic_portal_properties.json.rpmnew`)

``/opt/traffic_portal/public/resources/assets/css/custom.css``
--------------------------------------------------------------
This :abbr:`CSS (Cascading Style Sheets)` file is provided for users to insert CSS to override the default styling of Traffic Portal.

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

Using Traffic Portal
====================
.. toctree::

	./usingtrafficportal.rst
