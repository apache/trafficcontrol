/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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
import java.util.Arrays;
import java.util.Date;
import java.util.TimeZone;

import org.apache.wicket.ajax.AjaxSelfUpdatingTimerBehavior;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.markup.html.WebMarkupContainer;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.AbstractReadOnlyModel;
import org.apache.wicket.model.Model;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.publish.Stats;

public class StatusPanel extends Panel {
	private static final long serialVersionUID = 1L;

	public StatusPanel(final String id) {
		super(id);
		final WebMarkupContainer c = new WebMarkupContainer("container");
		add(c);

		//		final Model<String> count = new Model<String>("0");
		c.add(new Clock("clock", TimeZone.getTimeZone("America/Denver")));
		c.add(new AjaxSelfUpdatingTimerBehavior(Duration.seconds(1)));

		final ListView<Model<String>> props = new ListView<Model<String>>("props", CacheWatcher.getProps()) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<Model<String>> item) {
				final Model<String> val = item.getModelObject();
				item.add(new Label("value",val));
			}
		};
		c.add(props);

		final JSONObject stats = Stats.getVersionInfo().optJSONObject("stats");
		final String[] keys = (stats == null || stats.length() == 0)? new String[0] : JSONObject.getNames(stats);
		final ListView<String> props2 = new ListView<String>("versionInfo", Arrays.asList(keys)) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<String> item) {
				final String key = item.getModelObject();
				item.add(new Label("key", key));
				item.add(new Label("value", new Model<String>() {
					private static final long serialVersionUID = 1L;
					@Override
					public String getObject( ) {
						JSONObject stats = Stats.getVersionInfo().optJSONObject("stats");
						return stats.optString(key);
					}
				}));
			}
		};
		c.add(props2);
	}
	static class Clock extends Label {
		private static final long serialVersionUID = 1L;
		public Clock(final String id, final TimeZone tz) {
			super(id, new ClockModel(tz));
		}
		private static class ClockModel extends AbstractReadOnlyModel<String> {
			private static final long serialVersionUID = 1L;
			private final DateFormat df;

			public ClockModel(final TimeZone tz) {
				df = DateFormat.getDateTimeInstance(DateFormat.FULL, DateFormat.LONG);
				df.setTimeZone(tz);
			}

			@Override
			public String getObject() {
				return df.format(new Date());
			}
		}
	}

}


