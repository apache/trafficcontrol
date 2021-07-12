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

import org.junit.Test;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;
import java.util.regex.Pattern;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class RegionalGeoRuleTest {

    @Test
    public void testIsAllowedCoordinateRanges() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.INCLUDE;
        final Set<String> postals = new HashSet<String>();

        final NetworkNode whiteList = new NetworkNode.SuperNode();
        final String alternateUrl = "/alternate.m3u8";
        final ArrayList<RegionalGeoCoordinateRange> coordinateRanges = new ArrayList<>();
        RegionalGeoCoordinateRange coordinateRange = new RegionalGeoCoordinateRange();
        RegionalGeoCoordinateRange coordinateRange2 = new RegionalGeoCoordinateRange();
        coordinateRange.setMinLat(10.0);
        coordinateRange.setMinLon(165.0);
        coordinateRange.setMaxLat(22.0);
        coordinateRange.setMaxLon(179.0);
        coordinateRanges.add(coordinateRange);
        coordinateRange2.setMinLat(17.0);
        coordinateRange2.setMinLon(-20.0);
        coordinateRange2.setMaxLat(25.0);
        coordinateRange2.setMaxLon(19.0);
        coordinateRanges.add(coordinateRange2);

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, coordinateRanges);

        boolean allowed;

        allowed = urlRule.isAllowedCoordinates(11.0, 170.0);
        assertThat(allowed, equalTo(true));

        allowed = urlRule.isAllowedCoordinates(13.0, 162.0);
        assertThat(allowed, equalTo(false));

        allowed = urlRule.isAllowedCoordinates(23.0, 22.0);
        assertThat(allowed, equalTo(false));

        allowed = urlRule.isAllowedCoordinates(23.0, -12.0);
        assertThat(allowed, equalTo(true));

        allowed = urlRule.isAllowedCoordinates(9.0, 21.0);
        assertThat(allowed, equalTo(false));
    }

    @Test
    public void testMatchesUrl() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.INCLUDE;
        final Set<String> postals = new HashSet<String>();
        final NetworkNode whiteList = new NetworkNode.SuperNode();
        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, null);

        boolean matches;
        String url = "http://example.com/abc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(true));

        url = "http://example.com/AbC.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(true));

        url = "http://example.com/path/ABC.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(true));

        url = "http://example.com/cbaabc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(true));

        url = "http://example.com/cba.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(false));

        url = "http://example.com/abcabc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertThat(matches, equalTo(true));
    }

    @Test
    public void testIsAllowedPostalInclude() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.INCLUDE;
        final Set<String> postals = new HashSet<String>();
        postals.add("N6G");
        postals.add("N7G");
        final NetworkNode whiteList = new NetworkNode.SuperNode();
        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, null);

        boolean allowed;

        allowed = urlRule.isAllowedPostal("N6G");
        assertThat(allowed, equalTo(true));

        allowed = urlRule.isAllowedPostal("N7G");
        assertThat(allowed, equalTo(true));

        allowed = urlRule.isAllowedPostal("N8G");
        assertThat(allowed, equalTo(false));
    }

    @Test
    public void testIsAllowedPostalExclude() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.EXCLUDE;
        final Set<String> postals = new HashSet<String>();
        postals.add("N6G");
        postals.add("N7G");
        final NetworkNode whiteList = new NetworkNode.SuperNode();
        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, null);

        boolean allowed;

        allowed = urlRule.isAllowedPostal("N6G");
        assertThat(allowed, equalTo(false));

        allowed = urlRule.isAllowedPostal("N7G");
        assertThat(allowed, equalTo(false));

        allowed = urlRule.isAllowedPostal("N8G");
        assertThat(allowed, equalTo(true));
    }

    @Test
    public void testIsInWhiteList() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.EXCLUDE;
        final Set<String> postals = new HashSet<String>();
        final NetworkNode.SuperNode whiteList = new NetworkNode.SuperNode();
        final String location = RegionalGeoRule.WHITE_LIST_NODE_LOCATION;
        whiteList.add(new NetworkNode("10.74.50.0/24", location));
        whiteList.add(new NetworkNode("10.74.0.0/16", location));
        whiteList.add(new NetworkNode("192.168.250.1/32", location));
        whiteList.add(new NetworkNode("128.128.50.3/32", location));
        whiteList.add(new NetworkNode("128.128.50.3/22", location));
        whiteList.add6(new NetworkNode("2001:0db8:0:f101::1/64", location));
        whiteList.add6(new NetworkNode("2001:0db8:0:f101::1/48", location));

        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, null);

        boolean in;

        in = urlRule.isIpInWhiteList("10.74.50.12");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("10.75.51.12");
        assertThat(in, equalTo(false));

        in = urlRule.isIpInWhiteList("10.74.51.1");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("10.74.50.255");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("192.168.250.1");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("128.128.50.3");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("128.128.50.7");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("128.128.2.1");
        assertThat(in, equalTo(false));

        in = urlRule.isIpInWhiteList("2001:0db8:0:f101::2");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("2001:0db8:0:f102::1");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
        assertThat(in, equalTo(false));
    }

    @Test
    public void testIsInWhiteListInvalidParam() throws Exception {
        try {
            final String urlRegex = ".*abc.m3u8";
            final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.EXCLUDE;
            final Set<String> postals = new HashSet<String>();
            final NetworkNode.SuperNode whiteList = new NetworkNode.SuperNode();
            final String location = RegionalGeoRule.WHITE_LIST_NODE_LOCATION;
            whiteList.add(new NetworkNode("10.256.0.0/10", location));
            //whiteList.add(new NetworkNode("a.b.d.0/10", location));

            final String alternateUrl = "/alternate.m3u8";

            Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

            final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                    urlRegex, urlRegexPattern,
                    ruleType, postals,
                    whiteList, alternateUrl, null);

            boolean in;

            in = urlRule.isIpInWhiteList("10.74.50.12");
            assertThat(in, equalTo(false));

            in = urlRule.isIpInWhiteList("10.74.51.12");
            assertThat(in, equalTo(false));

            in = urlRule.isIpInWhiteList("1.1.50.1");
            assertThat(in, equalTo(false));

            in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
            assertThat(in, equalTo(false));
        } catch (NetworkNodeException e) {

        }
    }

    @Test
    public void testIsInWhiteListGlobalMatch() throws Exception {
        final String urlRegex = ".*abc.m3u8";
        final RegionalGeoRule.PostalsType ruleType = RegionalGeoRule.PostalsType.EXCLUDE;
        final Set<String> postals = new HashSet<String>();
        final NetworkNode.SuperNode whiteList = new NetworkNode.SuperNode();
        final String location = RegionalGeoRule.WHITE_LIST_NODE_LOCATION;
        whiteList.add(new NetworkNode("0.0.0.0/0", location));

        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl, null);

        boolean in;

        in = urlRule.isIpInWhiteList("10.74.50.12");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("10.74.51.12");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("1.1.50.1");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("222.254.254.254");
        assertThat(in, equalTo(true));

        in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
        assertThat(in, equalTo(false));
    }
}

