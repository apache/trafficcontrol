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

import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URI;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.function.BiConsumer;

import org.apache.http.HttpResponse;
import org.apache.http.client.methods.RequestBuilder;
import org.apache.http.entity.StringEntity;
import org.apache.trafficcontrol.client.RestApiSession;
import org.apache.trafficcontrol.client.exception.LoginException;
import org.apache.trafficcontrol.client.exception.OperationException;
import org.apache.trafficcontrol.client.trafficops.models.Alert;
import org.apache.trafficcontrol.client.trafficops.models.Response;

import com.google.auto.value.AutoValue;
import com.google.common.base.Charsets;
import com.google.common.base.Strings;

@AutoValue
public abstract class ResponseFuture<T extends Response> extends CompletableFuture<T> implements BiConsumer<HttpResponse, Throwable> {
	public static enum Method{
		POST, GET
	}
	
	public abstract Optional<BiConsumer<ResponseFuture<T>, Throwable>> handleException();
	public abstract Class<T> responseType();
	public abstract ResponseFuture.Method method();
	public abstract URI uri();
	public abstract RestApiSession session();
	public abstract Optional<String> body();
	
	private CompletableFuture<HttpResponse> subFuture;
	
	public static <T extends Response> ResponseFuture.Builder<T> builder(Class<T> response) {
		return new AutoValue_ResponseFuture.Builder<T>()
				.setResponseType(response);
	}
	public static ResponseFuture.Builder<Response> builder() {
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
				.setUri(this.uri());
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
		public abstract ResponseFuture.Builder<T> setHandleException(BiConsumer<ResponseFuture<T>, Throwable> function);
		public abstract ResponseFuture.Builder<T> setResponseType(Class<T> respone);
		public abstract ResponseFuture.Builder<T> setMethod(ResponseFuture.Method method);
		public abstract ResponseFuture.Builder<T> setUri(URI uri);
		public ResponseFuture.Builder<T> setUri(String uri){
			if(!Strings.isNullOrEmpty(uri)) {
				this.setUri(URI.create(uri));
			}
			return this;
		}
		
		public abstract ResponseFuture.Builder<T> setSession(RestApiSession session);
		public abstract ResponseFuture.Builder<T> setBody(String body);
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
			T resp = TOSession.gson.fromJson(r, responseType());
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