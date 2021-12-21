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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.List;
import java.util.Objects;
import java.util.Set;
import java.util.regex.Pattern;


public class RegionalGeoRule {
    private static final Logger LOGGER = LogManager.getLogger(RegionalGeoRule.class);

    public static final String WHITE_LIST_NODE_LOCATION = "w";

    public enum PostalsType {
        EXCLUDE, INCLUDE, UNDEFINED
    }

    private final RegionalGeoDsvc regionalGeoDsvc;
    private final String urlRegex;
    private final Pattern pattern;

    private final PostalsType postalsType;
    private final Set<String> postals;

    private final NetworkNode whiteListRoot;

    private final String alternateUrl; // if disallowed, client will be redirected to this url

    private final List<RegionalGeoCoordinateRange> coordinateRanges;

    public RegionalGeoRule(final RegionalGeoDsvc regionalGeoDsvc,
                           final String urlRegex, final Pattern urlRegexPattern, final PostalsType postalsType,
                           final Set<String> postals, final NetworkNode whiteListRoot,
                           final String alternateUrl, final List<RegionalGeoCoordinateRange> coordinateRanges) {
        this.regionalGeoDsvc = regionalGeoDsvc;
        this.urlRegex = urlRegex;
        this.pattern = urlRegexPattern;
        this.postalsType = postalsType;
        this.postals = postals;
        this.whiteListRoot = whiteListRoot;
        this.alternateUrl = alternateUrl;
        this.coordinateRanges = coordinateRanges;
    }

    public boolean matchesUrl(final String url) {
        return pattern.matcher(url).matches();
    }

    public boolean isAllowedPostal(final String postal) {

        if (postalsType == PostalsType.INCLUDE) {
            if (postals.contains(postal)) {
                return true;
            }
        } else { // EXCLUDE
            if (!postals.contains(postal)) {
                return true;
            }
        }

        return false;
    }

    public boolean isAllowedCoordinates(final double lat, final double lon) {
        if (coordinateRanges == null) {
            return false;
        }
        for (int i=0; i < coordinateRanges.size(); i++) {
            final RegionalGeoCoordinateRange coordinateRange = coordinateRanges.get(i);
            if ((lat >= coordinateRange.getMinLat() && lon >= coordinateRange.getMinLon()) &&
                    (lat <= coordinateRange.getMaxLat() && lon <= coordinateRange.getMaxLon())) {
                return true;
            }
        }
        return false;
    }

    public boolean isIpInWhiteList(final String ip) {
        if (whiteListRoot == null) {
            return false;
        }

        try {
            final NetworkNode nn = whiteListRoot.getNetwork(ip);
            if (Objects.equals(nn.getLoc(), WHITE_LIST_NODE_LOCATION)) {
                return true;
            }
        } catch (NetworkNodeException e) {
            LOGGER.warn("RegionalGeo: exception", e);
        }

        return false;
    }

    public String getUrlRegex() {
        return urlRegex;
    }

    public Pattern getPattern() {
        return pattern;
    }

    public PostalsType getPostalsType() {
        return postalsType;
    }

    public String getAlternateUrl() {
        return alternateUrl;
    }

    public String toString() {
        final StringBuilder sb = new StringBuilder();

        sb.append("RULE: dsvc ");
        sb.append(regionalGeoDsvc.getId());
        sb.append(", regex ");
        sb.append(urlRegex);
        sb.append(", alternate ");
        sb.append(alternateUrl);
        sb.append(", type ");
        sb.append(postalsType);

        sb.append(", postals ");
        for (final String s : postals) {
            sb.append(s);
            sb.append(',');
        }

        return sb.toString();
    }

}

