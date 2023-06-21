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

"""API Contract Test Case for coordinates endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_coordinate_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
	dict[str, object], list[object]]], dict[object, object]]],
	coordinate_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from coordinates response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param coordinate_post_data: Fixture to get sample coordinate data and actual coordinate response.
	"""
	# validate coordinate keys from coordinates get response
	logger.info("Accessing /coordinates endpoint through Traffic ops session.")

	coordinate_name = coordinate_post_data.get("name")
	if not isinstance(coordinate_name, str):
		raise TypeError("malformed coordinate in prerequisite data; 'name' not a string")

	coordinate_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_coordinates(query_params={"name": coordinate_name})
	try:
		coordinate_data = coordinate_get_response[0]
		if not isinstance(coordinate_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_coordinate = coordinate_data[0]
		if not isinstance(first_coordinate, dict):
			raise TypeError("malformed API response; first coordinate in response is not an dict")

		logger.info("Coordinate Api response %s", first_coordinate)
		coordinate_response_template = response_template_data.get("coordinates")
		if not isinstance(coordinate_response_template, dict):
			raise TypeError(
				f"coordinate response template data must be a dict, not '{type(coordinate_response_template)}'")

		# validate coordinate values from prereq data in coordinates get response.
		keys = ["name", "latitude", "longitude"]
		prereq_values = [coordinate_post_data[key] for key in keys]
		get_values = [first_coordinate[key] for key in keys]

		# validate keys, data types and values from coordinates get json response.
		assert validate(instance=first_coordinate, schema=coordinate_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for coordinates endpoint: API response was malformed")
