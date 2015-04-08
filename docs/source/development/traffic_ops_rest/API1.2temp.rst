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

API 1.2
=======

API 1.2 Roadmap
---------------

- Convert route parameters from unix time to ISO8601 standard
- Convert enpoints to the following URI convention:

GET /api/:version/scope/attribute/attribute-value/resource.json

Example: /api/1.2/deliveryservices/73/startDate/2015-03-01T00:00:00-07:00/endDate/2015-03-03T23:59:00-07:00/usage.json