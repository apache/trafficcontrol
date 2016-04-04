package com.comcast.cdn.traffic_control.traffic_router.api.controllers;

import com.comcast.cdn.traffic_control.traffic_router.core.util.DataExporter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/stats/zones")
public class ZonesController {
	@Autowired
	DataExporter dataExporter;


	@RequestMapping(value = "/caches")
	public @ResponseBody
	Map<String, Object> getAllCachesStats() {
		final Map<String, Object> statsMap = new HashMap<String, Object>();
		statsMap.put("dynamicZoneCaches", dataExporter.getDynamicZoneCacheStats());
		statsMap.put("staticZoneCaches", dataExporter.getStaticZoneCacheStats());
		return statsMap;
	}

	@RequestMapping(value = "/caches/{filter:static|dynamic}")
	public @ResponseBody
	Map<String, Object> getCachesStats(@PathVariable("filter") final String filter) {
		return "static".equals(filter) ? dataExporter.getStaticZoneCacheStats() : dataExporter.getDynamicZoneCacheStats();
	}
}
