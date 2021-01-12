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

package org.apache.trafficcontrol.client.trafficops;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.net.URI;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.apache.http.HttpResponse;
import org.apache.http.HttpVersion;
import org.apache.http.client.methods.RequestBuilder;
import org.apache.http.entity.StringEntity;
import org.apache.http.message.BasicStatusLine;
import org.apache.trafficcontrol.client.RestApiSession;
import org.apache.trafficcontrol.client.exception.LoginException;
import org.apache.trafficcontrol.client.trafficops.models.Response.CollectionResponse;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class TOSessionTest {
	private static final Logger LOG = LoggerFactory.getLogger(TOSessionTest.class);
	
	public static final URI baseUri = URI.create("http://trafficcontrol.apache.org:443");
	
	public static final String DeliveryService_Good_Response = "{\"response\": [{\"cachegroup\": \"us-co-denver\"}]}";
	
	private RestApiSession sessionMock;
	
	@Before
	public void before() {
		sessionMock = Mockito.mock(RestApiSession.class, Mockito.CALLS_REAL_METHODS);
	}
	@After
	public void after() {
		sessionMock=null;
	}

	@Test
	public void testBuild() {
		TOSession.builder()
			.setRestClient(sessionMock)
			.fromURI(baseUri)
			.build();
	}
	
	@Test(expected=LoginException.class)
	public void test401Response() throws Throwable {
		HttpResponse resp = Mockito.mock(HttpResponse.class);
		Mockito
			.when(resp.getStatusLine())
			.thenReturn(new BasicStatusLine(HttpVersion.HTTP_1_0, 401, "Not Auth"));
		
		final CompletableFuture<HttpResponse> f = new CompletableFuture<>();
		f.complete(resp);
		
		Mockito
			.doReturn(f)
			.when(sessionMock)
			.execute(Mockito.any(RequestBuilder.class));
		
		TOSession session = TOSession
				.builder()
				.fromURI(baseUri)
				.setRestClient(sessionMock)
				.build();
		
		try {
			session.getDeliveryServices().get();
		} catch(Throwable e) {
			throw e.getCause();
		}
	}
	
	@Test
	public void deliveryServices() throws Throwable {
		final HttpResponse resp = Mockito.mock(HttpResponse.class);
		Mockito
			.doReturn(new BasicStatusLine(HttpVersion.HTTP_1_0, 200, "Ok"))
			.when(resp)
			.getStatusLine();
		Mockito
			.doReturn(new StringEntity(DeliveryService_Good_Response))
			.when(resp)
			.getEntity();
		
		final CompletableFuture<HttpResponse> f = new CompletableFuture<>();
		f.complete(resp);
		
		Mockito
			.doReturn(f)
			.when(sessionMock)
			.execute(Mockito.any(RequestBuilder.class));
		
		final TOSession session = TOSession.builder()
				.fromURI(baseUri)
				.setRestClient(sessionMock)
				.build();
		
		CollectionResponse cResp = session.getDeliveryServices().get();
		
		assertNotNull(cResp);
		assertNotNull(cResp.getResponse());
		assertEquals(1, cResp.getResponse().size());
		
		final Map<String,?> service = cResp.getResponse().get(0);
		assertEquals("us-co-denver", service.get("cachegroup"));
		LOG.debug("Service: {}", service);
	}
}
