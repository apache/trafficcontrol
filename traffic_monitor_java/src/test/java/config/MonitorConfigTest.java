package config;

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


import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.isEmptyString;

public class MonitorConfigTest {
	@Test
	public void itAddsInSubsFromJson() throws Exception {
		JSONObject jsonConfig = new JSONObject()
			.put("tm.hostname", "example.com")
			.put("cdnName", "kabletown")
			.put("tm.auth.username", "superman")
			.put("tm.auth.password", "kryptonite");

		MonitorConfig monitorConfig = new MonitorConfig(jsonConfig);
		JSONObject configDoc = monitorConfig.getConfigDoc();

		JSONObject json = configDoc.getJSONObject("tm.crConfig.json.polling.url");
		assertThat(json.optString("defaultValue"), equalTo("https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json"));
		assertThat(monitorConfig.getCrConfigUrl(), equalTo("https://example.com/CRConfig-Snapshots/kabletown/CRConfig.json"));
		assertThat(json.optString("value"), equalTo("https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json"));

		json = configDoc.getJSONObject("tm.healthParams.polling.url");
		assertThat(json.optString("defaultValue"), equalTo("https://${tmHostname}/health/${cdnName}"));
		assertThat(monitorConfig.getHeathUrl(), equalTo("https://example.com/health/kabletown"));
		assertThat(json.optString("value"), equalTo("https://${tmHostname}/health/${cdnName}"));

		json = configDoc.getJSONObject("tm.auth.url");
		assertThat(json.optString("defaultValue"), equalTo( "https://${tmHostname}/login"));
		assertThat(monitorConfig.getAuthUrl(), equalTo( "https://example.com/login"));
		assertThat(json.optString("value"), equalTo("https://${tmHostname}/login"));

		json = configDoc.getJSONObject("tm.auth.username");
		assertThat(json.optString("defaultValue"), isEmptyString());
		assertThat(monitorConfig.getAuthUsername(), equalTo("superman"));
		assertThat(json.optString("value"), equalTo("superman"));

		json = configDoc.getJSONObject("tm.auth.password");
		assertThat(json.optString("defaultValue"), isEmptyString());
		assertThat(monitorConfig.getAuthPassword(), equalTo("kryptonite"));
		assertThat(json.optString("value"), equalTo("**********"));
	}
}
