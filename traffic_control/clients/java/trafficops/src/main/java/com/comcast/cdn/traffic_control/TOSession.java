package com.comcast.cdn.traffic_control;

import java.io.Closeable;
import java.io.IOException;
import java.net.URI;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.apache.http.NameValuePair;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.message.BasicNameValuePair;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.comcast.cdn.traffic_control.exception.LoginException;
import com.comcast.cdn.traffic_control.models.Response;
import com.google.auto.value.AutoValue;
import com.google.common.collect.ImmutableMap;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

@AutoValue
public abstract class TOSession implements Closeable {
	private static final Logger LOG = LoggerFactory.getLogger(TOSession.class);
	
	private static final String URL_FORMAT_STR = "/%s/%s/%s";
	
	public static final String DEFAULT_API_PATH = "api";
	public static final String DEFAULT_API_VERSION = "1.2";
	
	public abstract RestApiSession restClient();
	public abstract String host();
	public abstract int port();
	public abstract boolean ssl();
	public abstract String apiVersion();
	public abstract String apiBasePath();
	
	static final Gson gson = new GsonBuilder()
			.create();
	
	private boolean isLoggedIn = false;
	
	protected URIBuilder newUriBuilder(final String path) {
		final String _path = String.format(URL_FORMAT_STR, this.apiBasePath(), this.apiVersion(), path);
		
		return new URIBuilder()
				.setScheme(this.ssl() ? "https" : "http")
				.setHost(this.host())
				.setPort(this.port())
				.setPath(_path);
	}
	protected List<NameValuePair> toPairs(Map<String, ?> params){
		if(params != null && !params.isEmpty()) {
			List<NameValuePair> pairs = new LinkedList<>();
			for(Map.Entry<String, ?> param: params.entrySet()) {
				pairs.add(new BasicNameValuePair(param.getKey(), param.getValue().toString()));
			}
		}
		return Collections.emptyList();
	}
	
	public void close() throws IOException {
		this.restClient().close();
	}
	public boolean isLoggedIn() {
		return isLoggedIn;
	}
	
	public CompletableFuture<Boolean> login(final String username, final String password) {
		URI uri;
		try {
			uri = this.newUriBuilder("user/login.json")
					.build();
		} catch (Throwable e) {
			final CompletableFuture<Boolean> f = new CompletableFuture<>();
			f.completeExceptionally(e);
			return f;
		}
		
		LOG.debug("Logging into: {}", uri);
		return ResponseFuture.builder()
			.setHandleException((f,t)-> {
				f.completeExceptionally(new LoginException(String.format("Failed to login with username %s", username), t));
			})
			.setMethod(ResponseFuture.Method.POST)
			.setUri(uri)
			.setSession(this.restClient())
			.setBody(gson.toJson(ImmutableMap.<String,String>of("u", username, "p", password))).build()
			.thenApply(r->{
				isLoggedIn = true;
				return true;
			});
	}
	public CompletableFuture<Response.CollectionResponse> getServers(){
		return this.getServers(null);
	}
	public CompletableFuture<Response.CollectionResponse> getServers(final Map<String, ?> filterParams){
		URI uri;
		try {
			uri = this.newUriBuilder("servers.json")
					.setParameters(this.toPairs(filterParams))
					.build();
		} catch (Throwable e) {
			final CompletableFuture<Response.CollectionResponse> f = new CompletableFuture<>();
			f.completeExceptionally(e);
			return f;
		}
		return ResponseFuture.builder(Response.CollectionResponse.class)
				.setMethod(ResponseFuture.Method.GET)
				.setUri(uri)
				.setSession(this.restClient())
				.build();
	}
	
	public CompletableFuture<Response.CollectionResponse> getDeliveryServices(){
		return this.getDeliveryServices(null);
	}
	public CompletableFuture<Response.CollectionResponse> getDeliveryServices(final Map<String, ?> filterParams){
		URI uri;
		try {
			uri = this.newUriBuilder("deliveryservices.json")
					.setParameters(this.toPairs(filterParams))
					.build();
		} catch (Throwable e) {
			final CompletableFuture<Response.CollectionResponse> f = new CompletableFuture<>();
			f.completeExceptionally(e);
			return f;
		}
		LOG.debug("getDeliveryService url {}", uri);
		return ResponseFuture.builder(Response.CollectionResponse.class)
				.setMethod(ResponseFuture.Method.GET)
				.setUri(uri)
				.setSession(this.restClient())
				.build();
	}
	
	
	public static Builder builder() {
		return new AutoValue_TOSession.Builder()
				.setApiBasePath(DEFAULT_API_PATH)
				.setApiVersion(DEFAULT_API_VERSION);
	}
	public abstract Builder toBuilder();
	
	@AutoValue.Builder
	public abstract static class Builder {
		public abstract TOSession build();
		
		public abstract Builder setRestClient(RestApiSession restClient);
		public abstract RestApiSession.Builder restClientBuilder();

		public abstract Builder setHost(String host);
		public abstract Builder setPort(int port);
		public abstract Builder setApiVersion(String version);
		public abstract Builder setApiBasePath(String version);
		public abstract Builder setSsl(boolean ssl);
		
		public Builder fromURI(URI uri){
			return this.setSsl(uri.getScheme().equals("http") ? false: true)
					.setHost(uri.getHost())
					.setPort(uri.getPort());
		}
	}
}
