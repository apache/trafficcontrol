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

package org.apache.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.net.InetAddress;
import java.net.MalformedURLException;
import java.net.URL;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.net.InetAddresses;
import com.maxmind.geoip2.model.AnonymousIpResponse;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.apache.traffic_control.traffic_router.core.request.Request;
import org.apache.traffic_control.traffic_router.core.router.HTTPRouteResult;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;

public final class AnonymousIp {

	private static final Logger LOGGER = LogManager.getLogger(AnonymousIp.class);

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

	private String redirectUrl;

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

	private static void parseIPv4Whitelist(final JsonNode config, final AnonymousIp anonymousIp) throws JsonUtilsException {
		if (config.has("ip4Whitelist")) {
			try {
				anonymousIp.ipv4Whitelist = new AnonymousIpWhitelist();
				anonymousIp.ipv4Whitelist.init(JsonUtils.getJsonNode(config, "ip4Whitelist"));
			} catch (NetworkNodeException e) {
				LOGGER.error("Anonymous Ip ERR: Network node err ", e);
			}
		}
	}

	private static void parseIPv6Whitelist(final JsonNode config, final AnonymousIp anonymousIp) throws JsonUtilsException {
		if (config.has("ip6Whitelist")) {
			try {
				anonymousIp.ipv6Whitelist = new AnonymousIpWhitelist();
				anonymousIp.ipv6Whitelist.init(JsonUtils.getJsonNode(config, "ip6Whitelist"));
			} catch (NetworkNodeException e) {
				LOGGER.error("Anonymous Ip ERR: Network node err ", e);
			}
		}
	}

	@SuppressWarnings({ "PMD.NPathComplexity", "PMD.CyclomaticComplexity" })
	private static AnonymousIp parseConfigJson(final JsonNode config) {
		final AnonymousIp anonymousIp = new AnonymousIp();
		try {
			final JsonNode blockingTypes = JsonUtils.getJsonNode(config, "anonymousIp");

			anonymousIp.blockAnonymousIp = JsonUtils.getBoolean(blockingTypes, "blockAnonymousVPN");
			anonymousIp.blockHostingProvider = JsonUtils.getBoolean(blockingTypes, "blockHostingProvider");
			anonymousIp.blockPublicProxy = JsonUtils.getBoolean(blockingTypes, "blockPublicProxy");
			anonymousIp.blockTorExitNode = JsonUtils.getBoolean(blockingTypes, "blockTorExitNode");

			anonymousIp.enabled = AnonymousIp.currentConfig.enabled;

			parseIPv4Whitelist(config, anonymousIp);
			parseIPv6Whitelist(config, anonymousIp);

			if (config.has("redirectUrl")) {
				anonymousIp.redirectUrl = JsonUtils.getString(config, "redirectUrl");
			}

			return anonymousIp;
		} catch (Exception e) {
			LOGGER.error("AnonymousIp ERR: parsing config file failed", e);
		}

		return null;
	}

	@SuppressWarnings({ "PMD.NPathComplexity" })
	public static boolean parseConfigFile(final File f, final boolean verifyOnly) {
		JsonNode json = null;
		try {
			final ObjectMapper mapper = new ObjectMapper();
			json = mapper.readTree(f);
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
			if (AnonymousIp.getCurrentConfig().redirectUrl != null) {
				routeResult.addUrl(new URL(AnonymousIp.getCurrentConfig().redirectUrl));
			}
		}
	}
}
