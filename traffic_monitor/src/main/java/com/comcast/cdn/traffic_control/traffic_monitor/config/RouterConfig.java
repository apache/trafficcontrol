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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.TmListener;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Network;

public class RouterConfig {
	private static final Logger LOGGER = Logger.getLogger(RouterConfig.class);

	final private List<Cache> cacheList;
	final private Map<String, Peer> peerMap;
	final private JSONObject dsList;
	final private Collection<JSONObject> contentRouters;

	public RouterConfig(final JSONObject o, final HealthDeterminer myHealthDeterminer) throws JSONException {
		final ArrayList<Cache> al = new ArrayList<Cache>();
		final Map<String, Peer> pm = new HashMap<String, Peer>();
		final JSONObject caches = o.optJSONObject("contentServers");

		LOGGER.info("Processing new CrConfig");

		for(String id : JSONObject.getNames(caches)) {
			try {
				final JSONObject cjo = caches.getJSONObject(id);
//				if(!myHealthDeterminer.shouldMonitor(cjo)) { // now assumed with CrConfig instead of dataservers
//					continue;
//				}
				final Cache c = new Cache(id,cjo, this);
				myHealthDeterminer.setControls(c); // set the controls
				c.setCacheState((CacheState) CacheStateRegistry.getInstance().get(c.getHostname()));
				al.add(c);
			} catch (JSONException e) {
				LOGGER.warn("handleTmJson: ",e);
			}
		}
		cacheList = al;

		if (o.has("monitors")) {
			final JSONObject peers = o.optJSONObject("monitors");

			for (String id : JSONObject.getNames(peers)) {
				final JSONObject pjo = peers.getJSONObject(id);

				final String peerStatus = pjo.optString(HealthDeterminer.STATUS);
				final String peerIpAddress = pjo.getString("ip");

				if (Network.isIpAddressLocal(peerIpAddress)) {
					LOGGER.warn("Skipping monitor " + id + "; IP address " + peerIpAddress + " is local");
					continue;
				} else if (Network.isLocalName(pjo.getString("fqdn"))) {
					LOGGER.warn("Skipping monitor " + id + "; fqdn " + pjo.getString("fqdn") + " is the local fully qualified name");
					continue;
				} else if (Network.isLocalName(id)) {
					LOGGER.warn("Skipping monitor " + id + "; short name " + id + " is the local hostname");
					continue;
				} else if (peerStatus != null && peerStatus.equals("ONLINE")) {
					final Peer peer = new Peer(id, pjo);
					pm.put(peer.getId(), peer);
				}
			}

			PeerState.removeAllBut(pm.keySet());
		}

		peerMap = pm;

		final Collection<JSONObject> list = new ArrayList<JSONObject>();
		final JSONObject crsJo = o.getJSONObject("contentRouters");
		for(String key : JSONObject.getNames(crsJo)) {
			final JSONObject crJo = crsJo.getJSONObject(key);
			crJo.put("id", key);
			list.add(crJo);
		}
		contentRouters = list;
		dsList = o.optJSONObject("deliveryServices");
//		final ArrayList<DeliveryService> dsAl = new ArrayList<DeliveryService>();
//		for(String id : JSONObject.getNames(dsList)) {
//			try {
//				final JSONObject dsjo = dsList.getJSONObject(id);
//				final DeliveryService c = new DeliveryService(id,dsjo);
//				dsAl.add(c);
//			} catch (JSONException e) {
//				LOGGER.warn("handleTmJson: ",e);
//			}
//		}
//		DeliveryService.setList(dsAl);
	}
	public List<Cache> getCacheList() {
		return cacheList;
	}
	public Map<String, Peer> getPeerMap() {
		return peerMap;
	}
	public Collection<JSONObject> getContentRouters() {
		return contentRouters;
	}
	public JSONObject getDsList() {
		return dsList;
	}
	private static RouterConfig crConfig;
	public static TmListener getTmListener(final HealthDeterminer myHealthDeterminer) {
		return new TmListener() {
			@Override public void handleCrConfig(final JSONObject o) {
				try {
					crConfig = new RouterConfig(o, myHealthDeterminer);
				} catch (JSONException e) {
					if(LOGGER.isDebugEnabled()) {
						try {
							LOGGER.debug(o.toString(2));
						} catch (JSONException e1) {
							LOGGER.warn(e1,e1);
						}
						LOGGER.debug(e,e);
					} else {
						LOGGER.warn(e);
					}
				}

			}

		};
	}
	public static RouterConfig getCrConfig() {
		return crConfig;
	}
}
