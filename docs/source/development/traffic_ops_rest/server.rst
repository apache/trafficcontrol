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

.. _to-api-server:

Server
======
**GET /api/1.1/servers.json**

  Retrieves properties of CDN servers.

  Authentication Required: Yes

  **Response Properties**

  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  |     Parameter      |  Type  |                                                Description                                                 |
  +====================+========+============================================================================================================+
  | ``cachegroup``     | string | The cache group name (see :ref:`to-api-cachegroup`).                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``domainName``     | string | The domain name part of the FQDN of the cache.                                                             |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``hostName``       | string | The host name part of the cache.                                                                           |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``id``             | string | The server id (database row number).                                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpAddress``   | string | The IPv4 address of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpGateway``   | string | The IPv4 gateway address of the lights-out-management port.                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloIpNetmask``   | string | The IPv4 netmask of the lights-out-management port.                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloPassword``    | string | The password of the of the lights-out-management user (displays as ****** unless you are an 'admin' user). |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``iloUsername``    | string | The user name for lights-out-management.                                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``   | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``interfaceName``  | string | The network interface name used for serving traffic.                                                       |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Address``     | string | The IPv6 address/netmask for ``interfaceName``.                                                            |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``     | string | The IPv6 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipAddress``      | string | The IPv4 address for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipGateway``      | string | The IPv4 gateway for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``ipNetmask``      | string | The IPv4 netmask for ``interfaceName``.                                                                    |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``    | string | The Time and Date for the last update for this server.                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpAddress``  | string | The IPv4 address of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpGateway``  | string | The IPv4 gateway of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpNetmask``  | string | The IPv4 netmask of the management port (optional).                                                        |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``physLocation``   | string | The physical location name (see :ref:`to-api-phys-loc`).                                                   |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``profile``        | string | The assigned profile name (see :ref:`to-api-profile`).                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``rack``           | string | A string indicating rack location.                                                                         |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerHostName`` | string | The human readable name of the router.                                                                     |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``routerPortName`` | string | The human readable name of the router port.                                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``status``         | string | The Status string (See :ref:`to-api-status`).                                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``tcpPort``        | string | The default TCP port on which the main application listens (80 for a cache in most cases).                 |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``type``           | string | The name of the type of this server (see :ref:`to-api-type`).                                              |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``xmppId``         | string | Deprecated.                                                                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+
  | ``xmppPasswd``     | string | Deprecated.                                                                                                |
  +--------------------+--------+------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

   {
      "response": [
          {
              "cachegroup": "us-il-chicago",
              "domainName": "chi.kabletown.net",
              "hostName": "atsec-chi-00",
              "id": "19",
              "iloIpAddress": "172.16.2.6",
              "iloIpGateway": "172.16.2.1",
              "iloIpNetmask": "255.255.255.0",
              "iloPassword": "********",
              "iloUsername": "",
              "interfaceMtu": "9000",
              "interfaceName": "bond0",
              "ip6Address": "2033:D0D0:3300::2:2/64",
              "ip6Gateway": "2033:D0D0:3300::2:1",
              "ipAddress": "10.10.2.2",
              "ipGateway": "10.10.2.1",
              "ipNetmask": "255.255.255.0",
              "lastUpdated": "2015-03-08 15:57:32",
              "mgmtIpAddress": "",
              "mgmtIpGateway": "",
              "mgmtIpNetmask": "",
              "physLocation": "plocation-chi-1",
              "profile": "EDGE1_CDN1_421_SSL",
              "rack": "RR 119.02",
              "routerHostName": "rtr-chi.kabletown.net",
              "routerPortName": "2",
              "status": "ONLINE",
              "tcpPort": "80",
              "type": "EDGE",
              "xmppId": "atsec-chi-00-dummyxmpp",
              "xmppPasswd": "**********"
          },
          {
          ... more server data
          }
        ]
      "version": "1.1"
    }


|

**GET /api/1.1/servers/summary.json**

  Retrieves a count of CDN servers by type.

  Authentication Required: Yes

  **Response Properties**

  +-----------+--------+---------------------------------------------------------------------+
  | Parameter |  Type  |                             Description                             |
  +===========+========+=====================================================================+
  | ``count`` | int    | The number of servers of this type in this instance of Traffic Ops. |
  +-----------+--------+---------------------------------------------------------------------+
  | ``type``  | string | The name of the type of the server count (see :ref:`to-api-type`).  |
  +-----------+--------+---------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
          "count": 4,
          "type": "CCR"
        },
        {
          "count": 55,
          "type": "EDGE"
        },
        {
          "type": "MID",
          "count": 18
        },
        {
          "count": 0,
          "type": "REDIS"
        },
        {
          "count": 4,
          "type": "RASCAL"
        }
      "version": "1.1",
    }

|

**GET /api/1.1/servers/hostname/:name/details.json**

  Retrieves the details of a server.

  Authentication Required: Yes

  **Request Route Parameters**

  +----------+----------+----------------------------------+
  |   Name   | Required |           Description            |
  +==========+==========+==================================+
  | ``name`` | yes      | The host name part of the cache. |
  +----------+----------+----------------------------------+

  **Response Properties**

  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  |      Parameter       |  Type  |                                                 Description                                                 |
  +======================+========+=============================================================================================================+
  | ``cachegroup``       | string | The cache group name (see :ref:`to-api-cachegroup`).                                                        |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``deliveryservices`` | array  | Array of strings with the delivery service ids assigned (see :ref:`to-api-ds`).                             |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``domainName``       | string | The domain name part of the FQDN of the cache.                                                              |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``hardwareInfo``     | hash   | Hwinfo struct (see :ref:`to-api-hwinfo`).                                                                   |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``hostName``         | string | The host name part of the cache.                                                                            |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``id``               | string | The server id (database row number).                                                                        |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpAddress``     | string | The IPv4 address of the lights-out-management port.                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpGateway``     | string | The IPv4 gateway address of the lights-out-management port.                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloIpNetmask``     | string | The IPv4 netmask of the lights-out-management port.                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloPassword``      | string | The password of the of the lights-out-management user  (displays as ****** unless you are an 'admin' user). |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``iloUsername``      | string | The user name for lights-out-management.                                                                    |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``     | string | The Maximum Transmission Unit (MTU) to configure for ``interfaceName``.                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``interfaceName``    | string | The network interface name used for serving traffic.                                                        |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ip6Address``       | string | The IPv6 address/netmask for ``interfaceName``.                                                             |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``       | string | The IPv6 gateway for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipAddress``        | string | The IPv4 address for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipGateway``        | string | The IPv4 gateway for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``ipNetmask``        | string | The IPv4 netmask for ``interfaceName``.                                                                     |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``lastUpdated``      | string | The Time/Date of the last update for this server.                                                           |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpAddress``    | string | The IPv4 address of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpGateway``    | string | The IPv4 gateway of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``mgmtIpNetmask``    | string | The IPv4 netmask of the management port (optional).                                                         |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``physLocation``     | string | The physical location name (see :ref:`to-api-phys-loc`).                                                    |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``profile``          | string | The assigned profile name (see :ref:`to-api-profile`).                                                      |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``rack``             | string | A string indicating rack location.                                                                          |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``routerHostName``   | string | The human readable name of the router.                                                                      |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``routerPortName``   | string | The human readable name of the router port.                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``status``           | string | The Status string (See :ref:`to-api-status`).                                                               |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``tcpPort``          | string | The default TCP port on which the main application listens (80 for a cache in most cases).                  |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``type``             | string | The name of the type of this server (see :ref:`to-api-type`).                                               |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``xmppId``           | string | Deprecated.                                                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+
  | ``xmppPasswd``       | string | Deprecated.                                                                                                 |
  +----------------------+--------+-------------------------------------------------------------------------------------------------------------+

  **Response Example** ::
   
    {
      "response": {
        "cachegroup": "us-il-chicago",
        "deliveryservices": [
          "1",
          "2",
          "3",
          "4"
        ],
        "domainName": "chi.kabletown.net",
        "hardwareInfo": {
          "Physical Disk 0:1:3": "D1S2",
          "Physical Disk 0:1:2": "D1S2",
          "Physical Disk 0:1:15": "D1S2",
          "Power Supply.Slot.2": "04.07.15",
          "Physical Disk 0:1:24": "YS08",
          "Physical Disk 0:1:1": "D1S2",
          "Model": "PowerEdge R720xd",
          "Physical Disk 0:1:22": "D1S2",
          "Physical Disk 0:1:18": "D1S2",
          "Enterprise UEFI Diagnostics": "4217A5",
          "Lifecycle Controller": "1.0.8.42",
          "Physical Disk 0:1:8": "D1S2",
          "Manufacturer": "Dell Inc.",
          "Physical Disk 0:1:6": "D1S2",
          "SysMemTotalSize": "196608",
          "PopulatedDIMMSlots": "24",
          "Physical Disk 0:1:20": "D1S2",
          "Intel(R) Ethernet 10G 2P X520 Adapter": "13.5.7",
          "Physical Disk 0:1:14": "D1S2",
          "BACKPLANE FIRMWARE": "1.00",
          "Dell OS Drivers Pack, 7.0.0.29, A00": "7.0.0.29",
          "Integrated Dell Remote Access Controller": "1.57.57",
          "Physical Disk 0:1:5": "D1S2",
          "ServiceTag": "D6XPDV1",
          "PowerState": "2",
          "Physical Disk 0:1:23": "D1S2",
          "Physical Disk 0:1:25": "D903",
          "BIOS": "1.3.6",
          "Physical Disk 0:1:12": "D1S2",
          "System CPLD": "1.0.3",
          "Physical Disk 0:1:4": "D1S2",
          "Physical Disk 0:1:0": "D1S2",
          "Power Supply.Slot.1": "04.07.15",
          "PERC H710P Mini": "21.0.2-0001",
          "PowerCap": "689",
          "Physical Disk 0:1:16": "D1S2",
          "Physical Disk 0:1:10": "D1S2",
          "Physical Disk 0:1:11": "D1S2",
          "Lifecycle Controller 2": "1.0.8.42",
          "BP12G+EXP 0:1": "1.07",
          "Physical Disk 0:1:9": "D1S2",
          "Physical Disk 0:1:17": "D1S2",
          "Broadcom Gigabit Ethernet BCM5720": "7.2.20",
          "Physical Disk 0:1:21": "D1S2",
          "Physical Disk 0:1:13": "D1S2",
          "Physical Disk 0:1:7": "D1S2",
          "Physical Disk 0:1:19": "D1S2"
        },
        "hostName": "atsec-chi-00",
        "id": "19",
        "iloIpAddress": "172.16.2.6",
        "iloIpGateway": "172.16.2.1",
        "iloIpNetmask": "255.255.255.0",
        "iloPassword": "********",
        "iloUsername": "",
        "interfaceMtu": "9000",
        "interfaceName": "bond0",
        "ip6Address": "2033:D0D0:3300::2:2/64",
        "ip6Gateway": "2033:D0D0:3300::2:1",
        "ipAddress": "10.10.2.2",
        "ipGateway": "10.10.2.1",
        "ipNetmask": "255.255.255.0",
        "mgmtIpAddress": "",
        "mgmtIpGateway": "",
        "mgmtIpNetmask": "",
        "physLocation": "plocation-chi-1",
        "profile": "EDGE1_CDN1_421_SSL",
        "rack": "RR 119.02",
        "routerHostName": "rtr-chi.kabletown.net",
        "routerPortName": "2",
        "status": "ONLINE",
        "tcpPort": "80",
        "type": "EDGE",
        "xmppId": "atsec-chi-00-dummyxmpp",
        "xmppPasswd": "X"

      }
      "version": "1.1",

    }

|

**POST /api/1.1/servercheck**

  Post a server check result to the serverchecks table.

  Authentication Required: Yes

  **Request Route Parameters**

  +----------------------------+----------+-------------+
  |            Name            | Required | Description |
  +============================+==========+=============+
  | ``id``                     | yes      |             |
  +----------------------------+----------+-------------+
  | ``host_name``              | yes      |             |
  +----------------------------+----------+-------------+
  | ``servercheck_short_name`` | yes      |             |
  +----------------------------+----------+-------------+
  | ``value``                  | yes      |             |
  +----------------------------+----------+-------------+

  **Request Example** ::


    {
     "id": "",
     "host_name": "",
     "servercheck_short_name": "",
     "value": ""
    }

  Response Content Type: application/json

  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+
  | ``version`` | string |                                  |
  +-------------+--------+----------------------------------+
   
   
  **Response Example** ::


    Response Example:

    {
      "alerts":
        [
          { 
            "level": "success",
            "text": "Server Check was successfully updated."
          }
        ],
      "version": "1.1"
    }

