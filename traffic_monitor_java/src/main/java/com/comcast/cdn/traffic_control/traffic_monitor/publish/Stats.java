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

import java.io.InputStream;
import java.util.List;
import java.util.Properties;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.MonitorApplication;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.CacheDataModel;

public class Stats extends JsonPage {
	private static final Logger LOGGER = Logger.getLogger(Stats.class);
	private static final long serialVersionUID = 1L;

	/**
	 * Send out the json!!!!
	 */
	@Override
	public JSONObject getJson(final PageParameters pp) throws JSONException {
			return getVersionInfo();
	}

	static Properties props;
	public static JSONObject getVersionInfo() {
		synchronized(LOGGER) {
			final JSONObject o = new JSONObject();
			try {
				final InputStream stream = Stats.class.getResourceAsStream("/version.prop");
				if(props == null) {
					props = new Properties();
					try {
						props.load(stream);
						stream.close();
					} catch (Exception e) {
						LOGGER.warn(e,e);
						props = null;
					}
				}
				props.put("uptime", Long.toString(MonitorApplication.getUptime()));
				
				final List<CacheDataModel> cwProps = CacheWatcher.getProps();
				for(CacheDataModel m : cwProps) {
					props.put(m.getKey(), m.getValue());
				}
				o.put("stats", props);

			} catch (JSONException e) {
				LOGGER.warn(e,e);
			}
			return o;
		}
	}
}



