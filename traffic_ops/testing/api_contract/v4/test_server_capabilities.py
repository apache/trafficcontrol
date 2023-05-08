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
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["server_capabilities"], indirect=True)
def test_server_capabilities_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	server_capabilities_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from server_capabilities endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get request template data from a prerequisites file.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param server_capabilities_post_data: Fixture to get sample server_capabilities data 
    and actual server_capabilities response.
	"""
	# validate server_capabilities keys from server_capabilities get response
	logger.info("Accessing /server_capabilities endpoint through Traffic ops session.")

	server_capabilities = request_template_data[0]
	if not isinstance(server_capabilities, dict):
		raise TypeError("malformed server_capabilities in prerequisite data; not an object")

	server_capabilities_name = server_capabilities.get("name")
	if not isinstance(server_capabilities_name, str):
		raise TypeError("malformed server_capabilities in prerequisite data; 'name' not a string")

	server_capabilities_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_server_capabilities(query_params={"name": server_capabilities_name})
	try:
		server_capabilities_data = server_capabilities_get_response[0]
		if not isinstance(server_capabilities_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_server_capabilities = server_capabilities_data[0]
		if not isinstance(first_server_capabilities, dict):
			raise TypeError("malformed API response; first server_capabilities in response is not an object")
		server_capabilities_keys = set(first_server_capabilities.keys())

		logger.info("server_capabilities Keys from endpoint response %s", server_capabilities_keys)
		server_capabilities_response_template = response_template_data.get("server_capabilities")
		if not isinstance(server_capabilities_response_template, dict):
			raise TypeError(
				f"server_capabilities data must be a dict, not '{type(server_capabilities_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = server_capabilities_response_template.get("properties")
		# validate server_capabilities values from prereq data in api get response.
		prereq_values = [server_capabilities_post_data["name"]]
		get_values = [first_server_capabilities["name"]]
		get_types = {}
		for key, value in first_server_capabilities.items():
			get_types[key] = type(value).__name__
		logger.info("types from server_capabilities get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from server_capabilities response template %s", response_template_types)
		# validate keys, data types and values from server_capabilities get json response.
		assert server_capabilities_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for server_capabilities: API response was malformed")
	finally:
		# Delete server_capabilities after test execution to avoid redundancy.
		try:
			server_capability_name = server_capabilities_post_data["name"]
			to_session.delete_server_capabilities(query_params={"name": server_capability_name})
		except IndexError:
			logger.error("server_capabilities returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_capabilities_contract")
