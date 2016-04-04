package com.comcast.cdn.traffic_control.traffic_router.core;

import org.apache.catalina.Engine;
import org.apache.catalina.Host;
import org.apache.catalina.LifecycleException;
import org.apache.catalina.Wrapper;
import org.apache.catalina.connector.Connector;
import org.apache.catalina.core.StandardContext;
import org.apache.catalina.startup.Embedded;

public class EmbeddedTrafficRouter {
	private Embedded tomcat;

	public EmbeddedTrafficRouter() {
		this("src/main");
	}

	public EmbeddedTrafficRouter(String webappParent) {
		tomcat = new Embedded();
		tomcat.setName("Traffic Router");
		Host localhost = tomcat.createHost("localhost", webappParent);
		localhost.setAutoDeploy(false);

		StandardContext rootContext = (StandardContext) tomcat.createContext("/", "webapp");

		Wrapper defaultServlet = rootContext.createWrapper();
		defaultServlet.setName("default");
		defaultServlet.setServletClass("org.apache.catalina.servlets.DefaultServlet");
		defaultServlet.addInitParameter("debug", "0");
		defaultServlet.addInitParameter("listings", "false");
		defaultServlet.setLoadOnStartup(1);

		rootContext.addChild(defaultServlet);
		rootContext.setDefaultWebXml("web.xml");

		localhost.addChild(rootContext);

		Engine engine = tomcat.createEngine();
		engine.setName("Traffic Router Engine");
		engine.setDefaultHost(localhost.getName());
		engine.addChild(localhost);

		tomcat.addEngine(engine);
		tomcat.addConnector(tomcat.createConnector(localhost.getName(), 8888, false));
	}

	public void start() throws LifecycleException {
		tomcat.setAwait(true);
		tomcat.start();
	}

	public void stop() throws LifecycleException {
		tomcat.stop();
	}
}
