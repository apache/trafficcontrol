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

import java.util.ArrayList;
import java.util.List;

import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;
import org.apache.log4j.Logger;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.request.mapper.parameter.PageParameters;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.KeyValue;
import com.comcast.cdn.traffic_control.traffic_monitor.MonitorPage;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.MultiUpdatingTimerBehavior;

public class DsDetailsPage extends MonitorPage {
	private static final Logger LOGGER = Logger.getLogger(DsDetailsPage.class);
	private static final long serialVersionUID = 1L;

	public DsDetailsPage(final PageParameters pars) {
		this(pars.get("id").toString());
	}
	public DsDetailsPage(final String idStr) {
		final Behavior updater = new MultiUpdatingTimerBehavior(Duration.seconds(1));
		final Label hostname = new Label("id", idStr);
		hostname.add(updater);
		add(hostname);

		List<KeyValue> keyValues;
		try {
			keyValues = DeliveryServiceStateRegistry.getInstance().getModelList(idStr);
		} catch (Exception e) {
			LOGGER.warn(e,e);
			keyValues = new ArrayList<KeyValue>();
			keyValues.add(new KeyValue("Error", e.toString()));
		}
		final ListView<KeyValue> servers = new ListView<KeyValue>("params", keyValues ) {
			private static final long serialVersionUID = 1L;
			@Override
			protected void populateItem(final ListItem<KeyValue> item) {
				final KeyValue keyval = (KeyValue) item.getModelObject();
				item.add(new Label("key", keyval.getKey()));
				final Label v = new Label("value", keyval);
				v.add(updater);
				item.add(v);
			}
		};
		add(servers);
	}
}
