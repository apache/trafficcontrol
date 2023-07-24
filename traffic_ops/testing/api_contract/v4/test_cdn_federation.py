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

"""API Contract Test Case for cdn_federation endpoint."""
import logging
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_cdn_federation_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	cdn_federation_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from cdn_federation endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param cdn_federation_post_data: Fixture to get delivery services sslkeys data.
	"""
	# validate cdn_federation keys from api get response
	logger.info("Accessing /cdn_federation endpoint through Traffic ops session.")

	cdn_name = cdn_federation_post_data[0]
	cdn_federation_data_post = cdn_federation_post_data[1]
	if not isinstance(cdn_name, str):
		raise TypeError("malformed API response; 'cdn_name' property not a string")
	logger.info(cdn_name)

	cdn_federation_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federations_for_cdn(cdn_name=cdn_name)
	try:
		cdn_federation = cdn_federation_get_response[0]
		logger.info(cdn_federation)
		if not isinstance(cdn_federation, list):
			raise TypeError(
				"malformed API response; first cdn_federation in response is not an dict")
		first_cdn_federation = cdn_federation[1]
		logger.info(first_cdn_federation)
		logger.info("cdn_federation Api get response %s", first_cdn_federation)

		cdn_federation_response_template = response_template_data.get(
			"cdn_federation")
		if not isinstance(cdn_federation_response_template, dict):
			raise TypeError(f"cdn_federation response template data must be a dict, not '"
							f"{type(cdn_federation_response_template)}'")

		keys = ["cname", "ttl", "description"]
		prereq_values = [cdn_federation_data_post[key] for key in keys]
		get_values = [first_cdn_federation[key] for key in keys]

		assert validate(instance=first_cdn_federation,
		  schema=cdn_federation_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail(
			"API contract test failed for cdn_federation endpoint: API response was malformed")
