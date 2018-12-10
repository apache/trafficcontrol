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

.. _ciab:

************
CDN in a Box
************
"CDN in a Box" is a name given to the time-honored tradition of new Traffic Control developers/potential users attempting to set up their own miniature CDN to just see how it all fits together. Historically, this has been a nightmare of digging through leftover ``virsh`` scripts and manually configuring pretty hefty networking changes (don't even get me started on DNS) and just generally having a bad time. For a few years now, different people had made it to various stages of merging the project into Docker for ease of networking, but certain constraints hampered progress - until now. The project has finally reached a working state, and now getting a mock/test CDN running can be a very simple task (albeit rather time-consuming).

Getting Started
===============
Because it runs in Docker, the only true prerequisites are:

* Docker version >= 17.05.0-ce
* Docker Compose\ [1]_ version >= 1.9.0

Building
--------
The CDN in a Box directory is found within the Traffic Control repository at ``infrastructure/cdn-in-a-box/``. CDN in a Box relies on the presence of pre-built RPM files for the following Traffic Control components:

* Traffic Monitor - at ``infrastructure/cdn-in-a-box/traffic_monitor/traffic_monitor.rpm``
* Traffic Ops - at ``infrastructure/cdn-in-a-box/traffic_ops/traffic_ops.rpm``
* Traffic Portal - at ``infrastructure/cdn-in-a-box/traffic_portal/traffic_portal.rpm``
* Traffic Router - at ``infrastructure/cdn-in-a-box/traffic_router/traffic_router.rpm`` - also requires an Apache Tomcat RPM at ``infrastructure/cdn-in-a-box/traffic_router/tomcat.rpm``

.. note:: These can also be specified via the ``RPM`` variable to a direct Docker build of the component - with the exception of Traffic Router, which instead accepts ``JDK8_RPM`` to specify a Java Development Kit RPM,  ``TRAFFIC_ROUTER_RPM`` to specify a Traffic Router RPM, and  ``TOMCAT_RPM`` to specify an Apache Tomcat RPM.

These can all be supplied manually via the steps in :ref:`dev-building` (for Traffic Control component RPMs) or via some external source. Alternatively, the ``infrastructure/cdn-in-a-box/Makefile`` file contains recipes to build all of these - simply run ``make``\ [2]_ from the ``infrastructure/cdn-in-a-box/`` directory.
Once all RPM dependencies have been satisfied, run ``docker-compose build`` from the ``infrastructure/cdn-in-a-box/`` directory to construct the images needed to run CDN in a Box.

Usage
-----
In a typical scenario, if the steps in `Building`_ have been followed, all that's required to start the CDN in a Box is to run ``docker-compose up`` - optionally with ``-d``  to run without binding to the terminal - from the ``infrastructure/cdn-in-a-box/`` directory. This will start up the entire stack and should take care of any needed initial configuration. The services within the containers are exposed locally to the host on specific ports. These are configured within the ``infrastructure/cdn-in-a-box/docker-compose.yml`` file, but the default ports are shown in :ref:`ciab-service-info`. Some services have credentials associated, which are totally configurable in `variables.env`_.

.. _ciab-service-info:
.. table:: Service Info

	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Service                         | Ports exposed and their usage                                | Username                              | Password                                  |
	+=================================+==============================================================+=======================================+===========================================+
	| DNS                             | DNS name resolution on 9353                                  | N/A                                   | N/A                                       |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Edge Tier Cache                 | Apache Trafficserver HTTP caching reverse proxy on port 9000 | N/A                                   | N/A                                       |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Mid Tier Cache                  | Apache Trafficserver HTTP caching forward proxy on port 9100 | N/A                                   | N/A                                       |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Mock Origin Server              | Example web page served on port 9200                         | N/A                                   | N/A                                       |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Monitor                 | Web interface and API on port 80                             | N/A                                   | N/A                                       |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Ops                     | Main API endpoints on port 6443, with a direct route to the  | ``TO_ADMIN_USER`` in `variables.env`_ | ``TO_ADMIN_PASSWORD`` in `variables.env`_ |
	|                                 | Perl API on port 60443\ [3]_                                 |                                       |                                           |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Ops PostgresQL Database | PostgresQL connections accepted on port 5432 (database name: | ``DB_USER`` in `variables.env`_       | ``DB_USER_PASS`` in `variables.env`_      |
	|                                 | ``DB_NAME`` in `variables.env`_)                             |                                       |                                           |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Portal                  | Web interface on 443 (Javascript required)                   | ``TO_ADMIN_USER`` in `variables.env`_ | ``TO_ADMIN_PASSWORD`` in `variables.env`_ |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Router                  | Web interfaces on ports 3080 (HTTP) and 3443 (HTTPS), with a | N/A                                   | N/A                                       |
	|                                 | DNS service on 53 and an API on 3333                         |                                       |                                           |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+
	| Traffic Vault                   | Riak key-value store on port 8010                            | ``TV_ADMIN_USER`` in `variables.env`_ | ``TV_ADMIN_PASSWORD`` in `variables.env`_ |
	+---------------------------------+--------------------------------------------------------------+---------------------------------------+-------------------------------------------+

.. seealso:: :ref:`tr-api` and :ref:`tm-api`

While the components may be interacted with by the host using these ports, the true operation of the CDN can only truly be seen from within the Docker network. To see the CDN in action, connect to a container within the CDN in a Box project and use cURL to request the URL ``http://video.demo1.mycdn.ciab.test`` which will be resolved by the DNS container to the IP of the Traffic Router, which will provide a ``302 FOUND`` response pointing to the Edge-Tier cache. A typical choice for this is the "enroller" service, which has a very nuanced purpose not discussed here but already has the ``curl`` command line tool installed. For a more user-friendly interface into the CDN network, see `VNC Server`_.

.. code-block:: shell
	:caption: Example Command to See the CDN in Action

	sudo docker-compose exec enroller /usr/bin/curl -L "http://video.demo1.mycdn.ciab.test"

When the CDN is to be shut down, it is often best to do so using ``sudo docker-compose down -v`` due to the use of shared volumes in the system which might interfere with a proper initialization upon the next run.

variables.env
"""""""""""""
.. include:: ../../../../infrastructure/cdn-in-a-box/variables.env
	:code: shell
	:start-line: 16
	:tab-width: 4

.. note:: While these port settings can be changed without hampering the function of the CDN in a Box system, note that changing a port without also changing the matching port-mapping in ``infrastructure/cdn-in-a-box/docker-compose.yml`` for the affected service *will* make it unreachable from the host.

.. [1] It is perfectly possible to build and run all containers without Docker Compose, but it's not recommended and not covered in this guide.
.. [2] Consider ``make -j`` to build quickly, if your computer can handle multiple builds at once.
.. [3] Please do NOT use the Perl endpoints directly. The CDN will only work properly if everything hits the Go API, which will proxy to the Perl endpoints as needed.

X.509 SSL/TLS Certificates
==========================
All components in Apache Traffic Control utilize SSL/TLS secure communications by default. For SSL/TLS connections to properly validate within the "CDN in a Box" container network a shared self-signed X.509 Certificate Authority (CA) is generated at the first initial startup. Additional self-signed wildcard certificates are generated for each container service and all delivery services of the CDN. All certificates and keys are stored in the ``ca`` host volume which is located at ``infrastruture/cdn-in-a-box/traffic_ops/ca`` [4]_.

.. _ciab-x509-certificate-list:
.. table:: Self-Signed X.509 Certificate List

	+---------------------------+-----------------------------------+------------------------------+
	| Filename                  | Description                       | X.509 CN/SAN                 |
	+===========================+===================================+==============================+
	| CIAB-CA.crt               | Shared CA Certificate             | N/A                          |
	+---------------------------+-----------------------------------+------------------------------+
	| infra.ciab.test.crt       | Infrastruture Certificate         | \*.infra.ciab.test           |
	+---------------------------+-----------------------------------+------------------------------+
	| demo1.mycdn.ciab.test.crt | Demo1 Delivery Service Certificate| \*.demo1.mycdn.ciab.test     |
	+---------------------------+-----------------------------------+------------------------------+
	| demo2.mycdn.ciab.test.crt | Demo2 Delivery Service Certificate| \*.demo2.mycdn.ciab.test     |
	+---------------------------+-----------------------------------+------------------------------+
	| demo3.mycdn.ciab.test.crt | Demo3 Delivery Service Certificate| \*.demo3.mycdn.ciab.test     |
	+---------------------------+-----------------------------------+------------------------------+

.. [4] The ``ca`` volume is not purged with normal ``docker volume`` commands. This feature is by design to allow the existing shared SSL certificate to be trusted at the system level across restarts. To re-generate all SSL certificates and keys, remove the ``infrastructure/cdn-in-a-box/traffic_ops/ca`` directory before startup.

Trusting the CA
---------------
For developer and testing use-cases, it may be necessary to have full x509 CA validation by HTTPS clients [5]_. For x509 validation to work properly, the self-signed x509 CA certificate must be trusted either at the system level or by the client application itself. Procedures to import and trust the CA x.509 certificate are outlined below [6]_.

Importing the CA Certificate on OSX
-----------------------------------
#. Copy the CIAB root CA certificate from ``infrastructure/cdn-in-a-box/traffic_ops/ca/CIAB-CA.crt`` to the Mac.
#. Import the CIAB root CA certificate on the Mac.
#. Double-click the CIAB root CA certificate to open it in Keychain Access.
#. The CIAB root CA certificate appears in login.
#. Copy the CIAB root CA certificate to System.
#. Open the CIAB root CA certificate, expand Trust, select Use System Defaults, and save your changes.
#. Reopen the CIAB root CA certificate, expand Trust, select Always Trust, and save your changes.
#. Delete the CIAB root CA certificate from login.
#. Restart all HTTPS clients (browsers, etc).

Importing the CA certificate on Windows
---------------------------------------
#. Copy the CIAB root CA certificate from ``infrastructure/cdn-in-a-box/traffic_ops/ca/CIAB-CA.crt`` to Windows filesystem.
#. As Administrator, start the Microsoft Management Console.
#. Add the Certificates snap-in for the computer account and manage certificates for the local computer.
#. Import the CIAB root CA certificate into Trusted Root Certification Authorities > Certificates.
#. Restart all HTTPS clients (browsers, etc).

Importing the CA certificate on Linux/Centos7
---------------------------------------------
#. Copy the CIAB root CA certificate from ``infrastructure/cdn-in-a-box/traffic_ops/ca/CIAB-CA.crt`` to path ``/etc/pki/ca-trust/source/anchors``.
#. Run ``update-ca-trust-extract`` as the root user.
#. Restart all HTTPS clients (browsers, etc).

Importing the CA certificate on Linux/Ubuntu
--------------------------------------------
#. Copy the CIAB root CA certificate from ``infrastructure/cdn-in-a-box/traffic_ops/ca/CIAB-CA.crt`` to path ``/usr/local/share/ca-certificates``.
#. Run ``update-ca-certificates`` as the root user.
#. Restart all HTTPS clients (browsers, etc).

.. [5] All containers within CDN-in-a-Box start up with the self-signed CA already trusted.
.. [6] HTTP Client applications such as Google Chrome, Firefox, curl, and wget can also be individually configured to trust the CA certificate. Each application procedure can be found quickly online via Google.

Advanced Usage
==============
This section will be amended as functionality is added to the CDN in a Box project.

The Enroller
------------
The "enroller" provides an efficient way for Traffic Ops to be populated with data as CDN in a Box starts up. It connects to Traffic Ops as the admin user and processes files places in the docker volume shared between the containers. The enroller watches each directory within the ``/shared/enroller`` directory for new ``.json`` files to be created there. These files must follow the format outlined in the API guide for the ``POST`` method for each data type,  (e.g. for a ``tenant``, follow the guidelines for ``POST api/1.4/regions``). Of note,  the ``enroller`` does not require fields that reference database ids for other objects within the database.

The enroller runs within CDN in a Box using the ``-dir <dir>`` switch which provides the above behavior. It can also be run using the ``-http :<port>`` switch to instead have it listen on the indicated port. In this case, it accepts only POST requests with the JSON provided using the POST JSON method, e.g. ``curl -X POST https://enroller/api/1.4/regions -d @newregion.json``.  CDN in a Box does not currently use this method, but may be modified in the future to avoid using the shared volume approach.

Mock Origin Service
-------------------
The default "origin" service container provides a basic static file HTTP server as the central respository for content. Additional files can be added to the origin root content directory located at ``infrastructure/cdn-in-a-box/origin/content``. To request content directly from the origin directly and bypass the CDN:

* Origin Service URL: http://origin.infra.ciab.test/index.html
* Docker Host: http://localhost:9200/index.html

.. _ciab-optional-containers:

Optional Containers
===================

All optional containers that are not part of the core CDN-in-a-Box stack are located in the ``infrastructure/cdn-in-a-box/optional`` directory.

.. code-block:: shell

	infrastructure/cdn-in-a-box/optional/docker-compose.$NAME.yml
	infrastructure/cdn-in-a-box/optional/$NAME/Dockerfile

Multiple optional containers may be combined by using a shell alias:

.. code-block:: shell

	# From the infrastructure/cdn-in-a-box directory
	alias mydc="docker-compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.$NAME1.yml -f  $PWD/optional/docker-compose.$NAME2.yml"
	docker volume prune -f
	mydc build
	mydc up

VNC Server
----------
The TightVNC optional container provides a basic lightweight window manager (fluxbox), Firefox browser, xterm, and a few other utilities within the CDN-In-A-Box tcnet bridge network. This can be very helpful for quick demonstrations of CDN-in-a-Box that require the use of real container network FQDNs and full X.509 validation.

#. Download and install a VNC client. TightVNC client is preferred as it supports window resizing, host-to-vnc copy/pasting, and optimized frame buffer compression.
#. Set your VNC console password by adding the ``VNC_PASSWD`` environment variable to ``infrastructure/cdn-in-a-box/varibles.env``. The password needs to be at least six characters long. The default password is randomized for security.
#. Start up CDN-in-a-Box stack. It is recommended that this be done using a custom bash alias

	.. code-block:: shell
		:caption: CIAB Startup Using Bash Alias

		# From infrastructure/cdn-in-a-box
		alias mydc="docker-compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.vnc.yml"
		docker volume prune -f
		mydc build
		mydc kill
		mydc rm -fv
		mydc up

#. Connect with a VNC client to localhost port 9080.
#. When Traffic Portal becomes available, the Firefox within the VNC instance will subsequently be started.
#. An xterm with bash shell is also automatically spawned and minimized for convenience.

Socks Proxy
-----------
Dante's socks proxy is an optional container that can be used to provide browsers and other clients the ability to resolve DNS queries and network connectivity directly on the tcnet bridged interface. This is very helpful when running the CDN-In-A-Box stack on OSX/Windows docker host that lacks network bridge/IP-forward support. Below is the basic procedure to enable the Socks Proxy support for CDN-in-a-Box:

#. Start the CDN-in-a-Box stack at least once so that the x.509 self-signed certificate authority (CA) is created.
#. On the host, import and Trust the CA for your target OS. See `Trusting the CA`_
#. On the host, using either Firefox or Chrome, download the FoxyProxy Standard browser plugin which enables dynamic proxy support via URL regular expression
#. Once FoxyProxy is installed, click the Fox icon on the upper right hand of the browser window, select 'Options'
#. Once in Options Dialog, Click 'Add New Proxy' and navigate to the General tab:
#. Fill in the General tab according to the table

	.. table:: General Tab Values

		+------------+---------+
		| Name       |   Value |
		+============+=========+
		| Proxy Name |    CIAB |
		+------------+---------+
		| Color      |   Green |
		+------------+---------+

#. Fill in the Proxy Details tab according to the table

	.. table:: Proxy Details Tab Values

		+----------------------------+-----------+
		| Name                       |     Value |
		+============================+===========+
		| Manual Proxy Configuration |      CIAB |
		+----------------------------+-----------+
		| Host or IP Address         | localhost |
		+----------------------------+-----------+
		| Port                       |      9080 |
		+----------------------------+-----------+
		| Socks Proxy                |   checked |
		+----------------------------+-----------+
		| Socks V5                   |  selected |
		+----------------------------+-----------+

#. Go to URL Patterns tab, click Add New Pattern, and fill out form according to

	.. table:: URL Patters Tab Values

		+--------------+--------------+
		| Name         |        Value |
		+==============+==============+
		| Pattern Name | CIAB Pattern |
		+--------------+--------------+
		| URL Pattern  |   \*.test/\* |
		+--------------+--------------+
		| Whitelist    |     selected |
		+--------------+--------------+
		| Wildcards    |     selected |
		+--------------+--------------+

#. Enable dynamic 'pre-defined and patterns' mode by clicking the fox icon in the upper right of the browser. This mode only forwards URLs that match the wildcard ``\*.test/\*`` to the Socks V5 proxy.

10. On the docker host start up CDN-in-a-Box stack. It is recommended that this be done using a custom bash alias

	.. code-block:: shell
		:caption: CIAB Startup Using Bash Alias

		# From infrastructure/cdn-in-a-box
		alias mydc="docker-compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.socksproxy.yml"
		docker volume prune -f
		mydc build
		mydc kill
		mydc rm -fv
		mydc up

#. Once the CDN-in-a-box stack has started, use the aforementioned browser to access traffic portal via the socks proxy on the docker host.

.. seealso:: `The official Docker Compose documentation CLI reference <https://docs.docker.com/compose/reference/>`_ for complete instructions on how to pass service definition files to the ``docker-compose`` executable.
