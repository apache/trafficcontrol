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

*************************
Content Delivery Networks
*************************
The vast majority of today's Internet traffic is media files (often video or audio) being sent from a single source (the *Content Provider*) to many thousands or even millions of destinations (the *Content Consumers*). :abbr:`CDN (Content Delivery Network)`\ s are the technology that make that one-to-many distribution efficient. A :abbr:`CDN (Content Delivery Network)` is a distributed system of servers for delivering content over HTTP(S). These servers are deployed in multiple locations with the goal of optimizing the delivery of content to the end users, while minimizing the traffic on the network. A :abbr:`CDN (Content Delivery Network)` typically consists of the following:

:term:`cache servers`
	The :dfn:`cache server` is a server that both proxies the requests and caches the results for reuse. Traffic Control uses `Apache Traffic Server <http://trafficserver.apache.org/>`_ to provide :term:`cache servers`.

Content Router
	A :dfn:`content router` ensures that the end user is connected to the optimal :term:`cache server` for the location of the end user and content availability. Traffic Control uses :ref:`tr-overview` as a :dfn:`content router`.

Health Protocol
	The :ref:`health-proto` monitors the usage of the :term:`cache servers` and tenants in the :abbr:`CDN (Content Delivery Network)`.

Configuration Management System
	In many cases a :abbr:`CDN (Content Delivery Network)` encompasses hundreds or even thousands of servers across a large geographic area. In such cases, manual configuration of servers becomes impractical, and so a central authority on configuration is used to automate the tasks as much as possible. :ref:`to-overview` is the Traffic Control configuration management system, which is interacted with via :ref:`tp-overview`.

Log File Analysis System
	Statistics and analysis are extremely important to the management and administration of a :abbr:`CDN (Content Delivery Network)`. Transaction logs and usage statistics for a Traffic Control :abbr:`CDN (Content Delivery Network)` are gathered into :ref:`ts-overview`.
