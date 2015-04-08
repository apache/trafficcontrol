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

package com.comcast.cdn.traffic_control.traffic_monitor;

import java.text.DecimalFormat;
import java.text.NumberFormat;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

import org.apache.log4j.Logger;
import org.apache.wicket.Component;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.extensions.ajax.markup.html.tabs.AjaxTabbedPanel;
import org.apache.wicket.extensions.markup.html.tabs.AbstractTab;
import org.apache.wicket.markup.html.IHeaderContributor;
import org.apache.wicket.markup.html.WebMarkupContainer;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.model.Model;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.publish.Stats;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.MultiUpdatingTimerBehavior;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.CacheListPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.ConfigPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.DsListPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.EditConfigPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.EventLogPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.GraphPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.StatusPanel;

public class Index extends MonitorPage implements IHeaderContributor {
	private static final Logger LOGGER = Logger.getLogger(Index.class);
	private static final long serialVersionUID = 1L;

	final public static NumberFormat NUMBER_FORMAT = new DecimalFormat("#,###.00");

	public Index() {

		final Behavior updater = new MultiUpdatingTimerBehavior(Duration.seconds(1));//AjaxSelfUpdatingTimerBehavior


		final Model<Integer> serverListSize = getServerListSizeModel();
		final Label servers_count = new Label("servers_count", serverListSize);
		servers_count.setOutputMarkupId(true);
		add(servers_count);

		final Label servers_down = new Label("servers_down", getServersDownModel());
		
		servers_down.add(updater);
		add(servers_down);

		final Label totalBandwidth = new Label("totalBandwidth", getCacheStateSumModel("kbps"));
		totalBandwidth.add(updater);
		add(totalBandwidth);

		final Label totalBandwidthAvailable = new Label("totalBandwidthAvailable", getCacheStateSumModel("maxKbps"));
		totalBandwidthAvailable.add(updater);
		add(totalBandwidthAvailable);

		add(new Label("version", new Model<String>(getVersionStr())));

		final Label source = new Label("source", getSourceModel());
		source.add(updater);
		add(source);

		final Component[] updateList = new Component[] {servers_count}; //graph, 

		final Label servers_available = new Label("servers_available", getServersAvailableModel());
		servers_available.add(updater);
		add(servers_available);

		add(new CacheListPanel("serverList", updater, updateList));
		add(new EventLogPanel("eventLog"));

		add(new DsListPanel("dsList", updater, updateList));

		add(getTabbedPanel(updater));
	}
	private Model<Integer> getServerListSizeModel() {
		return new Model<Integer>() {
			private static final long serialVersionUID = 1L;
			@Override
			public Integer getObject( ) {
				return new Integer(CacheState.getCacheStates().size());
			}
		};
	}
	private Model<String> getServersDownModel() {
		return new Model<String>("") {
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject( ) {
				int cnt = 0;
				for(CacheState cs : CacheState.getCacheStates()) {
					//					CacheState cs = CacheState.get(server);
					if ( cs != null && cs.isError() ) { 
						cnt++;
					}
				}
				return String.valueOf(cnt);
			}
		};
	}
	private Model<String> getCacheStateSumModel(final String key) {
		return new Model<String>("") {
			private static final long serialVersionUID = 1L;
			@Override
			public String getObject( ) {
				long bw = 0;
				for(CacheState cs : CacheState.getCacheStates()) {
					bw += cs.getDouble(key);
				}
				return NUMBER_FORMAT.format(bw);
			}
		};
	}
	private Model<String> getSourceModel() {
		return new Model<String>(""){
			private static final long serialVersionUID = 1L;
			@Override
			public String getObject( ) {
				final MonitorConfig config = ConfigHandler.getConfig();
				if(config == null) { return "[no config]"; }
				final String host = config.getEffectiveProps().get("tm.hostname");
				final String cdnName = config.getEffectiveProps().get("cdnName");
				return host+"/"+cdnName;
			}
		};
	}
	protected static Model<String> getServersAvailableModel() {
		return new Model<String>("") {
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject( ) {
				int cnt = 0;
				for(CacheState cs : CacheState.getCacheStates()) {
					//					CacheState cs = CacheState.get(server);
					if ( cs != null && cs.isAvailable() ) { cnt++; }
				}
				return String.valueOf(cnt);
			}
		};
	}
	private String getVersionStr() {
		try {
			final JSONObject stats;
			stats = Stats.getVersionInfo().getJSONObject("stats");
			final String name = stats.getString("name");
			final String version = stats.getString("version");
			final String revision = stats.getString("git-revision").replace("${buildNumber}","");
			String dateStr = null;
			try {
				dateStr = " ("+(new SimpleDateFormat("yyyy-MM-dd").format(
						new Date(stats.getLong("buildTimestamp"))))+")";
			} catch (JSONException e) { 
				dateStr = "(dev build)"; 
			}
			return name+"-"+version+"-"+revision+dateStr;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
		return "";
	}
	private AjaxTabbedPanel<AbstractTab> getTabbedPanel(final Behavior updater) {
		final List<AbstractTab> tabs=new ArrayList<AbstractTab>();
		tabs.add(new AbstractTab(new Model<String>("Status")) {
			private static final long serialVersionUID = 1L;
			public WebMarkupContainer getPanel(final String panelId) {
				return new StatusPanel(panelId);
			}
		});
		tabs.add(new AbstractTab(new Model<String>("Config")) {
			private static final long serialVersionUID = 1L;
			public WebMarkupContainer getPanel(final String panelId) {
				return new ConfigPanel(panelId, updater);
			}
		});

		final MonitorConfig config = ConfigHandler.getConfig();
		if(config != null && config.allowConfigEdit()) {
			tabs.add(new AbstractTab(new Model<String>("Edit Config")) {
				private static final long serialVersionUID = 1L;
				public WebMarkupContainer getPanel(final String panelId) {
					return new EditConfigPanel(panelId);
				}
			});
		}
		return new AjaxTabbedPanel<AbstractTab>("tabs", tabs);
	}
}

