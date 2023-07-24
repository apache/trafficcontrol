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

"""API Contract Test Case for topologies endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

def test_topology_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive,
					 list[Union[Primitive, dict[str, object], list[object]]],
	dict[object, object]]], topology_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from topologies endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param topology_post_data: Fixture to get sample topology data and actual topology response.
	"""
	# validate topology keys from topologies get response
	logger.info("Accessing /topologies endpoint through Traffic ops session.")

	topology_name = topology_post_data.get("name")
	if not isinstance(topology_name, str):
		raise TypeError("malformed topology in prerequisite data; 'topology_name' not a string")

	topology_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_topologies(query_params={"name": topology_name})
	try:
		topology_data = topology_get_response[0]
		if not isinstance(topology_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_topology = topology_data[0]
		if not isinstance(first_topology, dict):
			raise TypeError("malformed API response; first topology in response is not an dict")
		logger.info("topology Api get response %s", first_topology)

		topology_response_template = response_template_data.get("topologies")
		if not isinstance(topology_response_template, dict):
			raise TypeError(
				f"topology response template data must be a dict, not '{type(topology_response_template)}'")

		# validate topology values from prereq data in topologies get response.
		prereq_values = [topology_post_data["name"], topology_post_data["description"]]
		get_values = [first_topology["name"], first_topology["description"]]

		assert validate(instance=first_topology, schema=topology_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for topologies endpoint: API response was malformed")
