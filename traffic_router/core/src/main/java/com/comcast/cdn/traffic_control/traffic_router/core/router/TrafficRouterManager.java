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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.io.IOException;
import java.net.UnknownHostException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import org.apache.commons.pool.ObjectPool;
import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.TrafficRouterException;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.NameServer;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;

public class TrafficRouterManager {
	private static final Logger LOGGER = Logger.getLogger(TrafficRouterManager.class);

	private JSONObject state;
	private TrafficRouter trafficRouter;
	private GeolocationService geolocationService;
	private GeolocationService geolocationService6;
	private ObjectPool hashFunctionPool;
	private StatTracker statTracker;
	private static final Map<String, Long> timeTracker = new ConcurrentHashMap<String, Long>();
	private NameServer nameServer;
	private TrafficOpsUtils trafficOpsUtils;

	public NameServer getNameServer() {
		return nameServer;
	}

	public static Map<String, Long> getTimeTracker() {
		return timeTracker;
	}

	public void trackEvent(final String event) {
		timeTracker.put(event, System.currentTimeMillis());
	}

	public void setNameServer(final NameServer nameServer) {
		this.nameServer = nameServer;
	}

	public boolean setState(final JSONObject jsonObject) throws UnknownHostException {
		trackEvent("lastCacheStateCheck");
		if(jsonObject == null) {
			return false;
		}
		trackEvent("lastCacheStateChange");
		synchronized(this) {
			this.state = jsonObject;
			if(trafficRouter != null) {
				trafficRouter.setState(state);
				return true;
			}
			return true;
		}
	}

	public TrafficRouter getTrafficRouter() {
		return trafficRouter;
	}

	public void setCacheRegister(final CacheRegister cacheRegister) throws IOException, JSONException, TrafficRouterException {
		trackEvent("lastConfigCheck");
		if(cacheRegister == null) {
			return;
		}

		final TrafficRouter tr = new TrafficRouter(cacheRegister, 
				geolocationService, 
				geolocationService6, 
				hashFunctionPool, 
				statTracker,
				trafficOpsUtils);
		synchronized(this) {
			if(state != null) {
				try {
					tr.setState(state);
				} catch (UnknownHostException e) {
					LOGGER.warn(e,e);
				}
			}
			this.trafficRouter = tr;
		}
		trackEvent("lastConfigChange");
	}
	public void setGeolocationService(final GeolocationService geolocationService) {
		this.geolocationService = geolocationService;
	}
	public void setGeolocationService6(final GeolocationService geolocationService) {
		this.geolocationService6 = geolocationService;
	}
	public void setHashFunctionPool(final ObjectPool hashFunctionPool) {
		this.hashFunctionPool = hashFunctionPool;
	}
	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}
}
