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

Requires Python Version >= 3.6
"""

# Core Modules
import logging
import sys
from requests import Response
from typing import Any, Dict, List, Tuple, Union

# Third-party Modules
import munch
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

		@api_request('get', 'cdns', ('3.0',))
		def get_cdns(self):
			pass


	.. code-block:: python3
		:caption: End-point with URL parameters and no query parameters

		@api_request('get', 'cdns/{cdn_id:d}', ('3.0',))
		def get_cdn_by_id(self, cdn_id=None):
			pass


	.. code-block:: python3
		:caption: End-point with no URL parameters but with query parameters

		@api_request('get', 'deliveryservices', ('3.0',))
		def get_deliveryservices(self, query_params=None):
			pass

	.. code-block:: python3
		:caption: End-point with URL parameters and query parameters

		@api_request('get', 'deliveryservices/xmlId/{xml_id}/sslkeys', ('3.0',))
		def get_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
			pass

	.. code-block:: python3
		:caption: End-point with request data

		@api_request('post', 'cdns', ('3.0',))
		def create_cdn(self, data=None):
			pass

	.. code-block:: python3
		:caption: End-point with URL parameters and request data

		@api_request('put', 'cdns', ('3.0',))
		def update_cdn_by_id(self, cdn_id=None, data=None):
			pass

	Calling end-point methods

	:meth:`get_cdns` calls endpoint :ref:`to-api-cdns` e.g. ``t.get_cdns()``

	:meth:`get_types` calls endpoint :ref:`to-api-types`, optionally with query parameters e.g. ``get_foo_data(id=45, query_params={'sort': 'asc'})`` calls endpoint ``GET api/2.x/foo/45?sort=asc`` (presumably)

	:meth:`cdns_queue_update` calls endpoint :ref:`to-api-cdns-id-queue_update`, with an ID path parameter and a JSON payload e.g. ``cdns_queue_update(id=1, data={'action': 'queue'})``

	.. note:: Only a small subset of the API endpoints are implemented. More can be implemented as needed.
	"""

	def __init__(self, host_ip, host_port=443, api_version='4.1', ssl=True, headers=None,
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
		                                api_base_path='api/{api_version}/',
		                                host_port=host_port, ssl=ssl, headers=headers,
		                                verify_cert=verify_cert)

		self._logged_in = False

		msg = 'TOSession instance {0:#0x} initialized: Details: {1}'
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
			self.post('user/login', data={'u': username, 'p': password})
			self._logged_in = True
		except rex.SSLError as e:
			logging.debug("%s", e, stack_info=True, exc_info=True)
			self.close()
			msg = ('{0}.  This system may have a self-signed certificate.  Try creating this'
			       ' TOSession object passing verify_cert=False. e.g. TOSession(..., '
			       'verify_cert=False).')
			msg = msg.format(e)
			logging.error(msg)
			logging.warning("disabling certificate verification is not recommended.")
			raise LoginError(msg) from e
		except OperationError as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			msg = 'Logging in to Traffic Ops has failed. Reason: {0}'.format(e)
			self.close()
			logging.error(msg)
			raise OperationError(msg) from e

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

		:return: The base url should match '[\\w\\+\\-\\.]+://[\\w\\+\\-\\.]+(:\\d+)?' e.g https://to.somedomain.net/api/4.0/
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
		if 'munchify' in kwargs:
			munchify = kwargs['munchify']

		while True:
			data, resp = self.get_deliveryserviceserver(query_params={'limit':limit, 'page': page},
			                                           *args, **kwargs)

			if not data:
				break

			result_set.extend(munch.munchify(data) if munchify else data)
			page += 1

		return result_set, resp  # Note: Return last response object received

	#
	# PUT ALL API DEFINITIONS BELOW AND UNDER ITS RESPECTIVE PAGE (whether it is 2.0 or 2.1, etc, if its
	# a CDN put it under CDN header and corresponding calls)
	#

	#
	#	API Capabilities
	#
	@api_request('get', 'api_capabilities', ('3.0',))
	def get_api_capabilities(self, query_params=None):
		"""
		Get all API-capability mappings
		:ref:`to-api-v3-api_capabilities`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# ASN
	#
	@api_request('get', 'asns', ('3.0', '4.0', '4.1', '5.0'))
	def get_asns(self, query_params=None):
		"""
		Get ASNs.
		:ref:`to-api-asns`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'asns', ('3.0', '4.0', '4.1', '5.0'))
	def create_asn(self, data=None):
		"""
		Create ASN
		:ref:`to-api-asns`
		:param data: The ASN data to use for ASN creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'asns', ('3.0', '4.0', '4.1', '5.0'))
	def update_asn(self, query_params=None):
		"""
		Update ASN
		:ref:`to-api-asns-id`
		:param asn_id: The ID of the ASN to update
		:type asn_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'asns', ('3.0', '4.0', '4.1', '5.0'))
	def delete_asn(self, query_params=None):
		"""
		Delete ASN
		:to-api-asns-id:
		:param asn_id: The ID of the ASN to delete
		:type asn_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Statistics
	#
	@api_request('get', 'cache_stats', ('3.0', '4.0', '4.1', '5.0'))
	def get_cache_stats(self, query_params=None):
		"""
		Retrieves statistics about the CDN.
		:ref:`to-api-cache_stats`
		:param query_params: See API page for more information on accepted params
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'caches/stats', ('3.0', '4.0', '4.1', '5.0'))
	def get_traffic_monitor_cache_stats(self):
		"""
		Retrieves cache stats from Traffic Monitor. Also includes rows for aggregates
		:ref:`to-api-caches-stats`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Cache Groups
	#
	@api_request('get', 'cachegroups', ('3.0', '4.0', '4.1', '5.0'))
	def get_cachegroups(self, query_params=None):
		"""
		Get Cache Groups.
		:ref:`to-api-cachegroups`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'cachegroups/{cache_group_id:d}/parameters', ('3.0',))
	def get_cachegroup_parameters(self, cache_group_id=None):
		"""
		Get a cache groups parameters. This endpoint has been
		deprecated and will no longer be available as of TO API v4.
		:ref:`to-api-v3-cachegroups-id-parameters`
		:param cache_group_id: The cache group Id
		:type cache_group_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'cachegroupparameters', ('3.0',))
	def get_all_cachegroup_parameters(self):
		"""
		A collection of all cache group parameters. This endpoint has been
		deprecated and will no longer be available as of TO API v4.
		:ref:`to-api-v3-cachegroupparameters`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cachegroups', ('4.0', '4.1', '5.0',))
	def create_cachegroups(self, data=None):
		"""
		Create a Cache Group
		:ref:`to-api-cachegroups`
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'cachegroups/{cache_group_id:d}', ('4.0', '4.1', '5.0',))
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

	@api_request('delete', 'cachegroups/{cache_group_id:d}', ('4.0', '4.1', '5.0',))
	def delete_cachegroups(self, cache_group_id=None):
		"""
		Delete a cache group
		:ref:`to-api-cachegroups-id`
		:param cache_group_id: The cache group id to update
		:type cache_group_id: Integer
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cachegroups/{cache_group_id:d}/queue_update', ('4.0', '4.1', '5.0',))
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

	@api_request('delete', 'cachegroupparameters/{cache_group_id:d}/{parameter_id:d}', ('3.0',))
	def delete_cache_group_parameters(self, cache_group_id=None, parameter_id=None):
		"""
		Delete a cache group parameter association. This endpoint has been
		deprecated and will no longer be available as of TO API v4.
		:ref:`to-api-v3-cachegroupparameters-id-parameterID`
		:param cache_group_id: The cache group id in which the parameter will be deleted
		:type cache_group_id: int
		:param parameter_id: The parameter id which will be disassociated
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Capabilities
	#
	@api_request('get', 'capabilities', ('3.0',))
	def get_capabilities(self, query_params=None):
		"""
		Retrieves capabilities
		:ref:`to-api-v3-capabilities`
		:param query_params: See API page for more information on accepted parameters
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN
	#
	@api_request('get', 'cdns', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdns(self, query_params=None):
		"""
		Get all CDNs.
		:ref:`to-api-cdns`
		:param query_params: See API page for more information on accepted parameters
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cdns', ('3.0', '4.0', '4.1', '5.0'))
	def create_cdn(self, data=None):
		"""
		Create a new CDN.
		:ref:`to-api-cdns`
		:param data: The parameter data to use for cdn creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'cdns/{cdn_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'cdns/{cdn_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_cdn_by_id(self, cdn_id=None):
		"""
		Delete a CDN by Id.
		:ref:`to-api-cdns-id`
		:param cdn_id: The CDN id
		:type cdn_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cdns/{cdn_id:d}/queue_update', ('3.0', '4.0', '4.1', '5.0'))
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
	# CDN Health/Usage
	#
	@api_request('get', 'cdns/health', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdns_health(self):
		"""
		Retrieves the health of all locations (cache groups) for all CDNs
		:ref:`to-api-cdns-health`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('get', 'cdns/{cdn_name:s}/health', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdn_health_by_name(self, cdn_name=None):
		"""
		Retrieves the health of all locations (cache groups) for a given CDN
		:ref:`to-api-cdns-name-health`
		:param cdn_name: The CDN name to find health for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'cdns/capacity', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'cdns/routing', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'cdns/domains', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdns_domains(self):
		"""
		Retrieves the different CDN domains
		:ref:`to-api-cdns-domains`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Notifications
	#
	@api_request('get', 'cdn_notifications', ('4.0', '4.1', '5.0'))
	def get_cdn_notifications(self, query_params=None):
		"""
		Get all CDN notifications.
		:ref:`to-api-cdn-notifications`
		:param query_params: See API page for more information on accepted parameters
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cdn_notifications', ('4.0', '4.1', '5.0'))
	def create_cdn_notification(self, data=None):
		"""
		Create a new CDN notification.
		:ref:`to-api-cdn-notifications`
		:param data: The parameter data to use for cdn creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'cdn_notifications', ('4.0', '4.1', '5.0'))
	def delete_cdn_notification(self, query_params=None):
		"""
		Deletes an existing CDN notification.
		:ref:`to-api-cdn-notifications`
		:param query_params: See API page for more information on accepted parameters
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN Topology
	#

	@api_request('get', 'cdns/{cdn_name:s}/configs/monitoring', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdn_monitoring_info(self, cdn_name=None):
		"""
		Retrieves CDN monitoring information
		:ref:`to-api-cdns-name-configs-monitoring`
		:param cdn_name: The CDN name to find configs for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# DNSSEC Keys
	#
	@api_request('get', 'cdns/name/{cdn_name:s}/dnsseckeys', ('3.0', '4.0', '4.1', '5.0'))
	def get_cdn_dns_sec_keys(self, cdn_name=None):
		"""
		Gets a list of dnsseckeys for a CDN and all associated Delivery Services
		:ref:`to-api-cdns-name-name-dnsseckeys`
		:param cdn_name: The CDN name to find dnsseckeys info for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'cdns/name/{cdn_name:s}/dnsseckeys', ('3.0', '4.0', '4.1', '5.0'))
	def delete_cdn_dns_sec_keys(self, cdn_name=None):
		"""
		Delete dnssec keys for a cdn and all associated delivery services
		:ref:`to-api-cdns-name-name-dnsseckeys`
		:param cdn_name: The CDN name to delete dnsseckeys info for
		:type cdn_name: String
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cnds/dnsseckeys/generate', ('3.0', '4.0', '4.1', '5.0'))
	def create_cdn_dns_sec_keys(self, data=None):
		"""
		Generates ZSK and KSK keypairs for a CDN and all associated Delivery Services
		:ref:`to-api-cdns-dnsseckeys-generate`
		:param data: The parameter data to use for cachegroup creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# CDN SSL Keys
	#
	@api_request('get', 'cdns/name/{cdn_name:s}/sslkeys', ('3.0', '4.0', '4.1', '5.0'))
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
	# CDN Lock
	#
	@api_request('get', 'cdn_locks', ('4.0', '4.1', '5.0'))
	def get_cdn_locks(self, query_params=None):
		"""
		Returns information for all CDN locks.
		:ref:`to-api-cdn-locks`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'cdn_locks', ('4.0', '4.1', '5.0'))
	def create_cdn_lock(self, data=None):
		"""
		Create a new CDN lock.
		:ref:`to-api-cdn-locks`
		:param data: The parameter data to use for cdn_lock creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'cdn_locks', ('4.0', '4.1', '5.0'))
	def delete_cdn_lock(self, query_params=None):
		"""
		Deletes an existing CDN Lock.
		:ref:`to-api-cdn-locks`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Change Logs
	#
	@api_request('get', 'logs', ('3.0', '4.0', '4.1', '5.0'))
	def get_change_logs(self, query_params=None):
		"""
		Retrieve all change logs from traffic ops
		:ref:`to-api-logs`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'logs/newcount', ('3.0', '4.0', '4.1', '5.0'))
	def get_change_logs_newcount(self):
		"""
		Get amount of new logs from traffic ops
		:ref:`to-api-logs-newcount`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service
	#
	@api_request('get', 'deliveryservices', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservices(self, query_params=None):
		"""
		Retrieves all delivery services (if admin or ops) or all delivery services assigned to user.
		:ref:`to-api-deliveryservices`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices', ('3.0', '4.0', '4.1', '5.0'))
	def create_deliveryservice(self, data=None):
		"""
		Allows user to create a delivery service.
		:ref:`to-api-deliveryservices`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'deliveryservices/{delivery_service_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('put', 'deliveryservices/{delivery_service_id:d}/safe', ('3.0', '4.0', '4.1', '5.0'))
	def update_deliveryservice_safe(self, delivery_service_id=None, data=None):
		"""
		Allows a user to edit limited fields of a Delivery Service.
		:ref:`to-api-deliveryservices-id-safe`
		:param delivery_service_id: The Delivery Service Id
		:type delivery_service_id: int
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('delete', 'deliveryservices/{delivery_service_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	# Delivery Service Requests
	#
	@api_request('get', 'deliveryservice_requests', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservice_requests(self, query_params=None):
		"""
		Retrieves all delivery service reuests.
		:ref:`to-api-deliveryservice-requests`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservice_requests', ('3.0', '4.0', '4.1', '5.0'))
	def create_deliveryservice_request(self, data=None):
		"""
		Creates a new delivery service request.
		:ref:`to-api-deliveryservice-requests`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'deliveryservice_requests', ('3.0', '4.0', '4.1', '5.0'))
	def update_deliveryservice_request(self, query_params=None, data=None):
		"""
		Updates an existing Delivery Service Request.
		:ref:`to-api-deliveryservice-requests`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'deliveryservice_requests', ('3.0', '4.0', '4.1', '5.0'))
	def delete_deliveryservice_request(self, query_params=None):
		"""
		Deletes a Delivery Service Request.
		:ref:`to-api-deliveryservice-requests`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Delivery Service Request Comments
	#
	@api_request('get', 'deliveryservice_request_comments', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservice_request_comments(self, query_params=None):
		"""
		Retrieves all delivery service reuest comments.
		:ref:`to-api-deliveryservice_request_comments`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservice_request_comments', ('3.0', '4.0', '4.1', '5.0'))
	def create_deliveryservice_request_comment(self, data=None):
		"""
		Creates a new delivery service request comment.
		:ref:`to-api-deliveryservice_request_comments`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'deliveryservice_request_comments', ('3.0', '4.0', '4.1', '5.0'))
	def update_deliveryservice_request_comment(self, query_params=None, data=None):
		"""
		Updates an existing Delivery Service Request comment.
		:ref:`to-api-deliveryservice_request_comments`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'deliveryservice_request_comments', ('3.0', '4.0', '4.1', '5.0'))
	def delete_deliveryservice_request_comment(self, query_params=None):
		"""
		Deletes a Delivery Service Request comment.
		:ref:`to-api-deliveryservice_request_comments`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

    #
	# Delivery Service Required capabilities
	#
	@api_request('get', 'deliveryservices_required_capabilities', ('3.0', '4.0', '4.1'))
	def get_deliveryservices_required_capabilities(self, query_params=None):
		"""
		Retrieves all delivery service required capabilities.
		:ref:`to-api-v4-deliveryservices_required_capabilities`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices_required_capabilities', ('3.0', '4.0', '4.1'))
	def create_deliveryservices_required_capabilities(self, data=None):
		"""
		Creates a new delivery service required capability.
		:ref:`to-api-v4-deliveryservices_required_capabilities`
		:param data: The request data structure for the API request
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'deliveryservices_required_capabilities', ('3.0', '4.0', '4.1'))
	def delete_deliveryservices_required_capabilities(self, query_params=None):
		"""
		Deletes a Delivery Service Required capability.
		:ref:`to-api-v4-deliveryservices_required_capabilities`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Health
	#
	@api_request('get', 'deliveryservices/{delivery_service_id:d}/health', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('get', 'deliveryservices/{delivery_service_id:d}/capacity', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('get', 'deliveryservices/{delivery_service_id:d}/routing', ('3.0', '4.0', '4.1', '5.0'))
	def get_delivery_service_routing(self, delivery_service_id=None):
		"""
		Retrieves the aggregated routing percentages for a given Delivery Service.
		:ref:`to-api-deliveryservices-id-routing`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Delivery Service Server
	#
	@api_request('get', 'deliveryserviceserver', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryserviceserver(self, query_params=None):
		"""
		Retrieves delivery service / server assignments. (Allows pagination and limits)
		:ref:`to-api-deliveryserviceserver`
		:param query_params: The required url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryserviceserver', ('3.0', '4.0', '4.1', '5.0'))
	def assign_deliveryservice_servers_by_ids(self, data=None):
		"""
		Assign servers by id to a Delivery Service. (New Method)
		:ref:`to-api-deliveryserviceserver`
		:param data: The required data to create server associations to a delivery service
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices/{xml_id}/servers', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'deliveryserviceserver/{delivery_service_id:d}/{server_id:d}',('3.0', '4.0', '4.1', '5.0'))
	def delete_deliveryservice_servers_by_id(self, delivery_service_id=None, server_id=None):
		"""
		Removes a server (cache) from a delivery service.
		:ref:`to-api-deliveryserviceserver-dsid-serverid`
		:param delivery_service_id: The delivery service id
		:type delivery_service_id: int
		:param server_id: The server id to remove from delivery service
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'deliveryservices/{delivery_service_id:d}/servers', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservice_servers(self, delivery_service_id=None):
		"""
		Retrieves properties of CDN EDGE or ORG servers assigned to a delivery service.
		:ref:`to-api-deliveryservices-id-servers`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'deliveryservices/{delivery_service_id:d}/servers/eligible', ('3.0', '4.0', '4.1', '5.0'))
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

	#
	# Delivery Service SSL Keys
	#
	@api_request('get', 'deliveryservices/xmlId/{xml_id}/sslkeys', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'deliveryservices/xmlId/{xml_id}/sslkeys', ('3.0', '4.0', '4.1', '5.0'))
	def delete_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
		"""
		Delete SSL keys for a Delivery Service by xmlId.
		:ref:`to-api-deliveryservices-xmlid-xmlid-sslkeys`
		:param xml_id: The Delivery Service xmlId
		:type xml_id: str
		:param query_params: The url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices/sslkeys/generate', ('3.0', '4.0', '4.1', '5.0'))
	def generate_deliveryservice_ssl_keys(self, data=None):
		"""
		Generate an SSL certificate. (self-signed)
		:ref:`to-api-deliveryservices-sslkeys-generate`
		:param data: The parameter data to use for Delivery Service SSL key generation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices/sslkeys/add', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('post', 'deliveryservices/xmlId/{xml_id}/urlkeys/generate', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'deliveryservices_regexes', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservices_regexes(self):
		"""
		Get RegExes for all Delivery Services.
		:ref:`to-api-deliveryservices_regexes`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'deliveryservices/{delivery_service_id:d}/regexes', ('3.0', '4.0', '4.1', '5.0'))
	def get_deliveryservice_regexes_by_id(self, delivery_service_id=None, query_params=None):
		"""
		Get RegExes for a Delivery Service by Id.
		:ref:`to-api-deliveryservices-id-regexes`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param query_params: The url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'deliveryservices/{delivery_service_id:d}/regexes', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('put', 'deliveryservices/{delivery_service_id:d}/regexes/{regex_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_deliveryservice_regexes(self, delivery_service_id=None, regex_id=None,
	                                   data=None):
		"""
		Update a regex for a delivery service
		:ref:`to-api-deliveryservices-id-regexes-rid`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param regex_id: The delivery service regex id
		:type regex_id: int
		:param data: The required data to update delivery service regexes
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'deliveryservices/{delivery_service_id:d}/regexes/'
	                        '{delivery_service_regex_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'deliveryservice_stats', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'divisions', ('3.0', '4.0', '4.1', '5.0'))
	def get_divisions(self, query_params=None):
		"""
		Get all divisions.
		:ref:`to-api-divisions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'divisions/{division_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('post', 'divisions', ('3.0', '4.0', '4.1', '5.0'))
	def create_division(self, data=None):
		"""
		Create a division
		:ref:`to-api-divisions`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'divisions/{division_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'federations', ('3.0', '4.0', '4.1', '5.0'))
	def get_federations(self):
		"""
		Retrieves a list of federation mappings (aka federation resolvers) for a the current user
		:ref:`to-api-federations`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'federations', ('3.0', '4.0', '4.1', '5.0'))
	def create_federation(self, data=None):
		"""
		Allows a user to add federations for their delivery service(s).
		:ref:`to-api-federations`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'federations', ('3.0', '4.0', '4.1', '5.0'))
	def delete_federation(self):
		"""
		Allows a user to delete federations for their delivery service(s).
		:ref:`to-api-federations`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('get', 'cdns/{cdn_name:s}/federations', ('3.0', '4.0', '4.1', '5.0'))
	def get_federations_for_cdn(self, cdn_name=None, query_params=None):
		"""
		Retrieves a list of federations for a cdn.
		:ref:`to-api-cdns-name-federations`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'cdns/{cdn_name:s}/federations', ('3.0', '4.0', '4.1', '5.0'))
	def create_federation_in_cdn(self, cdn_name=None, data=None):
		"""
		Create a federation.
		:ref:`to-api-cdns-name-federations`
		:param cdn_name: The CDN name to find federation
		:type cdn_name: String
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'cdns/{cdn_name:s}/federations/{federation_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_federation_in_cdn(self, cdn_name=None, federation_id=None, data=None):
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

	@api_request('delete', 'cdns/{cdn_name:s}/federations/{federation_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'federations/{federation_id:d}/deliveryservices', ('3.0', '4.0', '4.1', '5.0'))
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


	@api_request('post', 'federations/{federation_id:d}/deliveryservices', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'federations/{federation_id:d}/federation_resolvers', ('3.0', '4.0', '4.1', '5.0'))
	def get_federation_resolvers_by_id(self, federation_id=None):
		"""
		:ref:`to-api-federations-id-federation_resolvers`
		Retrieves federation resolvers assigned to a federation
		:param federation_id: The federation id
		:type federation_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'federations/{federation_id:d}/federation_resolvers', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'federation_resolvers', ('3.0', '4.0', '4.1', '5.0'))
	def get_federation_resolvers(self, query_params=None):
		"""
		Get federation resolvers.
		:ref:`to-api-federation_resolvers`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'federation_resolvers', ('3.0', '4.0', '4.1', '5.0'))
	def create_federation_resolver(self, data=None):
		"""
		Create a federation resolver.
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'federation_resolvers', ('3.0', '4.0', '4.1', '5.0'))
	def delete_federation_resolver(self, query_params=None):
		"""
		Delete a federation resolver.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Federation User
	#
	@api_request('get', 'federations/{federation_id:d}/users', ('3.0', '4.0', '4.1', '5.0'))
	def get_federation_users(self, federation_id=None):
		"""
		Retrieves users assigned to a federation.
		:ref:`to-api-federations-id-users`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'federations/{federation_id:d}/users', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'federations/{federation_id:d}/users/{user_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	# ISO
	#
	@api_request('get', 'osversions', ('3.0', '4.0', '4.1', '5.0'))
	def get_osversions(self):
		"""
		Get all OS versions for ISO generation and the directory where the kickstarter files are
		found. The values are retrieved from osversions.json found in either ``/var/www/files`` or in
		the location defined by the kickstart.files.location parameter (if defined).
		:ref:`to-api-osversions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#TODO: this currently doesn't work, as /isos wasn't rewritten yet
	@api_request('post', 'isos', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'jobs', ('3.0', '4.0', '4.1', '5.0'))
	def get_jobs(self, query_params=None):
		"""
		Get all content-invalidation jobs (tenancy permitting).
		:ref:`to-api-jobs`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'jobs', ('3.0', '4.0', '4.1', '5.0'))
	def create_job(self, data=None):
		"""
		Creates a new content-invalidation job sorted by start time.
		:ref:`to-api-jobs`
		:param data: The content-invalidation job object that will be created.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'jobs', ('3.0', '4.0', '4.1', '5.0'))
	def update_job(self, data=None, query_params=None):
		"""
		Replaces a content-invalidation job with the one passed.
		:param data: The content-invalidation job with which the identified job will be replaced.
		:type data: Dict[str, Any]
		:param query_params: 'id' is a required parameter, identifying the job being updated.
		:ref:`to-api-jobs`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'jobs', ('3.0', '4.0', '4.1', '5.0'))
	def delete_job(self, query_params=None):
		"""
		Deletes a content-invalidation job.
		:ref:`to-api-jobs`
		:param query_params: 'id' is a required parameter, identifying the job being deleted.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Parameter
	#
	@api_request('get', 'parameters', ('3.0', '4.0', '4.1', '5.0'))
	def get_parameters(self, query_params=None):
		"""
		Get all Parameters.
		:ref:`to-api-parameters`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'profiles/{profile_id:d}/parameters', ('3.0', '4.0', '4.1', '5.0'))
	def get_parameters_by_profile_id(self, profile_id=None):
		"""
		Get all Parameters associated with a Profile by Id.
		:ref:`to-api-profiles-id-parameters`
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'profiles/name/{profile_name}/parameters', ('3.0', '4.0', '4.1', '5.0'))
	def get_parameters_by_profile_name(self, profile_name=None):
		"""
		Get all Parameters associated with a Profile by Name.
		:ref:`to-api-profiles-name-name-parameters`
		:param profile_name: The profile name
		:type profile_name: str
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'parameters', ('3.0', '4.0', '4.1', '5.0'))
	def create_parameter(self, data=None):
		"""
		Create Parameter
		:ref:`to-api-parameters`
		:param data: The parameter(s) data to use for parameter creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'parameters/{parameter_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_parameter(self, parameter_id=None, query_params=None):
		"""
		Update Parameter
		:ref:`to-api-parameters-id`
		:param parameter_id: The parameter id to update
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('delete', 'parameters/{parameter_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_parameter(self, parameter_id=None):
		"""
		Delete Parameter
		:ref:`to-api-parameters-id`
		:param parameter_id: The parameter id to delete
		:type parameter_id: int
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Server_capabilities
	#
	@api_request('get', 'server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def get_server_capabilities(self, query_params=None):
		"""
		Get all Server_capabilities.
		:ref:`to-api-server_capabilities`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def create_server_capabilities(self, data=None):
		"""
		Create Server_capabilities
		:ref:`to-api-server_capabilities`
		:param data: The server_capabilities(s) data to use for server_capabilities creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('delete', 'server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def delete_server_capabilities(self, query_params=None):
		"""
		Delete server_capabilities
		:ref:`to-api-server_capabilities`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


    #
	# multiple_server_capabilities
	#
	@api_request('post', 'multiple_servers_capabilities', ('4.1', '5.0'))
	def create_multiple_servers_capabilities(self, data=None):
		"""
		Create multiple_servers_capabilities
		:ref:`to-api-multiple_servers_capabilities`
		:param data: The multiple_servers_capabilities data to use for multiple_servers_capabilities creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('delete', 'multiple_servers_capabilities', ('4.1', '5.0'))
	def delete_multiple_servers_capabilities(self, data=None):
		"""
		Delete multiple_servers_capabilities
		:ref:`to-api-multiple_servers_capabilities`
		:param data: The multiple_servers_capabilities data to use for multiple_servers_capabilities creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Server_server_capabilities
	#
	@api_request('get', 'server_server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def get_server_server_capabilities(self, query_params=None):
		"""
		Gets all associations of Server Capabilities to cache servers.
		:ref:`to-api-server_server_capabilities`
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'server_server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def associate_server_capability_to_server(self, data=None):
		"""
		Associate server capability to server.
		:ref:`to-api-server_server_capabilities`
		:param query_params: 'serverId' is a required parameter, identifying which server to associate server capability.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'server_server_capabilities', ('3.0', '4.0', '4.1', '5.0'))
	def delete_server_capability_association_to_server(self, query_params=None):
		"""
		Deletes an association of server capability to a server.
		:ref:`to-api-server_server_capabilities`
		:param query_params: 'serverId' is a required parameter, identifying where server capability association is being deleted.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Physical Location
	#
	@api_request('get', 'phys_locations', ('3.0', '4.0', '4.1', '5.0'))
	def get_physical_locations(self, query_params=None):
		"""
		Get Physical Locations.
		:ref:`to-api-phys_locations`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'phys_locations', ('3.0', '3.1', '4.0', '4.1', '5.0'))
	def create_physical_locations(self, data=None):
		"""
		Create a Physical Location
		:ref:`to-api-phys_locations`
		:param data: The parameter data to use for Physical Location creation
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'phys_locations/{physical_location_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_physical_location(self, physical_location_id=None, query_params=None):
		"""
		Update Physical Location by id
		:ref:`to-api-phys_locations-id`
		:param physical_location_id: The id to update
		:type physical_location_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'phys_locations/{physical_location_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	# Plugins
	#
	@api_request('get', 'plugins', ('3.0', '4.0', '4.1', '5.0'))
	def get_plugins(self):
		"""
		Retrieves the list of plugins.
		:ref:`to-api-plugins`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Profiles
	#
	@api_request('get', 'profiles', ('3.0', '4.0', '4.1', '5.0'))
	def get_profiles(self, query_params=None):
		"""
		Get Profiles.
		:ref:`to-api-profiles`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'profiles', ('3.0', '4.0', '4.1', '5.0'))
	def create_profile(self, data=None):
		"""
		Create a profile
		:ref:`to-api-profiles`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'profiles/name/{new_profile_name:s}/copy/{copy_profile_name:s}', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('put', 'profiles/{profile_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'profiles/{profile_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'profileparameters', ('3.0', '4.0', '4.1', '5.0'))
	def get_profile_parameters(self, data=None):
		"""
		Retrieve all Parameter/Profile assignments.
		:ref:`to-api-profileparameters`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'profileparameters', ('3.0', '4.0', '4.1', '5.0'))
	def associate_paramater_to_profile(self, data=None):
		"""
		Associate parameter to profile.
		:ref:`to-api-profileparameters`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'profiles/{profile_id:d}/parameters', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('post', 'profileparameter', ('3.0', '4.0', '4.1', '5.0'))
	def assign_profile_to_parameter_ids(self, data=None):
		"""
		Create one or more profile / parameter assignments.
		:ref:`to-api-profileparameter`
		:param data: The data to assign
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'parameterprofile', ('3.0', '4.0', '4.1', '5.0'))
	def assign_parameter_to_profile_ids(self, data=None):
		"""
		Create one or more parameter / profile assignments.
		:ref:`to-api-profileparameter`
		:param data: The data to assign
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'profiles/name/{profile_name}/parameters', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'profileparameters/{profile_id:d}/{parameter_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'regions', ('3.0', '4.0', '4.1', '5.0'))
	def get_regions(self, query_params=None):
		"""
		Get Regions.
		:ref:`to-api-regions`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'regions', ('3.0', '4.0', '4.1', '5.0'))
	def create_region(self, query_params=None, data=None):
		"""
		Create a region
		:ref:`to-api-regions`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'regions', ('3.0', '4.0', '4.1', '5.0'))
	def delete_region(self, query_params=None):
		"""
		Delete a region by name or ID as a query parameter
		:ref:`to-api-regions-id`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'regions/{region_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_region(self, region_id=None):
		"""
		Update a region
		:ref:`to-api-regions-id`
		:parma region_id: The region to update
		:type region_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Roles
	#
	@api_request('get', 'roles', ('3.0', '4.0', '4.1', '5.0'))
	def get_roles(self):
		"""
		Get Roles.
		:ref:`to-api-roles`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'roles', ('3.0', '4.0', '4.1', '5.0'))
	def create_role(self, data=None):
		"""
		Create a new Role.
		:ref:`to-api-roles`
		:param data: A new Role object to be created.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'roles', ('3.0', '4.0', '4.1', '5.0'))
	def update_role(self, data=None, query_params=None):
		"""
		Get Roles.
		:ref:`to-api-roles`
		:param data: A new Role object which will replace the one identified.
		:type data: Dict[str, Any]
		:param query_params: 'id' is a required parameter, defining the Role to be replaced.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'roles', ('3.0', '4.0', '4.1', '5.0'))
	def delete_role(self, query_params=None):
		"""
		Delete a Role.
		:ref:`to-api-roles`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Server
	#
	@api_request('get', 'servers', ('3.0', '4.0', '4.1', '5.0'))
	def get_servers(self, query_params=None):
		"""
		Get Servers.
		:ref:`to-api-servers`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'servers/{server_id:d}/deliveryservices', ('3.0', '4.0', '4.1', '5.0'))
	def get_server_delivery_services(self, server_id=None):
		"""
		Retrieves all delivery services assigned to the server
		:ref:`to-api-servers-id-deliveryservices`
		:param server_id: The server id to retrieve
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'servers/details?hostName={name}', ('3.0',))
	def get_server_details(self, name=None):
		"""
		Get servers/details
		:ref:`to-api-v3-servers-details`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]

		.. deprecated:: 3.0
			The endpoint this represents has been removed from APIv4 and clients should use get_servers instead.
		"""

	@api_request('post', 'servercheck', ('3.0', '4.0', '4.1', '5.0'))
	def create_servercheck(self, data=None):
		"""
		Post a server check result to the serverchecks table.
		:ref:`to-api-servercheck`
		:param data: The parameter data to use for server creation
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'servers', ('3.0', '4.0', '4.1', '5.0'))
	def create_server(self, data=None):
		"""
		Create a new Server.
		:ref:`to-api-servers`
		:param data: The parameter data to use for server creation
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'servers/{server_id:d}', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('put', 'servers/{server_id:d}/status', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'servers/{server_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_server_by_id(self, server_id=None):
		"""
		Delete a Server by Id.
		:ref:`to-api-servers-id`
		:param server_id: The server Id
		:type server_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'servers/{server_id:d}/queue_update', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('get', 'servers/{server_name}/update_status', ('3.0', '4.0', '4.1', '5.0'))
	def get_server_update_status(self, server_name=None):
		"""
		Gets the current update status of a server named ``server_name``.
		:ref:`to-api-servers-hostname-update_status`
		:param server_name: The (short) hostname of the server for which the update status will be fetched
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Service categories
	#
	@api_request('get', 'service_categories', ('3.0', '4.0', '4.1', '5.0'))
	def get_service_categories(self, query_params=None):
		"""
		Get service categories.
		:ref:`to-api-service_categories`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	@api_request('post', 'service_categories', ('3.0', '4.0', '4.1', '5.0'))
	def create_service_category(self, data=None):
		"""
		Create a service category.
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'service_categories/{service_category_name:s}', ('3.0', '4.0', '4.1', '5.0'))
	def update_service_category(self, service_category_name=None, data=None):
		"""
		Update a service category.
		:param service_category_name: The current name of the Service Category.
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'service_categories/{service_category_name:s}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_service_category(self, service_category_name=None):
		"""
		Delete a service category.
		:param service_category_name: The current name of the Service Category to be deleted.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Static DNS Entries
	#
	@api_request('get', 'staticdnsentries', ('3.0', '4.0', '4.1', '5.0'))
	def get_staticdnsentries(self, query_params=None):
		"""
		Get static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'staticdnsentries', ('3.0', '4.0', '4.1', '5.0'))
	def create_staticdnsentries(self, data=None):
		"""
		Create static DNS entries associated with the delivery service
		:ref:`to-api-staticdnsentries`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'staticdnsentries', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'staticdnsentries', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'statuses', ('3.0', '4.0', '4.1', '5.0'))
	def get_statuses(self, query_params=None):
		"""
		Retrieves a list of the server status codes available.
		:ref:`to-api-statuses`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'statuses', ('3.0', '3.1', '4.0', '4.1', '5.0'))
	def create_statuses(self, query_params=None):
		"""
		Create server status code.
		:ref:`to-api-statuses`
		:param data: A new status code created.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'statuses/{status_id:d}', ('3.0', '3.1', '4.0', '4.1', '5.0'))
	def delete_status_by_id(self, status_id=None):
		"""
		Delete a status
		:ref:`to-api-statuses-id`
		:param status_id: The status to delete
		:type status_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# steering_id_targets
	#
	@api_request('get', 'steering/{delivery_service_id:d}/targets', ('3.0', '4.0', '4.1', '5.0'))
	def get_steering_targets(self, delivery_service_id=None, query_params=None):
		"""
		Get all targets for a steering Delivery Service.
		:ref:`to-api-steering-id-targets`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'steering/{delivery_service_id:d}/targets', ('3.0', '4.0', '4.1', '5.0'))
	def create_steering_targets(self, delivery_service_id=None, data=None):
		"""
		Create a steering target.
		:ref:`to-api-steering-id-targets`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param data: The steering(s) data to use for steering creation.
		:type data: Union[Dict[str, Any], List[Dict[str, Any]]]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'steering/{delivery_service_id:d}/targets/{target_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_steering_targets(self, delivery_service_id=None, target_id=None):
		"""
		Removes a specific target mapping from a specific Delivery Service.
		:ref:`to-api-steering-id-targets`
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param target_id: The target Id
		:type target_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# System
	#
	@api_request('get', 'system/info', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'tenants', ('3.0', '4.0', '4.1', '5.0'))
	def get_tenants(self, query_params=None):
		"""
		Get all tenants.
		:ref:`to-api-tenants`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'tenants/{tenant_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_tenant(self, tenant_id=None):
		"""
		Update a tenant
		:ref:`to-api-tenants-id`
		:param tenant_id: The tenant to update
		:type tenant_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'tenants', ('3.0', '4.0', '4.1', '5.0'))
	def create_tenant(self, data=None):
		"""
		Create a tenant
		:ref:`to-api-tenants`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'tenants/{tenant_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_tenant(self, tenant_id=None):
		"""
		Delete a tenant
		:ref:`to-api-tenants-id`
		:param tenant_id: The tenant to delete
		:type tenant_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# TO Extensions
	#
	@api_request('get', 'servercheck/extensions', ('3.0', '4.0', '4.1', '5.0'))
	def get_servercheck_extensions(self):
		"""
		Retrieves the list of extensions.
		:ref:`to-api-servercheck_extensions`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'servercheck/extensions', ('3.0', '4.0', '4.1', '5.0'))
	def create_to_extension(self, data=None):
		"""
		Creates a Traffic Ops extension.
		:ref:`to-api-servercheck_extensions`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'servercheck/extensions/{extension_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_to_extension(self, extension_id=None):
		"""
		Deletes a Traffic Ops extension.
		:ref:`to-api-servercheck_extensions-id`
		:param extension_id: The extension id to delete
		:type extension_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Topologies
	#
	@api_request('post', 'topologies', ('3.0', '4.0', '4.1', '5.0'))
	def create_topology(self, data: Dict[str, Any]=None) -> Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]:
		"""
		Create a topology
		:ref:`to-api-topologies`
		:param data: The Topology data to use for Topology creation.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'topologies', ('3.0', '4.0', '4.1', '5.0'))
	def get_topologies(self, query_params: Dict[str, Any]=None) -> Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]:
		"""
		Get Topologies.
		:ref:`to-api-topologies`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'topologies?name={name:s}', ('3.0', '4.0', '4.1', '5.0'))
	def update_topology(self, name: str=None, data: Dict[str, Any]=None) -> Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]:
		"""
		Update a Topology
		:ref:`to-api-topologies`
		:param name: The name of the Topology
		:type name: str
		:param data: The new values for the Topology
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('delete', 'topologies?name={name:s}', ('3.0', '4.0', '4.1', '5.0'))
	def delete_topology(self, name: str=None) -> Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]:
		"""
		Delete a Topology
		:ref:`to-api-topologies`
		:param name: The name of the Topology to delete
		:type name: str
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'topologies/{name:s}/queue_update', ('3.0', '4.0', '4.1', '5.0'))
	def topologies_queue_update(self, name=None, data=None):
		"""
		Queue Updates by Topology name.
		:ref:`to-api-topologies-name-queue_update`
		:param name: The Topology name
		:param data: The update action.
		:type data: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""


	#
	# Types
	#
	@api_request('get', 'types', ('3.0', '4.0', '4.1', '5.0'))
	def get_types(self, query_params=None):
		"""
		Get Data Types.
		:ref:`to-api-types`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Users
	#
	@api_request('get', 'users', ('3.0', '4.0', '4.1', '5.0'))
	def get_users(self):
		"""
		Retrieves all users.
		:ref:`to-api-users`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'users/{user_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def get_user_by_id(self, user_id=None):
		"""
		Retrieves user by ID.
		:ref:`to-api-users-id`
		:param user_id: The user to retrieve
		:type user_id: int
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'users', ('3.0', '4.0', '4.1', '5.0'))
	def create_user(self, data=None):
		"""
		Create a user.
		:ref:`to-api-users`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'users/{user_id:d}', ('3.0', '4.0', '4.1', '5.0'))
	def update_user_by_id(self, user_id=None, data=None):
		"""
		Update a user.
		:ref:`to-api-users`
		:param data: The user update data payload.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'users/register', ('3.0', '4.0', '4.1', '5.0'))
	def create_user_with_registration(self, data=None):
		"""
		Register a user and send registration email
		:ref:`to-api-users-register`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('get', 'user/current', ('3.0', '4.0', '4.1', '5.0'))
	def get_authenticated_user(self):
		"""
		Retrieves the profile for the authenticated user.
		:ref:`to-api-user-current`
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'user/current', ('3.0', '4.0', '4.1', '5.0'))
	def replace_authenticated_user(self, data=None):
		"""
		Updates the currently authenticated user.
		:ref:`to-api-user-current`
		:param data: The new user information which will replace the current user's user information.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Snapshot CRConfig
	#
	@api_request('get', 'cdns/{cdn_name}/snapshot', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('get', 'cdns/{cdn_name}/snapshot/new', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('put', 'snapshot', ('3.0', '4.0', '4.1', '5.0'))
	def snapshot_crconfig(self, query_params=None):
		"""
		Snapshot CRConfig by CDN Name or ID.
		:ref:`to-api-snapshot`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Dict[str, Any], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	#
	# Coordinate
	#
	@api_request('get', 'coordinates', ('3.0', '4.0', '4.1', '5.0'))
	def get_coordinates(self, query_params=None):
		"""
		Get all coordinates associated with the cdn
		:ref:`to-api-coordinates`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'coordinates', ('3.0', '4.0', '4.1', '5.0'))
	def create_coordinates(self, data=None):
		"""
		Create coordinates
		:ref:`to-api-coordinates`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'coordinates', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'coordinates', ('3.0', '4.0', '4.1', '5.0'))
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
	@api_request('get', 'origins', ('3.0', '4.0', '4.1', '5.0'))
	def get_origins(self, query_params=None):
		"""
		Get origins associated with the delivery service
		:ref:`to-api-origins`
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('post', 'origins', ('3.0', '4.0', '4.1', '5.0'))
	def create_origins(self, data=None):
		"""
		Creates origins associated with a delivery service
		:ref:`to-api-origins`
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[str, Any]
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

	@api_request('put', 'origins', ('3.0', '4.0', '4.1', '5.0'))
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

	@api_request('delete', 'origins', ('3.0', '4.0', '4.1', '5.0'))
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




if __name__ == '__main__':
	# Sample usages
	import operator

	DEBUG = False

	logging.basicConfig(stream=sys.stderr, level=logging.INFO if not DEBUG else logging.DEBUG)

	# TOSession Class Examples
	#     TOSession is a class that allows you to create a session to a Traffic Ops instance
	#     and interact with the Traffic Ops API.

	# Traffic Ops System - for self-signed cert -> turn off cert verification
	TOS = TOSession(host_ip='to.somedomain.net', verify_cert=True)
	TOS.login('someuser', 'someuser123')

	# Objects get returned munch-ified by default which means you can access dictionary keys as
	# attributes names but you can still access the entries with keys as well e.g.
	# ``cdn.name`` is equivalent to ``cdn['name']``
	CDNS = TOS.get_cdns()[0]
	print(CDNS)
	for cdn in CDNS:
		print('CDN [{0}] has id [{1}]'.format(cdn.name, cdn.id))

	ALL_TYPES = TOS.get_types()[0]
	print('All Types are (sorted by useInTable, name):')
	print(ALL_TYPES)
	for atype in sorted(ALL_TYPES, key=operator.itemgetter('useInTable', 'name')):
		print('Type [{0}] for table [{1}]'.format(atype.name, atype.useInTable))

	print('Getting all cache groups (bulk)...')
	CACHE_GROUPS = TOS.get_cachegroups()[0]
	for cache_group in CACHE_GROUPS:
		print('Bulk cache group [{0}] has id [{1}]'.format(cache_group.name, cache_group.id))

		# Example with URL replacement parameters
		# e.g. TOSession.get_cachegroups_by_id() is the '/cachegroups/{id}' API endpoint
		#      See TOSession object for details.
		print('    Getting cachegroup by id [{0}]'.format(cache_group.id),
			  ' to demonstrate getting by id...')
		cg_id_list = TOS.get_cachegroup_by_id(cache_group_id=cache_group.id)[0]
		print('    Cache group [{0}] by id [{1}]'.format(cg_id_list[0].name, cg_id_list[0].id))

	# Example with URL query parameters
	SERVER_TYPES = TOS.get_types(query_params={'useInTable': 'server'})[0]
	print('Server Types are:')
	print(SERVER_TYPES)
	for stype in SERVER_TYPES:
		print('Type [{0}] for table [{1}]'.format(stype.name, stype.useInTable))
	TOS.close()
	print('Done!')
