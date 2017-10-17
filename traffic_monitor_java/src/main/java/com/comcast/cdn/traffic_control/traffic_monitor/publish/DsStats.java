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

import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.mapper.parameter.PageParameters;

public class DsStats extends JsonPage {
	private static final long serialVersionUID = 1L;
	private final DeliveryServiceStateRegistry deliveryServiceStateRegistry = DeliveryServiceStateRegistry.getInstance();


	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
		String str = pp.get("hc").toString();
		int hc = 0;
		try {
			hc = Integer.parseInt(str);
		} catch(Exception e) {hc = 0;}
		String[] stats = null;
		str = pp.get("stats").toString();
		if(str != null) {
			stats = str.split(",");
		}
		final boolean wildcard = pp.get("wildcard").toBoolean(false);
		final boolean hidden = pp.get("hidden").toBoolean(false);
		final String deliveryService = pp.get(0).toString();
		final JSONObject o = new JSONObject();
		o.put("date", new Date().toString());
		o.put("pp", pp);
		final JSONObject servers = new JSONObject();
		if(deliveryService != null && !deliveryService.equals("")) {
			if(deliveryServiceStateRegistry.has(deliveryService)) {
				servers.put(deliveryService, deliveryServiceStateRegistry.get(deliveryService).getStatsJson(hc, stats, wildcard, hidden));
			} else {
				o.put("error", "Delivery service not found: "+deliveryService);
			}
		} else {
			for(AbstractState s : deliveryServiceStateRegistry.getAll()) {
				servers.put(s.getId(),s.getStatsJson(hc, stats, wildcard, hidden));
			}
		}
		o.put("deliveryService", servers);
		return o;
	}

}

