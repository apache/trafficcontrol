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

package org.apache.traffic_control.traffic_router.core.util;

import java.io.IOException;
import java.net.HttpCookie;
import java.net.HttpURLConnection;

public class ProtectedFetcher extends Fetcher {
	private String authorizationEndpoint;
	private String data;
	private HttpCookie cookie;

	public ProtectedFetcher(final String authorizationEndpoint, final String data, final int timeout) {
		this.timeout = (timeout > 0) ? timeout : DEFAULT_TIMEOUT;
		this.setAuthorizationEndpoint(authorizationEndpoint);
		this.setData(data);
	}

	@Override
	protected HttpURLConnection getConnection(final String url, final String data, final String method, final long lastFetchedTime) throws IOException {

		if (isCookieValid()) {
			final HttpURLConnection connection = extractCookie(super.getConnection(url, data, method, lastFetchedTime));
			if (connection.getResponseCode() != HttpURLConnection.HTTP_UNAUTHORIZED) {
				return connection;
			}
		}

		extractCookie(super.getConnection(getAuthorizationEndpoint(), getData(), POST_STR, 0L));
		return extractCookie(super.getConnection(url, data, method, lastFetchedTime));
	}

	private HttpURLConnection extractCookie(final HttpURLConnection http) throws IOException {
		if ((http != null) &&  (http.getHeaderField("Set-Cookie") != null)) {
			setCookie(HttpCookie.parse(http.getHeaderField("Set-Cookie")).get(0));
		}
		return http;
	}

	private boolean isCookieValid() {
		if (cookie != null && !cookie.hasExpired()) {
			return true;
		} else {
			return false;
		}
	}

	private void setCookie(final HttpCookie cookie) {
		this.cookie = cookie;

		if (this.cookie != null) {
			requestProps.put("Cookie", this.cookie.toString());
		} else {
			requestProps.remove("Cookie");
		}
	}

	private String getAuthorizationEndpoint() {
		return authorizationEndpoint;
	}

	private void setAuthorizationEndpoint(final String authorizationEndpoint) {
		this.authorizationEndpoint = authorizationEndpoint;
	}

	private String getData() {
		return data;
	}

	private void setData(final String data) {
		this.data = data;
	}

	@Override
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.IfStmtsMustUseBraces"})
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;
		if (!super.equals(o)) return false;

		final ProtectedFetcher that = (ProtectedFetcher) o;

		if (authorizationEndpoint != null ? !authorizationEndpoint.equals(that.authorizationEndpoint) : that.authorizationEndpoint != null)
			return false;
		if (data != null ? !data.equals(that.data) : that.data != null) return false;
		return !(cookie != null ? !cookie.equals(that.cookie) : that.cookie != null);

	}

	@Override
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public int hashCode() {
		int result = super.hashCode();
		result = 31 * result + (authorizationEndpoint != null ? authorizationEndpoint.hashCode() : 0);
		result = 31 * result + (data != null ? data.hashCode() : 0);
		result = 31 * result + (cookie != null ? cookie.hashCode() : 0);
		return result;
	}
}
