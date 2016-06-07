package com.comcast.cdn.traffic_control.traffic_router.api.controllers;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.http.HttpServletRequest;
import java.net.URL;
import java.net.URLDecoder;
import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/deliveryservices")
public class DeliveryServicesController {
	@Autowired
	TrafficRouterManager trafficRouterManager;

	@RequestMapping
	public @ResponseBody
	ResponseEntity<Map<String, String>> getDeliveryService(final HttpServletRequest request, @RequestParam(name = "url") final String url) {
		final URL decodedUrl;
		try {
			decodedUrl = new URL(URLDecoder.decode(url, "UTF-8"));
		} catch (Exception e) {
			return ResponseEntity.badRequest().body(null);
		}

		final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
		final DeliveryService deliveryService = trafficRouter.getCacheRegister().getDeliveryService(new HTTPRequest(request, decodedUrl), true);

		if (deliveryService == null) {
			return ResponseEntity.status(HttpStatus.NOT_FOUND).body(null);
		}

		final Map<String, String> map = new HashMap<String, String>();
		map.put("id", deliveryService.getId());

		return ResponseEntity.ok(map);
	}
}
