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

Traffic Monitor Golang
**********************
Introduction
============
Traffic Monitor is an HTTP service application that monitors caches, provides health state information to Traffic Router, and collects statistics for use in tools such as Traffic Ops and Traffic Stats. The health state provided by Traffic Monitor is used by Traffic Router to control which caches are available on the CDN.

Software Requirements
=====================
To work on Traffic Monitor you need a \*nix (MacOS and Linux are most commonly used) environment that has the following installed:

* Golang

Project Tree Overview
=====================================

* ``traffic_control/traffic_monitor/`` - base directory for Traffic Monitor.

* ``cache/`` - Handler for processing cache results.
* ``config/`` - Application configuration; in-memory objects from ``traffic_monitor.cfg``.
* ``crconfig/`` - struct for deserlializing the CRConfig from JSON.
* ``deliveryservice/`` - aggregates delivery service data from cache results.
* ``deliveryservicedata/`` - deliveryservice structs. This exists separate from ``deliveryservice`` to avoid circular dependencies.
* ``enum/`` - enumerations and name alias types.
* ``health/`` - functions for calculating cache health, and creating health event objects.
* ``manager/`` - manager goroutines (microthreads).
	* ``health.go`` - Health request manager. Processes health results, from the health poller -> fetcher -> manager. The health poll is the "heartbeat" containing a small amount of stats, primarily to determine whether a cache is reachable as quickly as possible. Data is aggregated and inserted into shared threadsafe objects.
	* ``manager.go`` - Contains ``Start`` function to start all pollers, handlers, and managers.
	* ``monitorconfig.go`` - Monitor config manager. Gets data from the monitor config poller, which polls Traffic Ops for changes to which caches are monitored and how.
	* ``opsconfig.go`` - Ops config manager. Gets data from the ops config poller, which polls Traffic Ops for changes to monitoring settings.
	* ``peer.go`` - Peer manager. Gets data from the peer poller -> fetcher -> handler and aggregates it into the shared threadsafe objects.
	* ``stat.go`` - Stat request manager. Processes stat results, from the stat poller -> fetcher -> manager. The stat poll is the large statistics poll, containing all stats (such as HTTP codes, transactions, delivery service statistics, and more). Data is aggregated and inserted into shared threadsafe objects.
	* ``statecombiner.go`` - Manager for combining local and peer states, into a single combined states threadsafe object, for serving the CrStates endpoint.
* ``datareq/`` - HTTP routing, which has threadsafe health and stat objects populated by stat and health managers.
* ``peer/`` - Manager for getting and populating peer data from other Traffic Monitors
* ``srvhttp/`` - HTTP service. Given a map of endpoint functions, which are lambda closures containing aggregated data objects.
* ``static/`` - Web GUI HTML and javascript files
* ``threadsafe/`` - Threadsafe objects for storing aggregated data needed by multiple goroutines (typically the aggregator and HTTP server)
* ``trafficopsdata/`` - Struct for fetching and storing Traffic Ops data needed from the CRConfig. This is primarily mappings, such as delivery service servers, and server types.
* ``trafficopswrapper/`` - Threadsafe wrapper around the Traffic Ops client. The client used to not be threadsafe, however, it mostly (possibly entirely) is now. But, the wrapper also serves to overwrite the Traffic Ops ``monitoring.json`` values, which are live, with snapshotted CRConfig values.

Architecture
============
At the highest level, Traffic Monitor polls caches, aggregates their data and availability, and serves it at HTTP JSON endpoints.

In the code, the data flows thru microthread (goroutine) pipelines. All stages of the pipeline are independent running microthreads [#f1]_ . The pipelines are:

* **stat poll** - polls caches for all statistics data. This should be a slower poll, which gets a lot of data.
* **health poll** - polls caches for a tiny amount of data, typically system information. This poll is designed to be a heartbeat, determining quickly whether the cache is reachable. Since it's a small amount of data, it should poll more frequently.
* **peer poll** - polls Traffic Monitor peers for their availability data, and aggregates it with its own availability results and that of all other peers.
* **monitor config** - polls Traffic Ops for the list of Traffic Monitors and their info.
* **ops config** - polls for changes to the ops config file ``traffic_ops.cfg``, and sends updates to other pollers when the config file has changed.

  * The ops config manager also updates the shared Traffic Ops client, since it's the actor which becomes notified of config changes requiring a new client.

  * The ops config manager also manages, creates, and recreates the HTTP server, since ops config changes necessitate restarting the HTTP server.

All microthreads in the pipeline are started by ``manager/manager.go:Start()``.

::

  --------------------     --------------------     --------------------
  | ops config poller |-->| ops config handler |-->| ops config manager |-->-restart HTTP server-------------------------
   -------------------     --------------------     -------------------- |                                              |
                                                                         -->-ops config change subscriber-------------  |
                                                                         |                                           |  |
                                                                         -->-Traffic Ops client change subscriber--  |  |
                                                                                                                  |  |  |
      -------------------------------------------------------------------------------------------------------------  |  |
      |                                                                                                              |  |
      |   ------------------------------------------------------------------------------------------------------------  |
      |   |                                                                                                             |
      \/  \/                                                                                                            |
     -----------------------     ------------------------                                                               |
    | monitor config poller |-->| monitor config manager |-->-stat subscriber--------             -----------------------
     -----------------------     ------------------------ |                         |             |
                                                          |->-health subscriber---  |             \/                           _
                                                          |                      |  |       -------------                    _( )._
                                                          -->-peer subscriber--  |  |      | HTTP server |->-HTTP request-> (____)_)
                                                                              |  |  |       -------------
  -----------------------------------------------------------------------------  |  |              ^
  |                                                                              |  |              |
  |  -----------------------------------------------------------------------------  |              ------------------------
  |  |                                                                              |                                     |
  |  |  -----------------------------------------------------------------------------                                     |
  |  |  |                                                                                                                 ^
  |  |  |   -------------     --------------     --------------     --------------                            -----------------------
  |  |  -->| stat poller |-->| stat fetcher |-->| stat handler |-->| stat manager |->--------set shared data->| shared data         |
  |  |      ------------- |   --------------     --------------  |  --------------                            -----------------------
  |  |                    |   --------------     --------------  |                                            | events              |
  |  |                    |->| stat fetcher |-->| stat handler |-|                                            | toData              |
  |  |                    |   --------------     --------------  |                                            | errorCount          |
  |  |                    ...                                    ...                                          | healthIteration     |
  |  |                                                                                                        | fetchCount          |
  |  |     ---------------     ----------------     ----------------     ----------------                     | localStates         |
  |  ---->| health poller |-->| health fetcher |-->| health handler |-->| health manager |->-set shared data->| toSession           |
  |        --------------- |   ----------------     ----------------  |  ----------------                     | peerStates          |
  |                        |   ----------------     ----------------  |                                       | monitorConfig       |
  |                        |->| health fetcher |-->| health handler |-|                                       | combinedStates      |
  |                        |   ----------------     ----------------  |                                       | statInfoHistory     |
  |                        ...                                        ...                                     | statResultHistory   |
  |                                                                                                           | statMaxKbpses       |
  |       -------------     --------------     --------------     --------------                              | lastKbpsStats       |
  ------>| peer poller |-->| peer fetcher |-->| peer handler |-->| peer manager |->----------set shared data->| dsStats             |
          ------------- |   --------------     --------------  |  --------------                              | localCacheStatus    |
                        |   --------------     --------------  |                                              | lastHealthDurations |
                        |->| peer fetcher |-->| peer handler |-|                                              | healthHistory       |
                        |   --------------     --------------  |                                              -----------------------
                        ...                                    ...

.. [#f1] Technically, some stages which are one-to-one simply call the next stage as a function. For example, the Fetcher calls the Handler as a function in the same microthread. But this isn't architecturally significant.


Stat Pipeline
-------------

::

  ---------     ---------     ---------     ---------
  | poller |-->| fetcher |-->| handler |-->| manager |
   -------- |   ---------     ---------  |  ---------
            |   ---------     ---------  |
            |->| fetcher |-->| handler |-|
            |   ---------     ---------  |
            ...                          ...

* **poller** - ``common/poller/poller.go:HttpPoller.Poll()``. Listens for config changes (from the ops config manager), and starts its own internal microthreads, one for each cache to poll. These internal microthreads call the Fetcher at each cache's poll interval.

* **fetcher** - ``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Fetches the given URL, and passes the returned data to the Handler, along with any errors.


* **handler** - ``traffic_monitor/cache/cache.go:Handler.Handle()``. Takes the given result and does all data computation possible with the single result. Currently, this computation primarily involves processing the denormalized ATS data into Go structs, and processing System data into OutBytes, Kbps, etc. Precomputed data is then passed to its result channel, which is picked up by the Manager.

* **manager** - ``traffic_monitor/manager/stat.go:StartStatHistoryManager()``. Takes preprocessed results, and aggregates them. Aggregated results are then placed in shared data structures. The major data aggregated are delivery service statistics, and cache availability data. See :ref:`Aggregated Stat Data <agg-stat-data>` and :ref:`Aggregated Availability Data <agg-avail-data>`.


Health Pipeline
---------------

::

  ---------     ---------     ---------     ---------
  | poller |-->| fetcher |-->| handler |-->| manager |
   -------- |   ---------     ---------  |  ---------
            |   ---------     ---------  |
            |->| fetcher |-->| handler |-|
            |   ---------     ---------  |
            ...                          ...

* **poller** - ``common/poller/poller.go:HttpPoller.Poll()``. Same poller type as the Stat Poller pipeline, with a different handler object.

* **fetcher** - ``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Same fetcher type as the Stat Poller pipeline, with a different handler object.

* **handler** - ``traffic_monitor/cache/cache.go:Handler.Handle()``. Same handler type as the Stat Poller pipeline, but constructed with a flag to not precompute. The health endpoint is of the same form as the stat endpoint, but doesn't return all stat data. So, it doesn't precompute like the Stat Handler, but only processes the system data, and passes the processed result to its result channel, which is picked up by the Manager.

* **manager** - ``traffic_monitor/manager/health.go:StartHealthResultManager()``. Takes preprocessed results, and aggregates them. For the Health pipeline, only health availability data is aggregated. Aggregated results are then placed in shared data structures (lastHealthDurationsThreadsafe, lastHealthEndTimes, etc). See :ref:`Aggregated Availability Data <agg-avail-data>`.


Peer Pipeline
-------------

::

  ---------     ---------     ---------     ---------
  | poller |-->| fetcher |-->| handler |-->| manager |
   -------- |   ---------     ---------  |  ---------
            |   ---------     ---------  |
            |->| fetcher |-->| handler |-|
            |   ---------     ---------  |
            ...                          ...

* **poller** - ``common/poller/poller.go:HttpPoller.Poll()``. Same poller type as the Stat and Health Poller pipelines, with a different handler object. Its config changes come from the Monitor Config Manager, and it starts an internal microthread for each peer to poll.

* **fetcher** - ``common/fetcher/fetcher.go:HttpFetcher.Fetch()``. Same fetcher type as the Stat and Health Poller pipeline, with a different handler object.

* **handler** - ``traffic_monitor/cache/peer.go:Handler.Handle()``. Decodes the JSON result into an object, and without further processing passes to its result channel, which is picked up by the Manager.

* **manager** - ``traffic_monitor/manager/peer.go:StartPeerManager()``. Takes JSON peer Traffic Monitor results, and aggregates them. The availability of the Peer Traffic Monitor itself, as well as all cache availability from the given peer result, is stored in the shared ``peerStates`` object. Results are then aggregated via a call to the ``combineState()`` lambda, which signals the State Combiner microthread (which stores the combined availability in the shared object ``combinedStates``; See :ref:`State Combiner <state-combiner>`).


Monitor Config Pipeline
-----------------------

::

  ---------     ---------
  | poller |-->| manager |--> stat subscriber (Stat pipeline Poller)
   --------     --------- |
                          |-> health subscriber (Health pipeline Poller)
                          |
                          --> peer subscriber (Peer pipeline Poller)

* **poller** - ``common/poller/poller.go:MonitorConfigPoller.Poll()``. The Monitor Config poller, on its interval, polls Traffic Ops for the Monitor configuration, and writes the polled value to its result channel, which is read by the Manager.

* **manager** - ``traffic_monitor/manager/monitorconfig.go:StartMonitorConfigManager()``. Listens for results from the poller, and processes them. Cache changes are written to channels read by the Health, Stat, and Peer pollers. In the Shared Data objects, this also sets the list of new delivery services and removes ones which no longer exist, and sets the list of peer Traffic Monitors.


Ops Config Pipeline
-------------------
::

  ---------     ---------     ---------
  | poller |-->| handler |-->| manager |--> ops config change subscriber (Monitor Config Poller)
   --------     ---------     --------- |
                                        --> Traffic ops client change subscriber (Monitor Config Poller)

* **poller** - ``common/poller/poller.go:FilePoller.Poll()``. Polls for changes to the Traffic Ops config file ``traffic_ops.cfg``, and writes the changed config to its result channel, which is read by the Handler.

* **handler** - ``common/handler/handler.go:OpsConfigFileHandler.Listen()``. Takes the given raw config, unmarshalls the JSON into an object, and writes the object to its channel, which is read by the Manager, along with any error.

* **manager** - ``traffic_monitor/manager/monitorconfig.go:StartMonitorConfigManager()``. Listens for new configs, and processes them. When a new config is received, a new HTTP dispatch map is created via ``traffic_monitor/datareq/datareq.go:MakeDispatchMap()``, and the HTTP server is restarted with the new dispatch map. The Traffic Ops client is also recreated, and stored in its shared data object. The Ops Config change subscribers and Traffic Ops Client change subscribers (the Monitor Config poller) are also passed the new ops config and new Traffic Ops client.


Events
------
The ``events`` shared data object is passed to each pipeline microthread which needs to signal events. Most of them do. Events are then logged, and visible in the UI as well as an HTTP JSON endpoint. Most events are caches becoming available or unavailable, but include other things such as peer availability changes.


.. _state-combiner:

State Combiner
--------------
The State Combiner is a microthread started in ``traffic_monitor/manager/manager.go:Start()`` via ``traffic_monitor/manager/statecombiner.go:StartStateCombiner()``, which listens for signals to combine states. It should be signaled by any pipeline which updates the local or peer availability shared data objects, ``localStates`` and ``peerStates``. It holds the threadsafe shared data objects for local states and peer states, so no data is passed or returned, only a signal.

When a signal is received, it combines the local and peer states optimistically. That is, if a cache is marked available locally or by any peer, that cache is marked available in the combined states. There exists a variable to combine pessimistically, which may be set at compile time (it's unusual for a CDN to operate well with pessimistic cache availability). Combined data is stored in the threadsafe shared data object ``combinedStates``.


.. _agg-stat-data:

Aggregated Stat Data
--------------------
The Stat pipeline Manager is responsible for aggregating stats from all caches, into delivery services statistics. This is done via a call to ``traffic_monitor/deliveryservice/stat.go:CreateStats()``.


.. _agg-avail-data:

Aggregated Availability Data
----------------------------
Both the Stat and Health pipelines aggregate availability data received from caches. This is done via a call to ``traffic_monitor/deliveryservice/health.go:CalcAvailability()`` followed by a call to ``combineState()``. The ``CalcAvailability`` function calculates the availability of each cache from the result of polling it, that is, local availability. The ``combineState()`` function is a lambda passed to the Manager, which signals the State Combiner microthread, which will combine the local and peer Traffic Monitor availability data, and insert it into the shared data ``combinedStates`` object.


HTTP Data Requests
------------------
Data is provided to HTTP requests via the threadsafe shared data objects (see :ref:`Shared Data <shared-data>`). These objects are closed in lambdas created via ``traffic_monitor/datareq/datareq.go:MakeDispatchMap()``. This is called by the Ops Config Manager when it recreates the HTTP server.

Each HTTP endpoint is mapped to a function which closes around the shared data objects it needs, and takes the request data it needs (such as query parameters). Each endpoint function resides in its own file in ``traffic_monitor/datareq/``. Because each Go HTTP routing function must be a ``http.HandlerFunc``, wrapper functions take the endpoint functions and return ``http.HandlerFunc`` functions which call them, and which are stored in the dispatch map, to be registered with the HTTP server.


.. _shared-data:

Shared Data
-----------
Processed and aggregated data must be shared between the end of the stat and health processing pipelines, and HTTP requests. The CSP paradigm of idiomatic Go does not work efficiently with storing and sharing state. While not idiomatic Go, shared mutexed data structures are faster and simpler than CSP manager microthreads for each data object.

Traffic Monitor has many threadsafe shared data types and objects. All shared data objects can be seen in ``manager/manager.go:Start()``, where they are created and passed to the various pipeline stage microthreads that need them. Their respective types all include the word ``Threadsafe``, and can be found in ``traffic_monitor/threadsafe/`` as well as, for dependency reasons, various appropriate directories.

Currently, all Threadsafe shared data types use mutexes. In the future, these could be changed to lock-free or wait-free structures, if the performance needs outweighed the readability and correctness costs. They could also easily be changed to internally be manager microthreads and channels, if being idiomatic were deemed more important than readability or performance.



Formatting Conventions
======================
Go code should be formatted with ``gofmt``. See also ``CONTRIBUTING.md``.

Installing The Developer Environment
====================================
To install the Traffic Monitor Developer environment:

1. Install `go` version 1.7 or greater, from https://golang.org/doc/install and https://golang.org/doc/code.html
2. Clone the traffic_control repository using Git, into ``$GOPATH/src/github.com/apache/incubator-trafficcontrol``
3. Change directories into ``$GOPATH/src/github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor``
4. Run ``./build.sh``

Test Cases
==========
Tests can be executed by running ``go test ./...`` at the root of the ``traffic_monitor_golang`` project.

API
===

:ref:`reference-tm-api`

.. toctree:: 
  :hidden:
  :maxdepth: 1

  traffic_monitor/traffic_monitor_api
