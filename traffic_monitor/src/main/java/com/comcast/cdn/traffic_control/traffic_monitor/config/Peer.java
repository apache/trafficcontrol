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

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerState;

public class Peer implements Serializable {
	private static final long serialVersionUID = 1L;
	final private String hostname;
	private PeerState state;
	private String error;
	final private JSONObject json;
	final private Map<String, String> headerMap = new HashMap<String, String>();

	public Peer(final String id, final JSONObject o) throws JSONException {
		this.json = o;
		this.hostname = id;

		if (json.has("fqdn")) {
			headerMap.put("Host", json.getString("fqdn") + ":" + getPortString());
		}
	}

	public String getHostname() {
		return hostname;
	}

	public String getId() {
		return getHostname();
	}

	public PeerState getState() {
		return state;
	}

	public void setState(final PeerState state) {
		this.state = state;
	}

	public String getError() {
		return error;
	}

	public void setError(final String error) {
		this.error = error;
	}

	public String getProfile() {
		return json.optString("profile");
	}

	public String getFqdn() {
		return json.optString("fqdn");
	}

	public String getIpAddress() {
		return json.optString("ip");
	}

	public String getStatus() {
		return json.optString("status");
	}

	public String getLocation() {
		return json.optString("location");
	}

	public String getIp6Address() {
		return json.optString("ip6");
	}

	final public int getPort() {
		final int port = json.optInt("port");

		if (port == 0) {
			return(80);
		}

		return port;
	}

	final public String getPortString() {
		return String.valueOf(getPort());
	}

	public String toString() {
		return getHostname();
	}

	public Map<String, String> getHeaderMap() {
		return headerMap;
	}
}
