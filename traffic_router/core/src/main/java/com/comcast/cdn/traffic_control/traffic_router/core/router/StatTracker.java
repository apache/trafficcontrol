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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoResult;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.RouteType;

@SuppressWarnings("PMD.ExcessivePublicCount")
public class StatTracker {
	private String dnsRoutingName;
	private String httpRoutingName;

	public static class Tallies {

		public int getCzCount() {
			return czCount;
		}
		public void setCzCount(final int czCount) {
			this.czCount = czCount;
		}
		public int getGeoCount() {
			return geoCount;
		}
		public void setGeoCount(final int geoCount) {
			this.geoCount = geoCount;
		}
		public int getDsrCount() {
			return dsrCount;
		}
		public int getMissCount() {
			return missCount;
		}
		public void setMissCount(final int missCount) {
			this.missCount = missCount;
		}
		public int getErrCount() {
			return errCount;
		}
		public void setErrCount(final int errCount) {
			this.errCount = errCount;
		}
		public int getStaticRouteCount() {
			return staticRouteCount;
		}
		public void setStaticRouteCount(final int staticRouteCount) {
			this.staticRouteCount = staticRouteCount;
		}

		public int getFedCount() {
			return fedCount;
		}

		public void setFedCount(final int fedCount) {
			this.fedCount = fedCount;
		}

		public int getRegionalDeniedCount() {
			return regionalDeniedCount;
		}
		public void setRegionalDeniedCount(final int regionalDeniedCount) {
			this.regionalDeniedCount = regionalDeniedCount;
		}
		public int getRegionalAlternateCount() {
			return regionalAlternateCount;
		}
		public void setRegionalAlternateCount(final int regionalAlternateCount) {
			this.regionalAlternateCount = regionalAlternateCount;
		}

		public int czCount;
		public int geoCount;
		public int missCount;
		public int dsrCount;
		public int errCount;
		public int staticRouteCount;
		public int fedCount;
		public int regionalDeniedCount;
		public int regionalAlternateCount;
	}

	public static class Track {
		public static enum RouteType {
			DNS,HTTP
		}

		public static enum ResultType {
			ERROR, CZ, GEO, MISS, STATIC_ROUTE, DS_REDIRECT, DS_MISS, INIT, FED, RGDENY, RGALT, GEO_REDIRECT
		}

		public enum ResultDetails {
			NO_DETAILS, DS_NOT_FOUND, DS_NO_BYPASS, DS_BYPASS, DS_CZ_ONLY, DS_CLIENT_GEO_UNSUPPORTED, GEO_NO_CACHE_FOUND,
			REGIONAL_GEO_NO_RULE, REGIONAL_GEO_ALTERNATE_WITHOUT_CACHE, REGIONAL_GEO_ALTERNATE_WITH_CACHE
		}

		long time;
		RouteType routeType;
		String fqdn;
		ResultType result = ResultType.ERROR;
		ResultDetails resultDetails = ResultDetails.NO_DETAILS;
		Geolocation resultLocation;
		
		Geolocation clientGeolocation; // the GEO info always retrieved from GEO DB, not from Cache Location
		boolean isClientGeolocationQueried;

		RegionalGeoResult regionalGeoResult;

		public Track() {
			start();
		}
		public String toString() {
			return fqdn+" - "+result;
		}
		public void setRouteType(final RouteType routeType, final String fqdn) {
			this.routeType = routeType;
			this.fqdn = fqdn;
		}
		public void setResult(final ResultType result) {
			this.result = result;
		}
		public ResultType getResult() {
			return result;
		}
		public void setResultDetails(final ResultDetails resultDetails) {
			this.resultDetails = resultDetails;
		}
		public ResultDetails getResultDetails() {
			return resultDetails;
		}

		public void setResultLocation(final Geolocation resultLocation) {
			this.resultLocation = resultLocation;
		}

		public Geolocation getResultLocation() {
			return resultLocation;
		}

		public void setClientGeolocation(final Geolocation clientGeolocation) {
			this.clientGeolocation = clientGeolocation;
		}

		public Geolocation getClientGeolocation() {
			return clientGeolocation;
		}

		public void setClientGeolocationQueried(final boolean isClientGeolocationQueried) {
			this.isClientGeolocationQueried = isClientGeolocationQueried;
		}

		public boolean isClientGeolocationQueried() {
			return isClientGeolocationQueried;
		}

		public void setRegionalGeoResult(final RegionalGeoResult regionalGeoResult) {
			this.regionalGeoResult = regionalGeoResult;
		}
		public RegionalGeoResult getRegionalGeoResult() {
			return regionalGeoResult;
		}

		public final void start() {
			time = System.currentTimeMillis();
		}
		public final void end() {
			time = System.currentTimeMillis() - time;
		}
	}

	public static Track getTrack() {
		return new Track();
	}

	final private Map<String, Tallies> dnsMap = new HashMap<String, Tallies>();
	final private Map<String, Tallies> httpMap = new HashMap<String, Tallies>();
	public Map<String, Tallies> getDnsMap() {
		return dnsMap;
	}
	public Map<String, Tallies> getHttpMap() {
		return httpMap;
	}
	public int getTotalDnsCount() {
		return totalDnsCount;
	}
	public long getAverageDnsTime() {
		if(totalDnsCount==0) { return 0; }
		return totalDnsTime/totalDnsCount;
	}
	public int getTotalHttpCount() {
		return totalHttpCount;
	}
	public long getAverageHttpTime() {
		if(totalHttpCount==0) { return 0; }
		return totalHttpTime/totalHttpCount;
	}
	public int getTotalDsMissCount() {
		return totalDsMissCount;
	}
	public void setTotalDsMissCount(final int totalDsMissCount) {
		this.totalDsMissCount = totalDsMissCount;
	}

	private int totalDnsCount;
	private long totalDnsTime;
	private int totalHttpCount;
	private long totalHttpTime;
	private int totalDsMissCount = 0;
	public Map<String,Long> getUpdateTracker() {
		return TrafficRouterManager.getTimeTracker();
	}
	public long getAppStartTime() {
		return appStartTime;
	}

	private long appStartTime;

	public void saveTrack(final Track t) {
		if (t.result == ResultType.DS_MISS) {
			// don't tabulate this, it's for a DS that doesn't exist
			totalDsMissCount++;
			return;
		}

		t.end();

		synchronized(this) {
			Map<String,Tallies> map;
			if(t.routeType == RouteType.DNS) {
				totalDnsCount++;
				totalDnsTime += t.time;
				map = dnsMap;
			} else {
				totalHttpCount++;
				totalHttpTime += t.time;
				map = httpMap;
			}
			Tallies tallies = map.get(t.fqdn);
			if(tallies == null) {
				tallies = new Tallies();
				map.put((t.fqdn==null)?"null":t.fqdn, tallies);
			}
			incTally(t, tallies);
		}
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private static void incTally(final Track t, final Tallies tallies) {
		switch(t.result) {
		case ERROR:
			tallies.errCount++;
			break;
		case CZ:
			tallies.czCount++;
			break;
		case GEO:
			tallies.geoCount++;
			break;
		case MISS:
			tallies.missCount++;
			break;
		case DS_REDIRECT:
			tallies.dsrCount++;
			break;
		case STATIC_ROUTE:
			tallies.staticRouteCount++;
			break;
		case FED:
			tallies.fedCount++;
			break;
		case RGDENY:
			tallies.regionalDeniedCount++;
			break;
		case RGALT:
			tallies.regionalAlternateCount++;
			break;
		default:
			break;
		}
	}

	public void init() {
		appStartTime = System.currentTimeMillis();
	}

	public void initialize(final Map<String, List<String>> initMap, final CacheRegister cacheRegister) {
		for (final String dsId : initMap.keySet()) {
			final List<String> dsNames = initMap.get(dsId);
			final DeliveryService ds = cacheRegister.getDeliveryService(dsId);

			if (ds != null) {
				for (int i = 0; i < dsNames.size(); i++) {
					final Track t = getTrack();
					final StringBuffer dsName = new StringBuffer(dsNames.get(i));
					RouteType rt;

					if (ds.isDns()) {
						rt = RouteType.DNS;

						if (i == 0) {
							dsName.insert(0, getDnsRoutingName() + ".");
						} else {
							continue;
						}
					} else {
						rt = RouteType.HTTP;
						dsName.insert(0, getHttpRoutingName() + ".");
					}

					t.setRouteType(rt, dsName.toString());
					t.setResult(ResultType.INIT);
					t.end();

					saveTrack(t);
				}
			}
		}
	}
	private String getDnsRoutingName() {
		return dnsRoutingName;
	}
	public void setDnsRoutingName(final String dnsRoutingName) {
		this.dnsRoutingName = dnsRoutingName;
	}
	private String getHttpRoutingName() {
		return httpRoutingName;
	}
	public void setHttpRoutingName(final String httpRoutingName) {
		this.httpRoutingName = httpRoutingName;
	}
}
