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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.io.File;
import java.io.FileReader;
import java.text.DecimalFormat;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.Event.EventType;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Updatable;

public class HealthDeterminer {
	private static final Logger LOGGER = Logger.getLogger(HealthDeterminer.class);

	public static final String IS_AVAILABLE_KEY = "isAvailable";
	public static final String STATUS = "status";
	public static final String ERROR_STRING = "error-string";
	public static final String NO_ERROR_FOUND = "No error found";

	private JSONObject profiles;
	private JSONObject deliveryServices;

	public enum AdminStatus {
		ONLINE, OFFLINE, REPORTED, ADMIN_DOWN, STANDBY
	}

	public Updatable getUpdateHandler() {
		return new Updatable() {

			@Override
			public boolean update(final File newDB) {
				LOGGER.debug("enter: "+newDB);
				try {
					final String str = IOUtils.toString(new FileReader(newDB));
					final JSONObject o = new JSONObject(str);
					return update(o);
				} catch (Exception e) {
					LOGGER.warn("error on update: "+newDB, e);
					return false;
				}
			}
			public boolean update(final JSONObject o) throws JSONException {
				profiles = o.getJSONObject("profiles");
				deliveryServices = o.optJSONObject("deliveryServices");
				LOGGER.warn(o.toString(2));
				return true;
			}
		};
	}
	public boolean shouldMonitor(final Cache cache) {
		final String profile = cache.getProfile();
		final String type = cache.getType();
		if(profiles == null) { return false; }
		final JSONObject set = profiles.optJSONObject(type);
		if(set == null) { return false; }
		if(!set.has(profile)) { return false; }
//		final JSONObject controls = set.optJSONObject(profile);
		cache.setControls(this);
		return true;
	}
	public boolean shouldMonitor(final JSONObject o) throws JSONException {
		final String profile = o.getString("profile");
		final String type = o.getString("type"); 
		final JSONObject set = profiles.optJSONObject(type);
		if(set == null || !set.has(profile)) { return false; }
		return true;
	}
	public void setControls(final Cache c) {
		final String profile = c.getProfile();
		final String type = c.getType();
		if(profiles == null) { return; }
		final JSONObject set = profiles.optJSONObject(type);
		if(set == null) { return; }
		if(!set.has(profile)) { return; }
//		final JSONObject controls = set.optJSONObject(profile);
		c.setControls(this);
	}
	public JSONObject getControls(final Cache c) {
		final String profile = c.getProfile();
		final String type = c.getType();
		if(profiles == null) { return null; }
		final JSONObject set = profiles.optJSONObject(type);
		if(set == null) { return null; }
		if(!set.has(profile)) { return null; }
		return set.optJSONObject(profile);
	}

	public String getIp(final Cache c) {
		final CacheState state = c.getState();
		if(state != null) {
			final String ip = state.getLastValue("resolved-ip");
			if(ip != null) { return ip; }
		}
		return c.getIpAddress();
	}

	public static boolean getIsAvailable(final Cache c, final boolean isHealthy) {
		final String status = c.getStatus();
		try {
			switch(AdminStatus.valueOf(status)) {
			case ONLINE: return true;
			case ADMIN_DOWN: return false;
			case OFFLINE: return false;
			case REPORTED: return isHealthy;
			case STANDBY: return false;
			default: return true;
			}
		} catch(IllegalArgumentException e) {
			return false;
		}
	}
	public static boolean getIsAvailable(final String status, final boolean isHealthy) {
		try {
			switch(AdminStatus.valueOf(status)) {
			case ONLINE: return true;
			case ADMIN_DOWN: return false;
			case OFFLINE: return false;
			case REPORTED: return isHealthy;
			case STANDBY: return false;
			default: return true;
			}
		} catch(IllegalArgumentException e) {
			return false;
		}
	}
	public void setIsAvailable(final Cache cache, final CacheState state) {
		// first check ONLINEness
		final String status = cache.getStatus();
		final String error = getErrorString(cache, state);
		state.putDataPoint(STATUS, status);
		state.putDataPoint(ERROR_STRING, error);
		final boolean isHealthy = (error == null);
		final EventType type = EventType.CACHE_STATE_CHANGE;
		type.setType(cache.getType());
		state.setAvailable(type, getIsAvailable(cache, isHealthy), error);
	}

	public void setIsAvailable(final Cache cache, final String e, final CacheState state) {
		final String status = cache.getStatus();
		state.putDataPoint(STATUS, status);
		state.putDataPoint(ERROR_STRING, e);
		final EventType type = EventType.CACHE_STATE_CHANGE;
		type.setType(cache.getType());
		state.setAvailable(type, getIsAvailable(cache, false), e);
	}

	private boolean shouldClearData(final String status) {
		try {
			switch(AdminStatus.valueOf(status)) {
			case ONLINE: return true;
			case ADMIN_DOWN: return false;
			case OFFLINE: return true;
			case REPORTED: return false;
			default: return false;
			}
		} catch(IllegalArgumentException e) {
			return true;
		}

	}
	String getErrorString(final Cache cache, final CacheState state) {
		if(shouldClearData(cache.getStatus())) {
			state.putDataPoint("clearData", "true");
			return null;
		}

		// this is where all the intelligence goes
		final String loadStr = state.getLastValue("system.proc.loadavg");
		final String loadavg = loadStr.split(" ")[0];
		state.putDataPoint("loadavg", loadavg);

		final String str = state.getLastValue("system.proc.net.dev");
		String tx_bytes= "0";
		//		String rx_bytes= "0";
		if(str == null) {
			LOGGER.warn("system.proc.net.dev missing on: "+cache.getHostname());
		} else {
			for(String line : str.split("\\n")) {
				line = line.replace(":", " ").trim();
				final String[] parts = line.split("\\s+");
				if(parts.length < 11) { continue; }
				if(parts[0].equals(cache.getInterfaceName())) {
					tx_bytes=parts[9];
					//				rx_bytes=parts[1];
				}
			}
		}
		final Bandwidth currentTx = new Bandwidth(tx_bytes);

		final long speed = state.getLong("system.inf.speed");
		final long maxBW = speed * Bandwidth.BITS_IN_KBPS;
		//		if (BandwidthHALF_DUPLEX.equalsIgnoreCase(mode)) {
		//			maxBW = (maxBW / 2);
		//		}

		final double kbps = calculateCurrentBandwidth(cache.previousTx, currentTx);
		cache.previousTx = currentTx;
		final double availBandwidthKbps = (double) maxBW - kbps;
		final double availBandwidthMbps = availBandwidthKbps / 1000.0;

		DecimalFormat df = new DecimalFormat("0.00");
		state.putDataPoint("kbps", df.format(kbps));
		state.putDataPoint("bandwidth", df.format(kbps));
		state.putDataPoint("maxKbps", Long.toString(maxBW));
		state.putDataPoint("availableBandwidthInKbps", df.format(availBandwidthKbps));
		state.putDataPoint("availableBandwidthInMbps", df.format(availBandwidthMbps));

		return mapControlsToError(cache.getControls(), state, "");
	}
	private static String mapControlsToError(final JSONObject controls, final AbstractState state, final String propBase) {
		if(controls == null) { return null; }
		final String[] keys = JSONObject.getNames(controls);
		for(String key : keys) {
			try {
				if(!key.startsWith("health.threshold.")) { continue; }
				String value = controls.optString(key);
				key = key.replace("health.threshold.", "");
				key = propBase+key;
				boolean greater = false;
				if(value.startsWith(">")) {
					value = value.replace(">", "");
					greater = true;
				}
				double cv = 0.0;
				try {
					cv = Double.parseDouble(value);
				} catch (Exception e) {cv = 0;}
				final String vstr = state.getLastValue(key);
				if(vstr == null) {
					continue;
				}
				//			try {
				final double v = Double.parseDouble(vstr);
				//		} catch (Exception e) {
				//			return 0;
				//		}
				//			final double v = state.getDouble(key);
				if(!greater) {
					if(v > cv) {
						return String.format("%s too high (%f > %f)", key, v, cv);
					}
				} else {
					if(v < cv) {
						return String.format("%s too low (%f < %f)", key, v, cv);
					}
				}
			} catch(Exception e) {
				LOGGER.warn(e,e);
			}
		}
		//			health.threshold.availableBandwidthInMbps: ">200"
		//			health.polling.url: "http://${hostname}/_astats?application=&inf.name=${interface_name}"
		//			health.threshold.queryTime: "500"
		//			history.count: "30"
		//			health.threshold.loadavg: "8.0"
		return null;
	}
	double calculateCurrentBandwidth(final Bandwidth prev, final Bandwidth curr) {
		double currBW = 0.0;
		if (prev != null) {
			currBW = prev.calculateKbps(curr);
		}
		return currBW;
	}


	public JSONObject getJSONStats(final Cache cache, final boolean peerOptimistic, final boolean raw) throws JSONException {
		final JSONObject statsJson = new JSONObject();
		final boolean isAvailableKnown = cache.isAvailableKnown();
		final boolean isAvailable = cache.isAvailable();

		if (!raw && peerOptimistic && PeerState.isCacheAvailableOnAnyPeer(cache)) {
			statsJson.put(IS_AVAILABLE_KEY, getIsAvailable(cache, true)); // ensure status overrides peer
			return statsJson;
		}

		if (isAvailableKnown) {
			statsJson.put(IS_AVAILABLE_KEY, isAvailable);
		} else {
			statsJson.put(IS_AVAILABLE_KEY, "unknown");
		}

		if (raw) {

			String error = null;
			String status = cache.getStatus();

			if (cache.getState() != null) {
				error = cache.getState().getLastValue(ERROR_STRING);
				status = cache.getState().getLastValue(STATUS);
			}

			if (error == null) {
				error = NO_ERROR_FOUND;
			}

			statsJson.put(ERROR_STRING, error);
			statsJson.put(STATUS, status);
		}

		return statsJson;
	}

	public int getConnectionTimeout(final Cache cache, final int d) {
		final JSONObject jo = cache.getControls();
		if(jo == null) { return d; }
		final int r = jo.optInt("health.connection.timeout");
		if(r == 0) { return d; }
		return r;
	}
	public JSONObject getDsControls(final String id) {
		if(deliveryServices == null) {
			return null;
		}
		return deliveryServices.optJSONObject(id);
	}
	public static void setIsAvailable(final DsState dsState, final JSONObject dsControls) {
		final EventType type = EventType.DELIVERY_SERVICE_STATE_CHANGE;

		if (dsControls == null) {
			dsState.putDataPoint(STATUS, "ONLINE");
			dsState.setAvailable(type, getIsAvailable("ONLINE", true), null);
			return;
		}

		// first check ONLINEness
		final String status = dsControls.optString(STATUS);
		final String error = getErrorString(dsControls, dsState);
		dsState.putDataPoint(STATUS, status);
		dsState.putDataPoint(ERROR_STRING, error);
		final boolean isHealthy = (error == null);
		dsState.setAvailable(EventType.DELIVERY_SERVICE_STATE_CHANGE, getIsAvailable(status, isHealthy), error);
	}
	private static String getErrorString(final JSONObject dsControls, final DsState dsState) {
		return mapControlsToError(dsControls, dsState, "");
	}
	public static boolean setIsLocationAvailable(final DsState dsState, final EmbeddedStati loc, final JSONObject dsControls) {
		boolean isAvailable = true;
		String error = null;
		if(dsControls != null) {
			final JSONObject locControlSet = dsControls.optJSONObject("locations");
			if(locControlSet != null) {
				final JSONObject locControls = locControlSet.optJSONObject(loc.getId());
				if(locControls != null) {
					error = getErrorString(locControls, loc, dsState);
				}
			}
		}
		if(error!=null) {
			isAvailable = false;
		}
		dsState.putDataPoint("location."+loc.getId()+"."+ERROR_STRING, error );
		dsState.putDataPoint("location."+loc.getId()+"."+IS_AVAILABLE_KEY, String.valueOf(isAvailable));
		return isAvailable;
	}
	private static String getErrorString(final JSONObject locControls,
			final EmbeddedStati loc, final DsState dsState) {
		return mapControlsToError(locControls, dsState, "location."+loc.getId()+".");
	}

}
