#! /usr/bin/env python

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

import requests.exceptions as rex

# Local Modules
import common.restapi as restapi
import common.utils as utils


logger = logging.getLogger(__name__)

__all__ = [u'default_headers', u'TOSession']

# Miscellaneous Constants and/or Variables
default_headers = {u'Content-Type': u'application/json; charset=UTF-8'}


# TOSession Class
class TOSession(restapi.RestApiSession):
	"""
	Traffic Ops Session Class
	Once you login to the Traffic Ops API via the 'login' method, you can call one or more of the methods to retrieve,
	post, put, delete, etc. data to the API.  If you are not logged in, an exception will be thrown if you try
	to call any of the endpoint methods. e.g. get_servers, get_cachegroups, etc.

	This API client is simplistic and lightly structured on purpose but adding support for new end-points
	routinely takes seconds.  Another nice bit of convenience that result data is, by default, wrapped in
	munch.Munch objects, which provide attribute access to the returned dictionaries/hashes.

		e.g. "a_dict['a_key']" with munch becomes "a_dict.a_key" or "a_dict['a_key']"
			 "a_dict['a_key']['b_key']" with munch becomes "a_dict.a_key.b_key" or "a_dict['a_key']['b_key']"

	Also, the lack of rigid structure (loose coupling) means many changes to the Traffic Ops API,
	as it evolves, will probably go un-noticed (usually additions), which means fewer
	future problems to potentially fix in user applications.

	An area of improvement for later is defining classes to represent request data instead
	of loading up dictionaries for request data.

	As of now you can see the following URL for API details:
	   https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html #api for details

	Adding end-point methods: (See "Implemented Direct API URL Endpoint Methods" for actual examples)
		E.g. End-point with no URL parameters and no query parameters:
			given end-point URL: GET api/1.2/cdns
				@restapi.api_request(u'get', u'cdns', (u'1.1', u'1.2',))
				def get_cdns(self):
					pass

		E.g. End-point with URL parameters and no query parameters:
			given end-point URL: GET api/1.2/cdns/{cdn_id:d}
				 @restapi.api_request(u'get', u'cdns/{cdn_id:d}', (u'1.1', u'1.2',))
				 def get_cdn_by_id(self, cdn_id=None):
					 pass

		E.g. End-point with no URL parameters but with query parameters:
			given end-point URL: GET api/1.2/deliveryservices
				 @restapi.api_request(u'get', u'deliveryservices', (u'1.1', u'1.2',))
				 def get_deliveryservices(self, query_params=None):
					 pass

		E.g. End-point with URL parameters and query parameters:
			given end-point URL: GET api/1.2/deliveryservices/xmlId/{xml_id}/sslkeys
				 @restapi.api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys', (u'1.1', u'1.2',))
				 def get_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
					 pass

		E.g. End-point with request data:
			given end-point URL: POST api/1.2/cdns
				 @restapi.api_request(u'post', u'cdns', (u'1.1', u'1.2',))
				 def create_cdn(self, data=None):
					 pass

		E.g. End-point with URL parameters and request data:
			given end-point URL: PUT api/1.2/cdns/{cdn_id:d}
				 @restapi.api_request(u'put', u'cdns', (u'1.1', u'1.2',))
				 def update_cdn_by_id(self, cdn_id=None, data=None):
					 pass

	Calling end-point methods:

		E.g. Using no URL parameters and no query parameters:
			given end-point URL: GET api/1.2/cdns
			get_cdns() -> calls end-point: GET api/1.2/cdns

		E.g. Using no URL parameters but with query parameters:
			given end-point URL: GET api/1.2/types
			get_types(query_params={'useInTable': 'servers'}) -> calls end-point: GET api/1.2/types?useInTable=servers

		E.g. Using URL parameters and query parameters:
		   given end-point URL: GET api/1.2/foo/{id}
		   get_foo_data(id=45, query_params={'sort': 'asc'}) -> calls end-point: GET api/1.2/foo/45?sort=asc

		E.g. Using with required request data:
			given end-point URL: POST api/1.2/cdns/{id:d}/queue_update
			cdns_queue_update(...) -> calls end-point -> POST api/1.2/cdns/{id:d}/queue_update
			cdns_queue_update(id=1, data={'action': 'queue'}) -> calls end-point: POST api/1.2/cdns/1/queue_update
			   with json data '{"action": "queue"}'.

		   So,

		   dict_request = {'action': 'queue'}

		   or

		   Example with a namedtuple:
			   import collections
			   QueueUpdateRequest = collections.namedtuple('QueueUpdateRequest', ['action'])
			   request = QueueUpdateRequest(action='update')

		   Then:
			   cdns_queue_update(id=1, data=vars(request))     # Python 2.x
			   cdns_queue_update(id=1, data=request.asdict())  # Python 3.x
			   cdns_queue_update(id=1, data=dict_request)      # Python 2.x/3.x

		   NOTE: var(request)/request.asdict() transforms the namedtuple into a dictionary which is required
				 by the 'data' argument.

	NOTE: Only a small subset of the API endpoints are implemented.  More can be implemented as needed.
		  See the Traffic Ops API documentation for more detail:
					 https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html #api for details #api
	"""

	def __init__(self, host_ip, host_port=443, api_version=u'1.3', ssl=True, headers=default_headers,
				 verify_cert=True):
		"""
		The class initializer.
		:param host_ip: The dns name or ip address of the Traffic Ops host to use to talk to the API
		:type host_ip: Text
		:param host_port: The port to use when contacting the Traffic Ops API
		:type host_port: int
		:param api_version: The version of the API to use when calling end-points on the Traffic Ops API
		:type api_version: Text
		:param ssl: Should ssl be used? http vs. https
		:type ssl: bool
		:param headers:  The http headers to use when contacting the Traffic Ops API
		:type headers: Dict[Text, Text]
		:param verify_cert: Should the ssl certificates be verified when contacting the Traffic Ops API.
							You may want to set this to False for systems with self-signed certificates.
		:type verify_cert: bool
		"""
		super(TOSession, self).__init__(host_ip=host_ip, api_version=api_version,
										api_base_path=u'api/{api_version}/',
										host_port=host_port, ssl=ssl, headers=headers, verify_cert=verify_cert)

		self._logged_in = False

		msg = u'TOSession instance {0:#0x} initialized: Details: {1}'
		utils.log_with_debug_info(logging.DEBUG, msg.format(id(self), self.__dict__))

	def login(self, username, password):
		"""
		Login to the Traffic Ops API.
		:param username: Traffic Ops User Name
		:type username: Text
		:param password: Traffic Ops User Password
		:type password: Text
		:return: None
		:rtype: None
		:raises: trafficops.restapi.LoginError
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
			msg = (u'{0}.  This system may have a self-signed certificate.  Try creating this TOSession '
				   u'object passing verify_cert=False. e.g. TOSession(..., verify_cert=False). ')
			msg = msg.format(e)
			logging.error(msg)
			logging.warning("disabling certificate verification is not recommended.")
			raise restapi.LoginError(msg) from e
		except restapi.OperationError as e:
			logging.debug("%s", e, exc_info=True, stack_info=True)
			msg = u'Logging in to Traffic Ops has failed. Reason: {0}'.format(e)
			self.close()
			logging.error(msg)
			raise restapi.OperationError(msg) from e

		logging.info("Authenticated.")

	@property
	def to_url(self):
		"""
		The URL without the api portion. (read-only)
		:return: The url should be in the format of
				 '<protocol>://<hostname>[:<port>]'; [] = optional
				 e.g https://to.somedomain.net or https://to.somedomain.net:443
		:rtype: Text
		"""

		return self.server_url

	@property
	def base_url(self):
		"""
		Returns the base url. (read-only)
		:return: The base url should be in the format of
				 '<protocol>://<hostname>[:<port>]/api/<api version>/'; [] = optional
				 e.g https://to.somedomain.net/api/1.2/
		:rtype: Text
		"""

		return self._api_base_url

	@property
	def logged_in(self):
		"""
		Read-only property of boolean to determine if user is logged in to Traffic Ops. (read-only)
		:return: boolean if logged in or not.
		:rtype: bool
		"""

		return self.is_open and self._logged_in

	# Programmatic Endpoint Methods - These can be created when you need to employ "creative
	# methods" to form a correlated composite data set from one or more Traffic Ops API call(s) or
	# employ composite operations against the API.
	# Also, if the API requires you to retrieve the data via paging, these types of methods can be
	# useful to perform that type of work too.
	# These methods need to support similar method signatures as employed by the restapi.api_request decorator
	# method_name argument.
	def get_all_deliveryservice_servers(self, *args, **kwargs):
		"""
		Get all servers attached to all delivery services via the Traffic Ops API.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
		result_set = []
		response = None
		limit = 10000
		page = 1

		munchify = True  # Default to True
		if u'munchify' in kwargs:
			munchify = kwargs[u'munchify']

		while True:
			data, response = self.get_deliveryserviceserver(query_params={u'limit': limit, u'page': page},
															munchify=munchify, *args, **kwargs)

			if not data:
				break

			result_set.extend(munch.munchify(data) if munchify else data)
			page += 1

		return result_set, response  # Note: Return last response object received

	# Implemented Direct API URL Endpoint Methods
	# See https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html #api for detail
	@restapi.api_request(u'get', u'asns', (u'1.1', u'1.2', u'1.3',))
	def get_asns(self, query_params=None):
		"""
		Get ASNs.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'cachegroups', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroups(self, query_params=None):
		"""
		Get Cache Groups.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	# Example of URL parameter substitution via call parameter. You will need to pass the parameter
	# value as a keyword parameter with the proper type to match the str.format specification,
	# e.g. 'cachegroups/{cache_group_id:d}'.  In this case, ':d' specifies a decimal integer.  A specification
	# of 'cachegroups/{cache_group_id}' will try to convert any value passed to a string, which basically does
	# no type checking, unless of course the value cannot be cast to a string.
	# E.g. get_cachegroups_by_id(cache_group_id=23) -> call end-point .../api/1.2/cachegroups/23
	@restapi.api_request(u'get', u'cachegroups/{cache_group_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_cachegroup_by_id(self, cache_group_id=None):
		"""
		Get a Cache Group by Id.
		:param cache_group_id: The cache group Id
		:type cache_group_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservices(self, query_params=None):
		"""
		Get Delivery Services.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_by_id(self, delivery_service_id=None):
		"""
		Get a Delivery Service by Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
	@restapi.api_request(u'get', u'servers/hostname/{name}/details', (u'1.1', u'1.2', u'1.3',))
	def get_server_details(self, name=None):
		"""
		#GET /api/1.2/servers/hostname/:name/details
		Get server details from trafficOps
		https://traffic-control-cdn.readthedocs.io/en/latest/api/v12/server.html
		:param hostname: Server hostname
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
	@restapi.api_request(u'post', u'deliveryservices', (u'1.1', u'1.2', u'1.3',))
	def create_deliveryservice(self, data=None):
		"""
		Create a Delivery Service.
		:param data: The request data structure for the API request
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'put', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_deliveryservice_by_id(self, delivery_service_id=None, data=None):
		"""
		Update a Delivery Service by Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param data: The request data structure for the API request
		:type data: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'delete', u'deliveryservices/{delivery_service_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_by_id(self, delivery_service_id=None):
		"""
		Delete a Delivery Service by Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices/{delivery_service_id:d}/servers', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservices_servers(self, delivery_service_id=None):
		"""
		Get all servers associated with a Delivery Service Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryserviceserver', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryserviceserver(self, query_params=None):
		"""
		Get Servers for all defined Delivery Services.
		:param query_params: The required url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryserviceserver', (u'1.1', u'1.2', u'1.3',))
	def assign_deliveryservice_servers_by_ids(self, data=None):
		"""
		Assign servers by id to a Delivery Service. (New Method)
		:param data: The required data to create server associations to a delivery service
		:type data: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryservices/{xml_id}/servers', (u'1.1', u'1.2', u'1.3',))
	def assign_deliveryservice_servers_by_names(self, xml_id=None, data=None):
		"""
		Assign severs by name to a Delivery Service by xmlId. (Old Method)
		:param xml_id: The XML Id of the delivery service
		:type xml_id: Text
		:param data: The required data to assign servers to a delivery service
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices_regexes', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservices_regexes(self):
		"""
		Get RegExes for all Delivery Services.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices/{delivery_service_id:d}/regexes', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_regexes_by_id(self, delivery_service_id=None):
		"""
		Get RegExes for a Delivery Service by Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryservices/regexes', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_regexes(self, data=None):
		"""
		Delete RegExes.
		:param data: The required data to delete delivery service regexes
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'delete', u'deliveryservices/{delivery_service_id:d}/regexes/{delivery_service_regex_id:d}',
						 (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_regex_by_regex_id(self, delivery_service_id=None, delivery_service_regex_id=None):
		"""
		Delete a RegEx by Id for a Delivery Service by Id.
		:param delivery_service_id: The delivery service Id
		:type delivery_service_id: int
		:param delivery_service_regex_id: The delivery service regex Id
		:type delivery_service_regex_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys', (u'1.1', u'1.2', u'1.3',))
	def get_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
		"""
		Get SSL keys for a Delivery Service by xmlId.
		:param xml_id: The Delivery Service XML id
		:type xml_id: Text
		:param query_params: The url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'deliveryservices/xmlId/{xml_id}/sslkeys/delete', (u'1.1', u'1.2', u'1.3',))
	def delete_deliveryservice_ssl_keys_by_xml_id(self, xml_id=None, query_params=None):
		"""
		Delete SSL keys for a Delivery Service by xmlId.
		:param xml_id: The Delivery Service xmlId
		:type xml_id: Text
		:param query_params: The url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryservices/sslkeys/generate', (u'1.1', u'1.2', u'1.3',))
	def generate_deliveryservice_ssl_keys(self, data=None):
		"""
		Generate an SSL certificate. (self-signed)
		:param data: The parameter data to use for Delivery Service SSL key generation.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryservices/sslkeys/add', (u'1.1', u'1.2', u'1.3',))
	def add_ssl_keys_to_deliveryservice(self, data=None):
		"""
		Add SSL keys to a Delivery Service.
		:param data: The parameter data to use for adding SSL keys to a Delivery Service.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'deliveryservices/xmlId/{xml_id}/urlkeys/generate', (u'1.1', u'1.2', u'1.3',))
	def generate_deliveryservice_url_signature_keys(self, xml_id=None):
		"""
		Generate URL Signature Keys for a Delivery Service by xmlId.
		:param xml_id: The Delivery Service xmlId
		:type xml_id: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'cdns', (u'1.1', u'1.2', u'1.3',))
	def get_cdns(self):
		"""
		Get all CDNs.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'cdns', (u'1.1', u'1.2', u'1.3',))
	def create_cdn(self, data=None):
		"""
		Create a new CDN.
		:param data: The parameter data to use for cdn creation.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
	@restapi.api_request(u'get', u'cdns/{cdn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_cdn_by_id(self, cdn_id=None):
		"""
		Get a CDN by Id.
		:param cdn_id: The CDN id
		:type cdn_id: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'cdns/name/{cdn_name}', (u'1.1', u'1.2', u'1.3',))
	def get_cdn_by_name(self, cdn_name=None):
		"""
		Get a CDN by name.
		:param cdn_name: The CDN name
		:type cdn_name: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'put', u'cdns/{cdn_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_cdn_by_id(self, cdn_id=None, data=None):
		"""
		Update a CDN by Id.
		:param cdn_id: The CDN id
		:type cdn_id: int
		:param data: The parameter data to use for cdn update.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'servers', (u'1.1', u'1.2',u'1.3',))
	def get_servers(self, query_params=None):
		"""
		Get Servers.
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'servers', (u'1.1', u'1.2', u'1.3',))
	def create_server(self, data=None):
		"""
		Create a new Server.
		:param data: The parameter data to use for server creation
		:type data: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'put', u'servers/{server_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_server_by_id(self, server_id=None, data=None):
		"""
		Update a Server by Id.
		:param server_id: The server Id
		:type server_id: int
		:param data: The parameter data to edit
		:type data: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
	@restapi.api_request(u'put', u'servers/{server_id:d}/status', (u'1.1', u'1.2', u'1.3',))
	def update_server_status_by_id(self, server_id=None, data=None):
		"""
		Update server_status by Id.
		:param server_id: The server Id
		:type server_id: int
		:status: https://traffic-control-cdn.readthedocs.io/en/latest/api/v12/server.html
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'delete', u'servers/{server_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_server_by_id(self, server_id=None):
		"""
		Delete a Server by Id.
		:param server_id: The server Id
		:type server_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'parameters', (u'1.1', u'1.2', u'1.3',))
	def get_parameters(self):
		"""
		Get all Profile Parameters.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'profiles', (u'1.1', u'1.2', u'1.3',))
	def get_profiles(self, query_params=None):
		"""
		Get Profiles.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_profile_by_id(self, profile_id=None):
		"""
		Get Profile by Id.
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'put', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def update_profile_by_id(self, profile_id=None, data=None):
		"""
		Update Profile by Id.
		:param profile_id: The profile Id
		:type profile_id: int
		:param data: The parameter data to edit
		:type data: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'delete', u'profiles/{profile_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_profile_by_id(self, profile_id=None):
		"""
		Delete Profile by Id.
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'parameters/validate', (u'1.1', u'1.2', u'1.3',))
	def validate_parameter_exists(self, data=None):
		"""
		Validate that a Parameter exists.
		:param data: The parameter data to use for parameter validation.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'parameters/{parameter_id:d}', (u'1.1', u'1.2', u'1.3',))
	def get_parameter_by_id(self, parameter_id=None):
		"""
		Get a Parameter by Id.
		:param parameter_id: The parameter Id
		:type parameter_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'profiles/{id:d}/parameters', (u'1.1', u'1.2', u'1.3',))
	def get_parameters_by_profile_id(self, profile_id=None):
		"""
		Get all Parameters associated with a Profile by Id.
		:param profile_id: The profile Id
		:type profile_id: int
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'profiles/name/{profile_name}/parameters', (u'1.1', u'1.2', u'1.3',))
	def get_parameters_by_profile_name(self, profile_name=None):
		"""
		Get all Parameters associated with a Profile by Name.
		:param profile_name: The profile name
		:type profile_name: Text
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'parameters', (u'1.1', u'1.2', u'1.3',))
	def create_parameters(self, data=None):
		"""
		Create Parameters
		:param data: The parameter(s) data to use for parameter creation.
		:type data: Union[Dict[Text, Any], List[Dict[Text, Any]]]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'parameters/{parameter_id:d}/profiles', (u'1.1', u'1.2', u'1.3',))
	def get_associated_profiles_by_parameter_id(self, parameter_id=None):
		"""
		Get all Profiles associated to a Parameter by Id.
		:param parameter_id: The parameter id
		:type parameter_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'profiles/id/{profile_id:d}/parameters', (u'1.1', u'1.2', u'1.3',))
	def associate_parameters_by_profile_id(self, profile_id=None, data=None):
		"""
		Associate Parameters to a Profile by Id.
		:param profile_id: The profile id
		:type profile_id: int
		:param data: The parameter data to associate
		:type data: Union[Dict[Text, Any], List[Dict[Text, Any]]]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'profiles/name/{profile_name}/parameters', (u'1.1', u'1.2', u'1.3',))
	def associate_parameters_by_profile_name(self, profile_name=None, data=None):
		"""
		Associate Parameters to a Profile by Name.
		:param profile_name: The profile name
		:type profile_name: Text
		:param data: The parameter data to associate
		:type data: Union[Dict[Text, Any], List[Dict[Text, Any]]]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'delete', u'profileparameters/{profile_id:d}/{parameter_id:d}', (u'1.1', u'1.2', u'1.3',))
	def delete_profile_parameter_association_by_id(self, profile_id=None, parameter_id=None):
		"""
		Delete Parameter association by Id for a Profile by Id.
		:param profile_id: The profile id
		:type profile_id: int
		:param parameter_id: The parameter id
		:type parameter_id: int
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'phys_locations', (u'1.1', u'1.2', u'1.3',))
	def get_physical_locations(self, query_params=None):
		"""
		Get Physical Locations.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'users', (u'1.1', u'1.2', u'1.3',))
	def get_users(self):
		"""
		Get Users.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'roles', (u'1.1', u'1.2', u'1.3',))
	def get_roles(self):
		"""
		Get Roles.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'statuses', (u'1.1', u'1.2', u'1.3',))
	def get_statuses(self):
		"""
		Get Statuses.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'types', (u'1.1', u'1.2', u'1.3',))
	def get_types(self, query_params=None):
		"""
		Get Data Types.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'staticdnsentries', (u'1.1', u'1.2', u'1.3',))
	def get_static_dns_entries(self):
		"""
		Get Static DNS Entries.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'cdns/{cdn_id:d}/queue_update', (u'1.1', u'1.2', u'1.3',))
	def cdns_queue_update(self, cdn_id=None, data=None):
		"""
		Queue Updates by CDN Id.
		:param cdn_id: The CDN Id
		:type cdn_id: int
		:param data: The update action. QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'post', u'servers/{server_id:d}/queue_update', (u'1.1', u'1.2', u'1.3',))
	def servers_queue_update(self, server_id=None, data=None):
		"""
		Queue Updates by Server Id.
		:param server_id: The server Id
		:type server_id: int
		:param data: The update action.  QueueUpdateRequest() can be used for this argument also.
		:type data: Dict[Text, Any]
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'put', u'snapshot/{cdn_name}', (u'1.1', u'1.2', u'1.3',))
	def snapshot_crconfig(self, cdn_name=None):
		"""
		Snapshot CRConfig by CDN Name.
		:param cdn_name: The CDN name
		:type cdn_name: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'cdns/{cdn_name}/snapshot', (u'1.2', u'1.3',))
	def get_current_snapshot_crconfig(self, cdn_name=None):
		"""
		Retrieve the currently implemented CR Snapshot
		:param cdn_name: The CDN name
		:type cdn_name: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""
	@restapi.api_request(u'get', u'cdns/{cdn_name}/snapshot/new', (u'1.2', u'1.3',))
	def get_pending_snapshot_crconfig(self, cdn_name=None):
		"""
		Retrieve the pending CR Snapshot
		:param cdn_name: The CDN name
		:type cdn_name: Text
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'logs', (u'1.2', u'1.3',))
	def get_change_logs(self):
		"""
		Retrieve all change logs from traffic ops
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'logs/{days:d}/days', (u'1.2', u'1.3',))
	def get_change_logs_for_days(self, days=None):
		"""
		Retrieve all change logs from Traffic Ops
		:param days: The number of days to retrieve change logs
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'logs/newcount', (u'1.2', u'1.3',))
	def get_change_logs_newcount(self):
		"""
		Get amount of new logs from traffic ops
		:rtype: Tuple[Dict[Text, Any], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	###                                                                              ###
	###                                                                              ###
	### add version 3 endpoints from here                                            ###
	### ref: http://traffic-control-cdn.readthedocs.io/en/latest/api/v13/index.html  ###
	###                                                                              ###
	###                                                                              ###
	###                                                                              ###

	@restapi.api_request(u'get', u'coordinates', (u'1.3',))
	def get_coordinates(self, query_params=None):
		"""
		Get all coordinates associated with the cdn
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request("get", "servers/{servername:s}/configfiles/ats", ("1.3",))
	def getServerConfigFiles(self, servername=None):
		"""
		Fetches the configuration files for a server with the given short hostname
		:param servername: The short hostname of the server
		:returns: The response content and actual response object
		"""

	####################################################################################

	####                          Data Model Overrides                              ####

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

	@restapi.api_request(u'get', u'origins', (u'1.3',))
	def get_origins(self, query_params=None):
		"""
		Get origins associated with the delivery service
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""

	@restapi.api_request(u'get', u'staticdnsentries', (u'1.3',))
	def get_staticdnsentries(self, query_params=None):
		"""
		Get static DNS entries associated with the delivery service
		:param query_params: The optional url query parameters for the call
		:type query_params: Dict[Text, Any]
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]]], requests.Response]
		:raises: Union[trafficops.restapi.LoginError, trafficops.restapi.OperationError]
		"""


if __name__ == u'__main__':
	# Sample usages
	import sys
	import operator

	DEBUG = False

	logging.basicConfig(stream=sys.stderr, level=logging.INFO if not DEBUG else logging.DEBUG)

	# TOSession Class Examples
	#     TOSession is a class that allows you to create a session to a Traffic Ops instance
	#     and interact with the Traffic Ops API.

	# Traffic Ops System - for self-signed cert -> turn off cert verification
	tos = TOSession(host_ip=u'to.somedomain.net', verify_cert=True)
	tos.login(u'someuser', u'someuser123')

	# Objects get returned munch-ified by default which means you can access dictionary keys as
	# attributes names but you can still access the entries with keys as well.  E.g. cdn.name == cdn['name']
	cdns, response = tos.get_cdns()
	print(cdns)
	for cdn in cdns:
		print(u'CDN [{0}] has id [{1}]'.format(cdn.name, cdn.id))

	all_types, response = tos.get_types()
	print(u'All Types are (sorted by useInTable, name):')
	print(all_types)
	for atype in sorted(all_types, key=operator.itemgetter(u'useInTable', u'name')):
		print(u'Type [{0}] for table [{1}]'.format(atype.name, atype.useInTable))

	print(u'Getting all cache groups (bulk)...')
	cache_groups, response = tos.get_cachegroups()
	for cache_group in cache_groups:
		print(u'Bulk cache group [{0}] has id [{1}]'.format(cache_group.name, cache_group.id))

		# Example with URL replacement parameters
		# e.g. TOSession.get_cachegroups_by_id() == end-point 'api/1.2/cachegroups/{id}'
		#      See TOSession object for details.
		print(u'    Getting cachegroup by id [{0}] to demonstrate getting by id...'.format(cache_group.id))
		cg_id_list, response = tos.get_cachegroup_by_id(cache_group_id=cache_group.id)  # data returned is always a list
		print(u'    Cache group [{0}] by id [{1}]'.format(cg_id_list[0].name, cg_id_list[0].id))

	# Example with URL query parameters
	server_types, response = tos.get_types(query_params={u'useInTable': u'server'})
	print(u'Server Types are:')
	print(server_types)
	for stype in server_types:
		print(u'Type [{0}] for table [{1}]'.format(stype.name, stype.useInTable))
	tos.close()
	print(u'Done!')
