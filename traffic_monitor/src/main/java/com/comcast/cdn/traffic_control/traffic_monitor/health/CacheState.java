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

import java.io.IOException;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CancellationException;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.KeyValue;
import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher.CacheDataModel;
import com.ning.http.client.AsyncCompletionHandler;
import com.ning.http.client.AsyncHandler;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.AsyncHttpClient.BoundRequestBuilder;
import com.ning.http.client.AsyncHttpClientConfig;
import com.ning.http.client.ProxyServer;
import com.ning.http.client.Request;
import com.ning.http.client.Response;

public class CacheState extends AbstractState {
	private static final Logger LOGGER = Logger.getLogger(CacheState.class);
	private static final long serialVersionUID = 1L;
	private static Map<String, CacheState> states = new HashMap<String, CacheState>();

	transient private Future<Object> future;
	transient private Request request;
	transient private String usedIp;
	transient private int usedPort;
	transient private String usedUrl;
	transient private long requestTimeout;
	transient private UpdateHandler handler;
	transient private Cache cache;
	static private AsyncHttpClient asyncHttpClient;

	private CacheState(final String id) {
		super(id);
	}

	public static List<CacheState> getCacheStates() {
		synchronized (states) {
			return new ArrayList<CacheState>(states.values());
		}
	}

	public static CacheState getOrCreate(final Cache cache) {
		return getOrCreate(cache.getHostname(), cache);
	}

	public static CacheState getOrCreate(final String host, final Cache cache) {
		synchronized (states) {
			CacheState as = states.get(host);

			if (as == null) {
				as = new CacheState(host);
				states.put(host, as);
			}

			as.setCache(cache);
			return as;
		}
	}

	private void setCache(final Cache cache) {
		this.cache = cache;
	}

	public static CacheState getState(final String host) {
		synchronized (states) {
			return states.get(host);
		}
	}

	public static boolean has(final String host) {
		if (states.get(host) == null) {
			return false;
		}

		return true;
	}

	public static void removeAllBut(final List<CacheState> retList) {
		final List<String> hostnames = new ArrayList<String>();

		for (CacheState cs : retList) {
			hostnames.add(cs.getId());
		}

		synchronized (states) {
			for (String key : new ArrayList<String>(states.keySet())) {
				if (!hostnames.contains(key)) {
					states.remove(key);
				}
			}
		}
	}


	public Cache getCache() {
		return cache;
	}

	public void fetchAndUpdate(final HealthDeterminer myHealthDeterminer, final CacheDataModel fetchCount, final CacheDataModel errorCount, final AtomicInteger failCount) {
		if (!HealthDeterminer.shouldFetchStats(cache)) {
			synchronized (cache) {
				// TODO : clear states
				cache.setState(this, myHealthDeterminer);
			}
			return;
		}

		final AsyncHttpClient asyncClient = getAsyncHttpClient();
		final long time = System.currentTimeMillis();

		try {
			fetchCount.inc();
			this.startUpdate();
			final String url = getFetchUrl();
			this.put("_queryUrl", url);
			this.setHistoryTime(cache.getHistoryTime());
			requestTimeout = System.currentTimeMillis() + myHealthDeterminer.getConnectionTimeout(cache, 2000);
			future = asyncClient.executeRequest(getRequest(asyncClient, url), getAsyncHanlder(myHealthDeterminer, time, url, errorCount, failCount));
		} catch (IOException e) {
			LOGGER.warn(e, e);
		}
	}


	private AsyncHandler<Object> getAsyncHanlder(final HealthDeterminer myHealthDeterminer, final long time,
	                                             final String url, final CacheDataModel errorCount, final AtomicInteger failCount) {
		if (handler == null) {
			handler = new UpdateHandler(this, errorCount);
		}

		return handler.update(myHealthDeterminer, time, url, failCount);
	}

	private static class UpdateHandler extends AsyncCompletionHandler<java.lang.Object> {
		final private CacheState state;
		final private CacheDataModel errorCount;
		private long time;
		private HealthDeterminer myHealthDeterminer;
		private String url;
		private AtomicInteger failCount;

		public UpdateHandler(final CacheState cacheState, final CacheDataModel errorCount) {
			this.state = cacheState;
			this.errorCount = errorCount;
		}

		public AsyncHandler<Object> update(final HealthDeterminer myHealthDeterminer, final long time, final String url, AtomicInteger failCount) {
			this.myHealthDeterminer = myHealthDeterminer;
			this.time = time;
			this.url = url;
			this.failCount = failCount;
			return this;
		}

		@Override
		public Integer onCompleted(final Response response) throws JSONException, IOException {
			// Do something with the Response
			final int code = response.getStatusCode();
			state.put("queryTime", String.valueOf(System.currentTimeMillis() - time));

			if (code != 200) {
				synchronized (state.cache) {
					state.cache.setError(state, code + " - " + response.getStatusText(), myHealthDeterminer);
				}
				return code;
			}

			//			final long queryTime = System.currentTimeMillis() - time;
			final Map<String, String> stats = getMap(response.getResponseBody(), url);
			state.put(stats);

			synchronized (state.cache) {
				state.cache.setState(state, myHealthDeterminer);
			}

			return code;
		}

		@Override
		public void onThrowable(final Throwable t) {
			if (!(t instanceof CancellationException)) {
				LOGGER.warn(t + " : " + url);
				failCount.incrementAndGet();
			} else {
				LOGGER.warn("Request to " + url + " failed to complete in time");
			}

			state.put("queryTime", String.valueOf(System.currentTimeMillis() - time));

			try {
				errorCount.inc();
				synchronized (state.cache) {
					state.cache.setError(state, t.toString(), myHealthDeterminer);
				}
			} catch (Exception e2) {
				LOGGER.warn(e2, e2);
			}
		}
	}

	private Request getRequest(final AsyncHttpClient asyncClient, final String url) {
		if (request == null || !this.getCache().getQueryIp().equals(usedIp) || this.getCache().getQueryPort() != usedPort || !url.equals(usedUrl)) {
			if (request != null && !this.getCache().getQueryIp().equals(usedIp)) {
				LOGGER.info("Health polling IP change detected for " + url + " (new != old): " + this.getCache().getQueryIp() + " != " + usedIp);
			}

			if (request != null && this.getCache().getQueryPort() != usedPort) {
				LOGGER.info("Health polling port change detected for " + url + " (new != old): " + this.getCache().getQueryPort() + " != " + usedPort);
			}

			if (request != null && !url.equals(usedUrl)) {
				LOGGER.info("Health polling URL change detected for " + url + " (new != old): " + url + " != " + usedUrl);
			}

			final BoundRequestBuilder builder = asyncClient.prepareGet(url);
			usedIp = this.getCache().getQueryIp();
			usedPort = this.getCache().getQueryPort();
			usedUrl = url;

			if (usedPort == 0) {
				usedPort = 80;
			}

			final ProxyServer proxyServer = new ProxyServer(usedIp, usedPort);
			builder.setProxyServer(proxyServer);
			request = builder.build();
		}
		return request;
	}

	private AsyncHttpClient getAsyncHttpClient() {
		synchronized (LOGGER) {
			if (asyncHttpClient == null) {
				final AsyncHttpClientConfig cf = new AsyncHttpClientConfig.Builder()
					//			.setConnectionTimeoutInMs(myTimeout)
					//			.addRequestFilter(new ThrottleRequestFilter(10))
					.build();
				asyncHttpClient = new AsyncHttpClient(cf);
			}
		}

		return asyncHttpClient;
	}

	public boolean completeFetch(final HealthDeterminer myHealthDeterminer, final CacheDataModel errorCount, final AtomicInteger cancelCount, final AtomicInteger failCount) {
		if (future == null) {
			return true;
		}

		// Does this logic change make performance better????
		// if (future.isDone() || future.isCancelled()) {

		if (future.isDone()) {
			return true;
		}

		if (System.currentTimeMillis() > requestTimeout) {
			try {
				future.cancel(true);
				cancelCount.incrementAndGet();
				//				errorCount.inc();
			} catch (Exception e) {
				LOGGER.warn("Error on cancel: " + e);
			}

			return true;
		}

		return false;
	}

	private String getFetchUrl() {
		return HealthDeterminer.getStatusUrl(cache);
	}

	public static Map<String, String> getMap(final String jsonStr, final String stateUrl) throws JSONException {
		final Map<String, String> map = new HashMap<String, String>();
		final JSONObject json = new JSONObject(jsonStr);
		JSONObject global = json.optJSONObject("global");

		if (global == null) {
			global = json.optJSONObject("ats");
		}

		Iterator<?> keys = global.keys();

		while (keys.hasNext()) {
			final String key = (String) keys.next();
			map.put("ats." + key, String.valueOf(global.get(key)));
		}

		global = json.optJSONObject("system");
		keys = global.keys();

		while (keys.hasNext()) {
			final String key = (String) keys.next();
			map.put("system." + key, String.valueOf(global.get(key)));
		}

		//		map.put("queryTime", Long.toString(queryTime));
		map.put("stateUrl", stateUrl);

		return map;
	}

	public static void shutdown() {
		while (!asyncHttpClient.isClosed()) {
			LOGGER.warn("closing");
			asyncHttpClient.close();
		}
	}


	@Override
	protected KeyValue getKeyValue(final String key, final AbstractState state) {
		return new KeyValue(key, this) {
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject() {
				if (stateId != null) {
					return CacheState.get(stateId, key);
				}
				return val;
			}
		};
	}

	public static String get(final String stateId, final String key) {
		return getState(stateId).getLastValue(key);
	}
}