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

import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.Test;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoResult.RegionalGeoResultType;

public class RegionalGeoTest {
    private static final Logger LOGGER = Logger.getLogger(RegionalGeoTest.class);

    @Before
    public void setUp() throws Exception {
        final File dbFile = new File(getClass().getClassLoader().getResource("regional_geoblock.json").toURI());
        RegionalGeo.parseConfigFile(dbFile);
    }

    @Test
    public void testEnforceAllowed() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/TSN1";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALLOWED, result.getType());
        assertEquals(url, result.getUrl());
    }

    @Test
    public void testEnforceAlternateWithCache() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/TSN2";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITH_CACHE, result.getType());
        assertEquals("/path/redirect_T2", result.getUrl());
    }

    @Test
    public void testEnforceAlternateWithoutCache() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/TSN1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE, result.getType());
        assertEquals("http://example.com/redirect_T1", result.getUrl());
    }

    @Test
    public void testEnforceDeniedNoDsvc() {
        final String dsvcId = "ds-geoblock-no-exist";
        final String url = "http://ds1.example.com/TSN1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.DENIED, result.getType());
    }

    @Test
    public void testEnforceDeniedNoRegexMatch() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/TSN-NOT-EXIST";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.DENIED, result.getType());
    }

    @Test
    public void testEnforceAlternateToPathNoSlash() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/TSN3";
        final String postal = "V5D";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITH_CACHE, result.getType());
        assertEquals("/redirect_T3", result.getUrl());
    }

    @Test
    public void testEnforceAlternateNullPostal() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/TSN1";
        final String postal = null;
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE, result.getType());
        assertEquals("http://example.com/redirect_T1", result.getUrl());
    }

    @Test
    public void testEnforceAlternateEmptyPostalInclude() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/TSN2";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITH_CACHE, result.getType());
        assertEquals("/path/redirect_T2", result.getUrl());
    }

    @Test
    public void testEnforceAlternateEmptyPostalExclude() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/TSN1";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE, result.getType());
        assertEquals("http://example.com/redirect_T1", result.getUrl());
    }

    @Test
    public void testEnforceAllowLongPostal() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/TSN2";
        final String postal = "V5G 123";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALLOWED, result.getType());
        assertEquals(url, result.getUrl());
    }

    @Test
    public void testEnforceAlternateLongPotal() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/TSN2";
        final String postal = "N7G 123";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITH_CACHE, result.getType());
        assertEquals("/path/redirect_T2", result.getUrl());
    }

    @Test
    public void testEnforceWhiteListAllowed() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/TSN4";
        final String postal = null;
        final String ip = "129.100.254.2";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALLOWED, result.getType());
        assertEquals(url, result.getUrl());
    }

    @Test
    public void testEnforceNotInWhiteListAlternate() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/TSN4";
        final String postal = "N7G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALTERNATE_WITH_CACHE, result.getType());
        assertEquals("/redirect_T4", result.getUrl());
    }

    @Test
    public void testEnforceNotInWhiteListAllowedByPostal() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/TSN4";
        final String postal = "N6G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = new RegionalGeoResult();

        RegionalGeo.enforce(dsvcId, url, ip, postal, result);

        assertEquals(RegionalGeoResultType.ALLOWED, result.getType());
        assertEquals(url, result.getUrl());
    }
}

