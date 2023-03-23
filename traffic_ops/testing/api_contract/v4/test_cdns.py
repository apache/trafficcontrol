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
import os
import pytest
import requests
from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.fixture(name="cdn_prereq_data")
def get_cdn_prereq_data(pytestconfig: pytest.Config) -> list[dict[str, object] | list[object] | primitive]:
	"""
	PyTest Fixture to store POST request body data for cdns endpoint.

	:returns: Prerequisite data for cdns endpoint.
	"""
	prereq_path = pytestconfig.getoption("prerequisites")
	if not isinstance(prereq_path, str):
		# unlike the configuration file, this must be present
		raise ValueError("prereqisites path not configured")

	# Response keys for cdns endpoint
	data: dict[str, list[dict[str, object] | list[object] | primitive] | dict[str, object] | primitive] | list[object] | primitive = None
	with open(prereq_path, encoding="utf-8", mode="r") as prereq_file:
		data = json.load(prereq_file)
	if not isinstance(data, dict):
		raise TypeError(f"prerequisite data must be an object, not '{type(data)}'")
	cdn_data = data["cdns"]
	return cdn_data


def test_cdn_contract(to_session: TOSession, cdn_prereq_data: object, cdn_post_data: list[dict[str, str] | requests.Response]) -> None:
	"""
	Test step to validate keys, values and data types from cdns endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param get_cdn_data: Fixture to get CDN data from a prerequisites file.
	:param cdn_prereq: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing Cdn endpoint through Traffic ops session.")
	cdn_name = cdn_post_data[0]["name"]
	cdn_get_response = to_session.get_cdns(
		query_params={"name": cdn_name})
	try:
		cdn_data = cdn_get_response[0]
		cdn_keys = list(cdn_data[0].keys())
		logger.info("CDN Keys from cdns endpoint response %s", cdn_keys)
		# validate cdn values from prereq data in cdns get response.
		prereq_values = [cdn_post_data[0]["name"], cdn_post_data[0]["domainName"], cdn_post_data[0]["dnssecEnabled"]]
		get_values = [cdn_data[0]["name"], cdn_data[0]["domainName"], cdn_data[0]["dnssecEnabled"]]
		# validate data types for values from cdn get json response.
		for (prereq_value, get_value) in zip(prereq_values, get_values):
			assert isinstance(prereq_value, type(get_value))
		assert cdn_keys.sort() == list(cdn_prereq_data.keys()).sort()
		assert get_values == prereq_values
	except IndexError:
		logger.error("No CDN data from cdns get request")
		pytest.fail("Response from get request is empty, Failing test_get_cdn")
	finally:
		# Delete CDN after test execution to avoid redundancy.
		try:
			cdn_response = cdn_post_data[1]
			cdn_id = cdn_response["id"]
			to_session.delete_cdn_by_id(cdn_id=cdn_id)
		except IndexError:
			logger.error("CDN wasn't created")
			pytest.fail("Response from delete request is empty, Failing test_get_cdn")
