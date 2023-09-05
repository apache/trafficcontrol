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

package org.apache.traffic_control.traffic_router.core.util;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

public class TrafficOpsUtils {
	public static final String TO_API_VERSION = "5.0";

	private String username;
	private String password;
	private String hostname;
	private String cdnName;
	private JsonNode config;

	public String replaceTokens(final String input) {
		return input.replace("${tmHostname}", this.getHostname()).replace("${toHostname}", this.getHostname()).replace("${cdnName}", getCdnName());
	}

	public String getUrl(final String parameter) throws JsonUtilsException {
		return replaceTokens(JsonUtils.getString(config, parameter));
	}

	public String getUrl(final String parameter, final String defaultValue) {
		return config != null ? replaceTokens(JsonUtils.optString(config, parameter, defaultValue)) : defaultValue;
	}

	public JsonNode getAuthJSON() throws IOException {
		final Map<String, String> authMap = new HashMap<>();
		authMap.put("u", getUsername());
		authMap.put("p", getPassword());

		final ObjectMapper mapper = new ObjectMapper();
		final JsonNode data = mapper.valueToTree(authMap);

		return data;
	}

	public String getAuthUrl() {
		return getUrl("api.auth.url", "https://${toHostname}/api/"+TO_API_VERSION+"/user/login");
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

	public void setConfig(final JsonNode config) {
		this.config = config;
	}

	public long getConfigLongValue(final String name, final long defaultValue) {
		return JsonUtils.optLong(config, name, defaultValue);
	}
}
