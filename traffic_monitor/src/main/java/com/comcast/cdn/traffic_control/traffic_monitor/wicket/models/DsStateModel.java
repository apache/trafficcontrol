package com.comcast.cdn.traffic_control.traffic_monitor.wicket.models;

import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;

public class DsStateModel extends StateModel {
	private static final long serialVersionUID = 1L;

	public DsStateModel(final String stateName, final String key) {
		super(stateName, key);
	}

	@Override
	public String getObject() {
		return getObject(DeliveryServiceStateRegistry.getInstance().get(stateName));
	}
}
