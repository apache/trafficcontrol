/*
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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.io.FileReader;
import java.net.InetAddress;
import java.net.MalformedURLException;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.HTTPRouteResult;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.google.common.net.InetAddresses;
import com.maxmind.geoip2.model.AnonymousIpResponse;

public final class AnonymousIp {

	private static final Logger LOGGER = Logger.getLogger(AnonymousIp.class);

	private static AnonymousIp currentConfig = new AnonymousIp();

	// Feature flipper
	// This is set to true if the CRConfig parameters containing the MMDB URL
	// and the config url are present AND any delivery service has the feature
	// enabled
	public boolean enabled = false;

	private boolean blockAnonymousIp = true;
	private boolean blockHostingProvider = true;
	private boolean blockPublicProxy = true;
	private boolean blockTorExitNode = true;

	private AnonymousIpWhitelist ipv4Whitelist;
	private AnonymousIpWhitelist ipv6Whitelist;

	public final static int BLOCK_CODE = 403;
	public final static String WHITE_LIST_LOC = "w";

	private AnonymousIp() {
		try {
			ipv4Whitelist = new AnonymousIpWhitelist();
			ipv6Whitelist = new AnonymousIpWhitelist();
		} catch (NetworkNodeException e) {
			LOGGER.error("AnonymousIp ERR: Network node exception ", e);
		}
	}

	/*
	 * Returns the current anonymous ip object
	 */
	public static AnonymousIp getCurrentConfig() {
		return currentConfig;
	}

	/*
	 * Returns the list of subnets in the IPv4 whitelist
	 */
	public AnonymousIpWhitelist getIPv4Whitelist() {
		return ipv4Whitelist;
	}

	/*
	 * Returns the list of subnets in the IPv6 whitelist
	 */
	public AnonymousIpWhitelist getIPv6Whitelist() {
		return ipv6Whitelist;
	}

	private static void parseIPv4Whitelist(final JSONObject config, final AnonymousIp anonymousIp) throws JSONException {
		if (config.optJSONArray("ip4Whitelist") != null) {
			try {
				anonymousIp.ipv4Whitelist = new AnonymousIpWhitelist();
				anonymousIp.ipv4Whitelist.init(config.optJSONArray("ip4Whitelist"));
			} catch (NetworkNodeException e) {
				LOGGER.error("Anonymous Ip ERR: Network node err ", e);
			}
		}
	}

	private static void parseIPv6Whitelist(final JSONObject config, final AnonymousIp anonymousIp) throws JSONException {
		if (config.optJSONArray("ip6Whitelist") != null) {
			try {
				anonymousIp.ipv6Whitelist = new AnonymousIpWhitelist();
				anonymousIp.ipv6Whitelist.init(config.optJSONArray("ip6Whitelist"));
			} catch (NetworkNodeException e) {
				LOGGER.error("Anonymous Ip ERR: Network node err ", e);
			}
		}
	}

	@SuppressWarnings({ "PMD.NPathComplexity", "PMD.CyclomaticComplexity" })
	private static AnonymousIp parseConfigJson(final JSONObject config) {
		final AnonymousIp anonymousIp = new AnonymousIp();
		try {
			final JSONObject blockingTypes = config.getJSONObject("anonymousIp");

			anonymousIp.blockAnonymousIp = blockingTypes.getBoolean("blockAnonymousVPN");
			anonymousIp.blockHostingProvider = blockingTypes.getBoolean("blockHostingProvider");
			anonymousIp.blockPublicProxy = blockingTypes.getBoolean("blockPublicProxy");
			anonymousIp.blockTorExitNode = blockingTypes.getBoolean("blockTorExitNode");

			anonymousIp.enabled = AnonymousIp.currentConfig.enabled;

			parseIPv4Whitelist(config, anonymousIp);
			parseIPv6Whitelist(config, anonymousIp);

			return anonymousIp;
		} catch (Exception e) {
			LOGGER.error("AnonymousIp ERR: parsing config file failed", e);
		}

		return null;
	}

	@SuppressWarnings({ "PMD.NPathComplexity" })
	public static boolean parseConfigFile(final File f, final boolean verifyOnly) {
		JSONObject json = null;
		try {
			json = new JSONObject(new JSONTokener(new FileReader(f)));
		} catch (Exception e) {
			LOGGER.error("AnonymousIp ERR: json file exception " + f, e);
			return false;
		}

		final AnonymousIp anonymousIp = parseConfigJson(json);

		if (anonymousIp == null) {
			return false;
		}

		if (!verifyOnly) {
			currentConfig = anonymousIp; // point to the new parsed object
		}

		return true;
	}

	private static boolean inWhitelist(final String address) {
		// If the address is ipv4 check against the ipv4whitelist
		if (address.indexOf(':') == -1) {
			if (currentConfig.ipv4Whitelist.contains(address)) {
				return true;
			}
		}

		// If the address is ipv6 check against the ipv6whitelist
		else {
			if (currentConfig.ipv6Whitelist.contains(address)) {
				return true;
			}
		}

		return false;
	}

	@SuppressWarnings({ "PMD.CyclomaticComplexity", "PMD.NPathComplexity" })
	public static boolean enforce(final TrafficRouter trafficRouter, final String dsvcId, final String url, final String ip) {

		final InetAddress address = InetAddresses.forString(ip);

		if (inWhitelist(ip)) {
			return false;
		}

		final AnonymousIpResponse response = trafficRouter.getAnonymousIpDatabaseService().lookupIp(address);

		if (response == null) {
			return false;
		}

		// Check if the ip should be blocked by checking if the ip falls into a
		// specific policy
		if (AnonymousIp.getCurrentConfig().blockAnonymousIp && response.isAnonymousVpn()) {
			return true;
		}

		if (AnonymousIp.getCurrentConfig().blockHostingProvider && response.isHostingProvider()) {
			return true;
		}

		if (AnonymousIp.getCurrentConfig().blockPublicProxy && response.isPublicProxy()) {
			return true;
		}

		if (AnonymousIp.getCurrentConfig().blockTorExitNode && response.isTorExitNode()) {
			return true;
		}

		return false;
	}

	@SuppressWarnings({ "PMD.CyclomaticComplexity" })
	/*
	 * Enforces the anonymous ip blocking policies
	 * 
	 * If the Delivery Service has anonymous ip blocking enabled And the ip is
	 * in the anonymous ip database The ip will be blocked if it matches a
	 * policy defined in the config file
	 */
	public static void enforce(final TrafficRouter trafficRouter, final Request request, final DeliveryService deliveryService, final Cache cache,
			final HTTPRouteResult routeResult, final Track track) throws MalformedURLException {

		final HTTPRequest httpRequest = HTTPRequest.class.cast(request);

		// If the database isn't initialized dont block
		if (!trafficRouter.getAnonymousIpDatabaseService().isInitialized()) {
			return;
		}

		// Check if the ip is allowed
		final boolean block = enforce(trafficRouter, deliveryService.getId(), httpRequest.getRequestedUrl(), httpRequest.getClientIP());

		// Block the ip if it is not allowed
		if (block) {
			routeResult.setResponseCode(AnonymousIp.BLOCK_CODE);
			track.setResult(ResultType.ANON_BLOCK);
		}
	}
}
