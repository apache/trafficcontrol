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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

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
    @Before
    public void setUp() throws Exception {
        final File dbFile = new File("src/test/resources/regional_geoblock.json");
        RegionalGeo.parseConfigFile(dbFile, false);
    }

    @Test
    public void testEnforceAllowed() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceAlternateWithCache() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/live2";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/path/redirect_T2"));
    }

    @Test
    public void testEnforceAlternateWithoutCache() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceDeniedNoDsvc() {
        final String dsvcId = "ds-geoblock-no-exist";
        final String url = "http://ds1.example.com/live1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.DENIED));
    }

    @Test
    public void testEnforceDeniedNoRegexMatch() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live-not-exist";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.DENIED));
    }

    @Test
    public void testEnforceAlternateToPathNoSlash() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live3";
        final String postal = "V5D";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/redirect_T3"));
    }

    @Test
    public void testEnforceAlternateNullPostal() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = null;
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceAlternateEmptyPostalInclude() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/live2";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/path/redirect_T2"));
    }

    @Test
    public void testEnforceAlternateEmptyPostalExclude() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceWhiteListAllowed() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = null;
        final String ip = "129.100.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceNotInWhiteListAlternate() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = "N7G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/redirect_T4"));
    }

    @Test
    public void testEnforceNotInWhiteListAllowedByPostal() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = "N6G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }
}

