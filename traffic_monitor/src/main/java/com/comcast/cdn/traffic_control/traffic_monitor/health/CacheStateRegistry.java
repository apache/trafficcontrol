package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import org.apache.log4j.Logger;

import java.util.ArrayList;
import java.util.List;

import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.ADMIN_DOWN;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.OFFLINE;

public class CacheStateRegistry extends StateRegistry {
	private static final Logger LOGGER = Logger.getLogger(CacheStateRegistry.class);
	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private CacheStateRegistry() { }

	public static CacheStateRegistry getInstance() {
		return CacheStateRegistryHolder.REGISTRY;
	}

	private static class CacheStateRegistryHolder {
		private static final CacheStateRegistry REGISTRY = new CacheStateRegistry();
	}

	public CacheState getOrCreate(final Cache cache) {
		if (cache == null) {
			LOGGER.warn("Tried to create cache state from a null value cache");
			return null;
		}
		CacheState cacheState = (CacheState) getOrCreate(cache.getHostname());

		if (cacheState == null) {
			LOGGER.warn("getOrCreate returned a null CacheState");
		} else {
			cacheState.setCache(cache);
		}

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

 	public void removeAllBut(final List<CacheState> retList) {
		final List<String> hostnames = new ArrayList<String>();

		for (CacheState cs : retList) {
			hostnames.add(cs.getId());
		}

		synchronized (states) {
			for (String key : states.keySet()) {
				if (!hostnames.contains(key)) {
					states.remove(key);
				}
			}
		}
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
	protected AbstractState createState(final String id) {
		return new CacheState(id);
	}
}
