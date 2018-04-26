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

.. _to-api-v11-dns:

Static DNS Entries
==================

.. _to-api-v11-static-dns-route:

/api/1.1/staticdnsentries
+++++++++++++++++++++++++

**GET /api/1.1/staticdnsentries.json**

    Authentication Required: Yes

    Role(s) Required: None

    **Response Properties**

    +---------------------+-----------+------------------------------------------------------------+
    | Parameter           |  Type     |                             Description                    |
    +=====================+===========+============================================================+
    | ``deliveryservice`` | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+
    | ``ttl``             | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+
    | ``type``            | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+
    | ``address``         | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+
    | ``cachegroup``      | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+
    | ``host``            | string    |                                                            |
    +---------------------+-----------+------------------------------------------------------------+

    **Response Example** ::

       {
        "response": [
          {
            "deliveryservice": "foo-ds",
            "ttl": "30",
            "type": "CNAME_RECORD",
            "address": "bar.foo.baz.tv.",
            "cachegroup": "us-co-denver",
            "host": "mm"
          }
        ]
      }
