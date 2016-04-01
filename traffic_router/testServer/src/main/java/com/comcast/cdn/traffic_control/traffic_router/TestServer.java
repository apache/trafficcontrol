package com.comcast.cdn.traffic_control.traffic_router;

import org.apache.log4j.ConsoleAppender;
import org.apache.log4j.Level;
import org.apache.log4j.LogManager;
import org.apache.log4j.PatternLayout;
import org.eclipse.jetty.jmx.MBeanContainer;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.server.bio.SocketConnector;
import org.eclipse.jetty.server.handler.HandlerCollection;
import org.eclipse.jetty.webapp.WebAppContext;

import javax.management.MBeanServer;
import java.lang.management.ManagementFactory;

public class TestServer
{
    public static void main(String[] args) throws Exception
    {
        System.setProperty("deploy.dir", "core/src/test");
        System.setProperty("dns.zones.dir", "core/src/test/var/auto-zones");

        LogManager.getRootLogger().addAppender(new ConsoleAppender(new PatternLayout("%-5p %d{ISO8601} %c: %m%n")));
        LogManager.getRootLogger().setLevel(Level.INFO);
        LogManager.getLogger("org.springframework").setLevel(Level.WARN);

        Server server = new Server();
        SocketConnector connector = new SocketConnector();

        // Set some timeout options to make debugging easier.
        connector.setMaxIdleTime(3600 * 1000);
        connector.setSoLingerTime(-1);
        connector.setPort(8888);
        server.addConnector(connector);

        HandlerCollection handlers = new HandlerCollection();

        WebAppContext trafficRouterContext = new WebAppContext("core/src/main/webapp", "/trafficrouter");
        handlers.addHandler(trafficRouterContext);

        WebAppContext apiContext = new WebAppContext();
        apiContext.setWar("api/src/main/webapp");
        handlers.addHandler(apiContext);

        handlers.setServer(server);
        server.setHandler(handlers);

        MBeanServer mBeanServer = ManagementFactory.getPlatformMBeanServer();
        MBeanContainer mBeanContainer = new MBeanContainer(mBeanServer);
        server.getContainer().addEventListener(mBeanContainer);
        mBeanContainer.start();

        System.out.println(">>>>>>> Starting everything press 'Q' and enter to stop");

        server.start();

        while (System.in.read() != 'Q');

        System.out.println(">>>>>>>>>> Stopping Jetty");

        server.stop();
        server.join();
    }
}
