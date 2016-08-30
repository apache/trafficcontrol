package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
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

import java.net.URLEncoder;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsEqual.equalTo;

@Category(ExternalTest.class)
public class DeliveryServicesTest {
	private CloseableHttpClient closeableHttpClient;

	@Before
	public void before() throws LifecycleException {
		closeableHttpClient = HttpClientBuilder.create().build();
	}

	@After
	public void after() throws Exception{
		if (closeableHttpClient != null) closeableHttpClient.close();
	}

	@Test
	public void itReturnsIdOfValidDeliveryService() throws Exception {
		String encodedUrl = URLEncoder.encode("http://trafficrouter01.steering-target-1.thecdn.example.com/stuff", "utf-8");
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/deliveryservices?url="+encodedUrl);

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			String responseBody = EntityUtils.toString(response.getEntity());
			assertThat(responseBody, equalTo("{\"id\":\"steering-target-1\"}"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itReturnsNotFoundForNonexistentDeliveryService() throws Exception {
		String encodedUrl = URLEncoder.encode("http://trafficrouter01.somedeliveryservice.somecdn.domain.foo/stuff", "utf-8");
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/deliveryservices?url="+encodedUrl);

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(404));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itReturnsBadRequestForBadUrlQueryParameter() throws Exception {
		String encodedUrl = "httptrafficrouter01somedeliveryservicesomecdndomainfoo/stuff";
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/deliveryservices?url="+encodedUrl);

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(400));
		} finally {
			if (response != null) response.close();
		}
	}
}
