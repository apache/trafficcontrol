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
import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONObject;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Deque;
import java.util.Iterator;
import java.util.List;

public class DeliveryServiceStateRegistry extends StateRegistry<DsState> {
	private static final Logger LOGGER = Logger.getLogger(DeliveryServiceStateRegistry.class);

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private DeliveryServiceStateRegistry() { }

	public static DeliveryServiceStateRegistry getInstance() {
		return DeliveryServiceStateRegistryHolder.REGISTRY;
	}

	private static class DeliveryServiceStateRegistryHolder {
		private static final DeliveryServiceStateRegistry REGISTRY = new DeliveryServiceStateRegistry();
	}

	public void completeUpdateAll(final HealthDeterminer myHealthDeterminer, final JSONObject dsList, final long lenientTime) {
		for (CacheState cacheState : CacheStateRegistry.getInstance().getAll()) {
			if (cacheState.getCache() != null && cacheState.getCache().hasDeliveryServices()) {
				updateStates(cacheState, lenientTime);
			}
		}

		final Collection<String> toRemove = new ArrayList<String>();
		toRemove.addAll(states.keySet());

		if (dsList != null) {
			for (String dsId : JSONObject.getNames(dsList)) {
				toRemove.remove(dsId);
				try {
					final DsState dss = (DsState) getOrCreate(dsId);
					dss.completeRound(myHealthDeterminer.getDsControls(dss.getId()));
				} catch (Exception e) {
					LOGGER.warn(e, e);
				}
			}
		}

		for(String id : toRemove) {
			states.remove(id);
		}
	}

	private void updateStates(final CacheState cacheState, final long lenientTime) {
		final Cache cache = cacheState.getCache();

		for(String deliveryServiceId : cache.getDeliveryServiceIds()) {
			try {
				final List<String> fqdns = cache.getFqdns(deliveryServiceId);
				final DsState deliveryServiceState = (DsState) getOrCreate(deliveryServiceId);

				// Don't count the cache as reporting unless there were no errors and stats were read
				boolean error = false;
				boolean foundStats = false;

				for(String fqdn : fqdns) {
					final String propBase = "ats.plugin.remap_stats."+fqdn;
					final DsStati stati = createStati(propBase, cacheState, lenientTime, deliveryServiceId);

					deliveryServiceState.accumulate(stati, cache.getLocation(), cacheState);

					if (stati != null) {
						foundStats = true;

						if (stati.error) {
							error = true;
						}
					}
				}

				// Update cache counters
				deliveryServiceState.addCacheConfigured();

				if (cacheState.isAvailable()) {
					deliveryServiceState.addCacheAvailable();
				}

				if (foundStats && !error) {
					deliveryServiceState.addCacheReporting();
				}
			} catch(Exception e) {
				LOGGER.warn(e,e);
			}
		}
	}

	public void startUpdateAll() {
		synchronized(states) {
			for(AbstractState ds :states.values()) {
				ds.prepareStatisticsForUpdate();
			}
		}
	}

	@Override
	protected DsState createState(final String deliveryServiceId) {
		return new DsState(deliveryServiceId);
	}

	private DsStati createStati(final String propBase, final CacheState cacheState, final long leniency, final String dsId) {
		DsStati dsStati;

		synchronized (cacheState) {
			final Deque<DataPoint> dataPoints = cacheState.getDataPoints(propBase + ".out_bytes");

			if (dataPoints == null) {
				return null;
			}

			long lastIndex = dataPoints.getLast().getIndex();
			lastIndex = getLastGoodIndex(dataPoints, lastIndex);

			if (lastIndex < 0) {
				return null;
			}

			final long time = cacheState.getTime(lastIndex);

			if (time < leniency) {
				return null;
			}

			dsStati  = new DsStati(propBase, cacheState, lastIndex, dsId);

			final long prevIndex = getLastGoodIndex(dataPoints, lastIndex-1);

			if (prevIndex >= 0) {
				final DsStati priorDsStati = new DsStati(propBase, cacheState, prevIndex, dsId);

				if (!dsStati.calculateKbps(priorDsStati)) {
					if (LOGGER.isInfoEnabled()) {
						printDps(dataPoints, propBase);
					}
				}
			}
		}

		return dsStati;
	}

	public long getLastGoodIndex(final Deque<DataPoint> dataPoints, final long targetIndex) {
		if (targetIndex < 0) {
			return -1;
		}

		Iterator<DataPoint> dataPointIterator = dataPoints.descendingIterator();

		while (dataPointIterator.hasNext()) {
			DataPoint dataPoint = dataPointIterator.next();
			if (dataPoint.getValue() == null) {
				continue;
			}

			final long index = dataPoint.getIndex();
			final long span = dataPoint.getSpan();

			if (targetIndex <= (index-span)) {
				continue;
			}

			if (targetIndex < index) {
				return targetIndex;
			}

			return index;
		}

		return -1;
	}

	public boolean printDps(final Deque<DataPoint> dataPoints, final String id) {
		LOGGER.warn(id + ":");

		Iterator<DataPoint> dataPointIterator = dataPoints.descendingIterator();
		while (dataPointIterator.hasNext()) {
			DataPoint dataPoint = dataPointIterator.next();
			LOGGER.warn(String.format("\tindex: %d, span: %d, value: %s", dataPoint.getIndex(), dataPoint.getSpan(), dataPoint.getValue()));
		}

		return false;
	}

}
