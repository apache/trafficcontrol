<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# Apache Traffic Control

<picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://traffic-control-cdn.readthedocs.io/en/latest/_static/ATC-SVG-FULL-WHITE.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://trafficcontrol.apache.org/resources/Traffic-Control-Logo-FINAL-Black-Text.png">
    <img alt="Traffic Control Logo" src="https://trafficcontrol.apache.org/resources/Traffic-Control-Logo-FINAL-Black-Text.png">
</picture>

Apache Traffic Control allows you to build a large scale content delivery network using open source. Built around Apache Traffic Server as the caching software, Traffic Control implements all the core functions of a modern CDN.

[![Slack](https://img.shields.io/badge/slack-join_%23traffic--control-white.svg?logo=slack&style=social)](https://s.apache.org/tc-slack-request)
[![Twitter Follow](https://img.shields.io/twitter/follow/trafficctrlcdn?style=social&label=Follow%20@trafficctrlcdn)](https://twitter.com/intent/follow?screen_name=trafficctrlcdn)
[![Youtube Subscribe](https://img.shields.io/youtube/channel/subscribers/UC2zEj6sERinzx8w8uvyRBYg?style=social&label=Apache%20Traffic%20Control)](https://www.youtube.com/channel/UC2zEj6sERinzx8w8uvyRBYg)

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/apache/trafficcontrol)](https://github.com/apache/trafficcontrol/releases)
![Github commits since release](https://img.shields.io/github/commits-since/apache/trafficcontrol/latest/master)


__Build Status__ [^1]

[![Build Status](https://github.com/apache/trafficcontrol/workflows/CDN-in-a-Box%20CI/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/ciab.yaml?query=branch%3Amaster) 
[![Documentation Status](https://readthedocs.org/projects/traffic-control-cdn/badge/?version=latest)](http://traffic-control-cdn.readthedocs.io/en/latest/?badge=latest?query=branch%3Amaster)

__Code Status__ [^1]

[![Weasel License Checks](https://github.com/apache/trafficcontrol/workflows/Weasel%20License%20checks/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/weasel.yml?query=branch%3Amaster) 
[![Go Formatting](https://github.com/apache/trafficcontrol/workflows/Go%20Format/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/go.fmt.yml?query=branch%3Amaster) 
[![Go Vet](https://github.com/apache/trafficcontrol/workflows/Go%20Vet/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/go.vet.yml?query=branch%3Amaster)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            
[![CodeQL - C++](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20C++/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.cpp.yml?query=branch%3Amaster)
[![CodeQL - Go](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Go/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.go.yml?query=branch%3Amaster)
[![CodeQL - Java](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Java/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.java.yml?query=branch%3Amaster)
[![CodeQL - Javascript](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Javascript/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.javascript.yml?query=branch%3Amaster)
[![CodeQL - Python](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Python/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.python.yml?query=branch%3Amaster)

__Test Status__ [^1]

| Component      | Unit Tests                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Integration Tests                                                                                                                                                                                                                                                                                                                                                                                                                            | 
|:---------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| 
| Go Libraries   | [![Go Lib Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/go.lib.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/go.lib.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=golib_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/lib)                                        | -                                                                                                                                                                                                                                                                                                                                                                                                                                            | 
| Traffic Ops    | [![Traffic Ops Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/to.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/to.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=traffic_ops_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/traffic_ops) [![Traffic Ops API Contract Tests](https://github.com/apache/trafficcontrol/actions/workflows/to.api.contract.tests.yml/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/to.api.contract.tests.yml)                            | [![TO Go Client Integration Tests](https://github.com/apache/trafficcontrol/workflows/TO%20Go%20Client%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/traffic-ops.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=traffic_ops_integration)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/traffic_ops) | 
| Traffic Router | [![Traffic Router Tests](https://github.com/apache/trafficcontrol/workflows/Traffic%20Router%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tr.tests.yaml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=traffic_router_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/traffic_router)                                            | [![TR Ultimate Test Harness](https://github.com/apache/trafficcontrol/workflows/TR%20Ultimate%20Test%20Harness/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tr-ultimate-test-harness.yml?query=branch%3Amaster)                                                                                                                                                                                                    |
| Traffic Monitor| [![Traffic Monitor Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/tm.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/tm.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=traffic_monitor_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/traffic_monitor)                 | [![TM Integration Tests](https://github.com/apache/trafficcontrol/workflows/TM%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tm.integration.tests.yml?query=branch%3Amaster)                                                                                                                                                                                                                  |
| T3C            | [![T3C Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/cache-config.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/cache-config.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=t3c_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/cache-config)                        | [![T3C Integration Tests](https://github.com/apache/trafficcontrol/workflows/T3C%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/cache-config-tests.yml?query=branch%3Amaster)                                                                                                                                                                                                                  |
| Traffic Stats  | [![Traffic Stats Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/traffic-stats.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/traffic-stats.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=traffic_stats_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/traffic_stats) | -                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Grove          | [![Grove Unit Tests](https://github.com/apache/trafficcontrol/actions/workflows/grove.unit.tests.yml/badge.svg?branch=master)](https://github.com/apache/trafficcontrol/actions/workflows/grove.unit.tests.yml?query=branch%3Amaster) [![Codecov](https://codecov.io/gh/apache/trafficcontrol/branch/master/graph/badge.svg?flag=grove_unit)](https://app.codecov.io/github/apache/trafficcontrol/tree/master/grove)                                         | -                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Traffic Portal | -                                                                                                                                                                                                                                                                                                                                                                                                                                                            | [![TP Integration Tests](https://github.com/apache/trafficcontrol/workflows/TP%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tp.integration.tests.yml?query=branch%3Amaster)                                                                                                                                                                                                                  |
| TCHC           | -                                                                                                                                                                                                                                                                                                                                                                                                                                                            | [![TC Health Client Integration Tests](https://github.com/apache/trafficcontrol/workflows/TC%20Health%20Client%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/health-client-tests.yml?query=branch%3Amaster)                                                                                                                                                                                   |

---

## Documentation [^2]
* [Intro](http://traffic-control-cdn.readthedocs.io/en/latest/index.html)
* [CDN Basics](http://traffic-control-cdn.readthedocs.io/en/latest/basics/index.html)
* [Traffic Control Overview](http://traffic-control-cdn.readthedocs.io/en/latest/overview/index.html)
* [Administrator's Guide](http://traffic-control-cdn.readthedocs.io/en/latest/admin/index.html)
* [Developer's Guide](http://traffic-control-cdn.readthedocs.io/en/latest/development/index.html)
* [Traffic Ops API](https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html)

## Components [^2]
* [Traffic Ops](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_ops.html) is the RESTful API service for management and monitoring of all servers in the CDN.
* [Traffic Portal](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_portal.html) is the web GUI for managing and monitoring the CDN via the Traffic Ops API.
* [Traffic Router](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_router.html) uses DNS and HTTP302 to redirect clients to the closest available cache on the CDN.
* [Traffic Monitor](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_monitor.html) uses HTTP to poll the health of caches and provide this information to Traffic Router.
* [Traffic Stats](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_stats.html) acquires and stores real-time metrics and statistics into an InfluxDB for charting and alerting.

## Releases
* [Releases](https://github.com/apache/trafficcontrol/releases)

## Downloads
* [Downloads](https://www.apache.org/dyn/closer.cgi/trafficcontrol)

## Questions, Comments, Bugs and More
* [Frequently Asked Questions](https://traffic-control-cdn.readthedocs.io/en/latest/faq.html)
* [Found a bug or file a feature request](https://github.com/apache/trafficcontrol/issues)
* [Subscribe to our users list](mailto:users-subscribe@trafficcontrol.apache.org)
* [Subscribe to our dev list](mailto:dev-subscribe@trafficcontrol.apache.org)
* [Search the email archives](https://lists.apache.org/list.html?dev@trafficcontrol.apache.org)
* [Check out the wiki](https://cwiki.apache.org/confluence/display/TC/Traffic+Control+Home) for less formal documentation, design docs and roadmap discussions

[^1]: *Status links point to the unreleased *master* branch

[^2]: *Documentation links point to the __latest__ which is the unreleased master branch and are neither stable nor necessarily accurate for any given supported release.*
