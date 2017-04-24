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

import java.util.List;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DsState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerWatcher;

public class CrStates extends JsonPage {
	private static final Logger LOGGER = Logger.getLogger(CrStates.class);
	private static final long serialVersionUID = 1L;
	private static CacheWatcher myCacheWatcher;
	private static PeerWatcher myPeerWatcher;
	private static HealthDeterminer myHealthDeterminer;

	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
		if (myPeerWatcher == null) {
			return null;
		}

		final MonitorConfig config = ConfigHandler.getInstance().getConfig();
		final RouterConfig crConfig = RouterConfig.getCrConfig();

		if (crConfig == null || myCacheWatcher.getCycleCount() < config.getStartupMinCycles()) {
			return null;
		}

		final boolean raw = (pp.getPosition("raw") != -1);
		final String cacheType = pp.get("cacheType").toString();
		final JSONObject o = new JSONObject();
		o.put("caches", getCrStates(crConfig, raw, cacheType));

		if (ConfigHandler.getInstance().getConfig().getPublishDsStates()) {
			o.put("deliveryServices", getDsStates(crConfig));
		}

		return o;
	}

	private JSONObject getCrStates(final RouterConfig crConfig, final boolean raw, final String cacheType) {
		if (crConfig == null) {
			return null;
		}

		try {
			final JSONObject servers = new JSONObject();
			final List<Cache> caches = crConfig.getCacheList();

			for (Cache c : caches) {
				synchronized(c) {
					if (c.getControls() == null || (cacheType != null && !cacheType.equals(c.getType()))) {
						continue;
					}

					final MonitorConfig config = ConfigHandler.getInstance().getConfig();
					servers.put(c.getHostname(), myHealthDeterminer.getJSONStats(c, config.getPeerOptimistic(), raw));
				}
			}

			if (servers.length() == 0) {
				LOGGER.warn("no caches returned! ");
			}

			return servers;
		} catch (JSONException e) {
			LOGGER.warn(e, e);
		}

		return null;
	}

	private JSONObject getDsStates(final RouterConfig crConfig) {
		if (crConfig == null) {
			return null;
		}

		try {
			final JSONObject ret = new JSONObject();

			for (DsState dsState : DeliveryServiceStateRegistry.getInstance().getAll()) {
				final JSONObject dsJo = new JSONObject();

				if (!dsState.hasValue(AbstractState.IS_AVAILABLE_STR)) {
					continue;
				}

				dsJo.put(AbstractState.IS_AVAILABLE_STR, dsState.isAvailable());
				dsJo.put(DsState.DISABLED_LOCATIONS, dsState.getDisabledLocations());
				ret.put(dsState.getId(), dsJo);
			}

			return ret;
		} catch (JSONException e) {
			LOGGER.warn(e, e);
		}

		return null;
	}

	public static void init(final CacheWatcher cw, final PeerWatcher pw, final HealthDeterminer hd) {
		myCacheWatcher = cw;
		myPeerWatcher = pw;
		myHealthDeterminer = hd;
	}
}
