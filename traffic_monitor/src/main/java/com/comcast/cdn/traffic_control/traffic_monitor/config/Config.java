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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.lang.reflect.Method;
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
	private JSONObject overlayProps = new JSONObject();
	private JSONObject props = new JSONObject();
	private final JSONObject propDoc = new JSONObject();
	private boolean hasForcedPropCalls = false;

	public Config() {
	}
	public Config(final JSONObject o) throws JSONException {
		baseProps = o;
		LOGGER.debug(o.toString(2));
		props = new JSONObject(baseProps, JSONObject.getNames(baseProps));
	}
	protected String completePropString(final String pattern) {
		return pattern;
	}

	public void userOverrideBaseConfig(final JSONObject o) throws JSONException {
		baseProps = o;
		props = overlay(overlayProps);
	}
	private JSONObject overlay(final JSONObject o) throws JSONException {
		overlayProps = o;
		if(overlayProps == null) { return baseProps; }
		final JSONObject myprops = new JSONObject(baseProps, JSONObject.getNames(baseProps));
		LOGGER.warn(o.toString(2));
		final String[] names = JSONObject.getNames(o);
		if(names != null) {
			for(String name : names) {
				myprops.put(name, o.get(name));
			}
		}
		return myprops;
	}
	public void update(final JSONObject o) throws JSONException {
		LOGGER.warn("update, adding: "+o.toString(2));
		props = overlay(o);
	}
	@SuppressWarnings("unchecked")
	public Map<String, String> getBaseProps() {
		final Iterator<String> nameItr = baseProps.keys();
		final Map<String, String> outMap = new HashMap<String, String>();
		while(nameItr.hasNext()) {
			final String key = nameItr.next();
			outMap.put(key, baseProps.optString(key));
		}
		return outMap;
	}
	@SuppressWarnings("unchecked")
	public Map<String,String> getEffectiveProps() {
		final Iterator<String> nameItr = props.keys();
		final Map<String, String> outMap = new HashMap<String, String>();
		while(nameItr.hasNext()) {
			final String key = nameItr.next();
			final String url = completePropString(props.optString(key));
			outMap.put(key, url);
		}
		return outMap;
	}
	public String[] getPropNames() {
		return JSONObject.getNames(props);
	}
	public Object getProp(final String key) {
		return props.opt(key);
	}


	public String getString(final String key, final String defaultValue, final String description) {
		String ret = defaultValue;
		if(!propDoc.has(key)) {
			putDoc(key,defaultValue,description,"propString");
		}
		if(props.has(key)) {
			ret = props.optString(key);
		}
		putLast(key,String.valueOf(ret));
		return ret;
	}
	public String getPropertyString(final String key, final String defaultValue, final String description) {
		String ret = defaultValue;
		if(!propDoc.has(key)) {
			putDoc(key,defaultValue,description,"propString");
		}
		if(props.has(key)) {
			ret = props.optString(key);
		}
		putLast(key,String.valueOf(ret));
		return completePropString(ret);
	}
	public Long getLong(final String key, final long defaultValue, final String description) {
		long ret = defaultValue;
		if(!propDoc.has(key)) {
			putDoc(key,String.valueOf(defaultValue),description,"Long");
		}
		if(props.has(key)) {
			ret = props.optLong(key);
		}
		putLast(key,String.valueOf(ret));
		return ret;
	}
	public int getInt(final String key, final int defaultValue, final String description) {
		int ret = defaultValue;
		if(!propDoc.has(key)) {
			putDoc(key,String.valueOf(defaultValue),description,"Long");
		}
		if(props.has(key)) {
			ret = props.optInt(key);
		}
		putLast(key,String.valueOf(ret));
		return ret;
	}
	public boolean getBool(final String key, final boolean defaultValue, final String description) {
		boolean ret = defaultValue;
		if(!propDoc.has(key)) {
			putDoc(key,String.valueOf(defaultValue),description,"Long");
		}
		if(props.has(key)) {
			ret = props.optBoolean(key);
		}
		putLast(key,String.valueOf(ret));
		return ret;
	}
	private void putLast(final String key, final String value) {
		try {
			if (key.toLowerCase().contains("password")) {
				propDoc.getJSONObject(key).put("value", "**********");
			} else {
				propDoc.getJSONObject(key).put("value", value);
			}
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
	}
	private void putDoc(final String key, final String defaultValue, final String description, final String type) {
		try {
			final JSONObject o = new JSONObject();
			o.put("defaultValue", defaultValue);
			o.put("description", description);
			o.put("type", type);
			propDoc.put(key, o);
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		}
	}
	public JSONObject getConfigDoc() {
		if(!hasForcedPropCalls) {
			hasForcedPropCalls = true;
			forcePropCalls();
		}
		return propDoc;
	}
	protected void forcePropCalls() {
		for(Method m : this.getClass().getMethods()) {
			try {
				final Class<?> rtype = m.getReturnType();
				final Class<?>[] ptypes = m.getParameterTypes();
				if(!(rtype.equals(void.class)) && ptypes.length==0) {
					m.invoke(this, new Object[]{});
				}
			} catch (Exception e) {
				LOGGER.warn(e,e);
			}
		}
	}

}
