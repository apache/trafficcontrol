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
import org.apache.traffic_control.traffic_router.core.router.HTTPRouteResult;
import org.apache.traffic_control.traffic_router.core.router.StatTracker;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationException;
import org.junit.Before;
import org.junit.Test;

import java.io.File;
import java.net.MalformedURLException;
import java.net.URL;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.powermock.api.mockito.PowerMockito.mock;
import static org.powermock.api.mockito.PowerMockito.when;

public class RegionalGeoTest {
    @Before
    public void setUp() throws Exception {
        final File dbFile = new File("src/test/resources/regional_geoblock.json");
        RegionalGeo.parseConfigFile(dbFile, false);
    }

    @Test
    public void testEnforceAllowedCoordinateRange() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = null;
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 12.0, 55.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceAlternateWithCacheNoCoordinateRangeNoPostalCode() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/live2";
        final String postal = null;
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 12.0, 55.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
    }

    @Test
    public void testEnforceAllowed() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceAlternateWithCache() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/live2";
        final String postal = "N7G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/path/redirect_T2"));
    }

    @Test
    public void testEnforceAlternateWithoutCache() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceDeniedNoDsvc() {
        final String dsvcId = "ds-geoblock-no-exist";
        final String url = "http://ds1.example.com/live1";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.DENIED));
    }

    @Test
    public void testEnforceDeniedNoRegexMatch() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live-not-exist";
        final String postal = "V5G";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.DENIED));
    }

    @Test
    public void testEnforceAlternateToPathNoSlash() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live3";
        final String postal = "V5D";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/redirect_T3"));
    }

    @Test
    public void testEnforceAlternateNullPostal() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = null;
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceAlternateEmptyPostalInclude() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds2.example.com/live2";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/path/redirect_T2"));
    }

    @Test
    public void testEnforceAlternateEmptyPostalExclude() {
        final String dsvcId = "ds-geoblock-exclude";
        final String url = "http://ds1.example.com/live1";
        final String postal = "";
        final String ip = "10.0.0.1";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("http://example.com/redirect_T1"));
    }

    @Test
    public void testEnforceWhiteListAllowed() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = null;
        final String ip = "129.100.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceAllowedHttpsRedirect() {
        final String dsvcId = "ds-geoblock-redirect-https";
        final String url = "http://ds1.example.com/httpsredirect";
        final String postal = null;
        final String ip = "129.100.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("https://example.com/redirect_https"));
    }

    @Test
    public void testEnforceSteeringReDirect() {
        final String dsvcId = "ds-steering-1";
        final String url = "http://ds1.example.com/steering";
        final String postal = null;
        final String ip = "129.100.254.4";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITHOUT_CACHE));
        assertThat(result.getUrl(), equalTo("https://example.com/steering-test"));
    }


    @Test
    public void testEnforceNotInWhiteListAlternate() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = "N7G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALTERNATE_WITH_CACHE));
        assertThat(result.getUrl(), equalTo("/redirect_T4"));
    }

    @Test
    public void testEnforceNotInWhiteListAllowedByPostal() {
        final String dsvcId = "ds-geoblock-include";
        final String url = "http://ds1.example.com/live4";
        final String postal = "N6G";
        final String ip = "129.202.254.2";

        RegionalGeoResult result = RegionalGeo.enforce(dsvcId, url, ip, postal, 0.0, 0.0);

        assertThat(result.getType(), equalTo(RegionalGeoResultType.ALLOWED));
        assertThat(result.getUrl(), equalTo(url));
    }

    @Test
    public void testEnforceWhiteListAllowedRouteResultMultipleUrls() throws GeolocationException, MalformedURLException {
        String clientIp = "129.100.254.2";
        String requestUrl = "http://ds1.example.com/live4";

        HTTPRequest request = new HTTPRequest();
        request.setClientIP(clientIp);
        request.setHostname("ds1.example.com");
        request.applyUrl(new URL(requestUrl));

        StatTracker.Track track = new StatTracker.Track();

        Cache cache = mock(Cache.class);

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn("ds-geoblock-include");
        when(ds.createURIString(request, cache)).thenReturn(requestUrl);

        TrafficRouter tr = mock(TrafficRouter.class);
        when(tr.getClientGeolocation(clientIp, track, ds)).thenReturn(new Geolocation(42, -71));

        HTTPRouteResult routeResult = new HTTPRouteResult(true);
        String firstUrl = "http://example.com/url1.m3u8";
        routeResult.addUrl(new URL(firstUrl));

        RegionalGeo.enforce(tr, request, ds, cache, routeResult, track);

        assertThat(routeResult.getUrls().size(), equalTo(2));
        assertThat(routeResult.getUrls().get(1).toString(), equalTo(requestUrl));
        assertThat(routeResult.getUrls().get(0).toString(), equalTo(firstUrl));

    }
}

