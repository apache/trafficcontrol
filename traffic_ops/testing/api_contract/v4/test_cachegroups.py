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

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

primitive = bool | int | float | str | None

@pytest.mark.parametrize('request_template_data', ["cachegroup"], indirect=True)
def test_cachegroup_contract(to_session: TOSession, request_template_data:
	list[dict[str, object] | list[object] | primitive],
	response_template_data: list[dict[str, object] | list[object] | primitive],
	cachegroup_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from cachegroup endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param cachegroup_prereq_data: Fixture to get CDN data from a prerequisites file.
	:param cachegroup_post_data: Fixture to get sample CDN data and actual CDN response.
	"""
	# validate CDN keys from cdns get response
	logger.info("Accessing /cachegroup endpoint through Traffic ops session.")

	cachegroup = request_template_data[0]
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
		cachegroup_keys = set(first_cachegroup.keys())

		logger.info("Cache group Keys from cachegroup endpoint response %s", cachegroup_keys)
		response_template = response_template_data.get("cachegroup").get("properties")
		# validate cachegroup values from prereq data in cachegroup get response.
		prereq_values = [
			cachegroup_post_data["name"],
			cachegroup_post_data["shortName"],
			cachegroup_post_data["fallbackToClosest"],
			cachegroup_post_data["typeId"]
		]
		get_values = [
			first_cachegroup["name"],
	        first_cachegroup["shortName"],
	        first_cachegroup["fallbackToClosest"],
		    first_cachegroup["typeId"]
	    ]
		get_types = {}
		for key in first_cachegroup:
			get_types[key] = type(first_cachegroup[key]).__name__
		logger.info("types from cachegroup get response %s", get_types)
		response_template_types= {}
		for key in response_template:
			optional = response_template.get(key).get("optional")
			if optional:
				if key in cachegroup:
					response_template_types[key] = response_template.get(key).get("typeA")
				else:
					response_template_types[key] = response_template.get(key).get("typeB")
			else:
				response_template_types[key] = response_template.get(key).get("type")
		logger.info("types from cachegroup response template %s", response_template_types)
		# validate data types for values from cdn get json response.
		assert cachegroup_keys == set(response_template.keys())
		assert dict(sorted(get_types.items())) == dict(sorted(response_template_types.items()))
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("Either prerequisite data or API response was malformed")
	finally:
		# Delete Cache group after test execution to avoid redundancy.
		try:
			cachegroup_id = cachegroup_post_data["id"]
			to_session.delete_cachegroups(cache_group_id=cachegroup_id)
		except IndexError:
			logger.error("Cachegroup returned by Traffic Ops is missing an 'id' property")
			pytest.fail("Response from delete request is empty, Failing test_cachegroup_contract")
