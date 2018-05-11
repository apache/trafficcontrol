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

.. _reference-tr-api:


Traffic Router API
==================

**/crs/stats**

General stats.

|

**/crs/stats/ip/<ipaddress>**

Geolocation information for an IPv4 or IPv6 address.

|

**/crs/locations**

A list of configured cache groups.

|

**/crs/locations/caches**

A mapping of caches to cache groups and their current health state.

|

**/crs/locations/<location>/caches**

A list of caches for this cache group only.

|

**/crs/consistenthash/cache/coveragezone/?ip=<ip>&deliveryServiceId=<deliveryServiceId>&requestPath=<requestPath>**

The resulting cache of the consistent hash using coverage zone file for a given client IP, delivery service, and request path.

|

**/crs/consistenthash/cache/deep/coveragezone/?ip=<ip>&deliveryServiceId=<deliveryServiceId>&requestPath=<requestPath>**

The resulting cache of the consistent hash using deep coverage zone file (deep caching) for a given client IP, delivery service, and request path.

|

**/crs/consistenthash/cache/geolocation/?ip=<ip>&deliveryServiceId=<deliveryServiceId>&requestPath=<requestPath>**

The resulting cache of the consistent hash using geolocation for a given client IP, delivery service, and request path.

|

**/crs/consistenthash/deliveryservice/?deliveryServiceId=<deliveryServiceId>&requestPath=<requestPath>**

The resulting delivery service of the consistent hash for a given delivery service and request path -- used to test steering delivery services.

|

**/crs/coveragezone/caches/?deliveryServiceId=<deliveryServiceId>&cacheLocationId=<cacheLocationId>**

A list of caches for a given delivery service and cache location.

|

**/crs/coveragezone/cachelocation/?ip=<ip>&deliveryServiceId=<deliveryServiceId>**

The resulting cache location for a given client IP and delivery service.

|

**/crs/deepcoveragezone/cachelocation/?ip=<ip>&deliveryServiceId=<deliveryServiceId>**

The resulting cache location using deep coverage zone file (deep caching) for a given client IP and delivery service.

