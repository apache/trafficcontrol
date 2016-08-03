package com.comcast.cdn.traffic_control.traffic_router.protocol;

import org.apache.tomcat.util.net.ServerSocketFactory;
import org.apache.tomcat.util.net.jsse.JSSEFactory;

public class RouterJsseFactory extends JSSEFactory {
	final ServerSocketFactory serverSocketFactory = new RouterSslServerSocketFactory();
	@Override
	public ServerSocketFactory getSocketFactory() {
		return serverSocketFactory;
	}
}
