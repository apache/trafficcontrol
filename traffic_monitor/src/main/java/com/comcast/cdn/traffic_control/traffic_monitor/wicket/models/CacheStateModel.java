package com.comcast.cdn.traffic_control.traffic_monitor.wicket.models;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;

public class CacheStateModel extends StateModel {
	private static final long serialVersionUID = 1L;

	public CacheStateModel(final String stateName, final String key) {
		super(stateName, key);
	}

	@Override
	public String getObject() {
		return getObject(CacheStateRegistry.getInstance().get(stateName));
	}
}
