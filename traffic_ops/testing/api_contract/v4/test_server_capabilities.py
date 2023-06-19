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

"""API Contract Test Case for server_capabilities endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_server_capabilities_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	server_capabilities_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from server_capabilities endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param server_capabilities_post_data: Fixture to get sample server_capabilities data 
    and actual server_capabilities response.
	"""
	# validate server_capabilities keys from server_capabilities get response
	logger.info("Accessing /server_capabilities endpoint through Traffic ops session.")

	server_capabilities_name = server_capabilities_post_data.get("name")
	if not isinstance(server_capabilities_name, str):
		raise TypeError("malformed server_capabilities in prerequisite data; 'name' not a string")

	server_capabilities_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_server_capabilities(query_params={"name": server_capabilities_name})
	try:
		server_capabilities_data = server_capabilities_get_response[0]
		if not isinstance(server_capabilities_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_server_capabilities = server_capabilities_data[0]
		if not isinstance(first_server_capabilities, dict):
			raise TypeError("malformed API response; first server_capabilities in response is not an dict")
		logger.info("Server capabilities Api get response %s", first_server_capabilities)

		response_template = response_template_data.get("server_capabilities")

		# validate server_capabilities values from prereq data in api get response.
		prereq_values = server_capabilities_post_data["name"]
		get_values = first_server_capabilities["name"]

		# validate keys, data types and values from server_capabilities get json response.
		assert validate(instance=first_server_capabilities, schema=response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for server_capabilities: API response was malformed")
