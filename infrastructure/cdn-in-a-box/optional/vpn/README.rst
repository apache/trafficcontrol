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

**********************************
OpenVPN Container for CDN-In-A-Box
**********************************
This container provide an OpenVPN service.
It could let user or developer easily access CIAB network.

How to use it
=============
#. It is recommended that this be done using a custom bash alias.

    .. code-block:: shell

        # From infrastructure/cdn-in-a-box
        alias mydc="docker-compose -f $PWD/docker-compose.yml -f $PWD/optional/docker-compose.vpn.yml"
        mydc down -v
        mydc build
        mydc up

#. All certificates, keys, and client configuration are stored at ``infrastruture/cdn-in-a-box/optional/vpn/vpnca``. You just simply change ``REALHOSTIP`` and ``REALPORT`` of ``client.ovpn`` to fit your environment, and then you can connect to this OpenVPN server by it.

Private Subnet for Routing
==========================
Since ``docker-compose`` randomly create subnet, this container prepares 2 default private subnet for routing:

* 172.16.127.0/255.255.240.0
* 10.16.127.0/255.255.240.0

The strategy of choosing default private subnet is comparing the subnet prefix.
If the subnet prefix which ``docker-compose`` selected is 192. or 10.,
this container goes to select 172.16.127.0/255.255.240.0 for its routing subnet.
Otherwise, it selects 10.16.127.0/255.255.240.0.

Of course, you can decide which routing subnet subnet by supply environment
variable ``PRIVATE_NETWORK`` and ``PRIVATE_NETMASK``.

Pushed Settings
===============
Pushed settings are shown as follows:

* DNS
* A routing rule for CIAB subnet

.. note:: It will not change your default gateway.
