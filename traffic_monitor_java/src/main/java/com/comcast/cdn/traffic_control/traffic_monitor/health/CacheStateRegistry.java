package com.comcast.cdn.traffic_control.traffic_monitor.health;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;

import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.ADMIN_DOWN;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.OFFLINE;

public class CacheStateRegistry extends StateRegistry<CacheState> {
	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private CacheStateRegistry() { }

	public static CacheStateRegistry getInstance() {
		return CacheStateRegistryHolder.REGISTRY;
	}

	private static class CacheStateRegistryHolder {
		private static final CacheStateRegistry REGISTRY = new CacheStateRegistry();
	}

	public CacheState update(final Cache cache) {
		CacheState cacheState = getOrCreate(cache.getHostname());
		cacheState.setCache(cache);
		return cacheState;
	}

	public int getCachesDownCount() {
		int count = 0;
		for (AbstractState state : states.values()) {
			if (state.isError()) {
				count++;
			}
		}
		return count;
	}

	public int getCachesAvailableCount() {
		int count = 0;
		for (AbstractState state : states.values()) {
			if (state.isAvailable()) {
				count++;
			}
		}
		return count;
	}

	public long getCachesBandwidthInKbps() {
		return getSumOfLongStatistic("kbps");
	}

	public long getCachesMaxBandwidthInKbps() {
		return getSumOfLongStatistic("maxKbps");
	}

	public String getStatusString(final String hostname) {
		AbstractState cacheState = states.get(hostname);
		if (cacheState == null || cacheState.isAvailable()) {
			return " ";
		}

		final String status = cacheState.getLastValue(HealthDeterminer.STATUS);

		if (status == null) {
			return "error";
		}

		HealthDeterminer.AdminStatus adminStatus = HealthDeterminer.AdminStatus.valueOf(status);

		if (adminStatus == ADMIN_DOWN || adminStatus == OFFLINE) {
			return "warning";
		}

		return "error";
	}

	@Override
	protected CacheState createState(final String id) {
		return new CacheState(id);
	}
}
