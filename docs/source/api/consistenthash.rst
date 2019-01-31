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

.. _to-api-consistenthash:

***************
``consistenthash``
***************
Test pattern based consistent hashing for a Delivery Service using a regex and a request path

``POST``
=======
Queries database for an active Traffic Router on a given CDN and sends GET request to get the resulting path to consistent hash with a given regex and request path.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
:regex:       The regex to apply to the request path to get a resulting path that will be used for consistent hashing
:requestpath: The request path to use to test the regex against
:cdnid:       The unique identifier of a CDN that will be used to query for an active Traffic Router

Response Structure
------------------
:resultingPathToConsistentHash: The resulting path that Traffic Router will use for consistent hashing
:consistentHashRegex:           The regex used by Traffic Router derived from POST 'regex' parameter
:requestPath:                   The request path used by Traffic Router to test regex against

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"resultingPathToConsistentHash": "/path/mpd",
		"consistentHashRegex": "/.*?(/.*?/).*?(mpd)",
		"requestPath": "/test/path/asset.mpd"
	}}
