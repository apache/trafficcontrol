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

"""API Contract Test Case for steering endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

def test_steering_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive,
					 list[Union[Primitive, dict[str, object], list[object]]],
	dict[object, object]]], steering_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from steering endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param steering_post_data: Fixture to get sample steering data and actual steering response.
	"""
	# validate steering keys from steering get response
	logger.info("Accessing /steering endpoint through Traffic ops session.")

	target_id = steering_post_data.get("targetId")
	if not isinstance(target_id, int):
		raise TypeError("malformed target in prerequisite data; 'target_id' not a integer")
	deliveryservice_id = steering_post_data.get("deliveryServiceId")
	if not isinstance(deliveryservice_id, int):
		raise TypeError("malformed deliveryservice_id in prerequisite data; 'deliveryservice_id ' not a integer")

	steering_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_steering_targets(delivery_service_id = deliveryservice_id, query_params={"target":target_id})
	try:
		steering_data = steering_get_response[0]
		if not isinstance(steering_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_steering = steering_data[0]
		if not isinstance(first_steering, dict):
			raise TypeError("malformed API response; first steering in response is not an dict")
		logger.info("steering Api get response %s", first_steering)

		steering_response_template = response_template_data.get("steering")
		if not isinstance(steering_response_template, dict):
			raise TypeError(
				f"steering response template data must be a dict, not '{type(steering_response_template)}'")

		# validate steering values from prereq data in steering get response.
		keys = ["targetId", "value", "typeId"]
		prereq_values = [steering_post_data[key] for key in keys]
		get_values = [first_steering[key] for key in keys]

		assert validate(instance=first_steering, schema=steering_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for steering endpoint: API response was malformed")
