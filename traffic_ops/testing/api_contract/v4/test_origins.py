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

"""API Contract Test Case for origins endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_origin_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	origin_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from origins endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param origin_post_data: Fixture to get sample server data and actual origin response.
	"""
	# validate origin keys from origin get response
	logger.info("Accessing /origins endpoint through Traffic ops session.")

	origin_id = origin_post_data.get("id")
	if not isinstance(origin_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")

	origin_get_response: tuple[
		Union[dict[str, object], list[dict[str, object] | list[object] | Primitive], Primitive],
		requests.Response
	] = to_session.get_origins(query_params={"id": origin_id})
	try:
		origin_data = origin_get_response[0]
		if not isinstance(origin_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_origin = origin_data[0]
		if not isinstance(first_origin, dict):
			raise TypeError("malformed API response; first origin in response is not an dict")
		logger.info("origin Api get response %s", first_origin)
		origin_response_template = response_template_data.get("origins")
		if not isinstance(origin_response_template, dict):
			raise TypeError(
				f"origin response template data must be a dict, not '{type(origin_response_template)}'")

		keys = ["deliveryServiceId", "fqdn", "name", "port", "protocol", "tenantId"]
		prereq_values = [origin_post_data[key] for key in keys]
		get_values = [first_origin[key] for key in keys]

		assert validate(instance=first_origin, schema=origin_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for origin endpoint: API response was malformed")
