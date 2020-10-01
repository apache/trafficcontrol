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

.. _dev-traffic-monitor:

***************
Traffic Monitor
***************
Introduction
============
Traffic Monitor is an HTTP service application that monitors :term:`cache servers`, provides health state information to Traffic Router, and collects statistics for use in tools such as Traffic Portal and Traffic Stats. The health state provided by Traffic Monitor is used by Traffic Router to control which :term:`cache servers` are available on the CDN.

Software Requirements
=====================
To work on Traffic Monitor you need a Unix-like (MacOS and Linux are most commonly used) environment that has a working install of Go version 1.7+.

Project Tree Overview
=====================

``traffic_monitor/`` - base directory for Traffic Monitor.

* ``cache/`` - Handler for processing cache results.
* ``config/`` - Application configuration; in-memory objects from ``traffic_monitor.cfg``.
* ``crconfig/`` - data structure for deserlializing the CDN :term:`Snapshot` (historically named "CRConfig") from JSON.
* ``deliveryservice/`` - aggregates :term:`Delivery Service` data from cache results.
* ``deliveryservicedata/`` - :term:`Delivery Service` data structures. This exists separate from ``deliveryservice`` to avoid circular dependencies.
* ``enum/`` - enumerations and name alias types.
* ``health/`` - functions for calculating :term:`cache server` health, and creating health event objects.
* ``manager/`` - manager goroutines (microthreads).

	* ``health.go`` - Health request manager. Processes health results, from the health poller -> fetcher -> manager. The health poll is the "heartbeat" containing a small amount of stats, primarily to determine whether a :term:`cache server` is reachable as quickly as possible. Data is aggregated and inserted into shared thread-safe objects.
	* ``manager.go`` - Contains ``Start`` function to start all pollers, handlers, and managers.
	* ``monitorconfig.go`` - Monitor configuration manager. Gets data from the monitor configuration poller, which polls Traffic Ops for changes to which caches are monitored and how.
	* ``opsconfig.go`` - Ops configuration manager. Gets data from the ops configuration poller, which polls Traffic Ops for changes to monitoring settings.
	* ``peer.go`` - Peer manager. Gets data from the peer poller -> fetcher -> handler and aggregates it into the shared thread-safe objects.
	* ``stat.go`` - Stat request manager. Processes stat results, from the stat poller -> fetcher -> manager. The stat poll is the large statistics poll, containing all stats (such as HTTP response status codes, transactions, :term:`Delivery Service` statistics, and more). Data is aggregated and inserted into shared thread-safe objects.
	* ``statecombiner.go`` - Manager for combining local and peer states, into a single combined states thread-safe object, for serving the CrStates endpoint.

* ``datareq/`` - HTTP routing, which has thread-safe health and stat objects populated by stat and health managers.
* ``peer/`` - Manager for getting and populating peer data from other Traffic Monitors
* ``srvhttp/`` - HTTP(S) service. Given a map of endpoint functions, which are lambda closures containing aggregated data objects. If HTTPS, the HTTP service will redirect to HTTPS.
* ``static/`` - Web interface files (markup, styling and scripting)
* ``threadsafe/`` - Thread-safe objects for storing aggregated data needed by multiple goroutines (typically the aggregator and HTTP server)
* ``todata/`` - Data structure for fetching and storing Traffic Ops data needed from the CDN :term:`Snapshot`. This is primarily mappings, such as :term:`Delivery Service` servers, and server types.
* ``trafficopswrapper/`` - Thread-safe wrapper around the Traffic Ops client. The client used to not be thread-safe, however, it mostly (possibly entirely) is now. But, the wrapper also serves to overwrite the Traffic Ops :file:`monitoring.json` values, which are live, with values from the CDN :term:`Snapshot`.

Architecture
============
At the highest level, Traffic Monitor polls :term:`cache server` s, aggregates their data and availability, and serves it at HTTP endpoints in JSON format. In the code, the data flows through microthread (goroutine) pipelines. All stages of the pipeline are independently running microthreads\ [1]_ . The pipelines are:

stat poll
	Polls caches for all statistics data. This should be a slower poll, which gets a lot of data.
health poll
	Polls caches for a tiny amount of data, typically system information. This poll is designed to be a heartbeat, determining quickly whether the :term:`cache server` is reachable. Since it's a small amount of data, it should poll more frequently.
peer poll
	Polls Traffic Monitor peers for their availability data, and aggregates it with its own availability results and that of all other peers.
monitor config
	Polls Traffic Ops for the list of Traffic Monitors and their info.
ops config
	Polls for changes to the Traffic Ops configuration file :file:`traffic_ops.cfg`, and sends updates to other pollers when the configuration file has changed.

	* The ops config manager also updates the shared Traffic Ops client, since it's the actor which becomes notified of configuration changes requiring a new client.
	* The ops config manager also manages, creates, and recreates the HTTP server, since Traffic Ops configuration changes necessitate restarting the HTTP server.

All microthreads in the pipeline are started by ``manager/manager.go:Start()``.

.. figure:: traffic_monitor/Pipeline.*
	:align: center
	:width: 70%

	Pipeline Overview

.. [1] Technically, some stages which are one-to-one simply call the next stage as a function. For example, the Fetcher calls the Handler as a function in the same microthread. But this isn't architecturally significant.

Stat Pipeline
-------------
.. figure:: traffic_monitor/Stat_Pipeline.*
	:align: center
	:width: 70%

	The Stats Pipeline

poller
	``common/poller/poller.go:HttpPoller.Poll()``. Listens for configuration changes (from the Ops Configuration Manager), and starts its own, internal microthreads - one for each cache to poll. These internal microthreads call the Fetcher at each cache's poll interval.
fetcher
	``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Fetches the given URL, and passes the returned data to the Handler, along with any errors.
handler
	``traffic_monitor/cache/cache.go:Handler.Handle()``\ . Takes the given result and does all data computation possible with the single result. Currently, this computation primarily involves processing the de-normalized :abbr:`ATS (Apache Trafficserver)` data into Go ``struct``\ s, and processing System data into 'OutBytes', :abbr:`Kbps (kilobits per second)`, etc. Precomputed data is then passed to its result channel, which is picked up by the Manager.
manager
	``traffic_monitor/manager/stat.go:StartStatHistoryManager()``. Takes preprocessed results, and aggregates them. Aggregated results are then placed in shared data structures. The major data aggregated are :term:`Delivery Service` statistics, and :term:`cache server` availability data. See `Aggregated Stat Data`_ and `Aggregated Availability Data`_.


Health Pipeline
---------------
.. figure:: traffic_monitor/Health_Pipeline.*
	:align: center
	:width: 70%

	The Health Pipeline

poller
	``common/poller/poller.go:HttpPoller.Poll()``. Same poller type as the Stat Poller pipeline, with a different handler object.
fetcher
	``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Same fetcher type as the Stat Poller pipeline, with a different handler object.
handler
	``traffic_monitor/cache/cache.go:Handler.Handle()``. Same handler type as the Stat Poller pipeline, but constructed with a flag to not pre-compute anything. The health endpoint is of the same form as the stat endpoint, but doesn't return all stat data. So, it doesn't pre-compute like the Stat Handler, but only processes the system data, and passes the processed result to its result channel, which is picked up by the Manager.
manager
	``traffic_monitor/manager/health.go:StartHealthResultManager()``. Takes preprocessed results, and aggregates them. For the Health pipeline, only health availability data is aggregated. Aggregated results are then placed in shared data structures (lastHealthDurationsThreadsafe, lastHealthEndTimes, etc). See `Aggregated Availability Data`_.


Peer Pipeline
-------------
.. figure:: traffic_monitor/Peer_Pipeline.*
	:align: center
	:width: 70%

	The Peers Pipeline

poller
	``common/poller/poller.go:HttpPoller.Poll()``. Same poller type as the Stat and Health Poller pipelines, with a different handler object. Its configuration changes come from the Monitor Configuration Manager, and it starts an internal microthread for each peer to poll.
fetcher
	``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Same fetcher type as the Stat and Health Poller pipeline, with a different handler object.
handler
	``traffic_monitor/cache/peer.go:Handler.Handle()``. Decodes the JSON result into an object, and without further processing passes to its result channel, which is picked up by the Manager.
manager
	``traffic_monitor/manager/peer.go:StartPeerManager()``. Takes JSON peer Traffic Monitor results, and aggregates them. The availability of the Peer Traffic Monitor itself, as well as all :term:`cache server` availability from the given peer result, is stored in the shared ``peerStates`` object. Results are then aggregated via a call to the ``combineState()`` lambda, which signals the State Combiner microthread (which stores the combined availability in the shared object ``combinedStates``; See `State Combiner`_).


Monitor Config Pipeline
-----------------------
.. figure:: traffic_monitor/Monitor_Pipeline.*
	:align: center
	:width: 70%

	The Monitor Configuration Pipeline

poller
	``common/poller/poller.go:MonitorConfigPoller.Poll()``. The Monitor Configuration poller, on its interval, polls Traffic Ops for the Monitor configuration, and writes the polled value to its result channel, which is read by the Manager.
manager
	``traffic_monitor/manager/monitorconfig.go:StartMonitorConfigManager()``. Listens for results from the poller, and processes them. Cache changes are written to channels read by the Health, Stat, and Peer pollers. In the Shared Data objects, this also sets the list of new :term:`Delivery Services` and removes ones which no longer exist, and sets the list of peer Traffic Monitors.


Ops Config Pipeline
-------------------
.. figure:: traffic_monitor/Ops-Config_Pipeline.*
	:align: center
	:width: 70%

	The Ops Configuration Pipeline

poller
	``common/poller/poller.go:FilePoller.Poll()``. Polls for changes to the Traffic Ops configuration file :file:`traffic_ops.cfg`, and writes the changed configuration to its result channel, which is read by the Handler.
handler
	``common/handler/handler.go:OpsConfigFileHandler.Listen()``. Takes the given raw configuration, un-marshals the JSON into an object, and writes the object to its channel, which is read by the Manager, along with any error.
manager
	``traffic_monitor/manager/monitorconfig.go:StartMonitorConfigManager()``. Listens for new configurations, and processes them. When a new configuration is received, a new HTTP dispatch map is created via ``traffic_monitor/datareq/datareq.go:MakeDispatchMap()``, and the HTTP server is restarted with the new dispatch map. The Traffic Ops client is also recreated, and stored in its shared data object. The Ops Configuration change subscribers and Traffic Ops Client change subscribers (the Monitor Configuration poller) are also passed the new Traffic Ops configuration and new Traffic Ops client.

Events
------
The ``events`` shared data object is passed to each pipeline microthread which needs to signal events. Most of them do. Events are then logged, and visible in the UI as well as an HTTP JSON endpoint. Most events are :term:`cache server` becoming available or unavailable, but include other things such as peer availability changes.

State Combiner
--------------
The State Combiner is a microthread started in ``traffic_monitor/manager/manager.go:Start()`` via ``traffic_monitor/manager/statecombiner.go:StartStateCombiner()``, which listens for signals to combine states. It should be signaled by any pipeline which updates the local or peer availability shared data objects, ``localStates`` and ``peerStates``. It holds the thread-safe shared data objects for local states and peer states, so no data is passed or returned, only a signal. When a signal is received, it combines the local and peer states optimistically. That is, if a :term:`cache server` is marked available locally or by any peer, that :term:`cache server` is marked available in the combined states. There exists a variable to combine pessimistically, which may be set at compile time (it's unusual for a CDN to operate well with pessimistic :term:`cache server` availability). Combined data is stored in the thread-safe shared data object ``combinedStates``.

Aggregated Stat Data
--------------------
The Stat pipeline Manager is responsible for aggregating stats from all :term:`Edge-tier cache servers`, into :term:`Delivery Services` statistics. This is done via a call to ``traffic_monitor/deliveryservice/stat.go:CreateStats()``.

Aggregated Availability Data
----------------------------
Both the Stat and Health pipelines aggregate availability data received from caches. This is done via a call to ``traffic_monitor/deliveryservice/health.go:CalcAvailability()`` followed by a call to ``combineState()``. The ``CalcAvailability`` function calculates the availability of each :term:`cache server` from the result of polling it, that is, local availability. The ``combineState()`` function is a lambda passed to the Manager, which signals the State Combiner microthread, which will combine the local and peer Traffic Monitor availability data, and insert it into the shared data ``combinedStates`` object.

HTTP Data Requests
------------------
Data is provided to HTTP requests via the thread-safe shared data objects (see `Shared Data`_). These objects are closed in lambdas created via ``traffic_monitor/datareq/datareq.go:MakeDispatchMap()``. This is called by the Ops Configuration Manager when it recreates the HTTP(S) server. Each HTTP(S) endpoint is mapped to a function which closes around the shared data objects it needs, and takes the request data it needs (such as query parameters). Each endpoint function resides in its own file in ``traffic_monitor/datareq/``. Because each Go HTTP routing function must be a ``http.HandlerFunc``, wrapper functions take the endpoint functions and return ``http.HandlerFunc`` functions which call them, and which are stored in the dispatch map, to be registered with the HTTP(S) server.

Shared Data
-----------
Processed and aggregated data must be shared between the end of the stat and health processing pipelines, and HTTP requests. The :abbr:`CSP (Communicating Sequential Processes)` paradigm of idiomatic Go does not work efficiently with storing and sharing state. While not idiomatic Go, shared mutexed data structures are faster and simpler than :abbr:`CSP (Communicating Sequential Processes)` manager microthreads for each data object. Traffic Monitor has many thread-safe shared data types and objects. All shared data objects can be seen in ``manager/manager.go:Start()``, where they are created and passed to the various pipeline stage microthreads that need them. Their respective types all include the word ``Threadsafe``, and can be found in ``traffic_monitor/threadsafe/`` as well as, for dependency reasons, various appropriate directories. Currently, all thread-safe shared data types use mutexes. In the future, these could be changed to lock-free or wait-free structures, if the performance needs outweighed the readability and correctness costs. They could also easily be changed to internally be manager microthreads and channels, if being idiomatic were deemed more important than readability or performance.

Disk Backup
------------
The :ref:`Traffic Monitor configuration <to-api-cdns-name-snapshot>` and :ref:`CDN Snapshot <to-api-cdns-name-snapshot>` (see :term:`Snapshots`) are both stored as backup files (:file:`tmconfig.backup` and :file:`crconfig.backup` or whatever you set the values to in the configuration file). This allows the monitor to come up and continue serving even if Traffic Ops is down. These files are updated any time a valid configuration is received from Traffic Ops, so if Traffic Ops goes down and Traffic Monitor is restarted it can still serve the previous data. These files can also be manually edited and the changes will be reloaded into Traffic Monitor so that if Traffic Ops is down or unreachable for an extended period of time manual updates can be done. If on initial startup Traffic Ops is unavailable then Traffic Monitor will continue through its exponential back-off until it hits the max retry interval, at that point it will create an unauthenticated Traffic Ops session and use the data from disk. It will still poll Traffic Ops for updates though and if it successfully gets through then it will login at that point.

Formatting Conventions
======================
Go code should be formatted with ``gofmt``. See also :file:`CONTRIBUTING.md`.

Installing The Developer Environment
====================================
To install the Traffic Monitor Developer environment:

#. Install `Go <https://golang.org/doc/install>`_ version 1.7 or greater
#. Clone the `Traffic Control repository <https://github.com/apache/trafficcontrol>`_ using ``git``, into ``$GOPATH/src/github.com/apache/trafficcontrol``
#. Change directories into ``$GOPATH/src/github.com/apache/trafficcontrol/traffic_monitor``
#. Run ``./build.sh``

Test Cases
==========
Tests can be executed by running ``go test ./...`` at the root of the ``traffic_monitor`` project.

API
===

:ref:`tm-api`

.. toctree::
	:hidden:
	:maxdepth: 1

	traffic_monitor/traffic_monitor_api
