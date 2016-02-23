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

import java.util.HashMap;
import java.util.Map;

import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

public class DsState extends AbstractState {
	private static final long serialVersionUID = 1L;

	private DsStati currentDsStati;
	private int cachesConfigured = 0;
	private int cachesAvailable = 0;
	private int cachesReporting = 0;
	final private Map<String, EmbeddedStati> locs = new HashMap<String, EmbeddedStati>();
	final private Map<String, EmbeddedStati> cacheStatiMap = new HashMap<String, EmbeddedStati>();

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

		final StringBuilder sb = new StringBuilder();

		for (String locId : locs.keySet()) {
			final EmbeddedStati loc = locs.get(locId);
			final Map<String, String> stati = loc.completeRound();

			if (stati == null) {
				continue;
			}

			putDataPoints(stati);

			if (!HealthDeterminer.setIsAvailable(this, loc, dsControls)) {
				sb.append("\"").append(locId).append("\", ");
			}
		}

		putDataPoint("disabledLocations", sb.toString());

		for (String cacheId : cacheStatiMap.keySet()) {
			final EmbeddedStati cacheStat = cacheStatiMap.get(cacheId);
			final Map<String, String> stati = cacheStat.completeRound();

			if (stati == null) {
				continue;
			}

			addHiddenStats(stati.keySet());

			putDataPoints(stati);
		}

		return true;
	}

	public JSONArray getDisabledLocations() throws JSONException {
		return new JSONArray("["+this.getLastValue("disabledLocations")+"]");
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
