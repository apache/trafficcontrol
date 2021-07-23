/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.traffic_control.traffic_router.api.controllers;

import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/consistenthash")
public class ConsistentHashController {
	@Autowired
	TrafficRouterManager trafficRouterManager;
	final static int MAX_REQUEST_PATH_LENGTH = 28;
	final static String RESULTING_PATH_TO_HASH = "resultingPathToConsistentHash";
	final static String REQUEST_PATH = "requestPath";
	final static String CONSISTENT_HASH_REGEX = "consistentHashRegex";
	final static String DELIVERY_SERVICE_ID = "deliveryServiceId";

	@RequestMapping(value = "/cache/coveragezone")
	public @ResponseBody
	ResponseEntity hashCoverageZoneCache(@RequestParam(name="ip") final String ip,
	                                @RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
	                                @RequestParam(name = REQUEST_PATH) final String requestPath) {

		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashForCoverageZone(ip, deliveryServiceId, requestPath);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}

	@RequestMapping(value = "/cache/deep/coveragezone")
	public @ResponseBody
	ResponseEntity hashCoverageZoneDeepCache(@RequestParam(name="ip") final String ip,
										 @RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
										 @RequestParam(name = REQUEST_PATH) final String requestPath) {

		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashForCoverageZone(ip, deliveryServiceId, requestPath, true);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}

	@RequestMapping(value = "/cache/geolocation")
	public @ResponseBody
	ResponseEntity hashGeolocatedCache(@RequestParam(name="ip") final String ip,
	                                @RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
	                                @RequestParam(name = REQUEST_PATH) final String requestPath) {
		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashForGeolocation(ip, deliveryServiceId, requestPath);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}

	@RequestMapping(value = "/deliveryservice")
	public @ResponseBody
	ResponseEntity hashDeliveryService(@RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
	                                   @RequestParam(name = REQUEST_PATH) final String requestPath) {

		final DeliveryService deliveryService = trafficRouterManager.getTrafficRouter().consistentHashDeliveryService(deliveryServiceId, requestPath);

		if (deliveryService == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(deliveryService);
	}

	@RequestMapping(value = "/patternbased/regex")
	public @ResponseBody
	ResponseEntity<Map<String, String>> testPatternBasedRegex(@RequestParam(name = "regex") final String regex,
										 @RequestParam(name = REQUEST_PATH) final String requestPath) {

		// limit length of requestPath to protect against evil regexes
		if (requestPath != null && requestPath.length() > MAX_REQUEST_PATH_LENGTH) {
			final Map<String, String> map = new HashMap<String, String>();
			map.put("Bad Input", "Request Path length is restricted by API to " + MAX_REQUEST_PATH_LENGTH + " characters");
			return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(map);
		}

		final String pathToHash = trafficRouterManager.getTrafficRouter().buildPatternBasedHashString(regex, requestPath);

		if (pathToHash == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		final Map<String, String> map = new HashMap<String, String>();
		map.put(REQUEST_PATH, requestPath);
		map.put(CONSISTENT_HASH_REGEX, regex);
		map.put(RESULTING_PATH_TO_HASH, pathToHash);
		return ResponseEntity.ok(map);
	}

	@RequestMapping(value = "/patternbased/deliveryservice")
	public @ResponseBody
	ResponseEntity<Map<String, String>> testPatternBasedDeliveryService(@RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
																		@RequestParam(name = REQUEST_PATH) final String requestPath) {

		final String pathToHash = trafficRouterManager.getTrafficRouter().buildPatternBasedHashStringDeliveryService(deliveryServiceId, requestPath);

		if (pathToHash == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		final Map<String, String> map = new HashMap<String, String>();
		map.put(REQUEST_PATH, requestPath);
		map.put(DELIVERY_SERVICE_ID, deliveryServiceId);
		map.put(RESULTING_PATH_TO_HASH, pathToHash);
		return ResponseEntity.ok(map);
	}

	@RequestMapping(value = "/cache/coveragezone/steering")
	public @ResponseBody
	ResponseEntity hashSteeringCoverageZoneCache(@RequestParam(name="ip") final String ip,
										 @RequestParam(name = DELIVERY_SERVICE_ID) final String deliveryServiceId,
										 @RequestParam(name = REQUEST_PATH) final String requestPath) {

		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashSteeringForCoverageZone(ip, deliveryServiceId, requestPath);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}
}
