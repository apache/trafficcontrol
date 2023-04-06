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
import json
import logging

import pytest
import requests

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.fixture(name="cdn_prereq_data")
def get_cdn_prereq_data(
	pytestconfig: pytest.Config
) -> list[dict[str, object] | list[object] | primitive]:
	"""
	PyTest Fixture to store POST request body data for cdns endpoint.

	:returns: Prerequisite data for cdns endpoint.
	"""
	prereq_path = pytestconfig.getoption("prerequisites")
	if not isinstance(prereq_path, str):
		# unlike the configuration file, this must be present
		raise ValueError("prereqisites path not configured")

	# Response keys for cdns endpoint
	data: dict[
		str,
		list[dict[str, object] | list[object] | primitive] |\
			dict[object, object] |\
			primitive
		] |\
	primitive = None
	with open(prereq_path, encoding="utf-8", mode="r") as prereq_file:
		data = json.load(prereq_file)
	if not isinstance(data, dict):
		raise TypeError(f"prerequisite data must be an object, not '{type(data)}'")

	cdn_data = data["cdns"]
	if not isinstance(cdn_data, list):
		raise TypeError(f"cdns data must be a list, not '{type(cdn_data)}'")

	return cdn_data


def test_cdn_contract(
	to_session: TOSession,
	cdn_prereq_data: list[dict[str, object] | list[object] | primitive],
	cdn_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from cdns endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param get_cdn_data: Fixture to get CDN data from a prerequisites file.
	:param cdn_prereq: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing /cdns endpoint through Traffic ops session.")

	cdn = cdn_prereq_data[0]
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
		# validate cdn values from prereq data in cdns get response.
		prereq_values = [
			cdn_post_data["name"],
			cdn_post_data["domainName"],
			cdn_post_data["dnssecEnabled"]
		]
		get_values = [first_cdn["name"], first_cdn["domainName"], first_cdn["dnssecEnabled"]]
		# validate data types for values from cdn get json response.
		for (prereq_value, get_value) in zip(prereq_values, get_values):
			assert isinstance(prereq_value, type(get_value))
		assert cdn_keys == set(cdn_post_data.keys())
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Either prerequisite data or API response was malformed")
	finally:
		# Delete CDN after test execution to avoid redundancy.
		try:
			cdn_id = cdn_post_data["id"]
			to_session.delete_cdn_by_id(cdn_id=cdn_id)
		except IndexError:
			logger.error("CDN returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_get_cdn")
