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

package com.comcast.cdn.traffic_control.traffic_router.core.util;

import org.json.JSONException;
import org.json.JSONObject;

public class TrafficOpsUtils {
	private String username;
	private String password;
	private String hostname;
	private String cdnName;
	private JSONObject config;

	public String replaceTokens(final String input) {
		return input.replace("${tmHostname}", this.getHostname()).replace("${toHostname}", this.getHostname()).replace("${cdnName}", getCdnName());
	}

	public String getUrl(final String parameter) throws JSONException {
		return replaceTokens(config.getString(parameter));
	}

	public String getUrl(final String parameter, final String defaultValue) {
		return config != null ? replaceTokens(config.optString(parameter, defaultValue)) : defaultValue;
	}

	public JSONObject getAuthJSON() throws JSONException {
		final JSONObject data = new JSONObject();

		data.put("u", getUsername());
		data.put("p", getPassword());

		return data;
	}

	public String getAuthUrl() {
		return getUrl("api.auth.url", "https://${toHostname}/api/1.1/user/login");
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(final String username) {
		this.username = username;
	}

	public String getPassword() {
		return password;
	}

	public void setPassword(final String password) {
		this.password = password;
	}

	public String getHostname() {
		return hostname;
	}

	public void setHostname(final String hostname) {
		this.hostname = hostname;
	}

	public String getCdnName() {
		return cdnName;
	}

	public void setCdnName(final String cdnName) {
		this.cdnName = cdnName;
	}

	public void setConfig(final JSONObject config) {
		this.config = config;
	}

	public long getConfigLongValue(final String name, final long defaultValue) {
		return config != null ? config.optLong(name, defaultValue) : defaultValue;
	}
}
