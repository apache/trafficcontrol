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
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.extensions.ajax.markup.html.AjaxEditableChoiceLabel;
import org.apache.wicket.markup.html.WebMarkupContainer;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.Model;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;

public class ConfigPanel extends Panel {
//	private static final Logger LOGGER = Logger.getLogger(ConfigPanel.class);
	private static final long serialVersionUID = 1L;
	AjaxEditableChoiceLabel<String> cdnName;

	public ConfigPanel(final String id, final Behavior updater) {
		super(id);

		this.setOutputMarkupId(true);


		final WebMarkupContainer c = new WebMarkupContainer("configpanel");
		add(c);

		this.setOutputMarkupId(true);
		final MonitorConfig config = ConfigHandler.getConfig();

		final Map<String, String> effectiveProps = config.getEffectiveProps();
		final ListView<String> propView2 = new ListView<String>("propList", sort(effectiveProps.keySet())) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<String> item) {
				final String key = item.getModelObject();

				Label label = new Label("key", key);
				item.add(label);
				label = new Label("value", new ConfigModel(key));
				label.add(updater);
				item.add(label);
			}
		};

		c.add(propView2);

	}


	private List<String> sort(final Set<String> props) {
		final TreeSet<String> set = new TreeSet<String>(props);
 		return new ArrayList<String>(set);
	}

	static class ConfigModel extends Model<String> {
		private static final long serialVersionUID = 1L;
		final String key;
		public ConfigModel(final String key) {
			this.key = key;
		}
		@Override
		public String getObject( ) {
			final MonitorConfig config = ConfigHandler.getConfig();
			if(config == null) { return "[no config]"; }
			final String r = config.getEffectiveProps().get(key);
			if(r == null) { return "[null]"; }
			return r;
		}
	}
}

