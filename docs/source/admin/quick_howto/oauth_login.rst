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
.. _oauth_login:

*********************
Configure OAuth Login
*********************

An opt-in configuration for SSO using OAuth is supported and can be configured through the :file:`/opt/traffic_portal/public/traffic_portal_properties.json` and :file:`/opt/traffic_ops/app/conf/cdn.conf` files. OAuth uses a third party provider to authenticate the user. If enabled, the Traffic Portal Login page will still accept username and password for local accounts but will also allow authentication using OAuth. This will redirect to the ``oAuthUrl`` from :file:`/opt/traffic_portal/public/traffic_portal_properties.json` which will authenticate the user then redirect to the new ``/sso`` page with an authorization code. The new ``/sso`` page will then construct the full URL to exchange the authorization code for a JSON Web Token, and ``POST`` this information to the :ref:`to-api-user-login-oauth` API endpoint. The :ref:`to-api-user-login-oauth` API endpoint will ``POST`` to the URL provided and receive JSON Web Token. The :ref:`to-api-user-login-oauth` API endpoint will decode the token, validate that it is between the issued time and the expiration time, and validate that the public key set URL is allowed by the list of whitelisted URLs read from :file:`/opt/traffic_ops/app/conf/cdn.conf`. It will then authorize the user from the database and return a mojolicious cookie as per the normal login workflow.

.. Note:: Ensure that the user names in the Traffic Ops database match the value returned in the `sub` field in the response from the OAuth provider when setting up with the OAuth provider.  The `sub` field is used to reference the roles in the Traffic Ops database in order to authorize the user.

.. Note:: OAuth providers sometimes do not return the public key set URL but instead require a locally stored key. This functionality is not currently supported and will require further development.

.. Note:: The ``POST`` from the API to the OAuth provider to exchange the code for a token expects the response to have the token in JSON format with `access_token` as the desired field (and can include other fields).  It also supports a response with just the token itself as the body.  Further development work will need to be done to allow other resposne forms or other response fields.

.. Note:: Users must exist in both Traffic Ops as well as in the OAuth provider's system.  The user's rights are defined by the :term:`Role` assigned to the user.

To configure OAuth login:

- Set up authentication with a third party OAuth provider.

- Update :file:`/opt/traffic_portal/public/traffic_portal_properties.json` and ensure the following properties are set up correctly:

	.. table:: OAuth Configuration Property Definitions In traffic_portal_properties.json

		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+
		| Name                         | Type       | Description                                                                                                                               |
		+==============================+============+===========================================================================================================================================+
		| enabled                      | boolean    | Allow OAuth SSO login                                                                                                                     |
		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+
		| oAuthUrl                     | string     | URL to your OAuth provider                                                                                                                |
		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+
		| redirectUriParameterOverride | string     | Query parameter override if the oAuth provider requires a different key for the redirect_uri parameter, defaults to ``redirect_uri``      |
		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+
		| clientId                     | string     | Client id registered with OAuth provider, passed in with `client_id` parameter                                                            |
		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+
		| oAuthCodeTokenUrl            | string     | URL to your OAuth provider's endpoint for exchanging the code (from oAuthUrl) for a token                                                 |
		+------------------------------+------------+-------------------------------------------------------------------------------------------------------------------------------------------+


	.. code-block:: json
		:caption: Example OAuth Configuration Properties In traffic_portal_properties.json

		{
			"oAuth": {
				"_comment": "Opt-in OAuth properties for SSO login. See http://traffic-control-cdn.readthedocs.io/en/release-4.0.0/admin/quick_howto/oauth_login.html for more details. redirectUriParameterOverride defaults to redirect_uri if left blank.",
				"enabled": true,
				"oAuthUrl": "example.oauth.com",
				"redirectUriParameterOverride": "",
				"clientId": "",
				"oAuthCodeTokenUrl": "example.oauth.com/oauth/token"
			}
		}

- Update :file:`/opt/traffic_ops/app/conf/cdn.conf` property traffic_ops_golang.whitelisted_oauth_urls to contain all allowed domains for the JSON key set (Use ``*`` for wildcard):

	.. table:: OAuth Configuration Property Definitions In cdn.conf

		+--------------------------+--------------------+---------------------------------------------------------------------------------------------------------------------+
		| Name                     | Type               | Description                                                                                                         |
		+==========================+====================+=====================================================================================================================+
		| whitelisted_oauth_urls   | Array of strings   | List of whitelisted URLs for the JSON public key set returned by OAuth provider.  Can contain ``*`` wildcards.      |
		+--------------------------+--------------------+---------------------------------------------------------------------------------------------------------------------+
		| oauth_client_secret      | string             | Client secret registered with OAuth provider to verify client, passed in with `oauth_client_secret` parameter       |
		+--------------------------+--------------------+---------------------------------------------------------------------------------------------------------------------+
		| oauth_user_attribute     | string             | Client username registered with OAuth provider to verify client, passed in with `oauth_user_attribute` parameter    |
		+--------------------------+--------------------+---------------------------------------------------------------------------------------------------------------------+


	.. code-block:: json
		:caption: Example OAuth Configuration Properties In cdn.conf

		{
			"traffic_ops_golang": {
				"whitelisted_oauth_urls": [
					"oauth.example.com",
					"*.example.com",
					"username@email.com"
				],
				"oauth_client_secret": "secret"
			}
		}
