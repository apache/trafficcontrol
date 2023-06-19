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

"""API Contract Test Case for cdns endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_cdn_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive, dict[str, object],
	list[object]]], dict[object, object]]], cdn_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from cdns endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param cdn_post_data: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing /cdns endpoint through Traffic ops session.")

	cdn_name = cdn_post_data.get("name")
	if not isinstance(cdn_name, str):
		raise TypeError("malformed cdn in prerequisite data; 'name' not a string")

	cdn_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_cdns(query_params={"name": cdn_name})
	try:
		cdn_data = cdn_get_response[0]
		if not isinstance(cdn_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_cdn = cdn_data[0]
		if not isinstance(first_cdn, dict):
			raise TypeError("malformed API response; first CDN in response is not an object")
		logger.info("CDN Api get response %s", first_cdn)
		cdn_response_template = response_template_data.get("cdns")
		if not isinstance(cdn_response_template, dict):
			raise TypeError(
				f"Cdn response template data must be a dict, not '{type(cdn_response_template)}'")

		# validate cdn values from prereq data in cdns get response.
		keys = ["name", "domainName", "dnssecEnabled"]
		prereq_values = [cdn_post_data[key] for key in keys]
		get_values = [first_cdn[key] for key in keys]

		assert validate(instance=first_cdn, schema=cdn_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cdn endpoint: API response was malformed")
