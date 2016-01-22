package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.AsyncHttpClientConfig;
import com.ning.http.client.Request;
import org.apache.log4j.Logger;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Future;

public class CacheStatisticsClient {
	private static final Logger LOGGER = Logger.getLogger(CacheStatisticsClient.class);
	private final AsyncHttpClient asyncHttpClient = new AsyncHttpClient(new AsyncHttpClientConfig.Builder().build());
	private final Map<String, ProxiedRequest> statisticsRequestMap = new HashMap<String, ProxiedRequest>();

	public void fetchCacheStatistics(final Cache cache, final CacheStateUpdater cacheStateUpdater) {
		Request request = getRequest(cache);
		try {
			final Future<Object> future = asyncHttpClient.executeRequest(request, cacheStateUpdater);
			cacheStateUpdater.setFuture(future);
		} catch (IOException e) {
			LOGGER.warn(e,e);
		}
	}

	private Request getRequest(final Cache cache) {
		ProxiedRequest proxiedRequest;

		if (!statisticsRequestMap.containsKey(cache.getFqdn())) {
			proxiedRequest = new ProxiedRequest(cache.getQueryIp(), cache.getQueryPort(), cache.getStatisticsUrl(), asyncHttpClient);
		} else {
			proxiedRequest = statisticsRequestMap.get(cache.getFqdn()).updateForCache(cache, asyncHttpClient);
		}

		statisticsRequestMap.put(cache.getFqdn(), proxiedRequest);


		return proxiedRequest.getRequest();
	}


	public void shutdown() {
		while (!asyncHttpClient.isClosed()) {
			LOGGER.warn("closing");
			asyncHttpClient.close();
		}
	}
}
