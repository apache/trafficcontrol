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

.. _federations-qht:

*********************
Configure Federations
*********************

#. Create a user with a federations role (:menuselection:`User Admin --> Users --> '+' button`). This user will need the ability to:

	- Create/edit/delete federations
	- Add IPV4 resolvers
	- Add IPV6 resolvers

	.. figure:: federations/01.png
		:scale: 100%
		:align: center

#. As a user with administrative privileges, create a Federation Mapping by going to :menuselection:`Services --> :term:`Delivery Services` --> More --> Federations` and then clicking :guilabel:`Add Federation Mapping`.

#. Choose the :term:`Delivery Service` to which the federation will be mapped and assign it to the Federation-role user; click :guilabel:`Add`.

	.. figure:: federations/02.png
		:scale: 100%
		:align: center

#. After the Federation is added, Traffic Ops will display the Federation. Changes can be made at this time or the Federation can be deleted. Notice that no resolvers have been added to the Federation yet. This can only be done by the Federation-role user to whom the Federated :term:`Delivery Service` was assigned. If no further action is necessary, the :guilabel:`Close` button will close the window and display the list of all Federations.

	.. figure:: federations/03.png
		:scale: 100%
		:align: center

#. The federation user logs into either the Traffic Ops API or the Traffic Portal UI and stores the Mojolicious cookie. The Mojolicious cookie can be obtained manually using the debug tools on a web browser or via a command line utility like :manpage:`curl(1)`.

	.. code-block:: shell
		:caption: Example cURL Command

		curl -i -XPOST "http://localhost:3000/api/4.0/user/login" -H "Content-Type: application/json" -d '{ "u": "federation_user1", "p": "password" }'

	.. code-block:: http
		:caption: Example API Response

		HTTP/1.1 200 OK
		Date: Wed, 02 Dec 2015 21:12:06 GMT
		Content-Length: 65
		Access-Control-Allow-Credentials: true
		Content-Type: application/json
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
		X-Server-Name: traffic_ops_golang/
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
		Cache-Control: no-cache, no-store, max-age=0, must-revalidate
		Connection: keep-alive
		Access-Control-Allow-Origin: http://localhost:8080

		{"alerts":[{"level":"success","text":"Successfully logged in."}]}

#. The federation user sends a request to Traffic Ops to add IPV4 and/or IPV6 resolvers


	.. code-block:: shell
		:caption: Example cURL Command

		curl -ki -H "Cookie: mojolicious=eyJleHBpcmVzIjoxNDQ5MTA1MTI2LCJhdXRoX2RhdGEiOiJmZWRlcmF0aW9uX3VzZXIxIn0---06b4f870d809d82a91433e92eae8320875c3e8b0;" -XPUT 'http://localhost:3000/api/4.0/federations' -d '
		{"federations": [
			{ "deliveryService": "images-c1",
				"mappings":
				{ "resolve4": [ "8.8.8.8/32", "8.8.4.4/32" ],
					"resolve6": ["2001:4860:4860::8888/128", "2001:4860:4860::8844"]
				}
			}
		]}'

	.. code-block:: http
		:caption: Example API Response

		HTTP/1.1 200 OK
		Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
		X-Server-Name: traffic_ops_golang/
		Date: Wed, 02 Dec 2015 21:25:42 GMT
		Content-Length: 74
		Access-Control-Allow-Credentials: true
		Content-Type: application/json
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Cache-Control: no-cache, no-store, max-age=0, must-revalidate
		Access-Control-Allow-Origin: http://localhost:8080
		Connection: keep-alive
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept

		{"response":"federation_user1 successfully created federation resolvers."}

#. The resolvers added by the Federation-user will now visible in Traffic Portal.

	.. figure:: federations/04.png
		:scale: 100%
		:align: center

Any requests made from a client that resolves to one of the federation resolvers will now be given a :abbr:`CNAME (Canonical Name)` Record from Traffic Router.

	.. code-block:: shell
		:caption: Example DNS request (via ``dig``)

		dig @tr.kabletown.net foo.images-c1.kabletown.net

	.. code-block:: DNS
		:caption: Example Resolver Response

		; <<>> DiG 9.7.3-RedHat-9.7.3-2.el6 <<>> @tr.kabletown.net foo.images-c1.kabletown.net
		; (1 server found)
		;; global options: +cmd
		;; Got answer:
		;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 45110
		;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
		;; WARNING: recursion requested but not available

		;; QUESTION SECTION:
		;foo.images-c1.kabletown.net.	IN A

		;; ANSWER SECTION:
		foo.images-c1.kabletown.net.	30 IN CNAME img.mega-cdn.net.

		;; Query time: 9 msec
		;; SERVER: 10.10.10.10#53(10.10.10.10)
		;; WHEN: Wed Dec  2 22:05:26 2015
		;; MSG SIZE  rcvd: 84
