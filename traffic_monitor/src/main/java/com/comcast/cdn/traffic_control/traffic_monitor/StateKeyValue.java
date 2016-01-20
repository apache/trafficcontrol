package com.comcast.cdn.traffic_control.traffic_monitor;

import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.StateRegistry;

public class StateKeyValue extends KeyValue {
	private final String stateId;
	private final StateRegistry stateRegistry;

	public StateKeyValue(String key, AbstractState state, StateRegistry stateRegistry) {
		super(key, null);
		this.stateRegistry = stateRegistry;
		this.stateId = state.getId();
	}

	@Override
	public String getObject() {
		if (stateId != null) {
			return stateRegistry.get(stateId, getKey());
		}
		return super.getObject();
	}
}
