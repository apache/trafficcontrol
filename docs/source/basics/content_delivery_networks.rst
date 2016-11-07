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

.. index::
	Log File Analysis
	CDN
	Content Delivery Network 

Content Delivery Networks
=========================
The vast majority of today's Internet traffic is media files (often video or audio) being sent from a single source (the *Content Provider*) to many thousands or even millions of destinations (the *Content Consumers*).  Content Delivery Networks are the technology that make that one-to-many distribution possible in an economical way. A Content Delivery Network (CDN) is a distributed system of servers for delivering content over HTTP. These servers are deployed in multiple locations with the goal of optimizing the delivery of content to the end users, while minimizing the traffic on the network. A CDN typically consists of the following:

* **Caching Proxies**
	The proxy (cache or caching proxy) is a server that both proxies the requests and caches the results for reusing.  

* **Content Router**
    The Content Router ensures that the end user is connected to the optimal cache for the location of the end user and content availability.

* **Health Protocol** 
    The Health Protocol monitors the usage of the caches and tenants in the CDN.

* **Configuration Management System** 
    In many cases a CDN encompasses hundreds of servers across a large geographic area. The Configuration Management System allows an operator to manage these servers.

* **Log File Analysis System**
    Every transaction in the CDN gets logged. The Log File Analysis System aggregates all of the log entries from all of the servers to a central location for analysis and troubleshooting.


