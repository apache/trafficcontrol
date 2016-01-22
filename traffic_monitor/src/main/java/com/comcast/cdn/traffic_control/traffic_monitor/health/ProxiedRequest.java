package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.ProxyServer;
import com.ning.http.client.Request;
import org.apache.log4j.Logger;

public class ProxiedRequest {
	private static final Logger LOGGER = Logger.getLogger(ProxiedRequest.class);
	private final String url;
	private final int port;
	private final String ipAddress;
	private final Request request;

	public ProxiedRequest(final String ipAddress, final int port, final String url, AsyncHttpClient asyncHttpClient) {
		this.url = url;
		this.port = port;
		this.ipAddress = ipAddress;

		final AsyncHttpClient.BoundRequestBuilder builder = asyncHttpClient.prepareGet(this.url);
		final ProxyServer proxyServer = new ProxyServer(this.ipAddress, this.port);
		builder.setProxyServer(proxyServer);
		request = builder.build();
	}

	public String getIpAddress() {
		return ipAddress;
	}

	public int getPort() {
		return port;
	}

	public String getUrl() {
		return url;
	}

	public Request getRequest() {
		return request;
	}

	public ProxiedRequest updateForCache(final Cache cache, final AsyncHttpClient asyncHttpClient) {
		if (!(cache.getQueryIp().equals(ipAddress) && cache.getQueryPort() == port && cache.getStatisticsUrl().equals(url))) {
			return this;
		}

		String updatedIpAddress = ipAddress;
		int updatedPort = port;
		String updatedUrl = url;

		if (!cache.getQueryIp().equals(ipAddress)) {
			LOGGER.info("Health polling IP change detected for " + cache.getStatisticsUrl() + " (new != old): " + cache.getQueryIp() + " != " + ipAddress);
			updatedIpAddress = cache.getQueryIp();
		}

		if (cache.getQueryPort() != updatedPort) {
			LOGGER.info("Health polling port change detected for " + cache.getStatisticsUrl() + " (new != old): " + cache.getQueryPort() + " != " + updatedPort);
			updatedPort = port;
		}

		if (!cache.getStatisticsUrl().equals(updatedUrl)) {
			LOGGER.info("Health polling URL change detected for " + cache.getStatisticsUrl() + " (new != old): " + cache.getStatisticsUrl() + " != " + updatedUrl);
			updatedUrl = cache.getStatisticsUrl();
		}

		if (updatedPort == 0) {
			updatedPort = 80;
		}

		return new ProxiedRequest(updatedIpAddress, updatedPort, updatedUrl, asyncHttpClient);
	}
}
