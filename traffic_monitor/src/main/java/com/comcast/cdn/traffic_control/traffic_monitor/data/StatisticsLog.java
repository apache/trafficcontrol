package com.comcast.cdn.traffic_control.traffic_monitor.data;

import org.apache.log4j.Logger;

import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class StatisticsLog {
	private static final Logger LOGGER = Logger.getLogger(StatisticsLog.class);
	private final Map<String,List<DataPoint>> data = new HashMap<String,List<DataPoint>>();
	protected final Set<String> hiddenKeys = new HashSet<String>();
	private final List<Long> times = new LinkedList<Long>();
	private final List<Long> indexes = new LinkedList<Long>();
	private int index;

	public List<DataPoint> get(final String key) {
		return data.get(key);
	}

	public void putDataPoint(final String key, final String value) {
		List<DataPoint> statistics = data.get(key);

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
		statistics.add(dataPoint);
	}

	public DataPoint getLastDataPoint(final String key) {
		if (!hasValue(key)) {
			return null;
		}

		return data.get(key).get(data.size() - 1);
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

		List<DataPoint> dataPoints = get(key);

		for (int i = dataPoints.size()-1; i >= 0; i--) {
			final DataPoint dataPoint = dataPoints.get(i);

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

	protected Set<String> filterKeys(final String[] statList, final boolean wildcard) {
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

	public Map<String, List<DataPoint>> filter(final int hc, final String[] statList, final boolean wildcard, final boolean allowHidden) {
		final Map<String, List<DataPoint>> filteredStatistics = new HashMap<String,List<DataPoint>>();

		synchronized(data) {
			Set<String> statisticsKeys = filterKeys(statList, wildcard);

			for (String key : statisticsKeys) {

				if (!data.containsKey(key) || (!allowHidden && hiddenKeys.contains(key))) {
					continue;
				}

				final List<DataPoint> statistics = data.get(key);

				if (hc == 0 || statistics.size() <= 1) {
					filteredStatistics.put(key, statistics);
				} else {
					/*
					 * If fromIndex == toIndex, List.subList() will return an empty list.
					 * The only way they will be equal is if the list is empty or
					 * has a single item, which is handled above.
					 */

					final int toIndex = statistics.size();
					final int fromIndex = Math.max(0, toIndex - hc);

					filteredStatistics.put(key, statistics.subList(fromIndex, toIndex));
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

	public void prepareForUpdate(final long index, final long historyTime) {
		final long time = System.currentTimeMillis();
		final long removeTime = time - historyTime;
		int removeCount = 0;

		times.add(time);
		indexes.add(index);

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
	}

	public void clearNonMatchingDataPoints(final long index) {
		for(String key : data.keySet()) {
			if (getLastDataPoint(key).getIndex() != index) {
				putDataPoint(key, null);
			}
		}
	}

	public void removeOldest(final String stateId) {
		final long baseIndex = indexes.get(0);

		for(String key : data.keySet()) {
			final List<DataPoint> dataPoints = get(key);

			if (dataPoints.isEmpty()) {
				LOGGER.warn("list empty for " + key + " - " + stateId);
				continue;
			}

			while (dataPoints.get(0).getIndex() < baseIndex) {
				if (dataPoints.size() == 1) {
					LOGGER.warn(String.format("%s - %s: index %d < baseIndex %d", key, stateId, dataPoints.get(0).getIndex(), baseIndex));
					break;
				}

				dataPoints.remove(0);

				if (dataPoints.isEmpty()) {
					LOGGER.warn("list empty for " + key + " - " + stateId);
					break;
				}
			}
		}
	}

	public void prepareForUpdate(final String id, final long historyTime) {
		clearNonMatchingDataPoints(index);

		synchronized(data) {
			index++;
			prepareForUpdate(index, historyTime);
		}

		removeOldest(id);
	}
}