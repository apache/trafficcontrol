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
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.isEmptyString;
import static org.hamcrest.Matchers.nullValue;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;

public class CacheTest {

	private JSONObject cacheJson;

	@Before
	public void before() throws Exception {
		cacheJson = new JSONObject()
			.put("ip", "192.168.10.11")
			.put("status", "thestatus")
			.put("locationId", "the location id")
			.put("profile", "the cache profile")
			.put("fqdn", "cache1.kabletown.com")
			.put("type", "cache-type")
			.put("port", 1234);

	}
	@Test
	public void itCreatesCacheFromRequiredJson() throws Exception {
		Cache cache = new Cache("cache1", cacheJson);

		assertThat(cache.getHostname(), equalTo("cache1"));
		assertThat(cache.getIpAddress(), equalTo("192.168.10.11"));
		assertThat(cache.getInterfaceName(), isEmptyString());
		assertThat(cache.getStatus(), equalTo("thestatus"));
		assertThat(cache.getLocation(), equalTo("the location id"));
		assertThat(cache.getState(), nullValue());
		assertThat(cache.isAvailableKnown(), equalTo(false));
		assertThat(cache.isAvailable(), equalTo(true));
		assertThat(cache.getQueryIp(), equalTo("192.168.10.11"));
		assertThat(cache.getQueryPort(), equalTo(1234));
		assertThat(cache.getIp(), equalTo("192.168.10.11"));
		assertThat(cache.getType(), equalTo("cache-type"));
		assertThat(cache.getIp6(), isEmptyString());
		assertThat(cache.getControls(), nullValue());

		try {
			cache.getHistoryTime();
			fail("Should have thrown NullPointerException");
		} catch (NullPointerException e) {
			// expected
		}

		assertThat(cache.getProfile(), equalTo("the cache profile"));
		assertThat(cache.getFqdn(), equalTo("cache1.kabletown.com"));
		assertThat(cache.getDeliveryServices(), nullValue());

		try {
			cache.getDeliveryServiceIds();
			fail("getDeliveryServiceIds should have thrown null pointer exception");
		} catch (NullPointerException e) {
			// expected
		}

		try {
			cache.getFqdns("delivery-service-id");
			fail("getFqnds(deliveryServiceid) should have thrown null pointer exception");
		} catch (NullPointerException e) {
			// expected
		}
	}

	@Test
	public void itCreatesCacheWithOptionalJson() throws Exception {
		cacheJson
			.put("interfaceName", "eth0")
			.put("ip6", "fde5:acbf:f329::/48")
			.put("queryIp", "192.168.99.111")
			.put("queryPort", 4321)
			.put("deliveryServices", new JSONObject()
				.put("ds1", "something")
				.put("ds2", "somethingelse"));

		Cache cache = new Cache("cacheId", cacheJson);
		assertThat(cache.getInterfaceName(), equalTo("eth0"));
		assertThat(cache.getIp6(), equalTo("fde5:acbf:f329::/48"));
		assertThat(cache.getQueryIp(), equalTo("192.168.99.111"));
		assertThat(cache.getQueryPort(), equalTo(4321));
		assertThat(cache.getDeliveryServiceIds(), containsInAnyOrder("ds1", "ds2"));
	}

	@Test
	public void itUpdatesHealthDeterminerWithState() throws Exception {
		Cache cache = new Cache("cacheId", cacheJson);

		HealthDeterminer healthDeterminer = spy(new HealthDeterminer());

		final CacheState cacheState = mock(CacheState.class);
		cache.setState(cacheState, healthDeterminer);

		verify(healthDeterminer).setIsAvailable(cache, cacheState);
	}

	@Test
	public void itUpdatesHealthDeterminerWithError() throws Exception {
		HealthDeterminer healthDeterminer = spy(new HealthDeterminer());
		final CacheState cacheState = mock(CacheState.class);
		String errorString = "something went wrong";

		Cache cache = new Cache("cacheId", cacheJson);
		cache.setError(cacheState, errorString, healthDeterminer);

		verify(healthDeterminer).setIsAvailable(cache, "something went wrong", cacheState);
	}

	@Test
	public void itUsesStateToReportAvailabilityKnown() throws Exception {
		CacheState cacheState = spy(new CacheState("stateid"));
		Cache cache = new Cache("cacheId", cacheJson);

		cache.setCacheState(cacheState);
		assertThat(cache.isAvailableKnown(), equalTo(false));
		verify(cacheState).hasValue("isAvailable");

		doReturn(true).when(cacheState).hasValue("isAvailable");
		assertThat(cache.isAvailableKnown(), equalTo(true));
	}

	@Test
	public void itUsesStateToReportAvailability() throws Exception {
		CacheState cacheState = spy(new CacheState("stateid"));
		Cache cache = new Cache("cacheId", cacheJson);

		cache.setCacheState(cacheState);
		assertThat(cache.isAvailable(), equalTo(true));
		verify(cacheState).hasValue("isAvailable");

		doReturn(true).when(cacheState).hasValue("isAvailable");
		doReturn("").when(cacheState).getLastValue("isAvailable");
		assertThat(cache.isAvailable(), equalTo(false));

		doReturn("true").when(cacheState).getLastValue("isAvailable");
		assertThat(cache.isAvailable(), equalTo(true));
	}

	@Test
	public void itGetsControlsFromHealthDeterminer() throws Exception {
		Cache cache = new Cache("cacheid", cacheJson);

		JSONObject controls = new JSONObject();
		HealthDeterminer healthDeterminer = spy(new HealthDeterminer());
		doReturn(controls).when(healthDeterminer).getControls(cache);

		cache.setControls(healthDeterminer);

		assertThat(cache.getControls(), equalTo(controls));
	}

	@Test
	public void itGetsFqdnsForDeliveryServiceId() throws Exception {
		cacheJson
			.put("deliveryServices", new JSONObject()
				.put("ds1", new JSONArray(new Object[]{"foo", "bar", "baz"}))
				.put("ds2", "somethingelse"));

		Cache cache = new Cache("cacheId", cacheJson);

		try {
			cache.getFqdns("ds3");
			fail("Should have caught JSON Exception");
		} catch (JSONException e) {
			assertThat(e.getMessage(), equalTo("JSONObject[\"ds3\"] not found."));
		}

		assertThat(cache.getFqdns("ds1"), containsInAnyOrder("foo", "bar", "baz"));
		assertThat(cache.getFqdns("ds2"), containsInAnyOrder("somethingelse"));
	}
}
