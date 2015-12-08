/*
 * Copyright 2015 Cisco Systems, Inc.
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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.io.File;
import java.io.FileReader;
import java.util.regex.Pattern;
import java.util.Set;
import java.util.HashSet;
import java.util.List;
import java.util.ArrayList;

import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.Test;

public class RegionalGeoRuleTest {
    private static final Logger LOGGER = Logger.getLogger(RegionalGeoRuleTest.class);

    @Before
    public void setUp() throws Exception {
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
                whiteList, alternateUrl);

        boolean matches;
        String url = "http://example.com/abc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(true, matches);

        url = "http://example.com/AbC.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(true, matches);

        url = "http://example.com/path/ABC.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(true, matches);

        url = "http://example.com/cbaabc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(true, matches);

        url = "http://example.com/cba.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(false, matches);

        url = "http://example.com/abcabc.m3u8";
        matches = urlRule.matchesUrl(url);
        assertEquals(true, matches);
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
                whiteList, alternateUrl);

        boolean allowed;

        allowed = urlRule.isAllowedPostal("N6G");
        assertEquals(true, allowed);

        allowed = urlRule.isAllowedPostal("N7G");
        assertEquals(true, allowed);

        allowed = urlRule.isAllowedPostal("N8G");
        assertEquals(false, allowed);
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
                whiteList, alternateUrl);

        boolean allowed;

        allowed = urlRule.isAllowedPostal("N6G");
        assertEquals(false, allowed);

        allowed = urlRule.isAllowedPostal("N7G");
        assertEquals(false, allowed);

        allowed = urlRule.isAllowedPostal("N8G");
        assertEquals(true, allowed);
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
        LOGGER.debug("RegionalGeo: whitelist " + whiteList);

        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl);

        boolean in;

        in = urlRule.isIpInWhiteList("10.74.50.12");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("10.75.51.12");
        assertEquals(false, in);

        in = urlRule.isIpInWhiteList("10.74.51.1");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("10.74.50.255");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("192.168.250.1");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("128.128.50.3");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("128.128.50.7");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("128.128.2.1");
        assertEquals(false, in);

        in = urlRule.isIpInWhiteList("2001:0db8:0:f101::2");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("2001:0db8:0:f102::1");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
        assertEquals(false, in);
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
            LOGGER.debug("RegionalGeo: whitelist " + whiteList);

            final String alternateUrl = "/alternate.m3u8";

            Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

            final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                    urlRegex, urlRegexPattern,
                    ruleType, postals,
                    whiteList, alternateUrl);

            boolean in;

            in = urlRule.isIpInWhiteList("10.74.50.12");
            assertEquals(false, in);

            in = urlRule.isIpInWhiteList("10.74.51.12");
            assertEquals(false, in);

            in = urlRule.isIpInWhiteList("1.1.50.1");
            assertEquals(false, in);

            in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
            assertEquals(false, in);
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
        LOGGER.debug("RegionalGeo: whitelist " + whiteList);

        final String alternateUrl = "/alternate.m3u8";

        Pattern urlRegexPattern = Pattern.compile(urlRegex, Pattern.CASE_INSENSITIVE);

        final RegionalGeoRule urlRule = new RegionalGeoRule(null,
                urlRegex, urlRegexPattern,
                ruleType, postals,
                whiteList, alternateUrl);

        boolean in;

        in = urlRule.isIpInWhiteList("10.74.50.12");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("10.74.51.12");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("1.1.50.1");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("222.254.254.254");
        assertEquals(true, in);

        in = urlRule.isIpInWhiteList("2001:0db8:1:f101::3");
        assertEquals(false, in);
    }
}

