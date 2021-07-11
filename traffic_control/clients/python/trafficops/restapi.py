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
Module to help retrieve/create/update/delete data from/to any RESTful API (Base Class).

Requires Python Version >= 2.7 or >= 3.6
"""

# Core Modules
import json
import logging
import functools
from builtins import str

# Third-party Modules
import munch
import requests
import requests.adapters as ra

# Python 2 to Python 3 Compatibility
import requests.compat as compat
# Local Modules
from .utils import log_with_debug_info


try:
	from future.utils import iteritems
except ImportError:
	iteritems = lambda x: x.items()

__all__ = ['LoginError', 'OperationError', 'InvalidJSONError', 'api_request', 'RestApiSession',
           'DEFAULT_HEADERS']

LOGGER = logging.getLogger(__name__)

# Exception Classes
class LoginError(OSError):
	"""
	This represents an error that occurred during server login.
	"""
	def __init__(self, *args):
		OSError.__init__(self, *args)


class OperationError(IOError):
	"""
	This class represents a generic error, indicating something went wrong with the request or on
	the server.
	"""
	#: Contains the response object that generated the error
	resp = None
	def __init__(self, *args, resp=None):
		IOError.__init__(self, *args)
		self.resp = resp


class InvalidJSONError(ValueError):
	"""
	An error that occurs when an invalid JSON payload is passed to an endpoint.
	"""
	#: Contains the response object that generated the error
	resp = None
	def __init__(self, *args, resp=None):
		ValueError.__init__(self, *args)
		self.resp = resp

# Miscellaneous Constants and/or Variables
DEFAULT_HEADERS = {u'Content-Type': u'application/json; charset=UTF-8'}


# Helper Functions/Decorators
def api_request(method_name, api_path, supported_versions):
	"""
	This wrapper returns a decorator that routes the calls to the appropriate utility function that
	generates the RESTful API endpoint, performs the appropriate call to the endpoint and returns
	the data to the user.

	:param method_name: A method name defined on the Class, this decorator is decorating, that will be called to perform the operation. E.g. 'GET', 'POST', 'PUT', 'DELETE', etc. The method_name chosen must have the signature of ``<method>(self, api_path, **kwargs)`` e.g. ``def get(self, api_path, **kwargs): ...``
	:type method_name: str
	:param api_path: The path to the API end-point that you want to call which does not include the base url e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
	:type api_path: str
	:param supported_versions: A tuple of API versions that this route supports
	:type supported_versions: Tuple[str]
	:return: rtype int: A new function that replaces the original function with a boilerplate execution process.
	:rtype: Callable[str, Dict[str, Any]]
	"""

	def outer(func):
		@functools.wraps(func)
		def method_wrapper(self, *args, **kwargs):
			# Positional arguments, e.g. *args, are not being used. Keyword arguments are the
			# preferred way to pass the parameters needed by the helper functions
			if (self.api_version is None) or (self.api_version in supported_versions):
				msg = (u'Calling method [{0}] with keyword arguments [{1}] '
				       u'via API endpoint method [{2}]')
				log_with_debug_info(logging.DEBUG, msg.format(method_name, kwargs, func.__name__))

				return getattr(self, method_name)(api_path, **kwargs)

			# Client API version is not supported by the method being called
			msg = (u"Method [{0}] is not supported by this client's API version [{1}]; "
			       u'Supported versions: {2}')  # type: Text
			msg = msg.format(func.__name__, self.api_version, supported_versions)
			log_with_debug_info(logging.DEBUG, msg)
			raise OperationError(msg)

		return method_wrapper
	return outer


class RestApiSession(object):
	"""
	This class represents a login session with a generic REST API server. It provides base
	functionality inherited by :class:`TOSession`.
	"""
	def __init__(self, host_ip, api_version=None, api_base_path=u'api/', host_port=443, ssl=True,
	             headers=None, verify_cert=True, create_session=False, max_retries=5):
		"""
		The class initializer.

		:param host_ip: The dns name or ip address of the RESTful API host to use to talk to the API
		:type host_ip: str
		:param host_port: The port to use when contacting the RESTful API
		:type host_port: int
		:param api_version: The version of the API to make calls against. If supplied, endpoint version validation will be performed. If supplied as None, no version validation will be performed. :const:`None` is allowed so that non-versioned REST APIs can be implemented.
		:type api_version: Union[str, None]
		:param api_base_path: The part of the url that is the base path, from the web server root (which may include an API version), for all API endpoints without the server url portion e.g. 'api/', 'api/latest/'

			.. note:: To specify the base path with the passed ``api_version`` you can specify
				``api_base_path`` as ``api/{api_version}/`` and the API version will be substituted.
				If ``api_version`` is :const:`None` and '{api_version}' is specified in the
				``api_base_path`` string then an exception will be thrown.
				e.g. api_version=u'latest' -> 'api/{api_version}/' -> 'api/latest/'
				api_version=None   -> 'api/{api_version}/' -> Throws Exception

		:type api_base_path: str
		:param ssl: Should SSL be used? (http vs. https)
		:type ssl: bool
		:param headers: The HTTP headers to use when contacting the RESTful API
		:type headers: Dict[str, str]
		:param verify_cert: Should the SSL certificates be verified when contacting the RESTful API. You may want to set this to :const:`False` for systems with self-signed certificates.
		:type verify_cert: bool
		:param create_session: Should a session be created automatically?
		:type create_session: bool
		"""
		if headers is None:
			headers = DEFAULT_HEADERS

		self._session = None
		self._host_ip = host_ip
		self._host_port = host_port
		self._api_version = api_version
		self._api_base_path = api_base_path
		self._ssl = ssl
		self._headers = headers
		self._verify_cert = verify_cert
		self._create_session = create_session
		self._max_retries = max_retries

		# Setup API End-point Version validation, if enabled
		self.__api_version_format_name = u'api_version'
		self.__api_version_format_value = u'{{{0}}}'.format(self.__api_version_format_name)

		if self._api_version:
			# if api_base_path is supplied as 'api/{api_version}/' or some string
			# containing '{api_version}' then try to substitute the api_version supplied
			# by the user.

			version_params = {
				self.__api_version_format_name: self._api_version
			}
			self._api_base_path = self._api_base_path.format(**version_params)

		if not self._api_version and self.__api_version_format_value in self._api_base_path:
			msg = (u'{0} was specified in the API Base Path [{1}] '
			       u'but the replacement did not occur because the API Version '
			       u'was not supplied.')
			msg = msg.format(self.__api_version_format_value, self._api_base_path)
			log_with_debug_info(logging.ERROR, msg)
			raise OperationError(msg)

		# Setup some common URLs
		self._server_url = u'{0}://{1}{2}/'.format(u'https' if ssl else u'http',
		                                           host_ip,
		                                           u':{0}'.format(host_port) if host_port else u'')
		self._api_base_url = compat.urljoin(self._server_url, self._api_base_path)
		self._api_base_url = self._api_base_url.rstrip(u'/') + u'/'


		if not self._verify_cert:
			# Not verifying certs so let's disable the warning
			#pylint: disable=E1101
			requests.packages.urllib3.disable_warnings(\
			                            requests.packages.urllib3.exceptions.InsecureRequestWarning)
			#pylint: enable=E1101
			log_with_debug_info(logging.WARNING, u'Certificate verification warnings are disabled.')

		msg = u'RestApiSession instance {0:#0x} initialized: Details: {1}'
		log_with_debug_info(logging.DEBUG, msg.format(id(self), self.__dict__))

		if self._create_session:
			self.create()

	@property
	def is_open(self):
		"""
		Is the session open to the RESTful API? (Read-only Property)

		:return: :const:`True` if yes, otherwise, :const:`False`
		:rtype: bool
		"""
		return self._session is not None

	@property
	def session(self):
		"""
		The RESTful API session (Read-only Property)

		:return: The requests session
		:rtype: :class:`requests.Session`
		"""
		return self._session

	def create(self):
		"""
		Create the requests.Session to communicate with the RESTful API.

		:return: :const:`None`
		:rtype: NoneType
		"""
		if self._session:
			self.close()

		if not self._session:
			self._session = requests.Session()
			self._session.mount('http://', ra.HTTPAdapter(max_retries=self._max_retries))
			self._session.mount('https://', ra.HTTPAdapter(max_retries=self._max_retries))

			msg = u'Created internal requests Session instance {0:#0x}'
			log_with_debug_info(logging.DEBUG, msg.format(id(self._session)))

	def close(self):
		"""
		Close and cleanup the requests Session object.

		:return: :const:`None`
		:rtype: NoneType
		"""

		if self._session:
			sid = id(self._session)
			self._session.close()
			del self._session
			self._session = None

			msg = u'Internal requests Session instance 0x{0:x} closed and cleaned up'
			log_with_debug_info(logging.DEBUG, msg.format(sid))

	@property
	def server_url(self):
		"""
		The URL without the api portion. (read-only)

		:return: The URL should match '[\\w\\+\\-\\.]+://[\\w\\+\\-\\.]+(:\\d+)?' e.g. 'https://to.somedomain.net' or 'https://to.somedomain.net:443'
		:rtype: str
		"""

		return self._server_url

	@property
	def api_version(self):
		"""
		Returns the api version. (read-only)

		:return: The api version from which this instance will request endpoints.
		:rtype: str
		"""

		return self._api_version

	@property
	def api_base_url(self):
		"""
		Returns the base URL. (read-only)

		:return: The base URL should match '[\\w\\+\\-\\.]+://[\\w\\+\\-\\.]+(:\\d+)?' e.g. 'https://to.somedomain.net/api/0.1/'
		:rtype: str
		"""

		return self._api_base_url

	def _build_endpoint(self, api_path, params=None, query_params=None):
		"""
		Helper function to form API URL. The base URL is
		'<protocol>://<hostname>[:<port>]/<api base url>' e.g. 'https://to.somedomain.net/api/0.1/'

		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param params: If :meth:`str.format` field_name replacement field specifications exists in the ``api_path`` use this dictionary to perform replacements of the specifications with the value(s) in the dictionary that match the parameter name(s) e.g. ``{param_id}`` or ``{param_id:d}`` in ``api_string`` is replaced by value in ``params['param_id']``.
		:type params: Union[Dict[str, Any], None]
		:param query_params: URL query params to provide to the end-point e.g. ``{ 'sort': 'asc', 'maxresults': 200 }`` which translates to something like ``?sort=asc&maxresults=200`` which is appended to the request URL
		:type query_params: Union[Dict[str, Any], None]
		:return: The base url plus the passed and possibly substituted ``api_path`` to form a complete URL to the API resource to request
		:rtype: str
		:raises: ValueError
		"""

		new_api_path = api_path

		# Replace all parameters in the new_api_path path, if required
		try:
			# Make the parameters values safe for adding to URLs
			url_params = {k: compat.quote(str(v)) if isinstance(v, str)\
			                                      else v for k, v in iteritems(params)}

			log_with_debug_info(logging.DEBUG, u'URL parameters are: [{0}]'.format(url_params))

			qparams = u''
			if query_params:
				# Process the URL query parameters
				qparams = u'?{0}'.format(compat.urlencode(query_params))
				log_with_debug_info(logging.DEBUG,
				                    u'URL query parameters are: [{0}]'.format(qparams))

			new_api_path = api_path.format(**url_params) + qparams
		except KeyError as e:
			msg = (u'Expecting a value for keyword argument [{0}] for format field '
			       u'specification [{1!r}]')
			msg = msg.format(e, api_path)
			log_with_debug_info(logging.ERROR, msg)
			raise ValueError(msg)
		except ValueError as e:
			msg = (u'One or more values do not match the format field specification '
			       u'[{0!r}]; Supplied values: {1!r} ')
			msg = msg.format(api_path, params)
			log_with_debug_info(logging.ERROR, msg)
			raise ValueError(msg)

		retval = compat.urljoin(self.api_base_url, new_api_path)

		log_with_debug_info(logging.DEBUG, u'Built end-point to return: {0}'.format(retval))

		return retval

	def _do_operation(self,
	                  operation, api_path, query_params=None, munchify=True, debug_response=False,
	                  expected_status_codes=range(200, 300), *unused_args, **kwargs):
		"""
		Helper method to perform HTTP operation requests - This is a boilerplate process for HTTP
		operations.

		:param operation: Name of method to call on the :attr:`self._session` object to perform the HTTP request
		:type operation: str
		:param api_path: The path to the API end-point that you want to call which does not include the URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param: query_params: URL query parameters to provide to the endpoint e.g. ``{ 'sort': 'asc', 'maxresults': 200 }`` which translates to something like ``?sort=asc&maxresults=200`` which is appended to the request URL
		:type query_params: Union[Dict[str, Any], None]
		:param: munchify: If :const:`True` encapsulate data to be returned in a :class:`munch.Munch` object which allows keys in a Python dictionary to additionally have attribute access e.g. ``a_dict['a_key']`` with :mod:`munch` becomes ``a_dict['a_key']`` or ``a_dict.a_key``
		:type munchify: bool
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure e.g. a ``dict``. This method will convert it to JSON before sending it to the API endpoint.
		:type kwargs: Dict[str, Any]
		:param debug_response: If :const:`True`, the actual response data text will be added to the log if a JSON decoding exception is encountered.
		:type debug_response: bool
		:type expected_status_codes: Tuple[int]
		:param: expected_status_codes: expected success HTTP status codes. If the user needs to override the defaults this parameter can be passed e.g. ``(200, 204,)``
		:type munchify: bool
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[Text, Any], List[Dict[Text, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: miscellaneous.exceptions.OperationError
		"""

		if not self._session:
			msg = u'No session has been created for the API.  Have you called create() yet?'
			log_with_debug_info(logging.ERROR, msg)
			raise OperationError(msg)

		response = None
		retdata = None

		endpoint = self._build_endpoint(api_path, params=kwargs, query_params=query_params)

		params = {u'headers': self._headers, u'verify': self._verify_cert}

		if u'data' in kwargs:
			params[u'data'] = json.dumps(kwargs[u'data'])

		log_with_debug_info(logging.DEBUG, u'Call parameters: {0}'.format(params))

		# Call the API endpoint
		response = getattr(self._session, operation)(endpoint, **params)

		log_with_debug_info(logging.DEBUG, u'Response status: {0} {1}'.format(response.status_code,
		                    response.reason))

		if response.status_code not in expected_status_codes:
			try:
				retdata = response.json()
			except Exception as e:
				# Invalid JSON payload.
				msg = (u'HTTP Status Code: [{0}]; API response data for end-point [{1}] does not '
				       u'appear to be valid JSON. Cause: {2}.')
				msg = msg.format(response.status_code, endpoint, e)
				if debug_response:
					log_with_debug_info(logging.ERROR, msg + u' Data: [' + str(response.text) + u']')
				raise InvalidJSONError(msg, resp=response)
			msg = u'{0} request to RESTful API at [{1}] expected status(s) {2}; failed: {3} {4};'\
			      u' Response: {5}'
			msg = msg.format(operation.upper(), endpoint, expected_status_codes,
			                 response.status_code, response.reason, retdata)
			log_with_debug_info(logging.ERROR, msg)
			raise OperationError(msg, resp=response)

		try:
			if response.status_code in ('204',):
				# "204 No Content"
				retdata = {}
			else:
				# Decode the expected JSON
				retdata = response.json()
		except Exception as e:
			# Invalid JSON payload.
			msg = (u'HTTP Status Code: [{0}]; API response data for end-point [{1}] does not '
			       u'appear to be valid JSON. Cause: {2}.')
			msg = msg.format(response.status_code, endpoint, e)
			if debug_response:
				log_with_debug_info(logging.ERROR, msg + u' Data: [' + str(response.text) + u']')
			raise InvalidJSONError(msg, resp=response)
		retdata = munch.munchify(retdata) if munchify else retdata
		return (retdata[u'response'] if u'response' in retdata else retdata), response

	def get(self, api_path, *args, **kwargs):
		"""
		Perform http get requests

		:param api_path: The path to the API end-point that you want to call which does not include the base url e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'get', api_path, *args, **kwargs)

	def post(self, api_path, *args, **kwargs):
		"""
		Perform http post requests

		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'post', api_path, *args, **kwargs)

	def put(self, api_path, *args, **kwargs):
		"""
		Perform http put requests

		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'put', api_path, *args, **kwargs)

	def delete(self, api_path, *args, **kwargs):
		"""
		Perform HTTP DELETE requests

		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'delete', api_path, *args, **kwargs)

	def head(self, api_path, *args, **kwargs):
		"""
		Perform HTTP HEAD requests
		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'head', api_path, *args, **kwargs)

	def options(self, api_path, *args, **kwargs):
		"""
		Perform HTTP OPTIONS requests
		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'options', api_path, *args, **kwargs)

	def patch(self, api_path, *args, **kwargs):
		"""
		Perform HTTP PATCH requests
		:param api_path: The path to the API end-point that you want to call which does not include the base URL e.g. ``user/login``, ``servers``, etc. This string can contain substitution parameters as denoted by a valid field_name replacement field specification as per :meth:`str.format` e.g. ``cachegroups/{id}`` or ``cachegroups/{id:d}``
		:type api_path: str
		:param kwargs: Passed Keyword Parameters. If you need to send JSON data to the endpoint pass the keyword parameter ``data`` with the Python data structure. This method will convert it to JSON before sending it to the API endpoint. Use ``query_params`` to pass a dictionary of query parameters
		:type kwargs: Dict[str, Any]
		:return: Python data structure distilled from JSON from the API request.
		:rtype: Tuple[Union[Dict[str, Any], List[Dict[str, Any]], munch.Munch, List[munch.Munch]], requests.Response]
		:raises: Union[LoginError, OperationError]
		"""

		return self._do_operation(u'patch', api_path, *args, **kwargs)
