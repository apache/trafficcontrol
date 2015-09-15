/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.util;

import java.io.IOException;
import java.net.HttpCookie;
import java.net.HttpURLConnection;

import org.apache.log4j.Logger;

public class ProtectedFetcher extends Fetcher {
	private static final Logger LOGGER = Logger.getLogger(ProtectedFetcher.class);
	private String endpoint;
	private String data;
	private HttpCookie cookie;

	public ProtectedFetcher(final String endpoint, final String data, final int timeout) {
		this.timeout = (timeout > 0) ? timeout : DEFAULT_TIMEOUT;
		this.setEndpoint(endpoint);
		this.setData(data);
	}

	public ProtectedFetcher(final String endpoint, final String data) {
		this(endpoint, data, DEFAULT_TIMEOUT);
	}

	@Override
	protected HttpURLConnection getConnection(final String url, final String data, final String method) throws IOException {
		if (!isCookieValid()) {
			LOGGER.debug("Cookie is no longer valid; re-authenticating to " + getEndpoint() + " cookie = " + cookie);
			extractCookie(super.getConnection(getEndpoint(), getData(), POST_STR));
		}

		return extractCookie(super.getConnection(url, data, method));
	}

	private HttpURLConnection extractCookie(final HttpURLConnection http) throws IOException {
		if (http.getHeaderField("Set-Cookie") != null) {
			LOGGER.info("Storing cookie from: " + http.getURL().toString());
			setCookie(HttpCookie.parse(http.getHeaderField("Set-Cookie")).get(0));
			LOGGER.debug("cookie: "+ getCookie());
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

	private HttpCookie getCookie() throws IOException {
		return cookie;
	}

	private void setCookie(final HttpCookie cookie) {
		this.cookie = cookie;

		if (this.cookie != null) {
			requestProps.put("Cookie", this.cookie.toString());
		} else {
			requestProps.remove("Cookie");
		}
	}

	private String getEndpoint() {
		return endpoint;
	}

	private void setEndpoint(final String endpoint) {
		this.endpoint = endpoint;
	}

	private String getData() {
		return data;
	}

	private void setData(final String data) {
		this.data = data;
	}
}
