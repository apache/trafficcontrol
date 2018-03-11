package com.comcast.cdn.traffic_control.traffic_monitor.data;

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


import org.apache.log4j.Logger;

import java.util.Arrays;
import java.util.Deque;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class StatisticsLog {
	private static final Logger LOGGER = Logger.getLogger(StatisticsLog.class);
	private final Map<String,Deque<DataPoint>> data = new HashMap<String,Deque<DataPoint>>();
	protected final Set<String> hiddenKeys = new HashSet<String>();
	private final List<Long> times = new LinkedList<Long>();
	private final List<Long> indexes = new LinkedList<Long>();
	private int index;

	public Deque<DataPoint> get(final String key) {
		return data.get(key);
	}

	public void putDataPoint(final String key, final String value) {
		Deque<DataPoint> statistics = data.get(key);

		if (statistics == null) {
			statistics = new LinkedList<DataPoint>();
			synchronized(data) {
				data.put(key, statistics);
			}
		}

		DataPoint dataPoint = getLastDataPoint(key);

		if (dataPoint != null && dataPoint.matches(value)) {
			dataPoint.update(index);
			return;
		}

		dataPoint = new DataPoint(value, index);
		statistics.addLast(dataPoint);
	}

	private DataPoint getLastDataPoint(final String key) {
		if (!hasValue(key)) {
			return null;
		}

		return data.get(key).getLast();
	}

	public Set<String> getKeys() {
		return data.keySet();
	}

	public boolean hasValue(final String key) {
		return data.containsKey(key) && (!data.get(key).isEmpty());
	}

	public String getLastValue(final String key) {
		final DataPoint dataPoint = getLastDataPoint(key);
		return (dataPoint != null) ? dataPoint.getValue() : null;
	}

	public String getValue(final String key, final long targetIndex) {
		if (!hasValue(key)) {
			return null;
		}

		Iterator<DataPoint> dataPoints = get(key).descendingIterator();

		while (dataPoints.hasNext()) {
			final DataPoint dataPoint = dataPoints.next();
			if (targetIndex > dataPoint.getIndex()) {
				return null;
			}

			if (targetIndex <= (dataPoint.getIndex() - dataPoint.getSpan())) {
				continue;
			}

			return dataPoint.getValue();
		}

		return null;
	}

	public boolean getBool(final String key) {
		try {
			return Boolean.parseBoolean(getLastValue(key));
		} catch (Exception e) {
			return true;
		}
	}

	public long getLong(final String key) {
		try {
			return Long.parseLong(getLastValue(key));
		} catch (Exception e) {
			return 0;
		}
	}

	public double getDouble(final String key) {
		try {
			return Double.parseDouble(getLastValue(key));
		} catch (Exception e) {
			return 0;
		}
	}

	private Set<String> filterKeys(final String[] statList, final boolean wildcard) {
		Set<String> statisticsKeys;

		if (statList == null) {
			return getKeys();
		}

		if (!wildcard) {
			return new HashSet<String>(Arrays.asList(statList));
		}

		statisticsKeys = new HashSet<String>();

		for (String key : getKeys()) {
			for (String stat : statList) {
				if (key.toLowerCase().contains(stat.toLowerCase())) {
					statisticsKeys.add(key);
					break;
				}
			}
		}

		return statisticsKeys;
	}

	public Map<String, Deque<DataPoint>> filter(final int maxItems, final String[] statList, final boolean wildcard, final boolean allowHidden) {
		final Map<String, Deque<DataPoint>> filteredStatistics = new HashMap<String,Deque<DataPoint>>();

		synchronized(data) {
			Set<String> statisticsKeys = filterKeys(statList, wildcard);

			for (String key : statisticsKeys) {

				if (!data.containsKey(key) || (!allowHidden && hiddenKeys.contains(key))) {
					continue;
				}

				final LinkedList<DataPoint> statistics = (LinkedList<DataPoint>) data.get(key);

				if (maxItems == 0 || statistics.size() <= 1) {
					filteredStatistics.put(key, statistics);
				} else {
					/*
					 * If fromIndex == toIndex, List.subList() will return an empty list.
					 * The only way they will be equal is if the list is empty or
					 * has a single item, which is handled above.
					 */

					final int toIndex = statistics.size();
					final int fromIndex = Math.max(0, toIndex - maxItems);

					filteredStatistics.put(key, new LinkedList<DataPoint>(statistics.subList(fromIndex, toIndex)));
				}
			}
		}

		return filteredStatistics;
	}

	public void addHiddenStats(Set<String> keys) {
		hiddenKeys.addAll(keys);
	}

	public long getTime(final DataPoint dataPoint) {
		return getTime(dataPoint.getIndex());
	}

	public long getTime(final long targetIndex) {
		synchronized(times) {
			for (long index : indexes) {
				if (index == targetIndex) {
					return times.get(indexes.indexOf(index));
				}
			}
		}

		return 0;
	}

	public void prepareForUpdate(final String stateId, final long historyTime) {

		synchronized(times) {
			addNullDataForIndex(index);
			index++;
			final long time = System.currentTimeMillis();
			final long removeTime = time - historyTime;
			int removeCount = 0;

			times.add(time);
			indexes.add(new Long(index));

			while (times.get(0) < removeTime) {
				removeCount++;
				times.remove(0);
			}

			if (removeCount == 0) {
				return;
			}

			for (int i = 0; i < removeCount; i++) {
				indexes.remove(0);
			}

			removeOldest(stateId);
		}
	}

	private void addNullDataForIndex(final long index) {
		synchronized (data) {
			for(String key : data.keySet()) {
				DataPoint lastDataPoint = getLastDataPoint(key);
				if (lastDataPoint == null || lastDataPoint.getIndex() != index) {
					putDataPoint(key, null);
				}
			}
		}
	}

	private void removeOldest(final String stateId) {
		final long oldestIndex = indexes.get(0);

		for(String key : data.keySet()) {
			final Deque<DataPoint> dataPoints = get(key);

			if (dataPoints.isEmpty()) {
				LOGGER.warn("list empty for " + key + " - " + stateId);
				continue;
			}


			while (dataPoints.getFirst().getIndex() < oldestIndex) {
				if (dataPoints.size() == 1) {
					LOGGER.warn(String.format("%s - %s: index %d < baseIndex %d", key, stateId, dataPoints.getFirst().getIndex(), oldestIndex));
					break;
				}

				dataPoints.remove();

				if (dataPoints.isEmpty()) {
					LOGGER.warn(String.format("list empty for %s - %s", key, stateId));
				}
			}

		}
	}
}
