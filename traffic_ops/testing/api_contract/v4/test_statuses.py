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

"""API Contract Test Case for statuses endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_status_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive, dict[str, object],
	list[object]]], dict[object, object]]], status_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from statuses endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param status_post_data: Fixture to get sample Status data and actual Status response.
	"""
	# validate Status keys from statuses get response
	logger.info("Accessing /statuses endpoint through Traffic ops session.")

	status_name = status_post_data.get("name")
	if not isinstance(status_name, str):
		raise TypeError("malformed status in prerequisite data; 'name' not a string")

	status_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_statuses(query_params={"name": status_name})
	try:
		status_data = status_get_response[0]
		if not isinstance(status_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_status = status_data[0]
		if not isinstance(first_status, dict):
			raise TypeError("malformed API response; first Status in response is not an object")
		logger.info("Status Api get response %s", first_status)
		status_response_template = response_template_data.get("statuses")
		if not isinstance(status_response_template, dict):
			raise TypeError(
				f"Status response template data must be a dict, not '{type(status_response_template)}'")

		# validate status values from prereq data in statuses get response.
		keys = ["name", "description"]
		prereq_values = [status_post_data[key] for key in keys]
		get_values = [first_status[key] for key in keys]

		assert validate(instance=first_status, schema=status_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for status endpoint: API response was malformed")
