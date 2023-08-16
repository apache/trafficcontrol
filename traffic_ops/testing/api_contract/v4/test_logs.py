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

"""API Contract Test Case for logs endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_logs_contract(
	to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive, dict[str, object],
	list[object]]], dict[object, object]]], logs_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from logss endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param logs_post_data: Fixture to get sample Logs data and actual Logs response.
	"""
	# validate Log keys from logs get response
	logger.info("Accessing /logs endpoint through Traffic ops session.")

	change_log_id = logs_post_data.get("id")
	if not isinstance(change_log_id, int):
		raise TypeError("malformed log in prerequisite data; 'id' not an integer")

	change_logs_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_change_logs()
	try:
		change_logs_data = change_logs_get_response[0]
		if not isinstance(change_logs_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_change_log = change_logs_data[0]
		if not isinstance(first_change_log, dict):
			raise TypeError("malformed API response; first Log in response is not an object")
		logger.info("Logs Api get response %s", first_change_log)
		change_log_response_template = response_template_data.get("logs")
		if not isinstance(change_log_response_template, dict):
			raise TypeError(
				f"Log response template data must be a dict, not '{type(change_log_response_template)}'")

		# validate log values from prereq data in change logs get response.
		keys = ["id", "user"]
		prereq_values = [logs_post_data[key] for key in keys]
		get_values = [first_change_log[key] for key in keys]

		assert validate(instance=first_change_log, schema=change_log_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for logs endpoint: API response was malformed")
