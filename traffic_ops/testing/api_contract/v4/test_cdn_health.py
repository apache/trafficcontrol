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

"""API Contract Test Case for cdn health, capacity and routing endpoints."""
import logging
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_cdn_health_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
			dict[str, object], list[object]]], dict[object, object]]]
) -> None:
	"""
	Test step to validate keys, values and data types from cdn health and routing endpoints
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	"""
	# validate cdn_health keys from api get response
	logger.info("Accessing /cdn_health endpoints through Traffic ops session.")

	cdn_health_get_response : tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_cdns_health()

	cdn_routing_get_response : tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_cdns_routing()

	try:
		cdn_health_data = cdn_health_get_response[0]
		if not isinstance(cdn_health_data, dict):
			raise TypeError("malformed API response; 'response' property not an dict")
		logger.info("cdn_health Api response %s", cdn_health_data)
		
		cdn_routing_data = cdn_routing_get_response[0]
		if not isinstance(cdn_routing_data, dict):
			raise TypeError("malformed API response; 'response' property not an dict")
		logger.info("cdn_routing Api response %s", cdn_routing_data)

		cdn_health_response_template = response_template_data.get(
			"cdn_health")
		if not isinstance(cdn_health_response_template, dict):
			raise TypeError(f"cdn_health response template data must be a dict, not '"
							f"{type(cdn_health_response_template)}'")

		cdn_routing_response_template = response_template_data.get(
			"cdn_routing")
		if not isinstance(cdn_routing_response_template, dict):
			raise TypeError(f"cdn_routing response template data must be a dict, not '"
							f"{type(cdn_routing_response_template)}'")

		assert validate(instance=cdn_health_data, schema=cdn_health_response_template) is None
		assert validate(instance=cdn_routing_data, schema=cdn_routing_response_template) is None
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cdn_health endpoints:"
	        "API response was malformed")
