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

"""API Contract Test Case for roles endpoint."""
import logging
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["roles"], indirect=True)
def test_role_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	role_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from roles endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get role data from a prerequisites file.
	:param role_post_data: Fixture to get sample role data and actual role response.
	"""
	# validate Role keys from roles get response
	logger.info("Accessing /roles endpoint through Traffic ops session.")

	role = request_template_data[0]
	if not isinstance(role, dict):
		raise TypeError("malformed role in prerequisite data; not an object")

	role_name = role.get("name")
	if not isinstance(role_name, str):
		raise TypeError("malformed role in prerequisite data; 'name' not a string")

	role_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_roles(query_params={"name": role_name})
	try:
		role_data = role_get_response[0]
		if not isinstance(role_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_role = role_data[0]
		if not isinstance(first_role, dict):
			raise TypeError("malformed API response; first role in response is not an object")
		role_keys = set(first_role.keys())

		logger.info("Role Keys from roles endpoint response %s", role_keys)
		role_response_template = response_template_data.get("roles")
		if not isinstance(role_response_template, dict):
			raise TypeError(
				f"Role response template data must be a dict, not '{type(role_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = role_response_template.get("properties")
		# validate roles values from prereq data in roles get response.
		prereq_values = [
			role_post_data["name"],
			role_post_data["description"]
		]
		get_values = [
			first_role["name"],
	        first_role["description"]
	    ]
		get_types = {}
		for key, value in first_role.items():
			get_types[key] = type(value).__name__
		logger.info("types from role get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from role response template %s", response_template_types)
		# validate keys,data types for values and values from roles get json response.
		assert role_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for roles endpoint: API response was malformed")
	finally:
		# Delete Roel after test execution to avoid redundancy.
		try:
			role_name = role_post_data["name"]
			to_session.delete_role(query_params={"name": role_name})
		except IndexError:
			logger.error("Role returned by Traffic Ops is missing an 'name' property")
			pytest.fail("Response from delete request is empty, Failing test_role_contract")
