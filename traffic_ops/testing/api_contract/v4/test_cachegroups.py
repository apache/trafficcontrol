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

"""API Contract Test Case for cachegroup endpoint."""
import logging
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None


def test_cachegroup_contract(to_session: TOSession,
	request_template_data: list[dict[str, object] | list[object] | primitive],
	response_template_data: dict[str, primitive | list[primitive | dict[str, object]
						    | list[object]] | dict[object, object]],
	cachegroup_post_data: dict[str, object]
	) -> None:
	"""
	Test step to validate keys, values and data types from cachegroup endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param request_template_data: Fixture to get request template data from a prerequisites file.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param cachegroup_post_data: Fixture to get sample cachegroup data and actual cachegroup response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing /cachegroup endpoint through Traffic ops session.")

	cachegroup = request_template_data["cachegroup"][0]
	if not isinstance(cachegroup, dict):
		raise TypeError("malformed cachegroup in prerequisite data; not an object")

	cachegroup_name = cachegroup.get("name")
	if not isinstance(cachegroup_name, str):
		raise TypeError("malformed cachegroup in prerequisite data; 'name' not a string")

	cachegroup_get_response: tuple[
		dict[str, object] | list[dict[str, object] | list[object] | primitive] | primitive,
		requests.Response
	] = to_session.get_cachegroups(query_params={"name": str(cachegroup_name)})

	try:
		cachegroup_data = cachegroup_get_response[0]
		if not isinstance(cachegroup_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_cachegroup = cachegroup_data[0]
		if not isinstance(first_cachegroup, dict):
			raise TypeError("malformed API response; first Cache group in response is not an object")
		logger.info(f"Server Api response {first_cachegroup}")

		cachegroup_response_template = response_template_data.get("cachegroup")
		
		# validate cachegroup values from prereq data in cachegroup get response.
		prereq_values = [cachegroup_post_data["name"],cachegroup_post_data["shortName"],
		   cachegroup_post_data["fallbackToClosest"],cachegroup_post_data["typeId"]]

		get_values = [first_cachegroup["name"],first_cachegroup["shortName"],
		first_cachegroup["fallbackToClosest"],first_cachegroup["typeId"]]

		# validate keys,data types for values and actual values for cachegroup endpoint.
		assert validate(instance=first_cachegroup, schema=cachegroup_response_template) == None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for cachegroup endpoint: API response was malformed")
	finally:
		# Delete Cache group after test execution to avoid redundancy.
		try:
			cachegroup_id = cachegroup_post_data["id"]
			to_session.delete_cachegroups(cache_group_id=cachegroup_id)
		except IndexError:
			logger.error("Cachegroup returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_cachegroup_contract")
