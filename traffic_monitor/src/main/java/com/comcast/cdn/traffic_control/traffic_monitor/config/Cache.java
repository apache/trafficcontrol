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

import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.health.Bandwidth;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class Cache implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
//	protected List<Property> property;
	protected String hostname;
	private CacheState state;
	final private JSONObject json;

	public Bandwidth previousTx;
	private String error;
	public Cache(final String id, final JSONObject o) throws JSONException {
		json = o;
//		hostname = json.getString("host_name");
//		json.getString("ip_address");
//		json.optString("ip6_address");
//		json.getString("interface_name");
//		json.getString("status");
//		json.getString("location");
//		json.getString("profile");
//		json.getString("domain_name");
//		json.getString("host_name");
//		json.getString("type"); 
//		json.getInt("tcp_port"); 
		hostname = id;
		json.getString("ip");
		json.optString("ip6");
//		json.getString("interfaceName");
		json.getString("status");
		json.getString("locationId");
		json.getString("profile");
		json.getString("fqdn");
		json.getString("type"); 
		json.getInt("port"); 
	}
	public String getHostname() {
		return hostname;
	}
	public void setHostname(final String hostname) {
		this.hostname = hostname;
	}
	public String toString() {
		return "Cache Server: " + hostname;
	}
	public void setError(final String b, final HealthDeterminer myHealthDeterminer) {
		error = b;
	}
	public String getError() {
		return error;
	}
	public String getIpAddress() {
		return json.optString("ip");
	}
	public String getInterfaceName() {
		return json.optString("interfaceName");
	}
	public String getStatus() {
		return json.optString("status");
	}
	public String getLocation() {
		return json.optString("locationId");
	}
	public void setState(final CacheState state, final HealthDeterminer myHealthDeterminer) {
		myHealthDeterminer.setIsAvailable(this, state);
		this.state = state;
	}
	public void setError(final CacheState state, final String e, final HealthDeterminer myHealthDeterminer) {
		myHealthDeterminer.setIsAvailable(this, e, state);
		this.state = state;
	}
	public CacheState getState() {
		return state;
	}
//	public String getDomainName() {
//		return json.optString("domain_name");
//	}
	public boolean isAvailableKnown() {
		if(state==null) { return false; }
		final String v = state.getLastValue("isAvailable");
		if(v == null) { return false; }
		return true;
	}
	public boolean isAvailable() {
		if(state==null) { return true; }
		final String v = state.getLastValue("isAvailable");
		if(v == null) { return true; }
		return Boolean.parseBoolean(v);
	}
	
	public String getQueryIp() {
		final String ip = json.optString("queryIp");
		if(ip != null && !ip.equals("")) { 
			return ip; 
		}
		return getIp();
	}
	public int getQueryPort() {
		if(json.has("queryPort")) {
			return json.optInt("queryPort");
		}
		return json.optInt("port");
	}
	public String getIp() {
//		if(state != null) {
//			final String ip = state.getLastValue("resolved-ip");
//			if(ip != null) { return ip; }
//		}
		return getIpAddress();
	}
	public String getType() {
		return json.optString("type");
	}
	public String getIp6() {
		return json.optString("ip6");
	}
	
//	JSONObject controls;
	HealthDeterminer healthDeterminer;
	public void setControls(final HealthDeterminer healthDeterminer) {
		this.healthDeterminer = healthDeterminer;
	}
	public JSONObject getControls() {
		if(healthDeterminer==null) {
			return null;
		}
		return healthDeterminer.getControls(this);
	}
	public int getHistoryCount() {
		return getControls().optInt("history.count");
	}
	public void setCacheState(final CacheState cacheState) {
		state = cacheState;
	}
	public long getHistoryTime() {
		return getControls().optInt("history.time");
	}
	public String getProfile() {
		return json.optString("profile");
	}
//	public JSONObject getJson() {
//		return json;
//	}
	public String getFqdn() {
		return json.optString("fqdn");
	}

	public JSONObject getDeliveryServices() {
		return json.optJSONObject("deliveryServices");
	}

	public boolean hasDeliveryServices() {
		return json.has("deliveryServices");
	}

	public List<String> getDeliveryServiceIds() {
		return new ArrayList(Arrays.asList(JSONObject.getNames(getDeliveryServices())));
	}

	public List<String> getFqdns(final String deliveryServiceId) throws JSONException {
		final ArrayList<String> fqdns = new ArrayList<String>();

		final JSONObject deliveryServices = getDeliveryServices();

		if (!deliveryServices.has(deliveryServiceId)) {
			fqdns.add(deliveryServices.getString(deliveryServiceId));
			return fqdns;
		}

		final JSONArray ja = deliveryServices.optJSONArray(deliveryServiceId);

		if (ja != null) {
			for (int i = 0; i < ja.length(); i++) {
				fqdns.add(ja.getString(i));
			}
			return fqdns;
		}

		fqdns.add(deliveryServices.optString(deliveryServiceId));

		return fqdns;
	}
}
