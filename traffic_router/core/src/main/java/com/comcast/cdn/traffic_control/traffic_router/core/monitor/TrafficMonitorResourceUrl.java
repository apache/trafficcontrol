package com.comcast.cdn.traffic_control.traffic_router.core.monitor;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ResourceUrl;

class TrafficMonitorResourceUrl implements ResourceUrl {
	private final TrafficMonitorWatcher trafficMonitorWatcher;
	private final String urlTemplate;
	private int i = 0;
	public TrafficMonitorResourceUrl(final TrafficMonitorWatcher trafficMonitorWatcher, final String urlTemplate) {
		this.trafficMonitorWatcher = trafficMonitorWatcher;
		this.urlTemplate = urlTemplate;
	}
	@Override
	public String nextUrl() {
		final String[] hosts = trafficMonitorWatcher.getHosts();

		if (hosts == null || hosts.length == 0) {
			return urlTemplate;
		}

		i %= hosts.length;
		final String host = hosts[i];
		i++;
		return urlTemplate.replace("[host]", host);
	}
}
