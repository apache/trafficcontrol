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


def test_phys_locations_contract(to_session: TOSession,
	response_template_data: list[Union[dict[str, object], list[object], Primitive]],
	phys_locations_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from phys_locations response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param phys_location_post_data: Fixture to get sample phys_location data and actual 
	phys_location response.
	"""
	# validate phys_location keys from phys_locations get response
	logger.info("Accessing /phys_locations endpoint through Traffic ops session.")

	phys_location_name = phys_locations_post_data.get("name")
	if not isinstance(phys_location_name, str):
		raise TypeError("malformed phys_location in prerequisite data; 'name' not a string")

	phys_location_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_physical_locations(query_params={"name": phys_location_name})
	try:
		phys_location_data = phys_location_get_response[0]
		if not isinstance(phys_location_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_phys_location = phys_location_data[0]
		if not isinstance(first_phys_location, dict):
			raise TypeError("malformed API response; first phys_location in response is not an dict")

		logger.info("phys_location Api response %s", first_phys_location)
		phys_location_response_template = response_template_data.get("phys_locations")
		if not isinstance(phys_location_response_template, dict):
			raise TypeError(
				f"phys_loc response template data must be a dict, not'{type(phys_location_response_template)}'")

		# validate phys_location values from prereq data in phys_locations get response.
		keys = ["name", "address", "city", "zip", "comments", "email", "phone","poc",
	            "regionId", "shortName", "state"]
		prereq_values = [phys_locations_post_data[key] for key in keys]
		get_values = [first_phys_location[key] for key in keys]

		#validate keys, data types and values from regions get json response.
		assert validate(instance=first_phys_location, schema=phys_location_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for phys_locations endpoint: API response was malformed")
