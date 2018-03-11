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


import com.comcast.cdn.traffic_control.traffic_monitor.config.Config;

import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Test;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class ConfigTest {

	@Test
	public void itReturnsEmptyConfig() throws JSONException {
		Config config = new Config();
		assertThat(config.getConfigDoc().keys().hasNext(), equalTo(false));
		assertThat(config.getEffectiveProps().size(), equalTo(0));
	}

	@Test
	public void itReturnsDefaultValues() throws JSONException {
		Config config = new Config();

		assertThat(config.getString("foo", "default", "a description"), equalTo("default"));
		assertThat(config.getBool("some boolean", false, "boolean value"), equalTo(false));
		assertThat(config.getInt("an integer", -1, "integer property"), equalTo(-1));
		assertThat(config.getLong("a long", 1000000L, "long property"), equalTo(1000000L));

		JSONObject jsonConfig = config.getConfigDoc();

		assertThat(jsonConfig.getJSONObject("foo").getString("value"), equalTo("default"));
		assertThat(jsonConfig.getJSONObject("foo").getString("defaultValue"), equalTo("default"));
		assertThat(jsonConfig.getJSONObject("foo").getString("description"), equalTo("a description"));
		assertThat(jsonConfig.getJSONObject("foo").getString("type"), equalTo("propString"));

		assertThat(jsonConfig.getJSONObject("some boolean").getString("type"), equalTo("boolean"));
		assertThat(jsonConfig.getJSONObject("some boolean").getString("value"), equalTo("false"));
		assertThat(jsonConfig.getJSONObject("an integer").getString("type"), equalTo("integer"));
		assertThat(jsonConfig.getJSONObject("an integer").getString("value"), equalTo("-1"));
		assertThat(jsonConfig.getJSONObject("a long").getString("type"), equalTo("Long"));
		assertThat(jsonConfig.getJSONObject("a long").getString("value"), equalTo("1000000"));

		assertThat(config.getEffectiveProps().size(), equalTo(0));
	}

	@Test
	public void itGetsBackNewDefaultAndPreservesPropertiesDocDefaultValue() throws JSONException {
		Config config = new Config();

		assertThat(config.getString("foo", "originaldefault", "a description"), equalTo("originaldefault"));
		assertThat(config.getString("foo", "somethingelse", "a description"), equalTo("somethingelse"));
		assertThat(config.getConfigDoc().getJSONObject("foo").getString("defaultValue"), equalTo("originaldefault"));
	}

	@Test
	public void itGetsValues() throws JSONException {
		JSONObject properties = new JSONObject();
		properties.put("foo", "bar");
		properties.put("some boolean", true);
		properties.put("an integer", 1);
		properties.put("a long", 100L);

		Config config = new Config(properties);

		assertThat(config.getString("foo", "default", "a description"), equalTo("bar"));
		assertThat(config.getBool("some boolean", false, "boolean value"), equalTo(true));
		assertThat(config.getInt("an integer", -1, "integer property"), equalTo(1));
		assertThat(config.getLong("a long", 1000000L, "long property"), equalTo(100L));

		assertThat(config.getEffectiveProps().get("foo"), equalTo("bar"));
	}

	@Test
	public void itPreservesOriginalValuesWhenOverlayIsNull() throws JSONException {
		JSONObject properties = new JSONObject();
		properties.put("foo", "bar");
		properties.put("some boolean", true);
		properties.put("an integer", 1);
		properties.put("a long", 100L);

		Config config = new Config(properties);
		config.update(null);

		assertThat(config.getString("foo", "default", "a description"), equalTo("bar"));
		assertThat(config.getBool("some boolean", false, "boolean value"), equalTo(true));
		assertThat(config.getInt("an integer", -1, "integer property"), equalTo(1));
		assertThat(config.getLong("a long", 1000000L, "long property"), equalTo(100L));
	}

	@Test
	public void itReturnsUpdatedValues() throws JSONException {
		JSONObject properties = new JSONObject();
		properties.put("foo", "bar");
		properties.put("some boolean", true);
		properties.put("an integer", 1);
		properties.put("a long", 100L);

		Config config = new Config(properties);

		JSONObject updatedProperties = new JSONObject()
			.put("foo", "something new")
			.put("some boolean", false)
			.put("an integer", 1234)
			.put("a long", 4321L);

		config.update(updatedProperties);

		assertThat(config.getString("foo", "default", "a description"), equalTo("something new"));
		assertThat(config.getBool("some boolean", false, "boolean value"), equalTo(false));
		assertThat(config.getInt("an integer", -1, "integer property"), equalTo(1234));
		assertThat(config.getLong("a long", 1000000L, "long property"), equalTo(4321L));

		assertThat(config.getEffectiveProps().get("foo"), equalTo("something new"));
	}
}
