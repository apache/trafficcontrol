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

package com.comcast.cdn.traffic_control.traffic_monitor.publish;

import java.util.Date;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;

public class CacheStats extends JsonPage {
	private static final long serialVersionUID = 1L;
	private final CacheStateRegistry cacheStateRegistry = CacheStateRegistry.getInstance();

	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
		String[] stats = null;
		final String str = pp.get("stats").toString();
		final int hc = pp.get("hc").toInt(0);

		if (str != null) {
			stats = str.split(",");
		}

		final String type = pp.get("type").toString();
		final boolean wildcard = pp.get("wildcard").toBoolean(false);
		final boolean hidden = pp.get("hidden").toBoolean(false);
		final String host = pp.get(0).toString();
		final JSONObject o = new JSONObject();
		o.put("date", new Date().toString());
		o.put("pp", pp);
		final JSONObject servers = new JSONObject();

		if (host != null && !host.equals("")) {
			if (cacheStateRegistry.has(host)) {
				servers.put(host, cacheStateRegistry.get(host).getStatsJson(hc, stats, wildcard, hidden));
			} else {
				o.put("error", "Hostname not found: " + host);
			}
		} else {
			for (CacheState cacheState : cacheStateRegistry.getAll()) {
				if (type == null || type.equals(cacheState.getCache().getType())) {
					servers.put(cacheState.getId(), cacheState.getStatsJson(hc, stats, wildcard, hidden));
				}
			}
		}

		o.put("caches", servers);

		return o;
	}

}
