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

*********************
Environment Variables
*********************
Various :abbr:`ATC (Apache Traffic Control)` components and tools use specific environment variables to modify their behavior - in some cases they're even mandatory. What follows is a (hopefully) exhaustive list of the environment variables used in multiple places throughout the project - variables used by only a single tool or component can be found in the usage/development documentation for that tool or component. Developers are encouraged to expand the usage of variables mentioned here rather than create new ones to serve the same purpose in a different tool.

.. tip:: As environment variables can be set by any parent process, possibly unknown to the user or merely forgotten, they should generally have the lowest precedence of setting any given configuration (when multiple exist, e.g. configuration file parameters or command-line arguments) and should not be required if possible. It can be confusing to a user if ``some-tool --to-url="http://to.infra.ciab.test"`` actually uses some Traffic Ops URL other than the one they can see and specified explicitly. The suggested precedence order is:

	#. Command-line arguments
	#. Configuration file entries/parameters/lines
	#. Environment variables

.. envvar:: MOJO_MODE

	This is a legacy environment variable only used for compatibility's sake - new tools should not need to use it for anything, in general. It refers to the environment or "mode" of the Traffic Ops server from back in the old Perl days. Effectively, this chooses the set of configuration files it will consult. The default value is "development", and the possible values are:

	- development
	- integration
	- production
	- test

	:program:`admin` sets this to the value of the environment as specified by :option:`admin --env` (Default: ``development``)

.. envvar:: TO_PASSWORD

	This is used - typically in concert with :envvar:`TO_PASSWORD` and :envvar:`TO_USER` - to provide a password for a user as whom to authenticate with some Traffic Ops instance. This generally should not be validated, to avoid having to update validation to match the :ref:`to-api`'s own validation - because in general this will/should end up in the payload of an authentication request to :ref:`to-api-user-login`.

	.. caution:: For security reasons, the contents of this environment variable should not be stored anywhere for any length of time that isn't strictly necessary.

	The following list details which components/tools use this variable, how, and what, if any, restrictions they place upon its content.

	- :atc-file:`infrastructure/cdn-in-a-box/traffic_ops/to-access.sh` expects this password to authenticate the administrative-level user given by :envvar:`TO_USER`.
	- :ref:`atstccfg` uses this variable to authenticate the user used for fetching configuration information.
	- :ref:`toaccess-module` uses this variable to authenticate with the Traffic Ops instance before sending requests.

.. envvar:: TO_URL

	This is used - typically in concert with :envvar:`TO_USER` and :envvar:`TO_PASSWORD` - to identify a Traffic Ops instance for some purpose. In general, this should be able to handle the various ways to identify a network location, and need not be restricted to strictly a URL - for example, an :abbr:`FQDN (Fully Qualified Domain Name)` without a scheme should be acceptable, with or without port. When ports are not given, but scheme is, the port to use should be trivially deducible from a given supported scheme - 80 for ``http://`` and 443 for ``https://``. When the scheme is not given but the port is, the reverse assumptions should be made. When neither are given, or the port is non-standard, then ``https://`` and port 443 (if applicable) should be assumed.

	There are cases in the repository today where that is not true. The following list details the components/tools that use this variable and their restrictions and expected formats for it, if they differ from the description above.

	- :atc-file:`infrastructure/cdn-in-a-box/traffic_ops/to-access.sh` uses this variable to connect to the :ref:`ciab` Traffic Ops instance. It is passed directly to :manpage:`curl(1)` after path portions are appended, and so is subject to the restrictions and assumptions thereof.
	- The Traffic Ops :ref:`to-go-tests` for integration with the Go client use this to define the URL at which the Traffic Ops instance is running. It will override configuration file settings that specify the instance location.
	- :ref:`atstccfg` uses this variable to specify the Traffic Ops instance used to fetch configuration information.
	- :ref:`toaccess-module` uses this variable to identify the Traffic Ops instance to which requests will be sent.

.. envvar:: TO_USER

	This is used - typically in concert with :envvar:`TO_PASSWORD` and :envvar:`TO_USER` - to provide a user name as whom to authenticate with some Traffic Ops instance. This generally should not be validated, to avoid having to update validation to match the :ref:`to-api`'s own validation - because in general this will/should end up in the payload of an authentication request to :ref:`to-api-user-login`.

	The following list details which components/tools use this variable, how, and what, if any, restrictions they place upon its content.

	- :atc-file:`infrastructure/cdn-in-a-box/traffic_ops/to-access.sh` expects this to be the name of an administrative-level user.
	- :ref:`atstccfg` uses this variable to name the user as whom to authenticate for fetching configuration information.
	- :ref:`toaccess-module` uses this variable to authenticate with the Traffic Ops instance before sending requests.
