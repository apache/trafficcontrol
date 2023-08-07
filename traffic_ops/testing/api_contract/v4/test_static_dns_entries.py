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

"""API Contract Test Case for staticdnsentries endpoint."""
import logging
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_static_dns_entries_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	static_dns_entries_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from static_dns_entries endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param static_dns_entries_post_data: Fixture to get delivery services sslkeys data.
	"""
	# validate static_dns_entries keys from api get response
	logger.info("Accessing /static_dns_entries endpoint through Traffic ops session.")

	static_dns_entries_id = static_dns_entries_post_data["id"]
	if not isinstance(static_dns_entries_id, int):
		raise TypeError("malformed API response; 'id' property not a string")

	static_dns_entries_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_staticdnsentries(query_params={"id":static_dns_entries_id})
	try:
		static_dns_entries_data = static_dns_entries_get_response[0]
		if not isinstance(static_dns_entries_data, list):
			raise TypeError(
				"malformed API response; static_dns_entries in response is not an list")
		first_static_dns_entries = static_dns_entries_data[0]
		if not isinstance(first_static_dns_entries, dict):
			raise TypeError(
				"malformed API response; first static_dns_entries in response is not a dict")

		logger.info("static_dns_entries Api get response %s", first_static_dns_entries)

		static_dns_entries_response_template = response_template_data.get(
			"static_dns_entries")
		if not isinstance(static_dns_entries_response_template, dict):
			raise TypeError(f"static_dns_entries response template data must be a dict, not '"
							f"{type(static_dns_entries_response_template)}'")

		keys = ["deliveryserviceId", "address", "ttl", "typeId", "host"]
		prereq_values = [static_dns_entries_post_data[key] for key in keys]
		get_values = [first_static_dns_entries[key] for key in keys]

		assert validate(instance=first_static_dns_entries,
		  schema=static_dns_entries_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail(
			"API contract test failed for static_dns_entries endpoint: API response was malformed")
