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

package com.comcast.cdn.traffic_control.traffic_router.core;

import org.apache.catalina.Wrapper;
import org.apache.catalina.connector.Connector;
import org.apache.catalina.core.StandardContext;
import org.apache.catalina.core.StandardHost;
import org.apache.catalina.core.StandardService;
import org.apache.catalina.startup.Catalina;
import org.springframework.util.SocketUtils;

import java.util.logging.Level;

public class CatalinaTrafficRouter {
	Catalina catalina;

	public CatalinaTrafficRouter(String serverXmlPath, String appBase) {
		System.setProperty("java.util.logging.SimpleFormatter.format", "%1$tFT%1$tT.%1$tL [%4$s] %5$s %6$s%n");
		java.util.logging.Logger logger = java.util.logging.Logger.getLogger("");
		java.util.logging.Handler[] handlers = logger.getHandlers();
		for (java.util.logging.Handler handler : handlers) {
			handler.setLevel(Level.WARNING);
		}

		System.setProperty("dns.tcp.port", "1053");
		System.setProperty("dns.udp.port", "1053");

		System.setProperty("catalina.home", "");

		catalina = new Catalina();
		catalina.process(new String[] {"-config", serverXmlPath});
		catalina.load();

		// Override the port and app base property of server.xml
		StandardService trafficRouterService = (StandardService) catalina.getServer().findService("traffic_router_core");

		Connector[] connectors = trafficRouterService.findConnectors();
		for (Connector connector : connectors) {
			if (connector.getPort() == 80) {
				connector.setPort(Integer.parseInt(System.getProperty("routerHttpPort", "8888")));
			}

			SocketUtils.findAvailableTcpPort();

			if (connector.getPort() == 443) {
				connector.setPort(Integer.parseInt(System.getProperty("routerSecurePort", "8443")));
			}
			System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>>> Traffic Router listening on port " + connector.getPort() + " " + connector.getScheme());

		}

		StandardHost standardHost = (StandardHost) trafficRouterService.getContainer().findChild("localhost");
		standardHost.setAppBase(appBase);

		// We have to manually set up the default servlet, the Catalina class doesn't do this for us
		StandardContext rootContext = (StandardContext) standardHost.findChild("/");
		Wrapper defaultServlet = rootContext.createWrapper();
		defaultServlet.setName("default");
		defaultServlet.setServletClass("org.apache.catalina.servlets.DefaultServlet");
		defaultServlet.addInitParameter("debug", "0");
		defaultServlet.addInitParameter("listings", "false");
		defaultServlet.setLoadOnStartup(1);

		rootContext.addChild(defaultServlet);
	}

	public void start() {
		catalina.setAwait(false);
		catalina.start();
	}

	public void stop() {
		catalina.stop();
	}
}
