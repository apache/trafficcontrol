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

"""API Contract Test Case for delivery services endpoint."""
import logging
from random import randint
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_delivery_services_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	delivery_services_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from delivery_services endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param delivery_services_post_data: Fixture to get sample delivery service data and response.
	"""
	# validate delivery services keys from api get response
	logger.info("Accessing /delivery_services endpoint through Traffic ops session.")

	delivery_services_id = delivery_services_post_data["id"]
	if not isinstance(delivery_services_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")

	delivery_services_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_deliveryservices(query_params={"id": delivery_services_id})

	# Hitting delivery_services PUT method
	delivery_services_post_data["longDesc"] = "test" + str(randint(0, 1000)) 
	logger.info("Updated delivery services data to hit PUT method %s", delivery_services_post_data)
	put_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive], 
		requests.Response] = to_session.update_deliveryservice_by_id(
		data=delivery_services_post_data, delivery_service_id=delivery_services_id)

	try:
		delivery_services_data = delivery_services_get_response[0]
		if not isinstance(delivery_services_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_delivery_services = delivery_services_data[0]
		if not isinstance(first_delivery_services, dict):
			raise TypeError("malformed API response; first delivery_services in response is not an dict")
		logger.info("delivery_services Api get response %s", first_delivery_services)

		delivery_services_put_data = put_response[0]
		if not isinstance(delivery_services_put_data, list):
			raise TypeError("malformed API response; 'response' property not an array")
		delivery_services_put_response = delivery_services_put_data[0]
		if not isinstance(delivery_services_put_response, dict):
			raise TypeError("malformed API response; delivery_services in response is not an dict")
		logger.info("delivery_services Api put response %s", delivery_services_put_response)
		delivery_services_response_template = response_template_data.get("delivery_services")
		if not isinstance(delivery_services_response_template, dict):
			raise TypeError(f"delivery_services response template data must be a dict, not '"
							f"{type(delivery_services_response_template)}'")

		keys = ["cdnId", "profileId", "type", "typeId", "tenantId", "xmlId"]
		prereq_values = [delivery_services_post_data[key] for key in keys]
		get_values = [first_delivery_services[key] for key in keys]

		assert validate(instance=first_delivery_services,
		  schema=delivery_services_response_template) is None
		assert get_values == prereq_values
		assert validate(instance=delivery_services_put_response,
		  schema=delivery_services_response_template) is None
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for delivery_services endpoint: API response was malformed")
