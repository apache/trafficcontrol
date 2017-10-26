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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.ArrayList;
import java.util.Deque;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.comcast.cdn.traffic_control.traffic_monitor.data.StatisticsLog;
import com.comcast.cdn.traffic_control.traffic_monitor.health.Event.EventType;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.util.string.Strings;

import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;
import com.comcast.cdn.traffic_control.traffic_monitor.data.DataSummary;

abstract public class AbstractState {
	private static final Logger LOGGER = Logger.getLogger(AbstractState.class);
	public static final String IS_AVAILABLE_STR = "isAvailable";
	public static final String IS_HEALTHY_STR = "isHealthy";
	private long historyTime = 5*60*1000;

	final String id;
	protected final StatisticsLog statisticsLog = new StatisticsLog();
	private Event lastEvent = null;

	protected AbstractState(final String id) {
		this.id = id;
	}

	public String getId() {
		return id;
	}

	protected void putDataPoints(final Map<String, String> statistics) {
		if (statistics == null) {
			return;
		}

		synchronized(this) {
			for(String key : statistics.keySet()) {
				putDataPoint(key, statistics.get(key));
			}
		}
	}

	public void putDataPoint(final String key, final String v) {
		statisticsLog.putDataPoint(key, v);
	}

	public Set<String> getStatisticsKeys() {
		return statisticsLog.getKeys();
	}

	protected Map<String, Deque<DataPoint>> getStats(final int hc, final String[] statList, final boolean wildcard, final boolean hidden) {
		return statisticsLog.filter(hc, statList, wildcard, hidden);
	}

	protected Deque<DataPoint> getDataPoints(final String key) {
		return statisticsLog.get(key);
	}

	public boolean hasValue(final String key) {
		return statisticsLog.hasValue(key);
	}

	public String getLastValue(final String key) {
		return statisticsLog.getLastValue(key);
	}

	public String getValue(final String key, final long targetIndex) {
		return statisticsLog.getValue(key, targetIndex);
	}

	public boolean getBool(final String key) {
		return statisticsLog.getBool(key);
	}

	public long getLong(final String key) {
		return statisticsLog.getLong(key);
	}

	public double getDouble(final String key) {
		return statisticsLog.getDouble(key);
	}

	public boolean isAvailable() {
		return getBool(IS_AVAILABLE_STR);
	}

	public Map<String, DataSummary> getSummary(final long startTime, final long endTime, final String[] stats2, final boolean wildcard, final boolean hidden) {
		final Map<String, Deque<DataPoint>> map = getStats(0, stats2, wildcard, hidden);
		final Map<String, DataSummary> retMap = new  HashMap<String, DataSummary>();
		final long checkPeriod = 5000;
		for(String key : map.keySet()) {
			final Deque<DataPoint> dpList = map.get(key);
			final DataSummary ds = new DataSummary();
			retMap.put(key, ds);
			for(DataPoint dp : dpList) {
				final int span = dp.getSpan();
				final long lastTime = statisticsLog.getTime(dp);
				for(int i = 0; i < span; i++) {
					final long t = lastTime - (checkPeriod * (span - (i+1)));
					if(startTime != 0 && t < startTime) { continue; }
					if(endTime != -1 && t > endTime) { break; }
					try {
						ds.accumulate(dp, t);
					} catch(Exception e) {
						LOGGER.debug(e);
					}
				}
			}
		}
		return retMap;
	}

	protected void prepareStatisticsForUpdate() {
		statisticsLog.prepareForUpdate(id, historyTime);
	}

	public void setHistoryTime(final long t) {
		if (t > 0) {
			historyTime = t;
		}
	}

	public String getStatusString() {
		final List<String> errors = new ArrayList<String>();
		final String error = getLastValue("error-string");
		final StringBuilder status = new StringBuilder();

		status.append(getLastValue("status"));

		if (error != null) {
			errors.add(error);
		} else {
			errors.add(isAvailable() ? "available" : "unavailable");
		}

		if (getBool("clearData")) {
			errors.add("monitoring disabled");
		}

		if (!errors.isEmpty()) {
			status.append(" - ");
			status.append(Strings.join(", ", errors));
		}

		return status.toString();
	}

	public boolean isError() {
		return getLastValue("error-string") != null;
	}

	public void setAvailable(final EventType type, final boolean isAvailable, final String error) {
		final boolean isHealthy = (error == null);
		boolean logChange = true;

		if (!isHealthy) {
			LOGGER.debug(String.format("Error on '%s': %s", id, error));
		}

		if (lastEvent != null && getBool(IS_HEALTHY_STR) == isHealthy && getBool(IS_AVAILABLE_STR) == isAvailable) {
			logChange = false;
		}

		putDataPoint(IS_AVAILABLE_STR, String.valueOf(isAvailable));
		putDataPoint(IS_HEALTHY_STR, String.valueOf(isHealthy));

		if (logChange) {
			lastEvent = Event.logStateChange(this.getId(), type, isAvailable, getStatusString());
		}
	}

	public JSONObject getStatsJson(final int hc, final String[] statList, final boolean wildcard, final boolean hidden) throws JSONException {
		final Map<String, Deque<DataPoint>> map = getStats(hc, statList, wildcard, hidden);
		final JSONObject ret = new JSONObject();
		for(String key : map.keySet()) {
			final JSONArray a = new JSONArray();
			ret.put(key, a);
			for(DataPoint dp : map.get(key)) {
				if (dp == null)
					continue;

				final JSONObject o = new JSONObject();
				/* Technically, this belongs in DataPoint, however, other logic
				 * relies on a null DataPoint value.  Because the JSON library
				 * discards a key with a null value, let's handle it here. -jse
				 */
				final String value = (dp.getValue() == null) ? "false" : dp.getValue();
				o.put("value", value);
				o.put("span", dp.getSpan());
				o.put("time", statisticsLog.getTime(dp));
				o.put("index", dp.getIndex());
				a.put(o);
			}
		}
		return ret;
	}

	public void addHiddenStats(Set<String> keys) {
		statisticsLog.addHiddenStats(keys);
	}

	public long getTime(final long index) {
		return statisticsLog.getTime(index);
	}
}
