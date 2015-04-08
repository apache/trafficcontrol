.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

Traffic Ops APIs
****************

Working with APIs
=================

.. Path Structure 
.. --------------
.. dew!

User Requirements
-----------------
Verify you have the following:

* An API key.
* A Portal or Traffic Ops user account.

.. ODOL support email addy?

Using Unauthenticated API Endpoints
-----------------------------------
To use unauthenticated API endpoints, point to them and append your API key:

Example:

::

  curl -H "Accept: application/json" https://your.domain.com/cdn/customer/api/1.1/usage/overview.json?api_key=your_api_key

Using Authenticated API Endpoints
---------------------------------
To use authenticated API endpoints:

1. Authenticate with your Portal or Traffic Ops user account credentials:

::

  curl -H "Accept: application/json" -v -X POST --data '{ "u":"username", "p":"password" }' https://your.domain.com/cdn/customer/api/1.1/user/login?api_key=your_api_key

2. Upon successful user authentication, note the mojolicious cookie value in the response headers. 
3. Pass the Mojolicious cookie value, along with any subsequent calls to an authenticated API endpoint:

Examples:

::

  curl -H "Accept: application/json" -H "Cookie: mojolicious=enter-your-cookie" https://your.domain.com/cdn/customer/api/1.1/user/current.json?api_key=your_api_key

  curl -H "Accept: application/json" -H "Cookie: mojolicious=enter-your-cookie" https://your.domain.com/cdn/customer/api/1.1/deliveryservices/[deliveryservice-id].json?api_key=your_api_key

  curl -H "Accept: application/json" -H "Cookie: mojolicious=enter-your-cookie" -v -X POST --data \'{ "user": {"addressLine1":"", "addressLine2":"", "city":"", "company":"", "country":"", "email": "email@email.com", "fullName":"Full Name", "phoneNumber":"", "postalCode":"", "stateOrProvince":"", "username":"myusername" } }' https://your.domain.com/cdn/customer/api/1.1/user/current?api_key=your_api_key


API Reference 1.1
=================
:ref:`to-api-cachegroup`

* GET /api/1.1/cachegroups.json
* GET /api/1.1/cachegroups/trimmed.json
* GET /api/1.1/cachegroup/:parameter_id/parameter.json
* GET /api/1.1/cachegroupparameters.json
* GET /api/1.1/cachegroups/:parameter_id/parameter/available.json

**CDN**

* :ref:`to-api-cdn-health`
  
  * GET /api/1.1/cdns/health.json       
  * GET /api/1.1/cdns/:name/health.json
  * GET /api/1.1/cdns/usage/overview.json
  * GET /api/1.1/cdns/capacity.json

* :ref:`to-api-cdn-routing`

  * GET /api/1.1/cdns/routing.json

* :ref:`to-api-cdn-metrics`

  * GET /api/1.1/cdns/metric_types/:metric/start_date/:start/end_date/:end.json

* :ref:`to-api-cdn-domains`

  * GET /api/1.1/cdns/domains.json

* :ref:`to-api-cdn-topology`

  * GET /api/1.1/cdns/configs.json
  * GET /api/1.1/cdns/:name/configs/monitoring.json
  * GET /api/1.1/cdns/:name/configs/routing.json

* :ref:`to-api-cdn-dnsseckeys`

  * GET /api/1.1/cdns/name/:name/dnsseckeys.json
  * GET /api/1.1/cdns/name/:name/dnsseckeys/delete.json
  * POST /api/1.1/cdns/dnsseckeys/generate

:ref:`to-api-change-logs`

* GET /api/1.1/logs.json
* GET /api/1.1/logs/:days/days.json
* GET /api/1.1/logs/newcount.json

:ref:`to-api-asn`

* GET /api/1.1/asns.json

:ref:`to-api-hwinfo`

* GET /api/1.1/hwinfo.json

**Delivery Service**

* :ref:`to-api-ds`
  
  * GET /api/1.1/deliveryservices.json
  * GET /api/1.1/deliveryservices/:id.json

* :ref:`to-api-ds-health`
  
  * GET /api/1.1/deliveryservices/:id/capacity.json
  * GET /api/1.1/deliveryservices/:id/routing.json

.. * GET /api/1.1/deliveryservices/:id/state.json
.. * GET /api/1.1/deliveryservices/:id/health.json

* :ref:`to-api-ds-metrics`

  * GET /api/1.1/deliveryservices/:id/edge/metric_types/:metric/start_date/:start/end_date/:end/interval/:interval/window_start/:window_start/window_end/:window_end.json
  * GET /api/1.1/usage/deliveryservices/:ds/cachegroups/:name/metric_types/:metric/start_date/:start_date/end_date/:end_date/interval/:interval.json
  * GET /api/1.1/cdns/peakusage/:peak_usage_type/deliveryservice/:ds/cachegroup/:name/start_date/:start/end_date/:end/interval/:interval.json
  * GET /api/1.1/deliveryservices/:id/:server_type/metrics/:metric_type/:start/:end.json

* :ref:`to-api-ds-server`

  * GET /api/1.1/deliveryserviceserver.json

* :ref:`to-api-ds-sslkeys`

  * GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys.json
  * GET /api/1.1/deliveryservices/hostname/#hostname/sslkeys.json
  * GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete.json
  * POST /api/1.1/deliveryservices/sslkeys/generate
  * POST /api/1.1/deliveryservices/sslkeys/add


:ref:`to-api-users`

* GET /api/1.1/users.json
* GET /api/1.1/user/current.json
* POST /api/1.1/user/current/update
* GET /api/1.1/user/current/jobs.json
* POST/api/1.1/user/current/jobs
* POST /api/1.1/user/login
* GET /api/1.1/user/:id/deliveryservices/available.json
* POST /api/1.1/user/login/token
* POST /api/1.1/user/logout
* POST /api/1.1/user/reset_password


:ref:`to-api-server`

* GET /api/1.1/servers.json
* GET /api/1.1/servers/summary.json
* GET /api/1.1/servers/hostname/:name/details.json
* POST /api/1.1/servercheck


:ref:`to-api-parameter`

* GET /api/1.1/parameters.json
* GET /api/1.1/parameters/:profile_name.json


:ref:`to-api-phys-loc`

* GET /api/1.1/phys_locations.json
* GET /api/1.1/phys_locations/trimmed.json


:ref:`to-api-profile`

* GET /api/1.1/profiles.json
* GET /api/1.1/profile/trimmed.json

.. :ref:`to-api-params`

 * GET /api/1.1/profileparameters.json

:ref:`to-api-region`

* GET /api/1.1/regions.json


:ref:`to-api-roles`

* GET /api/1.1/roles.json

:ref:`to-api-ext`

* GET /api/1.1/to_extensions.json
* POST /api/1.1/to_extensions
* POST /api/1.1/to_extensions/:id/delete


:ref:`to-api-status`

* GET /api/1.1/statuses.json


:ref:`to-api-dns`

* GET /api/1.1/staticdnsentries.json


:ref:`to-api-sys`

* GET /api/1.1/system/info.json

:ref:`to-api-redis`

* GET /api/1.1/traffic_monitor/stats.json
* GET /api/1.1/redis/stats.json
* GET /api/1.1/redis/info/:host_name.json
* GET /api/1.1/redis/match/#match/start_date/:start_date/end_date/:end_date/interval/:interval.json


:ref:`to-api-type`

* GET /api/1.1/types.json
* GET /api/1.1/types/trimmed.json

:ref:`to-api-error`

.. toctree:: 
  :hidden:
  :maxdepth: 1

  traffic_ops_rest/API1.2temp
  traffic_ops_rest/asn
  traffic_ops_rest/cachegroup
  traffic_ops_rest/cdn
  traffic_ops_rest/changelog
  traffic_ops_rest/deliveryservice
  traffic_ops_rest/Errors
  traffic_ops_rest/hwinfo
  traffic_ops_rest/parameter
  traffic_ops_rest/phys_location
  traffic_ops_rest/profile
  traffic_ops_rest/redis
  traffic_ops_rest/region
  traffic_ops_rest/role
  traffic_ops_rest/server
  traffic_ops_rest/static_dns
  traffic_ops_rest/status
  traffic_ops_rest/system
  traffic_ops_rest/to_extension
  traffic_ops_rest/type
  traffic_ops_rest/user

..  API Reference 1.2
..  =================