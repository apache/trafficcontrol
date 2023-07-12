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

"""API Contract Test Case for delivery_services_regex endpoint."""
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


def test_delivery_services_regex_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	delivery_services_regex_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from delivery_services_regex endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param delivery_services_regex_post_data: Fixture to get sample delivery service data and response.
	"""
	# validate delivery_services_regex_regex keys from api get response
	logger.info("Accessing /delivery_services_regex endpoint through Traffic ops session.")
	delivery_services_id = delivery_services_regex_post_data[0]
	delivery_services_regex_data_post = delivery_services_regex_post_data[1]

	delivery_services_regex_id = delivery_services_regex_data_post["id"]
	if not isinstance(delivery_services_regex_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")
	delivery_services_regex_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_deliveryservice_regexes_by_id(delivery_service_id=delivery_services_id,
						  query_params={"id": delivery_services_regex_id})
	delivery_services_regex_data_post["pattern"] = ".*\\.test" + str(randint(0, 1000)) + "\\..*"
	logger.info("Updated delivery services regex data to hit PUT method %s",
	     delivery_services_regex_data_post)
	# Hitting delivery_services_regex PUT method
	put_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response] = to_session.update_deliveryservice_regexes(
		delivery_service_id=delivery_services_id, regex_id=delivery_services_regex_id,
	    data= delivery_services_regex_data_post)
	try:
		delivery_services_regex_data = delivery_services_regex_get_response[0]
		if not isinstance(delivery_services_regex_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_delivery_services_regex = delivery_services_regex_data[0]
		if not isinstance(first_delivery_services_regex, dict):
			raise TypeError(
				"malformed API response; first delivery_services_regex in response is not an dict")
		logger.info("delivery_services_regex Api get response %s", first_delivery_services_regex)
		delivery_services_regex_put_response = put_response[0]
		if not isinstance(delivery_services_regex_put_response, dict):
			raise TypeError("malformed API response; delivery_services in response is not an dict")
		logger.info("delivery_services_regex Api put response %s", delivery_services_regex_put_response)
		delivery_services_regex_response_template = response_template_data.get("delivery_services_regex")
		if not isinstance(delivery_services_regex_response_template, dict):
			raise TypeError(f"delivery_services_regex response template data must be a dict, not '"
							f"{type(delivery_services_regex_response_template)}'")

		keys = ["type", "setNumber"]
		prereq_values = [delivery_services_regex_data_post[key] for key in keys]
		get_values = [first_delivery_services_regex[key] for key in keys]
		assert validate(instance=first_delivery_services_regex,
		  schema=delivery_services_regex_response_template) is None
		assert get_values == prereq_values
		assert validate(instance=delivery_services_regex_put_response,
		  schema=delivery_services_regex_response_template) is None
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed : API response was malformed")
