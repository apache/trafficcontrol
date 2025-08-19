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

"""API Contract Test Case for phys_locations endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_types_contract(to_session: TOSession,
	response_template_data: list[Union[dict[str, object], list[object], Primitive]],
	types_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from types response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param types_post_data: Fixture to get sample types data and actual 
	types response.
	"""
	# validate types keys from types get response
	logger.info("Accessing /types endpoint through Traffic ops session.")

	type_name = types_post_data.get("name")
	if not isinstance(type_name, str):
		raise TypeError("malformed type in prerequisite data; 'name' not a string")

	types_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_types(query_params={"name": type_name})
	try:
		types_data = types_get_response[0]
		if not isinstance(types_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_type = types_data[0]
		if not isinstance(first_type, dict):
			raise TypeError("malformed API response; first type in response is not an dict")

		logger.info("types Api response %s", first_type)
		types_response_template = response_template_data.get("types")
		if not isinstance(types_response_template, dict):
			raise TypeError(
				f"Types response template data must be a dict, not'{type(types_response_template)}'")

		# validate types values from prereq data in types get response.
		keys = ["name", "description"]
		prereq_values = [types_post_data[key] for key in keys]
		get_values = [first_type[key] for key in keys]

		#validate keys, data types and values from regions get json response.
		assert validate(instance=first_type, schema=types_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for types endpoint: API response was malformed")
