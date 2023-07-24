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

"""API Contract Test Case for delivery_service_required_capabilities endpoint."""
import logging
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_delivery_service_required_capabilities_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	delivery_service_required_capabilities_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from delivery_service_required_capabilities endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param delivery_service_required_capabilities_post_data: Fixture to get sample data and response.
	"""
	# validate delivery_service_required_capabilities keys from api get response
	logger.info("Accessing /delivery_service_required_capabilities endpoint through Traffic ops session.")

	delivery_service_id = delivery_service_required_capabilities_post_data["deliveryServiceID"]
	if not isinstance(delivery_service_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")

	delivery_service_required_capabilities_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_deliveryservices_required_capabilities(
		query_params={"deliveryServiceID": delivery_service_id})
	try:
		delivery_service_required_capabilities_data = delivery_service_required_capabilities_get_response[0]
		if not isinstance(delivery_service_required_capabilities_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_delivery_service_required_capabilities = delivery_service_required_capabilities_data[0]
		if not isinstance(first_delivery_service_required_capabilities, dict):
			raise TypeError(
				"malformed API response; first delivery_service_required_capabilities in response is not an dict")
		logger.info(
			"delivery_service_required_capabilities Api response %s", first_delivery_service_required_capabilities)

		delivery_service_required_capabilities_response_template = response_template_data.get(
			"delivery_service_required_capabilities")
		if not isinstance(delivery_service_required_capabilities_response_template, dict):
			raise TypeError(f"delivery_service_required_capabilities response template data must be a dict, not '"
							f"{type(delivery_service_required_capabilities_response_template)}'")

		keys = ["deliveryServiceID", "requiredCapability"]
		prereq_values = [delivery_service_required_capabilities_post_data[key] for key in keys]
		get_values = [first_delivery_service_required_capabilities[key] for key in keys]

		assert validate(instance=first_delivery_service_required_capabilities,
		  schema=delivery_service_required_capabilities_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for delivery_service_required_capabilities endpoint:"
	        "API response was malformed")
