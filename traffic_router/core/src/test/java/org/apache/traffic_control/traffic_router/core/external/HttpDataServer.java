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

package org.apache.traffic_control.traffic_router.core.external;

import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpCookie;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.URI;
//import java.util.logging.Logger;

public class HttpDataServer implements HttpHandler {
	private HttpServer httpServer;
	private int testHttpServerPort;

	public HttpDataServer(int testHttpServerPort) {
		this.testHttpServerPort = testHttpServerPort;
	}
	private boolean receivedSteeringPost = false;
	private boolean receivedCertificatesPost = false;
	private boolean receivedCrConfig2Post = false;
	private boolean receivedCrConfig3Post = false;
	private boolean receivedCrConfig4Post = false;

// Useful for producing an access log
//	static {
//		Logger logger = Logger.getLogger("com.sun.net.httpserver");
//		logger.setLevel(java.util.logging.Level.ALL);
//
//		java.util.logging.Handler[] handlers = logger.getHandlers();
//		for (java.util.logging.Handler handler : handlers) {
//			handler.setLevel(java.util.logging.Level.ALL);
//		}
//	}

	public void start(int port) throws IOException {
		httpServer = HttpServer.create(new InetSocketAddress(InetAddress.getLoopbackAddress(), port),10);
		httpServer.createContext("/", this);
		httpServer.start();
		System.out.println(">>>>>>>>>>>>> Started Fake Http Data Server at " + port);
	}

	public void stop() {
		System.out.println(">>>>>>>>>>>>>> Stopping Fake Http Data Server");
		httpServer.stop(10);
		System.out.println(">>>>>>>>>>>>>> STOPPED Fake Http Data Server");
	}

	@Override
	public void handle(final HttpExchange httpExchange) throws IOException {

		new Thread(new Runnable() {
			@Override
			public void run() {
				if ("POST".equals(httpExchange.getRequestMethod()) ) {
					if (!receivedSteeringPost && "/steering".equals(httpExchange.getRequestURI().getPath())) {
						receivedSteeringPost = true;
					}

					if (!receivedCertificatesPost && "/certificates".equals(httpExchange.getRequestURI().getPath())) {
						receivedCertificatesPost = true;
					}

					if (!receivedCrConfig2Post && "/crconfig-2".equals(httpExchange.getRequestURI().getPath())) {
						receivedCrConfig2Post = true;
						receivedCrConfig3Post = false;
						receivedCrConfig4Post = false;
					}

					if (!receivedCrConfig3Post && "/crconfig-3".equals(httpExchange.getRequestURI().getPath())) {
						receivedCrConfig2Post = false;
						receivedCrConfig3Post = true;
						receivedCrConfig4Post = false;
					}

					if (!receivedCrConfig4Post && "/crconfig-4".equals(httpExchange.getRequestURI().getPath())) {
						receivedCrConfig2Post = false;
						receivedCrConfig3Post = false;
						receivedCrConfig4Post = true;
					}

					try {
						httpExchange.sendResponseHeaders(200,0);
					} catch (IOException e) {
						System.out.println(">>>>> failed acknowledging post");
					}
					return;
				}

				URI uri = httpExchange.getRequestURI();
				String path = uri.getPath();

				if (path.startsWith("/")) {
					path = path.substring(1);
				}

				String query = uri.getQuery();
				if ("json".equals(query)) {
					path += ".json";
				}

				if (("api/" + TrafficOpsUtils.TO_API_VERSION + "/user/login").equals(path)) {
					try {
						Headers headers = httpExchange.getResponseHeaders();
						headers.add("Content-length", Integer.toString(0));
						headers.set("Set-Cookie", new HttpCookie("mojolicious","fake-cookie").toString());
						httpExchange.sendResponseHeaders(200,0);

					} catch (Exception e) {
						System.out.println(">>>> Failed setting cookie");
					}
				}

				// Pretend that someone externally changed steering.json data
				if (receivedSteeringPost && ("api/" + TrafficOpsUtils.TO_API_VERSION + "/steering").equals(path)) {
					path = "api/" + TrafficOpsUtils.TO_API_VERSION + "/steering2";
				}

				// pretend certificates have not been updated
				if (!receivedCertificatesPost && ("api/" + TrafficOpsUtils.TO_API_VERSION + "/cdns/name/thecdn/sslkeys").equals(path)) {
					path = path.replace("/sslkeys", "/sslkeys-missing-1");
				}

				if (path.contains("CrConfig") && receivedCrConfig2Post) {
					path = path.replace("CrConfig", "CrConfig2");
				}

				if (path.contains("CrConfig") && receivedCrConfig3Post) {
					path = path.replace("CrConfig", "CrConfig3");
				}

				if (path.contains("CrConfig") && receivedCrConfig4Post) {
					path = path.replace("CrConfig", "CrConfig4");
				}

				InputStream inputStream = getClass().getClassLoader().getResourceAsStream(path);

				if (inputStream == null) {
					System.out.println(">>> " + path + " not found");
					String response = "404 (Not Found)\n";

					OutputStream os = null;
					try {
						httpExchange.sendResponseHeaders(404, response.length());
						os = httpExchange.getResponseBody();
						os.write(response.getBytes());
					} catch (Exception e) {
						System.out.println("Failed sending 404!: " + e.getMessage());
					} finally {
						if (os != null) try {
							os.close();
						} catch (IOException e) {
							System.out.println("Failed closing output stream!: " + e.getMessage());
						}
						return;
					}
				}

				if (path.contains("Geo")) {
					try (OutputStream os = httpExchange.getResponseBody()) {
						int bodySz = inputStream.available();
						httpExchange.getResponseHeaders().add("Content-length", Integer.toString(bodySz));
						httpExchange.sendResponseHeaders(200, bodySz);

						final byte[] buffer = new byte[0x10000];
						int count;

						while ((count = inputStream.read(buffer)) >= 0) {
							os.write(buffer, 0, count);
						}
					} catch (Exception e) {
						System.out.println("Failed sending data for " + path + " : " + e.getMessage());
					}
				} else {
					final byte[] buffer = new byte[0x10000];
					final StringBuilder stringBuilder = new StringBuilder();

					try {
						while (inputStream.read(buffer) >= 0) {
							stringBuilder.append(new String(buffer));
						}
					} catch (Exception e) {
						System.out.println("Failed to ingest input file associated with " + path);
					}

					String body = stringBuilder.toString();

					try {
						if (path.contains("CrConfig")) {
							body = body.replaceAll("localhost:8889" , "localhost:" + testHttpServerPort);
						}

						if (path.contains("json")) {
							httpExchange.getResponseHeaders().add("Content-type", "application/json");
						}

						final byte[] bodyBytes = body.getBytes();
						httpExchange.getResponseHeaders().add("Content-length", Integer.toString(bodyBytes.length));
						httpExchange.sendResponseHeaders(200, bodyBytes.length);
						httpExchange.getResponseBody().write(bodyBytes);
						httpExchange.getResponseBody().close();

					} catch (Exception e) {
						System.out.println("Failed sending data for " + path + " : " + e.getMessage());
					}
				}

				try {
					inputStream.close();
				} catch (Exception e) {
					System.out.println("Failed closing stream!: " + e.getMessage());
				}

			}
		}).start();
	}
}
