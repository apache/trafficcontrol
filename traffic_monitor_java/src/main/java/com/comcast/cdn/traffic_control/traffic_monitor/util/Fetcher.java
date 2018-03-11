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

package com.comcast.cdn.traffic_control.traffic_monitor.util;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpCookie;
import java.net.InetSocketAddress;
import java.net.Proxy;
import java.net.URL;
import java.net.URLConnection;
import java.net.URLEncoder;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Map;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.KeyManager;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;

public class Fetcher {
	private static final Logger LOGGER = Logger.getLogger(Fetcher.class);
	private static final String GET_STR = "GET";
	private static final String UTF8_STR = "UTF-8";

	static {
		try {
			final SSLContext ctx = SSLContext.getInstance("TLS");
			ctx.init(new KeyManager[0], new TrustManager[] {new DefaultTrustManager()}, new SecureRandom());
			SSLContext.setDefault(ctx);
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

	public static String fetchContent(final String link, final int timeout) throws IOException {
		return fetchContent(new URL(link).openConnection(), timeout);
	}

	public static String fetchContent(final String link, final Map<String, String> headerMap, final int timeout) throws IOException {
		final URL url = new URL(link);
		final URLConnection conn = url.openConnection();

		for (String key : headerMap.keySet()) {
			conn.setRequestProperty(key, headerMap.get(key));
		}

		return fetchContent(conn, timeout);
	}

	public static String fetchContent(final String stateUrl, final String ipStr, final int port, final int timeout) throws IOException {
		final URLConnection conn = new URL(stateUrl).openConnection(
				new Proxy(Proxy.Type.HTTP, new InetSocketAddress(ipStr, port)));
		if(timeout!=0) {
			conn.setConnectTimeout(timeout);
			conn.setReadTimeout(timeout);
		}
		conn.connect();

		return IOUtils.toString(new InputStreamReader(conn.getInputStream(), UTF8_STR));
	}

	public static String fetchContent(final URLConnection conn, final int timeout) throws IOException {
		conn.setAllowUserInteraction(true);
		if(timeout!=0) {
			conn.setConnectTimeout(timeout);
			conn.setReadTimeout(timeout);
		}
		conn.connect();
		return IOUtils.toString(new InputStreamReader(conn.getInputStream(), UTF8_STR));
	}

	public static String fetchDataFromServer(final String url) throws IOException {
		LOGGER.warn("__ENTERING fetchDataFromServer()");

		final URL u = new URL(url);
		final HttpsURLConnection http = (HttpsURLConnection)u.openConnection();
		http.setRequestMethod(GET_STR);
		return  fetchContent(http, 0);    
	}

	protected static String tmpPrefix = "loc";
	protected static String tmpSuffix = ".dat";
	public static File downloadFile(final String url) throws IOException {
		InputStream in = null;
		OutputStream out = null;
		try {
			LOGGER.info("downloadFile: " + url);
			final URL u = new URL(url);
			final URLConnection urlc = u.openConnection();
			if(urlc instanceof HttpsURLConnection) {
				final HttpsURLConnection http = (HttpsURLConnection)urlc;
				http.setInstanceFollowRedirects(false);
				http.setHostnameVerifier(new HostnameVerifier() {
					@Override
					public boolean verify(final String arg0, final SSLSession arg1) {
						return true;
					}
				});
				http.setRequestMethod(GET_STR);
				http.setAllowUserInteraction(true);
			}
			in = urlc.getInputStream();//new GZIPInputStream(dbURL.openStream());
			//		if(sourceCompressed) { in = new GZIPInputStream(in); }

			final File outputFile = File.createTempFile(tmpPrefix, tmpSuffix);
			out = new FileOutputStream(outputFile);

			IOUtils.copy(in, out);
			return outputFile;
		} finally {
			IOUtils.closeQuietly(in);
			IOUtils.closeQuietly(out);
		}
	}
	public static String fetchSecureContent(final String url, final int timeout) throws IOException {
		LOGGER.info("fetchSecureContent: " + url);
		final URL u = new URL(url);
		final URLConnection conn = u.openConnection();
		if(timeout!=0) {
			conn.setConnectTimeout(timeout);
			conn.setReadTimeout(timeout);
		}
		if(conn instanceof HttpsURLConnection) {
			final HttpsURLConnection http = (HttpsURLConnection)conn;
			http.setHostnameVerifier(new HostnameVerifier() {
				@Override
				public boolean verify(final String arg0, final SSLSession arg1) {
					return true;
				}
			});
			http.setRequestMethod(GET_STR);
			http.setAllowUserInteraction(true);
		}
		return IOUtils.toString(conn.getInputStream());
	}

	private static HttpCookie tmCookie;

	private static HttpCookie getTmCookie(final String url, final String username, final String password, final int timeout) throws IOException {
		if (tmCookie != null && !tmCookie.hasExpired()) {
			return tmCookie;
		}

		final String charset = UTF8_STR;
		final String query = String.format("u=%s&p=%s",
				URLEncoder.encode(username, charset),
				URLEncoder.encode(password, charset));
		final URLConnection connection = new URL(url).openConnection();

		if (!(connection instanceof HttpsURLConnection)) {
			return null;
		}

		final HttpsURLConnection http = (HttpsURLConnection) connection;

		http.setInstanceFollowRedirects(false);

		http.setHostnameVerifier(new HostnameVerifier() {
			@Override
			public boolean verify(final String arg0, final SSLSession arg1) {
				return true;
			}
		});

		http.setRequestMethod("POST");
		http.setAllowUserInteraction(true);

		if (timeout != 0) {
			http.setConnectTimeout(timeout);
			http.setReadTimeout(timeout);
		}

		http.setDoOutput(true); // Triggers POST.
		http.setRequestProperty("Accept-Charset", charset);
		http.setRequestProperty("Content-Type", "application/x-www-form-urlencoded;charset=" + charset);

		OutputStream output = null;

		try {
			output = http.getOutputStream();
			output.write(query.getBytes(charset));
		} finally {
			if (output != null) {
				try {
					output.close();
				} catch (IOException e) {
					LOGGER.debug(e,e);
				}
			}
		}

		LOGGER.info("fetching cookie: " + url);
		connection.connect();

		tmCookie = HttpCookie.parse(http.getHeaderField("Set-Cookie")).get(0);
		LOGGER.debug("cookie: "+ tmCookie);

		return tmCookie;
	}
	public static File downloadTM(final String url, final String authUrl, final String username, final String password) throws IOException {
		return downloadTM(url, authUrl, username, password, 0);
	}
	public static File downloadTM(final String url, final String authUrl, final String username, final String password, final int timeout) throws IOException {
		InputStream in = null;
		OutputStream out = null;

		try {
			final URL u = new URL(url);
			final URLConnection urlc = u.openConnection();

			if (timeout != 0) {
				urlc.setConnectTimeout(timeout);
				urlc.setReadTimeout(timeout);
			}

			if (urlc instanceof HttpsURLConnection) {
				final String cookie = getTmCookie(authUrl, username, password, timeout).toString();

				final HttpsURLConnection http = (HttpsURLConnection)urlc;
				http.setInstanceFollowRedirects(false);
				http.setHostnameVerifier(new HostnameVerifier() {
					@Override
					public boolean verify(final String arg0, final SSLSession arg1) {
						return true;
					}
				});
				http.setRequestMethod(GET_STR);
				http.setAllowUserInteraction(true);
				http.addRequestProperty("Cookie", cookie);
			}

			in = urlc.getInputStream();

			final File outputFile = File.createTempFile(tmpPrefix, tmpSuffix);
			out = new FileOutputStream(outputFile);

			IOUtils.copy(in, out);
			return outputFile;
		} finally {
			IOUtils.closeQuietly(in);
			IOUtils.closeQuietly(out);
		}
	}

	public static void clearTmCookie() {
		tmCookie = null;
	}

}
