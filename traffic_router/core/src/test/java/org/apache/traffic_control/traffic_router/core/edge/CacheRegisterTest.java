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

package org.apache.traffic_control.traffic_router.core.edge;

import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.JsonNodeFactory;
import com.fasterxml.jackson.databind.node.ObjectNode;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import org.apache.traffic_control.traffic_router.core.request.DNSRequest;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.apache.traffic_control.traffic_router.core.request.Request;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.junit.Before;
import org.junit.Test;
import org.xbill.DNS.Name;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.Type;

import java.util.HashMap;
import java.util.Map;
import java.util.TreeSet;

import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class CacheRegisterTest {
	private final CacheRegister cacheRegister = new CacheRegister();
	@Before
	public void before() {
		DeliveryService deliveryService1 = mock(DeliveryService.class);
		when(deliveryService1.getId()).thenReturn("delivery service 1");

		DeliveryService deliveryService2 = mock(DeliveryService.class);
		when(deliveryService2.getId()).thenReturn("delivery service 2");

		DeliveryServiceMatcher deliveryServiceMatcher1 = new DeliveryServiceMatcher(deliveryService1);
		deliveryServiceMatcher1.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher1.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher2 = new DeliveryServiceMatcher(deliveryService2);
		deliveryServiceMatcher2.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher2.addMatch(PATH, ".*abcde.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher3 = new DeliveryServiceMatcher(deliveryService2);
		deliveryServiceMatcher3.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher3.addMatch(PATH, ".*abcd.*", "");

		TreeSet<DeliveryServiceMatcher> deliveryServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
		deliveryServiceMatchers.add(deliveryServiceMatcher1);
		deliveryServiceMatchers.add(deliveryServiceMatcher2);
		deliveryServiceMatchers.add(deliveryServiceMatcher3);

		DeliveryServiceMatcher dnsMatcher1 = new DeliveryServiceMatcher(deliveryService1);
		dnsMatcher1.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");

		TreeSet<DeliveryServiceMatcher> dnsMatchers = new TreeSet<DeliveryServiceMatcher>();
		deliveryServiceMatchers.add(dnsMatcher1);

		cacheRegister.setDeliveryServiceMatchers(deliveryServiceMatchers);
	}

	@Test
	public void itPicksTheMostSpecificDeliveryService() {
		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.service01-kabletown.com");
		httpRequest.setPath("foo/abcde/bar");

		assertThat(cacheRegister.getDeliveryService(httpRequest).getId(), equalTo("delivery service 2"));

		Request request = new Request();
		request.setHostname("foo.service01-kabletown.com");
		assertThat(cacheRegister.getDeliveryService(request).getId(), equalTo("delivery service 1"));
	}

	@Test
	public void itReturnsNullForDeliveryServiceWhenItHasNoMatchers() {
		cacheRegister.setDeliveryServiceMatchers(null);

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.service01-kabletown.com");
		httpRequest.setPath("foo/abcde/bar");
		assertThat(cacheRegister.getDeliveryService(httpRequest), nullValue());
	}

	@Test
	public void itReturnsDeliveryServiceFromFQDNMapForHTTPRequest() throws JsonUtilsException {
		String requestName = "http://foo.service01.kabletown.com/";
		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.service01.kabletown.com");
		httpRequest.setRequestedUrl(requestName);
		Map<String, DeliveryService> map = new HashMap<>();

		ObjectNode node = JsonNodeFactory.instance.objectNode();
		ArrayNode domainNode = node.putArray("domains");
		domainNode.add("kabletown.com");
		node.put("routingName","foo");
		node.put("coverageZoneOnly", false);
		DeliveryService ds = new DeliveryService("service01", node);

		map.put("foo.service01.kabletown.com", ds);
		map.put("_.service01.kabletown.com", ds);
		cacheRegister.setFQDNToDeliveryServiceMap(map);

		DeliveryService answer = cacheRegister.getDeliveryService(httpRequest);
		assertThat("FQDNToDeliveryServiceMap was expected to have the key foo.service01.kabletown.com",
				cacheRegister.getFQDNToDeliveryServiceMap().containsKey("foo.service01.kabletown.com"));
		assertThat("Returned Delivery Service was expected to have the ID service01",
				answer.getId().equals("service01"));


		httpRequest.setRequestedUrl("http://_.service01.kabletown.com");
		answer = cacheRegister.getDeliveryService(httpRequest);
		assertThat("FQDNToDeliveryServiceMap was expected to have the key _.service01.kabletown.com",
				cacheRegister.getFQDNToDeliveryServiceMap().containsKey("_.service01.kabletown.com"));
		assertThat("Returned Delivery Service was expected to have the ID service01",
				answer.getId().equals("service01"));
	}

	@Test
	public void itReturnsDeliveryServiceFromFQDNMapForDNSRequest() throws JsonUtilsException, TextParseException {
		final Name name = Name.fromString("edge.example.com.");
		DNSRequest dnsRequest = new DNSRequest("example.com", name, Type.A);
		dnsRequest.setClientIP("10.10.10.10");
		dnsRequest.setHostname(name.relativize(Name.root).toString());

		Map<String, DeliveryService> map = new HashMap<>();

		ObjectNode node = JsonNodeFactory.instance.objectNode();
		ArrayNode domainNode = node.putArray("domains");
		domainNode.add("example.com");
		node.put("routingName","edge");
		node.put("coverageZoneOnly", false);
		DeliveryService ds = new DeliveryService("example", node);

		map.put("edge.example.com", ds);
		map.put("_.example.com", ds);
		cacheRegister.setFQDNToDeliveryServiceMap(map);

		DeliveryService answer = cacheRegister.getDeliveryService(dnsRequest);
		assertThat("FQDNToDeliveryServiceMap was expected to have the key edge.example.com",
				cacheRegister.getFQDNToDeliveryServiceMap().containsKey("edge.example.com"));
		assertThat("Returned Delivery Service was expected to have the ID example",
				answer.getId().equals("example"));

		final Name underscoreName = Name.fromString("_.example.com");
		dnsRequest = new DNSRequest("example.com", underscoreName, Type.A);
		dnsRequest.setClientIP("10.10.10.10");
		dnsRequest.setHostname(name.relativize(Name.root).toString());
		answer = cacheRegister.getDeliveryService(dnsRequest);
		assertThat("FQDNToDeliveryServiceMap was expected to have the key _.example.com",
				cacheRegister.getFQDNToDeliveryServiceMap().containsKey("_.example.com"));
		assertThat("Returned Delivery Service was expected to have the ID example",
				answer.getId().equals("example"));
	}
}
