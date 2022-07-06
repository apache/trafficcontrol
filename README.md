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

![Traffic Control Logo](https://traffic-control-cdn.readthedocs.io/en/latest/_static/ATC-SVG-FULL-WHITE.svg#gh-dark-mode-only)
![Traffic Control Logo](https://trafficcontrol.apache.org/resources/Traffic-Control-Logo-FINAL-Black-Text.png#gh-light-mode-only)

Apache Traffic Control allows you to build a large scale content delivery network using open source. Built around Apache Traffic Server as the caching software, Traffic Control implements all the core functions of a modern CDN.

[![Slack](https://img.shields.io/badge/slack-join_%23traffic--control-white.svg?logo=slack&style=social)](https://s.apache.org/tc-slack-request)
[![Twitter Follow](https://img.shields.io/twitter/follow/trafficctrlcdn?style=social&label=Follow%20@trafficctrlcdn)](https://twitter.com/intent/follow?screen_name=trafficctrlcdn)
[![Youtube Subscribe](https://img.shields.io/youtube/channel/subscribers/UC2zEj6sERinzx8w8uvyRBYg?style=social&label=Apache%20Traffic%20Control)](https://www.youtube.com/channel/UC2zEj6sERinzx8w8uvyRBYg)

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/apache/trafficcontrol)](https://github.com/apache/trafficcontrol/releases)

__Build Status__

[![Build Status](https://github.com/apache/trafficcontrol/workflows/CDN-in-a-Box%20CI/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/ciab.yaml) 
[![Documentation Status](https://readthedocs.org/projects/traffic-control-cdn/badge/?version=latest)](http://traffic-control-cdn.readthedocs.io/en/latest/?badge=latest)

__Code Status__

[![Weasel License Checks](https://github.com/apache/trafficcontrol/workflows/Weasel%20License%20checks/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/weasel.yml) 
[![Go Formatting](https://github.com/apache/trafficcontrol/workflows/Go%20Format/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/go.fmt.yml) 
[![Go Vet](https://github.com/apache/trafficcontrol/workflows/Go%20Vet/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/go.vet.yml)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            
[![CodeQL - C++](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20C++/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.cpp.yml)
[![CodeQL - Go](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Go/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.go.yml)
[![CodeQL - Java](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Java/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.java.yml)
[![CodeQL - Javascript](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Javascript/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.javascript.yml)
[![CodeQL - Python](https://github.com/apache/trafficcontrol/workflows/CodeQL%20-%20Python/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/codeql.python.yml)

__Test Status__

[![Go Unit Tests](https://github.com/apache/trafficcontrol/workflows/Go%20Unit%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/go.unit.tests.yaml)
[![Traffic Ops Integration Tests](https://github.com/apache/trafficcontrol/workflows/Traffic%20Ops%20Go%20client/API%20integration%20tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/traffic-ops.yml) 
[![TP Integration Tests](https://github.com/apache/trafficcontrol/workflows/TP%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tp.integration.tests.yml) 
[![TM Integration Tests](https://github.com/apache/trafficcontrol/workflows/TM%20Integration%20Tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tm.integration.tests.yml) 
[![TR Ultimate Test Harness](https://github.com/apache/trafficcontrol/workflows/TR%20Ultimate%20Test%20Harness/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/tr-ultimate-test-harness.yml) 
[![Traffic Control Cache Config integration tests](https://github.com/apache/trafficcontrol/workflows/Traffic%20Control%20Cache%20Config%20integration%20tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/cache-config-tests.yml)
[![Traffic Control Health Client integration tests](https://github.com/apache/trafficcontrol/workflows/Traffic%20Control%20Health%20Client%20integration%20tests/badge.svg)](https://github.com/apache/trafficcontrol/actions/workflows/health-client-tests.yml)

## Documentation[^1]
* [Intro](http://traffic-control-cdn.readthedocs.io/en/latest/index.html)
* [CDN Basics](http://traffic-control-cdn.readthedocs.io/en/latest/basics/index.html)
* [Traffic Control Overview](http://traffic-control-cdn.readthedocs.io/en/latest/overview/index.html)
* [Administrator's Guide](http://traffic-control-cdn.readthedocs.io/en/latest/admin/index.html)
* [Developer's Guide](http://traffic-control-cdn.readthedocs.io/en/latest/development/index.html)
* [Traffic Ops API](https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html)

## Components[^1]
* [Traffic Ops](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_ops.html) is the RESTful API service for management and monitoring of all servers in the CDN.
* [Traffic Portal](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_portal.html) is the web GUI for managing and monitoring the CDN via the Traffic Ops API.
* [Traffic Router](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_router.html) uses DNS and HTTP302 to redirect clients to the closest available cache on the CDN.
* [Traffic Monitor](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_monitor.html) uses HTTP to poll the health of caches and provide this information to Traffic Router.
* [Traffic Stats](https://traffic-control-cdn.readthedocs.io/en/latest/overview/traffic_stats.html) acquires and stores real-time metrics and statistics into an InfluxDB for charting and alerting.

## Releases
* [https://github.com/apache/trafficcontrol/releases](https://github.com/apache/trafficcontrol/releases)

## Downloads
* [https://www.apache.org/dyn/closer.cgi/trafficcontrol](https://www.apache.org/dyn/closer.cgi/trafficcontrol)

## Questions, Comments, Bugs and More
* [Frequently Asked Questions](https://traffic-control-cdn.readthedocs.io/en/latest/faq.html)
* [Found a bug or file a feature request](https://github.com/apache/trafficcontrol/issues)
* [Subscribe to our users list](mailto:users-subscribe@trafficcontrol.apache.org)
* [Subscribe to our dev list](mailto:dev-subscribe@trafficcontrol.apache.org)
* [Search the email archives](https://lists.apache.org/list.html?dev@trafficcontrol.apache.org)
* [Check out the wiki](https://cwiki.apache.org/confluence/display/TC/Traffic+Control+Home) for less formal documentation, design docs and roadmap discussions

[^1]: *Documentation links point to the __latest__ which is the unreleased master branch and are neither stable nor necessarily accurate for any given supported release.*