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

"""API Contract Test Case for federation_resolvers endpoint."""
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


def test_federation_resolvers_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	federation_resolver_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from federation_resolver endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param federation_resolver_post_data: Fixture to get delivery services sslkeys data.
	"""
	# validate federation_resolver keys from api get response
	logger.info("Accessing /federation_resolver endpoint through Traffic ops session.")

	federation_id = federation_resolver_post_data["id"]
	if not isinstance(federation_id, int):
		raise TypeError("malformed API response; 'federation_id' property not a integer")

	federation_resolver_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federation_resolvers(query_params={"id":federation_id})
    
	try:
		federation_resolver = federation_resolver_get_response[0]
		if not isinstance(federation_resolver, list):
			raise TypeError(
				"malformed API response; first federation_resolver in response is not an dict")
		first_federation_resolver = federation_resolver[0]

		federation_resolver_response_template = response_template_data.get(
			"federation_resolver")
		if not isinstance(federation_resolver_response_template, dict):
			raise TypeError(f"federation_resolver response template data must be a dict, not '"
							f"{type(federation_resolver_response_template)}'")

		assert validate(instance=first_federation_resolver, schema=federation_resolver_response_template) is None
		assert first_federation_resolver["ipAddress"] == federation_resolver_post_data["ipAddress"]
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail(
			"API contract test failed for federation_resolver endpoint: API response was malformed")
