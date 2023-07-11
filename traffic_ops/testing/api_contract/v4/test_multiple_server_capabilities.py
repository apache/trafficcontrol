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

"""API Contract Test Case for multiple_server_capabilities endpoint."""
import logging
from typing import Union

import pytest
from jsonschema import validate

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_multiple_server_capabilities_contract(response_template_data: dict[str, Union[Primitive,
	list[Union[Primitive, dict[str, object], list[object]]], dict[object, object]]],
	multiple_server_capabilities_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from multiple_server_capabilities endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param multiple_server_capabilities_post_data: Fixture to get sample data.
    and actual multiple_server_capabilities response.
	"""
	# validate multiple_server_capabilities keys from multiple_server_capabilities post response
	logger.info("Accessing /multiple_server_capabilities endpoint data.")

	try:
		if not isinstance(multiple_server_capabilities_post_data, dict):
			raise TypeError("malformed API response; multiple_server_capabilities is not an dict")
		logger.info("Multiple server capabilities Api response %s",
	      multiple_server_capabilities_post_data)
		response_template = response_template_data.get("multiple_servers_capabilities")

		# validate keys and data types from multiple_server_capabilities get json response.
		assert validate(instance=multiple_server_capabilities_post_data, schema=response_template) is None
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed : API response was malformed")
