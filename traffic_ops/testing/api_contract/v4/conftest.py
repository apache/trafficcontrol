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
This module is used to create a Traffic Ops session and to store prerequisite
data for endpoints.
"""

import json
import logging
import shutil
import sys
import os
from random import randint
from typing import Any, NamedTuple, Union, Optional, TypeAlias
from urllib.parse import urlparse
import munch
import psycopg2

import pytest
import requests
from trafficops.tosession import TOSession
from trafficops.restapi import OperationError

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

JSONData: TypeAlias = Union[dict[str, object], list[object], bool, int, float, Optional[str]]
JSONData.__doc__ = """An alias for the kinds of data that JSON can encode."""

class APIVersion(NamedTuple):
	"""Represents an API version."""
	major: int
	minor: int

	@staticmethod
	def from_string(ver_str: str) -> "APIVersion":
		"""
		Instantiates a new version from a string.

		>>> APIVersion.from_string("4.0")
		APIVersion(major=4, minor=0)
		>>> try:
		... 	APIVersion("not a version string")
		... except ValueError:
		... 	print("whoops")
		...
		whoops
		>>>
		>>> try:
		... 	APIVersion("4.Q")
		... except ValueError:
		... 	print("whoops")
		...
		whoops
		"""
		parts = ver_str.split(".", 1)
		if len(parts) != 2:
			raise ValueError("invalid version; must be of the form '{{major}}.{{minor}}'")
		return APIVersion(int(parts[0]), int(parts[1]))

	def __str__(self) -> str:
		"""
		Coalesces the version to a string.

		>>> print(APIVersion(4, 1))
		4.1
		"""
		return f"{self.major}.{self.minor}"

APIVersion.major.__doc__ = """The API's major version number."""
APIVersion.major.__doc__ = """The API's minor version number."""

class ArgsType(NamedTuple):
	"""Represents the configuration needed to create Traffic Ops session."""
	user: str
	password: str
	url: str
	port: int
	api_version: APIVersion

	def __str__(self) -> str:
		"""
		Formats the configuration as a string. Omits password and extraneous
		properties.

		>>> print(ArgsType("user", "password", "url", 420, APIVersion(4, 0)))
		User: 'user', URL: 'url'
		"""
		return f"User: '{self.user}', URL: '{self.url}'"


ArgsType.user.__doc__ = """The username used for authentication."""
ArgsType.password.__doc__ = """The password used for authentication."""
ArgsType.url.__doc__ = """The URL of the environment."""
ArgsType.port.__doc__ = """The port number on which to connect to Traffic Ops."""
ArgsType.api_version.__doc__ = """The version number of the API to use."""


class DbArgsType(NamedTuple):
	"""Represents the configuration needed to create Traffic Ops Database connection."""
	db_name: str
	user: str
	password: str
	hostname: str
	port: int
	sslmode: str

	def __str__(self) -> str:
		"""
		Formats the configuration as a string. Omits password and extraneous
		properties.

		>>> print(ArgsType("db_name", "user", "password", "hostname", "port", "sslmode"))
		Dbname: 'db_name', User: 'user'
		"""
		return f"User: '{self.db_name}', : '{self.user}'"

DbArgsType.db_name.__doc__ = """The DB name used for authentication."""
DbArgsType.user.__doc__ = """The DB username used for authentication."""
DbArgsType.password.__doc__ = """The DB password used for authentication."""
DbArgsType.hostname.__doc__ = """The DB hostname of the environment."""
DbArgsType.port.__doc__ = """The DB port number on which to connect to Traffic Ops."""
DbArgsType.sslmode.__doc__ = """Sslmode to use."""


@pytest.fixture(autouse=True, scope='function')
def delete_pytest_cache():
	"""
	Deletes cached data before every test case execution
	"""
	shutil.rmtree(".pytest_cache", ignore_errors=True)


def pytest_addoption(parser: pytest.Parser) -> None:
	"""
	Parses the Traffic Ops arguments from command line.
	:param parser: Parser to parse command line arguments.
	"""
	parser.addoption(
		"--to-user", action="store", help="User name for Traffic Ops Session."
	)
	parser.addoption(
		"--to-password", action="store", help="Password for Traffic Ops Session."
	)
	parser.addoption(
		"--to-url", action="store", help="Traffic Ops URL."
	)
	parser.addoption(
		"--to-db-name", action="store", help="Name for Traffic Ops Database."
	)
	parser.addoption(
		"--to-db-user", action="store", help="User name for Traffic Ops Database."
	)
	parser.addoption(
		"--to-db-password", action="store", help="Password for Traffic Ops Database."
	)
	parser.addoption(
		"--to-db-hostname", action="store", help="Hostname for Traffic Ops Database."
	)
	parser.addoption(
		"--to-db-port", action="store", help="Port for Traffic Ops Database."
	)
	parser.addoption(
		"--to-db-sslmode", action="store", help="Sslmode for Traffic Ops Database."
	)
	parser.addoption(
		"--config",
		help="Path to configuration file.",
		default=os.path.join(os.path.dirname(__file__), "data", "to_data.json")
	)
	parser.addoption(
		"--request-template",
		help="Path to request prerequisites file.",
		default=os.path.join(os.path.dirname(__file__), "data", "request_template.json")
	)
	parser.addoption(
		"--response-template",
		help="Path to response prerequisites file.",
		default=os.path.join(os.path.dirname(__file__), "data", "response_template.json")
	)

def coalesce_config(
	arg: Optional[object],
	file_key: str,
	file_contents: Optional[dict[str, Optional[object]]],
	env_key: str
) -> Optional[str]:
	"""
	Coalesces configuration retrieved from different sources into a single
	string.

	This will raise a ValueError if the type of the configuration value in the
	parsed configuration file is not a string.

	In order of descending precedence this checks the command-line argument
	value, the configuration file value, and then the environment variable
	value.

	:param arg: The command-line argument value.
	:param file_key: The key under which to look in the parsed JSON configuration file data.
	:param file_contents: The parsed JSON configuration file (if one was used).
	:param env_key: The environment variable name to look for a value if one wasn't provided elsewhere.
	:returns: The coalesced configuration value, or 'None' if no value could be determined.
	"""
	if isinstance(arg, str):
		return arg

	if file_contents:
		file_value = file_contents.get(file_key)
		if isinstance(file_value, str):
			return file_value
		if file_value is not None:
			raise ValueError(f"incorrect value; want: 'str', got: '{type(file_value)}'")

	return os.environ.get(env_key)

def parse_to_url(raw: str) -> tuple[APIVersion, int]:
	"""
	Parses the API version and port number from a raw URL string.

	>>> parse_to_url("https://trafficops.example.test:420/api/5.270")
	(APIVersion(major=5, minor=270), 420)
	>>> parse_to_url("trafficops.example.test")
	(APIVersion(major=4, minor=0), 443)
	"""
	parsed = urlparse(raw)
	if not parsed.netloc:
		raise ValueError("missing network location (hostname & optional port)")

	if parsed.scheme and parsed.scheme.lower() != "https":
		raise ValueError("invalid scheme; must use HTTPS")

	port = 443
	if ":" in parsed.netloc:
		port_str = parsed.netloc.split(":")[-1]
		try:
			port = int(port_str)
		except ValueError as e:
			raise ValueError(f"invalid port number: {port_str}") from e

	api_version = APIVersion(4, 0)
	if parsed.path and parsed.path != "/":
		ver_str = parsed.path.lstrip("/api/").split("/", 1)[0]
		if not ver_str:
			raise ValueError(f"invalid API path: {parsed.path} (should be e.g. '/api/4.0')")
		api_version = APIVersion.from_string(ver_str)
	else:
		logging.warning("using default API version: %s", api_version)

	return (api_version, port)


@pytest.fixture(name="to_db_args", scope="session")
def to_db_data(pytestconfig: pytest.Config) -> DbArgsType:
	"""
	PyTest fixture to store Traffic ops database arguments passed from command line.
	:param pytestconfig: Session-scoped fixture that returns the session's pytest.Config object.
	:returns: Configuration for connecting to Traffic Ops database.
	"""
	session_data: JSONData = None
	cfg_path = pytestconfig.getoption("--config")
	if isinstance(cfg_path, str):
		try:
			with open(cfg_path, encoding="utf-8", mode="r") as session_file:
				session_data = json.load(session_file)
		except (FileNotFoundError, PermissionError) as read_err:
			raise ValueError(f"could not read configuration file at '{cfg_path}'") from read_err

	if session_data is not None and not isinstance(session_data, dict):
		raise ValueError(
			f"invalid configuration file; expected top-level object, got: {type(session_data)}"
		)

	to_db_name = coalesce_config(pytestconfig.getoption("--to-db-name"), "db_name",
			      session_data, "TO_DB_NAME")
	if not to_db_name:
		raise ValueError(
			"Traffic Ops Database name is not configured - use '--to-db-name', the config file, or an "
			"environment variable to do so"
		)

	to_db_user = coalesce_config(pytestconfig.getoption("--to-db-user"), "db_user",
			      session_data, "TO_DB_USER")
	if not to_db_user:
		raise ValueError(
			"Traffic Ops Database Username is not configured - use '--to-db-user', the config file, or an "
			"environment variable to do so"
		)

	to_db_password = coalesce_config(pytestconfig.getoption("--to-db-password"), "db_password",
				  session_data, "TO_DB_PASSWORD")

	if not to_db_password:
		raise ValueError(
			"Traffic Ops Database password is not configured - use '--to-db-password', the config file,"
			"or an environment variable to do so"
		)

	to_db_hostname = coalesce_config(pytestconfig.getoption("--to-db-hostname"), "db_hostname",
				  session_data, "TO_DB_HOSTNAME")
	if not to_db_hostname:
		raise ValueError(
			"Traffic Ops Database Hostname is not configured - use '--to-db-hostname', the config file,"
			"or an environment variable to do so"
		)

	to_db_port = coalesce_config(pytestconfig.getoption("--to-db-port"), "db_port",
			      session_data, "TO_DB_PORT")
	if not to_db_port:
		raise ValueError(
			"Traffic Ops Database Port is not configured - use '--to-db-port', the config file, or an "
			"environment variable to do so"
		)
	port = int(to_db_port)

	to_db_sslmode = coalesce_config(pytestconfig.getoption("--to-db-sslmode"), "db_sslmode",
				 session_data, "TO_DB_SSLMODE")
	if not to_db_sslmode:
		raise ValueError(
			"Traffic Ops Database Sslmode is not configured - use '--to-db-sslmode', the config file, or an "
			"environment variable to do so"
		)

	return DbArgsType(
		to_db_name,
		to_db_user,
		to_db_password,
		to_db_hostname,
		port,
		to_db_sslmode
	)


@pytest.fixture(name="to_args")
def to_data(pytestconfig: pytest.Config) -> ArgsType:
	"""
	PyTest fixture to store Traffic ops arguments passed from command line.
	:param pytestconfig: Session-scoped fixture that returns the session's pytest.Config object.
	:returns: Configuration for connecting to Traffic Ops.
	"""
	session_data: JSONData = None
	cfg_path = pytestconfig.getoption("--config")
	if isinstance(cfg_path, str):
		try:
			with open(cfg_path, encoding="utf-8", mode="r") as session_file:
				session_data = json.load(session_file)
		except (FileNotFoundError, PermissionError) as read_err:
			raise ValueError(f"could not read configuration file at '{cfg_path}'") from read_err

	if session_data is not None and not isinstance(session_data, dict):
		raise ValueError(
			f"invalid configuration file; expected top-level object, got: {type(session_data)}"
		)

	to_user = coalesce_config(pytestconfig.getoption("--to-user"), "user", session_data, "TO_USER")
	if not to_user:
		raise ValueError(
			"Traffic Ops user is not configured - use '--to-user', the config file, or an "
			"environment variable to do so"
		)

	to_password = coalesce_config(
		pytestconfig.getoption("--to-password"),
		"password",
		session_data,
		"TO_PASSWORD"
	)

	if not to_password:
		raise ValueError(
			"Traffic Ops password is not configured - use '--to-password', the config file, or an "
			"environment variable to do so"
		)

	to_url = coalesce_config(pytestconfig.getoption("--to-url"), "url", session_data, "TO_USER")
	if not to_url:
		raise ValueError(
			"Traffic Ops URL is not configured - use '--to-url', the config file, or an "
			"environment variable to do so"
		)

	try:
		api_version, port = parse_to_url(to_url)
	except ValueError as e:
		raise ValueError("invalid Traffic Ops URL") from e

	return ArgsType(
		to_user,
		to_password,
		to_url,
		port,
		api_version
	)

@pytest.fixture(name="to_session")
def to_login(to_args: ArgsType) -> TOSession:
	"""
	PyTest Fixture to create a Traffic Ops session from Traffic Ops Arguments
	passed as command line arguments in to_args fixture in conftest.

	:param to_args: Fixture to get Traffic ops session arguments.
	:returns: An authenticated Traffic Ops session.
	"""
	# Create a Traffic Ops V4 session and login
	to_url = urlparse(to_args.url)
	to_host = to_url.hostname
	try:
		to_session = TOSession(
			host_ip=to_host,
			host_port=to_args.port,
			api_version=str(to_args.api_version),
			ssl=True,
			verify_cert=False
		)
		logger.info("Established Traffic Ops Session.")
	except OperationError as error:
		logger.debug("%s", error, exc_info=True, stack_info=True)
		logger.error("Failure in Traffic Ops session creation. Reason: %s", error)
		sys.exit(-1)

	# Login To TO_API
	to_session.login(to_args.user, to_args.password)
	logger.info("Successfully logged into Traffic Ops.")
	return to_session


@pytest.fixture(scope="session", name="db_connection")
def to_db_connection(to_db_args: DbArgsType) -> psycopg2.connect:
	"""
	Creates new traffic ops db connection.
	:returns: New Traffic ops database connection
	"""
	to_db_connection = None
	try:
		to_db_connection = psycopg2.connect(
            user=to_db_args.user,
            password=to_db_args.password,
            host=to_db_args.hostname,
            port=to_db_args.port,
            database=to_db_args.db_name,
            sslmode=to_db_args.sslmode
        )
		logger.info("Successfully connected to the Traffic Ops database.")
		yield to_db_connection
	except psycopg2.OperationalError as e:
		logger.error("Error connecting to the Traffic Ops database : %s", e)
	finally:
		if to_db_connection:
			to_db_connection.close()
			logger.info("Closed Traffic ops DB connection.")


@pytest.fixture(name="request_template_data", scope="session")
def request_prerequiste_data(pytestconfig: pytest.Config, request: pytest.FixtureRequest
			  ) -> list[Union[dict[str, object], list[object], Primitive]]:
	"""
	PyTest Fixture to store POST request template data for api endpoint.
	:param pytestconfig: Session-scoped fixture that returns the session's pytest.Config object.
	:param request: Fixture to access information about the requesting test function and its fixtures

	:returns: Prerequisite request data for api endpoint.
	"""
	request_template_path = pytestconfig.getoption("--request-template")
	if not isinstance(request_template_path, str):
		# unlike the configuration file, this must be present
		raise ValueError("prereqisites path not configured")

	# Response keys for api endpoint
	data: Union[dict[
		str,
		Union[list[Union[dict[str, object], list[object], Primitive]], dict[object, object], Primitive]
	], Primitive] = None
	with open(request_template_path, encoding="utf-8", mode="r") as prereq_file:
		data = json.load(prereq_file)
	if not isinstance(data, dict):
		raise TypeError(f"request template data must be an object, not '{type(data)}'")
	try:
		request_template = data[request.param]
		if not isinstance(request_template, list):
			raise TypeError(f"Request template data must be a list, not '{type(request_template)}'")
	except AttributeError:
		request_template = data
	return request_template


@pytest.fixture()
def response_template_data(pytestconfig: pytest.Config
			   ) -> dict[str, Union[Primitive, list[
	Union[Primitive, dict[str, object], list[object]]], dict[object, object]]]:
	"""
	PyTest Fixture to store response template data for api endpoint.
	:param pytestconfig: Session-scoped fixture that returns the session's pytest.Config object.
	:returns: Prerequisite response data for api endpoint.
	"""
	prereq_path = pytestconfig.getoption("--response-template")
	if not isinstance(prereq_path, str):
		# unlike the configuration file, this must be present
		raise ValueError("prereqisites path not configured")

	# Response keys for api endpoint
	response_template: Union[dict[
		str,
		Union[list[Union[dict[str, object], list[object], Primitive]], dict[object, object], Primitive]
	], Primitive] = None
	with open(prereq_path, encoding="utf-8", mode="r") as prereq_file:
		response_template = json.load(prereq_file)
	if not isinstance(response_template, dict):
		raise TypeError(f"Response template data must be an object, not '{type(response_template)}'")

	return response_template


def api_response_data(api_response: tuple[Union[Primitive, dict[str, object], list[
	Union[Primitive, dict[str, object], list[
		object]]]], requests.Response], create_check: bool) -> dict[str, object]:
	"""
	Checks API get/post response.
	:param api_response: Raw api response.
	:returns: Verified response data
	"""
	api_data = None
	if isinstance(api_response, tuple):
		api_response = api_response[0]
		if create_check and not isinstance(api_response, munch.Munch):
			raise ValueError("Malformed API response; 'response' property not an munch")
		if not create_check and not isinstance(api_response, list):
			raise ValueError("Malformed API response; 'response' property not an list")
	else:
		raise ValueError("Invalid API response format")

	if api_response:
		try:
			api_data = api_response
			if not create_check:
				api_data = api_response[0]
			if not isinstance(api_data, dict):
				raise ValueError("Malformed API response; 'response' property not a dict")
		except IndexError as e:
			raise ValueError(f"No response data from API request '{e.args[0]}'") from e

	return api_data


def get_existing_object(to_session: TOSession, object_type: str, query_params: Optional[
	dict[str, Any]])-> Union[dict[str, Any], None]:
	"""
	Check if the given endpoint with the given query params already exists.
	:param to_session: Fixture to get Traffic Ops session.
	:param object_type: api call name for get request.
	:param query_params: query params for api get request.
	:returns: Api data for the corresponding api request.
    """
	api_get_response: tuple[Union[dict[str, object], list[
		Union[dict[str, object], list[object], Primitive]], Primitive], requests.Response] = getattr(
		to_session, f"get_{object_type}")(query_params=query_params)
	return api_response_data(api_get_response, create_check = False)


def create_if_not_exists(to_session: TOSession, object_type: str,
			 data: dict[str, Any]) -> Union[dict[str, Any], None]:
	"""
	Hits Post request of the given endpoint with the given data.
	:param to_session: Fixture to get Traffic Ops session.
	:param object_type: api call name for post request.
	:param data: Post data for api post request.
	:returns: Api data for the corresponding api request.
	"""
	api_post_response: tuple[Union[dict[str, object], list[
		Union[dict[str, object], list[object], Primitive]], Primitive], requests.Response] = getattr(
		to_session, f"create_{object_type}")(data=data)
	return api_response_data(api_post_response, create_check = True)


def create_or_get_existing(to_session: TOSession, get_object_type: str, post_object_type: str, data:
	dict[str, Any], query_params: Optional[dict[str, Any]] = None) -> Union[dict[str, Any], None]:
	"""
	Get Api data of the given endpoint with the given query params if it exists. If not, create it.
	:param to_session: Fixture to get Traffic Ops session.
	:param get_object_type: api call name for get request.
	:param post_object_type: api call name for post request.
	:param query_params: query params for api get request.
	:returns: Api data for the corresponding api request.
	@param data: 
	"""
	existing_object = get_existing_object(to_session, get_object_type, query_params)
	return existing_object or create_if_not_exists(to_session, post_object_type, data)


def generate_unique_data(to_session: TOSession, base_name: str, object_type: str,
			 query_key=None)-> str:
	"""
	Generate unique data for the given endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param object_type: api call name for get request.
	:param base_name: Base name for get request.
	:param query_key: Hitting get request using specific query key.
	:returns: Unique name for the corresponding api request.
	@param data: 
	"""
	unique_name = base_name
	if query_key is None:
		query_key = "name"
	while True:
		try:
			response = getattr(to_session, f"get_{object_type}")(query_params={query_key:unique_name})
			# Check if query params works for the corresponding api.
			check_data = response[0]
			if len(check_data) > 1:
				logger.info("API response returns all api data, query params is not working.")
				if not any(data.get(query_key) == unique_name for data in check_data):
					return unique_name
			elif len(check_data) == 0:
				raise ValueError("No Api response with the unique data")
		except ValueError:
			return unique_name
		unique_name = base_name[:4] + str(randint(0, 1000))


def check_template_data(template_data: Union[list[JSONData], tuple[JSONData, requests.Response]],
						name: str) -> dict[str, object]:
	"""
	Checks API request/response template data.
	:param template_data: Fixture to get template data from a prerequisites file.
	:param name: Endpoint name
	:returns: Verified endpoint data
	"""
	try:
		endpoint = template_data[0]
	except IndexError as e:
		raise TypeError(
			f"malformed  data; no {name} present in {name} array property") from e

	if not isinstance(endpoint, dict):
		raise TypeError(f"malformed data; {name} must be objects, not '{type(endpoint)}'")
	return endpoint


@pytest.fixture(name="cdn_post_data")
def cdn_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  pytestconfig: pytest.Config) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cdns endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get CDN request template data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	cdn = check_template_data(request_template_data["cdns"], "cdns")

	# Return new post data and post response from cdns POST request
	randstr = str(randint(0, 1000))
	try:
		name = cdn["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		cdn_name = name[:4] + randstr
		cdn["name"] = generate_unique_data(to_session=to_session, base_name=cdn_name, object_type="cdns")
		domain_name = cdn["domainName"]
		if not isinstance(domain_name, str):
			raise TypeError(f"domainName must be str, not '{type(domain_name)}")
		domainname = domain_name[:5] + randstr
		cdn["domainName"] = generate_unique_data(to_session=to_session, base_name=domainname,
					   object_type="cdns", query_key="domainName")
	except KeyError as e:
		raise TypeError(f"missing CDN property '{e.args[0]}'") from e

	logger.info("New cdn data to hit POST method %s", cdn)
	# Hitting cdns POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_cdn(data=cdn)
	resp_obj = check_template_data(response, "cdns")
	pytestconfig.cache.set("cdnDomainName",resp_obj.get("domainName"))
	yield resp_obj
	cdn_id = resp_obj.get("id")
	msg = to_session.delete_cdn_by_id(cdn_id=cdn_id)
	logger.info("Deleting cdn data... %s", msg)
	if  msg is None:
		logger.error("Cdn returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="cache_group_post_data")
def cache_group_data_post(to_session: TOSession, request_template_data: list[JSONData],
			 pytestconfig: pytest.Config) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cachegroup endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get Cachegroup data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	cache_group = check_template_data(request_template_data["cachegroup"], "cachegroup")
	# Return new post data and post response from cachegroups POST request
	randstr = str(randint(0, 1000))
	try:
		name = cache_group["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		cache_group_name = name[:4] + randstr
		cache_group["name"] = generate_unique_data(to_session=to_session, base_name=cache_group_name,
					     object_type="cachegroups")
		short_name = cache_group["shortName"]
		if not isinstance(short_name, str):
			raise TypeError(f"shortName must be str, not '{type(short_name)}")
		cache_group["shortName"] = short_name[:5] + randstr
	except KeyError as e:
		raise TypeError(f"missing Cache group property '{e.args[0]}'") from e

	# Check if type already exists, otherwise create it
	type_id = pytestconfig.cache.get("typeId", default=None)
	if type_id:
		cache_group["typeId"] = type_id
	else:
		type_data = check_template_data(request_template_data["types"], "types")
		type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "cachegroup"})
		cache_group["typeId"] = type_object["id"]

	logger.info("New cachegroup data to hit POST method %s", cache_group)
	# Hitting cachegroup POST method
	response: tuple[JSONData, requests.Response] = to_session.create_cachegroups(data=cache_group)
	resp_obj = check_template_data(response, "cachegroup")
	yield resp_obj
	cachegroup_id = resp_obj.get("id")
	msg = to_session.delete_cachegroups(cache_group_id=cachegroup_id)
	logger.info("Deleting cachegroup data... %s", msg)
	if  msg is None:
		logger.error("Cachegroup returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="parameter_post_data")
def parameter_data_post(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for parameters endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get CDN request template data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	parameter = check_template_data(request_template_data["parameters"], "parameters")
	# Return new post data and post response from parameters POST request
	randstr = str(randint(0, 1000))
	try:
		name = parameter["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		parameter_name = name[:4] + randstr
		parameter["name"] = generate_unique_data(to_session=to_session, base_name=parameter_name,
					   object_type="parameters")
		value = parameter["value"]
		if not isinstance(value, str):
			raise TypeError(f"value must be str, not '{type(value)}")
		parameter_value = value[:5] + randstr
		parameter["value"] = generate_unique_data(to_session=to_session, base_name=parameter_value,
					    object_type="parameters", query_key="value")
	except KeyError as e:
		raise TypeError(f"missing Parameter property '{e.args[0]}'") from e

	logger.info("New parameter data to hit POST method %s", parameter)
	# Hitting cdns POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_parameter(data=parameter)
	resp_obj = check_template_data(response, "parameter")
	yield resp_obj
	parameter_id = resp_obj.get("id")
	msg = to_session.delete_parameter(parameter_id=parameter_id)
	logger.info("Deleting parameter data... %s", msg)
	if  msg is None:
		logger.error("Parameter returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="role_post_data")
def role_data_post(to_session: TOSession, request_template_data: list[JSONData]
			) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for roles endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get role data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	role = check_template_data(request_template_data["roles"], "roles")

	# Return new post data and post response from roles POST request
	randstr = str(randint(0, 1000))
	try:
		name = role["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		role_name = name[:4] + randstr
		role["name"] = generate_unique_data(to_session=to_session, base_name=role_name,
				      object_type="roles")
	except KeyError as e:
		raise TypeError(f"missing Role property '{e.args[0]}'") from e

	logger.info("New role data to hit POST method %s", role)
	# Hitting roles POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_role(data=role)
	resp_obj = check_template_data(response, "role")
	yield resp_obj
	role_name = resp_obj.get("name")
	msg = to_session.delete_role(query_params={"name": role_name})
	logger.info("Deleting role data... %s", msg)
	if msg is None:
		logger.error("Role returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="profile_post_data")
def profile_data_post(to_session: TOSession, request_template_data: list[JSONData],
		      cdn_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for profile endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	profile = check_template_data(request_template_data["profiles"], "profiles")
	# Return new post data and post response from profiles POST request
	randstr = str(randint(0, 1000))
	try:
		name = profile["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		profile_name = name[:4] + randstr
		profile["name"] = generate_unique_data(to_session=to_session, base_name=profile_name,
					 object_type="profiles")
	except KeyError as e:
		raise TypeError(f"missing Profile property '{e.args[0]}'") from e

	# Check if cdn already exists, otherwise create it
	profile["cdn"] = cdn_post_data["id"]
	logger.info("New profile data to hit POST method %s", profile)

	# Hitting profile POST method
	response: tuple[JSONData, requests.Response] = to_session.create_profile(data=profile)
	resp_obj = check_template_data(response, "profile")
	yield resp_obj
	profile_id = resp_obj.get("id")
	msg = to_session.delete_profile_by_id(profile_id=profile_id)
	logger.info("Deleting profile data... %s", msg)
	if msg is None:
		logger.error("Profile returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="tenant_post_data")
def tenant_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  pytestconfig: pytest.Config) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for tenants endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get tenant request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	tenant = check_template_data(request_template_data["tenants"], "tenants")

	# Return new post data and post response from tenants POST request
	randstr = str(randint(0, 1000))
	name = pytestconfig.cache.get("tenantName", default=None)
	if name:
		tenant["name"] = name
	else:
		try:
			name = tenant["name"]
			if not isinstance(name, str):
				raise TypeError(f"name must be str, not '{type(name)}'")
			tenant_name = name[:4] + randstr
			tenant["name"] = generate_unique_data(to_session=to_session, base_name=tenant_name,
					 object_type="tenants")
		except KeyError as e:
			raise TypeError(f"missing tenant property '{e.args[0]}'") from e

	logger.info("New tenant data to hit POST method %s", tenant)
	# Hitting tenants POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_tenant(data=tenant)
	resp_obj = check_template_data(response, "tenant")
	yield resp_obj
	tenant_id = resp_obj.get("id")
	msg = to_session.delete_tenant(tenant_id=tenant_id)
	logger.info("Deleting tenant data... %s", msg)
	if msg is None:
		logger.error("Tenant returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="server_capabilities_post_data")
def server_capabilities_data_post(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for server_capabilities endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get server_capabilities data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	server_capabilities = check_template_data(request_template_data["server_capabilities"],
					   "server_capabilities")

	# Return new post data and post response from server_capabilities POST request
	randstr = str(randint(0, 1000))
	try:
		name = server_capabilities["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		server_capabilities_name = name[:3] + randstr
		server_capabilities["name"] = generate_unique_data(to_session=to_session,
						     base_name=server_capabilities_name, object_type="server_capabilities")
	except KeyError as e:
		raise TypeError(f"missing server_capabilities property '{e.args[0]}'") from e

	logger.info("New server_capabilities data to hit POST method %s", server_capabilities)
	# Hitting server_capabilities POST method
	response: tuple[
		JSONData, requests.Response] = to_session.create_server_capabilities(data=server_capabilities)
	resp_obj = check_template_data(response, "server_capabilities")
	yield resp_obj
	server_capability_name = resp_obj.get("name")
	msg = to_session.delete_server_capabilities(query_params={"name": server_capability_name})
	logger.info("Deleting server_capabilities data %s", msg)
	if  msg is None:
		logger.error("server_capabilities returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_server_capabilities_contract")


@pytest.fixture(name="division_post_data")
def division_data_post(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for divisions endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get divisions data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	division = check_template_data(request_template_data["divisions"], "divisions")

	# Return new post data and post response from division POST request
	randstr = str(randint(0, 1000))
	try:
		name = division["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		division_name = name[:4] + randstr
		division["name"] = generate_unique_data(to_session=to_session, base_name=division_name,
					  object_type="divisions")
	except KeyError as e:
		raise TypeError(f"missing Parameter property '{e.args[0]}'") from e

	logger.info("New division data to hit POST method %s", division)
	# Hitting division POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_division(data=division)
	resp_obj = check_template_data(response, "divisions")
	yield resp_obj
	division_id = resp_obj.get("id")
	msg = to_session.delete_division(division_id=division_id)
	logger.info("Deleting division data... %s", msg)
	if msg is None:
		logger.error("Division returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="region_post_data")
def region_data_post(to_session: TOSession, request_template_data: list[JSONData],
			  division_post_data: dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for region endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get region data from a prerequisites file.
  	:returns: Sample POST data and the actual API response.
	"""

	region = check_template_data(request_template_data["regions"], "regions")

	# Return new post data and post response from regions POST request
	randstr = str(randint(0, 1000))
	try:
		name = region["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		region_name = name[:4] + randstr
		region["name"] = generate_unique_data(to_session=to_session, base_name=region_name,
					object_type="regions")
	except KeyError as e:
		raise TypeError(f"missing Region property '{e.args[0]}'") from e

	# Check if division already exists, otherwise create it
	region["division"] = division_post_data["id"]
	region["divisionName"] = division_post_data["name"]

	logger.info("New region data to hit POST method %s", region)
	# Hitting region POST method
	response: tuple[JSONData, requests.Response] = to_session.create_region(data=region)
	resp_obj = check_template_data(response, "regions")
	yield resp_obj
	region_name = resp_obj.get("name")
	msg = to_session.delete_region(query_params={"name": region_name})
	logger.info("Deleting region data... %s", msg)
	if msg is None:
		logger.error("Region returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="phys_locations_post_data")
def phys_locations_data_post(to_session: TOSession, request_template_data: list[JSONData],
			 region_post_data: dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for phys_location endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get phys_location data from a prerequisites file.
	:param region_post_data:
  	:returns: Sample POST data and the actual API response.
	"""

	phys_locations = check_template_data(request_template_data["phys_locations"], "phys_locations")

	# Return new post data and post response from phys_locations POST request
	randstr = str(randint(0, 1000))
	try:
		name = phys_locations["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		phys_locations_name = name[:4] + randstr
		phys_locations["name"] = generate_unique_data(to_session=to_session,
						base_name=phys_locations_name, object_type="physical_locations")
		short_name = phys_locations["shortName"]
		if not isinstance(short_name, str):
			raise TypeError(f"shortName must be str, not '{type(short_name)}'")
		shortname = short_name[:4] + randstr
		phys_locations["shortName"] = generate_unique_data(to_session=to_session,
						base_name=shortname, object_type="physical_locations", query_key="shortName")
	except KeyError as e:
		raise TypeError(f"missing Phys_location property '{e.args[0]}'") from e

	# Check if region already exists, otherwise create it
	region_id = region_post_data["id"]
	if not isinstance(region_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")
	phys_locations["regionId"] = region_id

	logger.info("New Phys_locations data to hit POST method %s", phys_locations)
	# Hitting region POST method
	response: tuple[JSONData, requests.Response] = to_session.create_physical_locations(
		data=phys_locations)
	resp_obj = check_template_data(response, "phys_locations")
	yield resp_obj
	phys_location_id = resp_obj.get("id")
	msg = to_session.delete_physical_location(physical_location_id=phys_location_id)
	logger.info("Deleting physical locations data... %s", msg)
	if msg is None:
		logger.error("Physical location returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="server_post_data")
def server_data_post(to_session: TOSession, request_template_data: list[JSONData],
		profile_post_data: dict[str, object], cache_group_post_data: dict[str, object],
		status_post_data: dict[str, object], phys_locations_post_data: dict[str, object],
		pytestconfig: pytest.Config)-> dict[str, object]:
	"""
	PyTest Fixture to create POST data for server endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile data from a prerequisites file.
	:param phys_locations_post_data:
	:returns: Sample POST data and the actual API response.
	"""
	server = check_template_data(request_template_data["servers"], "servers")

	# Check if type already exists, otherwise create it
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "server"})
	type_id = type_object["id"]
	edge_type_id = pytestconfig.cache.get("edgeTypeId", default=None)
	if edge_type_id:
		server["typeId"] = edge_type_id
	else:
		server["typeId"] = type_id
	
	if edge_type_id:
		server["typeId"] = edge_type_id
	else:
		server["typeId"] = type_id

	pytestconfig.cache.set("typeId", type_id)

	server["cachegroupId"]= cache_group_post_data["id"]

	# Check if cdn already exists, otherwise create it
	server["cdnId"] = profile_post_data["cdn"]
	server["domainName"] = pytestconfig.cache.get("cdnDomainName", default=None)

	# Check if profile with cdn already exists, otherwise create it
	server["profileNames"] = [profile_post_data["name"]]

	# Check if status already exists, otherwise create it
	server["statusId"] = status_post_data["id"]

	# Check if physical location with region already exists, otherwise create it
	physical_location_id = phys_locations_post_data["id"]
	if not isinstance(physical_location_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")
	server["physLocationId"] = physical_location_id

	logger.info("New server data to hit POST method %s", server)
	# Hitting server POST method
	response: tuple[JSONData, requests.Response] = to_session.create_server(data=server)
	resp_obj = check_template_data(response, "server")
	yield resp_obj
	server_id = resp_obj.get("id")
	msg = to_session.delete_server_by_id(server_id=server_id)
	logger.info("Deleting servers data... %s", msg)
	if msg is None:
		logger.error("Server returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="delivery_services_post_data")
def delivery_services_data_post(to_session: TOSession, request_template_data: list[JSONData],
				tenant_post_data: dict[str, object], profile_post_data: dict[str, object],
		      pytestconfig: pytest.Config) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for server endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	delivery_services = check_template_data(
		request_template_data["delivery_services"], "delivery_services")

	randstr = str(randint(0, 1000))
	try:
		xml_id = delivery_services["xmlId"]
		if not isinstance(xml_id, str):
			raise TypeError(f"xmlId must be str, not '{type(xml_id)}'")
		xmlid= xml_id[:4] + randstr
		delivery_services["xmlId"] = generate_unique_data(to_session=to_session, base_name=xmlid,
						    object_type="deliveryservices", query_key="xmlId")
	except KeyError as e:
		raise TypeError(f"missing delivery_services property '{e.args[0]}'") from e

	# Check if profile with cdn already exists, otherwise create it
	delivery_services["profileId"] = profile_post_data["id"]

	# Check if cdn already exists, otherwise create it
	delivery_services["cdnId"] = profile_post_data["cdn"]

	# Check if tenant already exists, otherwise create it
	pytestconfig.cache.set("tenantName", "root")
	delivery_services["tenantId"] = tenant_post_data["id"]

	# Check if type already exists, otherwise create it
	type_data = {"name": "HTTP", "useInTable":"deliveryservice"}
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"name": "HTTP", "useInTable":"deliveryservice"})
	delivery_services["typeId"] = type_object["id"]
	delivery_services["type"] = type_object["name"]

	logger.info("New delivery_services data to hit POST method %s", delivery_services)
	# Hitting delivery_services POST method
	response: tuple[JSONData, requests.Response] = to_session.create_deliveryservice(
		data=delivery_services)
	resp_obj = check_template_data(response[0], "delivery_services")
	yield resp_obj
	delivery_service_id = resp_obj.get("id")
	msg = to_session.delete_deliveryservice_by_id(delivery_service_id=delivery_service_id)
	logger.info("Deleting delivery service data... %s", msg)
	if msg is None:
		logger.error("delivery service returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="origin_post_data")
def origin_data_post(to_session: TOSession, request_template_data: list[JSONData],
		     delivery_services_post_data: dict[str, object], tenant_post_data: dict[str, object]
		      ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for origins endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	origin = check_template_data(request_template_data["origins"], "origins")

	randstr = str(randint(0, 1000))
	try:
		name = origin["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		origin_name = name[:4] + randstr
		origin["name"] = generate_unique_data(to_session=to_session, base_name=origin_name,
					object_type="origins")
	except KeyError as e:
		raise TypeError(f"missing origin property '{e.args[0]}'") from e

	# Check if delivery_service already exists, otherwise create it
	delivery_services_id = delivery_services_post_data["id"]
	if not isinstance(delivery_services_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")
	origin["deliveryServiceId"] = delivery_services_id

	# Check if tenant already exists, otherwise create it
	tenant_id = tenant_post_data["id"]
	if not isinstance(tenant_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")
	origin["tenantId"] = tenant_id

	logger.info("New origin data to hit POST method %s", origin)
	# Hitting origins POST method
	response: tuple[JSONData, requests.Response] = to_session.create_origins(data=origin)
	resp_obj = check_template_data(response, "origins")
	yield resp_obj
	origin_id = resp_obj.get("id")
	msg = to_session.delete_origins(query_params={"id": origin_id})
	logger.info("Deleting origin data... %s", msg)
	if msg is None:
		logger.error("Origin returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="status_post_data")
def status_data_post(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for statuses endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get Status request template data from a prerequisite file.
	:returns: Sample POST data and the actual API response.
	"""
	status = check_template_data(request_template_data["status"], "status")

	# Return new post data and post response from statuses POST request
	randstr = str(randint(0, 1000))
	try:
		name = status["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		status_name = name[:4] + randstr
		status["name"] = generate_unique_data(to_session=to_session, base_name=status_name,
					object_type="statuses")
	except KeyError as e:
		raise TypeError(f"missing Status property '{e.args[0]}'") from e

	logger.info("New status data to hit POST method %s", status)
	# Hitting statuses POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_statuses(data=status)
	resp_obj = check_template_data(response, "statuses")
	yield resp_obj
	status_id = resp_obj.get("id")
	msg = to_session.delete_status_by_id(status_id=status_id)
	logger.info("Deleting status data... %s", msg)
	if msg is None:
		logger.error("Status returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="asn_post_data")
def asn_data_post(to_session: TOSession, request_template_data: list[JSONData],
		      cache_group_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for asn endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get asn data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	asn = check_template_data(request_template_data["asns"], "asns")
	# Return new post data and post response from asns POST request
	randstr = randint(0, 1000)
	asn["asn"] = randstr

	# Check if cachegroup already exists, otherwise create it
	asn["cachegroupId"] = cache_group_post_data["id"]
	logger.info("New profile data to hit POST method %s", asn)

	# Hitting asns POST method
	response: tuple[JSONData, requests.Response] = to_session.create_asn(data=asn)
	resp_obj = check_template_data(response, "asn")
	yield resp_obj
	asn_id = resp_obj.get("id")
	msg = to_session.delete_asn(query_params={"id": asn_id})
	logger.info("Deleting asn data... %s", msg)
	if msg is None:
		logger.error("asn returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="job_post_data")
def job_data_post(to_session: TOSession, request_template_data: list[JSONData],
		     delivery_services_post_data: dict[str, object],
		      ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for jobss endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get job data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	job = check_template_data(request_template_data["jobs"], "jobs")

	# Check if delivery_service already exists, otherwise create it
	delivery_services_name = delivery_services_post_data["xmlId"]
	if not isinstance(delivery_services_name, str):
		raise TypeError("malformed API response; 'displayName' property not a string")
	job["deliveryService"] = delivery_services_name

	logger.info("New job data to hit POST method %s", job)
	# Hitting jobs POST method
	response: tuple[JSONData, requests.Response] = to_session.create_job(data=job)
	resp_obj = check_template_data(response, "jobs")
	yield resp_obj
	job_id = resp_obj.get("id")
	msg = to_session.delete_job(query_params={"id": job_id})
	logger.info("Deleting job data... %s", msg)
	if msg is None:
		logger.error("job returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="coordinate_post_data")
def coordinate_data_post(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for coordinates endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get coordinate request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	coordinate = check_template_data(request_template_data["coordinates"], "coordinates")

	# Return new post data and post response from coordinates POST request
	randstr = str(randint(0, 1000))
	try:
		name = coordinate["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		coordinate_name = name[:4] + randstr
		coordinate["name"] = generate_unique_data(to_session=to_session,
					    base_name=coordinate_name, object_type="coordinates")
	except KeyError as e:
		raise TypeError(f"missing coordinate property '{e.args[0]}'") from e

	logger.info("New coordinate data to hit POST method %s", coordinate)
	# Hitting coordinates POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_coordinates(data=coordinate)
	resp_obj = check_template_data(response, "coordinate")
	yield resp_obj
	coordinate_id = resp_obj.get("id")
	msg = to_session.delete_coordinates(query_params={"id": coordinate_id})
	logger.info("Deleting Coordinate data... %s", msg)
	if msg is None:
		logger.error("coordinate returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="user_post_data")
def user_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  tenant_post_data:dict[str, object], db_connection: psycopg2.connect) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for users endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get users request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	user = check_template_data(request_template_data["users"], "users")

	# Return new post data and post response from users POST request
	randstr = str(randint(0, 1000))
	try:
		username = user["username"]
		if not isinstance(username, str):
			raise TypeError(f"username must be str, not '{type(username)}'")
		unique_name = username[:4] + randstr
		user["username"] = generate_unique_data(to_session=to_session, base_name=unique_name,
					  object_type="users", query_key="username")
		user_email = user["email"]
		if not isinstance(user_email, str):
			raise TypeError(f"user email must be str, not '{type(user_email)}'")
		email = randstr + user_email
		user["email"] = generate_unique_data(to_session=to_session, base_name=email,
					  object_type="users", query_key="email")
	except KeyError as e:
		raise TypeError(f"missing user property '{e.args[0]}'") from e
	user["tenantId"] = tenant_post_data["id"]

	logger.info("New user data to hit POST method %s", user)
	# Hitting users POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_user(data=user)
	resp_obj = check_template_data(response, "user")
	yield resp_obj
	user_id = resp_obj.get("id")
	# Create a cursor object to interact with the database
	cursor = db_connection.cursor()
	cursor.execute("DELETE FROM tm_user WHERE id = %s;", (user_id,))
	# Commit the changes
	db_connection.commit()
	# Close the cursor
	cursor.close()


@pytest.fixture(name="topology_post_data")
def topology_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  server_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for topologies endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get topology request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	topology = check_template_data(request_template_data["topologies"], "topologies")

	# Return new post data and post response from topologies POST request
	randstr = str(randint(0, 1000))
	try:
		name = topology["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		topology_name = name[:4] + randstr
		topology["name"] = generate_unique_data(to_session=to_session, base_name=topology_name,
					  object_type="topologies")
	except KeyError as e:
		raise TypeError(f"missing topology property '{e.args[0]}'") from e

	cachegroup_name = server_post_data["cachegroup"]
	topology["nodes"][0]["cachegroup"] = cachegroup_name
	logger.info("New topology data to hit POST method %s", topology)
	# Hitting topology POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_topology(data=topology)
	resp_obj = check_template_data(response, "topology")
	yield resp_obj
	topology_name = resp_obj.get("name")
	msg = to_session.delete_topology(name=topology_name)
	logger.info("Deleting topology data... %s", msg)
	if msg is None:
		logger.error("topology returned by Traffic Ops is missing an 'name' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="cdn_lock_post_data")
def cdn_lock_data_post(to_session: TOSession, request_template_data: list[JSONData],
		user_post_data:dict[str, object], cdn_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cdn_locks endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get cdn_locks request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	cdn_lock = check_template_data(request_template_data["cdn_locks"], "cdn_locks")

	# Return new post data and post response from cdn_locks POST request
	cdn_lock["cdn"] = cdn_post_data["name"]
	cdn_lock["sharedUserNames"][0] = user_post_data["username"]
	logger.info("New cdn_lock data to hit POST method %s", cdn_lock)
	# Hitting cdn_locks POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_cdn_lock(data=cdn_lock)
	resp_obj = check_template_data(response, "cdn_lock")
	yield resp_obj
	cdn_name = resp_obj.get("cdn")
	msg = to_session.delete_cdn_lock(query_params={"cdn":cdn_name})
	logger.info("Deleting cdn_lock data... %s", msg)
	if msg is None:
		logger.error("cdn_lock returned by Traffic Ops is missing an 'cdn' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="cdn_notification_post_data")
def cdn_notification_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  cdn_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cdn_notifications endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get cdn_notification request template.
	:returns: Sample POST data and the actual API response.
	"""

	cdn_notification = check_template_data(
		request_template_data["cdn_notifications"], "cdn_notifications")
	# Return new post data and post response from cdn_notifications POST request
	cdn_notification["cdn"] = cdn_post_data["name"]
	logger.info("New cdn_notification data to hit POST method %s", cdn_notification)
	# Hitting cdn_notification POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_cdn_notification(
		data=cdn_notification)
	resp_obj = check_template_data(response, "cdn_notification")
	yield resp_obj
	notification_id = resp_obj.get("id")
	msg = to_session.delete_cdn_notification(query_params={"id":notification_id})
	logger.info("Deleting cdn_notification data... %s", msg)
	if msg is None:
		logger.error("cdn_notfication returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="deliveryservice_request_post_data")
def deliveryservice_request_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  delivery_services_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for deliveryservice_request endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get deliveryservice_request request template.
	:returns: Sample POST data and the actual API response.
	"""

	deliveryservice_request = check_template_data(
		request_template_data["deliveryservice_requests"], "deliveryservice_requests")

	# Return new post data and post response from deliveryservice_request POST request
	keys = ["displayName", "xmlId", "id", "cdnId", "tenantId", "type", "typeId"]
	for key in keys:
		deliveryservice_request["requested"][key] = delivery_services_post_data[key]
	logger.info("New deliveryservice_request data to hit POST method %s", deliveryservice_request)
	# Hitting deliveryservice_request POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_deliveryservice_request(
		data=deliveryservice_request)
	resp_obj = check_template_data(response, "deliveryservice_request")
	yield resp_obj
	deliveryservice_request_id = resp_obj.get("id")
	msg = to_session.delete_deliveryservice_request(query_params={"id":deliveryservice_request_id})
	logger.info("Deleting deliveryservice_request data... %s", msg)
	if msg is None:
		logger.error("deliveryservice_request returned by Traffic Ops is missing an 'id' property")


@pytest.fixture(name="steering_post_data")
def steering_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  delivery_services_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for steering endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get steering request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	steering = check_template_data(request_template_data["steering"], "steering")

	# Return new post data and post response from steering POST request
	ds_get_response = to_session.get_deliveryservices()
	ds_data = ds_get_response[0][0]
	delivery_service_id = ds_data.get("id")

	steering["targetId"] = delivery_services_post_data["id"]
	# Check if type already exists, otherwise create it
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "steering_target"})
	steering["typeId"]= type_object["id"]

	logger.info("New steering data to hit POST method %s", steering)
	# Hitting steering POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_steering_targets(
		delivery_service_id=delivery_service_id, data=steering)
	resp_obj = check_template_data(response, "steering")
	yield resp_obj
	deliveryservice_id = resp_obj.get("deliveryServiceId")
	target_id = resp_obj.get("targetId")
	msg = to_session.delete_steering_targets(
		delivery_service_id=deliveryservice_id, target_id=target_id)
	logger.info("Deleting Steering data... %s", msg)
	if msg is None:
		logger.error("Steering returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="delivery_service_sslkeys_post_data")
def delivery_service_sslkeys_data_post(to_session: TOSession, request_template_data: list[JSONData],
		cdn_post_data:dict[str, object], delivery_services_post_data:dict[str, object]
		) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for delivery_service_sslkeys_post_data endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get delivery_service_sslkeys request template.
	:returns: Sample POST data and the actual API response.
	"""

	delivery_service_sslkeys = check_template_data(
		request_template_data["delivery_service_sslkeys"], "delivery_service_sslkeys")

	# Return new post data and post response from delivery_service_sslkeys POST request
	delivery_service_sslkeys["key"] = delivery_services_post_data["xmlId"]
	delivery_service_sslkeys["cdn"] = cdn_post_data["name"]
	logger.info("New delivery_service_sslkeys data to hit POST method %s", delivery_service_sslkeys)
	# Hitting delivery_service_sslkeys POST methed
	response: tuple[JSONData, requests.Response] = to_session.generate_deliveryservice_ssl_keys(
		data=delivery_service_sslkeys)
	yield delivery_service_sslkeys
	deliveryservice_xml_id = delivery_service_sslkeys["key"]
	msg = to_session.delete_deliveryservice_ssl_keys_by_xml_id(xml_id=deliveryservice_xml_id)
	logger.info("Deleting delivery_service_sslkeys data... %s", msg)
	if msg is None:
		logger.error("delivery_service_sslkeys returned by Traffic Ops is missing an 'xmlId' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="delivery_services_regex_post_data")
def delivery_services_regex_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  delivery_services_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for delivery_services_regex endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get delivery_services_regex request template.
	:returns: Sample POST data and the actual API response.
	"""

	delivery_services_regex = check_template_data(
		request_template_data["delivery_services_regex"], "delivery_services_regex")

	# Return new post data and post response from delivery_services_regex POST request

	delivery_service_id = delivery_services_post_data["id"]
	# Check if type already exists, otherwise create it
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "regex"})
	delivery_services_regex["type"]= type_object["id"]

	logger.info("New delivery_services_regex data to hit POST method %s", delivery_services_regex)
	# Hitting delivery_services_regex POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_deliveryservice_regexes(
		delivery_service_id=delivery_service_id, data=delivery_services_regex)
	resp_obj = check_template_data(response, "delivery_services_regex")
	yield [delivery_service_id,resp_obj]
	regex_id = resp_obj.get("id")
	msg = to_session.delete_deliveryservice_regex_by_regex_id(
		delivery_service_id=delivery_service_id, delivery_service_regex_id=regex_id)
	logger.info("Deleting delivery_services_regex data... %s", msg)
	if msg is None:
		logger.error("delivery_services_regex returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="cdn_federation_post_data")
def cdn_federation_data_post(to_session: TOSession, request_template_data: list[JSONData],
		  cdn_post_data:dict[str, object], user_post_data:dict[str, object],
		  federation_resolver_post_data:dict[str, object],
		  delivery_services_post_data: dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cdn_name_federations endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get federations request template.
	:returns: Sample POST data and the actual API response.
	"""

	cdn_federation = check_template_data(request_template_data["cdn_federation"], "cdn_federation")
	# Return new post data and post response from cdn_federation POST request
	cdn_name = cdn_post_data["name"]

	logger.info("New federations data to hit POST method %s", cdn_federation)
	# Hitting cdn_federation POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_federation_in_cdn(cdn_name=cdn_name, data= cdn_federation)
	cdn_federation_resp_obj = check_template_data(response, "cdn_federation")
	federation_id = cdn_federation_resp_obj.get("id")

	#Assign created federation to a user
	user_id = user_post_data["id"]
	user_federation = check_template_data(request_template_data["user_federation"], "user_federation")
	user_federation["userIds"][0] = user_id
	response: tuple[JSONData, requests.Response] = to_session.create_federation_user(federation_id=federation_id, data=user_federation)
	user_federation_resp_obj = check_template_data(response, "user_federation")

	#Assign the federation to a delivery_service
	delivery_service_id = delivery_services_post_data["id"]
	delivery_service_federation = check_template_data(request_template_data["delivery_service_federation"], "delivery_service_federation")
	delivery_service_federation["dsIds"][0] = delivery_service_id
	response: tuple[JSONData, requests.Response] = to_session.assign_delivery_services_to_federations(federation_id=federation_id, data=delivery_service_federation)
	delivery_service_federation_resp_obj = check_template_data(response, "delivery_service_federation")

	#Assign a federation resolver to created federation
	federation_resolver_id = federation_resolver_post_data["id"]
	federation_federation_resolver = check_template_data(request_template_data["federation_federation_resolver"], "federation_federation_resolver")
	federation_federation_resolver["fedResolverIds"][0] = federation_resolver_id
	response: tuple[JSONData, requests.Response] = to_session.assign_federation_resolver_to_federations(federation_id=federation_id, data=federation_federation_resolver)
	federation_federation_resolver_resp_obj = check_template_data(response, "federation_federation_resolver")

	yield [cdn_name, cdn_federation_resp_obj, cdn_federation, federation_id, delivery_service_federation_resp_obj]

	msg = to_session.delete_federation_in_cdn(cdn_name=cdn_name, federation_id=federation_id)
	logger.info("Deleting cdn_federation dara... %s", msg)
	if msg is None:
		logger.error("cdn_federation returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="delivery_service_required_capabilities_post_data")
def delivery_service_required_capabilities_data_post(to_session: TOSession,
		request_template_data: list[JSONData], delivery_services_post_data:dict[str, object],
		server_capabilities_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for delivery_service_required_capabilities endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get delivery_service_required_capabilities request template.
	:returns: Sample POST data and the actual API response.
	"""

	delivery_service_required_capabilities = check_template_data(
		request_template_data["delivery_service_required_capabilities"], "delivery_service_required_capabilities")

	# Return new post data and post response from delivery_service_required_capabilities POST request
	deliveryServiceID = delivery_services_post_data["id"]
	requiredCapability = server_capabilities_post_data["name"]
	delivery_service_required_capabilities["deliveryServiceID"] = deliveryServiceID
	delivery_service_required_capabilities["requiredCapability"] = requiredCapability

	logger.info("New delivery_service_required_capabilities data to hit POST method %s",
	     delivery_service_required_capabilities)
	# Hitting delivery_service_required_capabilities POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_deliveryservices_required_capabilities(
		data=delivery_service_required_capabilities)
	resp_obj = check_template_data(response, "delivery_service_required_capabilities")
	yield resp_obj
	msg = to_session.delete_deliveryservices_required_capabilities(
		query_params={"deliveryServiceID":deliveryServiceID,"requiredCapability":requiredCapability})
	logger.info("Deleting delivery_service_required_capabilities data... %s", msg)
	if msg is None:
		logger.error(
		"delivery_service_required_capabilities returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="delivery_service_request_comments_post_data")
def delivery_service_request_comments_data_post(to_session: TOSession,
		request_template_data: list[JSONData],
		deliveryservice_request_post_data:dict[str, object]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for delivery_service_request_comments endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get delivery_service_request_comments request template.
	:returns: Sample POST data and the actual API response.
	"""

	delivery_service_request_comments = check_template_data(
		request_template_data["delivery_service_request_comments"], "delivery_service_request_comments")

	# Return new post data and post response from delivery_service_request_comments POST request
	delivery_service_request_id = deliveryservice_request_post_data["id"]
	delivery_service_request_comments["deliveryServiceRequestId"]= delivery_service_request_id

	logger.info("New delivery_service_request_comments data to hit POST method %s",
	     delivery_service_request_comments)
	# Hitting delivery_service_request_comments POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_deliveryservice_request_comment(
		data=delivery_service_request_comments)
	resp_obj = check_template_data(response, "delivery_service_request_comments")
	yield resp_obj
	request_comment_id = resp_obj.get("id")
	msg = to_session.delete_deliveryservice_request_comment(query_params={"id":request_comment_id})
	logger.info("Deleting delivery_service_request_comments data... %s", msg)
	if msg is None:
		logger.error("delivery_service_request_comments returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="profile_parameters_post_data")
def profile_parameters_post_data(to_session: TOSession, request_template_data: list[JSONData],
		profile_post_data:dict[str, object], parameter_post_data:dict[str, object]
		) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for profile parameters endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile parameters request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	profile_parameters = check_template_data(request_template_data["profile_parameters"], "profile_parameters")

	# Return new post data and post response from profile parameters POST request
	profile_get_response = to_session.get_profiles()
	profile_data = profile_get_response [0][0]
	profile_id = profile_data.get("id")

	profile_parameters["profileId"] = profile_post_data["id"]
	profile_parameters["parameterId"] = parameter_post_data["id"]

	logger.info("New profile_parameter data to hit POST method %s", profile_parameters)

	# Hitting profile parameters POST method
	response: tuple[JSONData, requests.Response] = to_session.associate_paramater_to_profile(profile_id=profile_id, data=profile_parameters)
	resp_obj = check_template_data(response, "profile_parameters")
	yield resp_obj
	profile_id = resp_obj.get("profileId")
	parameter_id = resp_obj.get("parameterId")
	msg = to_session.delete_profile_parameter_association_by_id(profile_id=profile_id, parameter_id=parameter_id)
	logger.info("Deleting Profile Parameters data... %s", msg)
	if msg is None:
		logger.error("Profile Parameter returned by Traffic Ops is missing a 'profile_id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="service_category_post_data")
def service_category_data_post(to_session: TOSession,
		request_template_data: list[JSONData]) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for service_category endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get service_category request template.
	:returns: Sample POST data and the actual API response.
	"""

	service_category = check_template_data(
		request_template_data["service_category"], "service_category")

	# Return new post data and post response from service_category POST request
	service_category_name = service_category["name"]
	service_category["name"] = service_category_name + str(randint(0,1000))

	logger.info("New service_category data to hit POST method %s", service_category)
	# Hitting service_category POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_service_category(
		data=service_category)
	resp_obj = check_template_data(response, "service_category")
	yield resp_obj
	service_category_name = resp_obj.get("name")
	msg = to_session.delete_service_category(service_category_name=service_category_name)
	logger.info("Deleting service_category data... %s", msg)
	if msg is None:
		logger.error("service_category returned by Traffic Ops is missing an 'name' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="federation_resolver_post_data")
def federation_resolver_data_post(to_session: TOSession, request_template_data: list[JSONData]
				  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for federation_resolver endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get federation_resolver request template.
	:returns: Sample POST data and the actual API response.
	"""
	randstr = str(randint(0, 10))
	federation_resolver = check_template_data(
		request_template_data["federation_resolver"], "federation_resolver")

	# Check if type already exists, otherwise create it
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "federation"})
	federation_resolver["typeId"] = type_object["id"]
	federation_resolver["ipAddress"] = ".".join(map(str, (randint(0, 255)
                        for _ in range(4))))

	logger.info("New federation_resolver data to hit POST method %s", federation_resolver)
	# Hitting federation_resolver POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_federation_resolver(data=federation_resolver)
	resp_obj = check_template_data(response, "federation_resolver")
	yield resp_obj
	resolver_id = resp_obj.get("id")
	msg = to_session.delete_federation_resolver(query_params={"id":resolver_id})
	logger.info("Deleting federation_resolver data... %s", msg)
	if msg is None:
		logger.error("federation_resolver returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="static_dns_entries_post_data")
def static_dns_entries_data_post(to_session: TOSession, request_template_data: list[JSONData],
				 delivery_services_post_data:dict[str, object]
				  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for static_dns_entries endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get static_dns_entries request template.
	:returns: Sample POST data and the actual API response.
	"""
	static_dns_entries = check_template_data(
		request_template_data["static_dns_entries"], "static_dns_entries")

	# Check if type already exists, otherwise create it
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"useInTable": "staticdnsentry"})
	static_dns_entries["typeId"] = type_object["id"]
	static_dns_entries["deliveryServiceId"] = delivery_services_post_data["id"]

	logger.info("New static_dns_entries data to hit POST method %s", static_dns_entries)
	# Hitting static_dns_entries POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_staticdnsentries(data=static_dns_entries)
	resp_obj = check_template_data(response, "static_dns_entries")
	yield resp_obj
	static_dns_entries_id = resp_obj.get("id")
	msg = to_session.delete_staticdnsentries(query_params={"id":static_dns_entries_id})
	logger.info("Deleting static_dns_entries data... %s", msg)
	if msg is None:
		logger.error("static_dns_entries returned by Traffic Ops is missing an 'id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")

@pytest.fixture(name="edge_type_data")
def edge_data_type(pytestconfig: pytest.Config, request_template_data: list[JSONData], to_session: TOSession):
	type_data = check_template_data(request_template_data["types"], "types")
	type_object = create_or_get_existing(to_session, "types", "type", type_data,
				      {"name":"EDGE" , "useInTable":"server"})
	type_id = type_object["id"]
	pytestconfig.cache.set("edgeTypeId", type_id)

@pytest.fixture(name="server_server_capabilities_post_data")
def server_server_capabilities_data_post(to_session: TOSession, edge_type_data:None, request_template_data: list[JSONData],
		server_post_data:dict[str, object], server_capabilities_post_data:dict[str, object], pytestconfig: pytest.Config
		) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for server server capabilities endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get Server Server Capabilities request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	server_server_capabilities = check_template_data(request_template_data["server_server_capabilities"], "server_server_capabilities")

	# Return new post data and post response from server server capabilities POST request

	server_id = server_post_data.get("id")
	serverCapability = server_capabilities_post_data.get("name")
	server_server_capabilities["serverId"] = server_id
	server_server_capabilities["serverCapability"] = serverCapability


	logger.info("New server_server_capabilities data to hit POST method %s", server_server_capabilities)

	# Hitting server server capabilities POST method
	response: tuple[JSONData, requests.Response] = to_session.associate_server_capability_to_server(server_id=server_id, data=server_server_capabilities)
	resp_obj = check_template_data(response, "server_server_capabilities")
	yield resp_obj
	server_id = resp_obj.get("serverId")
	msg = to_session.delete_server_capability_association_to_server(query_params={"serverId":server_id, "serverCapability":serverCapability})
	logger.info("Deleting Server Server Capability data... %s", msg)
	if msg is None:
		logger.error("Server Server Capability returned by Traffic Ops is missing a 'server_id' property")
		pytest.fail("Response from delete request is empty, Failing test_case")


@pytest.fixture(name="logs_data")
def logs_data(to_session: TOSession, request_template_data: list[JSONData],
        ) -> dict[str, object]:
	"""
    PyTest Fixture to retrieve log data from logs endpoint.
    :param to_session: Fixture to get Traffic Ops session.
    :returns: Log data obtained from logs endpoint.
    """

	change_logs = check_template_data(request_template_data["logs"], "logs")


    # Hitting logs GET methed
	response: tuple[JSONData, requests.Response] = to_session.get_change_logs(data=change_logs)
	resp_obj = check_template_data(response[0], "change_logs")
	change_log_id = resp_obj.get("id")

	yield [change_log_id, resp_obj]
