..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..   http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

************
Introduction
************
Traffic Control is a :abbr:`CDN (Content Delivery Network)` control plane. It is made up of a suite of applications which are used to configure, manage, and direct client traffic to a tiered system of HTTP caching proxy servers (herein referred to as :term:`cache servers`). In principle, a CDN may be implemented with any HTTP caching proxy. The caching software chosen for Traffic Control is `Apache Traffic Server <http://trafficserver.apache.org/>`_. Although the current release supports only :abbr:`ATS (Apache Traffic Server)` as a :term:`cache server` implementation, this may change with future releases.

Traffic Control was first developed at Comcast for internal use and released to Open Source in April of 2015. Traffic Control moved into the Apache Incubator in August of 2016.

Traffic Control implements the elements illustrated in green in the diagram below.


.. figure:: images/traffic.control.overview.*
	:align: center
	:width: 100%

	Apache Traffic Control Hierarchical Diagram


:ref:`to-overview`
	:dfn:`Traffic Ops` stores the configuration of :term:`cache servers` and CDN :term:`Delivery Services`. It also serves the :ref:`to-api` which can be used by tools, scripts, and programs to access and manipulate CDN data.

:ref:`tr-overview`
	:dfn:`Traffic Router` is used to route client requests to the closest healthy :term:`cache server` by analyzing the health, capacity, and state of the :term:`cache servers` according to the :ref:`health-proto` and relative geographic distance between each :term:`Cache Group` and the client.

:ref:`tm-overview`
	:dfn:`Traffic Monitor` does health polling of the :term:`cache servers` on a very short interval to keep track of which servers should be kept in rotation.

	.. seealso:: :ref:`health-proto`

:ref:`ts-overview`
	:dfn:`Traffic Stats` collects and stores real-time traffic statistics aggregated from each of the :term:`cache servers`. This data is used by the :ref:`tr-overview` to assess the available capacity of each :term:`cache server` which it uses to balance traffic load and prevent overload.

:ref:`tp-overview`
	:dfn:`Traffic Portal` is a web interface which uses the :ref:`to-api` to present CDN data and the controls to manipulate it in a user-friendly interface.

	.. versionadded:: 2.2
		As of Traffic Control 2.2, this is the recommended, official UI for the Traffic Control platform. In Traffic Control 3.x, the Traffic Ops UI has been deprecated and disabled by default, and it will be removed with the release of Traffic Control 4.0.

:ref:`tv-overview`
	:dfn:`Traffic Vault` is used as a secure key/value store for SSL private keys used by other Traffic Control components.
