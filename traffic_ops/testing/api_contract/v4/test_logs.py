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

"""
API Contract Test Case for logs endpoint.
"""

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
	list[object]]], dict[object, object]]], logs_data: list[Union[int, dict[str, object]]]) -> None:
	"""
	Test step to validate keys, values, and data types from logs endpoint response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from prerequisites file.
	"""
	# Validate Log keys from logs get response
	logger.info("Accessing /logs endpoint through Traffic ops session.")

	for log in logs_data:
		if isinstance(log, dict):
			change_log_id = log.get("id")
			if not isinstance(change_log_id, int):
				raise TypeError("Malformed log in prerequisite data; 'id' not an integer")

			# Hitting logs GET response
			change_logs_get_response: tuple[JSONData, requests.Response] = to_session.get_change_logs()

			logs_data, _ = change_logs_get_response

			try:
				logs_data = change_logs_get_response[0]
				if not isinstance(logs_data, list):
					raise TypeError("Malformed API response; 'response' property not an array")

				first_change_log = logs_data[0]
				if not isinstance(first_change_log, dict):
					raise TypeError("Malformed API response; first Log in response is not an object")
				logger.info("Logs API get response %s", first_change_log)

				change_log_response_template = response_template_data.get("logs")
				if not isinstance(change_log_response_template, dict):
					raise TypeError(
						f"Log response template data must be a dict, not '{type(change_log_response_template)}'")

				# Validate log values from prereq data in change logs get response.
				keys = ["id", "user"]
				prereq_values = [change_log_id, first_change_log["user"]]
				get_values = [first_change_log[key] for key in keys]

				assert validate(instance=first_change_log, schema=change_log_response_template) is None
				assert get_values == prereq_values
			except IndexError:
				logger.error("Either prerequisite data or API response was malformed")
				pytest.fail("API contract test failed for logs endpoint: API response was malformed")
		elif isinstance(log, int):
			# Handle the case where log_data is an integer
			pass
		else:
			raise TypeError(f"Malformed log in prerequisite data; expected dictionary or integer, got {type(log)}")
