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

"""API Contract Test Case for phys_locations endpoint."""
import logging
import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["phys_locations"], indirect=True)
def test_phys_locations_contract(
	to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: list[dict[str, object] | list[object] | primitive],
	phys_locations_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from phys_locations endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param get_phys_location_data: Fixture to get phys_location data from a prerequisites file.
	:param phys_location_prereq: Fixture to get sample phys_location data and actual 
	phys_location response.
	"""
	# validate phys_location keys from phys_locations get response
	logger.info("Accessing /phys_locations endpoint through Traffic ops session.")

	phys_location = request_template_data[0]
	if not isinstance(phys_location, dict):
		raise TypeError("malformed phys_location in prerequisite data; not an object")

	phys_location_name = phys_location.get("name")
	if not isinstance(phys_location_name, str):
		raise TypeError("malformed phys_location in prerequisite data; 'name' not a string")

	phys_location_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_physical_locations(query_params={"name": phys_location_name})
	try:
		phys_location_data = phys_location_get_response[0]
		if not isinstance(phys_location_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_phys_location = phys_location_data[0]
		if not isinstance(first_phys_location, dict):
			raise TypeError("malformed API response; first phys_location in response is not an object")
		phys_location_keys = set(first_phys_location.keys())

		logger.info("phys_location Keys from phys_locations endpoint response %s", phys_location_keys)
		phys_location_response_template = response_template_data.get("phys_locations")
		if not isinstance(phys_location_response_template, dict):
			raise TypeError(
				f"phys_loc response template data must be a dict, not'{type(phys_location_response_template)}'")
		response_template: dict[str, list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		]
		response_template = phys_location_response_template.get("properties")
		# validate phys_location values from prereq data in phys_locations get response.
		prereq_values = [
			phys_locations_post_data["name"],
			phys_locations_post_data["address"],
            phys_locations_post_data["city"],
            phys_locations_post_data["zip"],
			phys_locations_post_data["comments"],
            phys_locations_post_data["email"],
            phys_locations_post_data["phone"],
			phys_locations_post_data["poc"],
            phys_locations_post_data["regionId"],
			phys_locations_post_data["shortName"],
            phys_locations_post_data["state"]]

		get_values = [first_phys_location["name"], first_phys_location["address"],
         first_phys_location["city"], first_phys_location["zip"], first_phys_location["comments"],
         first_phys_location["email"], first_phys_location["phone"], first_phys_location["poc"],
         first_phys_location["regionId"], first_phys_location["shortName"], first_phys_location["state"]]

		get_types = {}
		for key, value in first_phys_location.items():
			get_types[key] = type(value).__name__
		logger.info("types from phys_location get response %s", get_types)
		response_template_types= {}
		for key, value in response_template.items():
			actual_type = value.get("type")
			if not isinstance(actual_type, str):
				raise TypeError(
					f"Type data must be a string, not '{type(actual_type)}'")
			response_template_types[key] = actual_type
		logger.info("types from phys_location response template %s", response_template_types)
		#validate keys, data types and values from regions get json response.
		assert phys_location_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for phys_locations endpoint: API response was malformed")
	finally:
		# Delete phys_location after test execution to avoid redundancy.
		try:
			physical_location_id = phys_locations_post_data["id"]
			to_session.delete_physical_location(physical_location_id=physical_location_id)
		except IndexError:
			logger.error("phys_location returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_get_phys_location")
