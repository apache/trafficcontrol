package com.comcast.cdn.traffic_control.traffic_router.neustar.data;

import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpUriRequest;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.log4j.Logger;

import java.io.IOException;

public class HttpClient {
	private final Logger LOGGER = Logger.getLogger(HttpClient.class);

	private CloseableHttpClient httpClient;

	public CloseableHttpResponse execute(HttpUriRequest request) {
		try {
			httpClient = HttpClientBuilder.create().build();
			return httpClient.execute(request);
		} catch (IOException e) {
			LOGGER.warn("Failed to execute http request " + request.getMethod() + " " + request.getURI() + ": " + e.getMessage());
			try {
				httpClient.close();
			} catch (IOException e1) {
				LOGGER.warn("After exception, Failed to close Http Client " + e1.getMessage());
			}
			return null;
		}
	}

	public void close() {
		try {
			httpClient.close();
		} catch (IOException e) {
			LOGGER.warn("Failed to close Http Client " + e.getMessage());
		}
	}
}
