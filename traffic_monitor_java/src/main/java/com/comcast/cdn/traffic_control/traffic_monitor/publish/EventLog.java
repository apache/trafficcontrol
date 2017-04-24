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

package com.comcast.cdn.traffic_control.traffic_monitor.publish;

import java.util.List;

import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.health.Event;

public class EventLog extends JsonPage {
	private static final long serialVersionUID = 1L;

	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
		final JSONObject o = new JSONObject();
		final List<JSONObject> list = Event.getEventLog();
		final JSONArray events = new JSONArray(list);
		o.put("events", events);
		return o;
	}

}



