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

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.Reader;
import java.net.Authenticator;
import java.net.InetSocketAddress;
import java.net.PasswordAuthentication;
import java.net.Proxy;
import java.net.URLConnection;
import java.util.HashMap;
import java.util.Map;

import java.net.URL;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
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

	private static final Map<String,PasswordAuthentication> passwds = new HashMap<String,PasswordAuthentication>();
	static {
		Authenticator.setDefault(new MyAuthenticator());

		try {
			final SSLContext ctx = SSLContext.getInstance("TLS");
			ctx.init(new KeyManager[0], new TrustManager[] {new DefaultTrustManager()}, new SecureRandom());
			SSLContext.setDefault(ctx);
		} catch (Exception e) {
			LOGGER.warn(e,e);
		}
	}
	static public void setUserPw(final String host, final String user, final String passwd) {
		passwds.put(host, new PasswordAuthentication(user, passwd.toCharArray()));
	}

	private static class DefaultTrustManager implements X509TrustManager {

		@Override
		public void checkClientTrusted(final X509Certificate[] arg0, final String arg1) throws CertificateException {}
		@Override
		public void checkServerTrusted(final X509Certificate[] arg0, final String arg1) throws CertificateException {}
		@Override
		public X509Certificate[] getAcceptedIssuers() { return null; }
	}


	public static String fetchContent(final String link) throws IOException {
		return fetchContent(new URL(link).openConnection());
	}

	public static class MyAuthenticator extends Authenticator {
		protected PasswordAuthentication getPasswordAuthentication() {
			return passwds.get(getRequestingHost());
		}
	}

	public static String fetchContent(final String stateUrl, final String ipStr) throws IOException {
		final URL url = new URL(stateUrl);
		final Proxy proxy = new Proxy(Proxy.Type.HTTP, new InetSocketAddress(
				ipStr, 80));

		final URLConnection conn = url.openConnection(proxy);

		final char[] buffer = new char[0x10000];
		final StringBuilder out = new StringBuilder();
		final Reader in = new InputStreamReader(conn.getInputStream(), "UTF-8");
		try {
			int read;
			do {
				read = in.read(buffer, 0, buffer.length);
				if (read>0) {
					out.append(buffer, 0, read);
				}
			} while (read>=0);
		} finally {
			in.close();
		}
		return out.toString();
	}

	public static String fetchContent(final URLConnection conn) throws IOException {
		conn.setAllowUserInteraction(true);
		final char[] buffer = new char[0x10000];
		final StringBuilder out = new StringBuilder();
		final Reader in = new InputStreamReader(conn.getInputStream(), "UTF-8");
		try {
			int read;
			do {
				read = in.read(buffer, 0, buffer.length);
				if (read>0) {
					out.append(buffer, 0, read);
				}
			} while (read>=0);
		} finally {
			in.close();
		}
		return out.toString();
	}

	public static String fetchDataFromServer(final String url) throws IOException {
		LOGGER.warn("__ENTERING fetchDataFromServer()");

		final URL u = new URL(url);
		final HttpsURLConnection http = (HttpsURLConnection)u.openConnection();
		http.setRequestMethod("GET");
		return  fetchContent(http);    
	}

	protected static String tmpPrefix = "loc";
	protected static String tmpSuffix = ".dat";
	public static File downloadFile(final String url) throws IOException {
		LOGGER.debug("Downloading file from: " + url);
		final URL u = new URL(url);
		final URLConnection urlc = u.openConnection();
		if(urlc instanceof HttpsURLConnection) {
			final HttpsURLConnection http = (HttpsURLConnection)urlc;
			http.setHostnameVerifier(new HostnameVerifier() {
				@Override
				public boolean verify(final String arg0, final SSLSession arg1) {
					return true;
				}
			});
			http.setRequestMethod("GET");
			http.setAllowUserInteraction(true);
		}
		final InputStream in = urlc.getInputStream();//new GZIPInputStream(dbURL.openStream());
		//		if(sourceCompressed) { in = new GZIPInputStream(in); }

		final File outputFile = File.createTempFile(tmpPrefix, tmpSuffix);
		final OutputStream out = new FileOutputStream(outputFile);

		IOUtils.copy(in, out);
		IOUtils.closeQuietly(in);
		IOUtils.closeQuietly(out);
		LOGGER.debug("Successfully downloaded file.");

		return outputFile;
	}

}
