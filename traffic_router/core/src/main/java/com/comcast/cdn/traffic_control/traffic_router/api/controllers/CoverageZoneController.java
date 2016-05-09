package com.comcast.cdn.traffic_control.traffic_router.api.controllers;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
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
		final CacheLocation cacheLocation = trafficRouterManager.getTrafficRouter().getCoverageZoneCacheLocation(ip, deliveryServiceId);
		if (cacheLocation == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		return ResponseEntity.ok(cacheLocation);
	}

	@RequestMapping(value = "/caches")
	public @ResponseBody
	ResponseEntity<List<Cache>> getCachesForDeliveryService(@RequestParam(name = "deliveryServiceId") final String deliveryServiceId,
	                                                        @RequestParam(name = "cacheLocationId") final String cacheLocationId) {
		final List<Cache> caches = trafficRouterManager.getTrafficRouter().selectCachesByCZ(deliveryServiceId, cacheLocationId, null);
		if (caches == null || caches.isEmpty()) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		return ResponseEntity.ok(caches);
	}
}
