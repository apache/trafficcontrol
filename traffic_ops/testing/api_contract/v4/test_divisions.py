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

"""API Contract Test Case for divisions endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_division_contract(to_session: TOSession,
	response_template_data: list[Union[dict[str, object], list[object], Primitive]],
	division_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from divisions endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param division_post_dat: Fixture to get sample division data and actual division response.
	"""
	# validate division keys from divisions get response
	logger.info("Accessing divisions endpoint through Traffic ops session.")

	division_name = division_post_data.get("name")
	if not isinstance(division_name, str):
		raise TypeError("malformed cdn in prerequisite data; 'name' not a string")

	division_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_divisions(query_params={"name": division_name})
	try:
		division_data = division_get_response[0]
		if not isinstance(division_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_division = division_data[0]
		if not isinstance(first_division, dict):
			raise TypeError("malformed API response; first division in response is not an dict")

		logger.info("Division Api get response %s", first_division)
		division_response_template = response_template_data.get("divisions")
		if not isinstance(division_response_template, dict):
			raise TypeError(
				f"Division response template data must be a dict, not '{type(division_response_template)}'")

		# validate division values from prereq data in divisions get response.
		prereq_values = division_post_data["name"]
		get_values = first_division["name"]

		assert validate(instance=first_division, schema=division_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Failed due to malformation")
