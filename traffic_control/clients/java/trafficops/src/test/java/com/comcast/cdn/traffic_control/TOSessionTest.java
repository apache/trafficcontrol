package com.comcast.cdn.traffic_control;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.apache.http.HttpResponse;
import org.apache.http.HttpVersion;
import org.apache.http.client.methods.RequestBuilder;
import org.apache.http.entity.StringEntity;
import org.apache.http.message.BasicStatusLine;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.comcast.cdn.traffic_control.exception.LoginException;
import com.comcast.cdn.traffic_control.models.Response.CollectionResponse;

public class TOSessionTest {
	private static final Logger LOG = LoggerFactory.getLogger(TOSessionTest.class);
	
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
		TOSession.builder().setRestClient(sessionMock).build();
	}
	
	@Test(expected=LoginException.class)
	public void test401Response() throws Throwable {
		HttpResponse resp = Mockito.mock(HttpResponse.class);
		Mockito.when(resp.getStatusLine()).thenReturn(new BasicStatusLine(HttpVersion.HTTP_1_0, 401, "Not Auth"));
		
		CompletableFuture<HttpResponse> f = new CompletableFuture<>();
		f.complete(resp);
		
		Mockito.doReturn(f).when(sessionMock).execute(Mockito.any(RequestBuilder.class));
		
		TOSession session = TOSession.builder().setRestClient(sessionMock).build();
		
		try {
			session.getDeliveryServices().get();
		} catch(Throwable e) {
			throw e.getCause();
		}
	}
	
	@Test
	public void deliveryServices() throws Throwable {
		HttpResponse resp = Mockito.mock(HttpResponse.class);
		Mockito.doReturn(new BasicStatusLine(HttpVersion.HTTP_1_0, 200, "Ok")).when(resp).getStatusLine();
		Mockito.doReturn(new StringEntity(DeliveryService_Good_Response)).when(resp).getEntity();
		
		CompletableFuture<HttpResponse> f = new CompletableFuture<>();
		f.complete(resp);
		
		Mockito.doReturn(f).when(sessionMock).execute(Mockito.any(RequestBuilder.class));
		
		TOSession session = TOSession.builder().setRestClient(sessionMock).build();
		CollectionResponse cResp = session.getDeliveryServices().get();
		
		assertNotNull(cResp);
		assertNotNull(cResp.getResponse());
		assertEquals(1, cResp.getResponse().size());
		
		final Map<String,Object> service = cResp.getResponse().get(0);
		assertEquals("us-co-denver", service.get("cachegroup"));
		LOG.debug("Service: {}", service);
	}
}
