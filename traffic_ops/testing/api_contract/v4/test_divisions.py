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

"""API Contract Test Case for divisions endpoint."""
import logging
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["divisions"], indirect=True)
def test_division_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: list[dict[str, object] | list[object] | primitive],
	division_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from divisions endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param get_division_data: Fixture to get division data from a prerequisites file.
	:param division_prereq: Fixture to get sample division data and actual division response.
	"""
	# validate division keys from divisions get response
	logger.info("Accessing divisions endpoint through Traffic ops session.")

	division = request_template_data[0]
	if not isinstance(division, dict):
		raise TypeError("malformed division in prerequisite data; not an object")

	division_name = division.get("name")
	if not isinstance(division_name, str):
		raise TypeError("malformed division in prerequisite data; 'name' not a string")

	division_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_divisions(query_params={"name": division_name})
	try:
		division_data = division_get_response[0]
		if not isinstance(division_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_division = division_data[0]
		if not isinstance(first_division, dict):
			raise TypeError("malformed API response; first division in response is not an object")
		division_keys = set(first_division.keys())

		logger.info("division Keys from divisions endpoint response %s", division_keys)
		response_template = response_template_data.get("divisions").get("properties")
		# validate division values from prereq data in divisions get response.
		prereq_values = [
			division_post_data["name"]]
		get_values = [first_division["name"]]
		get_types = {}
		for key in first_division:
			get_types[key] = first_division[key].__class__.__name__
		logger.info("types from division get response %s", get_types)
		response_template_types= {}
		for key in response_template:
			response_template_types[key] = response_template.get(key).get("type")
		logger.info("types from division response template %s", response_template_types)

		assert division_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Failed due to malformation")
	finally:
		# Delete division after test execution to avoid redundancy.
		try:
			division_id = division_post_data["id"]
			to_session.delete_division(division_id=division_id)
		except IndexError:
			logger.error("Division returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_get_division")
