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

import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

public class MonitorConfig extends Config  {
	private static final long serialVersionUID = 1L;

	public MonitorConfig() {
	}
	public MonitorConfig(final JSONObject jsonObject) throws JSONException {
		super(jsonObject);
	}
//	public String getDataServerUrl() {
//		return getPropertyString("tm.dataServer.polling.url", "https://${tmHostname}/dataserver", "Url for the list of servers (json)");
//	}
	public String getCrConfigUrl() {
		return getPropertyString("tm.crConfig.json.polling.url", "https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json", "Url for the cr-config (json)");
	}
	public String getHeathUrl() {
		return getPropertyString("tm.healthParams.polling.url", "https://${tmHostname}/health/${cdnName}", "The url for the heath params (json)");
	}
	public String getAuthUrl() {
		return getPropertyString("tm.auth.url", "https://${tmHostname}/login", "The url for the authentication form");
	}
	public String getAuthUsername() {
		return getPropertyString("tm.auth.username", null, "The username for the authentication form");
	}
	public String getAuthPassword() {
		return getPropertyString("tm.auth.password", null, "The password for the authentication form");
	}

	public Long getTmFrequency() {
		return getLong("tm.polling.interval", 10000, "The polling frequency for getting updates from TM");
	}
	@Override
	protected String completePropString(final String pattern) {
		if(pattern == null) { return null; }
		final String tmHostname = getString("tm.hostname", null, "TM hostname");
		final String cdnName = getString("cdnName", null, "Cluster/CDN name");
		return pattern.replace("${tmHostname}", tmHostname).replace("${cdnName}", cdnName);
	}
	public int getEventLogCount() {
		return getInt("health.event-count", 200, "The number of historical events that will be kept");
	}
	public int getHealthPollingInterval() {
		return getInt("health.polling.interval", 5000, "The polling frequency for getting the states from caches");
	}
	public long getHealthDsInterval() {
		return getInt("health.ds.interval", 1000, "The polling frequency for calculating the deliveryService states");
	}
	public long getDsCacheLeniency() {
		return getInt("health.ds.leniency", 30000, "The amount of time before the deliveryService disregards the last update from a non-responsive cache");
	}
	public String getTmHost() {
		return getString("tm.hostname", null, "Hostname of TM");
	}
	public boolean shouldForceSystemExit() {
		return getBool("hack.forceSystemExit", false, "Call System.exit on shutdown");
	}
	public String getPeerUrl() {
		return getString("peers.polling.url", "http://${hostname}/publish/CrStates?raw", "The url for current, unfiltered states from peer monitors");
	}
	public long getPeerPollingInterval() {
		return getInt("peers.polling.interval", 5000, "Polling frequency for getting states from peer monitors");
	}
	public int getPeerThreadPool() {
		return getInt("peers.threadPool", 1, "The number of threads given to the pool for querying peers");
	}
	public int getConnectionTimeout() {
		return getInt("default.connection.timeout", 2000, "Default connection time for all queries (cache, peers, TM)");
	}
	public int getCacheTimePad() {
		return getInt("health.timepad", 10, "A delay between each separate cache query");
	}
	@SuppressWarnings("PMD")
	public boolean getPeerOptimistic() {
		return getBool("hack.peerOptimistic", true, "The assumption of a caches availability when unknown by peers");
	}
	@SuppressWarnings("PMD")
	public boolean getPublishDsStates() {
		return getBool("hack.publishDsStates", true, "If true, the delivery service states will be included in the CrStates.json");
	}
	public String getAccessControlAllowOrigin() {
		return getString("default.accessControlAllowOrigin", "*", 
				"The value for the header: Access-Control-Allow-Origin for published jsons... should be narrowed down to TMs");
	}
	public int getStartupMinCycles() {
		return getInt("health.startupMinCycles", 2, "The number of query cycles that must be completed before this Traffic Monitor will start reporting");
	}
	public boolean allowConfigEdit() {
		return getBool("allow.config.edit", false, "Allow the running configuration to be edited through the UI");
	}

	@SuppressWarnings("PMD")
	public static void main(final String[] args) throws JSONException {
		final JSONObject doc = ConfigHandler.getConfig().getConfigDoc();
		System.out.println(doc.toString(2));
	}
}
