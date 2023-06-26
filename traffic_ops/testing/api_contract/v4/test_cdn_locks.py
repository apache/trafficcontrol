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

"""API Contract Test Case for cdn_locks endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

def test_cdn_lock_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive,
					 list[Union[Primitive, dict[str, object], list[object]]],
	dict[object, object]]], cdn_lock_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from cdn_locks endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param cdn_lock_post_data: Fixture to get sample cdn_lock data and actual cdn_lock response.
	"""
	# validate cdn_lock keys from cdn_locks get response
	logger.info("Accessing /cdn_locks endpoint through Traffic ops session.")

	cdn_name = cdn_lock_post_data.get("cdn")
	if not isinstance(cdn_name, str):
		raise TypeError("malformed cdn_name in prerequisite data; 'cdn_name' not a string")

	cdn_lock_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_cdn_locks(query_params={"name":cdn_name})
	try:
		cdn_lock_data = cdn_lock_get_response[0]
		if not isinstance(cdn_lock_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_cdn_lock = cdn_lock_data[0]
		if not isinstance(first_cdn_lock, dict):
			raise TypeError("malformed API response; first cdn_lock in response is not an dict")
		logger.info("cdn_lock Api get response %s", first_cdn_lock)

		cdn_lock_response_template = response_template_data.get("cdn_locks")
		if not isinstance(cdn_lock_response_template, dict):
			raise TypeError(
				f"cdn_lock response template data must be a dict, not '{type(cdn_lock_response_template)}'")

		# validate cdn_lock values from prereq data in cdn_locks get response.
		keys = ["cdn", "message", "soft"]
		prereq_values = [cdn_lock_post_data[key] for key in keys]
		get_values = [first_cdn_lock[key] for key in keys]

		assert validate(instance=first_cdn_lock, schema=cdn_lock_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cdn_locks endpoint: API response was malformed")
