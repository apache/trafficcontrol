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
import sys
import os
from random import randint
from typing import NamedTuple, Union, Optional, TypeAlias
from urllib.parse import urlparse

import pytest
import requests

from trafficops.tosession import TOSession
from trafficops.restapi import OperationError

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

JSONData: TypeAlias = Union[dict[str, object], list[object], bool, int, float, str | None]
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
		"--config",
		help="Path to configuration file.",
		default=os.path.join(os.path.dirname(__file__), "to_data.json")
	)
	parser.addoption(
		"--request-template",
		help="Path to request prerequisites file.",
		default=os.path.join(os.path.dirname(__file__), "request_template.json")
	)
	parser.addoption(
		"--response-template",
		help="Path to response prerequisites file.",
		default=os.path.join(os.path.dirname(__file__), "response_template.json")
	)

def coalesce_config(
	arg: object | None,
	file_key: str,
	file_contents: dict[str, object | None] | None,
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

	>>> parse_to_url("https://trafficops.example.test:420/api/5.270)
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
			"Traffic Ops password is not configured - use '--to-password', the config file, or an "
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


@pytest.fixture(name="request_template_data", scope="session")
def request_prerequiste_data(pytestconfig: pytest.Config, request: pytest.FixtureRequest
			  ) -> list[dict[str, object] | list[object] | primitive]:
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
	data: dict[
		str,
		list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		] |\
	primitive = None
	with open(request_template_path, encoding="utf-8", mode="r") as prereq_file:
		data = json.load(prereq_file)
	if not isinstance(data, dict):
		raise TypeError(f"request template data must be an object, not '{type(data)}'")
	request_template = data[request.param]
	if not isinstance(request_template, list):
		raise TypeError(f"Request template data must be a list, not '{type(request_template)}'")

	return request_template

@pytest.fixture()
def response_template_data(pytestconfig: pytest.Config
			   ) -> dict[str, primitive | list[primitive |
					  dict[str, object] | list[object]] | dict[object, object]]:
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
	response_template: dict[
		str,
		list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		] |\
	primitive = None
	with open(prereq_path, encoding="utf-8", mode="r") as prereq_file:
		response_template = json.load(prereq_file)
	if not isinstance(response_template, dict):
		raise TypeError(f"Response template data must be an object, not '{type(response_template)}'")

	return response_template


@pytest.fixture()
def cdn_post_data(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cdns endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get CDN request template data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	try:
		cdn = request_template_data[0]
	except IndexError as e:
		raise TypeError("malformed prerequisite data; no CDNs present in 'cdns' array property") from e

	if not isinstance(cdn, dict):
		raise TypeError(f"malformed prerequisite data; CDNs must be objects, not '{type(cdn)}'")

	# Return new post data and post response from cdns POST request
	randstr = str(randint(0, 1000))
	try:
		name = cdn["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		cdn["name"] = name[:4] + randstr
		domain_name = cdn["domainName"]
		if not isinstance(domain_name, str):
			raise TypeError(f"domainName must be str, not '{type(domain_name)}")
		cdn["domainName"] = domain_name[:5] + randstr
	except KeyError as e:
		raise TypeError(f"missing CDN property '{e.args[0]}'") from e

	logger.info("New cdn data to hit POST method %s", request_template_data)
	# Hitting cdns POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_cdn(data=cdn)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; cdn is not an object")
		return resp_obj
	except IndexError:
		logger.error("No CDN response data from cdns POST request.")
		sys.exit(1)


@pytest.fixture()
def cachegroup_post_data(to_session: TOSession, request_template_data: list[JSONData]
			 ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for cachegroup endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get Cachegroup data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""
	try:
		cachegroup = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Cache group present in 'cachegroup' array property") from e

	if not isinstance(cachegroup, dict):
		raise TypeError(
			f"malformed prerequisite data; Cache group must be objects, not '{type(cachegroup)}'")

	# Return new post data and post response from cachegroups POST request
	randstr = str(randint(0, 1000))
	try:
		name = cachegroup["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		cachegroup["name"] = name[:4] + randstr
		short_name = cachegroup["shortName"]
		if not isinstance(short_name, str):
			raise TypeError(f"shortName must be str, not '{type(short_name)}")
		cachegroup["shortName"] = short_name[:5] + randstr
	except KeyError as e:
		raise TypeError(f"missing Cache group property '{e.args[0]}'") from e
	# Hitting types GET method to access typeID for cachegroup POST data
	type_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_types(query_params={"useInTable": "cachegroup"})
	try:
		type_data = type_get_response[0]
		if not isinstance(type_data, list):
			raise TypeError("malformed API response; 'response' property not an array")
		first_type = type_data[0]
		if not isinstance(first_type, dict):
			raise TypeError("malformed API response; first Type in response is not an object")
		cachegroup["typeId"] = first_type["id"]
		type_id = cachegroup["typeId"]
		logger.info("extracted %s from %s", type_id, type_get_response)
	except KeyError as e:
		raise TypeError(f"missing Type property '{e.args[0]}'") from e

	logger.info("New cachegroup data to hit POST method %s", request_template_data)
	# Hitting cachegroup POST method
	response: tuple[JSONData, requests.Response] = to_session.create_cachegroups(data=cachegroup)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; cache group is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Cache group response data from cdns POST request.")
		sys.exit(1)


@pytest.fixture()
def parameter_post_data(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for parameters endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get parameter data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	try:
		parameter = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Parameters present in 'parameters' array property") from e

	if not isinstance(parameter, dict):
		raise TypeError(
			f"malformed prerequisite data; Paremeters must be objects, not '{type(parameter)}'")

	# Return new post data and post response from parameters POST request
	randstr = str(randint(0, 1000))
	try:
		name = parameter["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		parameter["name"] = name[:4] + randstr
		value = parameter["value"]
		if not isinstance(value, str):
			raise TypeError(f"value must be str, not '{type(value)}")
		parameter["value"] = value[:5] + randstr
	except KeyError as e:
		raise TypeError(f"missing Parameter property '{e.args[0]}'") from e

	logger.info("New parameter data to hit POST method %s", request_template_data)
	# Hitting cdns POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_parameter(data=parameter)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; parameter is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Parameter response data from parameters POST request.")
		sys.exit(1)


@pytest.fixture()
def role_post_data(to_session: TOSession, request_template_data: list[JSONData]
			) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for roles endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get role data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	try:
		role = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Roles present in 'roles' array property") from e

	if not isinstance(role, dict):
		raise TypeError(
			f"malformed prerequisite data; Roles must be objects, not '{type(role)}'")

	# Return new post data and post response from roles POST request
	randstr = str(randint(0, 1000))
	try:
		name = role["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		role["name"] = name[:4] + randstr
		description = role["description"]
		if not isinstance(description, str):
			raise TypeError(f"description must be str, not '{type(description)}")
		role["description"] = description[:5] + randstr
	except KeyError as e:
		raise TypeError(f"missing Role property '{e.args[0]}'") from e

	logger.info("New role data to hit POST method %s", request_template_data)
	# Hitting roles POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_role(data=role)
	logger.info(response[0])
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; role is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Role response data from roles POST request.")
		sys.exit(1)


@pytest.fixture()
def profile_post_data(to_session: TOSession, request_template_data: list[JSONData]
			  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for profile endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get profile data from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	try:
		profile = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Profile present in 'profile' array property") from e

	if not isinstance(profile, dict):
		raise TypeError(f"malformed prerequisite data; profile must be objects, not '{type(profile)}'")

	# Return new post data and post response from cachegroups POST request
	randstr = str(randint(0, 1000))
	try:
		name = profile["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		profile["name"] = name[:4] + randstr
	except KeyError as e:
		raise TypeError(f"missing Profile property '{e.args[0]}'") from e
	# Hitting types GET method to access typeID for profile POST data
	cdn_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_cdns()
	try:
		cdn_data = cdn_get_response[0]
		if not isinstance(cdn_data, list):
			raise TypeError("malformed API response; 'response' property not an array")
		first_cdn = cdn_data[0]
		if not isinstance(first_cdn, dict):
			raise TypeError("malformed API response; first CDN in response is not an object")
		profile["cdn"] = first_cdn["id"]
		cdn_id = profile["cdn"]
		logger.info("extracted %s from %s", cdn_id, cdn_get_response)
	except KeyError as e:
		raise TypeError(f"missing CDN property '{e.args[0]}'") from e

	logger.info("New profile data to hit POST method %s", request_template_data)
	# Hitting profile POST method
	response: tuple[JSONData, requests.Response] = to_session.create_profile(data=profile)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; profile is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Profile response data from cdns POST request.")
		sys.exit(1)


@pytest.fixture()
def tenant_post_data(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for tenants endpoint.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get tenant request template from a prerequisites file.
	:returns: Sample POST data and the actual API response.
	"""

	try:
		tenant = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Parameters present in 'tenants' array property") from e

	if not isinstance(tenant, dict):
		raise TypeError(
			f"malformed prerequisite data; tenants must be objects, not '{type(tenant)}'")

	# Return new post data and post response from tenants POST request
	randstr = str(randint(0, 1000))
	try:
		name = tenant["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		tenant["name"] = name[:4] + randstr
	except KeyError as e:
		raise TypeError(f"missing tenant property '{e.args[0]}'") from e

	logger.info("New tenant data to hit POST method %s", request_template_data)
	# Hitting tenants POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_tenant(data=tenant)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; parameter is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Parameter response data from parameters POST request.")
		sys.exit(1)

@pytest.fixture()
def server_capabilities_post_data(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for server_capabilities endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get server_capabilities data from a prerequisites file.
	  :returns: Sample POST data and the actual API response.
	"""
	try:
		server_capabilities = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no data present in 'server_capabilities' array property") from e

	if not isinstance(server_capabilities, dict):
		raise TypeError(
			f"malformed prerequisite data; data must be objects, not '{type(server_capabilities)}'")

	# Return new post data and post response from server_capabilities POST request
	randstr = str(randint(0, 1000))
	try:
		name = server_capabilities["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		server_capabilities["name"] = name[:3] + randstr
	except KeyError as e:
		raise TypeError(f"missing server_capabilities property '{e.args[0]}'") from e

	logger.info("New server_capabilities data to hit POST method %s", request_template_data)
	# Hitting server_capabilities POST method
	response: tuple[
		JSONData, requests.Response] = to_session.create_server_capabilities(data=server_capabilities)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; server_capabilities is not an object")
		return resp_obj
	except IndexError:
		logger.error("No server_capabilities response data from server_capabilities POST request.")
		sys.exit(1)

@pytest.fixture()
def division_post_data(to_session: TOSession, request_template_data: list[JSONData]
		  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for divisions endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get divisions request template data from
	request_template file.
	  :returns: Sample POST data and the actual API response.
	"""

	try:
		division = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no division present in 'division' array property") from e

	if not isinstance(division, dict):
		raise TypeError(
			f"malformed prerequisite data; divisions must be objects, not '{type(division)}'")

	# Return new post data and post response from division POST request
	randstr = str(randint(0, 1000))
	try:
		name = division["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		division["name"] = name[:4] + randstr
	except KeyError as e:
		raise TypeError(f"missing Parameter property '{e.args[0]}'") from e

	logger.info("New division data to hit POST method %s", request_template_data)
	# Hitting division POST methed
	response: tuple[JSONData, requests.Response] = to_session.create_division(data=division)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; division is not an object")
		return resp_obj
	except IndexError:
		logger.error("No division response data from division POST request.")
		sys.exit(1)

@pytest.fixture()
def region_post_data(to_session: TOSession, request_template_data: list[JSONData]
			  ) -> dict[str, object]:
	"""
	PyTest Fixture to create POST data for region endpoint.

	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get region data from a prerequisites file.
  	:returns: Sample POST data and the actual API response.
	"""

	try:
		region = request_template_data[0]
	except IndexError as e:
		raise TypeError(
			"malformed prerequisite data; no Region present in 'regions' array property") from e

	if not isinstance(region, dict):
		raise TypeError(f"malformed prerequisite data; region must be objects, not '{type(region)}'")

	# Return new post data and post response from regions POST request
	randstr = str(randint(0, 1000))
	try:
		name = region["name"]
		if not isinstance(name, str):
			raise TypeError(f"name must be str, not '{type(name)}'")
		region["name"] = name[:4] + randstr
	except KeyError as e:
		raise TypeError(f"missing Region property '{e.args[0]}'") from e
	# Hitting types GET method to access typeID for region POST data
	division_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_divisions()
	try:
		division_data = division_get_response[0]
		if not isinstance(division_data, list):
			raise TypeError("malformed API response; 'response' property not an array")
		first_division = division_data[0]
		if not isinstance(first_division, dict):
			raise TypeError("malformed API response; first Division in response is not an object")
		region["division"] = first_division["id"]
		division_id = region["division"]
		logger.info("extracted %s from %s", division_id, division_get_response)
	except KeyError as e:
		raise TypeError(f"missing Regions property '{e.args[0]}'") from e

	logger.info("New region data to hit POST method %s", request_template_data)
	# Hitting region POST method
	response: tuple[JSONData, requests.Response] = to_session.create_region(data=region)
	try:
		resp_obj = response[0]
		if not isinstance(resp_obj, dict):
			raise TypeError("malformed API response; region is not an object")
		return resp_obj
	except IndexError:
		logger.error("No Region response data from divisions POST request.")
	sys.exit(1)
