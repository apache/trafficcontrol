/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;

public class Event extends JSONObject implements Serializable {
	private static final Logger EVENT_LOGGER = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_monitor.event");
	private static final Logger LOGGER = Logger.getLogger(Event.class);
	private static final long serialVersionUID = 1L;
	static List<JSONObject> rollingLog = new LinkedList<JSONObject>();
	static int logIndex = 0;

	public enum EventType {
		CACHE_STATE_CHANGE,
		PEER_STATE_CHANGE("Peer"),
		DELIVERY_SERVICE_STATE_CHANGE("Delivery Service");

		private String type;

		EventType(final String type) {
			this.type = type;
		}

		EventType() {
		}

		@Override
		public String toString() {
			return (type != null) ? type : "UNKNOWN";
		}

		public void setType(final String type) {
			this.type = type;
		}
	}

	public static Event logStateChange(final String name, final EventType type, final boolean isAvailable, final String message) {
		final long currentTimeMillis = System.currentTimeMillis();
		final String timeString = String.format("%d.%03d", currentTimeMillis / 1000, currentTimeMillis % 1000);

		EVENT_LOGGER.info(String.format("%s host=\"%s\", type=%s, available=%s, msg=\"%s\"", timeString , name, type, String.valueOf(isAvailable), message));

		final Event ret = new Event(name, type, isAvailable, message);
		final int eventLogCount = ConfigHandler.getInstance().getConfig().getEventLogCount();

		synchronized (rollingLog) {
			rollingLog.add(0, ret);
			while(rollingLog.size() > eventLogCount) {
				rollingLog.remove(rollingLog.size()-1);
			}
		}
		return ret;
	}

	public static List<JSONObject> getEventLog() {
		synchronized (rollingLog) {
			return new ArrayList<JSONObject>(rollingLog);
		}
	}

	public Event(final String name, final EventType type, final boolean isAvailable, final String error) {
		try {
			this.put("hostname", name); // left this to preserve behavior for any dependent uses
			this.put("name", name);
			this.put("type", type.toString());
			this.put("time", System.currentTimeMillis());
			this.put("index", logIndex++);
			this.put("isAvailable", isAvailable);
			this.put("description", error);

		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
	}
}
