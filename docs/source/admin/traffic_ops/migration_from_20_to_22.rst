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

Traffic Ops - Migrating from 2.0 to 2.2
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

Per-DeliveryService Routing Names
---------------------------------
Before this release, DNS Delivery Services were hardcoded to use the name "edge", i.e. "edge.myds.mycdn.com", and HTTP Delivery Services use the name "tr" (or previously "ccr"), i.e. "tr.myds.mycdn.com". As of 2.2, Routing Names will default to "cdn" if left unspecified and can be set to any arbitrary non-dotted hostname.

Pre-2.2 the HTTP Routing Name is configurable via the http.routing.name option in in the Traffic Router http.properties config file. If your CDN uses that option to change the name from "tr" to something else, then you will need to perform the following steps for each CDN affected:

1. In Traffic Ops, create the following profile parameter (double-check for typos, trailing spaces, etc):

  * **name**: upgrade_http_routing_name
  * **config file**: temp
  * **value**: whatever value is used for the affected CDN's http.routing.name

2. Add this parameter to a single profile in the affected CDN

With those profile parameters in place Traffic Ops can be safely upgraded to 2.2. Before taking a post-upgrade snapshot, make sure to check your Delivery Service example URLs for unexpected Routing Name changes. Once Traffic Ops has been upgraded to 2.2 and a post-upgrade snapshot has been taken, your Traffic Routers can be upgraded to 2.2 (Traffic Routers must be upgraded after Traffic Ops so that they can work with custom per-DeliveryService Routing Names).

Apache Traffic Server 7.x (Cachekey Plugin)
-------------------------------------------
In Traffic Ops 2.2 we have added support for Apache Traffic Server 7.x. With 7.x comes support for the new cachekey plugin which replaces the cacheurl plugin which is now deprecated.  
While not needed immediately it is recommended to start replacing cacheurl usages with cachekey as soon as possible because ATS 6.x already supports the new cachekey plugin.

It is also recommended to thoroughly vet your cachekey replacement by comparing with an existing key value. There are inconsistencies in the 6.x version of cachekey which have been
fixed in 7.x (or require this patch(`cachekeypatch`_) on 6.x to match 7.x). So to ensure you have a matching key value you should use the xdebug plugin before fully implementing your cachekey replacement.

.. _cachekeypatch: https://github.com/apache/trafficserver/commit/244288fab01bdad823f9de19dcece62a7e2a0c11

First if you are currently using a regex for your delivery service you will have to remove that existing value. Then you will need to make a new DS profile and assign parameters in
it to the cachekey.config file.

Some common parameters are

.. highlight:: none

::

    static-prefix      - This is used for a simple domain replacement
    separator          - Used by cachekey and in general is always a single space
    remove-path        - Removes path information from the URL
    remove-all-params  - Removes parameters from the URL
    capture-prefix-uri - This is usually used in combination with remove-path and remove-all-params. 
                         Capture-prefix-uri will let you use your own full regex value for non simple cases

Examples of Cacheurl to Cachekey Replacements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**Original regex value:**

``http://test.net/(.*) http://test-cdn.net/$1``

**Cachekey parameters:**

+---------------+-----------------+---------------------------------+
| Parameter     |  File           |   Value                         |
+===============+=================+=================================+
| static-prefix | cachekey.config | ``http://test-cdn.net/``        |
+---------------+-----------------+---------------------------------+
| separator     | cachekey.config |   (empty space)                 |
+---------------+-----------------+---------------------------------+

**Original regex value:**

``http://([^?]+)(?:?|$) http://test-cdn.net/$1``

**Cachekey parameters:**

+-----------------------+-----------------+-----------------------------------------------------+
| Parameter             |  File           |   Value                                             |
+=======================+=================+=====================================================+
| remove-path           | cachekey.config | true                                                |
+-----------------------+-----------------+-----------------------------------------------------+
| remove-all-params     | cachekey.config |   true                                              |
+-----------------------+-----------------+-----------------------------------------------------+
| separator             | cachekey.config |    (empty space)                                    |
+-----------------------+-----------------+-----------------------------------------------------+
| capture-prefix-uri    | cachekey.config |  ``/https?:\/\/([^?]*)/http:\/\/test-cdn.net\/$1/`` |
+-----------------------+-----------------+-----------------------------------------------------+

Also note the ``s?`` used here so that both http and https requests will end up with the same key value

**Original regex value:**

``http://test.net/([^?]+)(?:\?|$) http://test-cdn.net/$1``

**Cachekey parameters:**

+-------------------+-----------------+---------------------------------+
| Parameter         |  File           |   Value                         |
+===================+=================+=================================+
| static-prefix     | cachekey.config | ``http://test-cdn.net/``        |
+-------------------+-----------------+---------------------------------+
| separator         | cachekey.config |   (empty space)                 |
+-------------------+-----------------+---------------------------------+
| remove-all-params | cachekey.config | true                            |
+-------------------+-----------------+---------------------------------+

.. _ApacheTrafficServerDocs: https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/cachekey.en.html

Further documentation on the cachekey plugin can be found at `ApacheTrafficServerDocs`_

Apache Traffic Server 7.x (Logging)
-------------------------------------------
Trafficserver v7 has changed the logging format.  Previously this was an xml file and now it is a lua file. Traffic Control compensates for this
automatically depending upon the filename used for the logging parameters.  Previously the file used was ``logs_xml.config``, for ATS 7 it is now
``logging.config``.  The same parameters will work this new file, ``LogFormat.Format``, ``LogFormat.Name``, ``LogObject.Format`` etc.


Traffic Ops Profile Modifications
-------------------------------------------
When upgrading to ATS 7.x, the Traffic Ops EDGE and MID cache profiles must be modified to provide new configuration values. Trafficserver recommends changes to the following parameters: https://cwiki.apache.org/confluence/display/TS/Upgrading+to+v7.0

Most users of Traffic Control have enough profiles where making these modifications manually is a tedious and time consuming process. A new utility `traffic_ops/install/bin/convert_profile/convert_profile` is provided to automatically convert an ATS6 profile to an ATS7 profile. This utility can be reused in the future for ATS7->8 as well. 

Usage example:
  1. Use Traffic Portal GUI to export profile to JSON
  2. Modify the trafficserver version numbers to match your current Traffic Server v6 RPM version and planned Traffic Server v7 RPM version
  3. Run ``convert_profile -input_profile <exported_file> -rules convert622to713.json -out <new_profile_name>``
  4. Review output messages and make manual updates if needed. If you have modified a default value the script also wants to change, it will prompt you to make the update manually. You may either do this directly in the JSON file or through the Traffic Portal GUI after import. 
  5. Use Traffic Portal GUI to import the newly created profile

