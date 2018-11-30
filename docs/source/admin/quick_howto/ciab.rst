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

While the components may be interacted with by the host using these ports, the true operation of the CDN can only truly be seen from within the Docker network. To see the CDN in action, connect to a container within the CDN in a Box project and use cURL to request the URL ``http://video.demo1.mycdn.ciab.test`` which will be resolved by the DNS container to the IP of the Traffic Router, which will provide a ``302 FOUND`` response pointing to the Edge-Tier cache. A typical choice for this is the "enroller" service, which has a very nuanced purpose not discussed here but already has the ``curl`` command line tool installed. For a more user-friendly interface into the CDN network, see `The Test Client`_.

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

Advanced Usage
==============
This section will be amended as functionality is added to the CDN in a Box project.

The Enroller
------------
The "enroller" provides an efficient way for Traffic Ops to be populated with data as CDN in a Box starts up.  It connects to Traffic Ops as the admin user and processes files places in the docker volume shared between the containers.  The enroller watches each directory within the ``/shared/enroller`` directory for new ``.json`` files to be created there.  These files must follow the format outlined in the API guide for the ``POST`` method for each data type,  (e.g. for a ``tenant``, follow the guidelines for ``POST api/1.4/regions``).  Of note,  the ``enroller`` does not require fields that reference database ids for other objects within the database.

The enroller runs within CDN in a Box using the ``-dir <dir>`` switch which provides the above behavior.  It can also be run using the ``-http :<port>`` switch to instead have it listen on the indicated port.  In this case, it accepts only POST requests with the JSON provided using the POST JSON method, e.g. ``curl -X POST https://enroller/api/1.4/regions -d @newregion.json``.   CDN in a Box does not currently use this method, but may be modified in the future to avoid using the shared volume approach.


The Test Client
---------------
The "testclient" service is an optional extension to CDN in a Box which provides several more user-friendly interfaces to the CDN network. Specifically it contains:

* A small HTTP proxy on port 7070
* An SSH server listening on port 2200
* A VNC server on port 5900
* A Socks5 Proxy on port 9090

Using this, the CDN can be directly utilized using the network's internal name resolution and a web browser. To use the test client, pass ``-f docker-compose.testclient.yml`` either to the same ``docker-compose`` command, or to a separate ``docker-compose`` command run after the CDN in a Box has already started. If the former is done, note that it is also necessary to pass ``-f docker-compose.yml`` to properly build the entire system along with the test client.

.. seealso:: `The official Docker Compose documentation CLI reference <https://docs.docker.com/compose/reference/>`_ for complete instructions on how to pass service definition files to the ``docker-compose`` executable.
