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

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.HashMap;
import java.util.Map;
import java.util.zip.GZIPInputStream;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class Fetcher {
	private static final Logger LOGGER = LogManager.getLogger(Fetcher.class);
	protected static final String GET_STR = "GET";
	protected static final String POST_STR = "POST";
	protected static final String UTF8_STR = "UTF-8";
	protected static final int DEFAULT_TIMEOUT = 10000;
	private static final String GZIP_ENCODING_STRING = "gzip";
	private static final String CONTENT_TYPE_STRING = "Content-Type";
	protected static final String CONTENT_TYPE_JSON = "application/json";
	protected int timeout = DEFAULT_TIMEOUT; // override if you want something different
	protected final Map<String, String> requestProps = new HashMap<String, String>();


	static {
		try {
			// TODO: make disabling self signed certificates configurable
			final SSLContext ctx = SSLContext.getInstance("SSL");
			ctx.init(null, new TrustManager[] {new DefaultTrustManager()}, new SecureRandom());
			SSLContext.setDefault(ctx);
			HttpsURLConnection.setDefaultSSLSocketFactory(ctx.getSocketFactory());
		} catch (Exception e) {
			LOGGER.warn(e,e);
		}
	}

	private static class DefaultTrustManager implements X509TrustManager {
		@Override
		public void checkClientTrusted(final X509Certificate[] arg0, final String arg1) throws CertificateException {}
		@Override
		public void checkServerTrusted(final X509Certificate[] arg0, final String arg1) throws CertificateException {}
		@Override
		public X509Certificate[] getAcceptedIssuers() { return null; }
	}

	protected HttpURLConnection getConnection(final String url, final String data, final String requestMethod, final long lastFetchTime) throws IOException {
		return getConnection(url, data, requestMethod, lastFetchTime, null);
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	protected HttpURLConnection getConnection(final String url, final String data, final String requestMethod, final long lastFetchTime, final String contentType) throws IOException {
		HttpURLConnection http = null;
		try {
			String method = GET_STR;

			if (requestMethod != null) {
				method = requestMethod;
			}

			LOGGER.info(method + "ing: " + url + "; timeout is " + timeout);

			final URLConnection connection = new URL(url).openConnection();

			connection.setIfModifiedSince(lastFetchTime);

			if (timeout != 0) {
				connection.setConnectTimeout(timeout);
				connection.setReadTimeout(timeout);
			}

			http = (HttpURLConnection) connection;

			if (connection instanceof HttpsURLConnection) {
				final HttpsURLConnection https = (HttpsURLConnection) connection;
				https.setHostnameVerifier(new HostnameVerifier() {
					@Override
					public boolean verify(final String arg0, final SSLSession arg1) {
						return true;
					}
				});
			}

			http.setInstanceFollowRedirects(false);
			http.setRequestMethod(method);
			http.setAllowUserInteraction(true);
			http.addRequestProperty("Accept-Encoding", GZIP_ENCODING_STRING);

			for (final String key : requestProps.keySet()) {
				http.addRequestProperty(key, requestProps.get(key));
			}

			if (contentType != null) {
				http.addRequestProperty(CONTENT_TYPE_STRING, contentType);
			}

			if (method.equals(POST_STR) && data != null) {
				http.setDoOutput(true); // Triggers POST.

				try (OutputStream output = http.getOutputStream()) {
					output.write(data.getBytes(UTF8_STR));
				}
			}

			connection.connect();

		} catch (IOException e) {
			// For IO exceptions disconnect the connection down and propagate the exception upward
			final String failureMessage = connectionFailed(http, url);
			LOGGER.error(failureMessage + " \"" + e.toString() + "\"");
			throw(e);

		} catch (Exception e) {
			// For other exceptions mimic existing functionality - attempt to disconnect the
			// connection but squelch the exception
			final String failureMessage = connectionFailed(http, url);
			LOGGER.error(failureMessage + " \"" + e.toString() + "\"");
		}

		return http;
	}

	public String fetchIfModifiedSince(final String url, final long lastFetchTime) throws IOException {
		return fetchIfModifiedSince(url, null, null, lastFetchTime);
	}

	public String fetch(final String url) throws IOException {
		return fetch(url, null, null);
	}

	private String fetchIfModifiedSince(final String url, final String data, final String method, final long lastFetchTime) throws IOException {
		String ifModifiedSince = null;
		final HttpURLConnection connection = getConnection(url, data, method, lastFetchTime);
		if (connection != null) {
			if (connection.getResponseCode() == HttpURLConnection.HTTP_NOT_MODIFIED) {
				return null;
			}

			if (connection.getResponseCode() > 399) {
				LOGGER.warn("Failed Http Request to " + url + " Status " + connection.getResponseCode());
				return null;
			}

			final StringBuilder sb = new StringBuilder();
			createStringBuilderFromResponse(sb, connection);
			ifModifiedSince = sb.toString();
		}
		return ifModifiedSince;
	}

	public int getIfModifiedSince(final String url, final long lastFetchTime, final StringBuilder stringBuilder) throws IOException {
		int status = 0;
		final HttpURLConnection connection = getConnection(url, null, "GET", lastFetchTime);
		if (connection != null) {
			status = connection.getResponseCode();

			if (status == HttpURLConnection.HTTP_NOT_MODIFIED) {
				return status;
			}

			if (connection.getResponseCode() > 399) {
				LOGGER.warn("Failed Http Request to " + url + " Status " + connection.getResponseCode());
				return status;
			}

			createStringBuilderFromResponse(stringBuilder, connection);

		}
		return status;
	}

	public String fetch(final String url, final String data, final String method) throws IOException {
		return fetchIfModifiedSince(url, data, method, 0L);
	}

	@Override
	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final Fetcher fetcher = (Fetcher) o;

		if (timeout != fetcher.timeout) return false;
		return !(requestProps != null ? !requestProps.equals(fetcher.requestProps) : fetcher.requestProps != null);

	}

	@Override
	public int hashCode() {
		int result = timeout;
		result = 31 * result + (requestProps != null ? requestProps.hashCode() : 0);
		return result;
	}

	public void createStringBuilderFromResponse (final StringBuilder sb, final HttpURLConnection connection) throws IOException {
		if (GZIP_ENCODING_STRING.equals(connection.getContentEncoding())) {
			try (GZIPInputStream zippedInputStream =  new GZIPInputStream(connection.getInputStream());
				 BufferedReader r = new BufferedReader(new InputStreamReader(zippedInputStream))) {
				String input;
				while ((input = r.readLine()) != null) {
					sb.append(input);
				}
			}
		} else {
			try (BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()))) {
				String input;

				while ((input = in.readLine()) != null) {
					sb.append(input);
				}
			}
		}
	}

	private String connectionFailed(final HttpURLConnection http, final String url) {
		String httpUrl = url;
		String responseCode = "(none)";

		try {
			httpUrl = http.getURL().toString();
			responseCode = String.valueOf(http.getResponseCode());
		} catch (Exception e2) {
			// Don't care
			LOGGER.debug("Exception during call attempt to retrieve url or responseCode from http");
		}

		if (http != null) {
			http.disconnect();
		}
		return String.format("Failed Http Request to %s, status code: %s", httpUrl, responseCode);
	}

}
