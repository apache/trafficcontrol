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

"""API Contract Test Case for regions endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_region_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
	dict[str, object], list[object]]], dict[object, object]]],
	region_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from regions response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param region_post_data: Fixture to get sample region data and actual region response.
	"""
	# validate region keys from regions get response
	logger.info("Accessing /regions endpoint through Traffic ops session.")

	region_name = region_post_data.get("name")
	if not isinstance(region_name, str):
		raise TypeError("malformed region in prerequisite data; 'name' not a string")

	region_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_regions(query_params={"name": region_name})
	try:
		region_data = region_get_response[0]
		if not isinstance(region_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_region = region_data[0]
		if not isinstance(first_region, dict):
			raise TypeError("malformed API response; first region in response is not an dict")

		logger.info("Region Api response %s", first_region)
		region_response_template = response_template_data.get("regions")
		if not isinstance(region_response_template, dict):
			raise TypeError(
				f"region response template data must be a dict, not '{type(region_response_template)}'")

		# validate region values from prereq data in regions get response.
		keys = ["name", "division", "divisionName"]
		prereq_values = [region_post_data[key] for key in keys]
		get_values = [first_region[key] for key in keys]

		# validate keys, data types and values from regions get json response.
		assert validate(instance=first_region, schema=region_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for regions endpoint: API response was malformed")
