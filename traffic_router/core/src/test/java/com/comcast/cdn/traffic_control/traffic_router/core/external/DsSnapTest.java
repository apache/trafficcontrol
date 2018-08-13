/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a shallowCopy of the License at
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
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpHead;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustSelfSignedStrategy;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.ssl.SSLContextBuilder;
import org.hamcrest.Matchers;
import org.junit.After;
import org.junit.Before;
import org.junit.FixMethodOrder;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.junit.runners.MethodSorters;
import org.xbill.DNS.DClass;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.Options;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Resolver;
import org.xbill.DNS.Section;
import org.xbill.DNS.SimpleResolver;
import org.xbill.DNS.Type;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SNIHostName;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLHandshakeException;
import javax.net.ssl.SSLParameters;
import javax.net.ssl.SSLSession;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.TrustManagerFactory;
import java.io.IOException;
import java.io.InputStream;
import java.security.KeyStore;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Set;

import static org.hamcrest.Matchers.*;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class DsSnapTest {
	private CloseableHttpClient httpClient;
	private final String cdnDomain = ".thecdn.example.com";

	private String deliveryServiceId;
	private String deliveryServiceDomain;
	private final List<String> validLocations = new ArrayList<>();

	private final String httpsOnlyId = "https-only-test";
	private final String httpsOnlyDomain = httpsOnlyId + cdnDomain;
	private final List<String> httpsOnlyLocations = new ArrayList<>();

	private final String additionalId = "https-additional";
	private final String additionalDomain = additionalId + cdnDomain;
	private final List<String> httpsAdditionalLocations = new ArrayList<>();

	private final String httpsNoCertsId = "https-nocert";
	private final String httpsNoCertsDomain = httpsNoCertsId + cdnDomain;
	private final List<String> httpsNoCertsLocations = new ArrayList<>();

	private final String httpAndHttpsId = "http-and-https-test";
	private final String httpAndHttpsDomain = httpAndHttpsId + cdnDomain;
	private final List<String> httpAndHttpsLocations = new ArrayList<>();

	private final String httpToHttpsId = "http-to-https-test";
	private final String httpToHttpsDomain = httpToHttpsId + cdnDomain;
	private final List<String> httpToHttpsLocations = new ArrayList<>();

	private final String httpOnlyId = "http-only-test";
	private final String httpOnlyDomain = httpOnlyId + cdnDomain;
	private final List<String> httpOnlyLocations = new ArrayList<>();

	private final String routerHttpPort = System.getProperty("routerHttpPort", "8888");
	private final String routerSecurePort = System.getProperty("routerSecurePort", "8443");
	private final String testHttpPort = System.getProperty("testHttpServerPort", "8889");
	private KeyStore trustStore;
	private static boolean doneBefore = false;

	@Before
	public void before() throws Exception {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

		String resourcePath = "internal/api/1.3/steering.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		Set<String> steeringDeliveryServices = new HashSet<String>();
		JsonNode steeringData = objectMapper.readTree(inputStream).get("response");
		Iterator<JsonNode> elements = steeringData.elements();

		while (elements.hasNext()) {
			JsonNode ds = elements.next();
			String dsId = ds.get("deliveryService").asText();
			steeringDeliveryServices.add(dsId);
		}

		resourcePath = "publish/DsSnap/CrConfig.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = objectMapper.readTree(inputStream);

		deliveryServiceId = null;

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext()) {
			String dsId = deliveryServices.next();

			if (steeringDeliveryServices.contains(dsId)) {
				continue;
			}

			JsonNode deliveryServiceNode = jsonNode.get("deliveryServices").get(dsId);
			Iterator<JsonNode> matchsets = deliveryServiceNode.get("matchsets").iterator();

			while (matchsets.hasNext() && deliveryServiceId == null) {
				if ("HTTP".equals(matchsets.next().get("protocol").asText())) {
					final boolean sslEnabled = JsonUtils.optBoolean(deliveryServiceNode, "sslEnabled");
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
				validLocations.add("http://" + cacheId + "." + deliveryServiceDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78&format=json");
			}

			if (cacheNode.get("deliveryServices").has(httpsOnlyId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpsOnlyLocations.add("https://" + cacheId + "." + httpsOnlyDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}
			if (cacheNode.get("deliveryServices").has(additionalId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpsAdditionalLocations.add("https://" + cacheId + "." + additionalDomain + portText + "/stuff?fakeClientIpAddress" +
						"=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(httpsNoCertsId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpsNoCertsLocations.add("https://" + cacheId + "." + httpsNoCertsDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(httpAndHttpsId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpAndHttpsLocations.add("https://" + cacheId + "." + httpAndHttpsDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");

				port = cacheNode.has("port") ? cacheNode.get("port").asInt(80) : 80;
				portText = (port == 80) ? "" : ":" + port;
				httpAndHttpsLocations.add("http://" + cacheId + "." + httpAndHttpsDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(httpToHttpsId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				httpToHttpsLocations.add("https://" + cacheId + "." + httpToHttpsDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(httpOnlyId)) {
				int port = cacheNode.has("port") ? cacheNode.get("port").asInt(80) : 80;

				String portText = (port == 80) ? "" : ":" + port;
				httpOnlyLocations.add("http://" + cacheId + "." + httpOnlyDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}
		}

		assertThat(validLocations.isEmpty(), equalTo(false));
		assertThat(httpsOnlyLocations.isEmpty(), equalTo(false));

		trustStore = KeyStore.getInstance(KeyStore.getDefaultType());
		InputStream keystoreStream = getClass().getClassLoader().getResourceAsStream("keystore.jks");
		trustStore.load(keystoreStream, "changeit".toCharArray());
		TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm()).init(trustStore);

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.https-only-test.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();

		// Pretend someone did a cr-config snapshot with a DsSnap crconfig
		if (!doneBefore) {
			doneBefore = true;
			HttpPost httpPost = new HttpPost("http://localhost:" + testHttpPort + "/crconfig-dssnap");
			httpClient.execute(httpPost).close();

			// Default interval for polling cr config is 10 seconds
			Thread.sleep(55 * 1000);
		}
	}

	@After
	public void after() throws IOException {
		if (httpClient != null) {
			httpClient.close();
		}
	}

	@Test
	public void itRedirectsValidHttpRequests() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "trvr." + deliveryServiceId + ".bar");
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
		httpGet.addHeader("Host", "fooswc." + deliveryServiceId + ".bar");
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
		httpGet.addHeader("Host", "trcvr." + deliveryServiceId + ".bar");
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

	@Test
	public void itRedirectsFromHttpToHttps() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpToHttpsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpToHttpsLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(httpToHttpsId + ".thecdn.example.com"));
			assertThat(header.getValue(), containsString("/stuff"));
		}

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.http-and-https-test.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();

		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpToHttpsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpToHttpsLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(httpToHttpsId + ".thecdn.example.com"));
			assertThat(header.getValue(), containsString("/stuff"));
		}
	}

	@Test
	public void itRejectsHttpsRequestsForHttpDeliveryService() throws Exception {
		HttpGet httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat("Response 503 expected got"+response.getStatusLine().getStatusCode(),response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itPreservesProtocolForHttpAndHttps() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpAndHttpsId + ".bar");
		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpAndHttpsLocations));
			assertThat(header.getValue(), startsWith("http://"));
			assertThat(header.getValue(), containsString(httpAndHttpsId + ".thecdn.example.com"));
			assertThat(header.getValue(), containsString("/stuff"));
		}

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.http-and-https-test.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();

		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + httpAndHttpsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(httpAndHttpsLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(httpAndHttpsId + ".thecdn.example.com"));
			assertThat(header.getValue(), containsString("/stuff"));
		}
	}

	@Test
	public void digDefaultCrConfig() throws Exception {
		Message response = lookupTest("edge.dns-test.thecdn.example.com",  Rcode.NOERROR);
		final String expectedIps[] = {"12.34.0.101","12.34.0.102"};
		recordTest(response, expectedIps, Section.ANSWER);
		final String expectedTrs[] = {"testing-tr-01.thecdn.example.com.","testing-tr-02.thecdn.example.com."};
		recordTest(response, expectedTrs, Section.AUTHORITY);
		lookupTest("edge.https-dns-test.thecdn.example.com",  Rcode.NXDOMAIN);
		lookupTest("edge.https-only-test.thecdn.example.com",  Rcode.NXDOMAIN);
	}


	private Message lookupTest(final String hostname, final int expectedResult) throws Exception {
		Options.set("verbose");
		Resolver resolver = new SimpleResolver("127.0.0.1");
		resolver.setPort(Integer.parseInt(System.getProperty("dns.udp.port")));
		Name name = Name.fromString(hostname, Name.root);
		Record rec = Record.newRecord(name, Type.A, DClass.IN);
		Message query = Message.newQuery(rec);
		Message response = resolver.send(query);
		assertThat("DNS Result '"+Rcode.string(expectedResult)+"' expected response but got "
						+response.getRcode(),response.getRcode(), equalTo(expectedResult));
		return response;
	}

	private void recordTest(final Message response, final String[] expectedRecords, final int sectionDex){
		final int rCnt = response.getSectionRRsets(sectionDex)[0].size();
		final int expectedCnt = expectedRecords.length;
		assertThat("Expected "+expectedCnt+" records but there were :"+rCnt,rCnt,equalTo(expectedCnt));
		response.getSectionRRsets(sectionDex)[0].rrs().forEachRemaining(rr-> {
			String rs = ((Record)rr).rdataToString();
			assertThat("Expected the record to be one of "+expectedRecords+" but IP was "+rs,
					rs, isIn(expectedRecords));
		});
	}

	@Test
	public void itRejectsCrConfigWithMissingCert() throws Exception {
		/* Test summary
		- Already on default crconfig from initialization
		- Do resolve httpOnly and get 302
		- Do resolve httpsNoCert and get 500+
		- Switch to crconfig2
		- crconfig2 does not get used because its waiting for SSL cert for httpsNoCert
		- Do resolve httpOnly and get 302
		- Do resolve httpsNoCert and get 500+ because still using crconfig 1
		- Switch to crconfig3
		- Do resolve httpOnly and get 302
		- Do resolve httpsNoCert and get 500+ because ssl=false
		- Switch to crconfig4
		- Switch to sslkeys-missing-1.json (get /certificates)
		- Do resolve https-additional and get 302
		- Do resolve httpsNoCerts and get 302
		- Do resolve httpOnly and get 302
		 */
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr1." + httpOnlyId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			assertThat(response.getFirstHeader("Location").getValue(), isOneOf(
				"http://edge-cache-000.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
				"http://edge-cache-001.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
				"http://edge-cache-002.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78"
			));
		}

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory(httpsNoCertsDomain))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();

		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/x?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr1." + httpsNoCertsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)){
			int code = response.getStatusLine().getStatusCode();
			assertThat("Expected to get an ssl handshake error! But got: "+code,
					code, greaterThan(500));
		} catch (SSLHandshakeException she) {
			// expected result
		}


		// Pretend someone did a cr-config snapshot that would have updated the location to be different
		HttpPost httpPost = new HttpPost("http://localhost:" + testHttpPort + "/crconfig-2");
		httpClient.execute(httpPost).close();

		// Default interval for polling cr config is 10 seconds
		Thread.sleep(15 * 1000);

		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr2." + httpOnlyId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, not(isOneOf(
					"http://edge-cache-010.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-011.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-012.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78"
			)));
		}

		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr2." + httpsNoCertsId + ".bar");

		// verify we do not yet use the new configuration
		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
			Header rdtlH = response.getFirstHeader("rdtl");
			if (rdtlH!=null) {
				assertThat(rdtlH.getValue(), is("DS_NOT_FOUND" ));
			}
		}

		// verify that if we get a new cr-config that turns off https for the problematic delivery service
		// that it's able to get through while TR is still concurrently trying to get certs

		String testHttpPort = System.getProperty("testHttpServerPort", "8889");
		httpPost = new HttpPost("http://localhost:"+ testHttpPort + "/crconfig-3");
		httpClient.execute(httpPost).close();

		// Default interval for polling cr config is 10 seconds
		Thread.sleep(30 * 1000);

		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr3." + httpOnlyId + ".bar");

		// verify we now use the new configuration
		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, isOneOf(
			"http://edge-cache-010.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					// Default interval for polling cr config is 10 seconds
				"http://edge-cache-011.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
				"http://edge-cache-012.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78"
			));
		}

		// assert that request gets rejected because SSL is turned off
		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr3." + httpsNoCertsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)){
			int code = response.getStatusLine().getStatusCode();
			assertThat("Expected to get an server error! But got: "+code,
					code, greaterThan(500));
		}
		catch (javax.net.ssl.SSLHandshakeException she)
		{
			// expected result
		}

		// Go back to the cr-config that makes the delivery service https again
		// Pretend someone did a cr-config snapshot that would have updated the location to be different
		httpPost = new HttpPost("http://localhost:" + testHttpPort + "/crconfig-4");
		httpClient.execute(httpPost).close();

		Thread.sleep(15 * 1000);

		// Update certificates so new ds is valid
		testHttpPort = System.getProperty("testHttpServerPort", "8889");
		httpPost = new HttpPost("http://localhost:"+ testHttpPort + "/certificates");
		httpClient.execute(httpPost).close();

		// Our initial test cr config data sets cert poller to 10 seconds
		Thread.sleep(25000L);

		httpClient = HttpClientBuilder.create()
				.setSSLSocketFactory(new DsSnapTest.ClientSslSocketFactory(additionalDomain))
				.setSSLHostnameVerifier(new DsSnapTest.TestHostnameVerifier())
				.disableRedirectHandling()
				.build();
		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr."+additionalId+".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, isOneOf(
			"https://edge-cache-011.https-additional.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78",
					"https://edge-cache-012.https-additional.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78"
			));
	    } catch (SSLHandshakeException e) {
			fail("Expected a 302 but got error: "+e.getMessage());
	    }

		httpClient = HttpClientBuilder.create()
				.setSSLSocketFactory(new ClientSslSocketFactory(httpsNoCertsDomain))
				.setSSLHostnameVerifier(new TestHostnameVerifier())
				.disableRedirectHandling()
				.build();
		httpGet = new HttpGet("https://localhost:" + routerSecurePort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr4." + httpsNoCertsId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, isOneOf(
				"https://edge-cache-090.https-nocert.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78",
				"https://edge-cache-091.https-nocert.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78",
				"https://edge-cache-092.https-nocert.thecdn.example.com/stuff?fakeClientIpAddress=12.34.56.78"
			));
		} catch (SSLHandshakeException e) {
			// TODO should not come here - fix
			fail(e.getMessage());
		}

		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr4." + httpOnlyId + ".bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, isOneOf(
				"http://edge-cache-010.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
				"http://edge-cache-011.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
				"http://edge-cache-012.http-only-test.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78"
			));
		}

		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr4.steering-target-three.bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();
			assertThat(location, isOneOf(
					"http://edge-cache-020.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-021.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-022.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-030.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-031.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-032.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-040.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-041.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78",
					"http://edge-cache-042.steering-target-3.thecdn.example.com:8090/stuff?fakeClientIpAddress=12.34.56.78"
			));
		}
		httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr4.steering-target-3.bar");

		try (CloseableHttpResponse response = httpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
		}
	}

	@Test
	public void itDoesUseLocationFormatResponse() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78&format=json");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));

			HttpEntity entity = response.getEntity();
			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

			assertThat(entity.getContent(), not(nullValue()));

			JsonNode json = objectMapper.readTree(entity.getContent());

			assertThat(json.has("location"), equalTo(true));
			assertThat(json.get("location").asText(), isIn(validLocations));
			assertThat(json.get("location").asText(), Matchers.startsWith("http://"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itDoesNotUseLocationFormatResponseForHead() throws IOException, InterruptedException {
		HttpHead httpHead = new HttpHead("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78&format=json");
		httpHead.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpHead);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));
			assertThat("Failed getting null body for HEAD request", response.getEntity(), nullValue());
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	// after itRejectsCrConfigWithMissingCert
	public void zdigCrConfig4() throws Exception {
		HttpPost httpPost = new HttpPost("http://localhost:" + testHttpPort + "/crconfig-4");
		httpClient.execute(httpPost).close();

		Thread.sleep(15 * 1000);

		// Update certificates so new ds is valid
		httpPost = new HttpPost("http://localhost:"+ testHttpPort + "/certificates");
		httpClient.execute(httpPost).close();

		// Our initial test cr config data sets cert poller to 10 seconds
		Thread.sleep(25000L);

		Message response =
				lookupTest("edge.dns-test.thecdn.example.com",  Rcode.NOERROR);
		final String expectedIps[] = {"12.34.0.100","12.34.0.101","12.34.0.102"};
		recordTest(response, expectedIps, Section.ANSWER);
		final String expectedTrs[] = {"testing-tr-01.thecdn.example.com.","testing-tr-02.thecdn.example.com."};
		recordTest(response, expectedTrs, Section.AUTHORITY);
		response = lookupTest("edge.https-dns-test.thecdn.example.com",  Rcode.NOERROR);
		final String expectedIPs[] = {"12.34.7.100","12.34.7.101","12.34.7.102"};
		recordTest(response, expectedIPs, Section.ANSWER);
		lookupTest("edge.https-only-test.thecdn.example.com",  Rcode.NXDOMAIN);
	}

	// This is a workaround to get HttpClient to do the equivalent of
	// curl -v --resolve 'tr.https-only-test.thecdn.cdnlab.example.com:8443:127.0.0.1' https://tr.https-only-test.thecdn.example.com:8443/foo.json
	class ClientSslSocketFactory extends SSLConnectionSocketFactory {
		private final String host;
		public ClientSslSocketFactory(String host) throws Exception {
			super(SSLContextBuilder.create().loadTrustMaterial(trustStore, new TrustSelfSignedStrategy()).build(),
					new DsSnapTest.TestHostnameVerifier());
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
			assertThat("s = "+s+", getPeerHost() = "+ sslSession.getPeerHost(), sslSession.getPeerHost(), equalTo(s));
			return sslSession.getPeerHost().equals(s);
		}
	}
}
