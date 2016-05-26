package com.comcast.cdn.traffic_control.traffic_router.core.util;

class DefaultResourceUrl implements ResourceUrl {
	private final String[] urla;
	private int i = 0;

	public DefaultResourceUrl(final String[] urla) {
		this.urla = urla;
	}

	@Override
	public String nextUrl() {
		i++;
		i %= urla.length;
		return urla[i];
	}
}
