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
import org.apache.http.Header;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustSelfSignedStrategy;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.ssl.SSLContextBuilder;
import org.apache.http.util.EntityUtils;
import org.hamcrest.Matchers;
import org.junit.After;
import org.junit.Before;
import org.junit.FixMethodOrder;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.junit.runners.MethodSorters;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SNIHostName;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLHandshakeException;
import javax.net.ssl.SSLParameters;
import javax.net.ssl.SSLSession;
import javax.net.ssl.SSLSocket;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.CoreMatchers.startsWith;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.endsWith;
import static org.hamcrest.Matchers.isIn;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class RouterTest {
	private CloseableHttpClient httpClient;
	private String deliveryServiceId;
	private List<String> validLocations = new ArrayList<>();
	private String deliveryServiceDomain;
	private final String httpsOnlyId = "https-only-test";
	private final String secureNoCertId = "https-nocert";
	private List<String> httpsOnlyLocations = new ArrayList<>();
	private List<String> noCertValidLocations = new ArrayList<>();
	private String httpsOnlyDomain = httpsOnlyId + ".thecdn.example.com";
	private String noCertsDeliveryServiceDomain = "https-nocert.thecdn.example.com";
	private String routerHttpPort = System.getProperty("routerHttpPort", "8888");
	private String routerSecurePort = System.getProperty("routerSecurePort", "8443");
	
	@Before
	public void before() throws Exception {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

		String resourcePath = "internal/api/1.2/steering.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode steeringNode = objectMapper.readTree(inputStream).get("response").get(0);

		String steeringDeliveryServiceId = steeringNode.get("deliveryService").asText();

		resourcePath = "publish/CrConfig.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = objectMapper.readTree(inputStream);

		deliveryServiceId = null;

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext()) {
			String dsId = deliveryServices.next();

			if (dsId.equals(steeringDeliveryServiceId)) {
				continue;
			}

			JsonNode deliveryServiceNode = jsonNode.get("deliveryServices").get(dsId);
			Iterator<JsonNode> matchsets = deliveryServiceNode.get("matchsets").iterator();

			while (matchsets.hasNext() && deliveryServiceId == null) {
				if ("HTTP".equals(matchsets.next().get("protocol").asText())) {
					final boolean sslEnabled = deliveryServiceNode.get("sslEnabled").asBoolean(false);
					if (!sslEnabled) {
						deliveryServiceId = dsId;
						deliveryServiceDomain = deliveryServiceNode.get("domains").get(0).asText();
					}
				}
			}
		}

		assertThat(deliveryServiceId, not(nullValue()));
		assertThat(deliveryServiceDomain, not(nullValue()));
		assertThat(httpsOnlyId, not(nullValue()));
		assertThat(httpsOnlyDomain, not(nullValue()));

		Iterator<String> cacheIds = jsonNode.get("contentServers").fieldNames();
		while (cacheIds.hasNext()) {
			String cacheId = cacheIds.next();
			JsonNode cacheNode = jsonNode.get("contentServers").get(cacheId);

			if (!cacheNode.has("deliveryServices")) {
				continue;
			}

			if (cacheNode.get("deliveryServices").has(deliveryServiceId)) {
				int port = cacheNode.get("port").asInt();
				String portText = (port == 80) ? "" : ":" + port;
				validLocations.add("http://" + cacheId + "." + deliveryServiceDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(httpsOnlyId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpsOnlyLocations.add("https://" + cacheId + "." + httpsOnlyDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(secureNoCertId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				noCertValidLocations.add("https://" + cacheId + "." + noCertsDeliveryServiceDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}
		}

		assertThat(validLocations.isEmpty(), equalTo(false));
		assertThat(httpsOnlyLocations.isEmpty(), equalTo(false));

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.https-only-test.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();
	}

	@After
	public void after() throws IOException {
	 	httpClient.close();
	}

	@Test
	public void itHasAHomePage() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/index.html");

		CloseableHttpResponse response = null;
		try {
			response = httpClient.execute(httpGet);
			assertThat(EntityUtils.toString(response.getEntity()), containsString("This is a test!"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsValidHttpRequests() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(validLocations));
			assertThat(header.getValue(), Matchers.startsWith("http://"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itDoesRoutingThroughPathsStartingWithCrs() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/crs/stats?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itConsistentlyRedirectsValidRequests() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();

			response.close();

			for (int i = 0; i < 100; i++) {
				response = httpClient.execute(httpGet);
				assertThat(response.getFirstHeader("Location").getValue(), equalTo(location));
				response.close();
			}
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRejectsInvalidRequests() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.invalid-delivery-service-id.bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsHttpsRequests() throws Exception {
		HttpGet httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpsOnlyId + ".thecdn.example.com");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpsOnlyLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(httpsOnlyId + ".thecdn.example.com/stuff"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRejectsHttpRequestsForHttpsOnlyDeliveryService() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpsOnlyId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	@Ignore // Ignore this test until we add explicit support for 'http to https' delivery service
	@Test
	public void itRedirectsFromHttpToHttps() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpsOnlyId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpsOnlyLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(httpsOnlyId + ".thecdn.example.com/stuff"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRejectsHttpsRequestsForHttpDeliveryService() throws Exception {
		HttpGet httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void z_itUpdatesCertsFromTrafficOps() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");

		httpGet.addHeader("Host", "tr.https-nocert.bar");
		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		}

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.https-nocert.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();

		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr.https-nocert.bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			fail("Should have gotten an SSL handshake failure!");
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} catch (SSLHandshakeException e) {
			assertThat(e.getMessage(), equalTo("Received fatal alert: handshake_failure"));
		}

		// Update certificates
		String testHttpPort = System.getProperty("testHttpServerPort", "8889");
		HttpPost httpPost = new HttpPost("http://localhost:"+ testHttpPort + "/certificates");
		httpClient.execute(httpPost).close();

		Thread.sleep(15000L);

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), startsWith("https://edge-cache-09"));
			assertThat(header.getValue(), endsWith("https-nocert.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78"));
		}
	}

	// This is a workaround to get HttpClient to do the equivalent of
	// curl -v --resolve 'tr.https-only-test.thecdn.cdnlab.example.com:8443:127.0.0.1' https://tr.https-only-test.thecdn.example.com:8443/foo.json
	class ClientSslSocketFactory extends SSLConnectionSocketFactory {
		private final String host;

		public ClientSslSocketFactory(String host) throws Exception {
			super(SSLContextBuilder.create().loadTrustMaterial(null, new TrustSelfSignedStrategy()).build(),
				new TestHostnameVerifier());
			this.host = host;
		}

		protected void prepareSocket(final SSLSocket sslSocket) throws IOException {
			SNIHostName serverName = new SNIHostName(host);
			List<SNIServerName> serverNames = new ArrayList<>(1);
			serverNames.add(serverName);

			SSLParameters params = sslSocket.getSSLParameters();
			params.setServerNames(serverNames);
			sslSocket.setSSLParameters(params);
		}
	}

	// This is a workaround for the same reason as above
	// org.apache.http.conn.ssl.SSLConnectionSocketFactory.verifyHostname(<socket>, 'localhost') normally fails
	class TestHostnameVerifier implements HostnameVerifier {
		@Override
		public boolean verify(String s, SSLSession sslSession) {
			return true;
		}
	}
}
