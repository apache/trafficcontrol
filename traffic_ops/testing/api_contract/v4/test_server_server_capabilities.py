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

"""API Contract Test Case for Server Server Capabilities endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_server_server_capabilities_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive, dict[str, object],
	list[object]]], dict[object, object]]], server_server_capabilities_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from server server capabilities endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param server_server_capabilities_post_data: Fixture to get sample server server capabilities data 
	and actual Server Server Capabilities response.
	"""
	# validate Server Server Capabilities keys from server server capabilities get response
	logger.info(
		"Accessing /server server capabilities endpoint through Traffic ops session.")

	serverCapability = server_server_capabilities_post_data.get("serverCapability")
	if not isinstance(serverCapability, str):
		raise TypeError(
			"malformed server server capabilities in prerequisite data; 'serverId' not an integer")

	server_server_capabilities_get_response: tuple[Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive], requests.Response
	] = to_session.get_server_server_capabilities(query_params={"serverCapability": serverCapability})
	try:
		server_server_capabilities_data = server_server_capabilities_get_response[0]
		if not isinstance(server_server_capabilities_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_server_server_capability = server_server_capabilities_data[0]
		if not isinstance(first_server_server_capability, dict):
			raise TypeError(
				"malformed API response; first Server Server Capability in response is not an object")
		logger.info("Server Server Capabilities Api get response %s", first_server_server_capability)
		server_server_capabilities_response_template = response_template_data.get("server_server_capabilities")
		if not isinstance(server_server_capabilities_response_template, dict):
			raise TypeError(
				f"Server Server Capability response template data must be a dict, not '{type(server_server_capabilities_response_template)}'")

		# validate server server capabilities values from prereq data in server server capabilities get response.
		keys = ["serverId","serverCapability"]
		prereq_values = [server_server_capabilities_post_data[key] for key in keys]
		get_values = [first_server_server_capability[key] for key in keys]

		assert validate(instance=first_server_server_capability,
						schema=server_server_capabilities_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for server_server_capabilities endpoint: API response was malformed")
