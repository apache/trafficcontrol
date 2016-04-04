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

import java.io.IOException;

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsEqual.equalTo;

@Category(ExternalTest.class)
public class RouterTest {
	private CloseableHttpClient httpClient;

	@Before
	public void before() throws IOException, InterruptedException, LifecycleException {
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

		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=113.203.235.227");
		httpGet.addHeader("Host", "foo.omg-04.bar");
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
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=113.203.235.227");
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
