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

"""API Contract Test Case for regions endpoint."""
import logging
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["regions"], indirect=True)
def test_region_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	region_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from regions endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get request template data from a prerequisites file.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param region_post_data: Fixture to get sample region data and actual region response.
	"""
	# validate region keys from regions get response
	logger.info("Accessing /regions endpoint through Traffic ops session.")

	region = request_template_data[0]
	if not isinstance(region, dict):
		raise TypeError("malformed region in prerequisite data; not an object")

	region_name = region.get("name")
	if not isinstance(region_name, str):
		raise TypeError("malformed region in prerequisite data; 'name' not a string")

	region_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_regions(query_params={"name": region_name})
	try:
		region_data = region_get_response[0]
		if not isinstance(region_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_region = region_data[0]
		if not isinstance(first_region, dict):
			raise TypeError("malformed API response; first region in response is not an object")
		region_keys = set(first_region.keys())

		logger.info("region Keys from regions endpoint response %s", region_keys)
		region_response_template = response_template_data.get("regions")
		if not isinstance(region_response_template, dict):
			raise TypeError(
				f"region response template data must be a dict, not '{type(region_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = region_response_template.get("properties")
		# validate region values from prereq data in regions get response.
		prereq_values = [region_post_data["name"], region_post_data["division"],
		region_post_data["divisionName"]]
		get_values = [first_region["name"], first_region["division"], first_region["divisionName"]]
		get_types = {}
		for key, value in first_region.items():
			get_types[key] = type(value).__name__
		logger.info("types from region get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from regions response template %s", response_template_types)
		# validate keys, data types and values from regions get json response.
		assert region_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for regions endpoint: API response was malformed")
	finally:
		# Delete region after test execution to avoid redundancy.
		try:
			region_name = region_post_data["name"]
			to_session.delete_region(query_params={"name": region_name})
		except IndexError:
			logger.error("region returned by Traffic Ops is missing a 'name' property")
			pytest.fail("Response from delete request is empty, Failing test_region_contract")
