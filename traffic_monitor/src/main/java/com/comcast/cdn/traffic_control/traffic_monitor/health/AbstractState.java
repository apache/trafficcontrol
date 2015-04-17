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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.KeyValue;
import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;
import com.comcast.cdn.traffic_control.traffic_monitor.data.DataSummary;

abstract public class AbstractState implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	private static final Logger LOGGER = Logger.getLogger(AbstractState.class);
	public static final String IS_AVAILABLE_STR = "isAvailable";
	public static final String IS_HEALTHY_STR = "isHealthy";
	private long historyTime = 5*60*1000;

	final String id;
	private final Map<String,List<DataPoint>> stats = new HashMap<String,List<DataPoint>>();
	protected final Set<String> hiddenStats = new HashSet<String>();
	private int index;
	private Event lastEvent = null;

	protected AbstractState(final String id) {
		this.id = id;
	}
	public String getId() {
		return id;
	}

	protected void put(final Map<String, String> stati) {
		// TODO : loop all stats, set to null anything not in atsJson
		synchronized(this) {
			for(String key : stati.keySet()) {
				setDp(key, stati.get(key));
			}
//			setDp("queryTime", queryTimeStr, time);
		}
	}
	public void put(final String key, final String v) {
		setDp(key, v);
	}

	protected void setDp(final String key, final String v) {
		List<DataPoint> hist = stats.get(key);
		if(hist == null) {
			hist = new LinkedList<DataPoint>();
			synchronized(this) {
				stats.put(key, hist);
			}
		}
		DataPoint dp = null;
		if(!hist.isEmpty()) {
			dp = hist.get(hist.size()-1);
		}
		if(dp != null && dp.matches(v)) {
			dp.update(index);
			return; 
		}
		dp = new DataPoint(v, index);
		hist.add(dp);
	}

	public List<KeyValue> getModelList() {
		final List<KeyValue> al = new ArrayList<KeyValue>();
		final TreeMap<String,List<DataPoint>> st = new TreeMap<String,List<DataPoint>>(stats);
		for(String key : st.keySet()) {
			al.add(getKeyValue(key, this));
		}
		return al;
	}

	abstract protected KeyValue getKeyValue(String key, AbstractState state);

	public Map<String, List<DataPoint>> getStats() {
		return stats;
	}
	protected Map<String, List<DataPoint>> getStats(final int hc, final String[] statList, final boolean wildcard, final boolean hidden) {
		final Map<String, List<DataPoint>> ret = new HashMap<String,List<DataPoint>>();

		synchronized(this) {
			Set<String> statSet;

			if (statList == null) {
				statSet = stats.keySet();
			} else {
				if (!wildcard) {
					statSet = new HashSet<String>(Arrays.asList(statList));
				} else {
					statSet = new HashSet<String>();

					for (String key : stats.keySet()) {
						for (String stat : statList) {
							if (key.toLowerCase().contains(stat.toLowerCase())) {
								statSet.add(key);
								break;
							}
						}
					}
				}
			}

			for (String key : statSet) {
				final List<DataPoint> list;

				if (stats.containsKey(key) && (hidden || !hiddenStats.contains(key))) {
					list = new ArrayList<DataPoint>(stats.get(key));
				} else {
					continue;
				}

				if (hc == 0 || list.size() <= 1) {
					ret.put((String) key, list);
				} else {
					/*
					 * If fromIndex == toIndex, List.subList() will return an empty list.
					 * The only way they will be equal is if the list is empty or
					 * has a single item, which is handled above.
					 */

					final int toIndex = list.size();
					final int fromIndex = Math.max(0, toIndex - hc);

					ret.put((String) key, list.subList(fromIndex, toIndex));
				}
			}
		}

		return ret;
	}
	public List<DataPoint> getDataPoints(final String key) {
		return stats.get(key);
	}
	public boolean hasValue(final String key) {
		final List<DataPoint> dps = getDataPoints(key);
		if(dps == null || dps.isEmpty()) { return false; }
		return true;
	}
	public boolean hasValue(final String key, final int index) {
		final List<DataPoint> dps = getDataPoints(key);
		if(dps == null || dps.isEmpty()) { return false; }
		if(dps.get(dps.size()-1).getIndex() >= index) {
			return true;
		}
		return false;
	}
	public String getLastValue(final String key) {
		final List<DataPoint> dps = getDataPoints(key);
		if(dps == null || dps.isEmpty()) { return null; }
		return dps.get(dps.size()-1).getValue();
	}
	public String getValue(final String key, final long targetIndex) {
		final List<DataPoint> dps = getDataPoints(key);
		if(dps == null || dps.isEmpty()) { return null; }
		for(int i = dps.size()-1; i >= 0; i--) {
			final long dpIndex = dps.get(i).getIndex();
			if(targetIndex > (dpIndex)) { return null; }
			final long span = dps.get(i).getSpan();
			if(targetIndex <= (dpIndex-span)) { continue; }
			final String v = dps.get(i).getValue();
			return v;
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
	public boolean isAvailable() {
		return getBool(IS_AVAILABLE_STR);
	}
	public Map<String, DataSummary> getSummary(final long startTime, final long endTime, final String[] stats2, final boolean wildcard, final boolean hidden) {
		final Map<String, List<DataPoint>> map = getStats(0, stats2, wildcard, hidden);
		final Map<String, DataSummary> retMap = new  HashMap<String, DataSummary>();
		final long checkPeriod = 5000;
		for(String key : map.keySet()) {
			final List<DataPoint> dpList = map.get(key);
			final DataSummary ds = new DataSummary();
			retMap.put(key, ds);
			for(DataPoint dp : dpList) {
				final int span = dp.getSpan();
				final long lastTime = getTime(dp);
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
	protected long getTime(final DataPoint dp) {
		return getTime(dp.getIndex());
	}
	protected long getTime(final long index) {
		synchronized(times) {
			for(int i = 0; i < indexes.size(); i++) {
				if(indexes.get(i).longValue() == index) {
					return times.get(i).longValue();
				}
			}
		}
		return 0;
	}
	private final List<Long> times = new LinkedList<Long>();
	private final List<Long> indexes = new LinkedList<Long>();
	protected void startUpdate() {
		bringUpToDate();

		final long time = System.currentTimeMillis();
		final long removeTime = time - historyTime;
		int removeCount = 0;
		synchronized(times) {
			index++;
			times.add(new Long(time));
			indexes.add(new Long(index));
			while(times.get(0).longValue() < removeTime) {
				removeCount++;
				times.remove(0);
			}
			if(removeCount == 0) {
				return;
			}
			for(int i = 0; i < removeCount; i++) {
				indexes.remove(0);
			}
		}
		final long baseIndex = indexes.get(0); 
		for(String key : stats.keySet()) {
			final List<DataPoint> l = stats.get(key);
			if(l.isEmpty()) {
				LOGGER.warn("list empty for "+key + " - "+this.id);
				continue;
			}
			while(l.get(0).getIndex() < baseIndex) {
				if(l.size() == 1) {
					LOGGER.warn(String.format("%s - %s: index %d < baseIndex %d", key, this.id, l.get(0).getIndex(), baseIndex));
					break;
				}
				synchronized(this) {
					l.remove(0);
				}
				if(l.isEmpty()) {
					LOGGER.warn("list empty for "+key + " - "+this.id);
					break;
				}
			}
		}
	}

	private void bringUpToDate() {
		for(String key : stats.keySet()) {
			final List<DataPoint> l = stats.get(key);
			final DataPoint dp = l.get(l.size()-1);
			if(dp.getIndex() != index) {
				put(key, null);
			}
		}
	}
	public void setHistoryTime(final long t) {
		if(t>0) {
			historyTime = t;
		}
	}
	public String getStatusString() {
		String error = getLastValue("error-string");
		if(error == null) {
			error = isAvailable()? "available" : "";
		}
		if(getBool("clearData")) { error = "No query"; }
		return getLastValue("status") + " - " + error;
	}

	public boolean isError() {
		if(getLastValue("error-string") != null) {
			return true;
		}
		return false;
	}

	public void setAvailable(final boolean isAvailable, final String error) {
		final boolean isHealthy = (error == null);
		boolean logChange = true;

		if (!isHealthy) {
			LOGGER.debug(String.format("Error on '%s': %s", id, error));
		}

		if (lastEvent != null && getBool(IS_HEALTHY_STR) == isHealthy && getBool(IS_AVAILABLE_STR) == isAvailable) {
			logChange = false;
		}

		put(IS_AVAILABLE_STR, String.valueOf(isAvailable));
		put(IS_HEALTHY_STR, String.valueOf(isHealthy));

		if (logChange) {
			lastEvent = Event.logStateChange(this.getId(), isAvailable, getStatusString());
		}
	}

	public JSONObject getStatsJson(final int hc, final String[] statList, final boolean wildcard, final boolean hidden) throws JSONException {
		final Map<String, List<DataPoint>> map = getStats(hc, statList, wildcard, hidden);
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
				final String value = (dp.getValue() == null) ? "" : dp.getValue();
				o.put("value", value);
				o.put("span", dp.getSpan());
				o.put("time", getTime(dp));
				o.put("index", dp.getIndex());
				a.put(o);
			}
		}
		return ret;
	}

	public int getCurrentIndex() {
		return index;
	}

	public static List<? extends AbstractState> getStates(final Class<? extends AbstractState> c) {
		return null;
	}

}
