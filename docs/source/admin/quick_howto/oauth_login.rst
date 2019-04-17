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
.. _oauth-login:

*********************
Configure OAuth Login
*********************

An opt-in configuration for SSO using OAuth is supported and can be configured through the :file:`/opt/traffic_portal/public/traffic_portal_properties.json` and :file:`/opt/traffic_ops/app/conf/cdn.conf` files.
OAuth uses a third party provider to authenticate the user. Once enabled, the Traffic Portal Login page will no longer accept username and password but instead will authenticate using OAuth.
This will redirect to the ``oAuthUrl`` from :file:`/opt/traffic_portal/public/traffic_portal_properties.json` which will authenticate the user then redirect to the new ``/sso`` page with a Json Web Token added as a query parameter.
The new ``/sso`` page will parse the token from the URL and post this information to the ``/api/1.4/user/login/oauth`` API endpoint. See :ref:`to-api-user-login-oauth`.
The ``/api/1.4/user/login/oauth`` API endpoint will decode the token, validate that it is between the issued time and the expiration time, and validate that the public key set URL is allowed by the list of whitelisted URLs read from :file:`/opt/traffic_ops/app/conf/cdn.conf`. It will then authorize the user from the database and return a mojolicious cookie as per the normal login workflow.

.. Note:: OAuth providers sometimes do not return the public key set URL but instead require a locally stored key. This functionality is not currently supported and will require further development.

To configure OAuth login:

Set up authentication with a third party OAuth provider.

- Update :file:`/opt/traffic_portal/public/traffic_portal_properties.json` and ensure the following properties are set up correctly:

.. table:: OAuth Configuration Property Definitions In traffic_portal_properties.json

	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+
	| Name                     | Type       | Description                                                                                                   |
	+==========================+============+===============================================================================================================+
	| enabled                  | boolean    | Allow OAuth SSO login                                                                                         |
	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+
	| oAuthUrl                 | string     | URL to your OAuth provider                                                                                    |
	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+
	| oAuthTokenQueryParam     | string     | Query parameter containing token from OAuth provider (returned in URL when redirected to ``/sso`` endpoint)   |
	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+


.. code-block:: text
    :caption: Example OAuth Configuration Properties In traffic_portal_properties.json

    "oAuth": {
        "enabled": true,
        "oAuthUrl": "example.oauth.com",
        "oAuthTokenQueryParam": "access_token",
    }

- Update :file:`/opt/traffic_ops/app/conf/cdn.conf` property traffic_ops_golang.whitelisted_oauth_urls to contain all allowed domains for the Json key set (Use `*` for wildcard):

.. table:: OAuth Configuration Property Definitions In cdn.conf

	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+
	| Name                     | Type       | Description                                                                                                   |
	+==========================+============+===============================================================================================================+
	| whitelisted_oauth_urls   | []string   | List of whitelisted URLs for the Json public key set returned by OAuth provider.  Can contain * wildcards.    |
	+--------------------------+------------+---------------------------------------------------------------------------------------------------------------+


.. code-block:: text
    :caption: Example OAuth Configuration Properties In cdn.conf

    "traffic_ops_golang": {
        ...
        "whitelisted_oauth_urls": [
            "example.oauth.com",
            "*.oauth.com"
        ]
    }