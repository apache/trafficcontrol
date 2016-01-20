package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.ADMIN_DOWN;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.OFFLINE;

public class CacheStateRegistry extends StateRegistry implements Serializable {
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
		CacheState cacheState = (CacheState) getOrCreate(cache.getHostname());
		cacheState.setCache(cache);
		return cacheState;
	}

	public int getCachesDownCount() {
		int count = 0;
		// Do we really allow calling code to store nulls in our registry???
		for(AbstractState state : states.values()) {
			if (state != null && state.isError()) {
				count++;
			}
		}
		return count;
	}

	public int getCachesAvailableCount() {
		int count = 0;
		// Do we really allow calling code to store nulls in our registry???
		for(AbstractState state : states.values()) {
			if (state != null && state.isAvailable()) {
				count++;
			}
		}
		return count;
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
