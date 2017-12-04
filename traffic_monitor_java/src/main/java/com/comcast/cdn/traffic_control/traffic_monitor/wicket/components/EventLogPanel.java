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

package com.comcast.cdn.traffic_control.traffic_monitor.wicket.components;


import java.text.DateFormat;
import java.util.ArrayList;
import java.util.Date;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.AjaxSelfUpdatingTimerBehavior;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.IModel;
import org.apache.wicket.model.Model;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.health.Event;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.UpdatingAttributeAppender;


public class EventLogPanel extends Panel {
	private static final Logger LOGGER = Logger.getLogger(EventLogPanel.class);
	private static final long serialVersionUID = 1L;

	static final DateFormat formatter = DateFormat.getDateTimeInstance(DateFormat.SHORT, DateFormat.SHORT);

	public EventLogPanel(final String id) {
		super(id);

		final IModel<ArrayList<JSONObject>> listModel =  new Model<ArrayList<JSONObject>>() {
			private static final long serialVersionUID = 1L;
			@Override
			@SuppressWarnings("PMD")
			public ArrayList<JSONObject> getObject( ) {
				return new ArrayList<JSONObject>(Event.getEventLog());
			}
		};
		final ListView<JSONObject> propView2 = new ListView<JSONObject>("events", listModel ) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<JSONObject> item) {
				final JSONObject jo = item.getModelObject();
				String errorClass = "";
				try {
					item.add(new Label("index", Long.toString(jo.getLong("index"))));
//					long time = System.currentTimeMillis();
					long time = jo.getLong("time");
					synchronized(formatter) {
						item.add(new Label("time", formatter.format(new Date(time))));
					}
					item.add(new Label("timeraw", Long.toString(jo.getLong("time"))));
					item.add(new Label("description", jo.getString("description")));
					item.add(new Label("name", jo.getString("name")));
					item.add(new Label("type", jo.getString("type")));
					String status = "available";
					if(!jo.getBoolean("isAvailable")) {
						status = "offline";
						errorClass = "error";
					}
					item.add(new Label("status", status));
				} catch (JSONException e) {
					LOGGER.warn(e,e);
				}
				item.add(new UpdatingAttributeAppender("class", new Model<String>(errorClass), " "));

			}
		};
		this.add(new AjaxSelfUpdatingTimerBehavior(Duration.seconds(5)));

		add(propView2);
	}

}
