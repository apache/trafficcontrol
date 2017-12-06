package com.comcast.cdn.traffic_control;

import java.io.Closeable;
import java.io.IOException;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.CookieStore;
import org.apache.http.client.config.CookieSpecs;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.HttpUriRequest;
import org.apache.http.client.methods.RequestBuilder;
import org.apache.http.client.protocol.HttpClientContext;
import org.apache.http.concurrent.FutureCallback;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.BasicCookieStore;
import org.apache.http.impl.nio.client.CloseableHttpAsyncClient;
import org.apache.http.impl.nio.client.HttpAsyncClients;
import org.apache.http.message.BasicHeader;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.comcast.cdn.traffic_control.exception.OperationException;
import com.google.auto.value.AutoValue;
import com.google.common.base.Charsets;
import com.google.common.collect.ImmutableList;

@AutoValue
public abstract class RestApiSession implements Closeable {
	private static final Logger LOG = LoggerFactory.getLogger(RestApiSession.class);

	private static final String URL_FORMAT_STR = "%s://%s:%s/%s/%s/%s";
	
	public static final String DEFAULT_API_PATH = "api";
	public static final String DEFAULT_API_VERSION = "1.2";
	public static final ImmutableList<Header> DEFAULT_HEADERS;
	static {
		DEFAULT_HEADERS = ImmutableList.<Header>builder()
				.add(new BasicHeader("Content-Type", "application/json; charset=UTF-8")).build();
	};

	private CloseableHttpAsyncClient httpclient;

	public void open() {
		if (httpclient == null) {
			RequestConfig globalConfig = RequestConfig.custom()
					.setCookieSpec(CookieSpecs.STANDARD) //User standard instead of default. Default will result in cookie parse exceptions with the Mojolicous cookie
					.setConnectTimeout(5000)
					.build();
			CookieStore cookieStore = new BasicCookieStore();
			HttpClientContext context = HttpClientContext.create();
			context.setCookieStore(cookieStore);

			httpclient = HttpAsyncClients.custom()
					.setDefaultRequestConfig(globalConfig)
					.setDefaultCookieStore(cookieStore)
					.build();
		}

		if (!httpclient.isRunning()) {
			httpclient.start();
		}
	}

	public boolean isRunning() {
		if (httpclient == null) {
			return false;
		} else {
			return httpclient.isRunning();
		}
	}

	public void close() throws IOException {
		if (httpclient != null) {
			httpclient.close();
			httpclient = null;
		}
	}

	public String buildUrl(String path) {
		return String.format(URL_FORMAT_STR, this.ssl() ?"https":"http", this.host(), this.port(), this.apiBasePath(),
				this.apiVersion(), path);
	}
	
	public CompletableFuture<HttpResponse> get(String url) {
		return execute(RequestBuilder.get(url));
	}
	
	public CompletableFuture<HttpResponse> post(String url, String body) {
		final HttpEntity e = new StringEntity(body, Charsets.UTF_8);
		return execute(RequestBuilder.post()
				.setUri(url)
				.setEntity(e));
	}
	
	public CompletableFuture<HttpResponse> execute(RequestBuilder request) {
		for(Header h: this.defaultHeaders()) {
			request.addHeader(h);
		}
		
		return this.execute(request.build());
	}
	
	private CompletableFuture<HttpResponse> execute(HttpUriRequest request) {
		final CompletableFutureCallback future = new CompletableFutureCallback();
		try {
			LOG.debug("Opening RestClient");
			this.open();
			
			LOG.debug("Dispatching request: {}", request);
			final Future<HttpResponse> reFuture = this.httpclient.execute(request, future);
			future.setReFuture(reFuture);
		} catch(Throwable e) {
			future.completeExceptionally(e);
		}
		
		return future;
	}
	
	private class CompletableFutureCallback extends CompletableFuture<HttpResponse> implements FutureCallback<HttpResponse>{
		private Future<HttpResponse> reFuture;
		
		@Override
		public boolean cancel(boolean mayInterruptIfRunning) {
			if(reFuture != null) {
				return reFuture.cancel(mayInterruptIfRunning);
			}
			return false;
		}
		
		@Override
		public void completed(HttpResponse result) {
			LOG.debug("Request Completed: {}", result);
			this.complete(result);
		}

		@Override
		public void failed(Exception ex) {
			LOG.debug("Request Failed", ex);
			this.completeExceptionally(ex);
		}

		@Override
		public void cancelled() {
			LOG.debug("Request Cancelled");
			this.completeExceptionally(new OperationException("HTTP Request was cancelled"));
		}

		public void setReFuture(Future<HttpResponse> reFuture) {
			this.reFuture = reFuture;
		}
	}

	public abstract String host();

	public abstract int port();

	public abstract String apiVersion();

	public abstract String apiBasePath();

	public abstract ImmutableList<Header> defaultHeaders();

	public abstract boolean ssl();

	static Builder builder() {
		return new AutoValue_RestApiSession.Builder()
				.setApiBasePath(DEFAULT_API_PATH)
				.setApiVersion(DEFAULT_API_VERSION)
				.setDefaultHeaders(DEFAULT_HEADERS);
	}
	
	public abstract Builder toBuilder();

	@AutoValue.Builder
	public abstract static class Builder {
		public abstract RestApiSession build();

		public abstract Builder setHost(String host);

		public abstract Builder setPort(int port);

		public abstract Builder setApiVersion(String version);

		public abstract Builder setApiBasePath(String version);

		public abstract Builder setSsl(boolean ssl);

		public abstract Builder setDefaultHeaders(ImmutableList<Header> headers);
		public abstract ImmutableList.Builder<Header> defaultHeadersBuilder();

	}
}
