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
import json
import logging
import os
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('api_prerequisite_data', ["roles"], indirect=True)
def test_role_contract(
	to_session: TOSession,
	api_prerequisite_data: list[dict[str, object] | list[object] | primitive],
	role_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from roles endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param api_prerequisite_data: Fixture to get role data from a prerequisites file.
	:param role_post_data: Fixture to get sample role data and actual role response.
	"""
	# validate Role keys from roles get response
	logger.info("Accessing /roles endpoint through Traffic ops session.")

	role = api_prerequisite_data[0]
	if not isinstance(role, dict):
		raise TypeError("malformed role in prerequisite data; not an object")

	role_name = role.get("name")
	if not isinstance(role_name, str):
		raise TypeError("malformed role in prerequisite data; 'name' not a string")

	role_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_roles(query_params={"name": role_name})
	response_template_path = os.path.join(os.path.dirname(__file__), "response_template.json")
	with open(response_template_path, encoding="utf-8", mode="r") as response_template_file:
		response_template = json.load(response_template_file)
	if not isinstance(response_template, dict):
		raise TypeError(f"response template data must be an object, not '{type(response_template)}'")
	try:
		role_data = role_get_response[0]
		if not isinstance(role_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_role = role_data[0]
		if not isinstance(first_role, dict):
			raise TypeError("malformed API response; first role in response is not an object")
		role_keys = set(first_role.keys())

		logger.info("Role Keys from roles endpoint response %s", role_keys)
		response_template_data = response_template.get("roles").get("properties")
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
		for key in first_role:
			get_types[key] = first_role[key].__class__.__name__
		logger.info("types from role get response %s", get_types)
		response_template_types= {}
		for key in response_template_data:
			response_template_types[key] = response_template_data.get(key).get("type")
		logger.info("types from role response template %s", response_template_types)
		# validate keys,data types for values and values from roles get json response.
		assert role_keys == set(role_post_data.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Either prerequisite data or API response was malformed")
	finally:
		# Delete Roel after test execution to avoid redundancy.
		try:
			role_name = role_post_data["name"]
			to_session.delete_role(query_params={"name": role_name})
		except IndexError:
			logger.error("Role returned by Traffic Ops is missing an 'name' property")
			pytest.fail("Response from delete request is empty, Failing test_role_contract")
