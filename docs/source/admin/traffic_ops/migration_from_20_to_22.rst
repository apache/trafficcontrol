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