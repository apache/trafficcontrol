package com.comcast.cdn.traffic_control.traffic_monitor.wicket.components;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import com.comcast.cdn.traffic_control.traffic_monitor.StatisticModel;
import com.comcast.cdn.traffic_control.traffic_monitor.MonitorPage;
import com.comcast.cdn.traffic_control.traffic_monitor.health.StateRegistry;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.MultiUpdatingTimerBehavior;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.util.time.Duration;

import java.util.ArrayList;
import java.util.List;

public abstract class StateDetailsPage extends MonitorPage {

	public StateDetailsPage(final String id, final String label) {
		final Behavior updater = new MultiUpdatingTimerBehavior(Duration.seconds(1));
		final Label wicketLabel = new Label(label, id);
		wicketLabel.add(updater);
		add(wicketLabel);

		final List<StatisticModel> statisticModels = new ArrayList<StatisticModel>();

		for (String key : getStateRegistry().get(id).getStatisticsKeys()) {
			statisticModels.add(new StatisticModel(key) {
				@Override
				public String getObject() {
					return (getStateRegistry().has(id)) ? getStateRegistry().get(id, getKey()) : "";
				}
			});
		}

		final ListView<StatisticModel> statesListView = new ListView<StatisticModel>("params", statisticModels) {
			private static final long serialVersionUID = 1L;
			@Override
			protected void populateItem(final ListItem<StatisticModel> item) {
				final StatisticModel statisticModel = item.getModelObject();

				final Label keyLabel = new Label("key", statisticModel.getKey());
				final Label valueLabel = new Label("value", statisticModel);

				valueLabel.add(updater);

				item.add(keyLabel);
				item.add(valueLabel);
			}
		};

		add(statesListView);
	}

	protected abstract StateRegistry getStateRegistry();
}
