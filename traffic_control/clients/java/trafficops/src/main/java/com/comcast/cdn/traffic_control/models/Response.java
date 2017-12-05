package com.comcast.cdn.traffic_control.models;

import java.util.List;
import java.util.Map;

public class Response {
	private List<Alert> alerts;

	public List<Alert> getAlerts() {
		return alerts;
	}
	public void setAlerts(List<Alert> alerts) {
		this.alerts = alerts;
	}
	
	public class CollectionResponse extends Response {
		private List<Map<String, Object>> response;

		public List<Map<String, Object>> getResponse() {
			return response;
		}

		public void setResponse(List<Map<String, Object>> response) {
			this.response = response;
		}
	}
}
