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

"""API Contract Test Case for servers endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_server_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]],
	dict[object, object]]], server_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from servers endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param server_post_data: Fixture to get sample server data and actual server response.
	"""
	# validate server keys from server get response
	logger.info("Accessing /servers endpoint through Traffic ops session.")
	profile_id = server_post_data[1]
	server_post_data = server_post_data[0]

	server_id = server_post_data.get("id")
	if not isinstance(server_id, int):
		raise TypeError("malformed API response; 'id' property not a integer")

	server_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_servers(query_params={"id": server_id})
	try:
		server_data = server_get_response[0]
		if not isinstance(server_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_server = server_data[0]
		if not isinstance(first_server, dict):
			raise TypeError("malformed API response; first Server in response is not an dict")
		logger.info("Server Api get response %s", first_server)
		server_response_template = response_template_data.get("servers")
		if not isinstance(server_response_template, dict):
			raise TypeError(
				f"Server response template data must be a dict, not '{type(server_response_template)}'")

		keys = ["cachegroupId", "cdnId", "domainName", "physLocationId", "profileNames",
	            "statusId", "typeId"]
		prereq_values = [server_post_data[key] for key in keys]
		get_values = [first_server[key] for key in keys]

		assert validate(instance=first_server, schema=server_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for server endpoint: API response was malformed")
	finally:
		# Delete Server after test execution to avoid redundancy.
		server_id = server_post_data.get("id")
		if to_session.delete_server_by_id(server_id=server_id) is None:
			logger.error("Server returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_contract")

		cachegroup_id = server_post_data.get("cachegroupId")
		if to_session.delete_cachegroups(cache_group_id=cachegroup_id) is None:
			logger.error("cachegroup returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_contract")

		if to_session.delete_profile_by_id(profile_id=profile_id) is None:
			logger.error("Profile returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_contract")

		cdn_id = server_post_data.get("cdnId")
		if to_session.delete_cdn_by_id(cdn_id=cdn_id) is None:
			logger.error("Cdn returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_contract")

		phys_loc_id = server_post_data.get("physLocationId")
		if to_session.delete_physical_location(physical_location_id=phys_loc_id) is None:
			logger.error("Physical location returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_server_contract")
