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

"""API Contract Test Case for cdns endpoint."""
import logging
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["cdns"], indirect=True)
def test_cdn_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	cdn_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from cdns endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get CDN request template data from a prerequisites file.
	:param response_template_data: Fixture to get CDN request template data from a prerequisites file.
	:param cdn_post_data: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing /cdns endpoint through Traffic ops session.")

	cdn = request_template_data[0]
	if not isinstance(cdn, dict):
		raise TypeError("malformed cdn in prerequisite data; not an object")

	cdn_name = cdn.get("name")
	if not isinstance(cdn_name, str):
		raise TypeError("malformed cdn in prerequisite data; 'name' not a string")

	cdn_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_cdns(query_params={"name": cdn_name})
	try:
		cdn_data = cdn_get_response[0]
		if not isinstance(cdn_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_cdn = cdn_data[0]
		if not isinstance(first_cdn, dict):
			raise TypeError("malformed API response; first CDN in response is not an object")
		cdn_keys = set(first_cdn.keys())
		logger.info("CDN Keys from cdns endpoint response %s", cdn_keys)

		cdn_response_template = response_template_data.get("cdns")
		if not isinstance(cdn_response_template, dict):
			raise TypeError(
				f"Cdn response template data must be a dict, not '{type(cdn_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = cdn_response_template.get("properties")
		if not isinstance(response_template, dict):
			raise TypeError(
				f"response template data must be a dict, not '{type(response_template)}'")
		# validate cdn values from prereq data in cdns get response.
		prereq_values = [cdn_post_data["name"], cdn_post_data["domainName"],
		   cdn_post_data["dnssecEnabled"]]
		get_values = [first_cdn["name"], first_cdn["domainName"], first_cdn["dnssecEnabled"]]
		get_types = {}
		for key, value in first_cdn.items():
			get_types[key] = type(value).__name__
		logger.info("types from cdn get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from cdn response template %s", response_template_types)

		assert cdn_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cdn endpoint: API response was malformed")
	finally:
		# Delete CDN after test execution to avoid redundancy.
		try:
			cdn_id = cdn_post_data["id"]
			to_session.delete_cdn_by_id(cdn_id=cdn_id)
		except IndexError:
			logger.error("CDN returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_get_cdn")
