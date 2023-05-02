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
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["tenants"], indirect=True)
def test_tenant_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	tenant_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from tenants endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get request template data from a prerequisites file.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param tenant_post_data: Fixture to get sample tenant data and actual tenant response.
	"""
	# validate tenant keys from tenants get response
	logger.info("Accessing /tenants endpoint through Traffic ops session.")

	tenant = request_template_data[0]
	if not isinstance(tenant, dict):
		raise TypeError("malformed tenant in prerequisite data; not an object")

	tenant_name = tenant.get("name")
	if not isinstance(tenant_name, str):
		raise TypeError("malformed tenant in prerequisite data; 'name' not a string")

	tenant_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_tenants(query_params={"name": tenant_name})
	try:
		tenant_data = tenant_get_response[0]
		if not isinstance(tenant_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_tenant = tenant_data[0]
		if not isinstance(first_tenant, dict):
			raise TypeError("malformed API response; first tenant in response is not an object")
		tenant_keys = set(first_tenant.keys())

		logger.info("tenant Keys from tenants endpoint response %s", tenant_keys)
		tenant_response_template = response_template_data.get("tenants")
		if not isinstance(tenant_response_template, dict):
			raise TypeError(
				f"tenant response template data must be a dict, not '{type(tenant_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = tenant_response_template.get("properties")
		# validate tenant values from prereq data in tenants get response.
		prereq_values = [tenant_post_data["name"], tenant_post_data["active"],
		tenant_post_data["parentId"]]
		get_values = [first_tenant["name"], first_tenant["active"], first_tenant["parentId"]]
		get_types = {}
		for key, value in first_tenant.items():
			get_types[key] = type(value).__name__
		logger.info("types from tenant get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from tenants response template %s", response_template_types)
		# validate keys, data types and values from tenants get json response.
		assert tenant_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for tenants endpoint: API response was malformed")
	finally:
		# Delete tenant after test execution to avoid redundancy.
		try:
			tenant_id = tenant_post_data["id"]
			to_session.delete_tenant(tenant_id=tenant_id)
		except IndexError:
			logger.error("Tenant returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_tenant_contract")
