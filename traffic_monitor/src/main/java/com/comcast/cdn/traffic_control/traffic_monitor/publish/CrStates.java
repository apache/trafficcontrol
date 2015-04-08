/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.servlet.http.HttpServletRequest;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.cycle.RequestCycle;
import org.apache.wicket.request.http.WebRequest;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DsState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerWatcher;

public class CrStates extends JsonPage {
	private static final Logger LOGGER = Logger.getLogger(CrStates.class);
	private static final long serialVersionUID = 1L;
	//	private static HealthDeterminer myHealthDeterminer;
	private static CacheWatcher myCacheWatcher;
	private static PeerWatcher myPeerWatcher;
	private static HealthDeterminer myHealthDeterminer;
	private static Map<String,Long> clientIps = new HashMap<String,Long>();

	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
		if(myPeerWatcher == null) {
			return null;
		}
		boolean raw = true;
		final MonitorConfig config = ConfigHandler.getConfig();
		final RouterConfig crConfig = RouterConfig.getCrConfig();
		if(crConfig == null || myCacheWatcher.getCycleCount() < config.getStartupMinCycles()) {
			return null;
		}
		if(pp == null || pp.getPosition("raw") == -1) {
			final WebRequest req = (WebRequest) RequestCycle.get().getRequest();
			final HttpServletRequest httpReq = (HttpServletRequest) req.getContainerRequest();
			clientIps.put(httpReq.getRemoteHost(), new Long(System.currentTimeMillis()));
			raw = false;
		}
		final JSONObject o = new JSONObject();
		o.put("caches", getCrStates(crConfig, raw));
		if(ConfigHandler.getConfig().getPublishDsStates()) {
			o.put("deliveryServices", getDsStates(crConfig));
		}
		return o;
	}
	public static Map<String, Long> getCrIps() {
		return new HashMap<String, Long>(clientIps);
	}

	// TODO: clean up/merge with PeerWatcher logic?
	/*private JSONObject getPeerSet(final JSONObject peers, final String id) {
		try {
			final JSONObject peerSet = new JSONObject();
			if(peers == null || peers.length() == 0) {
				return peerSet;
			}
			for(String peer : JSONObject.getNames(peers)) {
				JSONObject jo = peers.getJSONObject(peer);
				if(jo.has("caches")) {
					jo = jo.getJSONObject("caches");
				}
				if(jo.has(id)) {
					peerSet.put(peer, jo.getJSONObject(id));
				} else {
					LOGGER.warn("Cache ("+id+") not found in peer ("+peer+")");
				}
			}
			return peerSet;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
		return null;
	}*/

	private JSONObject getCrStates(final RouterConfig crConfig, final boolean raw) {
		if (crConfig == null) {
			return null;
		}
		try {
			final JSONObject servers = new JSONObject();
			final List<Cache> caches = crConfig.getCacheList();
			for(Cache c : caches) {
				synchronized(c) {
					if(c.getControls() == null) { continue; }
					final MonitorConfig config = ConfigHandler.getConfig();
					servers.put(c.getHostname(), myHealthDeterminer.getJSONStats(c, config.getPeerOptimistic(), raw));
				}
			}
			if(servers.length()==0) {
				LOGGER.warn("no caches returned! ");
			}
			return servers;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
		return null;
	}
	private JSONObject getDsStates(final RouterConfig crConfig) {
		if(crConfig == null) {
			return null;
		}
		try {
			final JSONObject ret = new JSONObject();
			final Collection<DsState> dsList = DsState.getDsStates();
			for(DsState ds : dsList) {
				final JSONObject dsJo = new JSONObject();
				if(!ds.hasValue(AbstractState.IS_AVAILABLE_STR)) {
					continue;
				}
				dsJo.put(AbstractState.IS_AVAILABLE_STR, ds.isAvailable());
				dsJo.put("disabledLocations", ds.getDisabledLocations());
				ret.put(ds.getId(), dsJo);
			}
			return ret;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
		return null;
	}

	public static void init(final CacheWatcher cw, final PeerWatcher pw, final HealthDeterminer hd) {
		myCacheWatcher = cw;
		myPeerWatcher = pw;
		myHealthDeterminer = hd;
	}
}

