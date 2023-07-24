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

"""API Contract Test Case for delivery_service_sslkeys_generate endpoint."""
import logging
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_delivery_service_sslkeys_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	delivery_service_sslkeys_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from delivery_service_sslkeys endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param delivery_service_sslkeys_post_data: Fixture to get delivery services sslkeys data.
	"""
	# validate delivery_service_sslkeys keys from api get response
	logger.info("Accessing /delivery_service_sslkeys endpoint through Traffic ops session.")

	delivery_service_sslkeys_xml_id = delivery_service_sslkeys_post_data["key"]
	if not isinstance(delivery_service_sslkeys_xml_id, str):
		raise TypeError("malformed API response; 'xmlId' property not a string")

	delivery_service_sslkeys_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_deliveryservice_ssl_keys_by_xml_id(xml_id=delivery_service_sslkeys_xml_id)
	try:
		first_delivery_service_sslkeys = delivery_service_sslkeys_get_response[0]
		if not isinstance(first_delivery_service_sslkeys, dict):
			raise TypeError(
				"malformed API response; first delivery_service_sslkeys in response is not an dict")
		logger.info("delivery_service_sslkeys Api get response %s", first_delivery_service_sslkeys)

		delivery_service_sslkeys_response_template = response_template_data.get(
			"delivery_service_sslkeys")
		if not isinstance(delivery_service_sslkeys_response_template, dict):
			raise TypeError(f"delivery_service_sslkeys response template data must be a dict, not '"
							f"{type(delivery_service_sslkeys_response_template)}'")

		keys = ["key", "version"]
		prereq_values = [delivery_service_sslkeys_post_data[key] for key in keys]
		get_values = [first_delivery_service_sslkeys[key] for key in keys]

		assert validate(instance=first_delivery_service_sslkeys,
		  schema=delivery_service_sslkeys_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail(
			"API contract test failed for delivery_service_sslkeys endpoint: API response was malformed")
