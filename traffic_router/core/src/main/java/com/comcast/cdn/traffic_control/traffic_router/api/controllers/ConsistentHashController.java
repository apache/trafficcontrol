package com.comcast.cdn.traffic_control.traffic_router.api.controllers;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
@RequestMapping("/consistenthash")
public class ConsistentHashController {
	@Autowired
	TrafficRouterManager trafficRouterManager;

	@RequestMapping(value = "/cache/coveragezone")
	public @ResponseBody
	ResponseEntity hashCoverageZoneCache(@RequestParam(name="ip") final String ip,
	                                @RequestParam(name = "deliveryServiceId") final String deliveryServiceId,
	                                @RequestParam(name = "requestPath") final String requestPath) {

		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashForCoverageZone(ip, deliveryServiceId, requestPath);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}

	@RequestMapping(value = "/cache/geolocation")
	public @ResponseBody
	ResponseEntity hashGeolocatedCache(@RequestParam(name="ip") final String ip,
	                                @RequestParam(name = "deliveryServiceId") final String deliveryServiceId,
	                                @RequestParam(name = "requestPath") final String requestPath) {
		final Cache cache = trafficRouterManager.getTrafficRouter().consistentHashForGeolocation(ip, deliveryServiceId, requestPath);

		if (cache == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(cache);
	}

	@RequestMapping(value = "/deliveryservice")
	public @ResponseBody
	ResponseEntity hashDeliveryService(@RequestParam(name = "deliveryServiceId") final String deliveryServiceId,
	                                   @RequestParam(name = "requestPath") final String requestPath) {

		final DeliveryService deliveryService = trafficRouterManager.getTrafficRouter().consistentHashDeliveryService(deliveryServiceId, requestPath);

		if (deliveryService == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body("{}");
		}

		return ResponseEntity.ok(deliveryService);
	}
}
