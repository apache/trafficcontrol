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

.. _to-api-v12-iso:

ISO
===

.. _to-api-v12-iso-route:

**GET /api/1.2/osversions**

  Get all OS versions for ISO generation and the directory where the kickstarter files are found. The values are retrieved from osversions.cfg found in either /var/www/files or in the location defined by the kickstart.files.location parameter (if defined).

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------------------------------------------------------------------------+
  | Parameter            | Description                                                              |
  +======================+==========================================================================+
  |``OS version name``   | OS version name. For example, "CentOS 7.2 vda".                          |
  +----------------------+--------------------------------------------------------------------------+
  |``OS version dir``    | The directory where the kickstarter ISO files are found. For example,    |
  |                      | centos72-netinstall.                                                     |
  +----------------------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response":
        {
           "CentOS 7.2": "centos72-netinstall"
           "CentOS 7.2 vda": "centos72-netinstall-vda"
        }
    }

|

**POST /api/1.2/isos**

  Generate an ISO.

  Authentication Required: Yes

  Role(s) Required: Operations

  **Request Properties**

  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | Parameter                     | Required | Description                                                                                     |
  +===============================+==========+=================================================================================================+
  | ``osversionDir``              | yes      | The directory name where the kickstarter ISO files are found.                                   |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``hostName``                  | yes      |                                                                                                 |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``domainName``                | yes      |                                                                                                 |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``rootPass``                  | yes      |                                                                                                 |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``dhcp``                      | yes      | Valid values are 'yes' or 'no'. If yes, other IP settings will be ignored.                      |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``interfaceMtu``              | yes      | 1500 or 9000                                                                                    |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``ipAddress``                 | yes|no   | Required if dhcp=no                                                                             |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``ipNetmask``                 | yes|no   | Required if dhcp=no                                                                             |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``ipGateway``                 | yes|no   | Required if dhcp=no                                                                             |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``ip6Address``                | no       | /64 is assumed if prefix is omitted.                                                            |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``ip6Gateway``                | no       | Ignored if an IPV4 gateway is specified.                                                        |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``interfaceName``             | no       | Typical values are bond0, eth4, etc. If you enter bond0, a LACP bonding config will be written. |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+
  | ``disk``                      | no       | Typical values are "sda"                                                                        |
  +-------------------------------+----------+-------------------------------------------------------------------------------------------------+

  **Request Example** ::

    {
        "osversionDir": "centos72-netinstall-vda",
        "hostName": "foo-bar",
        "domainName": "baz.com",
        "rootPass": "password",
        "dhcp": "no",
        "interfaceMtu": 1500,
        "ipAddress": "10.10.10.10",
        "ipNetmask": "255.255.255.252",
        "ipGateway": "10.10.10.10"
    }

|

  **Response Properties**

  +-----------------+--------+-------------------------------------------------------------------------------+
  | Parameter       | Type   | Description                                                                   |
  +=================+========+===============================================================================+
  |``isoURL``       | string | The URL location of the ISO. ISO locations can be found in cnd.conf file.     |
  +-----------------+--------+-------------------------------------------------------------------------------+

  **Response Example** ::

	{
		"response": {
			"isoURL": "https://traffic_ops.domain.net/iso/fqdn-centos72-netinstall.iso"
		},
		"alerts": [
			{
				"level": "success",
				"text": "Generate ISO was successful."
			}
		]
	}

|
