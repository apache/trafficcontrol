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

package org.apache.traffic_control.traffic_router.core.loc;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.apache.traffic_control.traffic_router.core.request.Request;
import org.apache.traffic_control.traffic_router.core.router.HTTPRouteResult;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.File;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.regex.Pattern;

import static org.apache.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType.ALLOWED;
import static org.apache.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE;
import static org.apache.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType.ALTERNATE_WITH_CACHE;
import static org.apache.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType.DENIED;


public final class RegionalGeo {
    private static final Logger LOGGER = LogManager.getLogger(RegionalGeo.class);
    public static final String HTTP_SCHEME = "http://";
    public static final String HTTPS_SCHEME = "https://";
    private boolean fallback = false;
    private final Map<String, RegionalGeoDsvc> regionalGeoDsvcs = new HashMap<String, RegionalGeoDsvc>();
    private static RegionalGeo currentConfig = new RegionalGeo();

    private RegionalGeo() {

    }

    public void setFallback(final boolean fallback) {
        this.fallback = fallback;
    }

    public boolean isFallback() {
        return fallback;
    }

    private RegionalGeoRule matchRule(final String dsvcId, final String url) {
        final RegionalGeoDsvc regionalGeoDsvc = regionalGeoDsvcs.get(dsvcId);
        if (regionalGeoDsvc == null) {
            LOGGER.debug("RegionalGeo: dsvc not found: " + dsvcId);
            return null;
        }

        final RegionalGeoRule rule = regionalGeoDsvc.matchRule(url);
        if (rule == null) {
            LOGGER.debug("RegionalGeo: no rule match for dsvc "
                         + dsvcId + " with url " + url);
            return null;
        }

        return rule;
    }

    private boolean addRule(final String dsvcId, final String urlRegex,
                            final RegionalGeoRule.PostalsType postalsType, final Set<String> postals,
                            final NetworkNode networkRoot, final String alternateUrl, final boolean isSteeringDS, final List<RegionalGeoCoordinateRange> coordinateRanges) {

        // Loop check for alternateUrl with fqdn against the regex before adding
        Pattern urlRegexPattern;

        try {
            LOGGER.info("RegionalGeo: compile regex for url " + urlRegex);
            urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);
        } catch (Exception e) {
            LOGGER.error("RegionalGeo ERR: Pattern.compile exception", e);
            return false;
        }

        if ((alternateUrl.toLowerCase().startsWith(HTTP_SCHEME) || alternateUrl.toLowerCase().startsWith(HTTPS_SCHEME))
            && urlRegexPattern.matcher(alternateUrl).matches()) {
            LOGGER.error("RegionalGeo ERR: possible LOOP detected, alternate fqdn url " + alternateUrl
                         + " matches regex " + urlRegex + " in dsvc " +  dsvcId);
            return false;
        }

        if (isSteeringDS && !(alternateUrl.toLowerCase().startsWith(HTTP_SCHEME) || alternateUrl.toLowerCase().startsWith(HTTPS_SCHEME))) {
            LOGGER.error("RegionalGeo ERR: Alternate URL for Steering delivery service: "
                    +  dsvcId + " must start with " + HTTP_SCHEME + " or " + HTTPS_SCHEME);
            return false;
        }

        RegionalGeoDsvc regionalGeoDsvc = regionalGeoDsvcs.get(dsvcId);
        if (regionalGeoDsvc == null) {
            regionalGeoDsvc = new RegionalGeoDsvc(dsvcId);
            regionalGeoDsvcs.put(dsvcId, regionalGeoDsvc);
        }

        final RegionalGeoRule urlRule = new RegionalGeoRule(regionalGeoDsvc,
                urlRegex, urlRegexPattern,
                postalsType, postals,
                networkRoot, alternateUrl, coordinateRanges);

        LOGGER.info("RegionalGeo: adding " + urlRule);
        regionalGeoDsvc.addRule(urlRule);
        return true;
    }

    /// static methods
    private static NetworkNode parseWhiteListJson(final JsonNode json)
        throws NetworkNodeException {

        final NetworkNode.SuperNode root = new NetworkNode.SuperNode();

        for (final JsonNode subnetNode : json) {
            final String subnet = subnetNode.asText();
            final NetworkNode node = new NetworkNode(subnet, RegionalGeoRule.WHITE_LIST_NODE_LOCATION);

            if (subnet.indexOf(':') == -1) { // ipv4 or ipv6
                root.add(node);
            } else {
                root.add6(node);
            }
        }

        return root;
    }

    private static boolean checkCoordinateRangeValidity (final RegionalGeoCoordinateRange cr) {
        if ((cr.getMinLat() < -90.0 || cr.getMinLat() > 90.0) ||
                (cr.getMaxLat() < -90.0 || cr.getMaxLat() > 90.0) ||
                (cr.getMinLon() < -180.0 || cr.getMinLon() > 180.0) ||
                (cr.getMaxLon() < -180.0 || cr.getMaxLon() > 180.0)) {
            LOGGER.error("The supplied coordinate range is invalid. Latitude must be between -90.0 and +90.0, Longitude must be between -180.0 and +180.0.");
            return false;
        }
        return true;
    }

    private static List<RegionalGeoCoordinateRange> parseLocationJsonCoordinateRange(final JsonNode locationJson) {
        final List<RegionalGeoCoordinateRange> coordinateRange =  new ArrayList<>();
        final JsonNode coordinateRangeJson = locationJson.get("coordinateRange");
        if (coordinateRangeJson == null) {
            return null;
        }
        final ObjectMapper mapper = new ObjectMapper();
        RegionalGeoCoordinateRange cr = new RegionalGeoCoordinateRange();
        for (final JsonNode cRange : coordinateRangeJson) {
            cr  = mapper.convertValue(cRange, RegionalGeoCoordinateRange.class);
            if (checkCoordinateRangeValidity(cr)) {
                coordinateRange.add(cr);
            }
        }
        return coordinateRange;
    }

    private static RegionalGeoRule.PostalsType parseLocationJson(final JsonNode locationJson,
        final Set<String> postals) {

        RegionalGeoRule.PostalsType postalsType = RegionalGeoRule.PostalsType.UNDEFINED;
        JsonNode postalsJson = locationJson.get("includePostalCode");
        
        if (postalsJson != null) {
            postalsType = RegionalGeoRule.PostalsType.INCLUDE;
        } else {
            postalsJson = locationJson.get("excludePostalCode");
            if (postalsJson == null) {
                LOGGER.error("RegionalGeo ERR: no include/exclude in geolocation");
                return RegionalGeoRule.PostalsType.UNDEFINED;
            }
        
            postalsType = RegionalGeoRule.PostalsType.EXCLUDE;
        }

        for (final JsonNode postal : postalsJson) {
            postals.add(postal.asText());
        }

        return postalsType;

    }

    @SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
    private static RegionalGeo parseConfigJson(final JsonNode json) {

        final RegionalGeo regionalGeo = new RegionalGeo();
        regionalGeo.setFallback(true);
        try {
            final JsonNode dsvcsJson = JsonUtils.getJsonNode(json, "deliveryServices");
            LOGGER.info("RegionalGeo: parse json with rule count " + dsvcsJson.size());

            for (final JsonNode ruleJson : dsvcsJson) {

                final String dsvcId = JsonUtils.getString(ruleJson, "deliveryServiceId");
                if (dsvcId.trim().isEmpty()) {
                    LOGGER.error("RegionalGeo ERR: deliveryServiceId empty");
                    return null;
                }
                Boolean isSteeringDS = false;
                try {
                    isSteeringDS = JsonUtils.getBoolean(ruleJson, "isSteeringDS");
                } catch (JsonUtilsException e) {
                    //It's not in the config so we can just keep it set as false.
                    LOGGER.debug("RegionalGeo ERR: isSteeringDS empty");
                }

                final String urlRegex = JsonUtils.getString(ruleJson, "urlRegex");
                if (urlRegex.trim().isEmpty()) {
                    LOGGER.error("RegionalGeo ERR: urlRegex empty");
                    return null;
                }

                final String redirectUrl = JsonUtils.getString(ruleJson, "redirectUrl");
                if (redirectUrl.trim().isEmpty()) {
                    LOGGER.error("RegionalGeo ERR: redirectUrl empty");
                    return null;
                }

                // FSAs (postal codes)
                final JsonNode locationJson = JsonUtils.getJsonNode(ruleJson, "geoLocation");
                final Set<String> postals = new HashSet<String>();
                final RegionalGeoRule.PostalsType postalsType = parseLocationJson(locationJson, postals);
                if (postalsType == RegionalGeoRule.PostalsType.UNDEFINED) {
                    LOGGER.error("RegionalGeo ERR: geoLocation empty");
                    return null;
                }
                // coordinate range
                final List<RegionalGeoCoordinateRange> coordinateRanges = parseLocationJsonCoordinateRange(locationJson);

                // white list
                NetworkNode whiteListRoot = null;
                final JsonNode whiteListJson = ruleJson.get("ipWhiteList");
                if (whiteListJson != null) {
                    whiteListRoot = parseWhiteListJson(whiteListJson);
                }


                // add the rule
                if (!regionalGeo.addRule(dsvcId, urlRegex, postalsType, postals, whiteListRoot, redirectUrl, isSteeringDS, coordinateRanges)) {
                    LOGGER.error("RegionalGeo ERR: add rule failed on parsing json file");
                    return null;
                }
            }

            regionalGeo.setFallback(false);
            return regionalGeo;
        } catch (Exception e) {
            LOGGER.error("RegionalGeo ERR: parse json file with exception", e);
        }

        return null;
    }

    public static boolean parseConfigFile(final File f, final boolean verifyOnly) {
        final ObjectMapper mapper = new ObjectMapper();
        JsonNode json = null;
        try {
            json = mapper.readTree(f);
        } catch (Exception e) {
            LOGGER.error("RegionalGeo ERR: json file exception " + f, e);
            currentConfig.setFallback(true);
            return false;
        }

        final RegionalGeo regionalGeo = parseConfigJson(json);
        if (regionalGeo == null) {
            currentConfig.setFallback(true);
            return false;
        }
        
        if (!verifyOnly) {
            currentConfig = regionalGeo; // point to the new parsed object
        }

        currentConfig.setFallback(false);
        LOGGER.debug("RegionalGeo: create instance from new json");
        return true;
    }


    public static RegionalGeoResult enforce(final String dsvcId, final String url,
                                            final String ip, final String postalCode, final double lat, final double lon) {

        final RegionalGeoResult result = new RegionalGeoResult();
        boolean allowed = false;
        RegionalGeoRule rule = null;

        result.setPostal(postalCode);
        result.setUsingFallbackConfig(currentConfig.isFallback());
        result.setAllowedByWhiteList(false);

        rule = currentConfig.matchRule(dsvcId, url);
        if (rule == null) {
            result.setHttpResponseCode(RegionalGeoResult.REGIONAL_GEO_DENIED_HTTP_CODE);
            result.setType(DENIED);
            LOGGER.debug("RegionalGeo: denied for dsvc " + dsvcId
                         + ", url " + url + ", postal " + postalCode);
            return result;
        }

        // first match whitelist, then FSA (postal)
        if (rule.isIpInWhiteList(ip)) {
            LOGGER.debug("RegionalGeo: allowing ip in whitelist");
            allowed = true;
            result.setAllowedByWhiteList(true);
        } else {
            if (postalCode == null || postalCode.isEmpty()) {
                LOGGER.warn("RegionalGeo: alternate a request with null or empty postal");
                allowed = rule.isAllowedCoordinates(lat, lon);
            } else {
                allowed = rule.isAllowedPostal(postalCode);
            }
        }

        final String alternateUrl = rule.getAlternateUrl();
        result.setRuleType(rule.getPostalsType());

        if (allowed) {
            result.setUrl(url);
            result.setType(ALLOWED);
        } else {
            // For a disallowed client, if alternateUrl starts with "http://" or "https://"
            // just redirect the client to this url without any cache selection;
            // if alternateUrl only has path and file name like "/path/abc.html",
            // then cache selection process will be needed, and hostname will be
            // added to make it like "http://cache01.example.com/path/abc.html" later.
            if (alternateUrl.toLowerCase().startsWith(HTTP_SCHEME) || alternateUrl.toLowerCase().startsWith(HTTPS_SCHEME)) {
                result.setUrl(alternateUrl);
                result.setType(ALTERNATE_WITHOUT_CACHE);
            } else {
                String redirectUrl;
                if (alternateUrl.startsWith("/")) { // add a '/' prefix if necessary for url path
                    redirectUrl = alternateUrl;
                } else {
                    redirectUrl = "/" + alternateUrl;
                }

                LOGGER.debug("RegionalGeo: alternate with cache url " + redirectUrl);
                result.setUrl(redirectUrl);
                result.setType(ALTERNATE_WITH_CACHE);
            }
        }

        LOGGER.debug("RegionalGeo: result " + result + " for dsvc " + dsvcId + ", url " + url + ", ip " + ip);

        return result;
    }

    public static void enforce(final TrafficRouter trafficRouter, final Request request, final DeliveryService deliveryService, final Cache cache,
                               final HTTPRouteResult routeResult, final Track track) throws MalformedURLException {
        enforce(trafficRouter, request, deliveryService, cache, routeResult, track, false);
    }

    @SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
    public static void enforce(final TrafficRouter trafficRouter, final Request request,
        final DeliveryService deliveryService, final Cache cache,
        final HTTPRouteResult routeResult, final Track track, final boolean isSteering) throws MalformedURLException {

        LOGGER.debug("RegionalGeo: enforcing");

        Geolocation clientGeolocation = null;
        try {
            clientGeolocation = trafficRouter.getClientGeolocation(request.getClientIP(), track, deliveryService);
        } catch (GeolocationException e) {
            LOGGER.warn("RegionalGeo: failed looking up Client GeoLocation: " + e.getMessage());
        }

        String postalCode = null;
        double lat = 0.0;
        double lon = 0.0;

        if (clientGeolocation != null) {
            postalCode = clientGeolocation.getPostalCode();

            // Get the first 3 chars in the postal code. These 3 chars are called FSA in Canadian postal codes.
            if (postalCode != null && postalCode.length() > 3) {
                postalCode = postalCode.substring(0, 3);
            } else {
                lat = clientGeolocation.getLatitude();
                lon = clientGeolocation.getLongitude();
            }
        }

        final HTTPRequest httpRequest = HTTPRequest.class.cast(request);
        final RegionalGeoResult result = enforce(deliveryService.getId(), httpRequest.getRequestedUrl(), 
                                                 httpRequest.getClientIP(), postalCode, lat, lon);

        if (cache == null && result.getType() == ALTERNATE_WITH_CACHE) {
            LOGGER.debug("RegionalGeo: denied for dsvc " + deliveryService.getId() + ", url " + httpRequest.getRequestedUrl() + ", postal " + postalCode + ". Relative re-direct URLs not allowed for Multi Route Delivery Services.");
            result.setHttpResponseCode(RegionalGeoResult.REGIONAL_GEO_DENIED_HTTP_CODE);
            result.setType(DENIED);
        }

        if (cache == null && result.getType() == ALLOWED) {
            LOGGER.debug("RegionalGeo: Client is allowed to access steering service, returning null re-direct URL");
            result.setUrl(null);
            updateTrack(track, result);
            return;
        }

        updateTrack(track, result);

        if (result.getType() == DENIED) {
            routeResult.setResponseCode(result.getHttpResponseCode());
        } else {
            final String redirectURIString = createRedirectURIString(httpRequest, deliveryService, cache, result);
            if(!"Denied".equals(redirectURIString)){
                routeResult.addUrl(new URL(redirectURIString));
            }else{
                LOGGER.warn("RegionalGeo: this needs a better error message, createRedirectURIString returned denied");
            }
        }
    }


    private static void updateTrack(final Track track, final RegionalGeoResult regionalGeoResult) {
        track.setRegionalGeoResult(regionalGeoResult);

        final RegionalGeoResultType resultType = regionalGeoResult.getType();

        if (resultType == DENIED) {
            track.setResult(ResultType.RGDENY);
            track.setResultDetails(ResultDetails.REGIONAL_GEO_NO_RULE);
            return;
        }

        if (resultType == ALTERNATE_WITH_CACHE) {
            track.setResult(ResultType.RGALT);
            track.setResultDetails(ResultDetails.REGIONAL_GEO_ALTERNATE_WITH_CACHE);
            return;
        }

        if (resultType == ALTERNATE_WITHOUT_CACHE) {
            track.setResult(ResultType.RGALT);
            track.setResultDetails(ResultDetails.REGIONAL_GEO_ALTERNATE_WITHOUT_CACHE);
            return;
        }

        // else ALLOWED, result & resultDetail shall be normal case, do not modify
    }   

    private static String createRedirectURIString(final HTTPRequest request, final DeliveryService deliveryService, 
        final Cache cache, final RegionalGeoResult regionalGeoResult) {

        if (regionalGeoResult.getType() == ALLOWED) {
            return deliveryService.createURIString(request, cache);
        }

        if (regionalGeoResult.getType() == ALTERNATE_WITH_CACHE) {
            return deliveryService.createURIString(request, regionalGeoResult.getUrl(), cache);
        }

        if (regionalGeoResult.getType() == ALTERNATE_WITHOUT_CACHE) {
            return regionalGeoResult.getUrl();
        }

        return "Denied"; // DENIED
    }

}

