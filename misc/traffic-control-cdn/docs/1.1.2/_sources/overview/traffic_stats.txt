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

.. _reference-label-tc-ts:

Traffic Stats
=============
Traffic Stats is a utility written in Go that mines metrics from Traffic Monitor's JSON APIs and stores the data locally in Redis for a short period of time. This data is inherently transient, rolls frequently, and is volatile due to the default in-memory nature of Redis. The transient nature of the data is acceptable, as this system's purpose is to land data in Redis for other tools to consume.

Once in Redis, the data can be extracted and prepared to be sent elsewhere for long term storage. Any number of Traffic Stats instances may run on a given CDN to collect metrics from Traffic Monitor, however, redundancy and integration with a long term metrics storage system is implementation dependent. Traffic Stats does not influence overall CDN operation, but is required in order to display charts in Traffic Operations.