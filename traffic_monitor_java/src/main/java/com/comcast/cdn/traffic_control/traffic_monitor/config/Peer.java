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

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerState;

public class Peer implements Serializable {
	private static final long serialVersionUID = 1L;
	final private String hostname;
	final private Map<String, String> headerMap = new HashMap<String, String>();
	final private String fqdn;
	final private String ip;
	final private String status;
	final private String location;
	final private int port;

	private PeerState state;
	private String error;

	public Peer(final String hostname, final JSONObject json) throws JSONException {
		this.hostname = hostname;

		fqdn = json.optString("fqdn");
		port = (json.optInt("port") == 0) ? 80 : json.optInt("port");

		if (fqdn != null && !fqdn.isEmpty()) {
			headerMap.put("Host", fqdn + ":" + getPortString());
		}

		ip = json.optString("ip");
		status = json.optString("status");
		location = json.optString("location");
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

	public String getFqdn() {
		return fqdn;
	}

	public String getIpAddress() {
		return ip;
	}

	public String getStatus() {
		return status;
	}

	public String getLocation() {
		return location;
	}

	final public int getPort() {
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
