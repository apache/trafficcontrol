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
import org.springframework.web.bind.annotation.GetMapping;
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

	@GetMapping
	public @ResponseBody
	Map<String, Object> getStats() {
		final Map<String, Object> map = new HashMap<String, Object>();

		map.put("app", dataExporter.getAppInfo());
		map.put("stats", dataExporter.getStatTracker());

		return map;
	}

	@GetMapping(value = "/ip/{ip:.+}")
	public @ResponseBody
	Map<String, Object> getCaches(@PathVariable("ip") final String ip,
	                              @RequestParam(name = "geolocationProvider", required = false, defaultValue = "maxmindGeolocationService") final String geolocationProvider) {
		return dataExporter.getCachesByIp(ip, geolocationProvider);
	}
}
