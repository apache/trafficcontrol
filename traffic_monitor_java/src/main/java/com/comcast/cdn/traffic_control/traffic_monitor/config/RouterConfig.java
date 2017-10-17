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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.TmListener;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Network;

public class RouterConfig {
	private static final Logger LOGGER = Logger.getLogger(RouterConfig.class);

	private static RouterConfig crConfig;

	final private List<Cache> cacheList = new ArrayList<Cache>();
	final private Map<String, Peer> peerMap = new HashMap<String, Peer>();
	final private JSONObject dsList;

	public RouterConfig(final JSONObject crConfigJson, final HealthDeterminer healthDeterminer) throws JSONException {
		final JSONObject cachesJson = crConfigJson.optJSONObject("contentServers");

		LOGGER.info("Processing new CrConfig");

		for (String id : JSONObject.getNames(cachesJson)) {
			try {
				final Cache cache = new Cache(id, cachesJson.getJSONObject(id));
				healthDeterminer.setControls(cache);
				cache.setCacheState(CacheStateRegistry.getInstance().get(cache.getHostname()));
				cacheList.add(cache);
			} catch (JSONException e) {
				LOGGER.warn("Failed processing json for cache " + id + ":", e);
			}
		}

		if (crConfigJson.has("monitors")) {
			final JSONObject peers = crConfigJson.optJSONObject("monitors");

			for (String id : JSONObject.getNames(peers)) {
				final Peer peer = new Peer(id, peers.optJSONObject(id));

				if (Network.isIpAddressLocal(peer.getIpAddress())) {
					LOGGER.debug("Skipping monitor " + id + "; IP address " + peer.getIpAddress() + " is local");
					continue;
				}

				if (Network.isLocalName(peer.getFqdn())) {
					LOGGER.debug("Skipping monitor " + id + "; fqdn " + peer.getFqdn() + " is the local fully qualified name");
					continue;
				}

				if (Network.isLocalName(id)) {
					LOGGER.debug("Skipping monitor " + id + "; short name " + id + " is the local hostname");
					continue;
				}

				if ("ONLINE".equals(peer.getStatus())) {
					peerMap.put(peer.getId(), peer);
				}
			}

			PeerState.removeAllBut(peerMap.keySet());
		}

		dsList = crConfigJson.optJSONObject("deliveryServices");
	}

	public List<Cache> getCacheList() {
		return cacheList;
	}

	public Map<String, Peer> getPeerMap() {
		return peerMap;
	}

	public JSONObject getDsList() {
		return dsList;
	}

	public static TmListener getTmListener(final HealthDeterminer healthDeterminer) {
		return new TmListener() {
			@Override
			public void handleCrConfig(final JSONObject crConfigJson) {
				try {
					crConfig = new RouterConfig(crConfigJson, healthDeterminer);
				} catch (Exception e) {
					LOGGER.warn("Failed Processing CrConfig json",e);
				}
			}
		};
	}

	public static RouterConfig getCrConfig() {
		return crConfig;
	}
}
