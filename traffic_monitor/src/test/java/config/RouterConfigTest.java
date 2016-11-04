package config;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.Peer;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Network;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.contains;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@PrepareForTest({RouterConfig.class, Cache.class, Network.class})
@RunWith(PowerMockRunner.class)
public class RouterConfigTest {

	private JSONObject crConfigJson;
	private Cache cache1;
	private Cache cache2;

	@Before
	public void before() throws Exception {
		JSONObject cache1Json = new JSONObject();
		cache1Json.put("foo", "bar");

		JSONObject cache2Json = new JSONObject();
		cache2Json.put("foo", "bar");

		JSONObject contentServersJson = new JSONObject();
		contentServersJson.put("cache1", cache1Json);
		contentServersJson.put("cache2", cache2Json);

		crConfigJson = new JSONObject();
		crConfigJson.put("contentServers", contentServersJson);

		cache1 = mock(Cache.class);
		when(cache1.getHostname()).thenReturn("cache1");

		cache2 = mock(Cache.class);
		when(cache2.getHostname()).thenReturn("cache2");

		whenNew(Cache.class).withArguments("cache1", cache1Json).thenReturn(cache1);
		whenNew(Cache.class).withArguments("cache2", cache2Json).thenReturn(cache2);
	}

	@Test
	public void itCreatesListOfCachesFromJson() throws Exception {
		RouterConfig routerConfig = new RouterConfig(crConfigJson, new HealthDeterminer());
		assertThat(routerConfig.getCacheList(), containsInAnyOrder(cache1, cache2));
	}

	@Test
	public void itCreatesPeerMapFromJson() throws Exception {
		Peer peer1 = mock(Peer.class);
		when(peer1.getIpAddress()).thenReturn("192.168.1.2");

		mockStatic(Network.class);
		when(Network.isIpAddressLocal("192.168.1.2")).thenReturn(true);

		Peer peer2 = mock(Peer.class);
		when(peer2.getIpAddress()).thenReturn("192.168.10.20");
		when(peer2.getFqdn()).thenReturn("peer2.kabletown.com");

		when(Network.isIpAddressLocal("192.168.10.20")).thenReturn(false);
		when(Network.isLocalName("peer2.kabletown.com")).thenReturn(true);

		Peer peer3 = mock(Peer.class);
		when(peer3.getIpAddress()).thenReturn("192.168.10.30");
		when(peer3.getFqdn()).thenReturn("peer3.kabletown.com");
		when(peer3.getId()).thenReturn("peer3");

		when(Network.isIpAddressLocal("192.168.10.30")).thenReturn(false);
		when(Network.isLocalName("peer3.kabletown.com")).thenReturn(false);
		when(Network.isLocalName("peer3")).thenReturn(true);

		Peer peer4 = mock(Peer.class);
		when(peer4.getIpAddress()).thenReturn("192.168.10.40");
		when(peer4.getFqdn()).thenReturn("peer4.kabletown.com");
		when(peer4.getId()).thenReturn("peer4");
		when(peer4.getStatus()).thenReturn("somethingelse");

		when(Network.isIpAddressLocal("192.168.10.40")).thenReturn(false);
		when(Network.isLocalName("peer4.kabletown.com")).thenReturn(false);
		when(Network.isLocalName("peer4")).thenReturn(false);

		Peer peer5 = mock(Peer.class);
		when(peer5.getIpAddress()).thenReturn("192.168.10.50");
		when(peer5.getFqdn()).thenReturn("peer5.kabletown.com");
		when(peer5.getId()).thenReturn("peer5");
		when(peer5.getStatus()).thenReturn("ONLINE");

		when(Network.isIpAddressLocal("192.168.10.50")).thenReturn(false);
		when(Network.isLocalName("peer5.kabletown.com")).thenReturn(false);
		when(Network.isLocalName("peer5")).thenReturn(false);

		JSONObject peer1Json = new JSONObject();
		peer1Json.put("fqdn", "peer1.kabletown.com");

		JSONObject peer2Json = new JSONObject();
		peer2Json.put("fqdn", "peer2.kabletown.com");

		JSONObject peer3Json = new JSONObject();
		peer3Json.put("fqdn", "peer3.kabletown.com");

		JSONObject peer4Json = new JSONObject();
		peer4Json.put("fqdn", "peer4.kabletown.com");

		JSONObject peer5Json = new JSONObject();
		peer5Json.put("fqdn", "peer5.kabletown.com");

		JSONObject monitorsJson = new JSONObject();
		monitorsJson.put("peer1", peer1Json);
		monitorsJson.put("peer2", peer2Json);
		monitorsJson.put("peer3", peer3Json);
		monitorsJson.put("peer4", peer4Json);
		monitorsJson.put("peer5", peer5Json);

		whenNew(Peer.class).withArguments("peer1", peer1Json).thenReturn(peer1);
		whenNew(Peer.class).withArguments("peer2", peer2Json).thenReturn(peer2);
		whenNew(Peer.class).withArguments("peer3", peer3Json).thenReturn(peer3);
		whenNew(Peer.class).withArguments("peer4", peer4Json).thenReturn(peer4);
		whenNew(Peer.class).withArguments("peer5", peer5Json).thenReturn(peer5);

		crConfigJson.put("monitors", monitorsJson);

		RouterConfig routerConfig = new RouterConfig(crConfigJson, new HealthDeterminer());
		assertThat(routerConfig.getPeerMap().keySet(), contains("peer5"));
	}

	@Test
	public void itReturnsDeliveryServicesJson() throws Exception {
		JSONObject deliveryServicesJson = new JSONObject().put("foo", "bar");
		crConfigJson.put("deliveryServices", deliveryServicesJson);

		assertThat(new RouterConfig(crConfigJson, new HealthDeterminer()).getDsList().getString("foo"), equalTo("bar"));
	}
}
