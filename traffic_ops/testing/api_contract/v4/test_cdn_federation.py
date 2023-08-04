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

"""API Contract Test Case for cdn_federation endpoint."""
import logging
from random import randint
from typing import Union
import pytest
import requests
from jsonschema import validate

from trafficops.tosession import TOSession

# Create and configure logger
logger = logging.getLogger()

Primitive = Union[bool, int, float, str, None]


def test_cdn_federation_contract(to_session: TOSession,
	response_template_data: dict[str, Union[Primitive, list[Union[Primitive,
							dict[str, object], list[object]]], dict[object, object]]],
	cdn_federation_post_data: dict[str, object]
) -> None:
	"""
	Test step to validate keys, values and data types from cdn_federation endpoint
	response.
	:param to_session: Fixture to get Traffic Ops session.
	:param response_template_data: Fixture to get response template data from a prerequisites file.
	:param cdn_federation_post_data: Fixture to get delivery services sslkeys data.
	"""
	# validate cdn_federation keys from api get response
	logger.info("Accessing /cdn_federation endpoint through Traffic ops session.")

	cdn_name = cdn_federation_post_data[0]
	cdn_federation_data_post = cdn_federation_post_data[1]
	cdn_federation_sample_data = cdn_federation_post_data[2]
	federation_id = cdn_federation_post_data[3]
	if not isinstance(cdn_name, str):
		raise TypeError("malformed API response; 'cdn_name' property not a string")

	cdn_federation_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federations_for_cdn(cdn_name=cdn_name)
	
	user_federation_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federation_users(federation_id=federation_id)

	delivery_service_federation_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federation_delivery_services(federation_id=federation_id)

	federation_resolvers_get_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive],
		requests.Response
	] = to_session.get_federation_resolvers_by_id(federation_id=federation_id)
	
    # Hitting cdn_federation PUT method
	cdn_federation_sample_data["description"] = "test" + str(randint(0, 1000)) 
	logger.info("Updated cdn federation data to hit PUT method %s", cdn_federation_sample_data)
	cdn_federation_put_response: tuple[
		Union[dict[str, object], list[Union[dict[str, object], list[object], Primitive]], Primitive], 
		requests.Response] = to_session.update_federation_in_cdn(
		cdn_name=cdn_name, federation_id=federation_id, data= cdn_federation_sample_data)
    
	try:
		cdn_federation = cdn_federation_get_response[0]
		if not isinstance(cdn_federation, list):
			raise TypeError(
				"malformed API response; first cdn_federation in response is not an dict")
		first_cdn_federation = cdn_federation[0]

		user_federation = user_federation_get_response[0]
		if not isinstance(user_federation, list):
			raise TypeError(
				"malformed API response; first user_federation in response is not an dict")
		first_user_federation = user_federation[0]
	
		delivery_service_federation = delivery_service_federation_get_response[0]
		if not isinstance(delivery_service_federation, list):
			raise TypeError(
				"malformed API response; first delivery_service_federation in response is not an dict")
		first_delivery_service_federation = delivery_service_federation[0]
		
		cdn_federation_put_data = cdn_federation_put_response[0]
		if not isinstance(cdn_federation_put_data, dict):
			raise TypeError(
				"malformed API response; first cdn_federation_put in response is not an dict")
		
		federation_federation_resolver= federation_resolvers_get_response[0]
		if not isinstance(federation_federation_resolver, list):
			raise TypeError(
				"malformed API response; first federation_federation_resolver in response is not an dict")
		first_federation_federation_resolver = federation_federation_resolver[0]

		cdn_federation_response_template = response_template_data.get(
			"cdn_federation")
		if not isinstance(cdn_federation_response_template, dict):
			raise TypeError(f"cdn_federation response template data must be a dict, not '"
							f"{type(cdn_federation_response_template)}'")

		user_federation_response_template = response_template_data.get(
			"user_federation")
		if not isinstance(user_federation_response_template, dict):
			raise TypeError(f"user_federation response template data must be a dict, not '"
							f"{type(user_federation_response_template)}'")

		delivery_service_federation_response_template = response_template_data.get(
			"delivery_service_federation")
		if not isinstance(delivery_service_federation_response_template, dict):
			raise TypeError(f"delivery_service_federation response template data must be a dict, not '"
							f"{type(user_federation_response_template)}'")
		
		federation_federation_resolver_response_template = response_template_data.get("federation_federation_resolver")
		if not isinstance(federation_federation_resolver_response_template, dict):
			raise TypeError(f"delivery_service_federation response template data must be a dict, not '"
							f"{type(user_federation_response_template)}'")

		keys = ["cname", "ttl"]
		prereq_values = [cdn_federation_data_post[key] for key in keys]
		get_values = [first_cdn_federation[key] for key in keys]

		assert validate(instance=first_cdn_federation, schema=cdn_federation_response_template) is None
		assert validate(instance=first_user_federation, schema=user_federation_response_template) is None
		assert validate(instance=first_delivery_service_federation, schema=delivery_service_federation_response_template) is None
		assert validate(instance=cdn_federation_put_data, schema=cdn_federation_response_template) is None
		assert validate(instance=first_federation_federation_resolver, schema=federation_federation_resolver_response_template) is None
		
		assert get_values == prereq_values
	except IndexError:
		logger.error("Either prerequisite data or API response was malformed")
		pytest.fail(
			"API contract test failed for cdn_federation endpoint: API response was malformed")
