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

"""API Contract Test Case for users endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

def test_user_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive,
					 list[Union[Primitive, dict[str, object], list[object]]],
	dict[object, object]]], user_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from users endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param user_post_data: Fixture to get sample user data and actual user response.
	"""
	# validate user keys from users get response
	logger.info("Accessing /users endpoint through Traffic ops session.")

	user_id = user_post_data.get("id")
	if not isinstance(user_id, int):
		raise TypeError("malformed user in prerequisite data; 'id' not a integer")

	user_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_user_by_id(user_id=user_id)
	try:
		user_data = user_get_response[0]
		if not isinstance(user_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_user = user_data[0]
		if not isinstance(first_user, dict):
			raise TypeError("malformed API response; first user in response is not an dict")
		logger.info("user Api get response %s", first_user)

		user_response_template = response_template_data.get("users")
		if not isinstance(user_response_template, dict):
			raise TypeError(
				f"user response template data must be a dict, not '{type(user_response_template)}'")

		# validate user values from prereq data in users get response.
		prereq_values = [user_post_data["username"], user_post_data["tenantId"]]
		get_values = [first_user["username"], first_user["tenantId"]]

		assert validate(instance=first_user, schema=user_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for users endpoint: API response was malformed")
