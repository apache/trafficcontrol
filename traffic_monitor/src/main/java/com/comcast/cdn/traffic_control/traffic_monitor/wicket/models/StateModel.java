package com.comcast.cdn.traffic_control.traffic_monitor.wicket.models;

import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import org.apache.wicket.model.Model;

public class StateModel extends Model<String> {
	String stateName;
	String key;

	public StateModel(final String stateName, final String key) {
		this.stateName = stateName;
		this.key = key;
	}

	protected String getObject(AbstractState state) {
		if (state == null) {
			return "err";
		}
		if ("_status_string_".equals(key)) {
			return state.getStatusString();
		}
		final boolean clearData = state.getBool("clearData");
		if (clearData) {
			return "-";
		}
		return state.getLastValue(key);
	}
}
