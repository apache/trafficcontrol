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

import org.apache.traffic_control.traffic_router.core.util.DataExporter;
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
