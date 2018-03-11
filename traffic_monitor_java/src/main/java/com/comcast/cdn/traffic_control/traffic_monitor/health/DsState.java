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

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.EmbeddedStati.StatType;

public class DsState extends AbstractState {
	private static final long serialVersionUID = 1L;
	final private Map<StatType, Map<String, EmbeddedStati>> aggregateStats = new HashMap<StatType, Map<String, EmbeddedStati>>();
	final public static String DISABLED_LOCATIONS = "disabledLocations";

	private DsStati currentDsStati;
	private int cachesConfigured = 0;
	private int cachesAvailable = 0;
	private int cachesReporting = 0;

	public DsState(final String id) {
		super(id);
	}

	public void accumulate(final DsStati stati, final String location, final CacheState state) {
		if (stati == null) {
			return;
		}

		if (currentDsStati == null) {
			currentDsStati = stati;
		} else {
			currentDsStati.accumulate(stati);
		}

		aggregateStats(StatType.LOCATION, location, stati);
		aggregateStats(StatType.CACHE, state.id, stati);
		aggregateStats(StatType.TYPE, state.getCache().getType(), stati);
	}

	private void aggregateStats(final StatType statType, final String statKey, final DsStati dsStat) {
		if (!aggregateStats.containsKey(statType)) {
			aggregateStats.put(statType, new HashMap<String, EmbeddedStati>());
		}

		final Map<String, EmbeddedStati> aggregate = aggregateStats.get(statType);

		if (!aggregate.containsKey(statKey)) {
			aggregate.put(statKey, new EmbeddedStati(statType, statKey));
		}

		aggregate.get(statKey).accumulate(dsStat);
	}

	public boolean completeRound(final JSONObject dsControls) {
		if (currentDsStati != null && currentDsStati.out_bytes != 0) {
			putDataPoints(currentDsStati.getStati("total"));
			currentDsStati = null;
		}

		putDataPoint("caches-configured", String.valueOf(cachesConfigured));
		putDataPoint("caches-available", String.valueOf(cachesAvailable));
		putDataPoint("caches-reporting", String.valueOf(cachesReporting));

		cachesConfigured = 0;
		cachesAvailable = 0;
		cachesReporting = 0;

		HealthDeterminer.setIsAvailable(this, dsControls);

		for (Map<String, EmbeddedStati> aggregate : aggregateStats.values()) {
			processDataPoints(aggregate.values(), dsControls);
		}

		return true;
	}

	private void processDataPoints(final Collection<EmbeddedStati> stats, final JSONObject dsControls) {
		final Map<StatType, StringBuilder> disabled = new HashMap<StatType, StringBuilder>();

		for (EmbeddedStati stat : stats) {
			final Map<String, String> points = stat.completeRound();

			if (points == null) {
				continue;
			}

			putDataPoints(points);

			if (stat.isHidden()) {
				addHiddenStats(points.keySet());
			}

			if (stat.getStatType() == StatType.LOCATION) {
				if (!disabled.containsKey(stat.getStatType())) {
					disabled.put(stat.getStatType(), new StringBuilder());
				}

				if (!HealthDeterminer.setIsLocationAvailable(this, stat, dsControls)) {
					disabled.get(stat.getStatType()).append("\"").append(stat.getId()).append("\", ");
				}
			}
		}

		for (StatType statType : disabled.keySet()) {
			final StringBuilder sb = disabled.get(statType);

			if (statType == StatType.LOCATION && sb != null) {
				final String s = sb.toString();
				putDataPoint(DISABLED_LOCATIONS, s);
			}
		}
	}

	public JSONArray getDisabledLocations() throws JSONException {
		return new JSONArray("["+this.getLastValue(DISABLED_LOCATIONS)+"]");
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
