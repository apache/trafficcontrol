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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

public class Config implements java.io.Serializable {
	private static final Logger LOGGER = Logger.getLogger(Config.class);
	private static final long serialVersionUID = 1L;

	private JSONObject baseProps = new JSONObject();
	private JSONObject props = new JSONObject();
	private final JSONObject propDoc = new JSONObject();

	public Config() {

	}

	public Config(final JSONObject properties) {
		baseProps = properties;
		props = new JSONObject(baseProps, JSONObject.getNames(baseProps));
	}

	public void update(final JSONObject overlayJson) throws JSONException {
		if (overlayJson == null) {
			LOGGER.warn("Skipping NULL overlay");
			props = baseProps;
			return;
		}

		LOGGER.info("update, adding: " + overlayJson.toString(2));
		final JSONObject properties = new JSONObject(baseProps, JSONObject.getNames(baseProps));

		Iterator<String> names = overlayJson.keys();
		while (names.hasNext()) {
			String name = names.next();
			properties.put(name, overlayJson.get(name));
		}

		props = properties;
	}

	public JSONObject getConfigDoc() {
		return propDoc;
	}

	@SuppressWarnings("unchecked")
	public Map<String,String> getEffectiveProps() {
		final Iterator<String> keys = props.keys();
		final Map<String, String> effectiveProperties = new HashMap<String, String>();

		while (keys.hasNext()) {
			String key = keys.next();
			effectiveProperties.put(key, props.optString(key));
		}

		return effectiveProperties;
	}

	public String getString(final String key, final String defaultValue, final String description) {
		String value = props.has(key) ? props.optString(key) : defaultValue;
		updatePropertiesDoc(key, defaultValue, value, description, "propString");
		return value;
	}

	public Long getLong(final String key, final long defaultValue, final String description) {
		long value = props.has(key) ? props.optLong(key) : defaultValue;
		updatePropertiesDoc(key, defaultValue, value, description, "Long");
		return value;
	}

	public int getInt(final String key, final int defaultValue, final String description) {
		int value = props.has(key) ? props.optInt(key) : defaultValue;
		updatePropertiesDoc(key, defaultValue, value, description, "integer");
		return value;
	}

	public boolean getBool(final String key, final boolean defaultValue, final String description) {
		boolean value = props.has(key) ? props.optBoolean(key) : defaultValue;
		updatePropertiesDoc(key, defaultValue, value, description, "boolean");
		return value;
	}

	private void updatePropertiesDoc(final String key, final Object defaultValue, final Object value, final String description, final String type) {
		if (!propDoc.has(key)) {
			try {
				String docDefaultValue = (defaultValue != null) ? String.valueOf(defaultValue) : "";
				JSONObject json = new JSONObject().put("defaultValue", docDefaultValue).put("description", description).put("type", type);
				propDoc.put(key, json);
			} catch (JSONException e) {
				LOGGER.warn(e,e);
			}
		}

		try {
			String s = String.valueOf(String.valueOf(value));
			s = (key.toLowerCase().contains("password")) ? "**********" : s;
			propDoc.getJSONObject(key).put("value", s);
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
	}
}
