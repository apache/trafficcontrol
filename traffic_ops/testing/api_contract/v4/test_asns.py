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

"""API Contract Test Case for asns endpoint."""
import logging
from typing import Union

import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]

def test_asn_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive,
					 list[Union[Primitive, dict[str, object], list[object]]],
	dict[object, object]]], asn_post_data: dict[str, object]) -> None:
	"""
	Test step to validate keys, values and data types from asns endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param asn_post_data: Fixture to get sample asn data and actual asn response.
	"""
	# validate asn keys from asns get response
	logger.info("Accessing /asns endpoint through Traffic ops session.")

	asn_id = asn_post_data.get("id")
	if not isinstance(asn_id, int):
		raise TypeError("malformed asn in prerequisite data; 'asn_id' not a integer")

	asn_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_asns(query_params={"id": asn_id})
	try:
		asn_data = asn_get_response[0]
		if not isinstance(asn_data, list):
			raise TypeError("malformed API response; 'response' property not an array")

		first_asn = asn_data[0]
		if not isinstance(first_asn, dict):
			raise TypeError("malformed API response; first asn in response is not an dict")
		logger.info("asn Api get response %s", first_asn)

		asn_response_template = response_template_data.get("asns")
		if not isinstance(asn_response_template, dict):
			raise TypeError(
				f"asn response template data must be a dict, not '{type(asn_response_template)}'")

		# validate asn values from prereq data in asns get response.
		prereq_values = asn_post_data["asn"]
		get_values = first_asn["asn"]

		assert validate(instance=first_asn, schema=asn_response_template) is None
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail("API contract test failed for asns endpoint: API response was malformed")
