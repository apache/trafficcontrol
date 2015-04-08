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

import java.text.DecimalFormat;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.KeyValue;
import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;

public class DsState extends AbstractState {
	private static final Logger LOGGER = Logger.getLogger(DsState.class);
	private static final long serialVersionUID = 1L;
	private static Map<String, DsState> states = new HashMap<String, DsState>();

	private DsStati currentDtati;
	private int cachesConfigured = 0;
	private int cachesAvailable = 0;
	private int cachesReporting = 0;
	final private Map<String, EmbeddedStati> locs = new HashMap<String, EmbeddedStati>();
	final private Map<String, EmbeddedStati> cacheStatiMap = new HashMap<String, EmbeddedStati>();

	public DsState(final String id) {
		super(id);
	}

	public static DsState getOrCreate(final String id) {
		synchronized(states) {
			DsState as = states.get(id);
			if(as == null) {
				as = new DsState(id);
				states.put(id, as);
			}
			return as;
		}
	}

	public static DsState getState(final String id) {
		synchronized(states) {
			return states.get(id);
		}
	}

	public void accumulate(final DsStati stati, final String location, final CacheState state) {
		if (stati == null) {
			return;
		}

		if (currentDtati == null) {
			currentDtati = stati;
		} else {
			currentDtati.accumulate(stati);
		}

		EmbeddedStati loc = locs.get(location);

		if (loc == null) {
			loc = new EmbeddedStati("location", location);
			locs.put(location,loc);
		}

		loc.accumulate(stati);

		EmbeddedStati cacheStati = cacheStatiMap.get(state.id);

		if (cacheStati == null) {
			cacheStati = new EmbeddedStati("cache", state.id);
			cacheStatiMap.put(state.id, cacheStati);
		}

		cacheStati.accumulate(stati);
	}

	public boolean completeRound(final JSONObject dsControls) {
		if (currentDtati != null && currentDtati.out_bytes != 0) {
			put(currentDtati.getStati("total"));
			currentDtati = null;
		}

		setDp("caches-configured", String.valueOf(cachesConfigured));
		setDp("caches-available", String.valueOf(cachesAvailable));
		setDp("caches-reporting", String.valueOf(cachesReporting));

		cachesConfigured = 0;
		cachesAvailable = 0;
		cachesReporting = 0;

		HealthDeterminer.setIsAvailable(this, dsControls);

		final StringBuilder sb = new StringBuilder();

		for (String locId : locs.keySet()) {
			final EmbeddedStati loc = locs.get(locId);
			final Map<String, String> stati = loc.completeRound();

			if (stati == null) {
				continue;
			}

			put(stati);

			if (!HealthDeterminer.setIsAvailable(this, loc, dsControls)) {
				sb.append("\"").append(locId).append("\", ");
			}
		}

		put("disabledLocations", sb.toString());

		for (String cacheId : cacheStatiMap.keySet()) {
			final EmbeddedStati cacheStat = cacheStatiMap.get(cacheId);
			final Map<String, String> stati = cacheStat.completeRound();

			if (stati == null) {
				continue;
			}

			hiddenStats.addAll(stati.keySet());

			put(stati);
		}

		return true;
	}
	public static void completeAll(final List<CacheState> crStates, final HealthDeterminer myHealthDeterminer, 
			final JSONObject dsList, final long lenientTime) {
		// loop all states
		for(CacheState crstate : crStates) {
			final Cache c = crstate.getCache();
			final JSONObject dsMap = c.getDeliveryServices();
			if(dsMap != null) {
				final String location = c.getLocation();
				for(String dsId : JSONObject.getNames(dsMap)) {
					try {
						final List<String> fqdns = getFqdns(dsId, dsMap);
						final DsState dss = DsState.getOrCreate(dsId);

						// Don't count the cache as reporting unless there were no errors and stats were read
						boolean error = false;
						boolean foundStats = false;

						for(String fqdn : fqdns) {
							final String propBase = "ats.plugin.remap_stats."+fqdn;
							final DsState.DsStati stati = DsState.createStati(propBase, crstate, lenientTime, dsId);

							dss.accumulate(stati, location, crstate);

							if (stati != null) {
								foundStats = true;

								if (stati.error) {
									error = true;
								}
							}
						}

						// Update cache counters
						dss.addCacheConfigured();

						if (crstate.isAvailable()) {
							dss.addCacheAvailable();
						}

						if (foundStats && !error) {
							dss.addCacheReporting();
						}
					} catch(Exception e) {
						LOGGER.warn(e,e);
					}
				}
			}
		}

		final Collection<String> toRemove = new ArrayList<String>();
		toRemove.addAll(states.keySet());
		for(String dsId : JSONObject.getNames(dsList)) {
			toRemove.remove(dsId);
			try {
				final DsState dss = getOrCreate(dsId);
				dss.completeRound(myHealthDeterminer.getDsControls(dss.getId()));
			} catch(Exception e) {
				LOGGER.warn(e,e);
			}
		}
		for(String id : toRemove) {
			states.remove(id);
		}
	}
	private static List<String> getFqdns(final String dsId, final JSONObject dsMap) throws JSONException {
		final org.apache.wicket.ajax.json.JSONArray ja = dsMap.optJSONArray(dsId);
		final ArrayList<String> fqdns = new ArrayList<String>();
		if(ja == null) {
			fqdns.add(dsMap.getString(dsId));
		} else {
			for (int i = 0; i < ja.length(); i++) {
				fqdns.add(ja.getString(i));
			}
		}
		return fqdns;
	}

	public static class EmbeddedStati implements java.io.Serializable {
		private static final long serialVersionUID = 1L;
		private DsStati currentDtati;
		private final String id;

		public EmbeddedStati(final String base, final String id, final String delimiter) {
			final StringBuilder statId = new StringBuilder();

			if (base != null) {
				statId.append(base);
				statId.append(delimiter);
			}

			statId.append(id);

			this.id = statId.toString();
		}

		public EmbeddedStati(final String base, final String id) {
			this(base, id, ".");
		}

		public EmbeddedStati(final String id) {
			this(null, id, ".");
		}

		public void accumulate(final DsStati stati) {
			if (currentDtati == null) {
				currentDtati = new DsStati(stati);
			} else {
				currentDtati.accumulate(stati);
			}
		}

		public Map<String, String> completeRound() {
			if (currentDtati == null) {
				return null;
			}

			final Map<String, String> r = new HashMap<String, String>();

			r.putAll(currentDtati.getStati(this.getId()));
			currentDtati = null;

			return r;
		}

		public String getId() {
			return id;
		}
	}

	public static DsStati createStati(final String propBase, final CacheState cs, final long leniency, final String dsId) {
		DsStati ds = null;
		synchronized (cs) {
			final List<DataPoint> dps = cs.getDataPoints(propBase+".out_bytes");
			if(dps == null) {
				return null;
			}
			long lastIndex = dps.get(dps.size()-1).getIndex();
			lastIndex = getLastGoodIndex(dps, lastIndex);
			if(lastIndex < 0) { return null; }
			final long time = cs.getTime(lastIndex);
			if(time < leniency) {
				return null;
			}
			ds  = new DsState.DsStati(propBase, cs, lastIndex, dsId);
			final long prevIndex = getLastGoodIndex(dps, lastIndex-1);
			if(prevIndex >= 0) {
				final DsStati priorDs = new DsState.DsStati(propBase, cs, prevIndex, dsId);
				if(!ds.calculateKbps(priorDs)) {
					if(LOGGER.isInfoEnabled()) {
						printDps(dps, propBase);
					}
				}
			}
		}
		return ds;
	}
	public static boolean printDps(final List<DataPoint> dps, final String id) {
		LOGGER.warn(id+":");
		for(int i = dps.size()-1; i >= 0; i--) {
			LOGGER.warn(String.format("\t%d - index: %d, span: %d, value: %s", i, 
					dps.get(i).getIndex(),
					dps.get(i).getSpan(),
					dps.get(i).getValue()
					));
		}
		return false;
	}
	private static long getLastGoodIndex(final List<DataPoint> dps, final long targetIndex) {
		if(targetIndex < 0) {
			return -1;
		}
		for(int i = dps.size()-1; i >= 0; i--) {
			if(dps.get(i).getValue()!=null) {
				final long index = dps.get(i).getIndex();
				final long span = dps.get(i).getSpan();
				if(targetIndex <= (index-span)) { continue; }
				if(targetIndex < index) {
					return targetIndex;
				}
				return index;
			}
		}
		return -1;
	}

	public static class DsStati implements java.io.Serializable {
		private static final long serialVersionUID = 1L;
		long csIndex = 0;
		long in_bytes; 
		long out_bytes; 
		long status_2xx;
		long status_3xx; 
		long status_4xx; 
		long status_5xx;
		boolean error = false;

		double kbps;
		double tps_2xx;
		double tps_3xx;
		double tps_4xx;
		double tps_5xx;
		double tps_total;

		String dsId;
		String csId;

		public static final int BITS_IN_BYTE = 8;
		//		public static final int BITS_IN_KBPS = 1000;
		public static final int MS_IN_SEC = 1000;
		public final long time;

		public DsStati(final String propBase, final CacheState cs, final long index, final String dsId) {
			this.csIndex = index;
			this.time = cs.getTime(index);
			String v = cs.getValue(propBase+".in_bytes", index);
			this.in_bytes = toLong(v);
			final String k = propBase+".out_bytes";
			v = cs.getValue(k, index);
			if(v == null) {
				LOGGER.warn("wtf: "+ cs.id + " - "+v);
				v = cs.getValue(k, index);
			}
			this.out_bytes = toLong(v);
			v = cs.getValue(propBase+".status_2xx", index);
			this.status_2xx = toLong(v);
			this.status_3xx = toLong(cs.getValue(propBase+".status_3xx", index));
			this.status_4xx = toLong(cs.getValue(propBase+".status_4xx", index));
			this.status_5xx = toLong(cs.getValue(propBase+".status_5xx", index));
			this.dsId = dsId;
			this.csId = cs.id;
		}
		public static boolean checkBytes(final List<DataPoint> dps, final String id) {
			long lastGoodIndex = -1;
			long goodValue = 0;
			for(int i = dps.size()-1; i >= 0; i--) {
				if(dps.get(i).getValue()==null) {
					continue;
				}
				if(lastGoodIndex == -1) {
					lastGoodIndex = dps.get(i).getIndex();
					goodValue = toLong(dps.get(i).getValue());
				} else {
					final long v = toLong(dps.get(i).getValue());
					if(v > goodValue) {
						LOGGER.warn(id+" - data error:" + v +" > "+ goodValue);
						return true;
					}
					break;
				}
			}
			return false;
		}
		public DsStati(final DsStati stati) {
			this.error = stati.error;
			this.time = stati.time;
			this.in_bytes = stati.in_bytes;
			this.out_bytes = stati.out_bytes;
			this.status_2xx = stati.status_2xx;
			this.status_3xx = stati.status_3xx;
			this.status_4xx = stati.status_4xx;
			this.status_5xx = stati.status_5xx;

			this.kbps = stati.kbps;
			this.tps_2xx = stati.tps_2xx;
			this.tps_3xx = stati.tps_3xx;
			this.tps_4xx = stati.tps_4xx;
			this.tps_5xx = stati.tps_5xx;
			this.tps_total = stati.tps_total;
		}
		private static long toLong(final String str) {
			if(str == null) {
				return 0;
			}
			return (long) Double.parseDouble(str);
		}
		void accumulate(final DsStati ds) {
			this.in_bytes += ds.in_bytes;
			this.out_bytes += ds.out_bytes;
			this.status_2xx += ds.status_2xx;
			this.status_3xx += ds.status_3xx;
			this.status_4xx += ds.status_4xx;
			this.status_5xx += ds.status_5xx;

			this.kbps += ds.kbps;
			this.tps_2xx += ds.tps_2xx;
			this.tps_3xx += ds.tps_3xx;
			this.tps_4xx += ds.tps_4xx;
			this.tps_5xx += ds.tps_5xx;
			this.tps_total += ds.tps_total;
		}

		public boolean calculateKbps(final DsStati prior) {
			if(prior == null) {
				LOGGER.warn("why is prior null");
				return false;
			}
			if(prior.time == 0) {
				LOGGER.warn("why is prior.time==0");
			}
			if((out_bytes == 0 || prior.out_bytes == 0) && out_bytes != prior.out_bytes) {
				LOGGER.warn(dsId+": throwing out "+csId+": out_bytes==0");
				if(prior.out_bytes != 0) {
					LOGGER.warn("\t prior.out_bytes="+prior.out_bytes);
				}
				return false;
			}
			final long deltaTimeMs = time - prior.time; // / MS_IN_SEC
			if(LOGGER.isDebugEnabled()) {
				LOGGER.debug(String.format("time delta: %d, index: %d -> %d", deltaTimeMs, prior.csIndex, this.csIndex));
			}
			if(deltaTimeMs == 0) {
				LOGGER.warn("time delta 0");
				return false;
			}
			final long delta = (out_bytes - prior.out_bytes);
			// as long as the numbers are not too large, dividing both num and denom by 1000 is a waste of time
			//			rates.kbps = (delta / BITS_IN_KBPS) * BITS_IN_BYTE / deltaTime;
			kbps = ((double) delta / (double) deltaTimeMs) * BITS_IN_BYTE;
			if(kbps < 0.0) {
				LOGGER.warn(dsId+": throwing out "+csId+": kbps="+ kbps);
				kbps = 0.0;
				return false;
			}

			final double deltaTime = (double) deltaTimeMs / (double) MS_IN_SEC;

			tps_2xx = ((double) status_2xx - (double) prior.status_2xx) / deltaTime;
			tps_3xx = ((double) status_3xx - (double) prior.status_3xx) / deltaTime;
			tps_4xx = ((double) status_4xx - (double) prior.status_4xx) / deltaTime;
			tps_5xx = ((double) status_5xx - (double) prior.status_5xx) / deltaTime;
			tps_total = tps_2xx + tps_3xx + tps_4xx + tps_5xx;

			return true;
		}
		Map<String, String> getStati(final String base) {
			final Map<String, String> r = new HashMap<String, String>();
			r.put(base+".in_bytes", String.valueOf(in_bytes));
			r.put(base+".out_bytes", String.valueOf(out_bytes));
			r.put(base+".status_2xx", String.valueOf(status_2xx));
			r.put(base+".status_3xx", String.valueOf(status_3xx));
			r.put(base+".status_4xx", String.valueOf(status_4xx));
			r.put(base+".status_5xx", String.valueOf(status_5xx));

			DecimalFormat df = new DecimalFormat("0.00");
			r.put(base+".kbps", df.format(kbps));
			r.put(base+".tps_2xx", df.format(tps_2xx));
			r.put(base+".tps_3xx", df.format(tps_3xx));
			r.put(base+".tps_4xx", df.format(tps_4xx));
			r.put(base+".tps_5xx", df.format(tps_5xx));
			r.put(base+".tps_total", df.format(tps_total));
			return r;
		}
	}
	public static Collection<DsState> getDsStates() {
		return states.values();
	}
	@Override
	protected KeyValue getKeyValue(final String key, final AbstractState state) {
		return new KeyValue(key,this) {
			private static final long serialVersionUID = 1L;
			@Override
			public String getObject( ) {
				if(stateId != null) {
					return DsState.get(stateId, key);
				}
				return val;
			}
		};
	}
	public static String get(final String stateId, final String key) {
		return get(stateId).getLastValue(key);
	}
	public static DsState get(final String host) {
		synchronized(states) {
			return states.get(host);
		}
	}
	public JSONArray getDisabledLocations() throws JSONException {
		return new JSONArray("["+this.getLastValue("disabledLocations")+"]");
	}

	public static boolean has(final String host) {
		if(states.get(host)==null) { return false; }
		return true;
	}
	public static void startUpdateAll() {
		synchronized(states) {
			for(DsState ds :states.values()) {
				ds.startUpdate();
			}
		}
	}

	public void addCacheConfigured() {
		this.cachesConfigured++;
	}

	public void addCacheAvailable() {
		this.cachesAvailable++;
	}

	public void addCacheReporting() {
		this.cachesReporting++;
	}

}
