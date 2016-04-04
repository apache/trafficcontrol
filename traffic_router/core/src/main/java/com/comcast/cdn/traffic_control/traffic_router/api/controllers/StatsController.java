package com.comcast.cdn.traffic_control.traffic_router.api.controllers;

import com.comcast.cdn.traffic_control.traffic_router.core.util.DataExporter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/stats")
public class StatsController {
	@Autowired
	private DataExporter dataExporter;

	@RequestMapping
	public @ResponseBody
	Map<String, Object> getStats() {
		final Map<String, Object> map = new HashMap<String, Object>();

		map.put("app", dataExporter.getAppInfo());
		map.put("stats", dataExporter.getStatTracker());

		return map;
	}

	@RequestMapping(value = "/ip/{ip:.+}")
	public @ResponseBody
	Map<String, Object> getCaches(@PathVariable("ip") final String ip,
	                              @RequestParam(name = "geolocationProvider", required = false, defaultValue = "maxmindGeolocationService") final String geolocationProvider) {
		return dataExporter.getCachesByIp(ip, geolocationProvider);
	}
}
