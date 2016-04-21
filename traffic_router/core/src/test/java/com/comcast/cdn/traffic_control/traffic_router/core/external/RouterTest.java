/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.catalina.LifecycleException;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.util.EntityUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.io.File;
import java.io.IOException;
import java.util.Iterator;

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;

@Category(ExternalTest.class)
public class RouterTest {
	private CloseableHttpClient httpClient;
	private String deliveryServiceId;

	@Before
	public void before() throws IOException, InterruptedException, LifecycleException {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

		System.out.println(System.getProperty("user.dir"));
		JsonNode jsonNode = objectMapper.readTree(new File("src/test/db/cr-config.json"));

		deliveryServiceId = null;

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext()) {
			String dsId = deliveryServices.next();
			Iterator<JsonNode> matchsets = jsonNode.get("deliveryServices").get(dsId).get("matchsets").iterator();
			while (matchsets.hasNext() && deliveryServiceId == null) {
				if ("HTTP".equals(matchsets.next().get("protocol").asText())) {
					deliveryServiceId = dsId;
				}
			}
		}

		assertThat(deliveryServiceId, not(nullValue()));

		httpClient = HttpClientBuilder.create().disableRedirectHandling().build();
	}

	@After
	public void after() throws IOException {
	 	httpClient.close();
	}

	@Test
	public void itHasAHomePage() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/index.html");

		CloseableHttpResponse response = null;
		try {
			response = httpClient.execute(httpGet);
			assertThat(EntityUtils.toString(response.getEntity()), containsString("This is a test!"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsValidRequests() throws IOException, InterruptedException {
		// Traffic Router will give us a 503 until it is ready to route
		// It also gives us a 503 when we don't make a valid routing request
		// The following request though *SHOULD* work so try and do this request multiple times
		// until we get a 302 to determine that all the application context is finished before
		// starting tests

		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		int triesLeft = 60;

		while (triesLeft > 0) {
			triesLeft--;
			try {
				response = httpClient.execute(httpGet);

				if (response.getStatusLine().getStatusCode() != 302) {
					Thread.sleep(500);
					continue;
				}

				triesLeft = 0;
			} finally {
				if (response != null) response.close();
			}
		}
	}

	@Test
	public void itRejectsInvalidRequests() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.invalid-delivery-service-id.bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}
}
