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
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_parameter_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]],
	dict[object, object]]], parameter_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from parameters endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param parameter_post_data: Fixture to get sample parameter data and actual parameter response.
	"""
	# validate Parameter keys from parameters get response
	logger.info("Accessing /parameters endpoint through Traffic ops session.")

	parameter_name = parameter_post_data.get("name")
	if not isinstance(parameter_name, str):
		raise TypeError("malformed parameter in prerequisite data; 'name' not a string")

	parameter_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_parameters(query_params={"name": parameter_name})
	try:
		parameter_data = parameter_get_response[0]
		if not isinstance(parameter_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_parameter = parameter_data[0]
		if not isinstance(first_parameter, dict):
			raise TypeError("malformed API response; first Parameter in response is not an dict")
		logger.info("Parameter Api get response %s", first_parameter)

		parameter_response_template = response_template_data.get("parameters")
		if not isinstance(parameter_response_template, dict):
			raise TypeError(
				f"Parameter response template data must be a dict, not '{type(parameter_response_template)}'")

		# validate parameter values from prereq data in parameters get response.
		keys = ["name", "value", "configFile", "secure"]
		prereq_values = [parameter_post_data[key] for key in keys]
		get_values = [first_parameter[key] for key in keys]

		# validate keys, data types and values from parameters get json response.
		assert validate(instance=first_parameter, schema=parameter_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cdn endpoint: API response was malformed")
