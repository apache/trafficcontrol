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

package com.comcast.cdn.traffic_control.traffic_monitor;

import java.text.DecimalFormat;
import java.text.NumberFormat;
import java.text.SimpleDateFormat;
import java.util.Date;

import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;
import org.apache.log4j.Logger;
import org.apache.wicket.Component;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.markup.html.IHeaderContributor;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.model.Model;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.publish.Stats;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.MultiUpdatingTimerBehavior;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.CacheListPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.DsListPanel;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.components.EventLogPanel;

public class Index extends MonitorPage implements IHeaderContributor {
	private static final Logger LOGGER = Logger.getLogger(Index.class);
	private static final long serialVersionUID = 1L;

	final public static NumberFormat NUMBER_FORMAT = new DecimalFormat("#,###.00");

	public Index() {
		final Behavior updater = new MultiUpdatingTimerBehavior(Duration.seconds(1));

		final Label servers_count = new Label("servers_count", getServerListSizeModel());
		servers_count.setOutputMarkupId(true);
		add(servers_count);

		final Label servers_down = new Label("servers_down", getServersDownModel());
		
		servers_down.add(updater);
		add(servers_down);

		final Label totalBandwidth = new Label("totalBandwidth",getCachesTotalBandwidthModel());
		totalBandwidth.add(updater);
		add(totalBandwidth);

		final Label totalBandwidthAvailable = new Label("totalBandwidthAvailable", getCachesTotalMaxBandwidthModel());
		totalBandwidthAvailable.add(updater);
		add(totalBandwidthAvailable);

		add(new Label("version", new Model<String>(getVersionStr())));

		final Label source = new Label("source", getSourceModel());
		source.add(updater);
		add(source);

		final Component[] updateList = new Component[] {servers_count};

		final Label servers_available = new Label("servers_available",getServersAvailableModel());

		servers_available.add(updater);
		add(servers_available);

		add(new CacheListPanel("serverList", updater, updateList));
		add(new EventLogPanel("eventLog"));

		add(new DsListPanel("dsList", updater, updateList));
	}

	private Model<Integer> getServerListSizeModel() {
		return new Model<Integer>() {
			@Override
			public Integer getObject() {
				return CacheStateRegistry.getInstance().size();
			}
		};
	}

	private Model<String> getServersDownModel() {
		return new Model<String>("") {

			@Override
			public String getObject() {
				return Integer.toString(CacheStateRegistry.getInstance().getCachesDownCount());
			}
		};
	}

	private Model<String> getCachesTotalBandwidthModel() {
		return new Model<String>("") {
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject( ) {
				return NUMBER_FORMAT.format(CacheStateRegistry.getInstance().getCachesBandwidthInKbps());
			}
		};
	}

	private Model<String> getCachesTotalMaxBandwidthModel() {
		return new Model<String>("") {
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject( ) {
				return NUMBER_FORMAT.format(CacheStateRegistry.getInstance().getCachesMaxBandwidthInKbps());
			}
		};
	}

	private Model<String> getSourceModel() {
		return new Model<String>(""){
			private static final long serialVersionUID = 1L;

			@Override
			public String getObject( ) {
				final MonitorConfig config = ConfigHandler.getInstance().getConfig();
				if(config == null) { return "[no config]"; }
				final String host = config.getEffectiveProps().get("tm.hostname");
				final String cdnName = config.getEffectiveProps().get("cdnName");
				return host+"/"+cdnName;
			}
		};
	}

	private String getVersionStr() {
		try {
			final JSONObject stats = Stats.getVersionInfo().getJSONObject("stats");
			final String name = stats.getString("name");
			final String version = stats.getString("version");
			final String revision = stats.getString("git-revision").replace("${buildNumber}","");
			String dateStr;

			try {
				dateStr = " (" + (new SimpleDateFormat("yyyy-MM-dd").format(new Date(stats.getLong("buildTimestamp")))) + ")";
			} catch (JSONException e) {
				dateStr = "(dev build)";
			}

			return name + "-" + version + "-" + revision + dateStr;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}

		return "";
	}

	private Model<String> getServersAvailableModel() {
		return new Model<String>("") {

			@Override
			public String getObject() {
				return Integer.toString(CacheStateRegistry.getInstance().getCachesAvailableCount());
			}
		};
	}
}

