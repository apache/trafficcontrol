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

"""API Contract Test Case for Profile Parameters endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_profile_parameters_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive, dict[str, object],
	list[object]]], dict[object, object]]], profile_parameters_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from profile parameters endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param profile_parameter_post_data: Fixture to get sample profile parameter data 
	and actual Profile Parameter response.
	"""
	# validate Profile Parameter keys from profile parameters get response
	logger.info(
		"Accessing /profile parameters endpoint through Traffic ops session.")

	profile_id = profile_parameters_post_data.get("profileId")
	if not isinstance(profile_id, int):
		raise TypeError(
			"malformed profile parameters in prerequisite data; 'profileId' not an integer")

	parameter_id = profile_parameters_post_data.get("parameterId")
	if not isinstance(parameter_id, int):
		raise TypeError(
			"malformed profile parameters in prerequisite data; 'parameterId' not an integer")

	profile_parameter_get_response: tuple[
		Union[dict[str, object],
			list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_profile_parameters(profile_id=profile_id, query_params={"profile": profile_id})
	try:
		profile_parameter_data = profile_parameter_get_response[0]
		if not isinstance(profile_parameter_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_profile_parameter = profile_parameter_data[0]
		if not isinstance(first_profile_parameter, dict):
			raise TypeError(
				"malformed API response; first Profile Parameter in response is not an object")
		logger.info("Profile Parameter Api get response %s", first_profile_parameter)
		profile_parameter_response_template = response_template_data.get("profile_parameters")
		if not isinstance(profile_parameter_response_template, dict):
			raise TypeError(
				f"Profile Parameter response template data must be a dict, not '{type(profile_parameter_response_template)}'")

		profile_parameters_post_data["profile"] = first_profile_parameter["profile"]
		profile_parameters_post_data["parameter"] = first_profile_parameter["parameter"]

		# validate profile_parameter values from prereq data in profile parameters get response.
		keys = ["profile", "parameter"]
		prereq_values = [profile_parameters_post_data[key] for key in keys]
		get_values = [first_profile_parameter[key] for key in keys]

		assert validate(instance=first_profile_parameter,
						schema=profile_parameter_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for profile_parameter endpoint: API response was malformed")
