package com.comcast.cdn.traffic_control.traffic_router.core;

import org.apache.catalina.Wrapper;
import org.apache.catalina.connector.Connector;
import org.apache.catalina.core.StandardContext;
import org.apache.catalina.core.StandardHost;
import org.apache.catalina.core.StandardService;
import org.apache.catalina.startup.Catalina;
import org.apache.juli.logging.LogFactory;

import java.util.logging.Level;

public class CatalinaTrafficRouter {
	Catalina catalina;

	public CatalinaTrafficRouter(String serverXmlPath, String appBase) {

		java.util.logging.Logger logger = java.util.logging.Logger.getLogger("");

		java.util.logging.Handler[] handlers = logger.getHandlers();
		for(java.util.logging.Handler handler : handlers) {
			handler.setLevel(Level.SEVERE);
		}

		System.setProperty("dns.tcp.port", "1053");
		System.setProperty("dns.udp.port", "1053");

		System.setProperty("catalina.home", "");

		catalina = new Catalina();
		catalina.process(new String[] {"-config", serverXmlPath});
		catalina.load();

		// Override the port and app base property of server.xml
		StandardService trafficRouterService = (StandardService) catalina.getServer().findService("traffic_router_core");

		Connector[] connectors =trafficRouterService.findConnectors();
		for (Connector connector : connectors) {
			if (connector.getPort() == 80) {
				connector.setPort(8888);
			}
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
