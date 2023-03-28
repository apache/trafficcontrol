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

"""API Contract Test Case for parameters endpoint."""
import json
import logging
import os
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('api_prerequisite_data', ["parameters"], indirect=True)
def test_parameter_contract(
	to_session: TOSession,
	api_prerequisite_data: list[dict[str, object] | list[object] | primitive],
	parameter_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from parameters endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param parameter_prereq_data: Fixture to get CDN data from a prerequisites file.
	:param parameter_post_data: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from parameters get response
	logger.info("Accessing /parameters endpoint through Traffic ops session.")

	parameter = api_prerequisite_data[0]
	if not isinstance(parameter, dict):
		raise TypeError("malformed parameter in prerequisite data; not an object")

	parameter_name = parameter.get("name")
	if not isinstance(parameter_name, str):
		raise TypeError("malformed parameter in prerequisite data; 'name' not a string")

	parameter_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_parameters(query_params={"name": parameter_name})
	response_template_path = os.path.join(os.path.dirname(__file__), "response_template.json")
	with open(response_template_path, encoding="utf-8", mode="r") as response_template_file:
		response_template = json.load(response_template_file)
	if not isinstance(response_template, dict):
		raise TypeError(f"response template data must be an object, not '{type(response_template)}'")
	try:
		parameter_data = parameter_get_response[0]
		if not isinstance(parameter_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_parameter = parameter_data[0]
		if not isinstance(first_parameter, dict):
			raise TypeError("malformed API response; first Parameter in response is not an object")
		parameter_keys = set(first_parameter.keys())

		logger.info("Parameter Keys from parameters endpoint response %s", parameter_keys)
		response_template_data = response_template.get("parameters").get("properties")
		# validate parameter values from prereq data in parameters get response.
		prereq_values = [
			parameter_post_data["name"],
			parameter_post_data["value"],
			parameter_post_data["configFile"],
			parameter_post_data["secure"]
		]
		get_values = [
			first_parameter["name"],
	        first_parameter["value"],
	        first_parameter["configFile"],
	        first_parameter["secure"]
	    ]
		get_types = {}
		for key in first_parameter:
			get_types[key] = first_parameter[key].__class__.__name__
		logger.info("types from parameter get response %s", get_types)
		response_template_types= {}
		for key in response_template_data:
			response_template_types[key] = response_template_data.get(key).get("type")
		logger.info("types from cdn response template %s", response_template_types)
		# validate data types for values from parameters get json response.
		assert parameter_keys == set(parameter_post_data.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Either prerequisite data or API response was malformed")
	finally:
		# Delete Parameter after test execution to avoid redundancy.
		try:
			parameter_id = parameter_post_data["id"]
			to_session.delete_parameter(parameter_id=parameter_id)
		except IndexError:
			logger.error("Parameter returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_parameter_contract")
