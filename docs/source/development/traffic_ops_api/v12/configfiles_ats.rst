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


.. _to-api-v12-configfiles_ats:

Config Files and Config File Metadata
=====================================

.. _to-api-v12-configfiles_ats-route:

/api/1.2/servers/:hostname/configfiles/ats
++++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/servers/:hostname/configfiles/ats**

  Authentication Required: Yes

  Role(s) Required: Operator

  **Request Query Parameters**

  **Response Properties**

  +-------------------+--------+-------------------------------------------------------------------------+
  |                   |        |           Info Section                                                  |
  +-------------------+--------+-------------------------------------------------------------------------+
  |    Parameter      |  Type  |                               Description                               |
  +===================+========+=========================================================================+
  | ``profileId``     |  int   | The ID of the profile assigned to the cache.                            |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``profileName``   | string | The name of the profile assigned to the cache.                          |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``toRevProxyUrl`` | string | The configured reverse proxy cache for configfile requests.             |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``toURL``         | string | The configured URL for Traffic Ops.                                     |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``serverIpv4``    | string | The configured IP address of the cache.                                 |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``serverName``    | string | The cache's short form hostname.                                        |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``serverId``      |  int   | The cache's Traffic Ops ID.                                             |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``cdnId``         |  int   | The ID of the cache's assigned CDN.                                     |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``cdnName``       | string | The name of the cache's assigned CDN.                                   |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``serverTcpPort`` |  int   | The configured port of the server's used by ATS.                        |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``fnameOnDisk``   | string | The filename of the configuration file as stored on the cache.          |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``location``      | string | The directory location of the configuration file as stored on the cache.|
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``apiUri``        | string | The path to generate the configuration file from Traffic Ops.           |
  +-------------------+--------+-------------------------------------------------------------------------+
  | ``scope``         | string | The scope of the configuration file.                                    |
  +-------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
      "info": {
        "profileId": 278,
        "toRevProxyUrl": "https://to.example.com:81",
        "toUrl": "https://to.example.com/",
        "serverIpv4": "192.168.1.5",
        "serverTcpPort": 80,
        "serverName": "cache-ats-01",
        "cdnId": 1,
        "cdnName": "CDN_1",
        "serverId": 21,
        "profileName": "EDGE_CDN_1_EXAMPLE"
      },
      "configFiles": [
        {
          "fnameOnDisk": "remap.config",
          "location": "/opt/trafficserver/etc/trafficserver",
          "apiUri": "/api/1.2/profiles/EDGE_CDN_1_EXAMPLE/configfiles/ats/remap.config",
          "scope": "profiles"
        },
        {
          "fnameOnDisk": "ssl_multicert.config",
          "location": "/opt/trafficserver/etc/trafficserver",
          "apiUri": "/api/1.2/cdns/CDN_1/configfiles/ats/ssl_multicert.config",
          "scope": "cdns"
        },
        {
          "fnameOnDisk": "parent.config",
          "location": "/opt/trafficserver/etc/trafficserver",
          "apiUri": "/api/1.2/servers/cache-ats-01/configfiles/ats/parent.config"
        }
      ]
    }


/api/1.2/servers/:hostname/configfiles/ats/:configfile
++++++++++++++++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/servers/:hostname/configfiles/ats/:configfile**
**GET /api/1.2/servers/:host_id/configfiles/ats/:configfile**


  Authentication Required: Yes

  Role(s) Required: Operator

  **Request Query Parameters**

  **Response Properties**

  Returns the requested configuration file for download.  If scope used is incorrect for the config file requested, returns a 404 with the correct scope.

  **Response Example** ::

    {
      "alerts": [
        {
          "level": "error",
          "text": "Error - incorrect file scope for route used.  Please use the profiles route."
        }
      ]
    }


/api/1.2/profiles/:profile_name/configfiles/ats/:configfile
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/profiles/:profile_name/configfiles/ats/:configfile**
**GET /api/1.2/profiles/:profile_id/configfiles/ats/:configfile**


  Authentication Required: Yes

  Role(s) Required: Operator

  **Request Query Parameters**

  **Response Properties**

  Returns the requested configuration file for download.  If scope used is incorrect for the config file requested, returns a 404 with the correct scope.

  **Response Example** ::

    {
      "alerts": [
        {
          "level": "error",
          "text": "Error - incorrect file scope for route used.  Please use the cdns route."
        }
      ]
    }


/api/1.2/cdns/:cdn_name/configfiles/ats/:configfile
+++++++++++++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/cdns/:cdn_name/configfiles/ats/:configfile**
**GET /api/1.2/cdns/:cdn_id/configfiles/ats/:configfile**


  Authentication Required: Yes

  Role(s) Required: Operator

  **Request Query Parameters**

  **Response Properties**

  Returns the requested configuration file for download.  If scope used is incorrect for the config file requested, returns a 404 with the correct scope.

  **Response Example** ::

    {
      "alerts": [
        {
          "level": "error",
          "text": "Error - incorrect file scope for route used.  Please use the servers route."
        }
      ]
    }

