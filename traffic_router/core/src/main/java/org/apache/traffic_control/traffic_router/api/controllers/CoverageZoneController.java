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
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.Node.IPVersions;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import java.util.List;

@Controller
@RequestMapping("/coveragezone")
public class CoverageZoneController {
	@Autowired
	TrafficRouterManager trafficRouterManager;

	@RequestMapping(value = "/cachelocation")
	public @ResponseBody
	ResponseEntity<CacheLocation> getCacheLocationForIp(@RequestParam(name = "ip") final String ip,
	                                    @RequestParam(name = "deliveryServiceId") final String deliveryServiceId) {
		final IPVersions requestVersion = ip.contains(":") ? IPVersions.IPV6ONLY : IPVersions.IPV4ONLY;
		final CacheLocation cacheLocation = trafficRouterManager.getTrafficRouter().getCoverageZoneCacheLocation(ip, deliveryServiceId, requestVersion);
		if (cacheLocation == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		return ResponseEntity.ok(cacheLocation);
	}

	@RequestMapping(value = "/caches")
	public @ResponseBody
	ResponseEntity<List<Cache>> getCachesForDeliveryService(@RequestParam(name = "deliveryServiceId") final String deliveryServiceId,
	                                                        @RequestParam(name = "cacheLocationId") final String cacheLocationId) {
		final List<Cache> caches = trafficRouterManager.getTrafficRouter().selectCachesByCZ(deliveryServiceId, cacheLocationId, null, IPVersions.ANY);
		if (caches == null || caches.isEmpty()) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		return ResponseEntity.ok(caches);
	}
}
