# -*- coding: utf-8 -*-

#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

"""
Module to help create/retrieve/update/delete data from/to the Traffic Ops API.

Requires Python Version >= 2.7 or >= 3.6
"""

# Core Modules
import logging
import sys

# Third-party Modules
import munch
from future.utils import raise_from
import requests.exceptions as rex

# Local Modules
from .restapi import LoginError, OperationError, api_request, RestApiSession
from .utils import log_with_debug_info

__all__ = ['TOSession']

LOGGER = logging.getLogger(__name__)

class TOSession(RestApiSession):
	"""
	Traffic Ops Session Class
	Once you login to the Traffic Ops API via :meth:`login`, you can call one or more of the methods
	to retrieve, POST, PUT, DELETE, etc. data to the API.  If you are not logged in, an	exception
	will be thrown if you try to call any of the endpoint methods. This API client is simplistic and
	lightly structured on purpose but adding support for new endpoints routinely takes seconds.
	Another nice bit of convenience that result data is, by default, wrapped in	:class:`munch.Munch`
	objects, which provide attribute access to the returned dictionaries/hashes - e.g.
	``a_dict['a_key']`` with :mod:`munch` becomes ``a_dict.a_key`` or ``a_dict['a_key']``. Also, the
	lack of rigid structure (loose coupling) means many changes to the Traffic Ops API,	as it
	evolves, will probably go un-noticed (usually additions), which means fewer	future problems to
	potentially fix in user applications.

	An area of improvement for later is defining classes to represent request data instead
	of loading up dictionaries for request data.

	Please see the :ref:`API documentation <to-api>` for the details of the API endpoints.

	Adding end-point methods

	.. code-block:: python3
		:caption: Endpoint with no URL parameters and no query parameters

		@api_request(u'get', u'cdns', (u'1.1', u'1.2',))
		def get_cdns(self):
			pass


	.. code-block:: python3
		:caption: End-point with URL parameters and no query parameters

		@api_request(u'get', u'cdns/{cdn_id:d}', (u'1.1', u'1.2',))
		def get_cdn_by_id(self, cdn_id=None):
			pass


	.. code-block:: python3
		:caption: End-point with no URL parameters but with query parameters

		@api_request(u'get', u'deliveryservices', (u'1.1', u'1.2',))
		def get_deliveryservices(self, query_params=None):
			pass

	.. code-block:: python3
		:caption: End-point with URL parameters and query parameters

		@api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys', (u'1.1', u'1.2',))
		def get_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
			pass

	.. code-block:: python3
		:caption: End-point with request data

		@api_request(u'post', u'cdns', (u'1.1', u'1.2',))
		def create_cdn(self, data=None):
			pass

	.. code-block:: python3
		:caption: End-point with URL parameters and request data

		@api_request(u'put', u'cdns', (u'1.1', u'1.2',))
		def update_cdn_by_id(self, cdn_id=None, data=None):
			pass

	Calling end-point methods

	:meth:`get_cdns` calls endpoint :ref:`to-api-cdns` e.g. ``t.get_cdns()``

	:meth:`get_types` calls endpoint :ref:`to-api-types`, optionally with query parameters e.g. ``get_foo_data(id=45, query_params={'sort': 'asc'})`` calls endpoint ``GET api/1.x/foo/45?sort=asc`` (presumably)

	:meth:`cdns_queue_update` calls endpoint :ref:`to-api-cdns-id-queue_update`, with an ID path parameter and a JSON payload e.g. ``cdns_queue_update(id=1, data={'action': 'queue'})``

	.. note:: Only a small subset of the API endpoints are implemented. More can be implemented as needed.
	"""

	def __init__(self, host_ip, host_port=443, api_version=u'1.3', ssl=True, headers=None,
	             verify_cert=True):
		"""
		The class initializer.
		:param host_ip: The dns name or ip address of the Traffic Ops host to use to talk to the API
		:type host_ip: str
		:param host_port: The port to use when contacting the Traffic Ops API
		:type host_port: int
		:param api_version: The version of the API to use when calling end-points on the Traffic Ops API
		:type api_version: str
		:param ssl: Should ssl be used? (http vs. https)
		:type ssl: bool
		:param headers:  The http headers to use when contacting the Traffic Ops API
		:type headers: Dict[str, str]
		:type verify_cert: bool
		"""
		super(TOSession, self).__init__(host_ip=host_ip, api_version=api_version,
		                                api_base_path=u'api/{api_version}/',
		                                host_port=host_port, ssl=ssl, headers=headers,
		                                verify_cert=verify_cert)

		self._logged_in = False

		msg = u'TOSession instance {0:#0x} initialized: Details: {1}'
		log_with_debug_info(logging.DEBUG, msg.format(id(self), self.__dict__))

	def login(self, username, password):
		"""
		Login to the Traffic Ops API.
		:param username: Traffic Ops User Name
		:type username: str
		:param password: Traffic Ops User Password
		:type password: str
		:return: None
		:rtype: None
		:raises: LoginError
		"""
		logging.info("Connecting to Traffic Ops at %s...", self.to_url)

		if not self.is_open:
			self.create()

		logging.info("Connected. Authenticating...")

		self._logged_in = False
		try:
			# Try to login to Traffic Ops
			self.post(u'user/login', data={u'u': username, u'p': password})
			self._logged_in = True
		except rex.SSLError as e:
			logging.debug("%s", e, stack_info=True, exc_info=True)
			self.close()
			msg = (u'{0}.  This system may have a self-signed certificate.  Try creating this'\
			       u' TOSession object passing verify_cert=False. e.g. TOSession(..., '\
			       u'verify_cert=False).')
			msg = msg.format(e)
			logging.error(msg)
			logging.warning("disabling certificate verification is not recommended.")
			raise_from(LoginError(msg), e)
		except OperationError as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			msg = u'Logging in to Traffic Ops has failed. Reason: {0}'.format(e)
			self.close()
			logging.error(msg)
			raise_from(OperationError(msg), e)

		logging.info("Authenticated.")

	@property
	def to_url(self):
		"""
		The URL without the api portion. (read-only)

		:return: The URL should match '[\\w\\+\\-\\.]+://[\\w\\+\\-\\.]+(:\\d+)?' e.g https://to.somedomain.net or https://to.somedomain.net:443
		:rtype: str
		"""

		return self.server_url

	@property
	def base_url(self):
		"""
		Returns the base url. (read-only)

		:return: The base url should match '[\\w\\+\\-\\.]+://[\\w\\+\\-\\.]+(:\\d+)?' e.g https://to.somedomain.net/api/1.2/
		:rtype: str
		"""

		return self._api_base_url

	@property
	def logged_in(self):
		"""
		Read-only property of to determine if user is logged in to Traffic Ops.
		:return: :const:`True` if connected and logged in, :const:`False` otherwise
		:rtype: bool
		"""

		return self.is_open and self._logged_in

	# Programmatic Endpoint Methods - These can be created when you need to employ "creative
	# methods" to form a correlated composite data set from one or more Traffic Ops API call(s) or
	# employ composite operations against the API.
	# Also, if the API requires you to retrieve the data via paging, these types of methods can be
	# useful to perform that type of work too.
	# These methods need to support similar method signatures as employed by the restapi.api_request
	# decorator method_name argument.
	def get_all_deliveryservice_servers(self, *args, **kwargs):
		"""
		Get all servers attached to all delivery services via the Traffic Ops API.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""
		result_set = []
		resp = None
		limit = 10000
		page = 1

		munchify = True  # Default to True
		if u'munchify' in kwargs:
			munchify = kwargs[u'munchify']

		while True:
			data, resp = self.get_deliveryserviceserver(query_params={u'limit':limit, u'page': page},
			                                           *args, **kwargs)

			if not data:
				break

			result_set.extend(munch.munchify(data) if munchify else data)
			page += 1

		return result_set, resp  # Note: Return last response object received

#
# PUT ALL API DEFINITIONS BELOW AND UNDER ITS RESPECTIVE PAGE (whether it is 1.2 or 1.3, etc, if its
# a CDN put it under CDN header and corresponding calls)
#

	#
	#	API Capabilities
	#
	@api_request(u'get', u'api_capabilities', (u'1.2', u'1.3',))
	def get_api_capabilities(self, query_params=None):
		"""
		Get all API-capability mappings
		:ref:`to-api-api_capabilities`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'api_capabilities/{id}', (u'1.2', u'1.3',))
	def get_api_capabilities_by_id(self, id=None):
		"""
		Get an API-capability mapping by ID
		:ref:`to-api-api_capabilities-id`
		:param id: The api-capabilities Id
		:type id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# ASN
	#
	@api_request(u'get', u'asns', (u'1.1', u'1.2', u'1.3',))
	def get_asns(self, query_params=None):
		"""
		Get ASNs.
		:ref:`to-api-asns`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'asns/{asn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_asn_by_id(self, asn_id=None):
		"""
		Get ASN by ID
		:ref:`to-api-asns-id`
		:param asn_id: The ID of the ASN to retrieve
		:type asn_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'asns/{asn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_asn(self, asn_id=None, query_params=None):
		"""
		Update ASN
		:ref:`to-api-asns-id`
		:param asn_id: The ID of the ASN to update
		:type asn_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'asns', (u'1.1', u'1.2', u'1.3',))
	def create_asn(self, data=None):
		"""
		Create ASN
		:ref:`to-api-asns`
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'asns/{asn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_asn(self, asn_id=None):
		"""
		Delete ASN
		:to-api-asns-id:
		:param asn_id: The ID of the ASN to delete
		:type asn_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Stats
	#
	@api_request(u'get', u'caches/stats', (u'1.1', u'1.2', u'1.3',))
	def get_traffic_monitor_cache_stats(self):
		"""
		Retrieves cache stats from Traffic Monitor. Also includes rows for aggregates
		:ref:`to-api-caches-stats`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Group
	#
	@api_request(u'get', u'cachegroups', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroups(self, query_params=None):
		"""
		Get Cache Groups.
		:ref:`to-api-cachegroups`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	# Example of URL parameter substitution via call parameter. You will need to pass the parameter
	# value as a keyword parameter with the proper type to match the str.format specification,
	# e.g. 'cachegroups/{cache_group_id:d}'.  In this case, ':d' specifies a decimal integer. A
	# specification of 'cachegroups/{cache_group_id}' will try to convert any value passed to a
	# string, which basically does no type checking, unless of course the value cannot be cast to a
	# string e.g. get_cachegroups_by_id(cache_group_id=23) -> call end-point .../api/1.2/cachegroups/23
	@api_request(u'get', u'cachegroups/{cache_group_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroup_by_id(self, cache_group_id=None):
		"""
		Get a Cache Group by Id.
		:ref:`to-api-cachegroups-id`
		:param cache_group_id: The cache group Id
		:type cache_group_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cachegroups/{cache_group_id:d}/parameters', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroup_parameters(self, cache_group_id=None):
		"""
		Get a cache groups parameters
		:ref:`to-api-cachegroups-id-parameters`
		:param cache_group_id: The cache group Id
		:type cache_group_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cachegroups/{cache_group_id:d}/unassigned_parameters',
	             (u'1.1', u'1.2', u'1.3'))
	def get_cachegroup_unassigned_parameters(self, cache_group_id=None):
		"""
		Get a cache groups unassigned parameters
		:ref:`to-api-cachegroups-id-unassigned_parameters`
		:param cache_group_id: The cache group Id
		:type cache_group_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cachegroup/{parameter_id:d}/parameter', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroup_parameters_by_id(self, parameter_id=None):
		"""
		Get a cache groups parameter by its ID
		:ref:`to-api-cachegroup-parameterID-parameter`
		:param parameter_id: The parameter Id
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cachegroupparameters', (u'1.1', u'1.2', u'1.3',))
	def get_all_cachegroup_parameters(self):
		"""
		A collection of all cache group parameters.
		:ref:`to-api-cachegroupparameters`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'cachegroups', (u'1.1', u'1.2', u'1.3',))
	def create_cachegroups(self, data=None):
		"""
		Create a Cache Group
		:ref:`to-api-cachegroups`
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'cachegroups/{cache_group_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_cachegroups(self, cache_group_id=None, data=None):
		"""
		Update a cache group
		:ref:`to-api-cachegroups-id`
		:param cache_group_id: The cache group id to update
		:type cache_group_id: Integer
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'cachegroups/{cache_group_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_cachegroups(self, cache_group_id=None):
		"""
		Delete a cache group
		:ref:`to-api-cachegroups-id`
		:param cache_group_id: The cache group id to update
		:type cache_group_id: Integer
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'cachegroups/{cache_group_id:d}/queue_update', (u'1.1', u'1.2', u'1.3',))
	def cachegroups_queue_update(self, cache_group_id=None, data=None):
		"""
		Queue Updates by Cache Group ID
		:ref:`to-api-cachegroups-id-queue_update`
		:param cache_group_id: The Cache Group Id
		:type cache_group_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'cachegroupparameters/{cache_group_id:d}/{parameter_id:d}',
	             (u'1.2', u'1.3'))
	def delete_cache_group_parameters(self, cache_group_id=None, parameter_id=None):
		"""
		Delete a cache group parameter association
		:ref:`to-api-cachegroupparameters-id-parameterID`
		:param cache_group_id: The cache group id in which the parameter will be deleted
		:type cache_group_id: int
		:param parameter_id: The parameter id which will be disassociated
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Group Fallback
	#
	@api_request(u'get', u'cachegroup_fallbacks', (u'1.2', u'1.3'))
	def get_cache_group_fallbacks(self, query_params=None):
		"""
		Retrieve fallback related configurations for a cache group
		:ref:`to-api-cachegroup_fallbacks`

		:param query_params: Either cacheGroupId or fallbackId must be used or can be used simultaneously
		:type query_params: Dict[str, int]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'cachegroup_fallbacks', (u'1.2', u'1.3'))
	def create_cache_group_fallbacks(self, data=None):
		"""
		Creates fallback configuration for the cache group. New fallbacks can be added only via POST.
		:ref:`to-api-cachegroup_fallbacks`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'cachegroup_fallbacks', (u'1.2', u'1.3'))
	def update_cache_group_fallbacks(self, data=None):
		"""
		Updates an existing fallback configuration for the cache group.
		:ref:`to-api-cachegroup_fallbacks`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'cachegroup_fallbacks', (u'1.2', u'1.3'))
	def delete_cache_group_fallbacks(self, query_params=None):
		"""
		Deletes an existing fallback related configurations for a cache group
		:ref:`to-api-cachegroup_fallbacks`

		:param query_params: Either cacheGroupId or fallbackId must be used or can be used simultaneously
		:type query_params: Dict[str, int]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Statistics
	#
	@api_request(u'get', u'cache_stats', (u'1.2', u'1.3',))
	def get_cache_stats(self, query_params=None):
		"""
		Retrieves statistics about the CDN.
		:ref:`to-api-cache_stats`
		:param query_params: See API page for more information on accepted params
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN
	#
	@api_request(u'get', u'cdns', (u'1.1', u'1.2', u'1.3',))
	def get_cdns(self):
		"""
		Get all CDNs.
		:ref:`to-api-cdns`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/{cdn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_cdn_by_id(self, cdn_id=None):
		"""
		Get a CDN by Id.
		:ref:`to-api-cdns-id`
		:param cdn_id: The CDN id
		:type cdn_id: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/name/{cdn_name}', (u'1.1', u'1.2', u'1.3',))
	def get_cdn_by_name(self, cdn_name=None):
		"""
		Get a CDN by name.
		:ref:`to-api-cdns-name-name`
		:param cdn_name: The CDN name
		:type cdn_name: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'cdns', (u'1.1', u'1.2', u'1.3',))
	def create_cdn(self, data=None):
		"""
		Create a new CDN.
		:ref:`to-api-cdns`
		:param data: The parameter data to use for cdn creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'cdns/{cdn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_cdn_by_id(self, cdn_id=None, data=None):
		"""
		Update a CDN by Id.
		:ref:`to-api-cdns-id`
		:param cdn_id: The CDN id
		:type cdn_id: int
		:param data: The parameter data to use for cdn update.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'cdns/{cdn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_cdn_by_id(self, cdn_id=None):
		"""
		Delete a CDN by Id.
		:ref:`to-api-cdns-id`
		:param cdn_id: The CDN id
		:type cdn_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'cdns/{cdn_id:d}/queue_update', (u'1.1', u'1.2', u'1.3',))
	def cdns_queue_update(self, cdn_id=None, data=None):
		"""
		Queue Updates by CDN Id.
		:ref:`to-api-cdns-id-queue_update`
		:param cdn_id: The CDN Id
		:type cdn_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Health
	#
	@api_request(u'get', u'cdns/health', (u'1.2', u'1.3',))
	def get_cdns_health(self):
		"""
		Retrieves the health of all locations (cache groups) for all CDNs
		:ref:`to-api-cdns-health`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'get', u'cdns/{cdn_name:s}/health', (u'1.2', u'1.3',))
	def get_cdn_health_by_name(self, cdn_name=None):
		"""
		Retrieves the health of all locations (cache groups) for a given CDN
		:ref:`to-api-cdns-name-health`
		:param cdn_name: The CDN name to find health for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/usage/overview', (u'1.2', u'1.3',))
	def get_cdns_usage(self):
		"""
		Retrieves the high-level CDN usage metrics.
		:ref:`to-api-cdns-usage-overview`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/capacity', (u'1.2', u'1.3',))
	def get_cdns_capacity(self):
		"""
		Retrieves the aggregate capacity percentages of all locations (cache groups) for a given CDN.
		:ref:`to-api-cdns-capacity`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Routing
	#
	@api_request(u'get', u'cdns/routing', (u'1.2', u'1.3',))
	def get_cdns_routing(self):
		"""
		Retrieves the aggregate routing percentages of all locations (cache groups) for a given CDN.
		:ref:`to-api-cdns-routing`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Domains
	#
	@api_request(u'get', u'cdns/domains', (u'1.2', u'1.3',))
	def get_cdns_domains(self):
		"""
		Retrieves the different CDN domains
		:ref:`to-api-cdns-domains`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Topology
	#

	# At the time of this writing (Tues. Jan 15 2019), this endpoint doesn't appear to exist.
	# @api_request(u'get', u'cdns/{cdn_name:s}/configs', (u'1.2', u'1.3',))
	# def get_cdn_config_info(self, cdn_name=None):
	# 	"""
	# 	Retrieves CDN config information
	# 	:ref:`to-api-cdns-name-configs`
	# 	:param cdn_name: The CDN name to find configs for
	# 	:type cdn_name: String
	# 	:rtype: Tuple[Dict[str, Any], requests.Response]
	# 	:raises: Union[LoginError, OperationError]
	# 	"""

	@api_request(u'get', u'cdns/{cdn_name:s}/configs/monitoring', (u'1.2', u'1.3',))
	def get_cdn_monitoring_info(self, cdn_name=None):
		"""
		Retrieves CDN monitoring information
		:ref:`to-api-cdns-name-configs-monitoring`
		:param cdn_name: The CDN name to find configs for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/{cdn_name:s}/configs/routing', (u'1.2', u'1.3',))
	def get_cdn_routing_info(self, cdn_name=None):
		"""
		Retrieves CDN routing information
		:ref:`to-api-cdns-name-configs-routing`
		:param cdn_name: The CDN name to find routing info for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN DNSSEC Keys
	#
	@api_request(u'get', u'cdns/name/{cdn_name:s}/dnsseckeys', (u'1.2', u'1.3',))
	def get_cdn_dns_sec_keys(self, cdn_name=None):
		"""
		Gets a list of dnsseckeys for a CDN and all associated Delivery Services
		:ref:`to-api-cdns-name-name-dnsseckeys`
		:param cdn_name: The CDN name to find dnsseckeys info for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/name/{cdn_name:s}/dnsseckeys/delete', (u'1.2', u'1.3',))
	def delete_cdn_dns_sec_keys(self, cdn_name=None):
		"""
		Delete dnssec keys for a cdn and all associated delivery services
		:ref:`to-api-cdns-name-name-dnsseckeys-delete`
		:param cdn_name: The CDN name to delete dnsseckeys info for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryservices/dnsseckeys/generate', (u'1.2', u'1.3',))
	def create_cdn_dns_sec_keys(self, data=None):
		"""
		Generates ZSK and KSK keypairs for a CDN and all associated Delivery Services
		:ref:`to-api-deliveryservices-dnsseckeys-generate`
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN SSL Keys
	#
	@api_request(u'get', u'cdns/name/{cdn_name:s}/sslkeys', (u'1.2', u'1.3', u'1.4'))
	def get_cdn_ssl_keys(self, cdn_name=None):
		"""
		Returns ssl certificates for all Delivery Services that are a part of the CDN.
		:ref:`to-api-cdns-name-name-sslkeys`
		:param cdn_name: The CDN name to find ssl keys for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Change Logs
	# https://traffic-control-cdn.readthedocs.io/en/latest/api/changelog.html#change-logs
	#

	@api_request(u'get', u'logs', (u'1.2', u'1.3',))
	def get_change_logs(self):
		"""
		Retrieve all change logs from traffic ops
		:ref:`to-api-logs`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'logs/{days:d}/days', (u'1.2', u'1.3',))
	def get_change_logs_for_days(self, days=None):
		"""
		Retrieve all change logs from Traffic Ops
		:ref:`to-api-logs-days-days`
		:param days: The number of days to retrieve change logs
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'logs/newcount', (u'1.2', u'1.3',))
	def get_change_logs_newcount(self):
		"""
		Get amount of new logs from traffic ops
		:ref:`to-api-logs-newcount`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Config Files and Config File Metadata
	#
	@api_request(u'get', u'servers/{host_name:s}/configfiles/ats', (u'1.2', u'1.3', u'1.4'))
	def get_server_config_files(self, host_name=None, query_params=None):
		"""
		Get the configuiration files for a given host name
		:ref:`to-api-servers-server-configfiles-ats`
		:param host_name: The host name to get config files for
		:type host_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/{host_name:s}/configfiles/ats/{config_file:s}', (u'1.2', u'1.3', u'1.4'))
	def get_server_specific_config_file(self, host_name=None, config_file=None, query_params=None):
		"""
		Get the configuiration files for a given host name and config file
		:ref:`to-api-servers-server-configfiles-ats-filename`
		:param host_name: The host name to get config files for
		:type host_name: String
		:param config_file: The config file name to retrieve for host
		:type config_file: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'profiles/{profile_name:s}/configfiles/ats/{config_file:s}',
	             (u'1.2', u'1.3',))
	def get_profile_specific_config_files(self, profile_name=None, config_file=None,
	                                      query_params=None):
		"""
		Get the configuiration files for a given profile name and config file
		:ref:`to-api-profiles-profile-configfiles-ats-filename`
		:param profile_name: The profile name to get config files for
		:type host_name: String
		:param config_file: The config file name to retrieve for host
		:type config_file: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/{cdn_name:s}/configfiles/ats/{config_file:s}', (u'1.2', u'1.3',))
	def get_cdn_specific_config_file(self, cdn_name=None, config_file=None, query_params=None):
		"""
		Get the configuiration files for a given cdn name and config file
		:ref:`to-api-cdns-cdn-configfiles-ats-filename`
		:param cdn_name: The cdn name to get config files for
		:type cdn_name: String
		:param config_file: The config file name to retrieve for host
		:type config_file: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Delivery Service
	#
	@api_request(u'get', u'deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservices(self, query_params=None):
		"""
		Retrieves all delivery services (if admin or ops) or all delivery services assigned to user.
		:ref:`to-api-deliveryservices`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_by_id(self, delivery_service_id=None):
		"""
		Retrieves a specific delivery service. If not admin / ops, delivery service must be assigned
		to user.
		:ref:`to-api-deliveryservices-id`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/servers',
	             (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_servers(self, delivery_service_id=None):
		"""
		Retrieves properties of CDN EDGE or ORG servers assigned to a delivery service.
		:ref:`to-api-deliveryservices-id-servers`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/unassigned_servers',
	             (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_unassigned_servers(self, delivery_service_id=None):
		"""
		Retrieves properties of CDN EDGE or ORG servers not assigned to a delivery service.
		(Currently call does not work)
		:ref:`to-api-deliveryservices-id-unassigned_servers`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/servers/eligible',
	             (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_ineligible_servers(self, delivery_service_id=None):
		"""
		Retrieves properties of CDN EDGE or ORG servers not eligible for assignment to a delivery
		service.
		:ref:`to-api-deliveryservices-id-servers-eligible`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def create_deliveryservice(self, data=None):
		"""
		Allows user to create a delivery service.
		:ref:`to-api-deliveryservices`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_deliveryservice_by_id(self, delivery_service_id=None, data=None):
		"""
		Update a Delivery Service by Id.
		:ref:`to-api-deliveryservices-id`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'deliveryservices/{delivery_service_id:d}/safe', (u'1.1', u'1.2', u'1.3',))
	def update_deliveryservice_safe(self, delivery_service_id=None, data=None):
		"""
		Allows a user to edit limited fields of an assigned delivery service.
		:ref:`to-api-deliveryservices-id-safe`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'delete', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_by_id(self, delivery_service_id=None):
		"""
		Allows user to delete a delivery service.
		:ref:`to-api-deliveryservices-id`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Health
	#
	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/state', (u'1.1', u'1.2', u'1.3',))
	def get_delivery_service_failover_state(self, delivery_service_id=None):
		"""
		Retrieves the failover state for a delivery service. Delivery service must be assigned to
		user if user is not admin or operations.
		:ref:`to-api-deliveryservices-id-state`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/health', (u'1.1', u'1.2', u'1.3',))
	def get_delivery_service_health(self, delivery_service_id=None):
		"""
		Retrieves the health of all locations (cache groups) for a delivery service. Delivery
		service must be assigned to user if user is not admin or operations.
		:ref:`to-api-deliveryservices-id-health`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/capacity',
	             (u'1.1', u'1.2', u'1.3',))
	def get_delivery_service_capacity(self, delivery_service_id=None):
		"""
		Retrieves the capacity percentages of a delivery service. Delivery service must be assigned
		to user if user is not admin or operations.
		:ref:`to-api-deliveryservices-id-capacity`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/routing',
	             (u'1.1', u'1.2', u'1.3',))
	def get_delivery_service_routing(self, delivery_service_id=None):
		"""
		Retrieves the routing method percentages of a delivery service. Delivery service must be
		assigned to user if user is not admin or operations.
		:ref:`to-api-deliveryservices-id-routing`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Server
	#
	@api_request(u'get', u'deliveryserviceserver', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryserviceserver(self, query_params=None):
		"""
		Retrieves delivery service / server assignments. (Allows pagination and limits)
		:ref:`to-api-deliveryserviceserver`
		:param query_params: The required url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryserviceserver', (u'1.1', u'1.2', u'1.3',))
	def assign_deliveryservice_servers_by_ids(self, data=None):
		"""
		Assign servers by id to a Delivery Service. (New Method)
		:ref:`to-api-deliveryserviceserver`
		:param data: The required data to create server associations to a delivery service
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryservices/{xml_id}/servers', (u'1.1', u'1.2', u'1.3',))
	def assign_deliveryservice_servers_by_names(self, xml_id=None, data=None):
		"""
		Assign servers by name to a Delivery Service by xmlId.
		:ref:`to-api-deliveryservices-xmlid-servers`
		:param xml_id: The XML Id of the delivery service
		:type xml_id: str
		:param data: The required data to assign servers to a delivery service
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'deliveryservice_server/{delivery_service_id:d}/{server_id:d}',
	             (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_servers_by_id(self, delivery_service_id=None, server_id=None):
		"""
		Removes a server (cache) from a delivery service.
		:ref:`to-api-deliveryservice_server-dsid-serverid`
		:param delivery_service_id: The delivery service id
		:type delivery_service_id: int
		:param server_id: The server id to remove from delivery service
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service User
	#
	@api_request(u'post', u'deliveryservice_user', (u'1.2', u'1.3'))
	def create_delivery_service_user_link(self, data=None):
		"""
		Create one or more user / delivery service assignments.
		:ref:`to-api-deliveryservice_user`
		:param data: The parameter data to use for Delivery Service SSL key generation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'deliveryservice_user/{delivery_service_id:d}/{user_id:d}',
	             (u'1.2', u'1.3'))
	def delete_delivery_service_user_link(self, delivery_service_id=None, user_id=None):
		"""
		Removes a delivery service from a user.
		:ref:`to-api-deliveryservice_user-dsid-userid`
		:param delivery_service_id: The delivery service id to dissasociate the user
		:type delivery_service_id: int
		:param user_id: The user id to dissassociate
		:type user_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Delivery Service SSL Keys
	#
	@api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
		"""
		Get SSL keys for a Delivery Service by xmlId.
		:ref:`to-api-deliveryservices-xmlid-xmlid-sslkeys`
		:param xml_id: The Delivery Service XML id
		:type xml_id: str
		:param query_params: The url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys/delete', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
		"""
		Delete SSL keys for a Delivery Service by xmlId.
		:ref:`to-api-deliveryservices-xmlid-xmlid-sslkeys-delete`
		:param xml_id: The Delivery Service xmlId
		:type xml_id: str
		:param query_params: The url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryservices/sslkeys/generate', (u'1.1', u'1.2', u'1.3',))
	def generate_deliveryservice_ssl_keys(self, data=None):
		"""
		Generate an SSL certificate. (self-signed)
		:ref:`to-api-deliveryservices-sslkeys-generate`
		:param data: The parameter data to use for Delivery Service SSL key generation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryservices/sslkeys/add', (u'1.1', u'1.2', u'1.3',))
	def add_ssl_keys_to_deliveryservice(self, data=None):
		"""
		Add SSL keys to a Delivery Service.
		:ref:`to-api-deliveryservices-sslkeys-add`
		:param data: The parameter data to use for adding SSL keys to a Delivery Service.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service URL Sig Keys
	#
	@api_request(u'post', u'deliveryservices/xmlId/{xml_id}/urlkeys/generate',
	             (u'1.1', u'1.2', u'1.3',))
	def generate_deliveryservice_url_signature_keys(self, xml_id=None):
		"""
		Generate URL Signature Keys for a Delivery Service by xmlId.
		:ref:`to-api-deliveryservices-xmlid-xmlid-urlkeys-generate`
		:param xml_id: The Delivery Service xmlId
		:type xml_id: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Regexes
	#
	@api_request(u'get', u'deliveryservices_regexes', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservices_regexes(self):
		"""
		Get RegExes for all Delivery Services.
		:ref:`to-api-deliveryservices_regexes`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/regexes',
	             (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_regexes_by_id(self, delivery_service_id=None):
		"""
		Get RegExes for a Delivery Service by Id.
		:ref:`to-api-deliveryservices-id-regexes`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'deliveryservices/{delivery_service_id:d}/regexes/{regex_id:d}',
	             (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_regexes_by_regex_id(self, delivery_service_id=None, regex_id=None):
		"""
		Retrieves a regex for a specific delivery service.
		:ref:`to-api-deliveryservices-id-regexes-rid`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param regex_id: The delivery service regex id
		:type regex_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'deliveryservices/{delivery_service_id:d}/regexes',
	             (u'1.1', u'1.2', u'1.3',))
	def create_deliveryservice_regexes(self, delivery_service_id=None, data=None):
		"""
		Create a regex for a delivery service
		:ref:`to-api-deliveryservices-id-regexes`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param data: The required data to create delivery service regexes
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'deliveryservices/{delivery_service_id:d}/regexes/{regex_id:d}',
	             (u'1.1', u'1.2', u'1.3',))
	def update_deliveryservice_regexes(self, delivery_service_id=None, regex_id=None,
	                                   query_params=None):
		"""
		Update a regex for a delivery service
		:ref:`to-api-deliveryservices-id-regexes-rid`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param regex_id: The delivery service regex id
		:type regex_id: int
		:param query_params: The required data to update delivery service regexes
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'deliveryservices/{delivery_service_id:d}/regexes/'\
	                        u'{delivery_service_regex_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_regex_by_regex_id(self, delivery_service_id=None,
	                                             delivery_service_regex_id=None):
		"""
		Delete a RegEx by Id for a Delivery Service by Id.
		:ref:`to-api-deliveryservices-id-regexes-rid`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param delivery_service_regex_id: The delivery service regex Id
		:type delivery_service_regex_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Statistics
	#
	@api_request(u'get', u'deliveryservice_stats', (u'1.1', u'1.2', u'1.3',))
	def get_delivery_service_stats(self, query_params=None):
		"""
		Retrieves statistics on the delivery services.
		:ref:`to-api-deliveryservice_stats`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Divisions
	#
	@api_request(u'get', u'divisions', (u'1.1', u'1.2', u'1.3',))
	def get_divisions(self):
		"""
		Get all divisions.
		:ref:`to-api-divisions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'divisions/{division_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_division_by_id(self, division_id=None):
		"""
		Get a division by division id
		:ref:`to-api-divisions-id`
		:param division_id: The division id to retrieve
		:type division_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'divisions/{division_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_division(self, division_id=None, query_params=None):
		"""
		Update a division by division id
		:ref:`to-api-divisions-id`
		:param division_id: The division id to update
		:type division_id: int
		:param query_params: The required data to update delivery service regexes
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'divisions', (u'1.1', u'1.2', u'1.3',))
	def create_division(self, data=None):
		"""
		Create a division
		:ref:`to-api-divisions`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'divisions/{division_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_division(self, division_id=None, query_params=None):
		"""
		Delete a division by division id
		:ref:`to-api-divisions-id`
		:param division_id: The division id to delete
		:type division_id: int
		:param query_params: The required data to update delivery service regexes
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation
	#
	@api_request(u'get', u'federations', (u'1.2', u'1.3'))
	def get_federations(self):
		"""
		Retrieves a list of federation mappings (aka federation resolvers) for a the current user
		:ref:`to-api-federations`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'federations', (u'1.2', u'1.3'))
	def create_federation(self, data=None):
		"""
		Allows a user to add federations for their delivery service(s).
		:ref:`to-api-federations`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'get', u'cdns/{cdn_name:s}/federations', (u'1.2', u'1.3'))
	def get_federations_for_cdn(self, cdn_name=None):
		"""
		Retrieves a list of federations for a cdn.
		:ref:`to-api-cdns-name-federations`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'get', u'cdns/{cdn_name:s}/federations/{federation_id:d}', (u'1.2', u'1.3'))
	def get_federation_for_cdn_by_id(self, cdn_name=None, federation_id=None):
		"""
		Retrieves a federation for a cdn.
		:ref:`to-api-cdns-name-federations-id`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'cdns/{cdn_name:s}/federations', (u'1.2', u'1.3'))
	def create_federation_in_cdn(self, cdn_name=None, data=None):
		"""
		Create a federation.
		:ref:`to-api-cdns-name-federations`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'cdns/{cdn_name:s}/federations/{federation_id:d}', (u'1.2', u'1.3'))
	def update_federation_in_cdn(self, cdn_name=None, federation_id=None, query_params=None):
		"""
		Update a federation.
		:ref:`to-api-cdns-name-federations-id`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:param federation_id: The federation id
		:type federation_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'cdns/{cdn_name:s}/federations/{federation_id:d}', (u'1.2', u'1.3'))
	def delete_federation_in_cdn(self, cdn_name=None, federation_id=None):
		"""
		Delete a federation.
		:ref:`to-api-cdns-name-federations-id`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:param federation_id: The federation id
		:type federation_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation Delivery Service
	#
	@api_request(u'get', u'federations/{federation_id:d}/deliveryservices', (u'1.2', u'1.3'))
	def get_federation_delivery_services(self, federation_id=None):
		"""
		Retrieves delivery services assigned to a federation
		:ref:`to-api-federations-id-deliveryservices`
		:param federation_id: The federation id
		:type federation_id: int
		:param federation_id: The federation id
		:type federation_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'federations/{federation_id:d}/deliveryservices', (u'1.2', u'1.3'))
	def assign_delivery_services_to_federations(self, federation_id=None, data=None):
		"""
		Create one or more federation / delivery service assignments.
		:ref:`to-api-federations-id-deliveryservices`
		:param federation_id: The federation id
		:type federation_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation Federation Resolver
	#
	@api_request(u'get', u'federations/{federation_id:d}/federation_resolvers', (u'1.2', u'1.3'))
	def get_federation_resolvers_by_id(self, federation_id=None):
		"""
		:ref:`to-api-federations-id-federation_resolvers`
		Retrieves federation resolvers assigned to a federation
		:param federation_id: The federation id
		:type federation_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'federations/{federation_id:d}/federation_resolvers', (u'1.2', u'1.3'))
	def assign_federation_resolver_to_federations(self, federation_id=None, data=None):
		"""
		Create one or more federation / federation resolver assignments.
		:ref:`to-api-federations-id-federation_resolvers`
		:param federation_id: The federation id
		:type federation_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation Resolver
	#
	@api_request(u'get', u'federation_resolvers', (u'1.2', u'1.3'))
	def get_federation_resolvers(self, query_params=None):
		"""
		Get federation resolvers.
		:ref:`to-api-federation_resolvers`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'federation_resolvers', (u'1.2', u'1.3'))
	def create_federation_resolver(self, data=None):
		"""
		Create a federation resolver.
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'federation_resolvers/{federation_resolver_id:d}', (u'1.2', u'1.3'))
	def delete_federation_resolver(self, federation_resolver_id=None):
		"""
		Delete a federation resolver.
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation User
	#
	@api_request(u'get', u'federations/{federation_id:d}/users', (u'1.2', u'1.3'))
	def get_federation_users(self, federation_id=None):
		"""
		Retrieves users assigned to a federation.
		:ref:`to-api-federations-id-users`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'federations/{federation_id:d}/users', (u'1.2', u'1.3'))
	def create_federation_user(self, federation_id=None, data=None):
		"""
		Create one or more federation / user assignments.
		:ref:`to-api-federations-id-users`
		:param federation_id: Federation ID
		:type federation_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'federations/{federation_id:d}/users/{user_id:d}', (u'1.2', u'1.3'))
	def delete_federation_user(self, federation_id=None, user_id=None):
		"""
		Delete one or more federation / user assignments.
		:ref:`to-api-federations-id-users-id`
		:param federation_id: Federation ID
		:type federation_id: int
		:param user_id: Federation User ID
		:type user_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Hardware Info
	#
	@api_request(u'get', u'hwinfo', (u'1.2', u'1.3'))
	def get_hwinfo(self):
		"""
		Get hwinfo for servers.
		:ref:`to-api-hwinfo`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# ISO
	#
	@api_request(u'get', u'osversions', (u'1.2', u'1.3'))
	def get_osversions(self):
		"""
		Get all OS versions for ISO generation and the directory where the kickstarter files are
		found. The values are retrieved from osversions.cfg found in either ``/var/www/files`` or in
		the location defined by the kickstart.files.location parameter (if defined).
		:ref:`to-api-osversions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'isos', (u'1.2', u'1.3',))
	def generate_iso(self, data=None):
		"""
		Generate an ISO
		:ref:`to-api-isos`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Jobs
	#
	@api_request(u'get', u'jobs', (u'1.2', u'1.3'))
	def get_jobs(self, query_params=None):
		"""
		Get all jobs (currently limited to invalidate content (PURGE) jobs) sorted by start time
		(descending).
		:ref:`to-api-jobs`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'jobs/{job_id:d}', (u'1.2', u'1.3'))
	def get_job_by_id(self, job_id=None):
		"""
		Get a job by ID (currently limited to invalidate content (PURGE) jobs).
		:ref:`to-api-jobs-id`
		:param job_id: The job id to retrieve
		:type job_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Parameter
	#
	@api_request(u'get', u'parameters', (u'1.2', u'1.3',))
	def get_parameters(self):
		"""
		Get all Profile Parameters.
		:ref:`to-api-parameters`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'parameters/{parameter_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_parameter_by_id(self, parameter_id=None):
		"""
		Get a Parameter by Id.
		:ref:`to-api-parameters-id`
		:param parameter_id: The parameter Id
		:type parameter_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'parameters/{parameter_id:d}/profiles', (u'1.1', u'1.2', u'1.3',))
	def get_associated_profiles_by_parameter_id(self, parameter_id=None):
		"""
		Get all Profiles associated to a Parameter by Id.
		:ref:`to-api-parameters-id-profiles`
		:param parameter_id: The parameter id
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'parameters/{parameter_id:d}/unassigned_profiles', (u'1.1', u'1.2', u'1.3',))
	def get_unassigned_profiles_by_parameter_id(self, parameter_id=None):
		"""
		Retrieves all profiles NOT assigned to the parameter.
		:ref:`to-api-parameters-id-unassigned_profiles`
		:param parameter_id: The parameter id
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'get', u'profiles/{id:d}/parameters', (u'1.1', u'1.2', u'1.3',))
	def get_parameters_by_profile_id(self, profile_id=None):
		"""
		Get all Parameters associated with a Profile by Id.
		:ref:`to-api-profiles-id-parameters`
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'profiles/{id:d}/unassigned_parameters', (u'1.1', u'1.2', u'1.3',))
	def get_unnassigned_parameters_by_profile_id(self, profile_id=None):
		"""
		Get all Parameters associated with a Profile by Id.
		:ref:`to-api-profiles-id-unassigned_parameters`
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'profiles/name/{profile_name}/parameters', (u'1.1', u'1.2', u'1.3',))
	def get_parameters_by_profile_name(self, profile_name=None):
		"""
		Get all Parameters associated with a Profile by Name.
		:ref:`to-api-profiles-name-name-parameters`
		:param profile_name: The profile name
		:type profile_name: str
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'parameters', (u'1.1', u'1.2', u'1.3',))
	def create_parameter(self, data=None):
		"""
		Create Parameter
		:ref:`to-api-parameters`
		:param data: The parameter(s) data to use for parameter creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'parameters/{parameter_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_parameter(self, parameter_id=None, query_params=None):
		"""
		Update Parameter
		:ref:`to-api-parameters-id`
		:param parameter_id: The parameter id to update
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'delete', u'parameters/{parameter_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_parameter(self, parameter_id=None):
		"""
		Delete Parameter
		:ref:`to-api-parameters-id`
		:param parameter_id: The parameter id to delete
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'parameters/validate', (u'1.1', u'1.2', u'1.3',))
	def validate_parameter_exists(self, data=None):
		"""
		Validate that a Parameter exists.
		:ref:`to-api-parameters-validate`
		:param data: The parameter data to use for parameter validation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Physical Location
	#
	@api_request(u'get', u'phys_locations', (u'1.1', u'1.2', u'1.3',))
	def get_physical_locations(self, query_params=None):
		"""
		Get Physical Locations.
		:ref:`to-api-phys_locations`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'phys_locations/trimmed', (u'1.1', u'1.2', u'1.3',))
	def get_trimmed_physical_locations(self):
		"""
		Get Physical Locations with name only
		:ref:`to-api-phys_locations-trimmed`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'phys_locations/{physical_location_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_physical_location_by_id(self, physical_location_id=None):
		"""
		Get Physical Location by id
		:ref:`to-api-phys_locations-id`
		:param physical_location_id: The id to retrieve
		:type physical_location_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'phys_locations/{physical_location_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_physical_location(self, physical_location_id=None, query_params=None):
		"""
		Update Physical Location by id
		:ref:`to-api-phys_locations-id`
		:param physical_location_id: The id to update
		:type physical_location_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'regions/{region_name:s}/phys_locations', (u'1.1', u'1.2', u'1.3',))
	def create_physical_location(self, region_name=None, query_params=None):
		"""
		Create physical location
		:ref:`to-api-regions-name-phys_locations`
		:param region_name: the name of the region to create physical location into
		:type region_name: String
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'phys_locations/{physical_location_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_physical_location(self, physical_location_id=None, query_params=None):
		"""
		Delete Physical Location by id
		:ref:`to-api-phys_locations-id`
		:param physical_location_id: The id to delete
		:type physical_location_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Profiles
	#
	@api_request(u'get', u'profiles', (u'1.1', u'1.2', u'1.3',))
	def get_profiles(self, query_params=None):
		"""
		Get Profiles.
		:ref:`to-api-profiles`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'profiles/trimmed', (u'1.1', u'1.2', u'1.3',))
	def get_trimmed_profiles(self):
		"""
		Get Profiles with names only
		:ref:`to-api-profiles-trimmed`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_profile_by_id(self, profile_id=None):
		"""
		Get Profile by Id.
		:ref:`to-api-profiles-id`
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'profiles', (u'1.1', u'1.2', u'1.3',))
	def create_profile(self, data=None):
		"""
		Create a profile
		:ref:`to-api-profiles`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'profiles/name/{new_profile_name:s}/copy/{copy_profile_name:s}',
	             (u'1.1', u'1.2', u'1.3',))
	def copy_profile(self, new_profile_name=None, copy_profile_name=None, data=None):
		"""
		Copy profile to a new profile. The new profile name must not exist
		:ref:`to-api-profiles-name-name-copy-copy`
		:param new_profile_name: The name of profile to copy to
		:type new_profile_name: String
		:param copy_profile_name: The name of profile copy from
		:type copy_profile_name: String
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_profile_by_id(self, profile_id=None, data=None):
		"""
		Update Profile by Id.
		:ref:`to-api-profiles-id`
		:param profile_id: The profile Id
		:type profile_id: int
		:param data: The parameter data to edit
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_profile_by_id(self, profile_id=None):
		"""
		Delete Profile by Id.
		:ref:`to-api-profiles-id`
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Profile Parameters
	#
	@api_request(u'post', u'profileparameters', (u'1.1', u'1.2', u'1.3',))
	def associate_paramater_to_profile(self, data=None):
		"""
		Associate parameter to profile.
		:ref:`to-api-profileparameters`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'profiles/{profile_id:d}/parameters', (u'1.1', u'1.2', u'1.3',))
	def associate_parameters_by_profile_id(self, profile_id=None, data=None):
		"""
		Associate Parameters to a Profile by Id.
		:ref:`to-api-profiles-id-parameters`
		:param profile_id: The profile id
		:type profile_id: int
		:param data: The parameter data to associate
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'profileparameter', (u'1.1', u'1.2', u'1.3',))
	def assign_profile_to_parameter_ids(self, data=None):
		"""
		Create one or more profile / parameter assignments.
		:ref:`to-api-profileparameter`
		:param data: The data to assign
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'parameterprofile', (u'1.1', u'1.2', u'1.3',))
	def assign_parameter_to_profile_ids(self, data=None):
		"""
		Create one or more parameter / profile assignments.
		:ref:`to-api-profileparameter`
		:param data: The data to assign
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request(u'post', u'profiles/name/{profile_name}/parameters', (u'1.1', u'1.2', u'1.3',))
	def associate_parameters_by_profile_name(self, profile_name=None, data=None):
		"""
		Associate Parameters to a Profile by Name.
		:ref:`to-api-profiles-name-name-parameters`
		:param profile_name: The profile name
		:type profile_name: str
		:param data: The parameter data to associate
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'profileparameters/{profile_id:d}/{parameter_id:d}',
	             (u'1.1', u'1.2', u'1.3',))
	def delete_profile_parameter_association_by_id(self, profile_id=None, parameter_id=None):
		"""
		Delete Parameter association by Id for a Profile by Id.
		:ref:`to-api-profileparameters-profileID-parameterID`
		:param profile_id: The profile id
		:type profile_id: int
		:param parameter_id: The parameter id
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Regions
	#
	@api_request(u'get', u'regions', (u'1.1', u'1.2', u'1.3',))
	def get_regions(self):
		"""
		Get Regions.
		:ref:`to-api-regions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'regions/{region_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_region_by_id(self, region_id=None):
		"""
		Get Region by ID
		:ref:`to-api-regions-id`
		:param region_id: The region id of the region to retrieve
		:type region_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'regions/{region_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_region(self, region_id=None):
		"""
		Update a region
		:ref:`to-api-regions-id`
		:parma region_id: The region to update
		:type region_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'divisions/{division_name:s}/regions', (u'1.1', u'1.2', u'1.3',))
	def create_region(self, division_name=None, data=None):
		"""
		Create a region
		:ref:`to-api-divisions-name-regions`
		:param division_name: The Division name in which region will reside
		:type division_name: String
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Roles
	#
	@api_request(u'get', u'roles', (u'1.1', u'1.2', u'1.3',))
	def get_roles(self):
		"""
		Get Roles.
		:ref:`to-api-roles`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Server
	#
	@api_request(u'get', u'servers', (u'1.1', u'1.2', u'1.3', u'1.4'))
	def get_servers(self, query_params=None):
		"""
		Get Servers.
		:ref:`to-api-servers`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/{server_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_server_by_id(self, server_id=None):
		"""
		Get Server by Server ID
		:ref:`to-api-servers-id`
		:param server_id: The server id to retrieve
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/{server_id:d}/deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def get_server_delivery_services(self, server_id=None):
		"""
		Retrieves all delivery services assigned to the server
		:ref:`to-api-servers-id-deliveryservices`
		:param server_id: The server id to retrieve
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/totals', (u'1.1', u'1.2', u'1.3',))
	def get_server_type_count(self):
		"""
		Retrieves a count of CDN servers by type
		:ref:`to-api-servers-totals`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/status', (u'1.1', u'1.2', u'1.3',))
	def get_server_status_count(self):
		"""
		Retrieves a count of CDN servers by status
		:ref:`to-api-servers-status`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/hostname/{name}/details', (u'1.1', u'1.2', u'1.3',))
	def get_server_details(self, name=None):
		"""
		Get server details from trafficOps
		:ref:`to-api-servers-hostname-name-details`
		:param hostname: Server hostname
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'servercheck', (u'1.1', u'1.2', u'1.3',))
	def create_servercheck(self, data=None):
		"""
		Post a server check result to the serverchecks table.
		:ref:`to-api-servercheck`
		:param data: The parameter data to use for server creation
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'servers', (u'1.1', u'1.2', u'1.3',))
	def create_server(self, data=None):
		"""
		Create a new Server.
		:ref:`to-api-servers`
		:param data: The parameter data to use for server creation
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'servers/{server_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_server_by_id(self, server_id=None, data=None):
		"""
		Update a Server by Id.
		:ref:`to-api-servers-id`
		:param server_id: The server Id
		:type server_id: int
		:param data: The parameter data to edit
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""
	@api_request(u'put', u'servers/{server_id:d}/status', (u'1.1', u'1.2', u'1.3',))
	def update_server_status_by_id(self, server_id=None, data=None):
		"""
		Update server_status by Id.
		:ref:`to-api-servers-id-status`
		:param server_id: The server Id
		:type server_id: int
		:status: https://traffic-control-cdn.readthedocs.io/en/latest/api/server.html
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'servers/{server_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_server_by_id(self, server_id=None):
		"""
		Delete a Server by Id.
		:ref:`to-api-servers-id`
		:param server_id: The server Id
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'servers/{server_id:d}/queue_update', (u'1.1', u'1.2', u'1.3',))
	def servers_queue_update(self, server_id=None, data=None):
		"""
		Queue Updates by Server Id.
		:ref:`to-api-servers-id-queue_update`
		:param server_id: The server Id
		:type server_id: int
		:param data: The update action.  QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'servers/{server_name}/update_status', (u'1.1', u'1.2', u'1.3', u'1.4'))
	def get_server_update_status(self, server_name=None):
		"""
		Gets the current update status of a server named ``server_name``.
		:ref:`to-api-servers-hostname-update_status`
		:param server_name: The (short) hostname of the server for which the update status will be fetched
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Static DNS Entries
	#
	@api_request(u'get', u'staticdnsentries', (u'1.1', u'1.2', ))
	def get_static_dns_entries(self):
		"""
		Get Static DNS Entries.
		:ref:`to-api-staticdnsentries`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'staticdnsentries', (u'1.1', u'1.2', u'1.3',))
	def get_staticdnsentries(self, query_params=None):
		"""
		Get static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'staticdnsentries', (u'1.3',))
	def create_staticdnsentries(self, data=None):
		"""
		Create static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'staticdnsentries', (u'1.3',))
	def update_staticdnsentries(self, data=None, query_params=None):
		"""
		Update static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'staticdnsentries', (u'1.3',))
	def delete_staticdnsentries(self, query_params=None):
		"""
		Delete static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Status
	#
	@api_request(u'get', u'statuses', (u'1.1', u'1.2', u'1.3', u'1.4'))
	def get_statuses(self):
		"""
		Retrieves a list of the server status codes available.
		:ref:`to-api-statuses`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'statuses/{status_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_statuses_by_id(self, status_id=None):
		"""
		Retrieves a server status by ID.
		:ref:`to-api-statuses-id`
		:param status_id: The status id to retrieve
		:type status_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# System
	#
	@api_request(u'get', u'system/info', (u'1.1', u'1.2', u'1.3',))
	def get_system_info(self):
		"""
		Get information on the traffic ops system.
		:ref:`to-api-system-info`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Tenants
	#
	@api_request(u'get', u'tenants', (u'1.1', u'1.2', u'1.3',))
	def get_tenants(self):
		"""
		Get all tenants.
		:ref:`to-api-tenants`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'tenants/{tenant_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_tenant_by_id(self, tenant_id=None):
		"""
		Get a tenant by ID.
		:ref:`to-api-tenants-id`
		:param tenant_id: The tenant to retrieve
		:type tenant_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'tenants/{tenant_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_tenant(self, tenant_id=None):
		"""
		Update a tenant
		:ref:`to-api-tenants-id`
		:param tenant_id: The tenant to update
		:type tenant_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'tenants', (u'1.1', u'1.2', u'1.3',))
	def create_tenant(self, data=None):
		"""
		Create a tenant
		:ref:`to-api-tenants`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# TO Extensions
	#
	@api_request(u'get', u'to_extensions', (u'1.1', u'1.2', u'1.3',))
	def get_to_extensions(self):
		"""
		Retrieves the list of extensions.
		:ref:`to-api-to_extensions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'to_extensions', (u'1.1', u'1.2', u'1.3',))
	def create_to_extension(self, data=None):
		"""
		Creates a Traffic Ops extension.
		:ref:`to-api-to_extensions`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'to_extensions/{extension_id:d}/delete', (u'1.1', u'1.2', u'1.3',))
	def delete_to_extension(self, extension_id=None):
		"""
		Deletes a Traffic Ops extension.
		:ref:`to-api-to_extensions-id-delete`
		:param extension_id: The extension id to delete
		:type extension_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Types
	#
	@api_request(u'get', u'types', (u'1.1', u'1.2', u'1.3',))
	def get_types(self, query_params=None):
		"""
		Get Data Types.
		:ref:`to-api-types`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'types/trimmed', (u'1.1', u'1.2', u'1.3',))
	def get_types_only_names(self):
		"""
		Get Data Types with only the Names
		:ref:`to-api-types-trimmed`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'types/{type_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_type_by_id(self, type_id=None):
		"""
		Get Data Type with the given type id
		:ref:`to-api-types-id`
		:param type_id: The ID of the type to retrieve
		:type type_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Users
	#
	@api_request(u'get', u'users', (u'1.1', u'1.2', u'1.3',))
	def get_users(self):
		"""
		Retrieves all users.
		:ref:`to-api-users`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'users/{user_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_user_by_id(self, user_id=None):
		"""
		Retrieves user by ID.
		:ref:`to-api-users-id`
		:param user_id: The user to retrieve
		:type user_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'users', (u'1.1', u'1.2', u'1.3',))
	def create_user(self, data=None):
		"""
		Create a user.
		:ref:`to-api-users`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'users/register', (u'1.1', u'1.2', u'1.3',))
	def create_user_with_registration(self, data=None):
		"""
		Register a user and send registration email
		:ref:`to-api-users-register`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'users/{user_id:d}/deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def get_user_delivery_services(self, user_id=None):
		"""
		Retrieves all delivery services assigned to the user.
		:ref:`to-api-users-id-deliveryservices`
		:param user_id: The user to retrieve
		:type user_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'user/current', (u'1.1', u'1.2', u'1.3',))
	def get_authenticated_user(self):
		"""
		Retrieves the profile for the authenticated user.
		:ref:`to-api-user-current`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'user/current/jobs', (u'1.1', u'1.2', u'1.3',))
	def get_authenticated_user_jobs(self):
		"""
		Retrieves the user's list of jobs.
		:ref:`to-api-user-current-jobs`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'user/current/jobs', (u'1.1', u'1.2', u'1.3',))
	def create_invalidation_job(self, data=None):
		"""
		Invalidating content on the CDN is sometimes necessary when the origin was mis-configured
		and something is cached in the CDN that needs to be removed. Given the size of a typical
		Traffic Control CDN and the amount of content that can be cached in it, removing the content
		from all the caches may take a long time. To speed up content invalidation, Traffic Ops will
		not try to remove the content from the caches, but it makes the content inaccessible using
		the `regex_revalidate ATS plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_revalidate.en.html>`_
		This forces a revalidation of the content, rather than a new get.
		:ref:`to-api-user-current-jobs`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Snapshot CRConfig
	#
	@api_request(u'get', u'cdns/{cdn_name}/snapshot', (u'1.2', u'1.3',))
	def get_current_snapshot_crconfig(self, cdn_name=None):
		"""
		Retrieves the CURRENT snapshot for a CDN which doesn't necessarily represent the current
		state of the CDN. The contents of this snapshot are currently used by Traffic Monitor and
		Traffic Router.
		:ref:`to-api-cdns-name-snapshot`
		:param cdn_name: The CDN name
		:type cdn_name: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'get', u'cdns/{cdn_name}/snapshot/new', (u'1.2', u'1.3',))
	def get_pending_snapshot_crconfig(self, cdn_name=None):
		"""
		Retrieves a PENDING snapshot for a CDN which represents the current state of the CDN. The
		contents of this snapshot are NOT currently used by Traffic Monitor and Traffic Router. Once
		a snapshot is performed, this snapshot will become the CURRENT snapshot and will be used by
		Traffic Monitor and Traffic Router.
		:ref:`to-api-cdns-name-snapshot-new`
		:param cdn_name: The CDN name
		:type cdn_name: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'snapshot/{cdn_name}', (u'1.2', u'1.3',))
	def snapshot_crconfig(self, cdn_name=None):
		"""
		Snapshot CRConfig by CDN Name.
		:ref:`to-api-snapshot-name`
		:param cdn_name: The CDN name
		:type cdn_name: str
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""



	#
	# Coordinate
	#
	@api_request(u'get', u'coordinates', (u'1.3',))
	def get_coordinates(self, query_params=None):
		"""
		Get all coordinates associated with the cdn
		:ref:`to-api-coordinates`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'coordinates', (u'1.3'))
	def create_coordinates(self, data=None):
		"""
		Create coordinates
		:ref:`to-api-coordinates`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'coordinates', (u'1.3'))
	def update_coordinates(self, query_params=None, data=None):
		"""
		Update coordinates
		:ref:`to-api-coordinates`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'coordinates', (u'1.3'))
	def delete_coordinates(self, query_params=None):
		"""
		Delete coordinates
		:ref:`to-api-coordinates`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Origin
	#
	@api_request(u'get', u'origins', (u'1.3',))
	def get_origins(self, query_params=None):
		"""
		Get origins associated with the delivery service
		:ref:`to-api-origins`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'post', u'origins', (u'1.3',))
	def create_origins(self, data=None):
		"""
		Creates origins associated with a delivery service
		:ref:`to-api-origins`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'put', u'origins', (u'1.3',))
	def update_origins(self, query_params=None):
		"""
		Updates origins associated with a delivery service
		:ref:`to-api-origins`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request(u'delete', u'origins', (u'1.3',))
	def delete_origins(self, query_params=None):
		"""
		Updates origins associated with a delivery service
		:ref:`to-api-origins`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	####################################################################################
	####                                                                            ####
	####                          Data Model Overrides                              ####
	####                                                                            ####
	####################################################################################

	def __enter__(self):
		"""
		Implements context-management for ToSessions. This will open the session by sending a
		connection request immediately, rather than waiting for login.

		:returns: The constructed object (:meth:`__init__` is called implicitly prior to this method)
		"""
		self.create()
		return self

	def __exit__(self, exc_type, exc_value, traceback):
		"""
		Implements context-management for TOSessions. This will close the underlying socket.
		"""
		self.close()

		if exc_type:
			logging.error("%s", exc_value)
			logging.debug("%s", exc_type, stack_info=traceback)




if __name__ == u'__main__':
	# Sample usages
	import operator

	DEBUG = False

	logging.basicConfig(stream=sys.stderr, level=logging.INFO if not DEBUG else logging.DEBUG)

	# TOSession Class Examples
	#     TOSession is a class that allows you to create a session to a Traffic Ops instance
	#     and interact with the Traffic Ops API.

	# Traffic Ops System - for self-signed cert -> turn off cert verification
	TOS = TOSession(host_ip=u'to.somedomain.net', verify_cert=True)
	TOS.login(u'someuser', u'someuser123')

	# Objects get returned munch-ified by default which means you can access dictionary keys as
	# attributes names but you can still access the entries with keys as well e.g.
	# ``cdn.name`` is equivalent to ``cdn['name']``
	CDNS = TOS.get_cdns()[0]
	print(CDNS)
	for cdn in CDNS:
		print(u'CDN [{0}] has id [{1}]'.format(cdn.name, cdn.id))

	ALL_TYPES = TOS.get_types()[0]
	print(u'All Types are (sorted by useInTable, name):')
	print(ALL_TYPES)
	for atype in sorted(ALL_TYPES, key=operator.itemgetter(u'useInTable', u'name')):
		print(u'Type [{0}] for table [{1}]'.format(atype.name, atype.useInTable))

	print(u'Getting all cache groups (bulk)...')
	CACHE_GROUPS = TOS.get_cachegroups()[0]
	for cache_group in CACHE_GROUPS:
		print(u'Bulk cache group [{0}] has id [{1}]'.format(cache_group.name, cache_group.id))

		# Example with URL replacement parameters
		# e.g. TOSession.get_cachegroups_by_id() == end-point 'api/1.2/cachegroups/{id}'
		#      See TOSession object for details.
		print(u'    Getting cachegroup by id [{0}]'.format(cache_group.id),
			  u' to demonstrate getting by id...')
		cg_id_list = TOS.get_cachegroup_by_id(cache_group_id=cache_group.id)[0]
		print(u'    Cache group [{0}] by id [{1}]'.format(cg_id_list[0].name, cg_id_list[0].id))

	# Example with URL query parameters
	SERVER_TYPES = TOS.get_types(query_params={u'useInTable': u'server'})[0]
	print(u'Server Types are:')
	print(SERVER_TYPES)
	for stype in SERVER_TYPES:
		print(u'Type [{0}] for table [{1}]'.format(stype.name, stype.useInTable))
	TOS.close()
	print(u'Done!')
