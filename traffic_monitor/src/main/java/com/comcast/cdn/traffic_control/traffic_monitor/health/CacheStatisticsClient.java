package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.AsyncHttpClientConfig;
import com.ning.http.client.ProxyServer;
import com.ning.http.client.Request;
import org.apache.log4j.Logger;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Future;

public class CacheStatisticsClient {
	private static final Logger LOGGER = Logger.getLogger(CacheStatisticsClient.class);
	private final AsyncHttpClient asyncHttpClient = new AsyncHttpClient();

	public void fetchCacheStatistics(final Cache cache, final CacheStateUpdater cacheStateUpdater) {
		int port = cache.getQueryPort() != 0 ? cache.getQueryPort() : 80;
		final ProxyServer proxyServer = new ProxyServer(cache.getQueryIp(), port);

		Request request = asyncHttpClient
			.prepareGet(cache.getStatisticsUrl())
			.setProxyServer(proxyServer)
			.build();

		try {
			final Future<Object> future = asyncHttpClient.executeRequest(request, cacheStateUpdater);
			cacheStateUpdater.setFuture(future);
		} catch (IOException e) {
			LOGGER.warn("Failed to fetch cache statistics from " + request.getUrl(),e);
		}
	}


	public void shutdown() {
		while (!asyncHttpClient.isClosed()) {
			LOGGER.warn("closing");
			asyncHttpClient.close();
		}
	}
}
