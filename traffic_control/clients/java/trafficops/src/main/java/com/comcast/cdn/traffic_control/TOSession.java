package com.comcast.cdn.traffic_control;

import java.io.Closeable;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.function.BiConsumer;

import org.apache.http.HttpResponse;
import org.apache.http.client.methods.RequestBuilder;
import org.apache.http.entity.StringEntity;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.comcast.cdn.traffic_control.exception.LoginException;
import com.comcast.cdn.traffic_control.exception.OperationException;
import com.comcast.cdn.traffic_control.models.Alert;
import com.comcast.cdn.traffic_control.models.Response;
import com.google.auto.value.AutoValue;
import com.google.common.base.Charsets;
import com.google.common.collect.ImmutableMap;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

@AutoValue
public abstract class TOSession implements Closeable {
	private static final Logger LOG = LoggerFactory.getLogger(TOSession.class);
	
	private static final Gson gson = new GsonBuilder()
			.create();
	
	private boolean isLoggedIn = false;
	
	public void close() throws IOException {
		this.restClient().close();
	}
	
	public String toUrl() {
		return this.restClient().buildUrl("/");
	}
	public boolean isLoggedIn() {
		return isLoggedIn;
	}
	
	public CompletableFuture<Boolean> login(final String username, final String password) {
		final String url = this.restClient().buildUrl("user/login.json");
		LOG.debug("Logging into: {}", url);
		return ResponseFuture.builder()
			.setHandleException((f,t)-> {
				f.completeExceptionally(new LoginException(String.format("Failed to login with username %s", username), t));
			})
			.setMethod(ResponseFuture.Method.POST)
			.setUrl(url)
			.setSession(this.restClient())
			.setBody(gson.toJson(ImmutableMap.<String,String>of("u", username, "p", password))).build()
			.thenApply(r->{
				isLoggedIn = true;
				return true;
			});
	}
	
	public CompletableFuture<Response.CollectionResponse> getServers(){
		final String url = this.restClient().buildUrl("servers.json");
		return ResponseFuture.builder(Response.CollectionResponse.class)
				.setMethod(ResponseFuture.Method.GET)
				.setUrl(url)
				.setSession(this.restClient())
				.build();
	}
	
	public CompletableFuture<Response.CollectionResponse> getDeliveryServices(){
		final String url = this.restClient().buildUrl("deliveryservices.json");
		LOG.debug("getDeliveryService url {}", url);
		return ResponseFuture.builder(Response.CollectionResponse.class)
				.setMethod(ResponseFuture.Method.GET)
				.setUrl(url)
				.setSession(this.restClient())
				.build();
	}
	
	
	@AutoValue
	protected abstract static class ResponseFuture<T extends Response> extends CompletableFuture<T> implements BiConsumer<HttpResponse, Throwable> {
		private static final Logger LOG = LoggerFactory.getLogger(ResponseFuture.class);
		public static enum Method{
			POST	, GET
		}
		
		public abstract Optional<BiConsumer<ResponseFuture<T>, Throwable>> handleException();
		public abstract Class<T> responseType();
		public abstract Method method();
		public abstract String url();
		public abstract RestApiSession session();
		public abstract Optional<String> body();
		
		private CompletableFuture<HttpResponse> subFuture;
		
		public static <T extends Response> Builder<T> builder(Class<T> response) {
			return new AutoValue_TOSession_ResponseFuture.Builder<T>()
					.setResponseType(response);
		}
		public static Builder<Response> builder() {
			return builder(Response.class);
		}
		
		@Override
		public boolean cancel(boolean mayInterruptIfRunning) {
			if(subFuture != null) {
				return subFuture.cancel(mayInterruptIfRunning);
			}
			return false;
		}
		
		public ResponseFuture<T> execute(){
			RequestBuilder rBuilder = RequestBuilder
					.create(this.method().toString())
					.setUri(this.url());
			if(this.body().isPresent()) {
				rBuilder.setEntity(new StringEntity(this.body().get(), Charsets.UTF_8));
			}
			subFuture = this.session().execute(rBuilder);
			subFuture.whenComplete(this);
			return this;
		}
		
		@AutoValue.Builder
		public abstract static class Builder<T extends Response> {
			public ResponseFuture<T> build(){
				return autoBuild().execute();
			}
			protected abstract ResponseFuture<T> autoBuild();
			public abstract Builder<T> setHandleException(BiConsumer<ResponseFuture<T>, Throwable> function);
			public abstract Builder<T> setResponseType(Class<T> respone);
			public abstract Builder<T> setMethod(Method method);
			public abstract Builder<T> setUrl(String url);
			public abstract Builder<T> setSession(RestApiSession session);
			public abstract Builder<T> setBody(String body);
		}
		
		@Override
		public void accept(HttpResponse res, Throwable u) {
			try {
				switch(res.getStatusLine().getStatusCode()) {
					case 200:
						break;
					case 401:
						_handleException(new LoginException("Login required"));
						return;
					default:
						_handleException(new OperationException(String.format("None 200 response: %s %s", res.getStatusLine().getStatusCode(), res.getStatusLine().getReasonPhrase())));
						return;
				}
				
				InputStreamReader r = new InputStreamReader(res.getEntity().getContent());
				T resp = gson.fromJson(r, responseType());
				if(resp.getAlerts() != null) {
					for(Alert a: resp.getAlerts()) {
						if("error".equals(a.getLevel())) {
							_handleException(new OperationException("Recieved error from server: "+ a.getText()));
							return;
						}
					}
				}
				
				this.complete(resp);
			} catch (UnsupportedOperationException | IOException e) {
				_handleException(new OperationException("Reading response failed", e));
				return;
			}
		}
		
		private void _handleException(Throwable t) {
			if(handleException().isPresent()) {
				handleException().get().accept(this, t);
			}
			
			if(!this.isDone()) {
				this.completeExceptionally(t);
			}
		}
	}
	
	public abstract RestApiSession restClient();
	
	public static Builder builder() {
		return new AutoValue_TOSession.Builder();
	}
	public abstract Builder toBuilder();
	
	@AutoValue.Builder
	public abstract static class Builder {
		public abstract TOSession build();
		
		public abstract Builder setRestClient(RestApiSession restClient);
		public abstract RestApiSession.Builder restClientBuilder();
	}
}
