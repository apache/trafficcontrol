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

"""API Contract Test Case for tenants endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_tenant_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	tenant_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from tenants endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param tenant_post_data: Fixture to get sample tenant data and actual tenant response.
	"""
	# validate tenant keys from tenants get response
	logger.info("Accessing /tenants endpoint through Traffic ops session.")

	tenant_name = tenant_post_data.get("name")
	if not isinstance(tenant_name, str):
		raise TypeError("malformed tenant in prerequisite data; 'name' not a string")

	tenant_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_tenants(query_params={"name": tenant_name})
	try:
		tenant_data = tenant_get_response[0]
		if not isinstance(tenant_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_tenant = tenant_data[0]
		if not isinstance(first_tenant, dict):
			raise TypeError("malformed API response; first tenant in response is not an dict")
		logger.info("Tenant Api get response %s", first_tenant)

		tenant_response_template = response_template_data.get("tenants")
		if not isinstance(tenant_response_template, dict):
			raise TypeError(
				f"tenant response template data must be a dict, not '{type(tenant_response_template)}'")

		# validate tenant values from prereq data in tenants get response.
		keys = ["name", "active", "parentId"]
		prereq_values = [tenant_post_data[key] for key in keys]
		get_values = [first_tenant[key] for key in keys]

		# validate keys, data types and values from tenants get json response.
		assert validate(instance=first_tenant, schema=tenant_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for tenants endpoint: API response was malformed")
