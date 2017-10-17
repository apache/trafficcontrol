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

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.markup.MarkupType;
import org.apache.wicket.markup.html.WebPage;
import org.apache.wicket.request.http.WebResponse;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;

public class JsonPage extends WebPage {
	private static final Logger LOGGER = Logger.getLogger(JsonPage.class);
	private static final long serialVersionUID = 1L;

	/**
	 * The default markup is html. We want it to be json.
	 */
	@Override
	public final MarkupType getMarkupType() {
		return new MarkupType("json", "application/json");
	}

	/**
	 * Send out the json!!!!
	 */
	@Override
	public final void renderPage() {
		final MonitorConfig config = ConfigHandler.getInstance().getConfig();
		final WebResponse response = (WebResponse)super.getResponse();
		response.setHeader("Access-Control-Allow-Origin", config.getAccessControlAllowOrigin() );
		try {
			final JSONObject jo = getJson(getPageParameters());
			if(jo == null) {
				response.sendError(503, "Not ready");
				return;
			}
			response.write(jo.toString(2));
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
	}
	protected static long getParam(final PageParameters pp, final String key, final long d) {
		try {
			return Long.parseLong(pp.get(key).toString());
		} catch(Exception e) {
			return d;
		}
	}


	protected JSONObject getJson(final PageParameters pp) throws JSONException {
		return new JSONObject();
	}
}
