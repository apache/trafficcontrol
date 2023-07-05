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

"""API Contract Test Case for deliveryservice_request endpoint."""
import json
import logging
from typing import Union
import pytest
import jsonref
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_deliveryservice_request_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	deliveryservice_request_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from deliveryservice_request endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param deliveryservice_request_post_data: Fixture to get sample delivery service data and response.
	"""
	# validate deliveryservice_request keys from api get response
	logger.info("Accessing /deliveryservice_request endpoint through Traffic ops session.")

	deliveryservice_request_id = deliveryservice_request_post_data["id"]
	if not isinstance(deliveryservice_request_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")

	deliveryservice_request_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_deliveryservice_requests(query_params={"id": deliveryservice_request_id})
	try:
		deliveryservice_request_data = deliveryservice_request_get_response[0]
		if not isinstance(deliveryservice_request_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_deliveryservice_request = deliveryservice_request_data[0]
		if not isinstance(first_deliveryservice_request, dict):
			raise TypeError("malformed API response; first deliveryservice_request in response is not an dict")
		logger.info("deliveryservice_request Api response %s", first_deliveryservice_request)

		# Resolve $ref references in the schema
		scheme_str = json.dumps(response_template_data)
		response_data = jsonref.loads(scheme_str)
		deliveryservice_request_response_template = response_data.get("deliveryservice_requests")
		if not isinstance(deliveryservice_request_response_template, dict):
			raise TypeError(f"deliveryservice_request response template data must be a dict, not '"
							f"{type(deliveryservice_request_response_template)}'")

		keys = ["displayName", "xmlId", "id", "cdnId", "tenantId", "type", "typeId"]
		prereq_values = [deliveryservice_request_post_data["requested"][key] for key in keys]
		get_values = [first_deliveryservice_request["requested"][key] for key in keys]

		assert validate(instance=first_deliveryservice_request,
		  schema=deliveryservice_request_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for deliveryservice_request endpoint: API response was malformed")
