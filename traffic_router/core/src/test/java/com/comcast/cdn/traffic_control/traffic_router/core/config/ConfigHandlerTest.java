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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheLocation.LocalizationMethod;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpConfigUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPoller;
import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPublisher;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Before;
import org.junit.Test;
import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.reflect.Whitebox;

import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyBoolean;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.anyLong;
import static org.mockito.Mockito.*;
import static org.mockito.Mockito.eq;


public class ConfigHandlerTest {
    private ConfigHandler handler;

    @Before
    public void before() throws Exception {
        handler = mock(ConfigHandler.class);
    }

    @Test
    public void itTestRelativeUrl() throws Exception {
        final String redirectUrl = "relative/url";
        final String dsId = "relative-url";
        final String[] urlType = {""};
        final String[] typeUrl = {""};

        Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getGeoRedirectUrl()).thenReturn(redirectUrl);
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                typeUrl[0] = (String)(args[0]);
                return null;
            }
        }).when(ds).setGeoRedirectFile(anyString());
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                urlType[0] = (String)args[0];
                return null;
            }
        }).when(ds).setGeoRedirectUrlType(anyString());

        dsMap.put(dsId, ds);
        
        CacheRegister register = PowerMockito.mock(CacheRegister.class);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo("DS_URL"));
        assertThat(typeUrl[0], equalTo(""));
    }

    @Test
    public void itTestRelativeUrlNegative() throws Exception {
        final String redirectUrl = "://invalid";
        final String dsId = "relative-url";
        final String[] urlType = {""};
        final String[] typeUrl = {""};

        Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getGeoRedirectUrl()).thenReturn(redirectUrl);

        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                typeUrl[0] = (String)(args[0]);
                return null;
            }
        }).when(ds).setGeoRedirectFile(anyString());
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                urlType[0] = (String)args[0];
                return null;
            }
        }).when(ds).setGeoRedirectUrlType(anyString());

        dsMap.put(dsId, ds);

        CacheRegister register = PowerMockito.mock(CacheRegister.class);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo(""));
        assertThat(typeUrl[0], equalTo(""));
    }

    @Test
    public void itTestNoSuchDsUrl() throws Exception {
        final String path = "/ds/url";
        final String redirectUrl = "http://test.com" + path;
        final String dsId = "relative-url";
        final String[] urlType = {""};
        final String[] typeUrl = {""};

        Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getGeoRedirectUrl()).thenReturn(redirectUrl);
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                typeUrl[0] = (String)(args[0]);
                return null;
            }
        }).when(ds).setGeoRedirectFile(anyString());
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                urlType[0] = (String)args[0];
                return null;
            }
        }).when(ds).setGeoRedirectUrlType(anyString());

        dsMap.put(dsId, ds);

        CacheRegister register = PowerMockito.mock(CacheRegister.class);

        when(register.getDeliveryService(any(HTTPRequest.class), anyBoolean())).thenReturn(null);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo(ConfigHandler.NOT_DS_URL));
        assertThat(typeUrl[0], equalTo(path));
    }

    @Test
    public void itTestNotThisDsUrl() throws Exception {
        final String path = "/ds/url";
        final String redirectUrl = "http://test.com" + path;
        final String dsId = "relative-ds";
        final String anotherId = "another-ds";
        final String[] urlType = {""};
        final String[] typeUrl = {""};

        Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getGeoRedirectUrl()).thenReturn(redirectUrl);
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                typeUrl[0] = (String)(args[0]);
                return null;
            }
        }).when(ds).setGeoRedirectFile(anyString());
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                urlType[0] = (String)args[0];
                return null;
            }
        }).when(ds).setGeoRedirectUrlType(anyString());

        dsMap.put(dsId, ds);

        DeliveryService anotherDs = mock(DeliveryService.class);
        when(anotherDs.getId()).thenReturn(anotherId);
        CacheRegister register = PowerMockito.mock(CacheRegister.class);

        when(register.getDeliveryService(any(HTTPRequest.class), anyBoolean())).thenReturn(anotherDs);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo(ConfigHandler.NOT_DS_URL));
        assertThat(typeUrl[0], equalTo(path));
    }

    @Test
    public void itTestThisDsUrl() throws Exception {
        final String path = "/ds/url";
        final String redirectUrl = "http://test.com" + path;
        final String dsId = "relative-ds";
        final String[] urlType = {""};
        final String[] typeUrl = {""};

        Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

        DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getGeoRedirectUrl()).thenReturn(redirectUrl);
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                typeUrl[0] = (String)(args[0]);
                return null;
            }
        }).when(ds).setGeoRedirectFile(anyString());
        doAnswer(new Answer<Void>() {
            public Void answer(InvocationOnMock invocation) {
                Object[] args = invocation.getArguments();
                urlType[0] = (String)args[0];
                return null;
            }
        }).when(ds).setGeoRedirectUrlType(anyString());

        dsMap.put(dsId, ds);

        CacheRegister register = PowerMockito.mock(CacheRegister.class);
        when(register.getDeliveryService(any(HTTPRequest.class), anyBoolean())).thenReturn(ds);
        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo("DS_URL"));
        assertThat(typeUrl[0], equalTo(path));
    }

	@Test
	public void parseRegionalGeoConfig() throws Exception {
		final long interval = 60000000000l;
		final String url =  "http://testing-tm-01.cdn.example";
		JsonNode config = mock(JsonNode.class);
		JsonNode geoInfo = mock(JsonNode.class);
		when(geoInfo.asText(anyString())).thenReturn(url);
		when(geoInfo.asText()).thenReturn(url);
		when(geoInfo.asLong(0)).thenReturn(interval);
		when(geoInfo.asLong()).thenReturn(interval);
		when(config.has(anyString())).thenReturn(true);
		when(config.get(anyString())).thenReturn(geoInfo);
		RegionalGeoUpdater rgu = mock(RegionalGeoUpdater.class);
		when(handler.getRegionalGeoUpdater()).thenReturn(rgu);
		SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
		DeliveryService ds = mock(DeliveryService.class);
		when(ds.isRegionalGeoEnabled()).thenReturn(true);
		Map<String, DeliveryService> modList = new HashMap<>();
		Map<String, DeliveryService> empty = new HashMap<>();
		modList.put("testDs", ds);
		when(snapshotEventsProcessor.getCreationEvents()).thenReturn(modList);
		when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(empty);
		Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
		verify(rgu).setDataBaseURL(eq(url),eq(interval ));
		when(snapshotEventsProcessor.getCreationEvents()).thenReturn(empty);
		when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(modList);
		Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
		verify(rgu, times(2)).setDataBaseURL(eq(url),eq(interval ));
		when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(empty);
		Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
		verify(rgu, times(2)).setDataBaseURL(eq(url),eq(interval ));
	}

    @Test
    public void testParseLocalizationMethods() throws Exception {
	    LocalizationMethod[] allMethods = new LocalizationMethod[]{
		    LocalizationMethod.CZ,
		    LocalizationMethod.DEEP_CZ,
		    LocalizationMethod.GEO,
	    };
	    Set<LocalizationMethod> expected = new HashSet<>();
	    expected.addAll(Arrays.asList(allMethods));

	    ObjectMapper mapper = new ObjectMapper();

	    String allMethodsString = "{\"localizationMethods\": [\"CZ\",\"DEEP_CZ\",\"GEO\"]}";
	    JsonNode allMethodsJson = mapper.readTree(allMethodsString);
	    Set<LocalizationMethod> actual = Whitebox
			    .invokeMethod(handler, "parseLocalizationMethods", "foo", allMethodsJson);
	    assertThat(actual, equalTo(expected));

	    String noMethodsString = "{}";
	    JsonNode noMethodsJson = mapper.readTree(noMethodsString);
	    actual = Whitebox.invokeMethod(handler, "parseLocalizationMethods", "foo", noMethodsJson);
	    assertThat(actual, equalTo(expected));

	    String nullMethodsString = "{\"localizationMethods\": null}";
	    JsonNode nullMethodsJson = mapper.readTree(nullMethodsString);
	    actual = Whitebox.invokeMethod(handler, "parseLocalizationMethods", "foo", nullMethodsJson);
	    assertThat(actual, equalTo(expected));

	    String CZMethodsString = "{\"localizationMethods\": [\"CZ\"]}";
	    JsonNode CZMethodsJson = mapper.readTree(CZMethodsString);
	    expected.clear();
	    expected.add(LocalizationMethod.CZ);
	    actual = Whitebox.invokeMethod(handler, "parseLocalizationMethods", "foo", CZMethodsJson);
	    assertThat(actual, equalTo(expected));
    }

    @Test
    public void parseRegionalGeoConfigNegative() throws Exception {
    	final long interval = 60000000000l;
    	final String url =  "http://testing-tm-01.cdn.example";
    	JsonNode config = mock(JsonNode.class);
	    JsonNode geoInfo = mock(JsonNode.class);
	    when(geoInfo.asText(anyString())).thenReturn(url);
	    when(geoInfo.asText()).thenReturn(url);
	    when(geoInfo.asLong(0)).thenReturn(interval);
	    when(geoInfo.asLong()).thenReturn(interval);
    	when(config.has(anyString())).thenReturn(true);
    	when(config.get(anyString())).thenReturn(geoInfo);
	    RegionalGeoUpdater rgu = mock(RegionalGeoUpdater.class);
        when(handler.getRegionalGeoUpdater()).thenReturn(rgu);
        SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
        DeliveryService nds = mock(DeliveryService.class);
	    when(nds.isRegionalGeoEnabled()).thenReturn(false);
        DeliveryService ds = mock(DeliveryService.class);
        when(ds.isRegionalGeoEnabled()).thenReturn(true);
        Map<String, DeliveryService> negativeList = new HashMap<>();
	    Map<String, DeliveryService> modList = new HashMap<>();
	    Map<String, DeliveryService> empty = new HashMap<>();
        modList.put("testDs", ds);
        negativeList.put("negTestDs",nds);
        when(snapshotEventsProcessor.getCreationEvents()).thenReturn(empty);
	    when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(empty);
        Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
        verify(rgu).cancelServiceUpdater();
	    when(snapshotEventsProcessor.getCreationEvents()).thenReturn(negativeList);
	    Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
	    verify(rgu, times(2)).cancelServiceUpdater();
	    when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(modList);
	    Whitebox.invokeMethod(handler, "parseRegionalGeoConfig", config, snapshotEventsProcessor);
	    verify(rgu).setDataBaseURL(eq(url),eq(interval ));
    };

    @Test
    public void parseAnonymousIpConfig() throws Exception {
	    final long interval = 60000000000l;
	    final String pollingUrl =  "http://testing-tm-01.cdn.example";
	    final String policyConfig =  "http://testing-config-01.cdn.example";
	    JsonNode config = mock(JsonNode.class);
	    JsonNode apcNode = mock(JsonNode.class);
	    JsonNode apuNode = mock(JsonNode.class);
	    when(apcNode.asText(anyString())).thenReturn(policyConfig);
	    when(apcNode.asText()).thenReturn(policyConfig);
	    when(apuNode.asText(anyString())).thenReturn(pollingUrl);
	    when(apuNode.asText()).thenReturn(pollingUrl);
	    when(apuNode.asLong()).thenReturn(interval);
	    when(apuNode.asLong(0l)).thenReturn(interval);
	    when(config.has(anyString())).thenReturn(true);
	    when(config.get("anonymousip.policy.configuration")).thenReturn(apcNode);
	    when(config.get("anonymousip.polling.url")).thenReturn(apuNode);
	    when(config.get("anonymousip.polling.interval")).thenReturn(apuNode);
	    SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
	    DeliveryService ds = mock(DeliveryService.class);
	    when(ds.isAnonymousIpEnabled()).thenReturn(true);
	    Map<String, DeliveryService> createList = new HashMap<>();
	    Map<String, DeliveryService> empty = new HashMap<>();
	    createList.put("testDs", ds);
	    when(snapshotEventsProcessor.getCreationEvents()).thenReturn(createList);
	    when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(empty);
	    AnonymousIpConfigUpdater anonymousIpConfigUpdater = mock(AnonymousIpConfigUpdater.class);
	    when(handler.getAnonymousIpConfigUpdater()).thenReturn(anonymousIpConfigUpdater);
	    AnonymousIpDatabaseUpdater anonymousIpDatabaseUpdater = mock(AnonymousIpDatabaseUpdater.class);
	    when(handler.getAnonymousIpDatabaseUpdater()).thenReturn(anonymousIpDatabaseUpdater);
        Whitebox.invokeMethod(handler, "parseAnonymousIpConfig", config, snapshotEventsProcessor);
        verify(anonymousIpConfigUpdater).setDataBaseURL(anyString(), anyLong());
        verify(anonymousIpDatabaseUpdater).setDataBaseURL(anyString(), anyLong());
    };

    @Test
    public void updateCertsPublisher() throws Exception {

	    CertificatesPublisher certificatesPublisher = mock(CertificatesPublisher.class);
	    when(handler.getCertificatesPublisher()).thenReturn(certificatesPublisher);
	    CertificatesPoller certificatesPoller = mock(CertificatesPoller.class);
	    when(handler.getCertificatesPoller()).thenReturn(certificatesPoller);
        SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
	    DeliveryService ds = mock(DeliveryService.class);
	    Map<String, DeliveryService> createList = new HashMap<>();
	    Map<String, DeliveryService> empty = new HashMap<>();
	    createList.put("testDs", ds);
	    when(snapshotEventsProcessor.getCreationEvents()).thenReturn(createList);
	    when(snapshotEventsProcessor.getUpdateEvents()).thenReturn(empty);
        Whitebox.invokeMethod(handler, "updateCertsPublisher", snapshotEventsProcessor);
	    verify(certificatesPoller).restart();
    };

    @Test
    // updates, creates and removes the DeliveryServices in cacheRegister
    public void parseDeliveryServiceMatchSets() throws Exception {

        JsonNode config = mock(JsonNode.class);
        SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
        CacheRegister cacheRegister = mock(CacheRegister.class);
        Whitebox.invokeMethod(handler, "parseDeliveryServiceMatchSets", snapshotEventsProcessor, cacheRegister);
    };

    @Test
    public void parseCacheConfig() throws Exception {

        JsonNode contentServers = mock(JsonNode.class);
        SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
	    CacheRegister cacheRegister = mock(CacheRegister.class);
        Whitebox.invokeMethod(handler, "parseCacheConfig", snapshotEventsProcessor, contentServers, cacheRegister);
    };

}
