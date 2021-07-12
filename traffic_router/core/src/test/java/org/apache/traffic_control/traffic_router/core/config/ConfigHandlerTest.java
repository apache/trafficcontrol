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

package org.apache.traffic_control.traffic_router.core.config;

import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.router.StatTracker;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation.LocalizationMethod;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.mockito.Mockito.doAnswer;

import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;
import org.powermock.api.mockito.PowerMockito;


import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.powermock.reflect.Whitebox;


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

        when(register.getDeliveryService(any(HTTPRequest.class))).thenReturn(null);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo("NOT_DS_URL"));
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
        when(ds.getId()).thenReturn(anotherId);
        CacheRegister register = PowerMockito.mock(CacheRegister.class);

        when(register.getDeliveryService(any(HTTPRequest.class))).thenReturn(anotherDs);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo("NOT_DS_URL"));
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

        when(register.getDeliveryService(any(HTTPRequest.class))).thenReturn(ds);

        Whitebox.invokeMethod(handler, "initGeoFailedRedirect", dsMap, register);
        assertThat(urlType[0], equalTo("DS_URL"));
        assertThat(typeUrl[0], equalTo(path));
    }

    @Test
    public void itParsesTheTopologiesConfig() throws Exception {
        /* Make the CacheLocation, add a Cache, and add the CacheLocation to the CacheRegister */
        final String cacheId = "edge";
        final Cache cache = new Cache(cacheId, cacheId, 0);
        final String location = "CDN_in_a_Box_Edge";
        final CacheLocation cacheLocation = new CacheLocation(location, new Geolocation(38.897663, 38.897663));
        cacheLocation.addCache(cache);
        final Set<CacheLocation> locations = new HashSet<>();
        locations.add(cacheLocation);
        final CacheRegister register = new CacheRegister();
        register.setConfiguredLocations(locations);

        /* Add a capability to the Cache */
        final String capability = "a-capability";
        final Set<String> capabilities = new HashSet<>();
        capabilities.add(capability);
        cache.addCapabilities(capabilities);

        /* Mock a DeliveryService and add it to our DeliveryService Map */
        final String dsId = "top-ds";
        final String routingName = "cdn";
        final String domain = "ds.site.com";
        final String topology = "foo";
        final String superHackedRegexp = "(.*\\.|^)" + dsId + "\\..*";
        final DeliveryService ds = mock(DeliveryService.class);
        when(ds.getId()).thenReturn(dsId);
        when(ds.getDomain()).thenReturn(domain);
        when(ds.getRemap(superHackedRegexp)).thenReturn(domain);
        when(ds.getRoutingName()).thenReturn(routingName);
        when(ds.getTopology()).thenReturn(topology);
        when(ds.hasRequiredCapabilities(capabilities)).thenReturn(true);
        when(ds.isDns()).thenReturn(false);
        final Map<String, DeliveryService> dsMap = new HashMap<>();
        dsMap.put(dsId, ds);

        final DeliveryServiceMatcher dsMatcher = new DeliveryServiceMatcher(ds);
        dsMatcher.addMatch(DeliveryServiceMatcher.Type.HOST, superHackedRegexp, "");
        final TreeSet<DeliveryServiceMatcher> dsMatchers = new TreeSet<>();
        dsMatchers.add(dsMatcher);
        register.setDeliveryServiceMap(dsMap);
        register.setDeliveryServiceMatchers(dsMatchers);

        /* Parse the Topologies config JSON */
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode allTopologiesJson = mapper.readTree("{\"" + topology + "\":{\"nodes\":[\"" + location + "\"]}}");
        Whitebox.setInternalState(handler, "statTracker", new StatTracker());
        Whitebox.invokeMethod(handler, "parseTopologyConfig", allTopologiesJson, dsMap, register);

        /* Assert that the DeliveryService was assigned to the Cache */
        Collection<Cache.DeliveryServiceReference> dsReferences = cache.getDeliveryServices();
        assertThat(dsReferences.size(), equalTo(1));
        assertThat(dsReferences.iterator().next().getDeliveryServiceId(), equalTo(dsId));
    }

    @Test
    public void testParseLocalizationMethods() throws Exception {
        LocalizationMethod[] allMethods = new LocalizationMethod[] {
                LocalizationMethod.CZ,
                LocalizationMethod.DEEP_CZ,
                LocalizationMethod.GEO,
        };
        Set<LocalizationMethod> expected = new HashSet<>();
        expected.addAll(Arrays.asList(allMethods));

        ObjectMapper mapper = new ObjectMapper();

        String allMethodsString = "{\"localizationMethods\": [\"CZ\",\"DEEP_CZ\",\"GEO\"]}";
        JsonNode allMethodsJson = mapper.readTree(allMethodsString);
        Set<LocalizationMethod> actual = Whitebox.invokeMethod(handler, "parseLocalizationMethods", "foo", allMethodsJson);
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
}
