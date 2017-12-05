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

import java.io.File;
import java.io.IOException;

import org.apache.commons.io.FileUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.markup.MarkupType;
import org.apache.wicket.markup.html.WebPage;
import org.apache.wicket.request.mapper.parameter.PageParameters;

import com.comcast.cdn.traffic_control.traffic_monitor.health.TmListener;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Updatable;

public class CrConfig extends WebPage  {
	private static final Logger LOGGER = Logger.getLogger(CrConfig.class);
	private static final long serialVersionUID = 1L;

	public CrConfig() {
		//		String jsonData = "{\"jsonKey\":\"jsonValue\"}"; 
		//		IResource jsonResource = new ByteArrayResource("text/plain", jsonData.getBytes()); 
		//		IRequestHandler requestHandler = new ResourceRequestHandler(jsonResource, null); 
		//		requestHandler.respond(getRequestCycle()); 
	}

	/**
	 * The default markup is html. We want it to be xml.
	 */
	@Override
	public final MarkupType getMarkupType() {
		return new MarkupType("xml", "text/xml");
	}
	/**
	 * Send out the xml!!!!
	 */
	@Override
	public final void renderPage() {
		final PageParameters pageParameters = this.getPageParameters();
		if(pageParameters != null && pageParameters.getPosition("json") != -1) {
			getResponse().write(crConfigJsonStr);
			return;
		}
		getResponse().write(crConfigStr);
	}
	private static void setCrConfigJsonStr(final String string) {
		crConfigJsonStr = string;
	}

	static private String crConfigStr = "";
	static private String crConfigJsonStr = "";
	static private JSONObject crConfigJson = null;

	public static String getCrConfigStr() {
		return crConfigStr;
	}

	public static void setCrConfigStr(final String crConfigStr) {
		CrConfig.crConfigStr = crConfigStr;
	}

	public static JSONObject getCrConfigJson() {
		return crConfigJson;
	}

	public static void setCrConfigJson(final JSONObject crConfigJson) {
		CrConfig.crConfigJson = crConfigJson;
	}

	public static class UpdateHandler implements Updatable {
		@Override
		public boolean update(final File newDB) {
			try {
				CrConfig.crConfigStr = FileUtils.readFileToString(newDB);
				return true;
			} catch (IOException e) {
				LOGGER.warn(e,e);
			}
			return false;
		}
		@Override
		public boolean update(final JSONObject jsonObject) throws JSONException {
			return false;
		}
	}

	public static TmListener getCrConfigListener() {
		return new TmListener() {
			@Override
			public void handleCrConfig(final JSONObject o) {
				try {
					setCrConfigJsonStr(o.toString(2));
				} catch (JSONException e) {
					LOGGER.warn(e,e);
				}
			}
			};
	}

}

