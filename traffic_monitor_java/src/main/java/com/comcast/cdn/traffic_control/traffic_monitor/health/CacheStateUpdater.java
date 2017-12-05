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


import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.CacheDataModel;
import com.ning.http.client.AsyncCompletionHandler;
import com.ning.http.client.Response;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import java.io.IOException;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.concurrent.CancellationException;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;

public class CacheStateUpdater extends AsyncCompletionHandler<Object> {
	final private static Logger LOGGER = Logger.getLogger(CacheStateUpdater.class);
	final private CacheState state;
	final private CacheDataModel errorCount;
	private long time;
	private long requestTimeout;
	private HealthDeterminer myHealthDeterminer;
	private AtomicInteger failCount;
	private Future<Object> future;

	public CacheStateUpdater(final CacheState cacheState, final CacheDataModel errorCount) {
		this.state = cacheState;
		this.errorCount = errorCount;
	}

	public CacheStateUpdater update(final HealthDeterminer myHealthDeterminer, AtomicInteger failCount, long requestTimeout) {
		this.myHealthDeterminer = myHealthDeterminer;
		this.time = System.currentTimeMillis();
		this.failCount = failCount;
		this.requestTimeout = requestTimeout;
		return this;
	}

	@Override
	public Integer onCompleted(final Response response) throws JSONException, IOException {
		final int code = response.getStatusCode();
		state.putDataPoint("queryTime", String.valueOf(System.currentTimeMillis() - time));

		if (code != 200) {
			state.setError(code + " - " + response.getStatusText());
			return code;
		}

		final Map<String, String> statisticsMap = new HashMap<String, String>();

		final JSONObject json = new JSONObject(response.getResponseBody());
		JSONObject ats = json.has("global") ? json.optJSONObject("global") : json.optJSONObject("ats");

		statisticsMap.putAll(jsonToPrefixedMap(ats, "ats."));
		statisticsMap.putAll(jsonToPrefixedMap(json.optJSONObject("system"), "system."));

		state.putDataPoints(statisticsMap);
		state.putDataPoint("stateUrl", state.getCache().getStatisticsUrl());

		synchronized (state.getCache()) {
			state.getCache().setState(state, myHealthDeterminer);
		}

		return code;
	}

	@Override
	public void onThrowable(final Throwable t) {
		if (!(t instanceof CancellationException)) {
			LOGGER.warn(t + " : " + state.getCache().getStatisticsUrl());
			failCount.incrementAndGet();
		} else {
			LOGGER.warn("Request to " + state.getCache().getStatisticsUrl() + " failed to complete in time");
		}

		state.putDataPoint("queryTime", String.valueOf(System.currentTimeMillis() - time));

		try {
			errorCount.inc();
			state.setError(t.toString());
		} catch (Exception e2) {
			LOGGER.warn(e2, e2);
		}
	}

	private Map<String, String> jsonToPrefixedMap(JSONObject json, final String prefix) {
		Map<String, String> map = new HashMap<String, String>();

		Iterator<?> keys = json.keys();

		while (keys.hasNext()) {
			final String key = (String) keys.next();
			map.put(prefix + key, String.valueOf(json.opt(key)));
		}

		return map;
	}

	public void setFuture(final Future<Object> future) {
		this.future = future;
	}

	public boolean completeFetchStatistics(final AtomicInteger cancelCount) {
		if (future == null || future.isDone() || future.isCancelled()) {
			return true;
		}

		if (System.currentTimeMillis() > requestTimeout) {
			try {
				future.cancel(true);
				cancelCount.incrementAndGet();
			} catch (Exception e) {
				LOGGER.warn("Error on cancel: " + e);
			}

			return true;
		}

		return false;
	}
}
